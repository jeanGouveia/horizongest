package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type ThemeHandler struct {
	themeService *service.ThemeService
}

func NewThemeHandler(themeService *service.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}

// GetTheme returns the theme configuration for the authenticated user
// GET /api/theme
func (h *ThemeHandler) GetTheme(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	theme, err := h.themeService.GetThemeForUser(r.Context(), userID)
	if err != nil {
		jsonError(w, "failed to load theme", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, theme)
}

// GetDefaultTheme returns the default PratoOnline theme
// GET /api/theme/default
func (h *ThemeHandler) GetDefaultTheme(w http.ResponseWriter, r *http.Request) {
	theme := h.themeService.GetDefaultTheme()
	jsonResponse(w, http.StatusOK, theme)
}

func (h *ThemeHandler) RegisterRoutes(r chi.Router) {
	r.Get("/theme", h.GetTheme)
	r.Get("/theme/default", h.GetDefaultTheme)
}
