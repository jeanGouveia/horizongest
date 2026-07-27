package service

import (
	"context"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// TestRBACService_HasRole tests role checking functionality
func TestRBACService_HasRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[2] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[2] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[2] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[2] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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
	mockUserRepo.users[2] = &domain.User{
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
	mockUserRepo.users[1] = &domain.User{
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

	// Test with Admin (should not be able to alter Admin role)
	mockUserRepo.users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	canAlter, err = rbacSvc.CanAlterAdminRole(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanAlterAdminRole failed: %v", err)
	}
	if canAlter {
		t.Error("Expected Admin not to alter Admin role")
	}
}

// TestRBACService_IsManager tests Manager role check
func TestRBACService_IsManager(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	isManager, err := rbacSvc.IsManager(context.Background(), 1)
	if err != nil {
		t.Fatalf("IsManager failed: %v", err)
	}
	if !isManager {
		t.Error("Expected user to be Manager")
	}

	// Test with non-manager
	isManager, err = rbacSvc.IsManager(context.Background(), 1)
	if err != nil {
		t.Fatalf("IsManager failed: %v", err)
	}
}

// TestRBACService_CanManageSettings tests settings management permission
func TestRBACService_CanManageSettings(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canManage, err := rbacSvc.CanManageSettings(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanManageSettings failed: %v", err)
	}
	if !canManage {
		t.Error("Expected Admin to manage settings")
	}

	// Test with Manager (should not be able to manage settings)
	mockUserRepo.users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	canManage, err = rbacSvc.CanManageSettings(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanManageSettings failed: %v", err)
	}
	if canManage {
		t.Error("Expected Manager not to manage settings")
	}
}

// TestRBACService_CanViewReports tests report viewing permission
func TestRBACService_CanViewReports(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	canView, err := rbacSvc.CanViewReports(context.Background(), 1)
	if err != nil {
		t.Fatalf("CanViewReports failed: %v", err)
	}
	if !canView {
		t.Error("Expected Manager to view reports")
	}

	// Test with Employee (should not be able to view reports)
	mockUserRepo.users[2] = &domain.User{
		ID:        2,
		Name:      "Test User 2",
		Email:     "test2@example.com",
		CompanyID: 123,
		Role:      domain.RoleEmployee,
		Active:    true,
	}
	canView, err = rbacSvc.CanViewReports(context.Background(), 2)
	if err != nil {
		t.Fatalf("CanViewReports failed: %v", err)
	}
	if canView {
		t.Error("Expected Employee not to view reports")
	}
}

// TestRBACService_GetUserRole tests getting user role
func TestRBACService_GetUserRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.users[1] = &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	role, err := rbacSvc.GetUserRole(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUserRole failed: %v", err)
	}
	if role != domain.RoleOwner {
		t.Errorf("expected role Owner, got %s", role)
	}

	// Test with non-existent user
	role, err = rbacSvc.GetUserRole(context.Background(), 999)
	if err != nil {
		t.Fatalf("GetUserRole failed: %v", err)
	}
	if role != "" {
		t.Error("expected empty role for non-existent user")
	}
}

// TestRBACService_HasRole_NonExistentUser tests role check for non-existent user
func TestRBACService_HasRole_NonExistentUser(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	rbacSvc := NewRBACService(mockUserRepo)

	hasRole, err := rbacSvc.HasRole(context.Background(), 999, domain.RoleOwner)
	if err != nil {
		t.Fatalf("HasRole failed: %v", err)
	}
	if hasRole {
		t.Error("expected false for non-existent user")
	}
}

// TestRBACService_HasAnyRole_NonExistentUser tests any role check for non-existent user
func TestRBACService_HasAnyRole_NonExistentUser(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	rbacSvc := NewRBACService(mockUserRepo)

	hasRole, err := rbacSvc.HasAnyRole(context.Background(), 999, domain.RoleOwner, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("HasAnyRole failed: %v", err)
	}
	if hasRole {
		t.Error("expected false for non-existent user")
	}
}

// TestRBACService_AllRoles tests all role types
func TestRBACService_AllRoles(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockUserRepo.users[1] = &domain.User{
		ID:        1,
		Name:      "Owner",
		Email:     "owner@example.com",
		CompanyID: 123,
		Role:      domain.RoleOwner,
		Active:    true,
	}
	mockUserRepo.users[2] = &domain.User{
		ID:        2,
		Name:      "Admin",
		Email:     "admin@example.com",
		CompanyID: 123,
		Role:      domain.RoleAdmin,
		Active:    true,
	}
	mockUserRepo.users[3] = &domain.User{
		ID:        3,
		Name:      "Manager",
		Email:     "manager@example.com",
		CompanyID: 123,
		Role:      domain.RoleManager,
		Active:    true,
	}
	mockUserRepo.users[4] = &domain.User{
		ID:        4,
		Name:      "Employee",
		Email:     "employee@example.com",
		CompanyID: 123,
		Role:      domain.RoleEmployee,
		Active:    true,
	}
	rbacSvc := NewRBACService(mockUserRepo)

	// Test Owner
	isOwner, _ := rbacSvc.IsOwner(context.Background(), 1)
	if !isOwner {
		t.Error("Expected user 1 to be Owner")
	}

	// Test Admin
	isAdmin, _ := rbacSvc.IsAdmin(context.Background(), 2)
	if !isAdmin {
		t.Error("Expected user 2 to be Admin")
	}

	// Test Manager
	isManager, _ := rbacSvc.IsManager(context.Background(), 3)
	if !isManager {
		t.Error("Expected user 3 to be Manager")
	}

	// Test Employee
	canManageOrders, _ := rbacSvc.CanManageOrders(context.Background(), 4)
	if !canManageOrders {
		t.Error("Expected Employee to manage orders")
	}
	canManageProducts, _ := rbacSvc.CanManageProducts(context.Background(), 4)
	if canManageProducts {
		t.Error("Expected Employee not to manage products")
	}
}
