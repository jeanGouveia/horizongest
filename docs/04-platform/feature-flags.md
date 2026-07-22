# Feature Flags Documentation

**HorizonGest Platform - Feature Flags**

---

## Overview

Feature flags allow dynamic enablement/disablement of platform modules without code deployment.

---

## Implementation

### Storage

`global_config` table

### Fields

- `enable_finance` - Enable finance module
- `enable_purchasing` - Enable purchasing module
- `enable_inventory` - Enable inventory module
- `enable_crm` - Enable CRM module
- `enable_calendar` - Enable calendar module
- `enable_pos` - Enable POS module
- `enable_ai` - Enable AI features
- `enable_delivery` - Enable delivery module
- `enable_marketplace` - Enable marketplace module

---

## Module Registry

### Purpose

Track available modules, dependencies, and status.

### Location

`internal/domain/module_registry.go`

### Sync

Must sync with feature flags in GlobalConfig.

---

## Usage

### Handler

Check feature flag before exposing routes.

### Service

Check feature flag before executing logic.

### Frontend

Check feature flag before showing UI.

---

## Benefits

- Gradual feature rollout
- A/B testing
- Controlling access to beta features
- Disabling problematic features

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
