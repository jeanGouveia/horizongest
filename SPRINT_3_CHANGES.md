# Sprint 3 - Lista de Alterações de Arquivos

**Data:** 19/07/2026  
**Versão:** 3.0  
**Status:** Planejamento

---

## Resumo

**Total de Arquivos:** 42  
**Arquivos a Criar:** 18  
**Arquivos a Modificar:** 18  
**Arquivos a Remover:** 6

---

## Backend - Domain Layer

### Arquivos a Criar

#### 1. `backend/internal/domain/platform_user.go`
**Status:** NOVO  
**Motivo:** Nova entidade para usuários da plataforma  
**Conteúdo:**
```go
package domain

import "time"

type PlatformUser struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool
	Role         PlatformRole
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
```

#### 2. `backend/internal/domain/platform_role.go`
**Status:** NOVO  
**Motivo:** Enum de roles da plataforma  
**Conteúdo:**
```go
package domain

type PlatformRole string

const (
	PlatformRoleAdmin    PlatformRole = "PlatformAdmin"
	PlatformRoleSupport  PlatformRole = "PlatformSupport"
)

func (r PlatformRole) IsValid() bool {
	switch r {
	case PlatformRoleAdmin, PlatformRoleSupport:
		return true
	default:
		return false
	}
}

func (r PlatformRole) String() string {
	return string(r)
}

func (r PlatformRole) DisplayName() string {
	switch r {
	case PlatformRoleAdmin:
		return "Administrador da Plataforma"
	case PlatformRoleSupport:
		return "Suporte da Plataforma"
	default:
		return "Desconhecido"
	}
}
```

### Arquivos a Modificar

#### 3. `backend/internal/domain/role.go`
**Status:** MODIFICAR  
**Motivo:** Adicionar novos roles e remover roles obsoletos  
**Alterações:**
- Remover: `RoleCashier`, `RoleKitchen`, `RoleWaiter`
- Adicionar: `RoleEmployee`
- Manter: `RoleOwner`, `RoleAdmin`, `RoleManager`

**Antes:**
```go
const (
	RoleOwner   Role = "owner"
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleCashier Role = "cashier"
	RoleKitchen Role = "kitchen"
	RoleWaiter  Role = "waiter"
)
```

**Depois:**
```go
const (
	RoleOwner    Role = "owner"
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)
```

#### 4. `backend/internal/domain/user.go`
**Status:** MODIFICAR  
**Motivo:** Remover nullable de CompanyID e Role  
**Alterações:**
- `CompanyID *uint` → `CompanyID uint`
- `Role *Role` → `Role Role`

**Antes:**
```go
type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool
	CompanyID    *uint
	Role         *Role
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
```

**Depois:**
```go
type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool
	CompanyID    uint
	Role         Role
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
```

#### 5. `backend/internal/domain/tenant_context.go`
**Status:** MODIFICAR  
**Motivo:** Remover suporte a CompanyID NULL  
**Alterações:**
- Remover campo `IsSystemAdmin` (não utilizado)
- `CompanyID *uint` → `CompanyID uint`
- Remover método `HasCompany()` (sempre true)

**Antes:**
```go
type TenantContext struct {
	UserID       uint
	CompanyID    *uint
	IsSystemAdmin bool
}

func (tc *TenantContext) HasCompany() bool {
	return tc.CompanyID != nil
}

func (tc *TenantContext) GetCompanyID() uint {
	if tc.CompanyID == nil {
		return 0
	}
	return *tc.CompanyID
}
```

**Depois:**
```go
type TenantContext struct {
	UserID    uint
	CompanyID uint
}

func (tc *TenantContext) GetCompanyID() uint {
	return tc.CompanyID
}
```

---

## Backend - Service Layer

### Arquivos a Criar

#### 6. `backend/internal/service/platform_auth_service.go`
**Status:** NOVO  
**Motivo:** Autenticação específica para platform users  
**Conteúdo:**
- Login para platform users
- Geração de JWT para platform users
- Validação de token de platform users
- Logout de platform users

