package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockRBACService is a mock implementation of RBACServiceInterface
type MockRBACService struct {
	canManageUsers      bool
	canAlterOwnerRole   bool
	canAlterAdminRole   bool
	manageUsersError    error
	alterOwnerRoleError error
	alterAdminRoleError error
}

func NewMockRBACService() *MockRBACService {
	return &MockRBACService{
		canManageUsers:    true,
		canAlterOwnerRole: true,
		canAlterAdminRole: true,
	}
}

func (m *MockRBACService) CanManageUsers(ctx context.Context, userID uint) (bool, error) {
	if m.manageUsersError != nil {
		return false, m.manageUsersError
	}
	return m.canManageUsers, nil
}

func (m *MockRBACService) CanAlterOwnerRole(ctx context.Context, userID uint) (bool, error) {
	if m.alterOwnerRoleError != nil {
		return false, m.alterOwnerRoleError
	}
	return m.canAlterOwnerRole, nil
}

func (m *MockRBACService) CanAlterAdminRole(ctx context.Context, userID uint) (bool, error) {
	if m.alterAdminRoleError != nil {
		return false, m.alterAdminRoleError
	}
	return m.canAlterAdminRole, nil
}

func (m *MockRBACService) CanApproveStockAdjustments(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (m *MockRBACService) HasRole(ctx context.Context, userID uint, role domain.Role) (bool, error) {
	return true, nil
}

func (m *MockRBACService) HasAnyRole(ctx context.Context, userID uint, roles ...domain.Role) (bool, error) {
	return true, nil
}

func (m *MockRBACService) IsOwner(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (m *MockRBACService) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (m *MockRBACService) IsManager(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (m *MockRBACService) CanManageCompany(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (m *MockRBACService) CanManageProducts(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (m *MockRBACService) CanManageOrders(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func TestUserManagementService_ListUsers(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	// Setup users from different companies
	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "User 1", Email: "user1@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "User 2", Email: "user2@example.com", CompanyID: 100, Role: domain.RoleAdmin, Active: true}
	mockUserRepo.users[3] = &domain.User{ID: 3, Name: "User 3", Email: "user3@example.com", CompanyID: 200, Role: domain.RoleManager, Active: true}

	users, err := svc.ListUsers(context.Background(), 100)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users for company 100, got %d", len(users))
	}
}

func TestUserManagementService_GetUser(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Test User", Email: "test@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	user, err := svc.GetUser(context.Background(), 100, 1)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if user.Name != "Test User" {
		t.Errorf("expected name 'Test User', got '%s'", user.Name)
	}
	if user.Role == nil || *user.Role != "owner" {
		t.Error("expected role 'owner'")
	}
}

func TestUserManagementService_GetUser_NotFound(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	_, err := svc.GetUser(context.Background(), 100, 999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
	if !errors.Is(err, ErrUserNotFound) && err.Error() != ErrUserNotFound.Error() {
		t.Errorf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserManagementService_GetUser_WrongCompany(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Test User", Email: "test@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	_, err := svc.GetUser(context.Background(), 200, 1)
	if err == nil {
		t.Error("expected error for user from different company")
	}
	if !errors.Is(err, ErrUserNotInCompany) && err.Error() != ErrUserNotInCompany.Error() {
		t.Errorf("expected ErrUserNotInCompany, got: %v", err)
	}
}

func TestUserManagementService_ChangeRole_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 2, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("ChangeRole failed: %v", err)
	}
	if mockUserRepo.users[2].Role != domain.RoleAdmin {
		t.Error("expected role to be updated to Admin")
	}
}

func TestUserManagementService_ChangeRole_AlterOwner(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 2, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("ChangeRole failed: %v", err)
	}
	if mockUserRepo.users[2].Role != domain.RoleAdmin {
		t.Error("expected role to be updated")
	}
}

func TestUserManagementService_ChangeRole_CannotAlterOwner(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	mockRBAC.canAlterOwnerRole = false
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleAdmin, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 2, domain.RoleAdmin)
	if err == nil {
		t.Error("expected error when non-Owner tries to alter Owner role")
	}
	if !errors.Is(err, ErrCannotAlterOwner) && err.Error() != ErrCannotAlterOwner.Error() {
		t.Errorf("expected ErrCannotAlterOwner, got: %v", err)
	}
}

func TestUserManagementService_ChangeRole_CannotAlterAdmin(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	mockRBAC.canAlterAdminRole = false
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleAdmin, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 2, domain.RoleManager)
	if err == nil {
		t.Error("expected error when non-Owner tries to alter Admin role")
	}
	if !errors.Is(err, ErrCannotAlterAdmin) && err.Error() != ErrCannotAlterAdmin.Error() {
		t.Errorf("expected ErrCannotAlterAdmin, got: %v", err)
	}
}

func TestUserManagementService_ChangeRole_PermissionDenied(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	mockRBAC.canManageUsers = false
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleEmployee, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 2, domain.RoleAdmin)
	if err == nil {
		t.Error("expected error when user lacks permission")
	}
	if !errors.Is(err, ErrPermissionDenied) && err.Error() != ErrPermissionDenied.Error() {
		t.Errorf("expected ErrPermissionDenied, got: %v", err)
	}
}

func TestUserManagementService_ChangeRole_TargetNotFound(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 999, domain.RoleAdmin)
	if err == nil {
		t.Error("expected error for non-existent target user")
	}
	if !errors.Is(err, ErrUserNotFound) && err.Error() != ErrUserNotFound.Error() {
		t.Errorf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserManagementService_ChangeRole_DifferentCompany(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 200, Role: domain.RoleManager, Active: true}

	err := svc.ChangeRole(context.Background(), 1, 2, domain.RoleAdmin)
	if err == nil {
		t.Error("expected error for user from different company")
	}
	if !errors.Is(err, ErrUserNotInCompany) && err.Error() != ErrUserNotInCompany.Error() {
		t.Errorf("expected ErrUserNotInCompany, got: %v", err)
	}
}

func TestUserManagementService_RemoveFromCompany_NotAllowed(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}

	err := svc.RemoveFromCompany(context.Background(), 1, 2)
	if err == nil {
		t.Error("expected error - RemoveFromCompany not allowed in Sprint 3")
	}
}

func TestUserManagementService_AddExistingUser_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "New User", Email: "new@example.com", CompanyID: 0, Role: "", Active: true}

	user, err := svc.AddExistingUser(context.Background(), 1, "new@example.com")
	if err != nil {
		t.Fatalf("AddExistingUser failed: %v", err)
	}
	if user.CompanyID == nil || *user.CompanyID != 100 {
		t.Error("expected user to be added to company 100")
	}
	if user.Role == nil || *user.Role != "manager" {
		t.Error("expected default role 'manager'")
	}
}

