# ETAPA 1: Mapeamento Completo dos Fluxos do Sistema
**Auditoria Forense - pratoOnline**
**Data:** 19/07/2026

---

## Árvore Completa de Fluxos

### 1. AUTENTICAÇÃO

#### 1.1 Registro (Register)
**Frontend:** `/frontend/src/routes/(auth)/register/+page.svelte`
**Backend:** `POST /api/auth/register`
**Handler:** `backend/internal/handler/auth_handler.go:Register (linha 24-56)`
**Service:** `backend/internal/service/auth_service.go:Register (linha 68-120)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:Create (linha 40-58)`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:Create (linha 40-58)`
**Model:** `backend/internal/domain/user.go`, `backend/internal/domain/company.go`
**SQL:** INSERT INTO users, INSERT INTO companies
**Middleware:** Nenhum (rota pública)

#### 1.2 Login
**Frontend:** `/frontend/src/routes/(auth)/login/+page.svelte`
**Backend:** `POST /api/auth/login`
**Handler:** `backend/internal/handler/auth_handler.go:Login (linha 58-97)`
**Service:** `backend/internal/service/auth_service.go:Login (linha 122-145)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByEmail (linha 60-70)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE email = ?
**Middleware:** Nenhum (rota pública)

#### 1.3 Logout
**Backend:** `POST /api/auth/logout`
**Handler:** `backend/internal/handler/auth_handler.go:Logout (linha 99-115)`
**Service:** `backend/internal/service/auth_service.go:Logout (linha 255-261)`
**Repository:** Nenhum (blacklist in-memory)
**Model:** Nenhum
**SQL:** Nenhum
**Middleware:** `authMw.Auth` (linha 109)

#### 1.4 Perfil (/me)
**Backend:** `GET /api/me`, `PUT /api/me`
**Handler:** `backend/internal/handler/auth_handler.go:Me (linha 117-135)`, `UpdateProfile (linha 137-167)`
**Service:** `backend/internal/service/auth_service.go:N/A (direto no handler)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByID (linha 72-82)`, `Update (linha 84-103)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE id = ?, UPDATE users
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 1.5 Alterar Senha
**Backend:** `POST /api/me/change-password`
**Handler:** `backend/internal/handler/auth_handler.go:ChangePassword (linha 169-197)`
**Service:** `backend/internal/service/auth_service.go:ChangePassword (linha 147-174)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:Update (linha 84-103)`
**Model:** `backend/internal/domain/user.go`
**SQL:** UPDATE users
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 2. GESTÃO DE EMPRESA

#### 2.1 Criar Empresa
**Backend:** `POST /api/companies`
**Handler:** `backend/internal/handler/company_handler.go:CreateCompany`
**Service:** `backend/internal/service/company_service.go:CreateCompany`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:Create`
**Model:** `backend/internal/domain/company.go`
**SQL:** INSERT INTO companies
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 2.2 Listar Empresas
**Backend:** `GET /api/companies`
**Handler:** `backend/internal/handler/company_handler.go:ListCompanies`
**Service:** `backend/internal/service/company_service.go:ListCompanies`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:List`
**Model:** `backend/internal/domain/company.go`
**SQL:** SELECT FROM companies
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 2.3 Obter Empresa
**Backend:** `GET /api/companies/{id}`
**Handler:** `backend/internal/handler/company_handler.go:GetCompany`
**Service:** `backend/internal/service/company_service.go:GetCompany`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:FindByID`
**Model:** `backend/internal/domain/company.go`
**SQL:** SELECT FROM companies WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 2.4 Atualizar Empresa
**Backend:** `PUT /api/companies/{id}`
**Handler:** `backend/internal/handler/company_handler.go:UpdateCompany`
**Service:** `backend/internal/service/company_service.go:UpdateCompany`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:Update`
**Model:** `backend/internal/domain/company.go`
**SQL:** UPDATE companies
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 2.5 Deletar Empresa
**Backend:** `DELETE /api/companies/{id}`
**Handler:** `backend/internal/handler/company_handler.go:DeleteCompany`
**Service:** `backend/internal/service/company_service.go:DeleteCompany`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:Delete`
**Model:** `backend/internal/domain/company.go`
**SQL:** UPDATE companies SET deleted_at = ? (soft delete)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 3. CONFIGURAÇÕES DA EMPRESA

