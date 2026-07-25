# DECISIONS.md

**HorizonGest Platform - Architectural Decisions History**

---

## Overview

This document records important architectural decisions made during the development of HorizonGest. Each decision includes the problem, the decision made, the justification, and the consequences.

---

## Decision 1: SQLite for Development, Oracle for Production

**Date:** Initial Development  
**Status:** Active

### Problem

Need a database solution that is:
- Easy to set up for development
- Production-ready for enterprise clients
- Supports migration path

### Decision

Use SQLite for development and Oracle for production.

### Justification

- **SQLite:** Zero configuration, file-based, perfect for local development
- **Oracle:** Enterprise-grade, required by many enterprise clients
- **Migration path:** GORM abstraction allows easy switch
- **Single file:** `connection.go` isolates database-specific code

### Consequences

- Development is fast and simple
- Production deployment requires Oracle infrastructure
- Database-specific code isolated in one file
- Migration requires testing on Oracle before production

---

## Decision 2: GORM as ORM

**Date:** Initial Development  
**Status:** Active

### Problem

Need an ORM for Go that:
- Is well-maintained
- Supports SQLite and Oracle
- Prevents SQL injection
- Provides caching

### Decision

Use GORM as the ORM.

### Justification

- **Mature:** Well-established, widely used
- **Multi-database:** Supports SQLite and Oracle
- **Security:** Parameterized queries prevent SQL injection
- **Features:** Built-in soft delete, hooks, associations
- **Community:** Large community, good documentation

### Consequences

- Database access is abstracted
- SQL injection risk minimized
- Soft delete built-in
- Learning curve for team unfamiliar with GORM

---

## Decision 3: CompanyID for Multi-Tenancy

**Date:** Initial Development  
**Status:** Active

### Problem

Need to isolate data between different companies (tenants) in a shared database.

### Decision

Use CompanyID field in all tenant entities for multi-tenancy isolation.

### Justification

- **Simple:** Easy to understand and implement
- **Efficient:** Single database, filtered queries
- **Scalable:** Can handle many tenants
- **Secure:** Automatic data isolation via Repository Layer
- **Standard:** Common pattern in multi-tenant SaaS

### Consequences

- All tenant entities must include CompanyID
- Repository Layer must filter by CompanyID
- Global entities (PlatformUser, etc.) don't have CompanyID
- Database queries always include CompanyID filter

---

## Decision 4: Feature Flags in GlobalConfig

**Date:** Sprint 3.6  
**Status:** Active

### Problem

Need to enable/disable modules dynamically without code deployment.

### Decision

Store feature flags in GlobalConfig table in database.

### Justification

- **Dynamic:** Can change without code deployment
- **Centralized:** All flags in one place
- **Queryable:** Easy to query and update
- **Auditable:** Changes can be tracked
- **Flexible:** Can enable/disable per module

### Consequences

- Module Registry must sync with feature flags
- Services must check flags before executing logic
- Handlers must check flags before exposing routes
- Frontend must check flags before showing UI

---

## Decision 5: Module Registry

**Date:** Sprint 3.6  
**Status:** Active

### Problem

Need to track available modules, their dependencies, and their status.

### Decision

Create ModuleRegistry in domain layer to track all modules.

### Justification

- **Centralized:** Single source of truth for modules
- **Dependencies:** Can track module dependencies
- **Validation:** Can validate module configurations
- **Documentation:** Self-documenting module structure
- **Extensible:** Easy to add new modules

### Consequences

- All modules must be registered
- Dependencies must be declared
- Module status must be maintained
- Registry must be kept in sync with feature flags

---

## Decision 6: PlatformBrandConfig vs GlobalConfig

**Date:** Sprint 3.6  
**Status**: Active

### Problem

Need to separate branding/institutional configuration from technical configuration.

### Decision

Create two separate tables: PlatformBrandConfig and GlobalConfig.

### Justification

- **Separation of Concerns:** Branding vs technical settings
- **Clarity:** Clear distinction between types of configuration
- **Security:** Branding can be public, technical config private
- **Flexibility:** Can update branding without affecting technical settings
- **White-label:** Easier to support multiple brands

### Consequences

- PlatformBrandConfig: name, logo, colors, email, copyright
- GlobalConfig: timezone, locale, upload limits, feature flags
- Services must know which config to use
- Frontend consumes PlatformBrandConfig via public endpoint

---

