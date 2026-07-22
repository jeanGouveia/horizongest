# Backend Documentation

**HorizonGest Platform - Backend Overview**

---

## Overview

The HorizonGest backend is built with Go 1.21+ following a strict layered architecture. It provides RESTful APIs for the frontend and handles all business logic.

---

## Architecture

**Pattern:** Layered Architecture

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
```

---

## Directory Structure

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── domain/                 # Pure Go structs
│   ├── ports/                  # Interfaces
│   ├── service/                # Business logic
│   ├── handler/                # HTTP handlers
│   ├── middleware/             # HTTP middleware
│   ├── infra/
│   │   ├── database/           # Database connection
│   │   └── repository/         # GORM implementations
│   └── util/                  # Utilities
└── migrations/                 # SQL migrations
```

---

## Key Components

### Handlers

HTTP handlers handle request/response, validation, and call services.

### Services

Services contain all business logic and call repositories.

### Repositories

Repositories handle data access with GORM and implement caching.

### Domain

Domain entities are pure Go structs without infrastructure dependencies.

---

## Tech Stack

- **Language:** Go 1.21+
- **Framework:** Chi (HTTP router)
- **ORM:** GORM
- **Database:** SQLite (dev), Oracle (prod)
- **Auth:** JWT (HS256)
- **Migrations:** Goose

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
