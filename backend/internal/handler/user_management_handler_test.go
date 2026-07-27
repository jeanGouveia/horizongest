package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserManagementTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&repository.GormUserModel{}, &repository.GormCompanyModel{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func setupUserContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, middleware.ContextKeyUserID, userID)
}

func TestUserManagementHandler_ListUsers(t *testing.T) {
	db := setupUserManagementTestDB(t)

	// Create company
	companyRepo := repository.NewGormCompanyRepository(db)
	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company",
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := companyRepo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("CreateCompany failed: %v", err)
	}

	// Create users
	userRepo := repository.NewGormUserRepository(db)
	user1 := &domain.User{
		Name:         "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
	}
	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &domain.User{
		Name:         "User 2",
		Email:        "user2@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleAdmin,
	}
	err = userRepo.Create(context.Background(), user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Create handler with real services
	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	ctx := setupUserContext(context.Background(), user1.ID)
	req := httptest.NewRequest("GET", "/api/company/users", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserManagementHandler_ListUsers_Unauthorized(t *testing.T) {
	db := setupUserManagementTestDB(t)
	userRepo := repository.NewGormUserRepository(db)
	companyRepo := repository.NewGormCompanyRepository(db)
	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	req := httptest.NewRequest("GET", "/api/company/users", nil)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestUserManagementHandler_GetUser(t *testing.T) {
	db := setupUserManagementTestDB(t)

	// Create company
	companyRepo := repository.NewGormCompanyRepository(db)
	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company",
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := companyRepo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("CreateCompany failed: %v", err)
	}

	// Create users
	userRepo := repository.NewGormUserRepository(db)
	user1 := &domain.User{
		Name:         "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
	}
	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &domain.User{
		Name:         "User 2",
		Email:        "user2@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleAdmin,
	}
	err = userRepo.Create(context.Background(), user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	ctx := setupUserContext(context.Background(), user1.ID)
	req := httptest.NewRequest("GET", "/api/company/users/2", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserManagementHandler_GetUser_NotFound(t *testing.T) {
	db := setupUserManagementTestDB(t)

	// Create company
	companyRepo := repository.NewGormCompanyRepository(db)
	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company",
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := companyRepo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("CreateCompany failed: %v", err)
	}

	// Create user
	userRepo := repository.NewGormUserRepository(db)
	user1 := &domain.User{
		Name:         "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
	}
	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	ctx := setupUserContext(context.Background(), user1.ID)
	req := httptest.NewRequest("GET", "/api/company/users/999", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestUserManagementHandler_ChangeRole(t *testing.T) {
	db := setupUserManagementTestDB(t)

	// Create company
	companyRepo := repository.NewGormCompanyRepository(db)
	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company",
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := companyRepo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("CreateCompany failed: %v", err)
	}

	// Create users
	userRepo := repository.NewGormUserRepository(db)
	user1 := &domain.User{
		Name:         "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
	}
	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &domain.User{
		Name:         "User 2",
		Email:        "user2@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleAdmin,
	}
	err = userRepo.Create(context.Background(), user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	body := `{"role":"manager"}`
	ctx := setupUserContext(context.Background(), user1.ID)
	req := httptest.NewRequest("PUT", "/api/company/users/2/role", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ChangeRole(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserManagementHandler_RemoveUser_NotAllowed(t *testing.T) {
	db := setupUserManagementTestDB(t)

	// Create company
	companyRepo := repository.NewGormCompanyRepository(db)
	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company",
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := companyRepo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("CreateCompany failed: %v", err)
	}

	// Create users
	userRepo := repository.NewGormUserRepository(db)
	user1 := &domain.User{
		Name:         "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
	}
	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &domain.User{
		Name:         "User 2",
		Email:        "user2@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleAdmin,
	}
	err = userRepo.Create(context.Background(), user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	ctx := setupUserContext(context.Background(), user1.ID)
	req := httptest.NewRequest("DELETE", "/api/company/users/2", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.RemoveUser(w, req)

	// RemoveFromCompany is not allowed in Sprint 3
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 (not allowed), got %d", w.Code)
	}
}

func TestUserManagementHandler_SetUserActive(t *testing.T) {
	db := setupUserManagementTestDB(t)

	// Create company
	companyRepo := repository.NewGormCompanyRepository(db)
	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company",
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := companyRepo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("CreateCompany failed: %v", err)
	}

	// Create users
	userRepo := repository.NewGormUserRepository(db)
	user1 := &domain.User{
		Name:         "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
	}
	err = userRepo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user2 := &domain.User{
		Name:         "User 2",
		Email:        "user2@example.com",
		PasswordHash: "hash",
		Active:       true,
		CompanyID:    company.ID,
		Role:         domain.RoleAdmin,
	}
	err = userRepo.Create(context.Background(), user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	handler := NewUserManagementHandler(userManagementSvc, userRepo)

	body := `{"active":false}`
	ctx := setupUserContext(context.Background(), user1.ID)
	req := httptest.NewRequest("PUT", "/api/company/users/2/active", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.SetUserActive(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
