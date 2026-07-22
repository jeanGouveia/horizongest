# Architecture Overview

**HorizonGest Platform - Architecture Documentation**

---

## Overview

HorizonGest follows a strict layered architecture with clear separation of concerns. The architecture is designed for multi-tenancy, white-label support, and long-term maintainability.

**Architecture Score:** 9.0/10 (Foundation Closed)

---

## Layered Architecture

### Structure

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
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

## Multi-Tenancy

### CompanyID Pattern

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

### Tenant Isolation

- Repository Layer filters by CompanyID
- Tenant Middleware verifies CompanyID
- Platform users cannot access tenant data
- Tenant users cannot access platform data

---

## Platform vs Tenant

### Platform Routes

**URL:** `/api/platform/*`
**Authentication:** Platform JWT
**Purpose:** Platform administration

### Tenant Routes

**URL:** `/api/*`
**Authentication:** Tenant JWT
**Purpose:** Business operations

### Separation

- Platform users cannot access tenant routes
- Tenant users cannot access platform routes
- Separate JWT secrets
- Separate authentication middleware

---

## Branding Architecture

### Platform Branding

**Table:** `platform_brand_config`
**Purpose:** Platform institutional identity
**Fields:** Name, logo, colors, email, copyright

### Tenant Branding

**Tables:** `themes`, `business_profiles`
**Purpose:** Customer-specific identity
**Fields:** Logo, colors, theme customization

### Public Endpoint

**URL:** `/api/public/brand`
**Authentication:** None
**Purpose:** Frontend branding without login
**Cache:** Repository Layer cache

### Dynamic Branding

- Backend: 100% dynamic (no hardcoded branding)
- Frontend: 95% dynamic (brandStore)
- Services: Use platform name dynamically
- JWT: Dynamic issuer from platform branding

---

## Configuration Architecture

### PlatformBrandConfig

**Purpose:** Branding/institutional
**Fields:** Name, logo, colors, email, copyright
**Access:** Public endpoint for frontend

### GlobalConfig

**Purpose:** Technical configuration
**Fields:** Timezone, locale, upload limits, feature flags
**Access:** Platform admin only

### Environment Variables

**Purpose:** Secrets and infrastructure
**Fields:** JWT secrets, DB password, DB host
**Access:** Application startup only

---

## Feature Flags

### Implementation

**Storage:** `global_config` table
**Fields:** `enable_finance`, `enable_purchasing`, etc.

### Usage

**Handler:** Check before exposing routes
**Service:** Check before executing logic
**Frontend:** Check before showing UI

### Module Registry

**Storage:** `module_registry` (domain layer)
**Purpose:** Track available modules and dependencies
**Sync:** Must sync with feature flags

---

## Cache Architecture

### Implementation

**Location:** Repository Layer
**Type:** In-memory with sync.RWMutex
**Pattern:** Cache-first logic

### Cache-First Logic

1. Check cache
2. If hit, return from cache
3. If miss, load from database
4. Update cache
5. Return data

### Invalidation

- Automatic on Update
- Explicit methods: `InvalidateCache()`, `ReloadCache()`
- Service Layer unaware of cache details

### Limitations

- Not distributed (single-instance only)
- Not persistent (lost on restart)
- TODO: Redis for distributed cache

---

## Security Architecture

### Authentication

**Method:** JWT (HS256)
**Expiration:** 24 hours
**Secrets:** Separate for platform and tenant
**Blacklist:** Token blacklist for logout

### Authorization

**Method:** RBAC (Role-Based Access Control)
**Roles:** Admin, Manager, Employee, Viewer
**Permissions:** Granular per operation
**Verification:** Service Layer

### Security Headers

- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security: max-age=31536000
- Content-Security-Policy: default-src 'self'

### Rate Limiting

**Platform routes:** 5 req/min per IP
**Tenant routes:** 30 req/hour per user
**Implementation:** Middleware

---

## Performance Architecture

### Database

- GORM for ORM
- Parameterized queries (prevent SQL injection)
- Indexes on foreign keys
- N+1 queries avoided (eager loading)

### Cache

- In-memory cache for frequently accessed data
- Cache-first logic
- Automatic invalidation

### Queries

- Paginated when appropriate
- Filtered by CompanyID
- No N+1 queries identified

---

## White-Label Architecture

### Current State

**Backend:** 100% ready
- Dynamic branding via PlatformBrandConfig
- Public endpoint for frontend
- Services use dynamic platform name

**Frontend:** 95% ready
- brandStore global
- Components updated
- TODO: Favicon, manifest, meta tags

### Multi-Brand Support

**Current:** Singleton pattern (ID=1)
**TODO:** Add brand_key for multi-brand
**TODO:** Map-based cache (keyed by brand)
**TODO:** Middleware for brand selection by domain

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

## Tech Stack

### Backend

- **Language:** Go 1.21+
- **Framework:** Chi (HTTP router)
- **ORM:** GORM
- **Database:** SQLite (dev), Oracle (prod)
- **Auth:** JWT (HS256)
- **Migrations:** Goose

### Frontend

- **Framework:** SvelteKit
- **Language:** TypeScript
- **UI:** Svelte 5 Runes
- **Icons:** Lucide Svelte
- **Styling:** CSS (TailwindCSS planned)

---

## Foundation Status

**Status:** ✅ **FOUNDATION CLOSED**

**Score:** 9.0/10

**Completed:**
- ✅ Layered architecture implemented
- ✅ Multi-tenancy with CompanyID
- ✅ Dynamic branding (backend 100%, frontend 95%)
- ✅ Feature flags implemented
- ✅ Cache in Repository Layer
- ✅ Security (JWT, RBAC, rate limiting)
- ✅ Unit tests for core services
- ✅ API documentation
- ✅ Security audit (8.5/10)
- ✅ Performance audit (8.0/10)
- ✅ White-label audit (8.0/10)

**Next Steps:**
- Focus on business modules
- Focus on UX improvements
- Focus on integrations
- No structural changes (except critical fixes)

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
