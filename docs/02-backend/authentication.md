# Authentication Documentation

**HorizonGest Platform - Authentication Overview**

---

## Overview

HorizonGest uses JWT (JSON Web Tokens) for authentication. Separate JWT secrets are used for platform and tenant authentication.

---

## JWT Configuration

### Platform JWT

**Secret:** `JWT_PLATFORM_SECRET` (environment variable)
**Expiration:** 24 hours
**Issuer:** Dynamic from PlatformBrandConfig

### Tenant JWT

**Secret:** `JWT_TENANT_SECRET` (environment variable)
**Expiration:** 24 hours
**Issuer:** Dynamic from PlatformBrandConfig

---

## Authentication Flow

1. User submits credentials
2. Backend validates credentials
3. Backend generates JWT token
4. Token includes UserID and CompanyID
5. Token is returned to client
6. Client includes token in Authorization header
7. Backend validates token on each request

---

## Token Blacklist

Tokens are blacklisted on logout to prevent reuse.

---

## Middleware

Authentication middleware validates JWT tokens on protected routes.

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
