package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type CompanyHandler struct {
	svc *service.CompanyService
}

func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// POST /api/companies
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var in service.CreateCompanyInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	c, err := h.svc.CreateCompany(r.Context(), in)
	if err != nil {
		if errors.Is(err, service.ErrSlugAlreadyExists) {
			jsonError(w, "slug já está em uso. Escolha outro identificador.", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível criar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, c)
}

// GET /api/companies
func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.svc.ListCompanies(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar as empresas. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, companies)
}

// GET /api/companies/{id}
// DEPRECATED: Use GET /api/me/company instead
// This endpoint now validates that the requested ID matches the tenant's CompanyID
// to prevent IDOR (Insecure Direct Object Reference) attacks
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da empresa inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}

	// Security check: Validate that requested ID matches tenant's CompanyID
	tenantCtx, ok := middleware.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}

	if tenantCtx.CompanyID == 0 {
		jsonError(w, "company ID não encontrado no contexto tenant", http.StatusUnauthorized)
		return
	}

	// Prevent IDOR: User can only access their own company
	if id != tenantCtx.CompanyID {
		jsonError(w, "acesso negado: empresa não pertence ao usuário", http.StatusForbidden)
		return
	}

	c, err := h.svc.GetCompany(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// GET /api/me/company
// Secure endpoint that returns the current user's company from tenant context
// Does not accept CompanyID from URL - uses only the CompanyID from JWT/tenant context
func (h *CompanyHandler) GetCurrentCompany(w http.ResponseWriter, r *http.Request) {
	// Get tenant context from middleware
	tenantCtx, ok := middleware.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}

	if tenantCtx.CompanyID == 0 {
		jsonError(w, "company ID não encontrado no contexto tenant", http.StatusUnauthorized)
		return
	}

	c, err := h.svc.GetCurrentCompany(r.Context(), tenantCtx.CompanyID)
	if err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// PUT /api/companies/{id}
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da empresa inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateCompanyInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	c, err := h.svc.UpdateCompany(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrSlugAlreadyExists) {
			jsonError(w, "slug já está em uso. Escolha outro identificador.", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível atualizar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// DELETE /api/companies/{id}
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da empresa inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteCompany(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível remover a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "empresa removida com sucesso"})
}
