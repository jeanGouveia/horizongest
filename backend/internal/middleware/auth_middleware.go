package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type contextKey string

const ContextKeyUserID contextKey = "user_id"
const ContextKeyClaims contextKey = "claims"

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[MIDDLEWARE] Request recebida: %s %s", r.Method, r.URL.Path)

		var token string

		// Estratégia 1: Cookie HttpOnly (produção)
		if cookie, err := r.Cookie("auth_token"); err == nil {
			token = cookie.Value
			log.Printf("[MIDDLEWARE] Token encontrado no cookie")
		} else {
			log.Printf("[MIDDLEWARE] Cookie não encontrado: %v", err)
		}

		// Estratégia 2: Authorization header (dev / Postman)
		if token == "" {
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				token = strings.TrimPrefix(h, "Bearer ")
				log.Printf("[MIDDLEWARE] Token encontrado no header Authorization")
			}
		}

		if token == "" {
			log.Printf("[MIDDLEWARE] Token não encontrado - retornando 401")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ValidateToken(r.Context(), token)
		if err != nil {
			log.Printf("[MIDDLEWARE] Token inválido: %v", err)
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Injeta UserID e claims completos no context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyClaims, claims)

		log.Printf("[MIDDLEWARE] UserID injetado no contexto: %d", claims.UserID)
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
