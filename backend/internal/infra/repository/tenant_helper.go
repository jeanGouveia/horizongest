package repository

import (
	"context"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"gorm.io/gorm"
)

// ApplyTenantFilter applies tenant filtering to a GORM query
// Sprint 3: CompanyID is always required, so filter is always applied
func ApplyTenantFilter(ctx context.Context, db *gorm.DB) *gorm.DB {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		// No tenant context, return query as-is (should not happen in authenticated routes)
		return db
	}

	// Apply tenant filter (CompanyID is always required in Sprint 3)
	return db.Where("company_id = ?", tenantCtx.GetCompanyID())
}

// ApplyTenantFilterWithID applies tenant filtering for queries by ID
// This ensures that id = X AND company_id = tenant (security requirement)
// Sprint 3: CompanyID is always required, so filter is always applied
func ApplyTenantFilterWithID(ctx context.Context, db *gorm.DB, id uint) *gorm.DB {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		// No tenant context, return query as-is (should not happen in authenticated routes)
		return db.Where("id = ?", id)
	}

	// Apply tenant filter with ID (CompanyID is always required in Sprint 3)
	return db.Where("id = ? AND company_id = ?", id, tenantCtx.GetCompanyID())
}

// GetCompanyIDFromContext extracts the CompanyID from the tenant context
// Sprint 3: CompanyID is always required
func GetCompanyIDFromContext(ctx context.Context) (uint, error) {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		return 0, fmt.Errorf("tenant context not found")
	}

	return tenantCtx.CompanyID, nil
}

// HasCompanyContext returns true if the user has a company assigned
// Sprint 3: Always returns true (CompanyID is required)
func HasCompanyContext(ctx context.Context) bool {
	_, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		return false
	}

	return true // CompanyID is always required in Sprint 3
}
