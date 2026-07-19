package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type UserManagementHandler struct {
	userManagementService *service.UserManagementService
	userRepo              ports.UserRepository
}

func NewUserManagementHandler(userManagementService *service.UserManagementService, userRepo ports.UserRepository) *UserManagementHandler {
	return &UserManagementHandler{
		userManagementService: userManagementService,
		userRepo:              userRepo,
	}
}

// getUserCompanyID gets the user's CompanyID from the repository
func (h *UserManagementHandler) getUserCompanyID(ctx context.Context, userID uint) (uint, error) {
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if user == nil || user.CompanyID == nil {
		return 0, errors.New("usuário não possui empresa")
	}
	return *user.CompanyID, nil
}

// ListUsers handles GET /api/company/users
func (h *UserManagementHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	companyID, err := h.getUserCompanyID(r.Context(), userID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusForbidden)
		return
	}

	users, err := h.userManagementService.ListUsers(r.Context(), companyID)
	if err != nil {
		jsonError(w, "não foi possível listar usuários", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, users)
}

// GetUser handles GET /api/company/users/{id}
func (h *UserManagementHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	companyID, err := h.getUserCompanyID(r.Context(), userID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusForbidden)
		return
	}

	targetIDStr := chi.URLParam(r, "id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := h.userManagementService.GetUser(r.Context(), companyID, uint(targetID))
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrUserNotInCompany {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar usuário", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, user)
}

// AddUser handles POST /api/company/users/add
func (h *UserManagementHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		jsonError(w, "e-mail é obrigatório", http.StatusBadRequest)
		return
	}

	user, err := h.userManagementService.AddExistingUser(r.Context(), userID, input.Email)
	if err != nil {
		if err == service.ErrUserNotFound {
			jsonError(w, "usuário não encontrado", http.StatusNotFound)
			return
		}
		if err == service.ErrPermissionDenied {
			jsonError(w, err.Error(), http.StatusForbidden)
			return
		}
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusOK, user)
}

// ChangeRole handles PUT /api/company/users/{id}/role
func (h *UserManagementHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	targetIDStr := chi.URLParam(r, "id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	if input.Role == "" {
		jsonError(w, "cargo é obrigatório", http.StatusBadRequest)
		return
	}

	role, ok := domain.ParseRole(input.Role)
	if !ok {
		jsonError(w, "cargo inválido", http.StatusBadRequest)
		return
	}

	if err := h.userManagementService.ChangeRole(r.Context(), actorUserID, uint(targetID), role); err != nil {
		if err == service.ErrUserNotFound || err == service.ErrUserNotInCompany {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == service.ErrPermissionDenied || err == service.ErrCannotAlterOwner || err == service.ErrCannotAlterAdmin {
			jsonError(w, err.Error(), http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível alterar cargo", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "cargo alterado com sucesso"})
}

// RemoveUser handles DELETE /api/company/users/{id}
func (h *UserManagementHandler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	targetIDStr := chi.URLParam(r, "id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.userManagementService.RemoveFromCompany(r.Context(), actorUserID, uint(targetID)); err != nil {
		if err == service.ErrUserNotFound || err == service.ErrUserNotInCompany {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == service.ErrPermissionDenied || err == service.ErrCannotRemoveOwner {
			jsonError(w, err.Error(), http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível remover usuário", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "usuário removido da empresa com sucesso"})
}
