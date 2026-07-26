package service

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// ProductServiceInterface defines the interface for product and ingredient operations
type ProductServiceInterface interface {
	// Produto
	CreateProduct(ctx context.Context, in CreateProductInput) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
	ListActiveProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	UpdateProduct(ctx context.Context, id uint, in UpdateProductInput) (*domain.Product, error)
	DuplicateProduct(ctx context.Context, id uint) (*domain.Product, error)
	ArchiveProduct(ctx context.Context, id uint) error

	// Ingrediente
	CreateIngredient(ctx context.Context, in CreateIngredientInput) (*domain.Ingredient, error)
	ListIngredients(ctx context.Context) ([]domain.Ingredient, error)
	GetIngredient(ctx context.Context, id uint) (*domain.Ingredient, error)
	UpdateIngredientStock(ctx context.Context, id uint, in UpdateStockInput) (*domain.Ingredient, error)
	UpdateIngredient(ctx context.Context, id uint, in UpdateIngredientInput) (*domain.Ingredient, error)
	DeleteIngredient(ctx context.Context, id uint) error

	// Ficha técnica
	SetProductIngredients(ctx context.Context, productID uint, in SetProductIngredientsInput) error
}

// OrderServiceInterface defines the interface for order operations
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, in CreateOrderInput) (*domain.Order, error)
	ListOrders(ctx context.Context) ([]domain.Order, error)
	GetOrder(ctx context.Context, id uint) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, in UpdateOrderStatusInput) (*domain.Order, error)
	UpdateOrder(ctx context.Context, id uint, in UpdateOrderInput) (*domain.Order, error)
}

// RBACServiceInterface defines the interface for role-based access control operations
type RBACServiceInterface interface {
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
	CanAlterOwnerRole(ctx context.Context, userID uint) (bool, error)
	CanAlterAdminRole(ctx context.Context, userID uint) (bool, error)
}

// AuthServiceInterface defines the interface for authentication operations
type AuthServiceInterface interface {
	Login(ctx context.Context, input LoginInput) (*LoginResult, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*JWTClaims, error)
	UpdateProfile(ctx context.Context, userID uint, input UpdateProfileInput) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uint, input ChangePasswordInput) error
	RequestPasswordReset(ctx context.Context, input RequestPasswordResetInput) error
	ResetPassword(ctx context.Context, input ResetPasswordInput) error
}

// UserManagementServiceInterface defines the interface for user management operations
type UserManagementServiceInterface interface {
	ListUsers(ctx context.Context, companyID uint) ([]UserOutput, error)
	GetUser(ctx context.Context, companyID uint, userID uint) (*UserOutput, error)
	ChangeRole(ctx context.Context, actorUserID uint, targetUserID uint, newRole domain.Role) error
	RemoveFromCompany(ctx context.Context, actorUserID uint, targetUserID uint) error
	AddExistingUser(ctx context.Context, actorUserID uint, email string) (*UserOutput, error)
	SetUserActive(ctx context.Context, actorUserID uint, targetUserID uint, active bool) error
}
