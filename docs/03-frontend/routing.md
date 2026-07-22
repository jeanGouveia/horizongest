# Routing Documentation

**HorizonGest Platform - Frontend Routing**

---

## Overview

HorizonGest uses SvelteKit's file-based routing system.

---

## Route Structure

```
src/routes/
├── (auth)/                    # Auth layout
│   ├── signin/
│   ├── forgot-password/
│   └── reset-password/
├── (platform)/                # Platform layout
│   ├── dashboard/
│   ├── menu/
│   ├── orders/
│   └── inventory/
└── +layout.svelte             # Root layout
```

---

## Layouts

### Auth Layout

Used for authentication pages (signin, forgot-password, reset-password).

### Platform Layout

Used for authenticated pages (dashboard, menu, orders, inventory).

---

## Route Protection

Route protection is implemented via `layout.server.ts` to verify authentication.

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
