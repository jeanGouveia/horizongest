# Sprint 3.2 - Project Audit Report

**Date**: 2026-07-20  
**Objective**: Complete audit of existing architecture before Sprint 3.2 implementation

## Backend Structure

### Domain Entities (28 files)

#### Platform Entities (Created in Sprint 3.1)
- `platform_user.go` - Platform user with roles (PlatformAdmin, PlatformSupport)
- `platform_session.go` - Platform session management
- `platform_audit.go` - Audit logging for platform actions
- `platform_role.go` - Platform role enumeration

#### Tenant Entities (Updated in Sprint 3.1)
- `user.go` - User with non-nullable CompanyID and Role
- `tenant_context.go` - Tenant context with non-nullable CompanyID, removed IsSystemAdmin
- `role.go` - Tenant roles (Owner, Admin, Manager, Staff)
- `company.go` - Company entity
- `category.go` - Category with non-nullable CompanyID
- `product.go` - Product with non-nullable CompanyID
- `order.go` - Order with non-nullable CompanyID
- `ingredient.go` - Ingredient with non-nullable CompanyID
- `media.go` - Media with non-nullable CompanyID
- `stock_adjustment_pending.go` - Stock adjustment with non-nullable CompanyID

#### Other Domain Entities
- `business_profile.go`
- `business_type.go`
- `dashboard.go`
- `dependency.go`
- `error_response.go`
- `invitation.go`
- `notifications.go`
- `order_item.go`
- `password_reset_token.go`
- `product_ingredient.go`
- `stock_validation.go`
- `system.go`
- `theme.go`
- `token_blacklist.go`

### Services (15 files)

#### Platform Services (Created in Sprint 3.1 - NOT WIRED)
- `platform_auth_service.go` - Platform authentication with JWT
- `platform_service.go` - Platform company management

#### Tenant Services
- `auth_service.go` - Tenant authentication (Register method removed)
- `business_service.go`
- `category_service.go`
- `company_service.go` - Company management (CreateCompany exists - needs review)
- `company_settings_service.go`
- `invitation_service.go`
- `media_service.go`
- `order_service.go`
- `product_service.go`
- `rbac_service.go`
- `stock_adjustment_service.go`
- `theme_service.go`
- `user_management_service.go`

### Handlers (16 files)

#### Platform Handlers (Created in Sprint 3.1 - NOT WIRED)
- `platform_auth_handler.go` - Platform login, logout, me
- `platform_company_handler.go` - Platform company management

#### Tenant Handlers
- `auth_handler.go` - Tenant auth (Register handler removed)
- `business_handler.go`
- `category_handler.go`
- `company_handler.go` - Company management (CreateCompany exists - needs review)
- `company_settings_handler.go`
- `dashboard_handler.go`
- `invitation_handler.go`
- `media_handler.go`
- `order_handler.go`
- `product_handler.go`
- `stock_adjustment_handler.go`
- `system_handler.go`
- `theme_handler.go`
- `user_management_handler.go`

### Repositories (13 files)

#### Platform Repositories (NOT CREATED)
- GORM platform user repository - MISSING
- GORM platform session repository - MISSING
- GORM platform audit repository - MISSING

#### Tenant Repositories
- `gorm_category_repository.go`
- `gorm_company_repository.go`
- `gorm_dashboard_repository.go`
- `gorm_invitation_repository.go`
- `gorm_media_repository.go`
- `gorm_notifications_repository.go`
- `gorm_order_repository.go`
- `gorm_password_reset_repository.go`
- `gorm_product_repository.go`
- `gorm_stock_adjustment_repository.go`
- `gorm_token_blacklist_repository.go`
- `gorm_user_repository.go`
- `tenant_helper.go`

### Middleware (5 files)

#### Platform Middleware (Created in Sprint 3.1 - NOT WIRED)
- `platform_auth_middleware.go` - Platform JWT authentication

#### Tenant Middleware
- `auth_middleware.go` - Tenant JWT authentication
- `error_middleware.go`
- `role_middleware.go` - RBAC enforcement
- `tenant_middleware.go` - Tenant context

### Database (2 files)
- `connection.go`
- `migrate.go`

### Main Entry Point
- `cmd/server/main.go`

## Frontend Structure

### Svelte Components (47 files)

#### Layout Components (4)
- `Footer.svelte`
- `Header.svelte`
- `Sidebar.svelte`
- `Workspace.svelte`