#### 3.1 Obter Configurações
**Frontend:** `/frontend/src/routes/(app)/settings/company/+page.svelte`
**Backend:** `GET /api/company/settings`
**Handler:** `backend/internal/handler/company_settings_handler.go:GetSettings`
**Service:** `backend/internal/service/company_settings_service.go:GetSettings`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:FindByID`
**Model:** `backend/internal/domain/company.go`
**SQL:** SELECT FROM companies WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 3.2 Atualizar Configurações
**Frontend:** `/frontend/src/routes/(app)/settings/company/+page.svelte`
**Backend:** `PUT /api/company/settings`
**Handler:** `backend/internal/handler/company_settings_handler.go:UpdateSettings`
**Service:** `backend/internal/service/company_settings_service.go:UpdateSettings`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:Update`
**Model:** `backend/internal/domain/company.go`
**SQL:** UPDATE companies
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 4. GESTÃO DE USUÁRIOS (Sprint 7)

#### 4.1 Listar Usuários da Empresa
**Frontend:** `/frontend/src/routes/(app)/settings/users/+page.svelte`
**Backend:** `GET /api/company/users`
**Handler:** `backend/internal/handler/user_management_handler.go:ListUsers (linha 42-63)`
**Service:** `backend/internal/service/user_management_service.go:ListUsers (linha 46-75)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:List (linha 105-120)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE deleted_at IS NULL
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 132)

#### 4.2 Obter Usuário
**Backend:** `GET /api/company/users/{id}`
**Handler:** `backend/internal/handler/user_management_handler.go:GetUser (linha 65-97)`
**Service:** `backend/internal/service/user_management_service.go:GetUser (linha 77-106)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByID (linha 72-82)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 132)

#### 4.3 Adicionar Usuário Existente
**Frontend:** `/frontend/src/routes/(app)/settings/users/+page.svelte`
**Backend:** `POST /api/company/users/add`
**Handler:** `backend/internal/handler/user_management_handler.go:AddUser (linha 99-135)`
**Service:** `backend/internal/service/user_management_service.go:AddExistingUser (linha 222-282)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByEmail (linha 60-70)`, `Update (linha 84-103)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE email = ?, UPDATE users
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 132)

#### 4.4 Alterar Cargo
**Frontend:** `/frontend/src/routes/(app)/settings/users/+page.svelte`
**Backend:** `PUT /api/company/users/{id}/role`
**Handler:** `backend/internal/handler/user_management_handler.go:ChangeRole (linha 137-185)`
**Service:** `backend/internal/service/user_management_service.go:ChangeRole (linha 108-171)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByID (linha 72-82)`, `Update (linha 84-103)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE id = ?, UPDATE users
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 132)

#### 4.5 Remover Usuário da Empresa
**Frontend:** `/frontend/src/routes/(app)/settings/users/+page.svelte`
**Backend:** `DELETE /api/company/users/{id}`
**Handler:** `backend/internal/handler/user_management_handler.go:RemoveUser (linha 187-217)`
**Service:** `backend/internal/service/user_management_service.go:RemoveFromCompany (linha 173-220)`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByID (linha 72-82)`, `Update (linha 84-103)`
**Model:** `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE id = ?, UPDATE users SET company_id = NULL
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 132)

---

### 5. CONVITES (Sprint 8)

