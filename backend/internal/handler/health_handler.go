package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/config"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/redis/go-redis/v9"
)

// HealthHandler handles health check endpoints
// FASE A.4: Comprehensive health checks for all dependencies
type HealthHandler struct {
	cfg   *config.Config
	db    *pgx.Conn
	redis *redis.Client
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		cfg: cfg,
	}
}

// SetDependencies sets the runtime dependencies for health checks
// FASE A.4: Set dependencies for comprehensive health checks
func (h *HealthHandler) SetDependencies(db *pgx.Conn, redis *redis.Client) {
	h.db = db
	h.redis = redis
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version,omitempty"`
	Checks    map[string]HealthCheck `json:"checks,omitempty"`
}

// HealthCheck represents a single health check result
type HealthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
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
// FASE A.4: Comprehensive readiness check for all dependencies
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]HealthCheck)
	allChecksOK := true

	// Check Database
	dbCheck := h.checkDatabase(ctx)
	checks["database"] = dbCheck
	if dbCheck.Status != "ok" {
		allChecksOK = false
	}

	// Check Redis
	redisCheck := h.checkRedis(ctx)
	checks["redis"] = redisCheck
	if redisCheck.Status != "ok" {
		allChecksOK = false
	}

	// Add correlation ID to response if available
	correlationID := middleware.GetCorrelationID(r.Context())
	requestID := middleware.GetRequestID(r.Context())

	if correlationID != "" {
		checks["correlation_id"] = HealthCheck{Status: "ok", Message: correlationID}
	}
	if requestID != "" {
		checks["request_id"] = HealthCheck{Status: "ok", Message: requestID}
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

// checkDatabase checks database connectivity
// FASE A.4: Database health check
func (h *HealthHandler) checkDatabase(ctx context.Context) HealthCheck {
	if h.db == nil {
		return HealthCheck{
			Status:  "not_ready",
			Message: "database connection not initialized",
		}
	}

	var result int
	err := h.db.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return HealthCheck{
			Status:  "not_ready",
			Message: err.Error(),
		}
	}

	if result != 1 {
		return HealthCheck{
			Status:  "not_ready",
			Message: "database query returned unexpected result",
		}
	}

	return HealthCheck{
		Status:  "ok",
		Message: "database connection healthy",
	}
}

// checkRedis checks Redis connectivity
// FASE A.4: Redis health check
func (h *HealthHandler) checkRedis(ctx context.Context) HealthCheck {
	if h.redis == nil {
		return HealthCheck{
			Status:  "not_ready",
			Message: "redis connection not initialized",
		}
	}

	result, err := h.redis.Ping(ctx).Result()
	if err != nil {
		return HealthCheck{
			Status:  "not_ready",
			Message: err.Error(),
		}
	}

	if result != "PONG" {
		return HealthCheck{
			Status:  "not_ready",
			Message: "redis ping returned unexpected result",
		}
	}

	return HealthCheck{
		Status:  "ok",
		Message: "redis connection healthy",
	}
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health/live", h.LivenessCheck)
	r.Get("/health/ready", h.ReadinessCheck)
	r.Get("/health", h.LivenessCheck) // Alias for liveness
}