#### UI Components (22)
- `Alert.svelte`
- `Badge.svelte`
- `Button.svelte`
- `Card.svelte`
- `Checkbox.svelte`
- `ConfirmDialog.svelte`
- `Divider.svelte`
- `EmptyState.svelte`
- `Input.svelte`
- `Loading.svelte`
- `Modal.svelte`
- `PageContainer.svelte`
- `PageHeader.svelte`
- `PhotoUpload.svelte`
- `ProductCard.svelte`
- `Section.svelte`
- `Select.svelte`
- `Skeleton.svelte`
- `TabNavigation.svelte`
- `Table.svelte`
- `Textarea.svelte`
- `Toast.svelte`

#### App Routes (15)
- `(app)/+layout.svelte`
- `(app)/dashboard/+page.svelte`
- `(app)/categories/+page.svelte`
- `(app)/ingredients/+page.svelte`
- `(app)/orders/+page.svelte`
- `(app)/orders/[id]/+page.svelte`
- `(app)/orders/new/+page.svelte`
- `(app)/products/+page.svelte`
- `(app)/products/[id]/+page.svelte`
- `(app)/products/[id]/edit/+page.svelte`
- `(app)/products/new/+page.svelte`
- `(app)/profile/+page.svelte`
- `(app)/settings/company/+page.svelte`
- `(app)/settings/invitations/+page.svelte`
- `(app)/settings/users/+page.svelte`
- `(app)/stock-adjustments/+page.svelte`

#### Auth Routes (3)
- `(auth)/forgot-password/+page.svelte`
- `(auth)/login/+page.svelte`
- `(auth)/reset-password/+page.svelte`

#### Other Routes (1)
- `invite/[token]/+page.svelte`

#### Root Layout
- `+layout.svelte`

### TypeScript Files (34 files)

#### API Clients (6)
- `category.ts`
- `client.ts`
- `media.ts`
- `order.ts`
- `product.ts`
- `stock-adjustment.ts`

#### Stores (4)
- `rbacStore.svelte.ts`
- `themeStore.svelte.ts`
- `toast.ts`
- `userStore.svelte.ts`

#### Types (12)
- `category.ts`
- `dashboard.ts`
- `dependency.ts`
- `ingredient.ts`
- `media.ts`
- `notifications.ts`
- `order.ts`
- `product.ts`
- `stock-adjustment.ts`
- `stock-validation.ts`
- `system.ts`
- `user.ts`

#### Theme (8)
- `animations.ts`
- `colors.ts`
- `dark-mode.ts`
- `radius.ts`
- `shadow.ts`
- `spacing.ts`
- `transitions.ts`
- `typography.ts`

#### Hooks (1)
- `useSystem.ts`

#### Services (1)
- `errorService.ts`

#### Component Indexes (2)
- `components/layout/index.ts`
- `components/ui/index.ts`

## Current Routes in main.go

### Public Routes
- `GET /api/health` - Health check
- `GET /api/system/*` - System endpoints
- `POST /api/auth/login` - Tenant login
- `POST /api/auth/logout` - Tenant logout
- `POST /api/auth/request-password-reset` - Password reset request
- `POST /api/auth/reset-password` - Password reset
- `GET /uploads/*` - Static file serving

### Protected Routes (Tenant Auth + Tenant Middleware)
- `GET /api/dashboard` - Dashboard
- `GET /api/me` - Current user info
- `PUT /api/me` - Update profile
- `POST /api/me/change-password` - Change password

### Company Routes (Tenant Level - NEEDS REVIEW)
- `POST /api/companies` - Create company (SHOULD BE REMOVED - platform only)
- `GET /api/companies` - List companies
- `GET /api/companies/{id}` - Get company
- `PUT /api/companies/{id}` - Update company
- `DELETE /api/companies/{id}` - Delete company

### Company Settings
- `GET /api/company/settings` - Get settings
- `PUT /api/company/settings` - Update settings

### User Management (RBAC: Owner, Admin)
- `GET /api/company/users` - List users
- `GET /api/company/users/{id}` - Get user
- `POST /api/company/users/add` - Add user
- `PUT /api/company/users/{id}/role` - Change role
- `PUT /api/company/users/{id}/active` - Set active status
- `DELETE /api/company/users/{id}` - Remove user

### Invitations (RBAC: Owner, Admin)
- `POST /api/company/invitations` - Create invitation
- `GET /api/company/invitations` - List invitations
- `GET /api/company/invitations/{id}` - Get invitation
- `DELETE /api/company/invitations/{id}` - Revoke invitation
- `POST /api/invitations/accept` - Accept invitation
- `GET /api/invitations/{token}` - Get invitation by token

