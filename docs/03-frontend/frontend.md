# Frontend Documentation

**HorizonGest Platform - Frontend Overview**

---

## Overview

The HorizonGest frontend is built with SvelteKit and TypeScript. It provides a modern, responsive user interface for restaurant management.

---

## Architecture

**Framework:** SvelteKit
**Language:** TypeScript
**UI:** Svelte 5 Runes
**Icons:** Lucide Svelte
**Styling:** CSS (TailwindCSS planned)

---

## Directory Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api/               # API client
│   │   ├── components/        # UI components
│   │   ├── stores/            # Svelte 5 Runes stores
│   │   └── types/             # TypeScript types
│   └── routes/                # Pages and layouts
└── package.json
```

---

## Key Concepts

### Stores

Global state management with Svelte 5 Runes:
- `brandStore` - Dynamic branding
- `userStore` - Authentication

### Routing

File-based routing with layouts for shared UI.

### API Client

Server-side API proxy to Go backend via `hooks.server.ts`.

### Branding

All branding from `/api/public/brand` (dynamic, no hardcoded).

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
