# HorizonGest Documentation

**Official Documentation for HorizonGest Platform**

---

## Overview

This is the official documentation for HorizonGest, a multi-tenant SaaS platform for restaurant management. This documentation serves as the single source of truth for developers, AI systems, and maintainers.

---

## Quick Start

### For New Developers

1. Read [01-overview/vision.md](01-overview/vision.md) - Project vision and goals
2. Read [01-overview/architecture.md](01-overview/architecture.md) - Architecture overview
3. Read [AI_CONTEXT.md](AI_CONTEXT.md) - Complete context for code modifications
4. Read [05-development/getting-started.md](05-development/getting-started.md) - Setup instructions

### For AI Systems

1. Read [AI_CONTEXT.md](AI_CONTEXT.md) - Complete context for code modifications
2. Read [05-development/architecture-rules.md](05-development/architecture-rules.md) - Mandatory architecture rules
3. Read [DECISIONS.md](DECISIONS.md) - Architectural decisions history

### For Users

- [manuals/administrator.md](manuals/administrator.md) - Platform administrator guide
- [manuals/company.md](manuals/company.md) - Company owner guide
- [manuals/employee.md](manuals/employee.md) - Employee guide

---

## Documentation Structure

```
docs/
├── README.md                    # This file
├── AI_CONTEXT.md                # Complete context for AI systems
├── DECISIONS.md                 # Architectural decisions history
├── 01-overview/                 # Project overview
│   ├── vision.md               # Vision and goals
│   ├── glossary.md             # Domain glossary
│   ├── architecture.md         # Architecture overview
│   └── business-rules.md       # Business rules
├── 02-backend/                  # Backend documentation
│   ├── backend.md              # Backend overview
│   ├── api.md                  # API documentation
│   ├── database.md             # Database documentation
│   ├── authentication.md       # Authentication
│   └── permissions.md          # Permissions/RBAC
├── 03-frontend/                 # Frontend documentation
│   ├── frontend.md             # Frontend overview
│   ├── routing.md              # Routing documentation
│   ├── ui.md                   # UI guidelines
│   └── stores.md               # State management
├── 04-platform/                # Platform features
│   ├── branding.md             # Branding system
│   ├── configuration.md        # Configuration
│   ├── multi-tenant.md         # Multi-tenancy
│   ├── white-label.md          # White-label support
│   └── feature-flags.md        # Feature flags
├── 05-development/              # Development guide
│   ├── getting-started.md     # Setup instructions
│   ├── coding-standards.md     # Coding standards
│   ├── architecture-rules.md   # Architecture rules
│   ├── testing.md              # Testing guide
│   └── deployment.md           # Deployment guide
├── 06-reference/                # Reference materials
│   ├── entities.md             # Entity reference
│   ├── modules.md              # Module reference
│   ├── environment.md          # Environment variables
│   ├── roadmap.md              # Roadmap
│   └── tech-debt.md            # Technical debt
└── manuals/                     # User manuals
    ├── administrator.md        # Platform administrator
    ├── company.md              # Company owner
    └── employee.md             # Employee
```

---

## Key Documents

### AI_CONTEXT.md

**Purpose:** Complete context for AI systems to modify code.

**Contents:**
- Project overview
- Architecture rules
- Creating new modules
- Creating migrations
- Creating APIs
- Creating frontend screens
- Creating tests
- Common patterns
- Mandatory rules

**Read before:** Any code modification

### DECISIONS.md

**Purpose:** History of architectural decisions.

**Contents:**
- Problem
- Decision
- Justification
- Consequences

**Read before:** Questioning architectural choices

### 05-development/architecture-rules.md

**Purpose:** Mandatory architecture rules.

**Contents:**
- Layer separation
- Business logic location
- CompanyID requirements
- Branding rules
- Configuration rules
- Feature flags
- Cache rules
- Migration rules
- Frontend rules
- Naming conventions
- Error handling
- Testing rules
- Performance rules
- Security rules
- White-label rules

**Read before:** Any code modification

---

## Project Status

### Foundation

**Status:** ✅ **CLOSED**

**Score:** 9.0/10

**Completed:**
- Layered architecture
- Multi-tenancy
- Dynamic branding
- Feature flags
- Cache system
- Security (JWT, RBAC)
- Unit tests
- API documentation
- Security audit
- Performance audit
- White-label audit

**Next Steps:**
- Business modules
- UX improvements
- Integrations

### Phases

- **Fase 1:** Foundation ✅ (CLOSED)
- **Fase 2:** Documentation & Knowledge Base 🔄 (IN PROGRESS)
- **Fase 3:** Business Modules ⏳ (PENDING)

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

## Architecture

**Pattern:** Layered Architecture

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
```

**Key Principles:**
- Strict layer separation
- Business logic in Service Layer
- Multi-tenancy with CompanyID
- Dynamic branding
- Feature flags
- Cache in Repository Layer

---

## Contributing

### Code Changes

1. Read [AI_CONTEXT.md](AI_CONTEXT.md)
2. Read [05-development/architecture-rules.md](05-development/architecture-rules.md)
3. Follow existing patterns
4. Write tests
5. Update documentation

### Documentation Changes

1. Keep documentation up to date
2. Use clear, concise language
3. Include examples
4. Update related documents
5. Review for consistency

---

## Support

### Documentation Issues

If you find documentation issues:
- Update the document directly
- Add TODOs for missing information
- Report critical issues to the team

### Code Issues

If you find code issues:
- Follow the architecture rules
- Create tests to reproduce
- Fix the root cause
- Update documentation if needed

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
