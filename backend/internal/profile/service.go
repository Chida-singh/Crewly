package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"crewly/backend/internal/auth"
	"crewly/backend/internal/platform/ids"
	"crewly/backend/internal/platform/storage"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrForbidden    = errors.New("forbidden")
	ErrPhotoLimit   = errors.New("photo limit reached")
)

const MaxPhotos = 12

type Service struct {
	repo    Repo
	storage storage.Presigner
}

func NewService(repo Repo, presigner storage.Presigner) *Service {
	return &Service{repo: repo, storage: presigner}
}

// Snapshot is the role-agnostic view of the caller's profile.
type Snapshot struct {
	Role   string
	Talent *TalentProfile
	Client *ClientProfile
	Photos []Photo
}

func (s *Service) Get(ctx context.Context, userID, role string) (Snapshot, error) {
	switch role {
	case auth.RoleTalent:
		talent, err := s.repo.GetTalent(ctx, userID)
		if err != nil {
			return Snapshot{}, err
		}
		photos, err := s.repo.ListPhotos(ctx, userID)
		if err != nil {
			return Snapshot{}, err
		}
		return Snapshot{Role: role, Talent: &talent, Photos: photos}, nil
	case auth.RoleClient:
		client, err := s.repo.GetClient(ctx, userID)
		if err != nil {
			return Snapshot{}, err
		}
		return Snapshot{Role: role, Client: &client}, nil
	default:
		return Snapshot{}, fmt.Errorf("%w: unknown role", ErrInvalidInput)
	}
}

type TalentInput struct {
	TalentType  string `json:"talent_type"`
	DisplayName string `json:"display_name"`
	City        string `json:"city"`
	Age         int32  `json:"age"`
	HeightCm    *int32 `json:"height_cm"`
	Bio         string `json:"bio"`
	RateMin     int32  `json:"rate_min"`
	RateMax     int32  `json:"rate_max"`
}

func (s *Service) UpsertTalent(ctx context.Context, userID string, in TalentInput) (Snapshot, error) {
	in.TalentType = strings.ToLower(strings.TrimSpace(in.TalentType))
	if in.TalentType == "" {
		in.TalentType = "model"
	}
	in.DisplayName = strings.TrimSpace(in.DisplayName)
	in.City = strings.TrimSpace(in.City)
	in.Bio = strings.TrimSpace(in.Bio)

	switch {
	case in.DisplayName == "":
		return Snapshot{}, fmt.Errorf("%w: display_name is required", ErrInvalidInput)
	case len(in.DisplayName) > 80:
		return Snapshot{}, fmt.Errorf("%w: display_name is too long", ErrInvalidInput)
	case in.City == "":
		return Snapshot{}, fmt.Errorf("%w: city is required", ErrInvalidInput)
	case in.Age < 18:
		return Snapshot{}, fmt.Errorf("%w: talent must be 18 or older", ErrInvalidInput)
	case in.Age > 100:
		return Snapshot{}, fmt.Errorf("%w: age is out of range", ErrInvalidInput)
	case in.HeightCm != nil && (*in.HeightCm < 100 || *in.HeightCm > 250):
		return Snapshot{}, fmt.Errorf("%w: height_cm is out of range", ErrInvalidInput)
	case in.RateMin < 0 || in.RateMax < 0:
		return Snapshot{}, fmt.Errorf("%w: rates cannot be negative", ErrInvalidInput)
	case in.RateMax < in.RateMin:
		return Snapshot{}, fmt.Errorf("%w: rate_max must be at least rate_min", ErrInvalidInput)
	case len(in.Bio) > 500:
		return Snapshot{}, fmt.Errorf("%w: bio is too long", ErrInvalidInput)
	}

	saved, err := s.repo.UpsertTalent(ctx, TalentProfile{
		UserID:      userID,
		TalentType:  in.TalentType,
		DisplayName: in.DisplayName,
		City:        in.City,
		Age:         in.Age,
		HeightCm:    in.HeightCm,
		Bio:         in.Bio,
		RateMin:     in.RateMin,
		RateMax:     in.RateMax,
	})
	if err != nil {
		return Snapshot{}, err
	}

	photos, err := s.repo.ListPhotos(ctx, userID)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Role: auth.RoleTalent, Talent: &saved, Photos: photos}, nil
}

type ClientInput struct {
	DisplayName string `json:"display_name"`
	OrgName     string `json:"org_name"`
}

func (s *Service) UpsertClient(ctx context.Context, userID string, in ClientInput) (Snapshot, error) {
	in.DisplayName = strings.TrimSpace(in.DisplayName)
	in.OrgName = strings.TrimSpace(in.OrgName)

	switch {
	case in.DisplayName == "":
		return Snapshot{}, fmt.Errorf("%w: display_name is required", ErrInvalidInput)
	case len(in.DisplayName) > 80:
		return Snapshot{}, fmt.Errorf("%w: display_name is too long", ErrInvalidInput)
	case len(in.OrgName) > 120:
		return Snapshot{}, fmt.Errorf("%w: org_name is too long", ErrInvalidInput)
	}

	saved, err := s.repo.UpsertClient(ctx, ClientProfile{
		UserID:      userID,
		DisplayName: in.DisplayName,
		OrgName:     in.OrgName,
	})
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Role: auth.RoleClient, Client: &saved}, nil
}

// AddPhoto reserves a photo slot and returns a presigned upload URL.
func (s *Service) AddPhoto(ctx context.Context, userID, role, contentType string) (Photo, string, error) {
	if role != auth.RoleTalent {
		return Photo{}, "", ErrForbidden
	}

	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if contentType == "" {
		contentType = "image/jpeg"
	}
	if !strings.HasPrefix(contentType, "image/") {
		return Photo{}, "", fmt.Errorf("%w: content_type must be an image", ErrInvalidInput)
	}

	count, err := s.repo.CountPhotos(ctx, userID)
	if err != nil {
		return Photo{}, "", err
	}
	if count >= MaxPhotos {
		return Photo{}, "", ErrPhotoLimit
	}

	key := "profiles/" + userID + "/" + ids.New() + photoExt(contentType)
	uploadURL, publicURL, err := s.storage.PresignPut(ctx, key, contentType)
	if err != nil {
		return Photo{}, "", fmt.Errorf("presign upload: %w", err)
	}

	photo, err := s.repo.CreatePhoto(ctx, userID, publicURL, int32(count))
	if err != nil {
		return Photo{}, "", err
	}
	return photo, uploadURL, nil
}

func photoExt(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}
