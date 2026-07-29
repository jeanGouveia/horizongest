package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type CompanySettingsHandler struct {
	companySettingsService *service.CompanySettingsService
}

func NewCompanySettingsHandler(companySettingsService *service.CompanySettingsService) *CompanySettingsHandler {
	return &CompanySettingsHandler{companySettingsService: companySettingsService}
}

// GetSettings handles GET /api/company/settings
func (h *CompanySettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	log.Printf("[FORENSIC] CompanySettingsHandler - GetSettings - Request recebida")

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("[FORENSIC] CompanySettingsHandler - UserID não encontrado no contexto")
		jsonError(w, "não autorizado. Faça login novamente.", http.StatusUnauthorized)
		return
	}

	// FORENSIC: Log UserID received
	log.Printf("[FORENSIC] CompanySettingsHandler - UserID recebido: %d", userID)

	// FORENSIC: Log TenantContext
	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if ok {
		log.Printf("[FORENSIC] CompanySettingsHandler - TenantContext - UserID: %d, CompanyID: %d", tenantCtx.UserID, tenantCtx.CompanyID)
	} else {
		log.Printf("[FORENSIC] CompanySettingsHandler - TenantContext NÃO encontrado no contexto")
	}

	// FORENSIC: Log Claims
	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if ok {
		log.Printf("[FORENSIC] CompanySettingsHandler - Claims - UserID: %d, CompanyID: %d, Email: %s, Name: %s", claims.UserID, claims.CompanyID, claims.Email, claims.Name)
	}

	settings, err := h.companySettingsService.GetSettings(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNoCompany) {
			jsonError(w, "usuário não possui uma empresa associada", http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível carregar as configurações da empresa", http.StatusInternalServerError)
		return
	}

	// FORENSIC: Log company settings loaded
	log.Printf("[FORENSIC] CompanySettingsHandler - Empresa carregada: Name=%s, Slug=%s", settings.Name, settings.Slug)

	jsonResponse(w, http.StatusOK, settings)
}

// UpdateSettings handles PUT /api/company/settings
func (h *CompanySettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado. Faça login novamente.", http.StatusUnauthorized)
		return
	}

	var input service.UpdateSettingsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Skip validation for now since all fields are optional
	// Validation will be handled by the service layer if needed

	if err := h.companySettingsService.UpdateSettings(r.Context(), userID, input); err != nil {
		if errors.Is(err, service.ErrUserNoCompany) {
			jsonError(w, "usuário não possui uma empresa associada", http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível atualizar as configurações da empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "configurações atualizadas com sucesso"})
}