#### 5.1 Criar Convite
**Frontend:** `/frontend/src/routes/(app)/settings/invitations/+page.svelte`
**Backend:** `POST /api/company/invitations`
**Handler:** `backend/internal/handler/invitation_handler.go:CreateInvitation (linha 42-96)`
**Service:** `backend/internal/service/invitation_service.go:CreateInvitation (linha 62-120)`
**Repository:** `backend/internal/infra/repository/gorm_invitation_repository.go:Create`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByEmail (linha 60-70)`
**Model:** `backend/internal/domain/invitation.go`, `backend/internal/domain/user.go`
**SQL:** SELECT FROM users WHERE email = ?, INSERT INTO invitations
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 144)

#### 5.2 Listar Convites
**Frontend:** `/frontend/src/routes/(app)/settings/invitations/+page.svelte`
**Backend:** `GET /api/company/invitations`
**Handler:** `backend/internal/handler/invitation_handler.go:ListInvitations (linha 98-119)`
**Service:** `backend/internal/service/invitation_service.go:ListInvitations (linha 122-134)`
**Repository:** `backend/internal/infra/repository/gorm_invitation_repository.go:ListByCompanyID`
**Model:** `backend/internal/domain/invitation.go`
**SQL:** SELECT FROM invitations WHERE company_id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 144)

#### 5.3 Obter Convite
**Backend:** `GET /api/company/invitations/{id}`
**Handler:** `backend/internal/handler/invitation_handler.go:GetInvitation (linha 121-153)`
**Service:** `backend/internal/service/invitation_service.go:GetInvitation (linha 136-152)`
**Repository:** `backend/internal/infra/repository/gorm_invitation_repository.go:FindByID`
**Model:** `backend/internal/domain/invitation.go`
**SQL:** SELECT FROM invitations WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 144)

#### 5.4 Revogar Convite
**Frontend:** `/frontend/src/routes/(app)/settings/invitations/+page.svelte`
**Backend:** `DELETE /api/company/invitations/{id}`
**Handler:** `backend/internal/handler/invitation_handler.go:RevokeInvitation (linha 155-190)`
**Service:** `backend/internal/service/invitation_service.go:RevokeInvitation (linha 154-192)`
**Repository:** `backend/internal/infra/repository/gorm_invitation_repository.go:FindByID`, `Update`
**Model:** `backend/internal/domain/invitation.go`
**SQL:** SELECT FROM invitations WHERE id = ?, UPDATE invitations SET status = 'revoked'
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110), `roleMw.RequireAny(Owner, Admin)` (linha 144)

#### 5.5 Obter Convite por Token (Público)
**Frontend:** `/frontend/src/routes/invite/[token]/+page.svelte`
**Backend:** `GET /api/invitations/{token}`
**Handler:** `backend/internal/handler/invitation_handler.go:GetInvitationByToken (linha 192-211)`
**Service:** `backend/internal/service/invitation_service.go:GetInvitationByToken (linha 194-214)`
**Repository:** `backend/internal/infra/repository/gorm_invitation_repository.go:FindByToken`
**Model:** `backend/internal/domain/invitation.go`
**SQL:** SELECT FROM invitations WHERE token = ?
**Middleware:** Nenhum (rota pública)

#### 5.6 Aceitar Convite
**Frontend:** `/frontend/src/routes/invite/[token]/+page.svelte`
**Backend:** `POST /api/invitations/accept`
**Handler:** `backend/internal/handler/invitation_handler.go:AcceptInvitation (linha 213-259)`
**Service:** `backend/internal/service/invitation_service.go:AcceptInvitation (linha 216-277)`
**Repository:** `backend/internal/infra/repository/gorm_invitation_repository.go:ListByCompanyID`, `Update`
**Repository:** `backend/internal/infra/repository/gorm_user_repository.go:FindByEmail`, `Update`
**Model:** `backend/internal/domain/invitation.go`, `backend/internal/domain/user.go`
**SQL:** SELECT FROM invitations WHERE token = ?, SELECT FROM users WHERE email = ?, UPDATE users, UPDATE invitations
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 6. PRODUTOS

#### 6.1 Criar Produto
**Frontend:** `/frontend/src/routes/(app)/products/new/+page.svelte`
**Backend:** `POST /api/products`
**Handler:** `backend/internal/handler/product_handler.go:CreateProduct (linha 24-41)`
**Service:** `backend/internal/service/product_service.go:CreateProduct (linha 164-200)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:CreateProduct (linha 96-155)`
**Model:** `backend/internal/domain/product.go`
**SQL:** INSERT INTO products
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 6.2 Listar Produtos
**Frontend:** `/frontend/src/routes/(app)/products/+page.svelte`
**Backend:** `GET /api/products`, `GET /api/products/active`
**Handler:** `backend/internal/handler/product_handler.go:ListProducts (linha 43-51)`, `ListActiveProducts (linha 53-61)`
**Service:** `backend/internal/service/product_service.go:ListProducts (linha 202-208)`, `ListActiveProducts (linha 210-216)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:ListProducts (linha 170-181)`, `ListActiveProducts (linha 183-194)`
**Model:** `backend/internal/domain/product.go`
**SQL:** SELECT FROM products WHERE deleted_at IS NULL
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 6.3 Obter Produto
**Frontend:** `/frontend/src/routes/(app)/products/[id]/+page.svelte`
**Backend:** `GET /api/products/{id}`
**Handler:** `backend/internal/handler/product_handler.go:GetProduct (linha 63-80)`
**Service:** `backend/internal/service/product_service.go:GetProduct (linha 218-233)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:FindProductByID (linha 157-168)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:GetProductIngredients (linha 439-456)`
**Model:** `backend/internal/domain/product.go`, `backend/internal/domain/product_ingredient.go`
**SQL:** SELECT FROM products WHERE id = ?, SELECT FROM product_ingredients
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 6.4 Atualizar Produto
**Frontend:** `/frontend/src/routes/(app)/products/[id]/edit/+page.svelte`
**Backend:** `PUT /api/products/{id}`
**Handler:** `backend/internal/handler/product_handler.go:UpdateProduct (linha 100-126)`
**Service:** `backend/internal/service/product_service.go:UpdateProduct (linha 246-300)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:UpdateProduct (linha 196-244)`
**Model:** `backend/internal/domain/product.go`
**SQL:** UPDATE products
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 6.5 Deletar Produto
**Frontend:** `/frontend/src/routes/(app)/products/+page.svelte`
**Backend:** `DELETE /api/products/{id}`
**Handler:** `backend/internal/handler/product_handler.go:DeleteProduct (linha 82-98)`
**Service:** `backend/internal/service/product_service.go:DeleteProduct (linha 235-244)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:DeleteProduct (linha 246-258)`
**Model:** `backend/internal/domain/product.go`
**SQL:** UPDATE products SET deleted_at = ? (soft delete)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 6.6 Definir Ingredientes do Produto (Ficha Técnica)
**Frontend:** `/frontend/src/routes/(app)/products/[id]/edit/+page.svelte`
**Backend:** `PUT /api/products/{id}/ingredients`
**Handler:** `backend/internal/handler/product_handler.go:SetProductIngredients (linha 254-279)`
**Service:** `backend/internal/service/product_service.go:SetProductIngredients (linha 387-415)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:SetProductIngredients (linha 428-437)`
**Model:** `backend/internal/domain/product_ingredient.go`
**SQL:** DELETE FROM product_ingredients, INSERT INTO product_ingredients
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 6.7 Obter Ingredientes do Produto
**Backend:** `GET /api/products/{id}/ingredients`
**Handler:** `backend/internal/handler/product_handler.go:GetProductIngredients (linha 281-298)`
**Service:** `backend/internal/service/product_service.go:GetProduct (linha 218-233)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:GetProductIngredients (linha 439-456)`
**Model:** `backend/internal/domain/product_ingredient.go`
**SQL:** SELECT FROM product_ingredients WHERE product_id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 7. INGREDIENTES

#### 7.1 Criar Ingrediente
**Frontend:** `/frontend/src/routes/(app)/ingredients/+page.svelte`
**Backend:** `POST /api/ingredients`
**Handler:** `backend/internal/handler/product_handler.go:CreateIngredient (linha 130-147)`
**Service:** `backend/internal/service/product_service.go:CreateIngredient (linha 304-313)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:CreateIngredient (linha 260-273)`
**Model:** `backend/internal/domain/ingredient.go`
**SQL:** INSERT INTO ingredients
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 7.2 Listar Ingredientes
**Frontend:** `/frontend/src/routes/(app)/ingredients/+page.svelte`
**Backend:** `GET /api/ingredients`
**Handler:** `backend/internal/handler/product_handler.go:ListIngredients (linha 149-157)`
**Service:** `backend/internal/service/product_service.go:ListIngredients (linha 315-317)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:ListIngredients (linha 275-286)`
**Model:** `backend/internal/domain/ingredient.go`
**SQL:** SELECT FROM ingredients WHERE deleted_at IS NULL
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 7.3 Obter Ingrediente
**Backend:** `GET /api/ingredients/{id}`
**Handler:** `backend/internal/handler/product_handler.go:GetIngredient (linha 159-176)`
**Service:** `backend/internal/service/product_service.go:GetIngredient (linha 319-328)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:FindIngredientByID (linha 288-299)`
**Model:** `backend/internal/domain/ingredient.go`
**SQL:** SELECT FROM ingredients WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 7.4 Atualizar Ingrediente
**Frontend:** `/frontend/src/routes/(app)/ingredients/+page.svelte`
**Backend:** `PUT /api/ingredients/{id}`
**Handler:** `backend/internal/handler/product_handler.go:UpdateIngredient (linha 206-232)`
**Service:** `backend/internal/service/product_service.go:UpdateIngredient (linha 354-372)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:UpdateIngredient (linha 301-312)`
**Model:** `backend/internal/domain/ingredient.go`
**SQL:** UPDATE ingredients
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 7.5 Deletar Ingrediente
**Frontend:** `/frontend/src/routes/(app)/ingredients/+page.svelte`
**Backend:** `DELETE /api/ingredients/{id}`
**Handler:** `backend/internal/handler/product_handler.go:DeleteIngredient (linha 234-250)`
**Service:** `backend/internal/service/product_service.go:DeleteIngredient (linha 374-383)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:DeleteIngredient (linha 314-326)`
**Model:** `backend/internal/domain/ingredient.go`
**SQL:** UPDATE ingredients SET deleted_at = ? (soft delete)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 7.6 Ajustar Estoque
**Frontend:** `/frontend/src/routes/(app)/ingredients/+page.svelte`
**Backend:** `PATCH /api/ingredients/{id}/stock`
**Handler:** `backend/internal/handler/product_handler.go:UpdateIngredientStock (linha 178-204)`
**Service:** `backend/internal/service/product_service.go:UpdateIngredientStock (linha 330-352)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:UpdateIngredient (linha 301-312)`
**Model:** `backend/internal/domain/ingredient.go`
**SQL:** UPDATE ingredients SET stock_quantity = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 8. CATEGORIAS

#### 8.1 Criar Categoria
**Backend:** `POST /api/categories`
**Handler:** `backend/internal/handler/category_handler.go:CreateCategory`
**Service:** `backend/internal/service/category_service.go:CreateCategory`
**Repository:** `backend/internal/infra/repository/gorm_category_repository.go:Create`
**Model:** `backend/internal/domain/category.go`
**SQL:** INSERT INTO categories
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 8.2 Listar Categorias
**Backend:** `GET /api/categories`
**Handler:** `backend/internal/handler/category_handler.go:ListCategories`
**Service:** `backend/internal/service/category_service.go:ListCategories`
**Repository:** `backend/internal/infra/repository/gorm_category_repository.go:List`
**Model:** `backend/internal/domain/category.go`
**SQL:** SELECT FROM categories
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 8.3 Obter Categoria
**Backend:** `GET /api/categories/{id}`
**Handler:** `backend/internal/handler/category_handler.go:GetCategory`
**Service:** `backend/internal/service/category_service.go:GetCategory`
**Repository:** `backend/internal/infra/repository/gorm_category_repository.go:FindByID`
**Model:** `backend/internal/domain/category.go`
**SQL:** SELECT FROM categories WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 8.4 Atualizar Categoria
**Backend:** `PUT /api/categories/{id}`
**Handler:** `backend/internal/handler/category_handler.go:UpdateCategory`
**Service:** `backend/internal/service/category_service.go:UpdateCategory`
**Repository:** `backend/internal/infra/repository/gorm_category_repository.go:Update`
**Model:** `backend/internal/domain/category.go`
**SQL:** UPDATE categories
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 8.5 Deletar Categoria
**Backend:** `DELETE /api/categories/{id}`
**Handler:** `backend/internal/handler/category_handler.go:DeleteCategory`
**Service:** `backend/internal/service/category_service.go:DeleteCategory`
**Repository:** `backend/internal/infra/repository/gorm_category_repository.go:Delete`
**Model:** `backend/internal/domain/category.go`
**SQL:** UPDATE categories SET deleted_at = ? (soft delete)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 9. PEDIDOS

#### 9.1 Criar Pedido
**Frontend:** `/frontend/src/routes/(app)/orders/new/+page.svelte`
**Backend:** `POST /api/orders`
**Handler:** `backend/internal/handler/order_handler.go:CreateOrder (linha 20-41)`
**Service:** `backend/internal/service/order_service.go:CreateOrder (linha 77-151)`
**Repository:** `backend/internal/infra/repository/gorm_order_repository.go:CreateOrder (linha 66-145)`
**Repository:** `backend/internal/infra/repository/gorm_product_repository.go:FindProductByID`, `GetProductIngredients`, `FindIngredientByID`
**Model:** `backend/internal/domain/order.go`, `backend/internal/domain/order_item.go`
**SQL:** BEGIN TRANSACTION, INSERT INTO orders, INSERT INTO order_items, UPDATE ingredients (baixa de estoque), COMMIT
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 9.2 Listar Pedidos
**Frontend:** `/frontend/src/routes/(app)/orders/+page.svelte`
**Backend:** `GET /api/orders`
**Handler:** `backend/internal/handler/order_handler.go:ListOrders (linha 43-51)`
**Service:** `backend/internal/service/order_service.go:ListOrders (linha 204-210)`
**Repository:** `backend/internal/infra/repository/gorm_order_repository.go:ListOrders (linha 147-159)`
**Model:** `backend/internal/domain/order.go`
**SQL:** SELECT FROM orders WHERE deleted_at IS NULL
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 9.3 Obter Pedido
**Frontend:** `/frontend/src/routes/(app)/orders/[id]/+page.svelte`
**Backend:** `GET /api/orders/{id}`
**Handler:** `backend/internal/handler/order_handler.go:GetOrder (linha 53-70)`
**Service:** `backend/internal/service/order_service.go:GetOrder (linha 212-221)`
**Repository:** `backend/internal/infra/repository/gorm_order_repository.go:FindOrderByID (linha 161-173)`
**Model:** `backend/internal/domain/order.go`
**SQL:** SELECT FROM orders WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 9.4 Atualizar Status do Pedido
**Frontend:** `/frontend/src/routes/(app)/orders/[id]/+page.svelte`
**Backend:** `PATCH /api/orders/{id}/status`
**Handler:** `backend/internal/handler/order_handler.go:UpdateOrderStatus (linha 72-102)`
**Service:** `backend/internal/service/order_service.go:UpdateOrderStatus (linha 223-283)`
**Repository:** `backend/internal/infra/repository/gorm_order_repository.go:UpdateOrderStatus`, `UpdateOrderStatusWithAdjustments`
**Model:** `backend/internal/domain/order.go`
**SQL:** UPDATE orders SET status = ?, INSERT INTO stock_adjustments (seja cancelamento)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 10. AJUSTES DE ESTOQUE

#### 10.1 Listar Ajustes Pendentes
**Frontend:** `/frontend/src/routes/(app)/stock-adjustments/+page.svelte`
**Backend:** `GET /api/stock-adjustments/pending`
**Handler:** `backend/internal/handler/stock_adjustment_handler.go:ListPendingAdjustments`
**Service:** `backend/internal/service/stock_adjustment_service.go:ListPendingAdjustments`
**Repository:** `backend/internal/infra/repository/gorm_stock_adjustment_repository.go:ListPending`
**Model:** `backend/internal/domain/stock_adjustment_pending.go`
**SQL:** SELECT FROM stock_adjustments WHERE status = 'pending'
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 10.2 Aprovar Ajuste
**Frontend:** `/frontend/src/routes/(app)/stock-adjustments/+page.svelte`
**Backend:** `POST /api/stock-adjustments/{id}/approve`
**Handler:** `backend/internal/handler/stock_adjustment_handler.go:ApproveAdjustment`
**Service:** `backend/internal/service/stock_adjustment_service.go:ApproveAdjustment`
**Repository:** `backend/internal/infra/repository/gorm_stock_adjustment_repository.go:Approve`
**Model:** `backend/internal/domain/stock_adjustment_pending.go`
**SQL:** UPDATE stock_adjustments SET status = 'approved', UPDATE ingredients (aplicar ajuste)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 10.3 Rejeitar Ajuste
**Frontend:** `/frontend/src/routes/(app)/stock-adjustments/+page.svelte`
**Backend:** `POST /api/stock-adjustments/{id}/reject`
**Handler:** `backend/internal/handler/stock_adjustment_handler.go:RejectAdjustment`
**Service:** `backend/internal/service/stock_adjustment_service.go:RejectAdjustment`
**Repository:** `backend/internal/infra/repository/gorm_stock_adjustment_repository.go:Reject`
**Model:** `backend/internal/domain/stock_adjustment_pending.go`
**SQL:** UPDATE stock_adjustments SET status = 'rejected'
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 11. TEMA (White Label)

#### 11.1 Obter Tema
**Backend:** `GET /api/theme`
**Handler:** `backend/internal/handler/theme_handler.go:GetTheme`
**Service:** `backend/internal/service/theme_service.go:GetTheme`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:FindByID`
**Model:** `backend/internal/domain/theme.go`, `backend/internal/domain/company.go`
**SQL:** SELECT FROM companies WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 11.2 Obter Tema Padrão
**Backend:** `GET /api/theme/default`
**Handler:** `backend/internal/handler/theme_handler.go:GetDefaultTheme`
**Service:** `backend/internal/service/theme_service.go:GetDefaultTheme`
**Repository:** Nenhum (retorna valores hardcoded)
**Model:** Nenhum
**SQL:** Nenhum
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 12. BUSINESS PROFILE

