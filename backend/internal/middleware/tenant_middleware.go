package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
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
		// FORENSIC: Log tenant middleware authentication
		log.Printf("[FORENSIC MIDDLEWARE] AUTH_ATTEMPT - URL: %s %s", r.Method, r.URL.Path)
		log.Printf("[FORENSIC MIDDLEWARE] AUTHORIZATION - Header: %s", r.Header.Get("Authorization"))
		log.Printf("[FORENSIC MIDDLEWARE] COOKIE - Cookies: %v", r.Cookies())

		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			log.Printf("[FORENSIC MIDDLEWARE] JWT_FOUND - NÃO (UserID not in context)")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		log.Printf("[FORENSIC MIDDLEWARE] JWT_FOUND - SIM (UserID: %d)", userID)

		// FORENSIC: Log claims before database query
		claims, _ := GetClaimsFromContext(r.Context())
		if claims != nil {
			log.Printf("[FORENSIC MIDDLEWARE] CLAIMS - UserID: %d, CompanyID: %d", claims.UserID, claims.CompanyID)
		}

		// Load user to get CompanyID
		user, err := m.userRepo.FindByID(r.Context(), userID)
		if err != nil {
			log.Printf("[FORENSIC MIDDLEWARE] AUTH_RESULT - FALHA (Failed to load user: %v)", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		if user == nil {
			log.Printf("[FORENSIC MIDDLEWARE] AUTH_RESULT - FALHA (User not found)")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		log.Printf("[FORENSIC MIDDLEWARE] USER_AUTHENTICATED - UserID: %d, CompanyID: %d, Name: %s", user.ID, user.CompanyID, user.Name)
		log.Printf("[FORENSIC MIDDLEWARE] COMPANY_FOUND - CompanyID: %d", user.CompanyID)

		// FORENSIC: Check if CompanyID changed
		if claims != nil && claims.CompanyID != user.CompanyID {
			log.Printf("[FORENSIC MIDDLEWARE] ⚠️ MUDANÇA DETECTADA - claims.CompanyID=%d, user.CompanyID=%d", claims.CompanyID, user.CompanyID)
		}

		// Create TenantContext
		tenantCtx := &domain.TenantContext{
			UserID:    userID,
			CompanyID: user.CompanyID, // Always required - Sprint 3
		}

		// Inject TenantContext into request context
		ctx := context.WithValue(r.Context(), ContextKeyTenant, tenantCtx)

		log.Printf("[FORENSIC MIDDLEWARE] AUTH_RESULT - SUCESSO (Tenant context populated)")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTenantContextFromContext extracts the TenantContext injected by the Tenant middleware
func GetTenantContextFromContext(ctx context.Context) (*domain.TenantContext, bool) {
	tenantCtx, ok := ctx.Value(ContextKeyTenant).(*domain.TenantContext)
	return tenantCtx, ok
}