func TestUserManagementService_AddExistingUser_AlreadyInCompany(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Existing", Email: "existing@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}

	_, err := svc.AddExistingUser(context.Background(), 1, "existing@example.com")
	if err == nil {
		t.Error("expected error when user already in company")
	}
	if !errors.Is(err, ErrUserAlreadyInCompany) && err.Error() != ErrUserAlreadyInCompany.Error() {
		t.Errorf("expected ErrUserAlreadyInCompany, got: %v", err)
	}
}

func TestUserManagementService_AddExistingUser_UserNotFound(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	_, err := svc.AddExistingUser(context.Background(), 1, "nonexistent@example.com")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
	if !errors.Is(err, ErrUserNotFound) && err.Error() != ErrUserNotFound.Error() {
		t.Errorf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserManagementService_AddExistingUser_UserInOtherCompany(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Other", Email: "other@example.com", CompanyID: 200, Role: domain.RoleManager, Active: true}

	_, err := svc.AddExistingUser(context.Background(), 1, "other@example.com")
	if err == nil {
		t.Error("expected error when user belongs to another company")
	}
}

func TestUserManagementService_AddExistingUser_PermissionDenied(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	mockRBAC.canManageUsers = false
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleEmployee, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "New", Email: "new@example.com", CompanyID: 0, Role: "", Active: true}

	_, err := svc.AddExistingUser(context.Background(), 1, "new@example.com")
	if err == nil {
		t.Error("expected error when user lacks permission")
	}
	if !errors.Is(err, ErrPermissionDenied) && err.Error() != ErrPermissionDenied.Error() {
		t.Errorf("expected ErrPermissionDenied, got: %v", err)
	}
}

func TestUserManagementService_SetUserActive_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleManager, Active: false}

	err := svc.SetUserActive(context.Background(), 1, 2, true)
	if err != nil {
		t.Fatalf("SetUserActive failed: %v", err)
	}
	if !mockUserRepo.users[2].Active {
		t.Error("expected user to be activated")
	}
}

