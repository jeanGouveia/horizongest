package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type SystemHandler struct {
	startTime time.Time
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		startTime: time.Now(),
	}
}

func (h *SystemHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.GetHealth)
	r.Get("/version", h.GetVersion)
	r.Get("/capabilities", h.GetCapabilities)
}

func (h *SystemHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	health := &domain.Health{
		Status:   "healthy",
		Database: "connected",
		Storage:  "available",
		Version:  "1.0.0",
		Uptime:   time.Since(h.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

func (h *SystemHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	version := &domain.Version{
		Version:     "1.0.0",
		Commit:      "dev",
		Build:       time.Now().Format("20060102-150405"),
		Environment: "development",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}

func (h *SystemHandler) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	capabilities := &domain.Capabilities{
		Upload:          true,
		SEO:             true,
		Marketplace:     false,
		IFood:           false,
		PIX:             false,
		Fiscal:          false,
		Delivery:        false,
		CardapioDigital: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(capabilities)
}
