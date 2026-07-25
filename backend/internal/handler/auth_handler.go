package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
	"github.com/jeanGouveia/horizongest/backend/internal/util"
)

var validate = validator.New()

type AuthHandler struct {
	authService *service.AuthService
	userRepo    ports.UserRepository
	sanitizer   *util.Sanitizer
}

func NewAuthHandler(authService *service.AuthService, userRepo ports.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		sanitizer:   util.NewSanitizer(),
	}
}

// --- POST /api/auth/register (REMOVED - Sprint 3) ---
// Public registration has been removed. Companies are now created by platform administrators only.

// --- POST /api/auth/login ---

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Sanitize inputs (Sprint 3.4 - Security Hardening)
	sanitizedEmail, err := h.sanitizer.SanitizeEmail(input.Email)
	if err != nil {
		jsonError(w, fmt.Sprintf("email inválido: %s", err.Error()), http.StatusBadRequest)
		return
	}
	input.Email = sanitizedEmail

	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	result, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			jsonError(w, "e-mail ou senha incorretos. Verifique suas credenciais.", http.StatusUnauthorized)
			return
		}
		jsonError(w, "não foi possível fazer login. Tente novamente.", http.StatusInternalServerError)
		return
	}

	// Seta o JWT como Cookie HttpOnly — nunca exposto ao JavaScript
	secureCookie := os.Getenv("ENVIRONMENT") == "production"
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"id":    result.User.ID,
		"name":  result.User.Name,
		"email": result.User.Email,
	})
}

// --- POST /api/auth/logout ---

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract token from cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		// Blacklist the token on server
		_ = h.authService.Logout(r.Context(), cookie.Value)
	}

	// Zera o cookie com expiração no passado
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	jsonResponse(w, http.StatusOK, map[string]string{"message": "logout realizado"})
}

// --- GET /api/me (rota protegida) ---

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado. Faça login novamente.", http.StatusUnauthorized)
		return
	}

	// DEBUG: Log JWT claims from context
	claims, _ := middleware.GetClaimsFromContext(r.Context())
	if claims != nil {
		log.Printf("[DEBUG] /api/me - JWT recebido: UserID=%d, CompanyID=%d, Email=%s, Name=%s, IsImpersonating=%v",
			claims.UserID, claims.CompanyID, claims.Email, claims.Name, claims.IsImpersonating)
	}
	log.Printf("[DEBUG] /api/me - UserID do contexto: %d", userID)

	// Get full user data to include CompanyID
	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil || user == nil {
		jsonError(w, "não foi possível carregar dados do usuário", http.StatusInternalServerError)
		return
	}

	// DEBUG: Log user loaded from database
	log.Printf("[DEBUG] /api/me - Usuário carregado do banco: ID=%d, Nome=%s, CompanyID=%d, Email=%s",
		user.ID, user.Name, user.CompanyID, user.Email)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"company_id": user.CompanyID,
	})
}

// --- PUT /api/me (rota protegida) ---

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado. Faça login novamente.", http.StatusUnauthorized)
		return
	}

	var input service.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Sanitize inputs (Sprint 3.4 - Security Hardening)
	sanitizedName, err := h.sanitizer.SanitizeName(input.Name)
	if err != nil {
		jsonError(w, fmt.Sprintf("nome inválido: %s", err.Error()), http.StatusBadRequest)
		return
	}
	input.Name = sanitizedName

	sanitizedEmail, err := h.sanitizer.SanitizeEmail(input.Email)
	if err != nil {
		jsonError(w, fmt.Sprintf("email inválido: %s", err.Error()), http.StatusBadRequest)
		return
	}
	input.Email = sanitizedEmail

	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	user, err := h.authService.UpdateProfile(r.Context(), claims.UserID, input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			jsonError(w, "este e-mail já está cadastrado por outro usuário.", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível atualizar o perfil. Tente novamente.", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// --- POST /api/me/change-password (rota protegida) ---

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		jsonError(w, "não autorizado. Faça login novamente.", http.StatusUnauthorized)
		return
	}

	var input service.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	err := h.authService.ChangePassword(r.Context(), claims.UserID, input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			jsonError(w, "senha atual incorreta. Verifique e tente novamente.", http.StatusUnauthorized)
			return
		}
		jsonError(w, "não foi possível alterar a senha. Tente novamente.", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "senha alterada com sucesso"})
}

// --- POST /api/auth/request-password-reset ---

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var input service.RequestPasswordResetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Sanitize inputs (Sprint 3.4 - Security Hardening)
	sanitizedEmail, err := h.sanitizer.SanitizeEmail(input.Email)
	if err != nil {
		jsonError(w, fmt.Sprintf("email inválido: %s", err.Error()), http.StatusBadRequest)
		return
	}
	input.Email = sanitizedEmail

	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	err = h.authService.RequestPasswordReset(r.Context(), input)
	if err != nil {
		jsonError(w, "não foi possível solicitar recuperação de senha. Tente novamente.", http.StatusInternalServerError)
		return
	}

	// Always return success to avoid email enumeration
	jsonResponse(w, http.StatusOK, map[string]string{"message": "se o e-mail estiver cadastrado, você receberá instruções para recuperar sua senha"})
}

// --- POST /api/auth/reset-password ---

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input service.ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	err := h.authService.ResetPassword(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidResetToken) {
			jsonError(w, "token inválido ou expirado. Solicite uma nova recuperação de senha.", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrResetTokenAlreadyUsed) {
			jsonError(w, "este token já foi utilizado. Solicite uma nova recuperação de senha.", http.StatusBadRequest)
			return
		}
		jsonError(w, "não foi possível redefinir a senha. Tente novamente.", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "senha redefinida com sucesso"})
}

// --- helpers de resposta ---

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

func jsonValidationError(w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		jsonError(w, "dados inválidos", http.StatusBadRequest)
		return
	}
	fields := make(map[string]string, len(ve))
	for _, fe := range ve {
		fieldName := fe.Field()
		tag := fe.Tag()
		param := fe.Param()

		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s é obrigatório", fieldName)
		case "min":
			message = fmt.Sprintf("%s deve ter no mínimo %s caracteres", fieldName, param)
		case "max":
			message = fmt.Sprintf("%s deve ter no máximo %s caracteres", fieldName, param)
		case "gt":
			message = fmt.Sprintf("%s deve ser maior que %s", fieldName, param)
		case "gte":
			message = fmt.Sprintf("%s deve ser maior ou igual a %s", fieldName, param)
		case "lt":
			message = fmt.Sprintf("%s deve ser menor que %s", fieldName, param)
		case "lte":
			message = fmt.Sprintf("%s deve ser menor ou igual a %s", fieldName, param)
		case "email":
			message = fmt.Sprintf("%s deve ser um email válido", fieldName)
		case "oneof":
			message = fmt.Sprintf("%s deve ser um dos seguintes valores: %s", fieldName, param)
		default:
			message = fmt.Sprintf("%s falhou na validação: %s", fieldName, tag)
		}
		fields[fieldName] = message
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "dados inválidos",
		"fields":  fields,
		"details": "Verifique os campos marcados em vermelho",
	})
}