## Decision 7: Layered Architecture (Handler → Service → Repository)

**Date:** Initial Development  
**Status:** Active

### Problem

Need a clear separation of concerns to maintain code quality and testability.

### Decision

Implement strict layered architecture: Handler → Service → Repository → Database.

### Justification

- **Separation:** Clear separation of concerns
- **Testability:** Each layer can be tested independently
- **Maintainability:** Easy to locate and fix bugs
- **Scalability:** Can scale layers independently
- **Standard:** Well-established pattern

### Consequences

- Strict rules about layer dependencies
- Handler cannot call Repository
- Service cannot access database directly
- Repository cannot call Service
- All business logic in Service Layer

---

## Decision 8: Branding Separation (Platform vs Tenant)

**Date:** Sprint 3.6  
**Status**: Active

### Problem

Need to separate platform branding (institutional) from tenant branding (customer-specific).

### Decision

Platform branding in PlatformBrandConfig, tenant branding in Theme and BusinessProfile.

### Justification

- **White-label:** Platform can be rebranded
- **Multi-tenant:** Each tenant can have own branding
- **Separation:** Clear distinction between platform and tenant
- **Flexibility:** Platform branding independent of tenant branding
- **Security:** Platform branding public, tenant branding private

### Consequences

- Platform branding: identity of the platform
- Tenant branding: identity of the customer
- Frontend consumes platform branding via public endpoint
- Frontend consumes tenant branding via authenticated endpoint

---

## Decision 9: Cache in Repository Layer

**Date:** Sprint 3.6  
**Status**: Active

### Problem

Need to improve performance for frequently accessed data (PlatformBrandConfig, GlobalConfig).

### Decision

Implement cache in Repository Layer with sync.RWMutex for thread-safety.

### Justification

- **Performance:** Reduces database queries
- **Thread-safe:** sync.RWMutex prevents race conditions
- **Simple:** In-memory cache, easy to implement
- **Service-agnostic:** Service Layer doesn't know about cache
- **Invalidation:** Automatic invalidation on Update

### Consequences

- Repository Layer manages cache
- Cache-first logic for frequently accessed data
- Service Layer unaware of cache details
- Cache invalidation on Update
- Not distributed (single-instance only)

---

## Decision 10: Soft Delete with DeletedAt

**Date:** Initial Development  
**Status**: Active

### Problem

Need to preserve data for audit and historical purposes while allowing logical deletion.

### Decision

Use DeletedAt field for soft delete in all entities.

### Justification

- **Audit:** Preserves historical data
- **Recovery:** Can recover deleted data
- **Compliance:** Meets regulatory requirements
- **Standard:** Common pattern in enterprise systems
- **GORM:** Built-in soft delete support

### Consequences

- All entities must have DeletedAt
- Queries filter DeletedAt IS NULL by default
- Hard delete only in exceptional cases
- Database size grows over time
- Need cleanup strategy for old deleted records

---

## Decision 11: JWT for Authentication

**Date:** Initial Development  
**Status**: Active

### Problem

Need stateless authentication that works across multiple instances.

### Decision

Use JWT (JSON Web Tokens) with HS256 for authentication.

### Justification

- **Stateless:** No server-side session storage
- **Scalable:** Works across multiple instances
- **Standard:** Well-established standard
- **Secure:** Cryptographically signed
- **Flexible:** Can include claims (UserID, CompanyID)

### Consequences

- Tokens have expiration (24 hours)
- Need token blacklist for logout
- Separate secrets for platform and tenant
- Tokens must be validated on each request
- Token revocation requires blacklist

---

## Decision 12: RBAC for Authorization

**Date:** Sprint 3.6  
**Status**: Active

### Problem

Need fine-grained access control for different user roles.

### Decision

Implement Role-Based Access Control (RBAC) with granular permissions.

### Justification

- **Granular:** Fine-grained permissions
- **Flexible:** Can assign permissions to roles
- **Standard:** Common pattern in enterprise systems
- **Scalable:** Easy to add new permissions
- **Auditable:** Can track permission changes

### Consequences

- RBACService manages permissions
- Middleware checks permissions
- Services verify permissions for operations
- Roles: admin, manager, employee, viewer
- Permissions must be defined for each operation

---

## Decision 13: Chi as HTTP Router

**Date:** Initial Development  
**Status**: Active

### Problem

Need an HTTP router for Go that is lightweight and flexible.

### Decision

