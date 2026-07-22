# Entities Reference

**HorizonGest Platform - Entity Reference**

---

## Global Entities (No CompanyID)

### PlatformUser
Platform administrator users.

### PlatformSession
Platform authentication sessions.

### PlatformAudit
Platform audit logs.

### Plan
Subscription plans for companies.

### PlatformBrandConfig
Platform branding configuration.

### GlobalConfig
Global technical configuration.

---

## Tenant Entities (With CompanyID)

### Company
Tenant companies.

### User
Tenant users.

### Product
Menu products.

### Order
Customer orders.

### Ingredient
Inventory ingredients.

---

## Entity Fields

### Standard Fields

All entities include:
- `id` - Primary key
- `company_id` - Tenant ID (except global entities)
- `created_at` - Creation timestamp
- `updated_at` - Update timestamp
- `deleted_at` - Soft delete timestamp

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
