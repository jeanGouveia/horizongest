package repository

import (
	"context"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"gorm.io/gorm"
)

// ApplyTenantFilter applies tenant filtering to a GORM query
// If CompanyID is nil (Core V1 user), no filter is applied
// If CompanyID is set, WHERE company_id = ? is added automatically
func ApplyTenantFilter(ctx context.Context, db *gorm.DB) *gorm.DB {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		// No tenant context, return query as-is (should not happen in authenticated routes)
		return db
	}

	if !tenantCtx.HasCompany() {
		// Core V1 user, no tenant filtering
		return db
	}

	// Apply tenant filter
	return db.Where("company_id = ?", tenantCtx.GetCompanyID())
}

// ApplyTenantFilterWithID applies tenant filtering for queries by ID
// This ensures that id = X AND company_id = tenant (security requirement)
func ApplyTenantFilterWithID(ctx context.Context, db *gorm.DB, id uint) *gorm.DB {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		// No tenant context, return query as-is (should not happen in authenticated routes)
		return db.Where("id = ?", id)
	}

	if !tenantCtx.HasCompany() {
		// Core V1 user, no tenant filtering
		return db.Where("id = ?", id)
	}

	// Apply tenant filter with ID
	return db.Where("id = ? AND company_id = ?", id, tenantCtx.GetCompanyID())
}

// GetCompanyIDFromContext extracts the CompanyID from the tenant context
// Returns 0 if no company (Core V1 user)
func GetCompanyIDFromContext(ctx context.Context) (*uint, error) {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("tenant context not found")
	}

	return tenantCtx.CompanyID, nil
}

// HasCompanyContext returns true if the user has a company assigned
func HasCompanyContext(ctx context.Context) bool {
	tenantCtx, ok := middleware.GetTenantContextFromContext(ctx)
	if !ok {
		return false
	}

	return tenantCtx.HasCompany()
}