Use Chi as the HTTP router.

### Justification

- **Lightweight:** Minimal overhead
- **Flexible:** Composable middleware
- **Standard:** Well-established in Go community
- **Context:** Native context support
- **Fast:** Good performance

### Consequences

- Routes defined in handlers
- Middleware composable
- Context propagation
- Learning curve for team unfamiliar with Chi

---

## Decision 14: SvelteKit for Frontend

**Date:** Initial Development  
**Status**: Active

### Problem

Need a modern frontend framework with good developer experience.

### Decision

Use SvelteKit with Svelte 5 Runes.

### Justification

- **Modern:** Latest web standards
- **Performant:** Compiled to vanilla JS
- **Developer Experience:** Great DX with Svelte 5 Runes
- **TypeScript:** Native TypeScript support
- **SSR:** Server-side rendering built-in

### Consequences

- File-based routing
- Server-side API proxy
- Svelte 5 Runes for reactivity
- Learning curve for team unfamiliar with Svelte

---

## Decision 15: Product Live vs Product Sold

**Date:** Initial Development  
**Status**: Active

### Problem

Need to preserve historical product data in orders while allowing product catalog changes.

### Decision

Create snapshots of products in orders (Product Sold) separate from live catalog (Product Live).

### Justification

- **Historical Accuracy:** Orders reflect state at time of sale
- **Flexibility:** Catalog can change without affecting history
- **Audit:** Complete historical record
- **Compliance:** Meets regulatory requirements
- **Domain Principle:** History is immutable

### Consequences

- Orders contain product snapshots
- Product catalog can change freely
- Snapshots never reference live product
- Increased storage for snapshots
- More complex order creation logic

---

## Decision 16: Active vs DeletedAt

**Date:** Initial Development  
**Status**: Active

### Problem

Need to distinguish between "available for use" and "logically removed".

### Decision

Use Active field for availability, DeletedAt for logical removal.

### Justification

- **Separation:** Clear distinction between availability and removal
- **Flexibility:** Can deactivate without deleting
- **Audit:** Preserves deleted records
- **Standard:** Common pattern in enterprise systems
- **Domain Principle:** Single responsibility

### Consequences

- Active: can be used by business
- DeletedAt: logically removed
- Never have active=true and deleted_at!=NULL
- Two fields to manage
- Clear semantics

---

## Decision 17: Domain-Driven Design Principles

**Date:** Initial Development  
**Status**: Active

### Problem

Need to ensure business logic is correctly implemented and maintainable.

### Decision

Follow Domain-Driven Design principles with domain as source of truth.

### Justification

- **Business-First:** Domain drives implementation
- **Maintainable:** Clear business rules
- **Evolution:** System can evolve without rewrite
- **Quality:** Architecture prioritized over speed
- **Standard:** Well-established methodology

### Consequences

- Domain entities pure (no infrastructure)
- Business logic in Service Layer
- Database follows domain
- No premature generalization
- Continuous evolution reduces technical debt

---

## Decision 18: Environment Variables for Secrets

**Date:** Initial Development  
**Status**: Active

### Problem

Need to securely manage secrets (JWT, DB password) without hardcoding.

### Decision

Store secrets in environment variables.

### Justification

- **Security:** Secrets not in code
- **Flexibility:** Different values per environment
- **Standard:** 12-factor app methodology
- **Simple:** Easy to implement
- **Auditable:** Can track environment changes

### Consequences

- Secrets in environment variables
- Configuration files for different environments
- Secrets never logged
- Need to manage environment variables in production

---

## Decision 19: APP_VERSION in Environment Variable

**Date:** Sprint 3.6  
**Status**: Active

### Problem

Need to track application version without storing in database.

### Decision

Store APP_VERSION in environment variable, use helper util.PlatformVersion().

### Justification

- **Simple:** Easy to update
- **Deployment:** Can set per deployment
- **No Database:** Version not in database
- **Standard:** Common pattern
- **Flexible:** Can include build info

### Consequences

- Version in environment variable
- Helper function to retrieve version
- Version not persisted in database
- Must update environment variable on deployment

---

## Decision 20: Public Endpoint for Branding

**Date:** Sprint 3.7  
**Status**: Active

### Problem

Frontend needs to access branding without authentication for white-label support.

### Decision

Create public endpoint `/api/public/brand` accessible without authentication.

### Justification

