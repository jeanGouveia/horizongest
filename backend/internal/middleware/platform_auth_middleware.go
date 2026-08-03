package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type PlatformAuthMiddleware struct {
	platformAuthService *service.PlatformAuthService
}

func NewPlatformAuthMiddleware(platformAuthService *service.PlatformAuthService) *PlatformAuthMiddleware {
	return &PlatformAuthMiddleware{platformAuthService: platformAuthService}
}

const ContextKeyPlatformUserID contextKey = "platformUserID"
const ContextKeyPlatformRole contextKey = "platformRole"

// Auth validates the JWT token and adds platform user ID and role to context
func (m *PlatformAuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("===================================================")
		log.Println("TENANT AUTH MIDDLEWARE")
		log.Println("REQUEST:", r.Method, r.URL.Path)
		log.Println("===================================================")

		// Check auth_token
		authTokenCookie, err := r.Cookie("auth_token")
		if err != nil {
			log.Println("auth_token: NÃO ENCONTRADO")
		} else {
			log.Println("auth_token ENCONTRADO")
			log.Println("tamanho:", len(authTokenCookie.Value))
		}

		// Check platform_auth_token
		platformTokenCookie, err := r.Cookie("platform_auth_token")
		if err != nil {
			log.Println("platform_auth_token: NÃO ENCONTRADO")
		} else {
			log.Println("platform_auth_token ENCONTRADO")
			log.Println("tamanho:", len(platformTokenCookie.Value))
		}

		// FORENSIC: Log authentication attempt
		log.Printf("[FORENSIC MIDDLEWARE] AUTH_ATTEMPT - URL: %s %s", r.Method, r.URL.Path)

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			log.Println("MOTIVO DO 401: missing authorization header")
			log.Printf("[FORENSIC MIDDLEWARE] JWT_FOUND - NÃO")
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		log.Printf("[FORENSIC MIDDLEWARE] JWT_FOUND - SIM")

		// Validate token
		userID, role, err := m.platformAuthService.ValidateToken(r.Context(), token)
		if err != nil {
			log.Println("MOTIVO DO 401: invalid token -", err.Error())
			log.Printf("[FORENSIC MIDDLEWARE] JWT_VALID - NÃO - Error: %v", err)
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("[FORENSIC MIDDLEWARE] JWT_VALID - SIM")
		log.Printf("[FORENSIC MIDDLEWARE] CLAIMS - UserID: %d, Role: %s", userID, role)
		log.Printf("[FORENSIC MIDDLEWARE] TOKEN_TYPE - Platform")
		log.Printf("[FORENSIC MIDDLEWARE] AUTH_RESULT - SUCESSO")

		// Add platform user ID and role to context
		ctx := context.WithValue(r.Context(), "platformUserID", userID)
		ctx = context.WithValue(ctx, "platformRole", role)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logJWTClaims(claims jwt.MapClaims) {
	if uid, ok := claims["uid"].(float64); ok {
		log.Printf("uid = %d", int(uid))
	}
	if cid, ok := claims["cid"].(float64); ok {
		log.Printf("cid = %d", int(cid))
	}
	if imp, ok := claims["imp"].(bool); ok {
		log.Printf("imp = %v", imp)
	}
	if opuid, ok := claims["opuid"].(float64); ok {
		log.Printf("opuid = %d", int(opuid))
	}
	if iss, ok := claims["iss"].(string); ok {
		log.Printf("iss = %s", iss)
	}
	if sub, ok := claims["sub"].(string); ok {
		log.Printf("sub = %s", sub)
	}
	if exp, ok := claims["exp"].(float64); ok {
		log.Printf("exp = %d", int(exp))
	}
}

// RequireAdmin checks if the platform user has admin role
func (m *PlatformAuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("platformRole").(domain.PlatformRole)
		if !ok {
			http.Error(w, "missing role in context", http.StatusInternalServerError)
			return
		}

		if role != domain.PlatformRoleAdmin {
			http.Error(w, "admin role required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetPlatformUserID extracts platform user ID from context
func GetPlatformUserID(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(ContextKeyPlatformUserID).(uint)
	if !ok {
		return 0, errors.New("platform user ID not found in context")
	}
	return userID, nil
}

// GetPlatformRole extracts platform role from context
func GetPlatformRole(ctx context.Context) (domain.PlatformRole, error) {
	role, ok := ctx.Value(ContextKeyPlatformRole).(domain.PlatformRole)
	if !ok {
		return domain.PlatformRoleSupport, errors.New("platform role not found in context")
	}
	return role, nil
}
