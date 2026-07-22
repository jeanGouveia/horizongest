# AI_CONTEXT.md

**HorizonGest Platform - AI Context for Code Modifications**

---

## Purpose

This document provides all necessary context for any AI to modify HorizonGest code. It replaces large prompts and serves as the single source of truth for architectural rules, patterns, and conventions.

**Read this document before making ANY code changes.**

---

## Project Overview

### What is HorizonGest?

HorizonGest is a multi-tenant SaaS platform for restaurant management. It provides:
- Menu management
- Order processing
- Inventory tracking
- Financial management
- Customer relationship management
- Multi-tenant architecture (white-label ready)

### Tech Stack

**Backend:**
- Language: Go 1.21+
- Framework: Chi (HTTP router)
- ORM: GORM
- Database: SQLite (development), Oracle (production - migration path prepared)
- Auth: JWT (HS256)
- Migrations: Goose

**Frontend:**
- Framework: SvelteKit
- Language: TypeScript
- UI: Svelte 5 Runes
- Icons: Lucide Svelte
- Styling: TailwindCSS (planned)

### Architecture

**Backend follows strict layered architecture:**

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
```

**Frontend follows:**
- Pages in `src/routes/`
- Components in `src/lib/components/`
- Stores in `src/lib/stores/`
- API client in `src/lib/api/`

---

## Backend Architecture

### Directory Structure

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── domain/                 # Pure Go structs (no GORM tags)
│   ├── ports/                  # Interfaces (contracts)
│   ├── service/                # Business logic
│   ├── handler/                # HTTP handlers
│   ├── middleware/             # HTTP middleware
│   ├── infra/
│   │   ├── database/           # Database connection
│   │   └── repository/         # GORM implementations
│   └── util/                  # Utilities
└── migrations/                 # SQL migrations (Goose)
```

### Layer Responsibilities

**Handler Layer:**
- HTTP request/response handling
- Input validation (via validator)
- JSON parsing
- Calls Service Layer
- Returns HTTP status codes

**Service Layer:**
- ALL business logic
- Business validation
- Calls Repository Layer
- Returns domain entities

**Repository Layer:**
- Data access
- GORM queries
- Cache management (in-memory with sync.RWMutex)
- Returns domain entities

**Domain Layer:**
- Pure Go structs
- No GORM tags
- No dependencies on infrastructure

### Strict Rules

**NEVER:**
- Handler calling Repository
- Handler calling Handler
- Handler accessing database directly
- Service calling Handler
- Service accessing database directly
- Repository calling Service
- Repository calling Handler

**ALWAYS:**
- Handler → Service → Repository → Database
- Business logic in Service Layer
- Data access in Repository Layer
- Domain entities pure (no infrastructure)

---

## Frontend Architecture

### Directory Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api/               # API client
│   │   ├── components/        # UI components
│   │   ├── stores/            # Svelte 5 Runes stores
│   │   └── types/             # TypeScript types
│   └── routes/                # Pages and layouts
└── package.json
```

### Key Concepts

**Stores:**
- Global state management with Svelte 5 Runes
- `brandStore` for dynamic branding
- `userStore` for authentication
- Derived stores for computed values

**Routing:**
- File-based routing
- Layouts for shared UI
- Server-side API proxy via `hooks.server.ts`

**Branding:**
- All branding from `/api/public/brand`
- No hardcoded branding in frontend
- Dynamic colors, logos, names

---

## Creating New Modules

### Step 1: Register in Module Registry

**File:** `backend/internal/domain/module_registry.go`

Add module definition:

```go
modules["your_module"] = Module{
    Key:         "your_module",
    Name:        "Your Module",
    Description: "Module description",
    Category:    "business",
    Status:      "active",
    Dependencies: []string{"inventory"}, // If depends on other modules
    Features: map[string]bool{
        "feature1": true,
        "feature2": false,
    },
}
```

### Step 2: Add Feature Flag

**File:** `backend/internal/domain/global_config.go`

Add field:

```go
EnableYourModule bool `gorm:"column:enable_your_module" json:"enableYourModule"`
```

**File:** `backend/migrations/`

Create migration to add column:

```sql
-- YYYYMMDD_add_your_module_flag.sql
+goose Up
ALTER TABLE global_config ADD COLUMN enable_your_module BOOLEAN DEFAULT FALSE;

