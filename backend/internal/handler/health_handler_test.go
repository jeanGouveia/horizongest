package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_LivenessCheck(t *testing.T) {
	handler := NewHealthHandler()

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
	handler := NewHealthHandler()

	req := httptest.NewRequest("GET", "/health/ready", nil)
	rr := httptest.NewRecorder()

	handler.ReadinessCheck(rr, req)

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
}
