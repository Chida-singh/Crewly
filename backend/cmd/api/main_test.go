package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"crewly/backend/internal/platform/httpx"
)

func TestHealthEndpointDegradedWhenNoDB(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":   "degraded",
			"database": "down",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json response: %v", err)
	}

	if res["status"] != "degraded" || res["database"] != "down" {
		t.Errorf("expected degraded/down, got %v", res)
	}
}

