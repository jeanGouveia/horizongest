package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

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
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		// Validate token
		userID, role, err := m.platformAuthService.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Add platform user ID and role to context
		ctx := context.WithValue(r.Context(), "platformUserID", userID)
		ctx = context.WithValue(ctx, "platformRole", role)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
