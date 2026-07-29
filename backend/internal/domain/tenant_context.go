package domain

import "context"

// TenantContext represents the tenant context for a request
// It contains information about the authenticated user and their company
// This context is populated by the Tenant middleware and is available throughout the request
// Sprint 3: CompanyID is always required (NOT NULL) - no more Core V1 users without company
type TenantContext struct {
	UserID    uint
	CompanyID uint // Always required - Sprint 3
}

// GetCompanyID returns the company ID
func (tc *TenantContext) GetCompanyID() uint {
	return tc.CompanyID
}

// ContextKeyTenant is the key used to store TenantContext in context
type contextKey string

const ContextKeyTenant contextKey = "tenant_context"

// GetTenantContextFromContext extracts the TenantContext from the request context
// This is used by repositories to apply tenant filtering
func GetTenantContextFromContext(ctx context.Context) (*TenantContext, bool) {
	tenantCtx, ok := ctx.Value(ContextKeyTenant).(*TenantContext)
	return tenantCtx, ok
}
