package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crewly/backend/internal/platform/httpx"
)

func TestJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := map[string]string{"message": "hello"}

	httpx.JSON(rec, http.StatusOK, payload)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json header, got %s", rec.Header().Get("Content-Type"))
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res["message"] != "hello" {
		t.Errorf("expected hello, got %s", res["message"])
	}
}

func TestError(t *testing.T) {
	rec := httptest.NewRecorder()

	httpx.Error(rec, http.StatusBadRequest, "invalid payload")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	var res httpx.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res.Error != "invalid payload" {
		t.Errorf("expected 'invalid payload', got %s", res.Error)
	}
}

