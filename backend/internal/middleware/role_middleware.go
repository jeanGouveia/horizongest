package middleware

import (
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

// RoleMiddleware provides role-based access control middleware
// Usage examples:
//
//	roleMw.Require(domain.RoleOwner)
//	roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin)
type RoleMiddleware struct {
	rbacService *service.RBACService
}

func NewRoleMiddleware(rbacService *service.RBACService) *RoleMiddleware {
	return &RoleMiddleware{rbacService: rbacService}
}

// Require creates middleware that requires a specific role
// During impersonation, platform admins are granted Owner permissions automatically
func (m *RoleMiddleware) Require(role domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				jsonError(w, "não autorizado", http.StatusUnauthorized)
				return
			}

			// Check if user is impersonating - if so, grant Owner permissions
			isImpersonating, _ := GetIsImpersonating(r.Context())
			if isImpersonating {
				// During impersonation, platform admin has Owner permissions
				next.ServeHTTP(w, r)
				return
			}

			hasRole, err := m.rbacService.HasRole(r.Context(), userID, role)
			if err != nil {
				jsonError(w, "erro ao verificar permissões", http.StatusInternalServerError)
				return
			}

			if !hasRole {
				jsonError(w, "permissão insuficiente", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAny creates middleware that requires any of the specified roles
// During impersonation, platform admins are granted Owner permissions automatically
func (m *RoleMiddleware) RequireAny(roles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				jsonError(w, "não autorizado", http.StatusUnauthorized)
				return
			}

			// Check if user is impersonating - if so, grant Owner permissions
			isImpersonating, _ := GetIsImpersonating(r.Context())
			if isImpersonating {
				// During impersonation, platform admin has Owner permissions
				next.ServeHTTP(w, r)
				return
			}

			hasAnyRole, err := m.rbacService.HasAnyRole(r.Context(), userID, roles...)
			if err != nil {
				jsonError(w, "erro ao verificar permissões", http.StatusInternalServerError)
				return
			}

			if !hasAnyRole {
				jsonError(w, "permissão insuficiente", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission creates middleware that requires a specific permission
// This is a convenience wrapper around RBACService methods
func (m *RoleMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				jsonError(w, "não autorizado", http.StatusUnauthorized)
				return
			}

			var hasPermission bool
			var err error

			switch permission {
			case "manage_company":
				hasPermission, err = m.rbacService.CanManageCompany(r.Context(), userID)
			case "manage_products":
				hasPermission, err = m.rbacService.CanManageProducts(r.Context(), userID)
			case "manage_orders":
				hasPermission, err = m.rbacService.CanManageOrders(r.Context(), userID)
			case "manage_users":
				hasPermission, err = m.rbacService.CanManageUsers(r.Context(), userID)
			case "manage_settings":
				hasPermission, err = m.rbacService.CanManageSettings(r.Context(), userID)
			case "view_reports":
				hasPermission, err = m.rbacService.CanViewReports(r.Context(), userID)
			default:
				jsonError(w, "permissão desconhecida", http.StatusBadRequest)
				return
			}

			if err != nil {
				jsonError(w, "erro ao verificar permissões", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				jsonError(w, "permissão insuficiente", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions for JSON responses (should be moved to a shared package)
func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
