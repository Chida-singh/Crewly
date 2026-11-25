package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"crewly/backend/internal/platform/db"
	"crewly/backend/internal/platform/ids"
)

var ErrNotFound = errors.New("profile not found")

type TalentProfile struct {
	UserID      string
	TalentType  string
	DisplayName string
	City        string
	Age         int32
	HeightCm    *int32
	Bio         string
	RateMin     int32
	RateMax     int32
}

type ClientProfile struct {
	UserID      string
	DisplayName string
	OrgName     string
}

type Photo struct {
	ID       string
	URL      string
	Position int32
}

type Repo interface {
	GetTalent(ctx context.Context, userID string) (TalentProfile, error)
	UpsertTalent(ctx context.Context, p TalentProfile) (TalentProfile, error)
	GetClient(ctx context.Context, userID string) (ClientProfile, error)
	UpsertClient(ctx context.Context, p ClientProfile) (ClientProfile, error)
	CountPhotos(ctx context.Context, talentID string) (int64, error)
	CreatePhoto(ctx context.Context, talentID, url string, position int32) (Photo, error)
	ListPhotos(ctx context.Context, talentID string) ([]Photo, error)
}

type PGRepo struct {
	q *db.Queries
}

func NewPGRepo(pool *pgxpool.Pool) *PGRepo {
	return &PGRepo{q: db.New(pool)}
}

func (p *PGRepo) GetTalent(ctx context.Context, userID string) (TalentProfile, error) {
	uid, err := ids.ToPG(userID)
	if err != nil {
		return TalentProfile{}, err
	}
	row, err := p.q.GetTalentProfile(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		return TalentProfile{}, ErrNotFound
	}
	if err != nil {
		return TalentProfile{}, fmt.Errorf("get talent profile: %w", err)
	}
	return fromTalent(row), nil
}

func (p *PGRepo) UpsertTalent(ctx context.Context, t TalentProfile) (TalentProfile, error) {
	uid, err := ids.ToPG(t.UserID)
	if err != nil {
		return TalentProfile{}, err
	}
	height := pgtype.Int4{}
	if t.HeightCm != nil {
		height = pgtype.Int4{Int32: *t.HeightCm, Valid: true}
	}
	row, err := p.q.UpsertTalentProfile(ctx, db.UpsertTalentProfileParams{
		UserID:      uid,
		TalentType:  t.TalentType,
		DisplayName: t.DisplayName,
		City:        t.City,
		Age:         t.Age,
		HeightCm:    height,
		Bio:         t.Bio,
		RateMin:     t.RateMin,
		RateMax:     t.RateMax,
	})
	if err != nil {
		return TalentProfile{}, fmt.Errorf("upsert talent profile: %w", err)
	}
	return fromTalent(row), nil
}

func (p *PGRepo) GetClient(ctx context.Context, userID string) (ClientProfile, error) {
	uid, err := ids.ToPG(userID)
	if err != nil {
		return ClientProfile{}, err
	}
	row, err := p.q.GetClientProfile(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		return ClientProfile{}, ErrNotFound
	}
	if err != nil {
		return ClientProfile{}, fmt.Errorf("get client profile: %w", err)
	}
	return fromClient(row), nil
}

func (p *PGRepo) UpsertClient(ctx context.Context, c ClientProfile) (ClientProfile, error) {
	uid, err := ids.ToPG(c.UserID)
	if err != nil {
		return ClientProfile{}, err
	}
	row, err := p.q.UpsertClientProfile(ctx, db.UpsertClientProfileParams{
		UserID:      uid,
		DisplayName: c.DisplayName,
		OrgName:     c.OrgName,
	})
	if err != nil {
		return ClientProfile{}, fmt.Errorf("upsert client profile: %w", err)
	}
	return fromClient(row), nil
}

func (p *PGRepo) CountPhotos(ctx context.Context, talentID string) (int64, error) {
	uid, err := ids.ToPG(talentID)
	if err != nil {
		return 0, err
	}
	return p.q.CountProfilePhotos(ctx, uid)
}

func (p *PGRepo) CreatePhoto(ctx context.Context, talentID, url string, position int32) (Photo, error) {
	uid, err := ids.ToPG(talentID)
	if err != nil {
		return Photo{}, err
	}
	row, err := p.q.CreateProfilePhoto(ctx, db.CreateProfilePhotoParams{
		TalentID: uid,
		Url:      url,
		Position: position,
	})
	if err != nil {
		return Photo{}, fmt.Errorf("create profile photo: %w", err)
	}
	return Photo{ID: ids.FromPG(row.ID), URL: row.Url, Position: row.Position}, nil
}

func (p *PGRepo) ListPhotos(ctx context.Context, talentID string) ([]Photo, error) {
	uid, err := ids.ToPG(talentID)
	if err != nil {
		return nil, err
	}
	rows, err := p.q.ListProfilePhotos(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("list profile photos: %w", err)
	}
	photos := make([]Photo, 0, len(rows))
	for _, row := range rows {
		photos = append(photos, Photo{ID: ids.FromPG(row.ID), URL: row.Url, Position: row.Position})
	}
	return photos, nil
}

func fromTalent(t db.TalentProfile) TalentProfile {
	var height *int32
	if t.HeightCm.Valid {
		h := t.HeightCm.Int32
		height = &h
	}
	return TalentProfile{
		UserID:      ids.FromPG(t.UserID),
		TalentType:  t.TalentType,
		DisplayName: t.DisplayName,
		City:        t.City,
		Age:         t.Age,
		HeightCm:    height,
		Bio:         t.Bio,
		RateMin:     t.RateMin,
		RateMax:     t.RateMax,
	}
}

func fromClient(c db.ClientProfile) ClientProfile {
	return ClientProfile{
		UserID:      ids.FromPG(c.UserID),
		DisplayName: c.DisplayName,
		OrgName:     c.OrgName,
	}
}
