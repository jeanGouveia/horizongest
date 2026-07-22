package handler

import (
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type PlatformDashboardHandler struct {
	platformService *service.PlatformService
}

func NewPlatformDashboardHandler(platformService *service.PlatformService) *PlatformDashboardHandler {
	return &PlatformDashboardHandler{platformService: platformService}
}

// GET /api/platform/dashboard/stats
func (h *PlatformDashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.platformService.GetDashboardStats(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar estatísticas", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, stats)
}