func TestUserManagementService_SetUserActive_DeactivateSuccess(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}

	err := svc.SetUserActive(context.Background(), 1, 2, false)
	if err != nil {
		t.Fatalf("SetUserActive failed: %v", err)
	}
	if mockUserRepo.users[2].Active {
		t.Error("expected user to be deactivated")
	}
}

func TestUserManagementService_SetUserActive_CannotDeactivateSelf(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	err := svc.SetUserActive(context.Background(), 1, 1, false)
	if err == nil {
		t.Error("expected error when trying to deactivate self")
	}
	if !errors.Is(err, ErrCannotDeactivateSelf) && err.Error() != ErrCannotDeactivateSelf.Error() {
		t.Errorf("expected ErrCannotDeactivateSelf, got: %v", err)
	}
}

func TestUserManagementService_SetUserActive_CannotDeactivateOwner(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Owner", Email: "owner@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	err := svc.SetUserActive(context.Background(), 1, 2, false)
	if err == nil {
		t.Error("expected error when trying to deactivate Owner")
	}
}

func TestUserManagementService_SetUserActive_PermissionDenied(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	mockRBAC.canManageUsers = false
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleEmployee, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 100, Role: domain.RoleManager, Active: true}

	err := svc.SetUserActive(context.Background(), 1, 2, false)
	if err == nil {
		t.Error("expected error when user lacks permission")
	}
	if !errors.Is(err, ErrPermissionDenied) && err.Error() != ErrPermissionDenied.Error() {
		t.Errorf("expected ErrPermissionDenied, got: %v", err)
	}
}

func TestUserManagementService_SetUserActive_TargetNotFound(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}

	err := svc.SetUserActive(context.Background(), 1, 999, true)
	if err == nil {
		t.Error("expected error for non-existent target user")
	}
	if !errors.Is(err, ErrUserNotFound) && err.Error() != ErrUserNotFound.Error() {
		t.Errorf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserManagementService_SetUserActive_DifferentCompany(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockRBAC := NewMockRBACService()
	svc := NewUserManagementService(mockUserRepo, mockCompanyRepo, mockRBAC)

	mockUserRepo.users[1] = &domain.User{ID: 1, Name: "Actor", Email: "actor@example.com", CompanyID: 100, Role: domain.RoleOwner, Active: true}
	mockUserRepo.users[2] = &domain.User{ID: 2, Name: "Target", Email: "target@example.com", CompanyID: 200, Role: domain.RoleManager, Active: true}

	err := svc.SetUserActive(context.Background(), 1, 2, false)
	if err == nil {
		t.Error("expected error for user from different company")
	}
	if !errors.Is(err, ErrUserNotInCompany) && err.Error() != ErrUserNotInCompany.Error() {
		t.Errorf("expected ErrUserNotInCompany, got: %v", err)
	}
}