#### 7. `backend/internal/service/platform_service.go`
**Status:** NOVO  
**Motivo:** Gestão de empresas pela plataforma  
**Conteúdo:**
- Criar empresa
- Editar empresa
- Suspender empresa
- Excluir empresa
- Listar empresas
- Obter empresa por ID

#### 8. `backend/internal/service/platform_user_service.go`
**Status:** NOVO  
**Motivo:** Gestão de platform users  
**Conteúdo:**
- Criar platform user
- Listar platform users
- Atualizar platform user
- Desativar platform user
- Excluir platform user

### Arquivos a Modificar

#### 9. `backend/internal/service/auth_service.go`
**Status:** MODIFICAR  
**Motivo:** Remover auto-criação de empresas no Register  
**Alterações:**
- Remover método `Register()` completamente (linhas 61-120)
- Ou modificar para retornar 403 com mensagem de cadastro indisponível
- Remover dependência de `companyRepo` se Register for removido

**Linhas a Remover:** 61-120 (método Register completo)

#### 10. `backend/internal/service/invitation_service.go`
**Status:** MODIFICAR  
**Motivo:** Remover validação "usuário já pertence a outra empresa"  
**Alterações:**
- Remover validação em `CreateInvitation()` (linhas 82-85)
- Remover validação em `AcceptInvitation()` (linhas 254-257)
- Usuários sem empresa não existem mais, então convites sempre funcionam

**Linhas a Remover:**
- Linhas 82-85 em `CreateInvitation()`
- Linhas 254-257 em `AcceptInvitation()`

#### 11. `backend/internal/service/user_management_service.go`
**Status:** MODIFICAR  
**Motivo:** Remover método AddExistingUser (CompanyID NULL não existe mais)  
**Alterações:**
- Remover método `AddExistingUser()` (linhas 223-283)
- Remover método `RemoveFromCompany()` (linhas 174-221) - usuários não podem ficar sem empresa
- Adicionar validação: CompanyID sempre vem de `ctx.User.CompanyID`

**Linhas a Remover:** 174-283

---

## Backend - Handler Layer

### Arquivos a Criar

#### 12. `backend/internal/handler/platform_auth_handler.go`
**Status:** NOVO  
**Motivo:** Handlers de autenticação da plataforma  
**Conteúdo:**
- `Login()` - POST /platform/auth/login
- `Logout()` - POST /platform/auth/logout
- `Me()` - GET /platform/auth/me

#### 13. `backend/internal/handler/platform_handler.go`
**Status:** NOVO  
**Motivo:** Handlers de gestão da plataforma  
**Conteúdo:**
- `GetDashboard()` - GET /platform
- `RegisterRoutes()` - Registro de rotas

#### 14. `backend/internal/handler/platform_company_handler.go`
**Status:** NOVO  
**Motivo:** Handlers de gestão de empresas pela plataforma  
**Conteúdo:**
- `CreateCompany()` - POST /platform/companies
- `ListCompanies()` - GET /platform/companies
- `GetCompany()` - GET /platform/companies/{id}
- `UpdateCompany()` - PUT /platform/companies/{id}
- `DeleteCompany()` - DELETE /platform/companies/{id}
- `SuspendCompany()` - PATCH /platform/companies/{id}/suspend
- `CreateOwner()` - POST /platform/companies/{id}/owner

#### 15. `backend/internal/handler/platform_user_handler.go`
**Status:** NOVO  
**Motivo:** Handlers de gestão de platform users  
**Conteúdo:**
- `CreatePlatformUser()` - POST /platform/users
- `ListPlatformUsers()` - GET /platform/users
- `GetPlatformUser()` - GET /platform/users/{id}
- `UpdatePlatformUser()` - PUT /platform/users/{id}
- `DeactivatePlatformUser()` - PATCH /platform/users/{id}/active
- `DeletePlatformUser()` - DELETE /platform/users/{id}

### Arquivos a Modificar

#### 16. `backend/internal/handler/auth_handler.go`
**Status:** MODIFICAR  
**Motivo:** Remover ou desabilitar Register  
**Alterações:**
- Remover método `Register()` (linhas 31-57)
- Ou modificar para retornar HTTP 403 com mensagem

**Linhas a Remover:** 31-57

