package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

type DashboardHandler struct {
	dashboardRepo ports.DashboardRepository
}

func NewDashboardHandler(dashboardRepo ports.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{dashboardRepo: dashboardRepo}
}

func (h *DashboardHandler) RegisterRoutes(r chi.Router) {
	r.Get("/dashboard", h.GetDashboard)
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard, err := h.dashboardRepo.GetDashboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboard)
}
