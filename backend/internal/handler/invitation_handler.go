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

type InvitationHandler struct {
	invitationService *service.InvitationService
	userRepo          ports.UserRepository // To get CompanyID from user
}

func NewInvitationHandler(invitationService *service.InvitationService, userRepo ports.UserRepository) *InvitationHandler {
	return &InvitationHandler{
		invitationService: invitationService,
		userRepo:          userRepo,
	}
}

// getUserCompanyID gets the user's CompanyID from the repository
func (h *InvitationHandler) getUserCompanyID(ctx context.Context, userID uint) (uint, error) {
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("usuário não encontrado")
	}
	return user.CompanyID, nil
}

// CreateInvitation handles POST /api/company/invitations
func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		jsonError(w, "e-mail é obrigatório", http.StatusBadRequest)
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

	invitation, err := h.invitationService.CreateInvitation(r.Context(), userID, companyID, input.Email, role)
	if err != nil {
		if err == service.ErrPermissionDenied {
			jsonError(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrDuplicateInvitation || err == service.ErrUserAlreadyInCompany || err == service.ErrUserBelongsToOtherCompany {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonError(w, "não foi possível criar convite", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, invitation)
}

// ListInvitations handles GET /api/company/invitations
func (h *InvitationHandler) ListInvitations(w http.ResponseWriter, r *http.Request) {
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

	invitations, err := h.invitationService.ListInvitations(r.Context(), companyID)
	if err != nil {
		jsonError(w, "não foi possível listar convites", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, invitations)
}

// GetInvitation handles GET /api/company/invitations/{id}
func (h *InvitationHandler) GetInvitation(w http.ResponseWriter, r *http.Request) {
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

	invitationIDStr := chi.URLParam(r, "id")
	invitationID, err := strconv.ParseUint(invitationIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	invitation, err := h.invitationService.GetInvitation(r.Context(), companyID, uint(invitationID))
	if err != nil {
		if err == service.ErrInvitationNotFound {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar convite", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, invitation)
}

// RevokeInvitation handles DELETE /api/company/invitations/{id}
func (h *InvitationHandler) RevokeInvitation(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	companyID, err := h.getUserCompanyID(r.Context(), actorUserID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusForbidden)
		return
	}

	invitationIDStr := chi.URLParam(r, "id")
	invitationID, err := strconv.ParseUint(invitationIDStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.invitationService.RevokeInvitation(r.Context(), actorUserID, companyID, uint(invitationID)); err != nil {
		if err == service.ErrInvitationNotFound {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == service.ErrPermissionDenied {
			jsonError(w, err.Error(), http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível revogar convite", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "convite revogado com sucesso"})
}

// GetInvitationByToken handles GET /api/invitations/:token (public endpoint)
func (h *InvitationHandler) GetInvitationByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		jsonError(w, "token é obrigatório", http.StatusBadRequest)
		return
	}

	invitation, err := h.invitationService.GetInvitationByToken(r.Context(), token)
	if err != nil {
		if err == service.ErrInvitationNotFound || err == service.ErrInvitationExpired {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar convite", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, invitation)
}

// AcceptInvitation handles POST /api/invitations/accept (requires authentication)
func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado. Faça login para aceitar o convite.", http.StatusUnauthorized)
		return
	}

	// Get user email to validate against invitation email
	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil || user == nil {
		jsonError(w, "não foi possível carregar dados do usuário", http.StatusInternalServerError)
		return
	}

	var input struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	if input.Token == "" {
		jsonError(w, "token é obrigatório", http.StatusBadRequest)
		return
	}

	if err := h.invitationService.AcceptInvitation(r.Context(), input.Token, user.Email); err != nil {
		if err == service.ErrInvitationNotFound || err == service.ErrInvitationExpired || err == service.ErrInvitationRevoked || err == service.ErrInvitationAlreadyUsed {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err == service.ErrUserNotFound {
			jsonError(w, "usuário não encontrado. Por favor, realize o cadastro primeiro.", http.StatusNotFound)
			return
		}
		if err == service.ErrUserBelongsToOtherCompany {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonError(w, "não foi possível aceitar convite", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "convite aceito com sucesso"})
}