#### 17. `backend/internal/handler/invitation_handler.go`
**Status:** MODIFICAR  
**Motivo:** Simplificar validações (usuários sem empresa não existem mais)  
**Alterações:**
- Remover validação de "usuário já pertence a outra empresa" em `AcceptInvitation()`
- Simplificar lógica de aceitação

---

## Backend - Middleware Layer

### Arquivos a Criar

#### 18. `backend/internal/middleware/platform_auth_middleware.go`
**Status:** NOVO  
**Motivo:** Autenticação específica para platform users  
**Conteúdo:**
- `Auth()` - Valida token de platform user
- Injeta `PlatformUserID` no contexto
- Verifica se usuário é platform user

#### 19. `backend/internal/middleware/platform_role_middleware.go`
**Status:** NOVO  
**Motivo:** RBAC para platform users  
**Conteúdo:**
- `RequirePlatformRole()` - Exige role específica da plataforma
- `RequirePlatformAdmin()` - Exige PlatformAdmin
- `RequirePlatformSupport()` - Exige PlatformSupport

### Arquivos a Modificar

#### 20. `backend/internal/middleware/auth_middleware.go`
**Status:** MODIFICAR  
**Motivo:** Suportar ambos os tipos de usuários (platform e company)  
**Alterações:**
- Adicionar lógica para distinguir platform users de company users
- Injetar `UserType` no contexto ("platform" ou "company")
- Manter compatibilidade com JWT existente

#### 21. `backend/internal/middleware/tenant_middleware.go`
**Status:** MODIFICAR  
**Motivo:** Remover suporte a CompanyID NULL  
**Alterações:**
- Remover verificação de `user.CompanyID == nil`
- Assumir sempre que CompanyID existe
- Remover logs de "Core V1 users"

**Linhas a Modificar:** 47-52

#### 22. `backend/internal/middleware/role_middleware.go`
**Status:** MODIFICAR  
**Motivo:** Distinguir PlatformAdmin de Owner  
**Alterações:**
- Adicionar verificação de `UserType` no contexto
- PlatformAdmin não deve ter permissões de Owner
- Owner não deve ter permissões de PlatformAdmin

---

## Backend - Repository Layer

### Arquivos a Criar

#### 23. `backend/internal/ports/platform_user_repository.go`
**Status:** NOVO  
**Motivo:** Interface para repository de platform users  
**Conteúdo:**
```go
package ports

import (
	"context"
	
	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type PlatformUserRepository interface {
	Create(ctx context.Context, user *domain.PlatformUser) error
	FindByID(ctx context.Context, id uint) (*domain.PlatformUser, error)
	FindByEmail(ctx context.Context, email string) (*domain.PlatformUser, error)
	List(ctx context.Context) ([]*domain.PlatformUser, error)
	Update(ctx context.Context, user *domain.PlatformUser) error
	Delete(ctx context.Context, id uint) error
}
```

#### 24. `backend/internal/infra/repository/gorm_platform_user_repository.go`
**Status:** NOVO  
**Motivo:** Implementação GORM de PlatformUserRepository  
**Conteúdo:**
- Implementação completa da interface
- Queries usando GORM
- Soft delete suportado

---

## Backend - Router

### Arquivos a Modificar

#### 25. `backend/cmd/server/main.go`
**Status:** MODIFICAR  
**Motivo:** Adicionar rotas de plataforma e remover registro público  
**Alterações:**

**Adicionar (após linha 101):**
```go
// --- Rotas da Plataforma (Nível 1) ---
r.Route("/platform", func(r chi.Router) {
	platformAuthHandler.RegisterRoutes(r)
	platformHandler.RegisterRoutes(r)
	platformCompanyHandler.RegisterRoutes(r)
	platformUserHandler.RegisterRoutes(r)
})
```

**Modificar (linha 104):**
```go
// REMOVER: r.Post("/register", authHandler.Register)
// OU: r.Post("/register", authHandler.RegisterDisabled)
```

