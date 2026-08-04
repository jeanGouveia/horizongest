package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/config"
)

func TestHealthHandler_LivenessCheck(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
	}
	handler := NewHealthHandler(cfg)

	req := httptest.NewRequest("GET", "/health/live", nil)
	rr := httptest.NewRecorder()

	handler.LivenessCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("expected status ok, got %s", response.Status)
	}

	if response.Service != "horizongest" {
		t.Errorf("expected service horizongest, got %s", response.Service)
	}

	if response.Timestamp == "" {
		t.Error("expected timestamp to be set")
	}
}

func TestHealthHandler_ReadinessCheck(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
	}
	handler := NewHealthHandler(cfg)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	rr := httptest.NewRecorder()

	handler.ReadinessCheck(rr, req)

	// Should return 503 when dependencies are not set
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rr.Code)
	}

	var response HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "not_ready" {
		t.Errorf("expected status not_ready, got %s", response.Status)
	}

	if response.Service != "horizongest" {
		t.Errorf("expected service horizongest, got %s", response.Service)
	}
}

func TestHealthHandler_ReadinessCheck_WithDependencies(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
	}
	handler := NewHealthHandler(cfg)

	// Note: In a real test, you would set up mock dependencies
	// For now, we test that the handler can be created with dependencies
	handler.SetDependencies(nil, nil)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	rr := httptest.NewRecorder()

	handler.ReadinessCheck(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rr.Code)
	}
}
