package domain

// TenantContext represents the tenant context for a request
// It contains information about the authenticated user and their company
// This context is populated by the Tenant middleware and is available throughout the request
type TenantContext struct {
	UserID       uint
	CompanyID    *uint // nil for Core V1 users (no company)
	IsSystemAdmin bool // reserved for future RBAC implementation
}

// HasCompany returns true if the user has a company assigned
func (tc *TenantContext) HasCompany() bool {
	return tc.CompanyID != nil
}

// GetCompanyID returns the company ID, or 0 if nil
func (tc *TenantContext) GetCompanyID() uint {
	if tc.CompanyID == nil {
		return 0
	}
	return *tc.CompanyID
}