**Adicionar DI (após linha 62):**
```go
platformUserRepo := repository.NewGormPlatformUserRepository(db)
platformAuthSvc := service.NewPlatformAuthService(platformUserRepo, tokenBlacklistRepo)
platformSvc := service.NewPlatformService(companyRepo)
platformUserSvc := service.NewPlatformUserService(platformUserRepo)

platformAuthHandler := handler.NewPlatformAuthHandler(platformAuthSvc)
platformHandler := handler.NewPlatformHandler()
platformCompanyHandler := handler.NewPlatformCompanyHandler(platformSvc, userRepo)
platformUserHandler := handler.NewPlatformUserHandler(platformUserSvc)

platformAuthMw := middleware.NewPlatformAuthMiddleware(platformAuthSvc)
platformRoleMw := middleware.NewPlatformRoleMiddleware()
```

---

## Backend - Database

### Arquivos a Criar

#### 26. `backend/internal/infra/database/migrations/003_create_platform_users.sql`
**Status:** NOVO  
**Motivo:** Migration para criar tabela platform_users  
**Conteúdo:**
```sql
-- Criar tabela platform_users
CREATE TABLE platform_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT 1,
  role TEXT NOT NULL CHECK (role IN ('PlatformAdmin', 'PlatformSupport')),
  deleted_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices
CREATE INDEX idx_platform_users_email ON platform_users(email);
CREATE INDEX idx_platform_users_role ON platform_users(role);
CREATE INDEX idx_platform_users_active ON platform_users(active);
```

#### 27. `backend/internal/infra/database/migrations/004_alter_users_not_null.sql`
**Status:** NOVO  
**Motivo:** Migration para alterar users.company_id e users.role para NOT NULL  
**Conteúdo:**
```sql
-- NOTA: Esta migration deve ser executada APÓS migrar dados existentes
-- Para desenvolvimento, pode ser executada diretamente se não houver dados

-- Alterar company_id para NOT NULL
-- SQLite não suporta ALTER COLUMN diretamente, então precisamos:
-- 1. Criar nova tabela com NOT NULL
-- 2. Migrar dados
-- 3. Drop tabela antiga
-- 4. Renomear nova tabela

CREATE TABLE users_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT 1,
  company_id INTEGER NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'manager', 'employee')),
  deleted_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  FOREIGN KEY (company_id) REFERENCES companies(id)
);

INSERT INTO users_new (id, name, email, password_hash, active, company_id, role, deleted_at, created_at, updated_at)
SELECT id, name, email, password_hash, active, 
       COALESCE(company_id, 0) as company_id,  -- Valor padrão temporário
       COALESCE(role, 'employee') as role,  -- Valor padrão temporário
       deleted_at, created_at, updated_at
FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

-- Recriar índices
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(active);
```

#### 28. `backend/internal/infra/database/migrations/005_create_platform_admin.sql`
**Status:** NOVO  
**Motivo:** Migration para criar platform admin inicial  
**Conteúdo:**
```sql
-- Criar platform admin inicial
-- Senha: admin123 (deve ser alterada no primeiro login)
INSERT INTO platform_users (name, email, password_hash, role, active, created_at, updated_at)
VALUES (
  'Platform Admin',
  'admin@pratoonline.com',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',  -- bcrypt de "admin123"
  'PlatformAdmin',
  1,
  datetime('now'),
  datetime('now')
);
```

---

## Frontend - Routes

### Arquivos a Remover

#### 29. `frontend/src/routes/(auth)/register/+page.svelte`
**Status:** REMOVER  
**Motivo:** Cadastro público não existe mais  

#### 30. `frontend/src/routes/(auth)/register/+page.server.ts`
**Status:** REMOVER  
**Motivo:** Server actions de cadastro não existem mais  

### Arquivos a Criar

#### 31. `frontend/src/routes/(platform)/+page.svelte`
**Status:** NOVO  
**Motivo:** Dashboard da plataforma  
**Conteúdo:**
- Visão geral da plataforma
- Métricas de empresas
- Links para gestão de empresas
- Links para gestão de platform users

#### 32. `frontend/src/routes/(platform)/auth/+page.svelte`
**Status:** NOVO  
**Motivo:** Login da plataforma  
**Conteúdo:**
- Formulário de login específico para plataforma
- Diferenciação visual do login de empresa
- Redirecionamento para /platform após login

