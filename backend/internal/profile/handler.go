package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	"crewly/backend/internal/auth"
	"crewly/backend/internal/platform/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type photoResponse struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Position int32  `json:"position"`
}

type talentResponse struct {
	Role        string          `json:"role"`
	TalentType  string          `json:"talent_type"`
	DisplayName string          `json:"display_name"`
	City        string          `json:"city"`
	Age         int32           `json:"age"`
	HeightCm    *int32          `json:"height_cm"`
	Bio         string          `json:"bio"`
	RateMin     int32           `json:"rate_min"`
	RateMax     int32           `json:"rate_max"`
	Photos      []photoResponse `json:"photos"`
}

type clientResponse struct {
	Role        string `json:"role"`
	DisplayName string `json:"display_name"`
	OrgName     string `json:"org_name"`
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.FromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	snapshot, err := h.svc.Get(r.Context(), user.ID, user.Role)
	if err != nil {
		respondError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, snapshotJSON(snapshot))
}

func (h *Handler) PutProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.FromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var (
		snapshot Snapshot
		err      error
	)
	switch user.Role {
	case auth.RoleTalent:
		var in TalentInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid json body")
			return
		}
		snapshot, err = h.svc.UpsertTalent(r.Context(), user.ID, in)
	case auth.RoleClient:
		var in ClientInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid json body")
			return
		}
		snapshot, err = h.svc.UpsertClient(r.Context(), user.ID, in)
	default:
		httpx.Error(w, http.StatusForbidden, "forbidden for this role")
		return
	}

	if err != nil {
		respondError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, snapshotJSON(snapshot))
}

type photoRequest struct {
	ContentType string `json:"content_type"`
}

func (h *Handler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.FromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body photoRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid json body")
			return
		}
	}

	photo, uploadURL, err := h.svc.AddPhoto(r.Context(), user.ID, user.Role, body.ContentType)
	if err != nil {
		respondError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]any{
		"photo":      photoResponse{ID: photo.ID, URL: photo.URL, Position: photo.Position},
		"upload_url": uploadURL,
	})
}

func snapshotJSON(s Snapshot) any {
	if s.Talent != nil {
		photos := make([]photoResponse, 0, len(s.Photos))
		for _, p := range s.Photos {
			photos = append(photos, photoResponse{ID: p.ID, URL: p.URL, Position: p.Position})
		}
		return talentResponse{
			Role:        s.Role,
			TalentType:  s.Talent.TalentType,
			DisplayName: s.Talent.DisplayName,
			City:        s.Talent.City,
			Age:         s.Talent.Age,
			HeightCm:    s.Talent.HeightCm,
			Bio:         s.Talent.Bio,
			RateMin:     s.Talent.RateMin,
			RateMax:     s.Talent.RateMax,
			Photos:      photos,
		}
	}
	if s.Client != nil {
		return clientResponse{Role: s.Role, DisplayName: s.Client.DisplayName, OrgName: s.Client.OrgName}
	}
	return map[string]string{"role": s.Role}
}

func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrForbidden):
		httpx.Error(w, http.StatusForbidden, "only talent can add photos")
	case errors.Is(err, ErrPhotoLimit):
		httpx.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "profile not found")
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal error")
	}
}
