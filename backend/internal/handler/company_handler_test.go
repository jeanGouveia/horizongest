package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCompanyTestDB(t *testing.T) *gorm.DB {
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

	// Auto-migrate the schema
	err = db.AutoMigrate(&repository.GormCompanyModel{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Clean up before each test - drop and recreate schema
	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")
	db.Exec("GRANT ALL ON SCHEMA public TO prato")
	db.Exec("GRANT ALL ON SCHEMA public TO public")
	db.AutoMigrate(&repository.GormCompanyModel{})

	return db
}

func setupTenantContext(ctx context.Context, companyID uint) context.Context {
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: companyID,
	}
	return context.WithValue(ctx, domain.ContextKeyTenant, tenantCtx)
}

func TestCompanyHandler_CreateCompany(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	body := `{"name":"Test Company","slug":"test-company"}`
	req := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCompany(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestCompanyHandler_CreateCompany_InvalidJSON(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	body := `invalid json`
	req := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCompany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCompanyHandler_CreateCompany_SlugExists(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	// Create first company
	body1 := `{"name":"Test Company","slug":"test-company"}`
	req1 := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.CreateCompany(w1, req1)

	// Try to create second company with same slug
	body2 := `{"name":"Another Company","slug":"test-company"}`
	req2 := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	handler.CreateCompany(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w2.Code)
	}
}

func TestCompanyHandler_ListCompanies(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	// Create companies
	for i := 1; i <= 2; i++ {
		body := `{"name":"Test Company ` + string(rune('0'+i)) + `","slug":"test-company-` + string(rune('0'+i)) + `"}`
		req := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.CreateCompany(w, req)
	}

	req := httptest.NewRequest("GET", "/api/companies", nil)
	w := httptest.NewRecorder()

	handler.ListCompanies(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCompanyHandler_GetCompany(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	// Create company
	body := `{"name":"Test Company","slug":"test-company"}`
	reqCreate := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	handler.CreateCompany(wCreate, reqCreate)

	// Get company by ID
	ctx := setupTenantContext(context.Background(), 1)
	req := httptest.NewRequest("GET", "/api/companies/1", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCompany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCompanyHandler_GetCompany_IDOR(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	ctx := setupTenantContext(context.Background(), 100)
	req := httptest.NewRequest("GET", "/api/companies/200", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "200")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCompany(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403 (IDOR protection), got %d", w.Code)
	}
}

func TestCompanyHandler_GetCompany_NotFound(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	ctx := setupTenantContext(context.Background(), 100)
	req := httptest.NewRequest("GET", "/api/companies/100", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "100")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCompany(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestCompanyHandler_GetCurrentCompany(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	// Create company
	body := `{"name":"Test Company","slug":"test-company"}`
	reqCreate := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	handler.CreateCompany(wCreate, reqCreate)

	ctx := setupTenantContext(context.Background(), 1)
	req := httptest.NewRequest("GET", "/api/me/company", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCurrentCompany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCompanyHandler_GetCurrentCompany_NoTenantContext(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	req := httptest.NewRequest("GET", "/api/me/company", nil)
	w := httptest.NewRecorder()

	handler.GetCurrentCompany(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestCompanyHandler_UpdateCompany(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	// Create company
	body := `{"name":"Test Company","slug":"test-company"}`
	reqCreate := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	handler.CreateCompany(wCreate, reqCreate)

	body = `{"name":"Updated Company","slug":"updated-company"}`
	req := httptest.NewRequest("PUT", "/api/companies/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCompany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCompanyHandler_UpdateCompany_NotFound(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	body := `{"name":"Updated Company","slug":"updated-company"}`
	req := httptest.NewRequest("PUT", "/api/companies/100", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "100")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCompany(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestCompanyHandler_DeleteCompany(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	// Create company
	body := `{"name":"Test Company","slug":"test-company"}`
	reqCreate := httptest.NewRequest("POST", "/api/companies", bytes.NewBufferString(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	handler.CreateCompany(wCreate, reqCreate)

	req := httptest.NewRequest("DELETE", "/api/companies/1", nil)
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.DeleteCompany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCompanyHandler_DeleteCompany_NotFound(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := repository.NewGormCompanyRepository(db)
	svc := service.NewCompanyService(repo)
	handler := NewCompanyHandler(svc)

	req := httptest.NewRequest("DELETE", "/api/companies/100", nil)
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "100")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.DeleteCompany(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