### Theme
- `GET /api/theme` - Get theme
- `GET /api/theme/default` - Get default theme

### Business Profile
- `GET /api/business/profile` - Get business profile
- `GET /api/business/profile/default` - Get default business profile

### Products
- `POST /api/products` - Create product
- `GET /api/products` - List products
- `GET /api/products/active` - List active products
- `GET /api/products/{id}` - Get product
- `PUT /api/products/{id}` - Update product
- `DELETE /api/products/{id}` - Delete product
- `POST /api/products/{id}/duplicate` - Duplicate product
- `POST /api/products/{id}/archive` - Archive product
- `PUT /api/products/{id}/ingredients` - Set product ingredients
- `GET /api/products/{id}/ingredients` - Get product ingredients

### Ingredients
- `POST /api/ingredients` - Create ingredient
- `GET /api/ingredients` - List ingredients
- `GET /api/ingredients/{id}` - Get ingredient
- `PUT /api/ingredients/{id}` - Update ingredient
- `DELETE /api/ingredients/{id}` - Delete ingredient
- `PATCH /api/ingredients/{id}/stock` - Update ingredient stock

### Categories
- `POST /api/categories` - Create category
- `GET /api/categories` - List categories
- `GET /api/categories/{id}` - Get category
- `PUT /api/categories/{id}` - Update category
- `DELETE /api/categories/{id}` - Delete category

### Orders
- `POST /api/orders` - Create order
- `GET /api/orders` - List orders
- `GET /api/orders/{id}` - Get order
- `PUT /api/orders/{id}` - Update order
- `PATCH /api/orders/{id}/status` - Update order status

### Stock Adjustments
- `GET /api/stock-adjustments/pending` - List pending adjustments
- `POST /api/stock-adjustments/{id}/approve` - Approve adjustment
- `POST /api/stock-adjustments/{id}/reject` - Reject adjustment

### Media
- `POST /api/media/upload` - Upload media
- `GET /api/media/{id}` - Get media
- `DELETE /api/media/{id}` - Delete media
- `GET /api/media/entity/{entity_type}/{entity_id}` - Get media by entity

## Platform Routes (NOT WIRED - Created in Sprint 3.1)

### Platform Auth
- `POST /api/platform/auth/login` - Platform login
- `POST /api/platform/auth/logout` - Platform logout
- `GET /api/platform/auth/me` - Get current platform user

### Platform Companies
- `POST /api/platform/companies` - Create company (PlatformAdmin only)
- `GET /api/platform/companies` - List companies (PlatformAdmin only)
- `GET /api/platform/companies/:id` - Get company (PlatformAdmin only)
- `PUT /api/platform/companies/:id` - Update company (PlatformAdmin only)
- `POST /api/platform/companies/:id/deactivate` - Deactivate company (PlatformAdmin only)
- `POST /api/platform/companies/:id/activate` - Activate company (PlatformAdmin only)

## Issues Identified

### Critical Issues
1. **Platform repositories not created** - GORM repositories for PlatformUser, PlatformSession, PlatformAudit are missing
2. **Platform services not instantiated** - platform_auth_service and platform_service are not created in main.go
3. **Platform handlers not wired** - platform_auth_handler and platform_company_handler are not registered in main.go
4. **Platform middleware not used** - platform_auth_middleware is not used in main.go
5. **Platform routes not defined** - No platform routes are defined in main.go
6. **Tenant company creation endpoint exists** - `POST /api/companies` should be removed (platform only)

### Self-Service Code to Remove
1. `POST /api/companies` - Tenant-level company creation (platform only)
2. Any remaining references to public registration
3. Auto company creation logic (if any exists)

### Missing Platform Features
1. Platform user repository (GORM)
2. Platform session repository (GORM)
3. Platform audit repository (GORM)
4. Platform dashboard frontend
5. Platform login frontend
6. Platform companies management frontend
7. Platform owner management frontend
8. Platform plans structure
9. Platform suspension logic
10. Platform backup functionality
11. Platform export functionality
12. Platform audit viewer frontend
13. Login as company (support) feature

## Next Steps

1. Create platform repositories (GORM)
2. Wire platform services in main.go
3. Wire platform handlers in main.go
4. Wire platform middleware in main.go
5. Define platform routes in main.go
6. Remove tenant company creation endpoint
7. Create platform frontend screens
8. Implement remaining platform features