- **White-label:** Frontend can load branding before login
- **Security:** Only public-safe information exposed
- **Performance:** Cached in Repository Layer
- **Standard:** Common pattern in white-label systems
- **Flexibility:** Easy to extend for other public data

### Consequences

- Public endpoint returns only safe information
- Frontend consumes endpoint on load
- Branding cached for performance
- No authentication required
- Must be careful what information is exposed

---

## Decision 21: Singleton Pattern for PlatformBrandConfig

**Date:** Sprint 3.6  
**Status**: Active (with TODO for multi-brand)

### Problem

Platform branding is global, only one configuration needed.

### Decision

Use singleton pattern (ID=1) for PlatformBrandConfig.

### Justification

- **Simple:** Only one configuration to manage
- **Performance:** Cache single record
- **Clear:** Easy to understand
- **Sufficient:** For single-brand platform

### Consequences

- Only one platform brand configuration
- Repository assumes ID=1
- TODO: Add brand_key for multi-brand support
- Not suitable for multi-brand platform without changes

---

## Decision 22: Dynamic Issuer in JWT

**Date:** Sprint 3.7  
**Status**: Active

### Problem

JWT issuer should reflect platform name for white-label support.

### Decision

Use dynamic issuer from PlatformBrandConfig in JWT generation.

### Justification

- **White-label:** Issuer reflects platform name
- **Dynamic:** Changes with platform branding
- **Standard:** JWT issuer should identify issuer
- **Flexible:** Can change without code deployment
- **Security:** Still cryptographically signed

### Consequences

- Isser from PlatformBrandConfig
- JWT generation uses platform name
- Issuer changes with branding
- Tokens from different brands have different issuers

---

## Decision 23: Dynamic PlatformName in Services

**Date:** Sprint 3.7  
**Status**: Active

### Problem

Services need to use platform name dynamically for white-label support.

### Decision

Pass platformName to EmailService and BackupService for dynamic usage.

### Justification

- **White-label:** Services reflect platform name
- **Dynamic:** Changes with platform branding
- **Email:** Email templates use platform name
- **Backup:** Backup filenames use platform name
- **Flexible:** Can change without code deployment

### Consequences

- EmailService accepts platformName
- BackupService accepts platformName
- Email templates use dynamic name
- Backup filenames use dynamic name
- Services initialized with platform branding

---

## Decision 24: Goose for Migrations

**Date:** Initial Development  
**Status**: Active

### Problem

Need a database migration tool that works with Go and supports multiple databases.

### Decision

Use Goose for database migrations.

### Justification

- **Go-native:** Written in Go
- **Multi-database:** Supports SQLite and Oracle
- **Versioned:** Timestamp-based versioning
- **Up/Down:** Supports rollback
- **Standard:** Well-established in Go community

### Consequences

- Migrations in SQL files
- Timestamp-based naming
- Up and Down migrations required
- Must be idempotent
- Down must completely revert Up

---

## Decision 25: Naming Conventions

**Date:** Initial Development  
**Status**: Active

### Problem

Need consistent naming conventions across codebase for maintainability.

### Decision

Establish naming conventions: Go (PascalCase/camelCase), Database (snake_case), API (kebab-case/camelCase).

### Justification

- **Consistency:** Uniform naming across codebase
- **Standards:** Follow language/framework conventions
- **Readability:** Easy to read and understand
- **Tooling:** Tools expect these conventions
- **Best Practice:** Industry standard

### Consequences

- Go: PascalCase for structs/methods, camelCase for variables
- Database: snake_case for tables/columns
- API: kebab-case for routes, camelCase for JSON
- Must follow conventions strictly
- Linters enforce conventions

---

## Decision 26: Secure Impersonation for Platform Admins

**Date:** July 2026  
**Status:** Active

### Problem

Platform administrators need to temporarily access company dashboards to troubleshoot issues and provide support without:
- Requiring company owner credentials
- Creating permanent user accounts
- Bypassing security and audit trails
- Breaking existing authentication and authorization systems

### Decision

Implement secure impersonation using JWT claims with temporary sessions and comprehensive audit logging.

### Justification

- **Security:** No password sharing or credential exposure
- **Auditability:** All actions during impersonation are logged with platform admin identity
- **Temporary:** Sessions are time-limited (24 hours) and can be ended immediately
- **Non-invasive:** Doesn't alter existing authentication/authorization architecture
- **Clear Indication:** Frontend banner shows impersonation status to users
- **RBAC Integration:** Platform admins get Owner permissions only during impersonation

