# Permissions Documentation

**HorizonGest Platform - RBAC Overview**

---

## Overview

HorizonGest uses Role-Based Access Control (RBAC) for authorization. Permissions are granular and verified in the Service Layer.

---

## Roles

### Platform Roles

- **Platform Admin:** Full platform access
- **Platform Manager:** Platform configuration
- **Platform Viewer:** Read-only platform access

### Tenant Roles

- **Admin:** Full company access
- **Manager:** User and settings management
- **Employee:** Orders and inventory
- **Viewer:** Read-only access

---

## Permissions

Permissions are defined per operation and verified in the Service Layer.

---

## Authorization Flow

1. User authenticates
2. Token includes role
3. Middleware checks authentication
4. Service checks permissions
5. Operation allowed or denied

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
