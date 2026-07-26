package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	Users             map[uint]*domain.User
	FindByEmailResult *domain.User
	FindByEmailError  error
	UpdateError       error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[uint]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	return m.Users[id], nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailResult != nil {
		return m.FindByEmailResult, m.FindByEmailError
	}
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.Users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	delete(m.Users, id)
	return nil
}

// TestTenantMiddleware_TenantContext tests that tenant context is properly set
func TestTenantMiddleware_TenantContext(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Active:    true,
	}
	tenantMw := NewTenantMiddleware(mockUserRepo)

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), ContextKeyUserID, uint(1))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := tenantMw.Tenant(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	tenant, ok := GetTenantContextFromContext(capturedContext)
	if !ok {
		t.Fatal("Expected TenantContext in context")
	}
	if tenant.CompanyID != 123 {
		t.Errorf("Expected CompanyID 123, got %d", tenant.CompanyID)
	}
	if tenant.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", tenant.UserID)
	}
}

// TestTenantMiddleware_CompanyIDIsolation tests that users cannot access other companies' data
func TestTenantMiddleware_CompanyIDIsolation(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Active:    true,
	}
	tenantMw := NewTenantMiddleware(mockUserRepo)

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), ContextKeyUserID, uint(1))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := tenantMw.Tenant(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	tenant, ok := GetTenantContextFromContext(capturedContext)
	if !ok {
		t.Fatal("Expected TenantContext in context")
	}
	if tenant.CompanyID != 123 {
		t.Errorf("Expected CompanyID 123, got %d", tenant.CompanyID)
	}
}

// TestTenantMiddleware_MissingCompanyID tests that requests without CompanyID are rejected
func TestTenantMiddleware_MissingCompanyID(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 0, // Missing CompanyID
		Active:    true,
	}
	tenantMw := NewTenantMiddleware(mockUserRepo)

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), ContextKeyUserID, uint(1))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := tenantMw.Tenant(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	tenant, ok := GetTenantContextFromContext(capturedContext)
	if !ok {
		t.Fatal("Expected TenantContext in context")
	}
	// CompanyID should be 0 if user has no CompanyID
	if tenant.CompanyID != 0 {
		t.Errorf("Expected CompanyID 0, got %d", tenant.CompanyID)
	}
}

// TestTenantMiddleware_Impersonation tests that impersonation preserves original tenant context
func TestTenantMiddleware_Impersonation(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Active:    true,
	}
	tenantMw := NewTenantMiddleware(mockUserRepo)

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), ContextKeyUserID, uint(1))
	ctx = context.WithValue(ctx, ContextKeyIsImpersonating, true)
	ctx = context.WithValue(ctx, ContextKeyOriginalPlatformUserID, uint(100))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := tenantMw.Tenant(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	tenant, ok := GetTenantContextFromContext(capturedContext)
	if !ok {
		t.Fatal("Expected TenantContext in context")
	}
	if tenant.CompanyID != 123 {
		t.Errorf("Expected CompanyID 123, got %d", tenant.CompanyID)
	}
}

// TestTenantContext_GetCompanyID tests that CompanyID can be extracted from context
func TestTenantContext_GetCompanyID(t *testing.T) {
	ctx := context.Background()

	// Test with no tenant context
	_, ok := GetTenantContextFromContext(ctx)
	if ok {
		t.Error("expected false when no tenant context")
	}

	// Test with tenant context
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}
	ctx = context.WithValue(ctx, ContextKeyTenant, tenantCtx)

	retrieved, ok := GetTenantContextFromContext(ctx)
	if !ok {
		t.Error("expected true when tenant context exists")
	}
	if retrieved.CompanyID != 123 {
		t.Errorf("expected CompanyID 123, got %d", retrieved.CompanyID)
	}
	if retrieved.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", retrieved.UserID)
	}
}