#### 12.1 Obter Business Profile
**Backend:** `GET /api/business/profile`
**Handler:** `backend/internal/handler/business_handler.go:GetBusinessProfile`
**Service:** `backend/internal/service/business_service.go:GetBusinessProfile`
**Repository:** `backend/internal/infra/repository/gorm_company_repository.go:FindByID`
**Model:** `backend/internal/domain/business_profile.go`, `backend/internal/domain/company.go`
**SQL:** SELECT FROM companies WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 12.2 Obter Business Profile Padrão
**Backend:** `GET /api/business/profile/default`
**Handler:** `backend/internal/handler/business_handler.go:GetDefaultBusinessProfile`
**Service:** `backend/internal/service/business_service.go:GetDefaultBusinessProfile`
**Repository:** Nenhum (retorna valores hardcoded)
**Model:** Nenhum
**SQL:** Nenhum
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 13. MÍDIA

#### 13.1 Upload de Mídia
**Backend:** `POST /api/media/upload`
**Handler:** `backend/internal/handler/media_handler.go:UploadMedia`
**Service:** `backend/internal/service/media_service.go:UploadMedia`
**Repository:** `backend/internal/infra/repository/gorm_media_repository.go:Create`
**Model:** `backend/internal/domain/media.go`
**SQL:** INSERT INTO media
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 13.2 Obter Mídia
**Backend:** `GET /api/media/{id}`
**Handler:** `backend/internal/handler/media_handler.go:GetMedia`
**Service:** `backend/internal/service/media_service.go:GetMedia`
**Repository:** `backend/internal/infra/repository/gorm_media_repository.go:FindByID`
**Model:** `backend/internal/domain/media.go`
**SQL:** SELECT FROM media WHERE id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 13.3 Deletar Mídia
**Backend:** `DELETE /api/media/{id}`
**Handler:** `backend/internal/handler/media_handler.go:DeleteMedia`
**Service:** `backend/internal/service/media_service.go:DeleteMedia`
**Repository:** `backend/internal/infra/repository/gorm_media_repository.go:Delete`
**Model:** `backend/internal/domain/media.go`
**SQL:** UPDATE media SET deleted_at = ? (soft delete)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 13.4 Obter Mídia por Entidade
**Backend:** `GET /api/media/entity/{entity_type}/{entity_id}`
**Handler:** `backend/internal/handler/media_handler.go:GetMediaByEntity`
**Service:** `backend/internal/service/media_service.go:GetMediaByEntity`
**Repository:** `backend/internal/infra/repository/gorm_media_repository.go:FindByEntity`
**Model:** `backend/internal/domain/media.go`
**SQL:** SELECT FROM media WHERE entity_type = ? AND entity_id = ?
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

