package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type BusinessHandler struct {
	businessService *service.BusinessService
}

func NewBusinessHandler(businessService *service.BusinessService) *BusinessHandler {
	return &BusinessHandler{businessService: businessService}
}

// GetBusinessProfile returns the business profile for the authenticated user
// GET /api/business/profile
func (h *BusinessHandler) GetBusinessProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	profile, err := h.businessService.GetBusinessProfile(r.Context(), userID)
	if err != nil {
		jsonError(w, "failed to load business profile", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, profile)
}

// GetDefaultBusinessProfile returns the default tenant business profile
// GET /api/business/profile/default
func (h *BusinessHandler) GetDefaultBusinessProfile(w http.ResponseWriter, r *http.Request) {
	profile := h.businessService.GetDefaultBusinessProfile()
	jsonResponse(w, http.StatusOK, profile)
}

func (h *BusinessHandler) RegisterRoutes(r chi.Router) {
	r.Get("/business/profile", h.GetBusinessProfile)
	r.Get("/business/profile/default", h.GetDefaultBusinessProfile)
}