#### 33. `frontend/src/routes/(platform)/companies/+page.svelte`
**Status:** NOVO  
**Motivo:** Listagem de empresas  
**Conteúdo:**
- Tabela com todas as empresas
- Filtros (ativas, suspensas, por plano)
- Ações: Criar, Editar, Suspender, Excluir
- Link para criar owner

#### 34. `frontend/src/routes/(platform)/companies/create/+page.svelte`
**Status:** NOVO  
**Motivo:** Criação de empresa  
**Conteúdo:**
- Formulário: Nome, CNPJ, Email, Telefone, Plano, Status
- Validação de campos
- Criação de empresa
- Redirecionamento para criar owner

#### 35. `frontend/src/routes/(platform)/companies/[id]/+page.svelte`
**Status:** NOVO  
**Motivo:** Detalhes da empresa  
**Conteúdo:**
- Informações da empresa
- Estatísticas da empresa
- Ações: Editar, Suspender, Excluir
- Link para criar owner
- Lista de usuários da empresa

#### 36. `frontend/src/routes/(platform)/companies/[id]/owner/+page.svelte`
**Status:** NOVO  
**Motivo:** Criação de owner  
**Conteúdo:**
- Formulário: Nome, Email, Senha
- Validação de campos
- Criação de owner
- Envio de email com credenciais
- Redirecionamento para detalhes da empresa

#### 37. `frontend/src/routes/(platform)/users/+page.svelte`
**Status:** NOVO  
**Motivo:** Gestão de platform users  
**Conteúdo:**
- Lista de platform users
- Ações: Criar, Editar, Desativar, Excluir
- Filtros por role

#### 38. `frontend/src/routes/(platform)/users/create/+page.svelte`
**Status:** NOVO  
**Motivo:** Criação de platform user  
**Conteúdo:**
- Formulário: Nome, Email, Senha, Role
- Validação de campos
- Criação de platform user

### Arquivos a Modificar

#### 39. `frontend/src/routes/(auth)/login/+page.svelte`
**Status:** MODIFICAR  
**Motivo:** Adicionar opção de login da plataforma  
**Alterações:**
- Adicionar toggle/tabs para "Login Empresa" vs "Login Plataforma"
- Redirecionar para /platform/auth se selecionar plataforma
- Manter fluxo existente para login de empresa

#### 40. `frontend/src/routes/(auth)/login/+page.server.ts`
**Status:** MODIFICAR  
**Motivo:** Suportar redirecionamento para plataforma  
**Alterações:**
- Adicionar lógica para redirecionar para /platform/auth
- Manter lógica existente para /dashboard

---

## Frontend - API Client

### Arquivos a Criar