#### 13.5 Servir Arquivo Estático
**Backend:** `GET /uploads/*`
**Handler:** `backend/internal/handler/media_handler.go:ServeFile`
**Service:** Nenhum
**Repository:** Nenhum
**Model:** Nenhum
**SQL:** Nenhum
**Middleware:** Nenhum (rota pública)

---

### 14. DASHBOARD

#### 14.1 Obter Dashboard
**Frontend:** `/frontend/src/routes/(app)/dashboard/+page.svelte`
**Backend:** `GET /api/dashboard`
**Handler:** `backend/internal/handler/dashboard_handler.go:GetDashboard`
**Service:** Nenhum (direto no handler)
**Repository:** `backend/internal/infra/repository/gorm_dashboard_repository.go:GetDashboard`
**Model:** `backend/internal/domain/dashboard.go`
**SQL:** SELECT FROM orders, SELECT FROM products, SELECT FROM ingredients (agregações)
**Middleware:** `authMw.Auth` (linha 109), `tenantMw.Tenant` (linha 110)

---

### 15. SYSTEM

#### 15.1 Health Check
**Backend:** `GET /api/health`
**Handler:** `backend/cmd/server/main.go:linha 91-95`
**Service:** Nenhum
**Repository:** Nenhum
**Model:** Nenhum
**SQL:** Nenhum
**Middleware:** Nenhum (rota pública)

