package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"crewly/backend/internal/platform/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/otp", h.requestOTP)
	r.Post("/verify", h.verifyOTP)
	return r
}

type otpRequest struct {
	Phone string `json:"phone"`
}

func (h *Handler) requestOTP(w http.ResponseWriter, r *http.Request) {
	var body otpRequest
	if err := decode(w, r, &body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.svc.RequestOTP(r.Context(), body.Phone); err != nil {
		respondError(w, err)
		return
	}
	httpx.JSON(w, http.StatusAccepted, map[string]string{"message": "otp sent"})
}

type verifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
	Role  string `json:"role"`
}

type userResponse struct {
	ID         string `json:"id"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

func (h *Handler) verifyOTP(w http.ResponseWriter, r *http.Request) {
	var body verifyRequest
	if err := decode(w, r, &body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	user, token, err := h.svc.VerifyOTP(r.Context(), body.Phone, body.Code, body.Role)
	if err != nil {
		respondError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  userResponse{ID: user.ID, Phone: user.Phone, Role: user.Role, IsVerified: user.IsVerified},
	})
}

func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidPhone), errors.Is(err, ErrInvalidCode), errors.Is(err, ErrInvalidRole):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrRateLimited):
		httpx.Error(w, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, ErrUnauthorized):
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal error")
	}
}

func decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	return json.NewDecoder(r.Body).Decode(dst)
}
