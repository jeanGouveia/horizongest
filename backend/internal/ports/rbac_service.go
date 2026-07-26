package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// RBACService defines the interface for role-based access control operations
type RBACService interface {
	HasRole(ctx context.Context, userID uint, role domain.Role) (bool, error)
	HasAnyRole(ctx context.Context, userID uint, roles ...domain.Role) (bool, error)
	IsOwner(ctx context.Context, userID uint) (bool, error)
	IsAdmin(ctx context.Context, userID uint) (bool, error)
	IsManager(ctx context.Context, userID uint) (bool, error)
	CanManageCompany(ctx context.Context, userID uint) (bool, error)
	CanManageProducts(ctx context.Context, userID uint) (bool, error)
	CanManageOrders(ctx context.Context, userID uint) (bool, error)
	CanManageUsers(ctx context.Context, userID uint) (bool, error)
	CanApproveStockAdjustments(ctx context.Context, userID uint) (bool, error)
}
