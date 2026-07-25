package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type contextKey string

const ContextKeyUserID contextKey = "user_id"
const ContextKeyClaims contextKey = "claims"
const ContextKeyIsImpersonating contextKey = "is_impersonating"
const ContextKeyOriginalPlatformUserID contextKey = "original_platform_user_id"

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[FORENSIC] AuthMiddleware - Request recebida: %s %s", r.Method, r.URL.Path)
		log.Printf("[FORENSIC] AuthMiddleware - Request ID: %s", r.Header.Get("X-Request-ID"))

		// FORENSIC: Log ALL cookies received
		log.Printf("[FORENSIC] AuthMiddleware - TODOS os cookies recebidos:")
		for _, c := range r.Cookies() {
			log.Printf("[FORENSIC] AuthMiddleware - [COOKIE] %s=%s", c.Name, c.Value)
		}

		// FORENSIC: Log Authorization header
		authHeader := r.Header.Get("Authorization")
		log.Printf("[FORENSIC] AuthMiddleware - Authorization Header: %s", authHeader)

		var token string
		var tokenSource string

		// Estratégia 1: Cookie HttpOnly (produção)
		if cookie, err := r.Cookie("auth_token"); err == nil {
			token = cookie.Value
			tokenSource = "cookie auth_token"
			// FORENSIC: Log token bruto
			log.Printf("[FORENSIC] AuthMiddleware - Token encontrado no cookie auth_token: %s", token)
		} else {
			log.Printf("[FORENSIC] AuthMiddleware - Cookie auth_token não encontrado: %v", err)
		}

		// Estratégia 2: Authorization header (dev / Postman)
		if token == "" {
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				token = strings.TrimPrefix(h, "Bearer ")
				tokenSource = "Authorization header"
				log.Printf("[FORENSIC] AuthMiddleware - Token encontrado no header Authorization: %s", token)
			}
		}

		if token == "" {
			log.Printf("[FORENSIC] AuthMiddleware - Token não encontrado em nenhum lugar - retornando 401")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// FORENSIC: Log which token was chosen
		log.Printf("[FORENSIC] AuthMiddleware - Token escolhido - Origem: %s, Token: %s", tokenSource, token)

		claims, err := m.authService.ValidateToken(r.Context(), token)
		if err != nil {
			log.Printf("[FORENSIC] AuthMiddleware - Token inválido: %v", err)
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// FORENSIC: Log claims
		log.Printf("[FORENSIC] AuthMiddleware - Claims - UserID: %d, CompanyID: %d, Email: %s, Name: %s, IsImpersonating: %v, OriginalPlatformUserID: %d",
			claims.UserID, claims.CompanyID, claims.Email, claims.Name, claims.IsImpersonating, claims.OriginalPlatformUserID)

		// Injeta UserID e claims completos no context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyClaims, claims)

		// Inject impersonation context if present
		if claims.IsImpersonating {
			ctx = context.WithValue(ctx, ContextKeyIsImpersonating, true)
			ctx = context.WithValue(ctx, ContextKeyOriginalPlatformUserID, claims.OriginalPlatformUserID)
			log.Printf("[FORENSIC] AuthMiddleware - Impersonation detected - PlatformUserID: %d, CompanyID: %d", claims.OriginalPlatformUserID, claims.CompanyID)
		}

		log.Printf("[FORENSIC] AuthMiddleware - UserID injetado no contexto: %d", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extrai o UserID injetado pelo middleware.
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uint)
	return id, ok
}

// GetClaimsFromContext extrai os claims completos (UserID, Email, Name).
func GetClaimsFromContext(ctx context.Context) (*service.JWTClaims, bool) {
	claims, ok := ctx.Value(ContextKeyClaims).(*service.JWTClaims)
	return claims, ok
}

// GetIsImpersonating extracts the impersonation flag from context
func GetIsImpersonating(ctx context.Context) (bool, bool) {
	isImpersonating, ok := ctx.Value(ContextKeyIsImpersonating).(bool)
	return isImpersonating, ok
}

// GetOriginalPlatformUserID extracts the original platform user ID from context during impersonation
func GetOriginalPlatformUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(ContextKeyOriginalPlatformUserID).(uint)
	return userID, ok
}
