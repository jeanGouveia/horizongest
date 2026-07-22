# Multi-Tenant Documentation

**HorizonGest Platform - Multi-Tenancy**

---

## Overview

HorizonGest is a multi-tenant SaaS platform where each company (tenant) has isolated data.

---

## CompanyID Pattern

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

---

## Tenant Isolation

### Repository Layer

All repository queries filter by CompanyID for tenant entities.

### Middleware

Tenant middleware verifies CompanyID in JWT token.

### Platform vs Tenant

- Platform users cannot access tenant data
- Tenant users cannot access platform data
- Separate JWT secrets
- Separate authentication middleware

---

## Tenant Configuration

Each tenant can configure:
- Branding (logo, colors, theme)
- Settings (timezone, locale)
- Users (invite, roles, permissions)
- Feature access (based on plan)

---

## Plans

Companies are assigned to plans that define:
- Feature limits
- User limits
- Module access

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
