package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type PlatformCompanyHandler struct {
	platformService *service.PlatformService
}

func NewPlatformCompanyHandler(platformService *service.PlatformService) *PlatformCompanyHandler {
	return &PlatformCompanyHandler{platformService: platformService}
}

// POST /api/platform/companies
func (h *PlatformCompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CompanyData service.PlatformCreateCompanyInput `json:"company"`
		OwnerEmail  string                             `json:"owner_email" validate:"required,email"`
		OwnerName   string                             `json:"owner_name" validate:"required,min=2"`
		Password    string                             `json:"password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Password hashing is handled by the service
	output, err := h.platformService.CreateCompany(r.Context(), platformUserID, input.CompanyData, input.OwnerEmail, input.Password, input.OwnerName)
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyAlreadyExists) {
			jsonError(w, "empresa com este slug já existe", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível criar empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, output)
}

// GET /api/platform/companies
func (h *PlatformCompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	companies, err := h.platformService.ListCompanies(r.Context(), platformUserID)
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível listar empresas", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"companies": companies,
	})
}

// GET /api/platform/companies/:id
func (h *PlatformCompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	company, err := h.platformService.GetCompany(r.Context(), platformUserID, uint(companyID))
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível obter empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, company)
}

// PUT /api/platform/companies/:id
func (h *PlatformCompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var input service.PlatformCreateCompanyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.UpdateCompany(r.Context(), platformUserID, uint(companyID), input); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrCompanyAlreadyExists) {
			jsonError(w, "empresa com este slug já existe", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível atualizar empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "empresa atualizada com sucesso",
	})
}

// POST /api/platform/companies/:id/deactivate
func (h *PlatformCompanyHandler) DeactivateCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.DeactivateCompany(r.Context(), platformUserID, uint(companyID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível desativar empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "empresa desativada com sucesso",
	})
}

// POST /api/platform/companies/:id/activate
func (h *PlatformCompanyHandler) ActivateCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.ActivateCompany(r.Context(), platformUserID, uint(companyID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível ativar empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "empresa ativada com sucesso",
	})
}

// GET /api/platform/companies/:id/owner
func (h *PlatformCompanyHandler) GetCompanyOwner(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	owner, err := h.platformService.GetCompanyOwner(r.Context(), platformUserID, uint(companyID))
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível obter owner", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, owner)
}

// POST /api/platform/companies/:id/owner/reset-password
func (h *PlatformCompanyHandler) ResetOwnerPassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.ResetOwnerPassword(r.Context(), platformUserID, uint(companyID), input.Password); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível redefinir senha", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "senha redefinida com sucesso",
	})
}

// POST /api/platform/users/:id/block
func (h *PlatformCompanyHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get user ID from URL parameter
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.BlockUser(r.Context(), platformUserID, uint(userID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível bloquear usuário", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "usuário bloqueado com sucesso",
	})
}

// POST /api/platform/users/:id/unblock
func (h *PlatformCompanyHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get user ID from URL parameter
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.UnblockUser(r.Context(), platformUserID, uint(userID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível desbloquear usuário", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "usuário desbloqueado com sucesso",
	})
}

// POST /api/platform/companies/:id/login-as
func (h *PlatformCompanyHandler) LoginAsCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	ownerEmail, err := h.platformService.LoginAsCompany(r.Context(), platformUserID, uint(companyID))
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível fazer login como empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message":     "login como empresa iniciado",
		"owner_email": ownerEmail,
	})
}

// POST /api/platform/companies/:id/trial
func (h *PlatformCompanyHandler) SetCompanyTrial(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TrialEndsAt string `json:"trial_ends_at" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	// Parse trial end date
	trialEndsAt, err := time.Parse(time.RFC3339, input.TrialEndsAt)
	if err != nil {
		jsonError(w, "data de fim de trial inválida", http.StatusBadRequest)
		return
	}

	if err := h.platformService.SetCompanyTrial(r.Context(), platformUserID, uint(companyID), trialEndsAt); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível definir trial", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "trial definido com sucesso",
	})
}

// POST /api/platform/companies/:id/suspend
func (h *PlatformCompanyHandler) SuspendCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.SuspendCompany(r.Context(), platformUserID, uint(companyID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível suspender empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "empresa suspensa com sucesso",
	})
}

// POST /api/platform/companies/:id/cancel
func (h *PlatformCompanyHandler) CancelCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.CancelCompany(r.Context(), platformUserID, uint(companyID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível cancelar empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "empresa cancelada com sucesso",
	})
}

// POST /api/platform/companies/:id/reactivate
func (h *PlatformCompanyHandler) ReactivateCompany(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get company ID from URL parameter
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	if err := h.platformService.ReactivateCompany(r.Context(), platformUserID, uint(companyID)); err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			jsonError(w, "permissão negada", http.StatusForbidden)
			return
		}
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível reativar empresa", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "empresa reativada com sucesso",
	})
}
