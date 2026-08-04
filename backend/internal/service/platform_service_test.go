package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/database"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupPlatformTestDB(t *testing.T) *gorm.DB {
	db, err := database.Connect(database.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		DBName:   "horizongest",
		User:     "horizongest_user",
		Password: "horizongest_secure_password",
		SSLMode:  "disable",
	})
	require.NoError(t, err)

	// Clean up before test
	db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")
	db.Exec("DELETE FROM companies WHERE slug LIKE 'test-%'")
	db.Exec("DELETE FROM platform_users WHERE email LIKE 'test-%'")

	return db
}

func TestPlatformService_CreateCompany_Success(t *testing.T) {
	db := setupPlatformTestDB(t)
	defer func() {
		// Clean up after test
		db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")
		db.Exec("DELETE FROM companies WHERE slug LIKE 'test-%'")
		db.Exec("DELETE FROM platform_users WHERE email LIKE 'test-%'")
	}()

	companyRepo := repository.NewGormCompanyRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	platformUserRepo := repository.NewGormPlatformUserRepository(db)
	platformAuditRepo := repository.NewGormPlatformAuditRepository(db)
	emailSvc := NewEmailService(false, "test@example.com", "TestPlatform")

	platformSvc := NewPlatformService(db, companyRepo, userRepo, platformUserRepo, platformAuditRepo, emailSvc)

	// Create platform admin user
	ctx := context.Background()
	platformAdmin := &domain.PlatformUser{
		Name:         "Test Platform Admin",
		Email:        "test-admin@platform.com",
		PasswordHash: "hashed_password",
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
	}
	err := platformUserRepo.Create(ctx, platformAdmin)
	require.NoError(t, err)

	// Create company with owner
	input := PlatformCreateCompanyInput{
		Name:         "Test Company",
		Slug:         "test-company",
		Description:  "Test Description",
		BusinessType: "restaurant",
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}

	output, err := platformSvc.CreateCompany(ctx, platformAdmin.ID, input, "test-owner@example.com", "password123", "Test Owner")

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Greater(t, output.CompanyID, uint(0))
	assert.Greater(t, output.UserID, uint(0))

	// Verify company was created
	company, err := companyRepo.FindByID(ctx, output.CompanyID)
	require.NoError(t, err)
	require.NotNil(t, company)
	assert.Equal(t, "Test Company", company.Name)
	assert.Equal(t, "test-company", company.Slug)
	assert.True(t, company.Active)

	// Verify owner was created
	owner, err := userRepo.FindByID(ctx, output.UserID)
	require.NoError(t, err)
	require.NotNil(t, owner)
	assert.Equal(t, "Test Owner", owner.Name)
	assert.Equal(t, "test-owner@example.com", owner.Email)
	assert.Equal(t, output.CompanyID, owner.CompanyID)
	assert.Equal(t, domain.RoleOwner, owner.Role)
	assert.True(t, owner.Active)
}

func TestPlatformService_CreateCompany_OwnerFailure_Rollback(t *testing.T) {
	db := setupPlatformTestDB(t)
	defer func() {
		// Clean up after test
		db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")
		db.Exec("DELETE FROM companies WHERE slug LIKE 'test-%'")
		db.Exec("DELETE FROM platform_users WHERE email LIKE 'test-%'")
	}()

	companyRepo := repository.NewGormCompanyRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	platformUserRepo := repository.NewGormPlatformUserRepository(db)
	platformAuditRepo := repository.NewGormPlatformAuditRepository(db)
	emailSvc := NewEmailService(false, "test@example.com", "TestPlatform")

	// Create platform admin user
	ctx := context.Background()
	platformAdmin := &domain.PlatformUser{
		Name:         "Test Platform Admin",
		Email:        "test-admin-rollback@platform.com",
		PasswordHash: "hashed_password",
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
	}
	err := platformUserRepo.Create(ctx, platformAdmin)
	require.NoError(t, err)

	// Create a failing user repository that simulates owner creation failure
	failingUserRepo := &FailingUserRepository{delegate: userRepo}

	platformSvcWithFail := NewPlatformService(db, companyRepo, failingUserRepo, platformUserRepo, platformAuditRepo, emailSvc)

	input := PlatformCreateCompanyInput{
		Name:         "Test Company Rollback",
		Slug:         "test-company-rollback",
		Description:  "Test Description",
		BusinessType: "restaurant",
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}

	// Attempt to create company - should fail
	_, err = platformSvcWithFail.CreateCompany(ctx, platformAdmin.ID, input, "test-owner-fail@example.com", "password123", "Test Owner")

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "simulated owner creation failure")

	// Verify company was NOT created (rollback occurred)
	company, err := companyRepo.FindBySlug(ctx, "test-company-rollback")
	require.NoError(t, err)
	assert.Nil(t, company)

	// Verify owner was NOT created
	owner, err := userRepo.FindByEmail(ctx, "test-owner-fail@example.com")
	require.NoError(t, err)
	assert.Nil(t, owner)
}

