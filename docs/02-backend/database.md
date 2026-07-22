# Database Documentation

**HorizonGest Platform - Database Overview**

---

## Overview

HorizonGest uses SQLite for development and Oracle for production. Database access is abstracted through GORM.

---

## Connection

### Development

**Database:** SQLite
**File:** `app.db`
**Location:** Project root

### Production

**Database:** Oracle
**Connection:** Configured via environment variables
**File:** `internal/infra/database/connection.go`

---

## Migrations

### Tool

**Tool:** Goose
**Location:** `backend/migrations/`

### Naming Convention

`YYYYMMDD_description.sql`

### Running Migrations

```bash
cd backend
goose -dir migrations sqlite "app.db" up
```

---

## Schema

### Tables

- `platform_users` - Platform administrators
- `platform_sessions` - Platform sessions
- `platform_audit` - Platform audit logs
- `companies` - Tenant companies
- `users` - Tenant users
- `products` - Menu products
- `orders` - Customer orders
- `ingredients` - Inventory ingredients
- `platform_brand_config` - Platform branding
- `global_config` - Global configuration
- `plans` - Subscription plans

---

## CompanyID Pattern

All tenant tables include `company_id` for multi-tenancy isolation.

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