#### 41. `frontend/src/lib/api/platform_client.ts`
**Status:** NOVO  
**Motivo:** Client API específico para plataforma  
**Conteúdo:**
```typescript
export const platformApi = {
  auth: {
    login: (body: { email: string; password: string }) => request('/platform/auth/login', { method: 'POST', body: JSON.stringify(body) }),
    logout: () => request('/platform/auth/logout', { method: 'POST' }),
    me: () => request('/platform/auth/me'),
  },
  companies: {
    create: (body: CompanyInput) => request('/platform/companies', { method: 'POST', body: JSON.stringify(body) }),
    list: () => request('/platform/companies'),
    get: (id: number) => request(`/platform/companies/${id}`),
    update: (id: number, body: CompanyInput) => request(`/platform/companies/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
    delete: (id: number) => request(`/platform/companies/${id}`, { method: 'DELETE' }),
    suspend: (id: number) => request(`/platform/companies/${id}/suspend`, { method: 'PATCH' }),
    createOwner: (id: number, body: OwnerInput) => request(`/platform/companies/${id}/owner`, { method: 'POST', body: JSON.stringify(body) }),
  },
  users: {
    create: (body: PlatformUserInput) => request('/platform/users', { method: 'POST', body: JSON.stringify(body) }),
    list: () => request('/platform/users'),
    get: (id: number) => request(`/platform/users/${id}`),
    update: (id: number, body: PlatformUserInput) => request(`/platform/users/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
    deactivate: (id: number, active: boolean) => request(`/platform/users/${id}/active`, { method: 'PATCH', body: JSON.stringify({ active }) }),
    delete: (id: number) => request(`/platform/users/${id}`, { method: 'DELETE' }),
  },
};
```

### Arquivos a Modificar

#### 42. `frontend/src/lib/api/client.ts`
**Status:** MODIFICAR  
**Motivo:** Remover endpoint de registro público  
**Alterações:**
- Remover `auth.register` do objeto `api`
- Adicionar comentário explicando que registro foi desabilitado

---

## Resumo por Camada

### Domain Layer
- Criar: 2 arquivos
- Modificar: 3 arquivos
- Total: 5 arquivos

### Service Layer
- Criar: 3 arquivos
- Modificar: 3 arquivos
- Total: 6 arquivos

### Handler Layer
- Criar: 4 arquivos
- Modificar: 2 arquivos
- Total: 6 arquivos

### Middleware Layer
- Criar: 2 arquivos
- Modificar: 3 arquivos
- Total: 5 arquivos

### Repository Layer
- Criar: 2 arquivos
- Modificar: 0 arquivos
- Total: 2 arquivos

### Router
- Criar: 0 arquivos
- Modificar: 1 arquivo
- Total: 1 arquivo

### Database
- Criar: 3 arquivos
- Modificar: 0 arquivos
- Total: 3 arquivos

### Frontend Routes
- Criar: 7 arquivos
- Modificar: 2 arquivos
- Remover: 2 arquivos
- Total: 11 arquivos

### Frontend API
- Criar: 1 arquivo
- Modificar: 1 arquivo
- Total: 2 arquivos

---

## Ordem Sugerida de Implementação

### Fase 1: Backend Foundation (Dia 1-2)
1. Criar domain entities (platform_user.go, platform_role.go)
2. Modificar role.go
3. Modificar user.go
4. Modificar tenant_context.go
5. Criar migrations
6. Executar migrations

### Fase 2: Backend Services (Dia 2-3)
7. Criar platform_auth_service.go
8. Criar platform_service.go
9. Criar platform_user_service.go
10. Modificar auth_service.go
11. Modificar invitation_service.go
12. Modificar user_management_service.go

### Fase 3: Backend Handlers (Dia 3-4)
13. Criar platform_auth_handler.go
14. Criar platform_handler.go
15. Criar platform_company_handler.go
16. Criar platform_user_handler.go
17. Modificar auth_handler.go
18. Modificar invitation_handler.go

### Fase 4: Backend Middleware (Dia 4)
19. Criar platform_auth_middleware.go
20. Criar platform_role_middleware.go
21. Modificar auth_middleware.go
22. Modificar tenant_middleware.go
23. Modificar role_middleware.go

### Fase 5: Backend Repository (Dia 4)
24. Criar platform_user_repository.go
25. Criar gorm_platform_user_repository.go

### Fase 6: Backend Router (Dia 5)
26. Modificar main.go

### Fase 7: Frontend Foundation (Dia 5-6)
27. Remover telas de cadastro
28. Modificar login
29. Criar platform_client.ts
30. Modificar client.ts

### Fase 8: Frontend Platform UI (Dia 6-7)
31. Criar dashboard plataforma
32. Criar login plataforma
33. Criar listagem empresas
34. Criar criação empresa
35. Criar detalhes empresa
36. Criar criação owner
37. Criar gestão platform users

### Fase 9: Testes (Dia 7-8)
38. Testes de integração
39. Testes de regressão
40. Testes manuais

---

## Notas Importantes

1. **Migration de Dados:** A migration 004 deve ser testada cuidadosamente em ambiente de desenvolvimento antes de produção
2. **Platform Admin Inicial:** A migration 005 cria um admin com senha "admin123" que deve ser alterada imediatamente
3. **Compatibilidade:** Manter compatibilidade com JWT existente para não quebrar sessões ativas
4. **Logs:** Adicionar logs detalhados em middlewares para facilitar debugging
5. **Rollback:** Ter script de rollback pronto caso algo dê errado
6. **Testes:** Executar testes completos após cada fase para identificar problemas cedo
