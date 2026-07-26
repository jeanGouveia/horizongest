package service

import (
	"context"
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

// TestRBACService_HasRole tests role checking functionality
func TestRBACService_HasRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	hasRole, err := rbacSvc.HasRole(context.Background(), 1, domain.RoleOwner)
	if err != nil {
		t.Fatalf("HasRole failed: %v", err)
	}
	if !hasRole {
		t.Error("Expected user to have Owner role")
	}

	hasRole, err = rbacSvc.HasRole(context.Background(), 1, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("HasRole failed: %v", err)
	}
	if hasRole {
		t.Error("Expected user not to have Admin role")
	}
}

// TestRBACService_HasAnyRole tests checking for any of multiple roles
func TestRBACService_HasAnyRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	hasRole, err := rbacSvc.HasAnyRole(context.Background(), 1, domain.RoleOwner, domain.RoleManager)
	if err != nil {
		t.Fatalf("HasAnyRole failed: %v", err)
	}
	if !hasRole {
		t.Error("Expected user to have at least one of the roles")
	}

	hasRole, err = rbacSvc.HasAnyRole(context.Background(), 1, domain.RoleOwner, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("HasAnyRole failed: %v", err)
	}
	if hasRole {
		t.Error("Expected user not to have any of the roles")
	}
}

// TestRBACService_IsOwner tests Owner role check
func TestRBACService_IsOwner(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	isOwner, err := rbacSvc.IsOwner(context.Background(), 1)
	if err != nil {
		t.Fatalf("IsOwner failed: %v", err)
	}
	if !isOwner {
		t.Error("Expected user to be Owner")
	}
}

// TestRBACService_IsAdmin tests Admin role check
func TestRBACService_IsAdmin(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	isAdmin, err := rbacSvc.IsAdmin(context.Background(), 1)
	if err != nil {
		t.Fatalf("IsAdmin failed: %v", err)
	}
	if !isAdmin {
		t.Error("Expected user to be Admin")
	}
}

// TestRBACService_CanManageCompany tests company management permission
func TestRBACService_CanManageCompany(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canManage, err := rbacSvc.CanManageCompany(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanManageCompany failed: %v", err)
	}
	if !canManage {
		t.Error("Expected Admin to manage company")
	}

	// Test with Manager (should not be able to manage company)
	mockUserRepo.Users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	canManage, err = rbacSvc.CanManageCompany(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanManageCompany failed: %v", err)
	}
	if canManage {
		t.Error("Expected Manager not to manage company")
	}
}

// TestRBACService_CanManageProducts tests product management permission
func TestRBACService_CanManageProducts(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canManage, err := rbacSvc.CanManageProducts(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanManageProducts failed: %v", err)
	}
	if !canManage {
		t.Error("Expected Manager to manage products")
	}

	// Test with Employee (should not be able to manage products)
	mockUserRepo.Users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleEmployee,
		Active:    true,
	}
	canManage, err = rbacSvc.CanManageProducts(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanManageProducts failed: %v", err)
	}
	if canManage {
		t.Error("Expected Employee not to manage products")
	}
}

// TestRBACService_CanManageOrders tests order management permission
func TestRBACService_CanManageOrders(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleEmployee,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canManage, err := rbacSvc.CanManageOrders(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanManageOrders failed: %v", err)
	}
	if !canManage {
		t.Error("Expected Employee to manage orders")
	}
}

// TestRBACService_CanManageUsers tests user management permission
func TestRBACService_CanManageUsers(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canManage, err := rbacSvc.CanManageUsers(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanManageUsers failed: %v", err)
	}
	if !canManage {
		t.Error("Expected Owner to manage users")
	}

	// Test with Admin (should not be able to manage users)
	mockUserRepo.Users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	canManage, err = rbacSvc.CanManageUsers(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanManageUsers failed: %v", err)
	}
	if canManage {
		t.Error("Expected Admin not to manage users")
	}
}

// TestRBACService_CanApproveStockAdjustments tests stock adjustment approval permission
func TestRBACService_CanApproveStockAdjustments(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canApprove, err := rbacSvc.CanApproveStockAdjustments(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanApproveStockAdjustments failed: %v", err)
	}
	if !canApprove {
		t.Error("Expected Admin to approve stock adjustments")
	}

	// Test with Manager (should not be able to approve)
	mockUserRepo.Users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	canApprove, err = rbacSvc.CanApproveStockAdjustments(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanApproveStockAdjustments failed: %v", err)
	}
	if canApprove {
		t.Error("Expected Manager not to approve stock adjustments")
	}
}

// TestRBACService_CanAlterOwnerRole tests Owner role alteration permission
func TestRBACService_CanAlterOwnerRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canAlter, err := rbacSvc.CanAlterOwnerRole(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanAlterOwnerRole failed: %v", err)
	}
	if !canAlter {
		t.Error("Expected Owner to alter Owner role")
	}

	// Test with Admin (should not be able to alter Owner role)
	mockUserRepo.Users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	canAlter, err = rbacSvc.CanAlterOwnerRole(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanAlterOwnerRole failed: %v", err)
	}
	if canAlter {
		t.Error("Expected Admin not to alter Owner role")
	}
}

// TestRBACService_CanAlterAdminRole tests Admin role alteration permission
func TestRBACService_CanAlterAdminRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.Users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canAlter, err := rbacSvc.CanAlterAdminRole(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanAlterAdminRole failed: %v", err)
	}
	if !canAlter {
		t.Error("Expected Owner to alter Admin role")
	}
}
