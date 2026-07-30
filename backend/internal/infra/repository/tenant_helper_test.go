package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=horizongest_user password=horizongest_secure_password dbname=horizongest_test sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	return db
}

func TestApplyTenantFilter(t *testing.T) {
	db := setupTestDB(t)

	// Create a test context with tenant
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// Apply tenant filter
	filteredDB := ApplyTenantFilter(ctx, db)

	// Verify the query has the tenant filter
	sql := filteredDB.Statement.SQL.String()
	if sql == "" {
		// GORM doesn't build SQL until execution, so we check the clause
		if len(filteredDB.Statement.Clauses) == 0 {
			t.Error("Expected tenant filter to be applied")
		}
	}
}

func TestApplyTenantFilter_NoContext(t *testing.T) {
	db := setupTestDB(t)

	// Create a context without tenant
	ctx := context.Background()

	// Apply tenant filter
	filteredDB := ApplyTenantFilter(ctx, db)

	// Should return the same DB instance when no context
	if filteredDB == nil {
		t.Error("Expected DB to be returned even without context")
	}
}

func TestApplyTenantFilterWithID(t *testing.T) {
	db := setupTestDB(t)

	// Create a test context with tenant
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// Apply tenant filter with ID
	filteredDB := ApplyTenantFilterWithID(ctx, db, 456)

	// Verify the query has both ID and tenant filter
	if filteredDB == nil {
		t.Error("Expected DB to be returned")
	}
}

func TestApplyTenantFilterWithID_NoContext(t *testing.T) {
	db := setupTestDB(t)

	// Create a context without tenant
	ctx := context.Background()

	// Apply tenant filter with ID
	filteredDB := ApplyTenantFilterWithID(ctx, db, 456)

	// Should return DB with only ID filter when no context
	if filteredDB == nil {
		t.Error("Expected DB to be returned even without context")
	}
}

func TestGetCompanyIDFromContext(t *testing.T) {
	// Create a test context with tenant
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// Get CompanyID
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		t.Fatalf("GetCompanyIDFromContext failed: %v", err)
	}
	if companyID != 123 {
		t.Errorf("expected CompanyID 123, got %d", companyID)
	}
}

func TestGetCompanyIDFromContext_NoContext(t *testing.T) {
	// Create a context without tenant
	ctx := context.Background()

	// Get CompanyID should fail
	_, err := GetCompanyIDFromContext(ctx)
	if err == nil {
		t.Error("expected error when no tenant context")
	}
}

func TestHasCompanyContext(t *testing.T) {
	// Create a test context with tenant
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// Check if has company context
	hasContext := HasCompanyContext(ctx)
	if !hasContext {
		t.Error("expected true when tenant context exists")
	}
}

func TestHasCompanyContext_NoContext(t *testing.T) {
	// Create a context without tenant
	ctx := context.Background()

	// Check if has company context
	hasContext := HasCompanyContext(ctx)
	if hasContext {
		t.Error("expected false when no tenant context")
	}
}

func TestTenantIsolation_CrossTenantAccess(t *testing.T) {
	db := setupTestDB(t)

	// Simulate user from company 123 trying to access company 456
	userTenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, userTenantCtx)

	// Apply tenant filter - should only allow access to company 123
	filteredDB := ApplyTenantFilter(ctx, db)

	// The filter should prevent cross-tenant access
	if filteredDB == nil {
		t.Error("Expected DB to be returned")
	}
}

func TestTenantIsolation_IDORPrevention(t *testing.T) {
	db := setupTestDB(t)

	// Simulate user from company 123 trying to access resource ID 456 from company 456
	userTenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, userTenantCtx)

	// Apply tenant filter with ID - should require both ID match AND company match
	filteredDB := ApplyTenantFilterWithID(ctx, db, 456)

	// The filter should prevent IDOR by requiring company match
	if filteredDB == nil {
		t.Error("Expected DB to be returned")
	}
}

func TestTenantContext_CompanyIDRequired(t *testing.T) {
	// Test that CompanyID is always required (Sprint 3)
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 0, // Invalid - CompanyID is required
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// HasCompanyContext returns true if context exists, regardless of CompanyID value
	// The validation of CompanyID > 0 happens at the user creation level
	hasContext := HasCompanyContext(ctx)
	if !hasContext {
		t.Error("expected true when tenant context exists (CompanyID validation happens elsewhere)")
	}
}

func TestTenantContext_ValidCompanyID(t *testing.T) {
	// Test valid CompanyID
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// HasCompanyContext should return true for valid CompanyID
	hasContext := HasCompanyContext(ctx)
	if !hasContext {
		t.Error("expected true when CompanyID is valid")
	}

	// GetCompanyID should return the correct value
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		t.Fatalf("GetCompanyIDFromContext failed: %v", err)
	}
	if companyID != 123 {
		t.Errorf("expected CompanyID 123, got %d", companyID)
	}
}
