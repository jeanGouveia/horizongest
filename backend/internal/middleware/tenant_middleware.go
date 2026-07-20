package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

const ContextKeyTenant contextKey = "tenant_context"

type TenantMiddleware struct {
	userRepo ports.UserRepository
}

func NewTenantMiddleware(userRepo ports.UserRepository) *TenantMiddleware {
	return &TenantMiddleware{userRepo: userRepo}
}

// Tenant loads the user's company and populates the TenantContext
// This middleware must run AFTER AuthMiddleware
func (m *TenantMiddleware) Tenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			log.Printf("[TENANT MIDDLEWARE] UserID not found in context - returning 401")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Load user to get CompanyID
		user, err := m.userRepo.FindByID(r.Context(), userID)
		if err != nil {
			log.Printf("[TENANT MIDDLEWARE] Failed to load user: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		if user == nil {
			log.Printf("[TENANT MIDDLEWARE] User not found - returning 401")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Create TenantContext
		tenantCtx := &domain.TenantContext{
			UserID:    userID,
			CompanyID: user.CompanyID, // Always required - Sprint 3
		}

		// Inject TenantContext into request context
		ctx := context.WithValue(r.Context(), ContextKeyTenant, tenantCtx)

		log.Printf("[TENANT MIDDLEWARE] TenantContext populated - UserID: %d, CompanyID: %v", userID, user.CompanyID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTenantContextFromContext extracts the TenantContext injected by the Tenant middleware
func GetTenantContextFromContext(ctx context.Context) (*domain.TenantContext, bool) {
	tenantCtx, ok := ctx.Value(ContextKeyTenant).(*domain.TenantContext)
	return tenantCtx, ok
}
