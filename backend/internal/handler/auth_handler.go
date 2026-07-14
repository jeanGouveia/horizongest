package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

var validate = validator.New()

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// --- POST /api/auth/register ---

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "body inválido", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	user, err := h.authService.Register(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			jsonError(w, "e-mail já cadastrado", http.StatusConflict)
			return
		}
		jsonError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// --- POST /api/auth/login ---

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "body inválido", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	result, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			jsonError(w, "e-mail ou senha inválidos", http.StatusUnauthorized)
			return
		}
		jsonError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	// Seta o JWT como Cookie HttpOnly — nunca exposto ao JavaScript
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true em produção (HTTPS)
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
	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"id":    claims.UserID,
		"name":  claims.Name,
		"email": claims.Email,
	})
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
		"error":  "dados inválidos",
		"fields": fields,
	})
}
