package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Service   string            `json:"service"`
	Version   string            `json:"version,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// LivenessCheck returns whether the service is alive
// This is a simple check that always returns OK if the service is running
func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "horizongest",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessCheck returns whether the service is ready to handle requests
// This checks dependencies like database, redis, etc.
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	// FASE A: Basic readiness check - can be extended to check DB, Redis, etc.
	// For now, we'll return OK if the service is running
	// In a full implementation, you would check:
	// - Database connectivity
	// - Redis connectivity
	// - RabbitMQ connectivity
	// - Other external dependencies

	checks := make(map[string]string)
	allChecksOK := true

	// Add correlation ID to response if available
	correlationID := middleware.GetCorrelationID(r.Context())
	requestID := middleware.GetRequestID(r.Context())

	if correlationID != "" {
		checks["correlation_id"] = correlationID
	}
	if requestID != "" {
		checks["request_id"] = requestID
	}

	status := "ok"
	if !allChecksOK {
		status = "not_ready"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "horizongest",
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health/live", h.LivenessCheck)
	r.Get("/health/ready", h.ReadinessCheck)
	r.Get("/health", h.LivenessCheck) // Alias for liveness
}
