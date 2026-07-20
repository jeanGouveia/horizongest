package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

var validatePlatform = validator.New()

type PlatformAuthHandler struct {
	platformAuthService *service.PlatformAuthService
}

func NewPlatformAuthHandler(platformAuthService *service.PlatformAuthService) *PlatformAuthHandler {
	return &PlatformAuthHandler{platformAuthService: platformAuthService}
}

// POST /api/platform/auth/login
func (h *PlatformAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.PlatformLoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	output, err := h.platformAuthService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrPlatformInvalidCredentials) {
			jsonError(w, "credenciais inválidas", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, service.ErrPlatformUserInactive) {
			jsonError(w, "usuário desativado", http.StatusForbidden)
			return
		}
		jsonError(w, "não foi possível fazer login. Tente novamente.", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, output)
}

// POST /api/platform/auth/logout
func (h *PlatformAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		jsonError(w, "token não fornecido", http.StatusBadRequest)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.platformAuthService.Logout(r.Context(), token); err != nil {
		jsonError(w, "não foi possível fazer logout", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "logout realizado com sucesso",
	})
}

// GET /api/platform/auth/me
func (h *PlatformAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get platform user ID from context (set by middleware)
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	user, err := h.platformAuthService.GetPlatformUser(r.Context(), platformUserID)
	if err != nil {
		if errors.Is(err, service.ErrPlatformUserNotFound) {
			jsonError(w, "usuário não encontrado", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível obter usuário", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, user)
}
