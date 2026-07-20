package domain

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