---

## Resumo de Rotas por Middleware

### Rotas Públicas (Sem Autenticação)
- GET /api/health
- POST /api/auth/register
- POST /api/auth/login
- POST /api/auth/logout
- GET /api/invitations/{token}
- GET /uploads/*

### Rotas Protegidas (Com Auth + Tenant)
- GET /api/dashboard
- GET /api/me
- PUT /api/me
- POST /api/me/change-password
- POST /api/companies
- GET /api/companies
- GET /api/companies/{id}
- PUT /api/companies/{id}
- DELETE /api/companies/{id}
- GET /api/company/settings
- PUT /api/company/settings
- GET /api/theme
- GET /api/theme/default
- GET /api/business/profile
- GET /api/business/profile/default
- POST /api/products
- GET /api/products
- GET /api/products/active
- GET /api/products/{id}
- PUT /api/products/{id}
- DELETE /api/products/{id}
- PUT /api/products/{id}/ingredients
- GET /api/products/{id}/ingredients
- POST /api/ingredients
- GET /api/ingredients
- GET /api/ingredients/{id}
- PUT /api/ingredients/{id}
- DELETE /api/ingredients/{id}
- PATCH /api/ingredients/{id}/stock
- POST /api/categories
- GET /api/categories
- GET /api/categories/{id}
- PUT /api/categories/{id}
- DELETE /api/categories/{id}
- POST /api/orders
- GET /api/orders
- GET /api/orders/{id}
- PATCH /api/orders/{id}/status
- GET /api/stock-adjustments/pending
- POST /api/stock-adjustments/{id}/approve
- POST /api/stock-adjustments/{id}/reject
- POST /api/media/upload
- GET /api/media/{id}
- DELETE /api/media/{id}
- GET /api/media/entity/{entity_type}/{entity_id}
- POST /api/invitations/accept

### Rotas Protegidas (Com Auth + Tenant + RBAC: Owner/Admin)
- GET /api/company/users
- GET /api/company/users/{id}
- POST /api/company/users/add
- PUT /api/company/users/{id}/role
- DELETE /api/company/users/{id}
- POST /api/company/invitations
- GET /api/company/invitations
- GET /api/company/invitations/{id}
- DELETE /api/company/invitations/{id}
