# Stores Documentation

**HorizonGest Platform - State Management**

---

## Overview

HorizonGest uses Svelte 5 Runes for global state management.

---

## Available Stores

### brandStore

**Purpose:** Dynamic branding data

**State:**
- `platformName` - Platform name
- `primaryColor` - Primary color
- `logo` - Logo URL
- And more branding fields

**Usage:**
```typescript
import { platformName, primaryColor } from '$lib/stores/brandStore';
```

### userStore

**Purpose:** Authentication state

**State:**
- `user` - Current user
- `authenticated` - Authentication status

**Usage:**
```typescript
import { user, authenticated } from '$lib/stores/userStore';
```

---

## Derived Stores

Derived stores can be created for computed values using Svelte 5 Runes.

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