### Consequences

- JWT claims include impersonation flags (IsImpersonating, OriginalPlatformUserID)
- New impersonation_audit table for comprehensive logging
- Middleware updated to recognize and handle impersonation context
- RBAC middleware grants Owner permissions during impersonation
- Frontend banner component shows impersonation status
- Platform-only feature (tenants cannot impersonate)
- All impersonation operations are auditable

---

## Decision 27: Professional Session Management Architecture

**Date:** July 2026  
**Status:** Active

### Problem

The HorizonGest lacked a formal session management policy, causing:
- Dashboard opening without session validation
- Old JWTs surviving after backend restart
- Stores persisting across company switches
- Active impersonation sessions remaining after logout
- 409 errors for "already existing impersonation session"
- Inconsistent company switching
- User insecurity due to unpredictable session behavior

### Decision

Implement a professional session management architecture with centralized managers and strict lifecycle policies.

### Justification

- **Centralization:** Single point of truth for session operations
- **Consistency:** Context always validated before navigation
- **Security:** Automatic 401 handling and session destruction
- **Prevention:** Flags prevent race conditions
- **Professional:** Equivalent to Stripe, Notion, GitHub, Slack, Linear
- **Maintainability:** Easy to test and evolve
- **User Trust:** Predictable and secure session behavior

### Consequences

- **SessionManager:** Centralized Platform and Tenant session management
- **TenantSessionManager:** Centralized company switching lifecycle
- **Validation:** Backend validation on app initialization
- **401 Handling:** Automatic session destruction on 401
- **Stale Sessions:** Auto-cleanup of impersonations older than 24 hours
- **Cache Invalidation:** Complete cache clearing on company switch
- **Store Reset:** All stores have clear/reset methods
- **No Hacks:** No setTimeout, no reload, no scattered clearing
- **Single Point:** All session operations go through managers

---

## Decision 28: Storage Keys Centralization and Browser Compatibility

**Date:** July 24, 2026  
**Status:** Active

### Problem

The codebase had scattered storage key strings throughout, causing:
- Typos in cookie and localStorage keys
- Difficulty refactoring when keys needed to change
- No documentation of storage key purposes
- Use of experimental Navigation API not supported by all browsers
- Broad cache clearing (`sessionStorage.clear()`) affecting platform data

### Decision

Centralize all storage keys in a constants file and implement granular cache clearing with browser-compatible APIs.

### Justification

- **Type Safety:** Constants prevent typos and provide TypeScript type safety
- **Maintainability:** Single source of truth for all storage keys
- **Documentation:** Each key's purpose is documented
- **Refactoring:** Easy to change keys in one place
- **Browser Compatibility:** Removed experimental Navigation API
- **Granular Clearing:** Only tenant-specific data is cleared, preserving platform session
- **Cross-browser:** Uses only APIs supported by Chrome, Firefox, Safari, Edge

### Consequences

- **New File:** `frontend/src/lib/constants/storage-keys.ts` with all cookie/localStorage/sessionStorage keys
- **Rule:** Never use string literals for storage keys
- **Rule:** Never use `sessionStorage.clear()` or `localStorage.clear()` globally
- **Rule:** Only use APIs supported by all target browsers
- **Rule:** Never use `@ts-ignore` to silence TypeScript errors
- **Implementation:** Granular methods like `clearTenantSessionStorage()` remove only specific keys
- **Error Handling:** Typed error classes (InfrastructureError, SessionValidationError, BackendError, UIError)
- **Documentation:** Comprehensive JSDoc comments on cache lifecycle methods

### Files Changed

- `frontend/src/lib/constants/storage-keys.ts` (new)
- `frontend/src/lib/managers/sessionManager.ts`
- `frontend/src/lib/managers/tenantSessionManager.ts`
- `frontend/src/routes/platform/signin/+page.svelte`
- `frontend/src/routes/platform/admin/+page.svelte`
- `frontend/src/routes/platform/companies/+page.svelte`
- `frontend/src/routes/platform/companies/[id]/+page.svelte`
- `frontend/src/routes/platform/companies/[id]/owner/+page.svelte`
- `docs/05-development/SESSION_MANAGEMENT.md`
- `docs/05-development/SESSION_TESTING.md`
- `docs/AI_CONTEXT.md`

---

## Decision 29: Backend-Driven Impersonation Management