func TestPlatformService_CreateCompany_SlugConflict(t *testing.T) {
	db := setupPlatformTestDB(t)
	defer func() {
		// Clean up after test
		db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")
		db.Exec("DELETE FROM companies WHERE slug LIKE 'test-%'")
		db.Exec("DELETE FROM platform_users WHERE email LIKE 'test-%'")
	}()

	companyRepo := repository.NewGormCompanyRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	platformUserRepo := repository.NewGormPlatformUserRepository(db)
	platformAuditRepo := repository.NewGormPlatformAuditRepository(db)
	emailSvc := NewEmailService(false, "test@example.com", "TestPlatform")

	platformSvc := NewPlatformService(db, companyRepo, userRepo, platformUserRepo, platformAuditRepo, emailSvc)

	// Create platform admin user
	ctx := context.Background()
	platformAdmin := &domain.PlatformUser{
		Name:         "Test Platform Admin",
		Email:        "test-admin-conflict@platform.com",
		PasswordHash: "hashed_password",
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
	}
	err := platformUserRepo.Create(ctx, platformAdmin)
	require.NoError(t, err)

	// Create first company
	input1 := PlatformCreateCompanyInput{
		Name:         "Test Company 1",
		Slug:         "test-company-conflict",
		Description:  "Test Description",
		BusinessType: "restaurant",
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}

	output1, err := platformSvc.CreateCompany(ctx, platformAdmin.ID, input1, "test-owner1@example.com", "password123", "Test Owner 1")
	require.NoError(t, err)
	require.NotNil(t, output1)

	// Attempt to create second company with same slug - should fail
	input2 := PlatformCreateCompanyInput{
		Name:         "Test Company 2",
		Slug:         "test-company-conflict", // Same slug
		Description:  "Test Description",
		BusinessType: "restaurant",
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}

	_, err = platformSvc.CreateCompany(ctx, platformAdmin.ID, input2, "test-owner2@example.com", "password123", "Test Owner 2")

	// Assertions
	require.Error(t, err)
	assert.Equal(t, ErrCompanyAlreadyExists, err)

	// Verify only one company exists
	companies, err := companyRepo.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(companies))
}

func TestPlatformService_CreateCompany_PermissionDenied(t *testing.T) {
	db := setupPlatformTestDB(t)
	defer func() {
		// Clean up after test
		db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")
		db.Exec("DELETE FROM companies WHERE slug LIKE 'test-%'")
		db.Exec("DELETE FROM platform_users WHERE email LIKE 'test-%'")
	}()

	companyRepo := repository.NewGormCompanyRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	platformUserRepo := repository.NewGormPlatformUserRepository(db)
	platformAuditRepo := repository.NewGormPlatformAuditRepository(db)
	emailSvc := NewEmailService(false, "test@example.com", "TestPlatform")

	platformSvc := NewPlatformService(db, companyRepo, userRepo, platformUserRepo, platformAuditRepo, emailSvc)

	// Create platform user with non-admin role
	ctx := context.Background()
	platformUser := &domain.PlatformUser{
		Name:         "Test Platform User",
		Email:        "test-user@platform.com",
		PasswordHash: "hashed_password",
		Role:         domain.PlatformRoleSupport, // Not admin
		Active:       true,
	}
	err := platformUserRepo.Create(ctx, platformUser)
	require.NoError(t, err)

	input := PlatformCreateCompanyInput{
		Name:         "Test Company Permission",
		Slug:         "test-company-permission",
		Description:  "Test Description",
		BusinessType: "restaurant",
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}

	// Attempt to create company - should fail due to permission
	_, err = platformSvc.CreateCompany(ctx, platformUser.ID, input, "test-owner@example.com", "password123", "Test Owner")

	// Assertions
	require.Error(t, err)
	assert.Equal(t, ErrPermissionDenied, err)

	// Verify company was NOT created
	company, err := companyRepo.FindBySlug(ctx, "test-company-permission")
	require.NoError(t, err)
	assert.Nil(t, company)
}

// FailingUserRepository simulates a failure in CreateBootstrapOwner
type FailingUserRepository struct {
	delegate ports.UserRepository
}

func (f *FailingUserRepository) Create(ctx context.Context, user *domain.User) error {
	return f.delegate.Create(ctx, user)
}

func (f *FailingUserRepository) CreateBootstrapOwner(ctx context.Context, user *domain.User, companyID uint) error {
	return fmt.Errorf("simulated owner creation failure")
}

func (f *FailingUserRepository) CreateBootstrapOwnerWithTx(ctx context.Context, user *domain.User, companyID uint, tx *gorm.DB) error {
	return fmt.Errorf("simulated owner creation failure")
}

func (f *FailingUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.delegate.FindByEmail(ctx, email)
}

func (f *FailingUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	return f.delegate.FindByID(ctx, id)
}

func (f *FailingUserRepository) Update(ctx context.Context, user *domain.User) error {
	return f.delegate.Update(ctx, user)
}

func (f *FailingUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	return f.delegate.List(ctx)
}
