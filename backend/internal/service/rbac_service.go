package service

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

// RBACService centralizes all role-based access control logic.
// All permission checks should go through this service to maintain consistency.
type RBACService struct {
	userRepo ports.UserRepository
}

func NewRBACService(userRepo ports.UserRepository) *RBACService {
	return &RBACService{userRepo: userRepo}
}

// HasRole checks if a user has a specific role
func (s *RBACService) HasRole(ctx context.Context, userID uint, role domain.Role) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return user.Role == role, nil
}

// HasAnyRole checks if a user has any of the specified roles
func (s *RBACService) HasAnyRole(ctx context.Context, userID uint, roles ...domain.Role) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	for _, role := range roles {
		if user.Role == role {
			return true, nil
		}
	}
	return false, nil
}

// IsOwner checks if user is an Owner
func (s *RBACService) IsOwner(ctx context.Context, userID uint) (bool, error) {
	return s.HasRole(ctx, userID, domain.RoleOwner)
}

// IsAdmin checks if user is an Admin
func (s *RBACService) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	return s.HasRole(ctx, userID, domain.RoleAdmin)
}

// IsManager checks if user is a Manager
func (s *RBACService) IsManager(ctx context.Context, userID uint) (bool, error) {
	return s.HasRole(ctx, userID, domain.RoleManager)
}

// CanManageCompany checks if user can manage company settings
// Owner and Admin can manage company settings
func (s *RBACService) CanManageCompany(ctx context.Context, userID uint) (bool, error) {
	return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin)
}

// CanManageProducts checks if user can manage products
// Owner, Admin, and Manager can manage products
func (s *RBACService) CanManageProducts(ctx context.Context, userID uint) (bool, error) {
	return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin, domain.RoleManager)
}

// CanManageOrders checks if user can manage orders
// Owner, Admin, Manager, and Employee can manage orders
func (s *RBACService) CanManageOrders(ctx context.Context, userID uint) (bool, error) {
	return s.HasAnyRole(ctx, userID,
		domain.RoleOwner,
		domain.RoleAdmin,
		domain.RoleManager,
		domain.RoleEmployee,
	)
}

// CanManageUsers checks if user can manage other users
// Only Owner can manage users (including changing roles)
func (s *RBACService) CanManageUsers(ctx context.Context, userID uint) (bool, error) {
	return s.HasRole(ctx, userID, domain.RoleOwner)
}

// CanManageSettings checks if user can manage settings
// Owner and Admin can manage settings
func (s *RBACService) CanManageSettings(ctx context.Context, userID uint) (bool, error) {
	return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin)
}

// CanViewReports checks if user can view reports
// Owner, Admin, and Manager can view reports
func (s *RBACService) CanViewReports(ctx context.Context, userID uint) (bool, error) {
	return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin, domain.RoleManager)
}

// CanAlterOwnerRole checks if user can alter Owner role
// Only Owner can alter Owner role (Admin cannot)
func (s *RBACService) CanAlterOwnerRole(ctx context.Context, userID uint) (bool, error) {
	return s.HasRole(ctx, userID, domain.RoleOwner)
}

// CanAlterAdminRole checks if user can alter Admin role
// Only Owner can alter Admin role
func (s *RBACService) CanAlterAdminRole(ctx context.Context, userID uint) (bool, error) {
	return s.HasRole(ctx, userID, domain.RoleOwner)
}

// GetUserRole returns the user's role
func (s *RBACService) GetUserRole(ctx context.Context, userID uint) (domain.Role, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}
	return user.Role, nil
}

// CanApproveStockAdjustments checks if user can approve stock adjustments
// Owner and Admin can approve stock adjustments
func (s *RBACService) CanApproveStockAdjustments(ctx context.Context, userID uint) (bool, error) {
	return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin)
}