**Date:** July 24, 2026  
**Status:** Active

### Problem

The previous implementation returned HTTP 409 "já existe uma sessão de impersonation ativa" when a platform admin tried to switch companies with an active impersonation less than 24 hours old. This contradicted the HorizonGest architecture, which was designed to allow platform admins to freely switch between companies without logout or manual impersonation termination.

**Previous Behavior:**
- If impersonation active AND > 24h: Auto-end and create new ✅
- If impersonation active AND < 24h: Return HTTP 409 ❌

**Expected Flow:**
Platform → Company A → Platform → Company B → Platform → Company C (without any 409 errors)

### Decision

Make the backend fully responsible for impersonation management. The endpoint `POST /api/platform/impersonation/start` is now idempotent and automatically ends any active impersonation before creating a new one, regardless of age.

**New Behavior:**
- If impersonation active (any age): Auto-end and create new ✅
- No HTTP 409 for normal company switching

### Justification

- **User Experience:** Platform admins can freely switch companies without manual intervention
- **Simplicity:** Frontend no longer needs to call `endPreviousImpersonation()` before `enterCompany()`
- **Idempotency:** Endpoint is now idempotent for the same platform admin
- **Architectural Consistency:** Aligns with HorizonGest's design goal of seamless company switching
- **Audit Integrity:** Previous impersonation receives EndedAt correctly, new one is created
- **No Race Conditions:** Backend ensures atomic operation (end previous + create new)

### Consequences

- **Backend Changes:**
  - `impersonation_service.go`: Removed age check, auto-end any active impersonation
  - `impersonation_handler.go`: Removed `ErrAlreadyImpersonating` error handling
  - Endpoint no longer returns HTTP 409 for active impersonations

- **Frontend Changes:**
  - `tenantSessionManager.ts`: Removed `endPreviousImpersonation()` call from `enterCompany()`
  - Simplified flow: just call `requestTenantJWT()`, backend handles the rest

- **HTTP 409 Usage:**
  - Only for truly invalid scenarios (non-existent user, company, permissions, deactivated company, business rule violations)
  - Never for "already has active impersonation"

- **Audit Guarantees:**
  - Previous impersonation always receives EndedAt
  - New impersonation is always created
  - Never two active impersonations for same PlatformUser
  - History remains intact

### Files Changed

- `backend/internal/service/impersonation_service.go`
- `backend/internal/handler/impersonation_handler.go`
- `frontend/src/lib/managers/tenantSessionManager.ts`
- `docs/05-development/SESSION_MANAGEMENT.md`

### Testing Validation

Required flow validation:
- Company A → Company B → Company C → Company D → Company A → Company B
- No HTTP 409 errors
- No logout required
- No browser restart required
- No cookie clearing required

---

## Summary

| Decision | Status | Impact |
|----------|--------|--------|
| SQLite/Oracle | Active | Database abstraction |
| GORM | Active | ORM choice |
| CompanyID | Active | Multi-tenancy |
| Feature Flags | Active | Dynamic module control |
| Module Registry | Active | Module tracking |
| PlatformBrandConfig/GlobalConfig | Active | Config separation |
| Layered Architecture | Active | Code organization |
| Branding Separation | Active | White-label support |
| Cache in Repository | Active | Performance |
| Soft Delete | Active | Data preservation |
| JWT | Active | Authentication |
| RBAC | Active | Authorization |
| Chi | Active | HTTP router |
| SvelteKit | Active | Frontend framework |
| Product Live vs Sold | Active | Historical accuracy |
| Active vs DeletedAt | Active | Availability vs removal |
| DDD Principles | Active | Business-first |
| Environment Variables | Active | Secret management |
| APP_VERSION | Active | Version tracking |
| Public Branding Endpoint | Active | White-label frontend |
| Singleton PlatformBrand | Active | Simple branding (TODO multi-brand) |
| Dynamic JWT Issuer | Active | White-label JWT |
| Dynamic PlatformName | Active | White-label services |
| Goose | Active | Migrations |
| Naming Conventions | Active | Code consistency |
| Secure Impersonation | Active | Platform admin access to companies |
| Professional Session Management | Active | Centralized session lifecycle |
| Storage Keys Centralization | Active | Type-safe storage management |
| Backend-Driven Impersonation | Active | Idempotent company switching |

---

**Last Updated:** July 24, 2026 - Backend-Driven Impersonation Management
