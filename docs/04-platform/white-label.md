# White-Label Documentation

**HorizonGest Platform - White-Label Support**

---

## Overview

HorizonGest is designed for white-label deployment, allowing the platform to be rebranded for different markets.

---

## Current State

### Backend

**Status:** 100% ready

- Dynamic branding via PlatformBrandConfig
- Public endpoint for frontend
- Services use dynamic platform name
- JWT issuer dynamic from platform branding

### Frontend

**Status:** 95% ready

- brandStore global
- Components updated
- TODO: Favicon, manifest, meta tags

---

## Multi-Brand Support

### Current Implementation

Singleton pattern (ID=1) for PlatformBrandConfig.

### Future Implementation

**TODO:**
- Add `brand_key` for multi-brand
- Map-based cache (keyed by brand)
- Middleware for brand selection by domain
- Repository support for multiple configurations

---

## Platform vs Tenant Branding

### Platform Branding

- Identity of the platform
- Configured in PlatformBrandConfig
- Public endpoint `/api/public/brand`

### Tenant Branding

- Identity of the customer
- Configured in Theme and BusinessProfile
- Authenticated endpoint

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