+goose Down
ALTER TABLE global_config DROP COLUMN enable_your_module;
```

### Step 3: Create Domain Entity

**File:** `backend/internal/domain/your_entity.go`

```go
package domain

import "time"

type YourEntity struct {
    ID          uint      `json:"id"`
    CompanyID   uint      `json:"companyId"` // REQUIRED for tenant entities
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Active      bool      `json:"active"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
```

**Rules:**
- Use `CompanyID` for ALL tenant entities
- Use `DeletedAt` for soft delete
- Use `Active` for availability
- Use `CreatedAt`, `UpdatedAt` for timestamps

### Step 4: Create Repository

**File:** `backend/internal/infra/repositorygorm_your_entity_repository.go`

```go
package repository

import (
    "context"
    "sync"

    "github.com/jeanGouveia/horizongest/backend/internal/domain"
    "gorm.io/gorm"
)

type GormYourEntityRepository struct {
    db    *gorm.DB
    cache *domain.YourEntity
    mutex sync.RWMutex
}

func NewGormYourEntityRepository(db *gorm.DB) *GormYourEntityRepository {
    return &GormYourEntityRepository{db: db}
}

func (r *GormYourEntityRepository) Get(ctx context.Context, id uint) (*domain.YourEntity, error) {
    // Implementation with cache-first logic
}

func (r *GormYourEntityRepository) Create(ctx context.Context, entity *domain.YourEntity) error {
    // Implementation
}

func (r *GormYourEntityRepository) Update(ctx context.Context, entity *domain.YourEntity) error {
    // Implementation with cache invalidation
}

func (r *GormYourEntityRepository) Delete(ctx context.Context, id uint) error {
    // Implementation (soft delete)
}

func (r *GormYourEntityRepository) List(ctx context.Context, companyID uint) ([]*domain.YourEntity, error) {
    // Implementation with CompanyID filter
}
```

**Rules:**
- Always filter by `CompanyID` for tenant entities
- Use cache-first logic for frequently accessed data
- Use `sync.RWMutex` for thread-safety
- Invalidate cache on Update

### Step 5: Create Service

**File:** `backend/internal/service/your_entity_service.go`

```go
package service

import (
    "context"
    "errors"

    "github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type YourEntityService struct {
    repo          domain.YourEntityRepository
    globalConfig  domain.GlobalConfigRepository
}

func NewYourEntityService(repo domain.YourEntityRepository, globalConfig domain.GlobalConfigRepository) *YourEntityService {
    return &YourEntityService{
        repo:         repo,
        globalConfig: globalConfig,
    }
}

func (s *YourEntityService) Get(ctx context.Context, id uint) (*domain.YourEntity, error) {
    // Check feature flag
    config, err := s.globalConfig.Get(ctx)
    if err != nil {
        return nil, err
    }
    if !config.EnableYourModule {
        return nil, errors.New("module not enabled")
    }

    return s.repo.Get(ctx, id)
}

func (s *YourEntityService) Create(ctx context.Context, entity *domain.YourEntity, userID uint) error {
    // Business validation
    if entity.Name == "" {
        return errors.New("name is required")
    }

    return s.repo.Create(ctx, entity)
}
```

**Rules:**
- Check feature flag before executing logic
- Business validation in Service Layer
- Call Repository Layer for data access

### Step 6: Create Handler

**File:** `backend/internal/handler/your_entity_handler.go`

```go
package handler

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jeanGouveia/horizongest/backend/internal/service"
)

type YourEntityHandler struct {
    service *service.YourEntityService
}

func NewYourEntityHandler(service *service.YourEntityService) *YourEntityHandler {
    return &YourEntityHandler{service: service}
}

func (h *YourEntityHandler) RegisterRoutes(r chi.Router) {
    r.Route("/your-entities", func(r chi.Router) {
        r.Use(h.authMiddleware)
        r.Get("/", h.List)
        r.Post("/", h.Create)
        r.Get("/{id}", h.Get)
        r.Put("/{id}", h.Update)
        r.Delete("/{id}", h.Delete)
    })
}

func (h *YourEntityHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Validate input
    // Call service
    // Return response
}
```

**Rules:**
- Input validation in Handler
- Call Service Layer
- Return appropriate HTTP status codes
- Use middleware for auth/authorization

### Step 7: Add Routes

**File:** `backend/cmd/server/main.go`

```go
yourEntityHandler := handler.NewYourEntityHandler(yourEntityService)
yourEntityHandler.RegisterRoutes(r)
```

### Step 8: Create Migration

**File:** `backend/migrations/YYYYMMDD_create_your_entity.sql`

```sql
+goose Up
CREATE TABLE your-entities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (company_id) REFERENCES companies(id)
);

CREATE INDEX idx_your_entities_company_id ON your-entities(company_id);

+goose Down
DROP TABLE your-entities;
```

**Rules:**
- Use `YYYYMMDD_description.sql` naming
- Include `+goose Up` and `+goose Down`
- Down must completely revert Up
- Use `CREATE TABLE IF NOT EXISTS` for idempotency
- Add indexes for foreign keys

---

## Creating Migrations

### Naming Convention

`YYYYMMDD_description.sql`

Example: `20240120_create_products.sql`

### Structure

```sql
-- +goose Up
-- Your migration SQL here

-- +goose Down
-- Rollback SQL here
```

### Idempotency

Use:
- `CREATE TABLE IF NOT EXISTS`
- `INSERT OR IGNORE`
- `ALTER TABLE IF NOT EXISTS`

### Foreign Keys

Always add indexes for foreign keys:

```sql
CREATE INDEX idx_table_foreign_key ON table(foreign_key_id);
```

### Running Migrations

```bash
cd backend
goose -dir migrations sqlite "app.db" up
```

---

## Creating APIs

### REST Conventions

**Routes:** kebab-case
- `/api/users`
- `/api/platform-brand`

**HTTP Methods:**
- `GET` - Retrieve
- `POST` - Create
- `PUT` - Update (full)
- `PATCH` - Update (partial)
- `DELETE` - Delete

**Response Format:**

```json
{
  "id": 1,
  "name": "Example",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Error Format:**

```json
{
  "error": "error message"
}
```

### Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Not authenticated
- `403 Forbidden` - Not authorized
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Creating Frontend Screens

### Page Structure

**File:** `frontend/src/routes/(platform)/your-page/+page.svelte`

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { api } from '$lib/api/client';
    
    let data = $state([]);
    let loading = $state(false);

    onMount(async () => {
        loading = true;
        try {
            const response = await api.yourModule.list();
            data = response;
        } catch (error) {
            console.error(error);
        } finally {
            loading = false;
        }
    });
</script>

{#if loading}
    <p>Loading...</p>
{:else if data.length === 0}
    <p>No data</p>
{:else}
    <!-- Render data -->
{/if}
```

### Component Structure

**File:** `frontend/src/lib/components/YourComponent.svelte`

```svelte
<script lang="ts">
    interface Props {
        data: YourDataType;
        onSave: (data: YourDataType) => void;
    }

    let { data, onSave }: Props = $props();
</script>

<!-- Component UI -->
```

### Using Branding

```svelte
<script lang="ts">
    import { platformName, primaryColor } from '$lib/stores/brandStore';
</script>

<h1>{$platformName}</h1>
<div style="color: {$primaryColor}">
    <!-- Content -->
</div>
```

---

## Creating Tests

### Unit Tests (Service)

**File:** `backend/internal/service/your_service_test.go`

```go
package service

import (
    "context"
    "errors"
    "testing"

    "github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type mockRepository struct {
    data *domain.YourEntity
    err  error
}

func (m *mockRepository) Get(ctx context.Context, id uint) (*domain.YourEntity, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.data, nil
}

func TestYourService_Get(t *testing.T) {
    tests := []struct {
        name    string
        repo    *mockRepository
        wantErr bool
    }{
        {
            name: "success",
            repo: &mockRepository{
                data: &domain.YourEntity{ID: 1},
            },
            wantErr: false,
        },
        {
            name: "error",
            repo: &mockRepository{
                err: errors.New("error"),
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := NewYourService(tt.repo, nil)
            _, err := svc.Get(context.Background(), 1)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Running Tests

```bash
cd backend
go test ./internal/service/...
```

---

## Adding Entities

### Domain Entity

**File:** `backend/internal/domain/your_entity.go`

```go
package domain

import "time"

type YourEntity struct {
    ID          uint      `json:"id"`
    CompanyID   uint      `json:"companyId"` // REQUIRED for tenant entities
    Name        string    `json:"name"`
    Active      bool      `json:"active"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
```

### Migration

**File:** `backend/migrations/YYYYMMDD_create_your_entity.sql`

```sql
+goose Up
CREATE TABLE your_entities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (company_id) REFERENCES companies(id)
);

CREATE INDEX idx_your_entities_company_id ON your_entities(company_id);

+goose Down
DROP TABLE your_entities;
```

### Repository

**File:** `backend/internal/infra/repositorygorm_your_entity_repository.go`

Implement standard CRUD operations with:
- CompanyID filtering
- Cache-first logic
- Soft delete

---

## Using Feature Flags

### Check Feature Flag in Service

```go
config, err := s.globalConfig.Get(ctx)
if err != nil {
    return nil, err
}
if !config.EnableYourModule {
    return nil, errors.New("module not enabled")
}
```

### Check Feature Flag in Handler

```go
config, err := h.globalConfigService.Get(r.Context())
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
if !config.EnableYourModule {
    http.Error(w, "module not enabled", http.StatusForbidden)
    return
}
```

### Check Feature Flag in Frontend

```typescript
const config = await api.globalConfig.get();
if (!config.enableYourModule) {
    // Hide UI or show message
}
```

---

## Mandatory Architecture Rules

### Layer Separation

**STRICT:**
- Handler → Service → Repository → Database
- NO direct Handler → Repository
- NO direct Service → Database
- NO Repository → Service

### Business Logic

**ALL business logic MUST be in Service Layer:**
- Calculations
- Validations
- Permissions
- Rules

**NEVER in:**
- Handler
- Repository
- Frontend

### CompanyID

**ALL tenant entities MUST have CompanyID:**
- Required field
- Filtered in all queries
- Used for tenant isolation

**Global entities WITHOUT CompanyID:**
- PlatformUser
- PlatformSession
- PlatformAudit
- Plan
- PlatformBrandConfig
- GlobalConfig

### Soft Delete

**ALL entities MUST have DeletedAt:**
- Use soft delete by default
- Hard delete only in exceptional cases
- Queries filter `DeletedAt IS NULL`

### Branding

**ALL branding MUST be dynamic:**
- From PlatformBrandConfig (platform)
- From Theme/BusinessProfile (tenant)
- NO hardcoded branding in code
- Frontend consumes `/api/public/brand`

### Configuration

**PlatformBrandConfig:**
- Branding/institutional
- Name, logo, colors, email, copyright

**GlobalConfig:**
- Technical configuration
- Timezone, locale, upload limits, feature flags

**Environment Variables:**
- Secrets (JWT, DB password)
- Infrastructure (DB host, port)

### Naming Conventions

**Go:**
- Structs: PascalCase (User, Product)
- Methods: PascalCase (GetUser, CreateProduct)
- Variables: camelCase (userID, productName)
- Constants: PascalCase or UPPER_CASE

**Database:**
- Tables: snake_case (users, platform_brand_config)
- Columns: snake_case (user_id, platform_name)

**API:**
- Routes: kebab-case (/api/users)
- JSON: camelCase (userId, platformName)

### Error Handling

**Return errors, don't panic:**
- Specific errors (ErrUserNotFound)
- Propagate to Handler
- Log in Service Layer

**HTTP Status Codes:**
- 400: Validation errors
- 401: Authentication errors
- 403: Authorization errors
- 404: Not found
- 500: Server errors

---

## Common Patterns

### Repository with Cache

```go
type GormYourRepository struct {
    db    *gorm.DB
    cache *domain.YourEntity
    mutex sync.RWMutex
}

func (r *GormYourRepository) Get(ctx context.Context, id uint) (*domain.YourEntity, error) {
    r.mutex.RLock()
    if r.cache != nil && r.cache.ID == id {
        r.mutex.RUnlock()
        return r.cache, nil
    }
    r.mutex.RUnlock()

    var entity domain.YourEntity
    err := r.db.WithContext(ctx).First(&entity, id).Error
    if err != nil {
        return nil, err
    }

    r.mutex.Lock()
    r.cache = &entity
    r.mutex.Unlock()

    return &entity, nil
}
```

### Service with Feature Flag

```go
func (s *YourService) DoSomething(ctx context.Context) error {
    config, err := s.globalConfig.Get(ctx)
    if err != nil {
        return err
    }
    if !config.EnableYourModule {
        return errors.New("module not enabled")
    }

    // Business logic
    return nil
}
```

### Handler with Validation

```go
func (h *YourHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if req.Name == "" {
        http.Error(w, "name is required", http.StatusBadRequest)
        return
    }

    entity := &domain.YourEntity{Name: req.Name}
    if err := h.service.Create(r.Context(), entity); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(entity)
}
```

---

## Domain Principles

### Identity is Immutable

- Entity IDs never change
- Never reuse IDs
- All other attributes can evolve

### Product Live vs Sold

- **Product Live:** Current catalog (can change)
- **Product Sold:** Historical snapshot (immutable)
- Sales create snapshots, don't reference live product

### History is Immutable

- Past events cannot be altered
- Changes to products/ingredients don't affect old orders
- Snapshots preserve historical state

### Active vs Deleted

- **Active:** Can be used by business
- **DeletedAt:** Logically removed
- Never have `active=true` and `deleted_at!=NULL`

### Stock

- Represents current availability
- Not historical
- Movements recorded separately
- Decrements only during sales

---

## Environment Variables

### Required

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=pratodb

# JWT
JWT_PLATFORM_SECRET=
JWT_TENANT_SECRET=

# Application
APP_VERSION=1.0.0
PORT=8080
```

### Optional

```bash
# Email
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
EMAIL_ENABLED=false
```

---

## Testing Commands

### Backend

```bash
cd backend
go test ./...
go test ./internal/service/... -v
go test ./internal/domain/... -v
```

### Frontend

```bash
cd frontend
npm run check
npm run dev
```

---

## Common Issues

### CompanyID Missing

**Error:** Missing CompanyID in entity

**Fix:** Add `CompanyID uint` to domain entity and migration

### Hardcoded Branding

**Error:** Hardcoded "HorizonGest" in code

**Fix:** Use PlatformBrandConfig from database

### Business Logic in Handler

**Error:** Business logic in Handler

**Fix:** Move to Service Layer

### N+1 Queries

**Error:** N+1 query problem

**Fix:** Use GORM Preload for relationships

---

## Summary

**Before modifying code:**
1. Read this document
2. Read ARCHITECTURE_RULES.md
3. Check existing patterns
4. Follow naming conventions
5. Respect layer separation
6. Use feature flags
7. Add CompanyID for tenant entities
8. Use soft delete
9. Make branding dynamic
10. Write tests

**After modifying code:**
1. Run tests
2. Check for violations
3. Update documentation
4. Create migration if needed

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
