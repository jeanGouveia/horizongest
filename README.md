# HorizonGest - Full Stack Application

Full-stack SaaS application with Go backend and SvelteKit frontend for restaurant management.

## Prerequisites

- Go 1.21+
- Node.js 18+
- npm

## Backend Setup

### Install Go Dependencies

```bash
cd backend
go mod init github.com/seu-usuario/my-app/backend
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get golang.org/x/crypto
go get github.com/golang-jwt/jwt/v5
go get github.com/pressly/goose/v3
go get github.com/joho/godotenv
```

### Configure Environment

```bash
cp .env.example .env
# Edit .env with your settings
```

### Run Migrations

```bash
goose -dir migrations sqlite "app.db" up
```

### Run Backend

```bash
go run cmd/server/main.go
```

Backend will run on `http://localhost:8080`

## Frontend Setup

### Install Node Dependencies

```bash
cd frontend
npm install
```

### Run Frontend

```bash
npm run dev
```

Frontend will run on `http://localhost:5173`

## Project Structure

### Backend (Go)
- `cmd/server/main.go` - Entry point
- `internal/domain/` - Domain entities (pure Go structs)
- `internal/ports/` - Interfaces (contracts for decoupling)
- `internal/service/` - Business logic
- `internal/handler/` - HTTP handlers (Chi router)
- `internal/middleware/` - HTTP middleware
- `internal/infra/database/` - Database connection (single file to change for Oracle migration)
- `internal/infra/repository/` - GORM repository implementations
- `migrations/` - SQL migrations (Goose)

### Frontend (SvelteKit)
- `src/lib/stores/` - Svelte 5 Runes stores
- `src/lib/api/` - API client with credentials
- `src/routes/` - Pages and layouts
- `src/hooks.server.ts` - API proxy to Go backend

## Architecture Notes

### Backend
- Clean Architecture with ports/adapters pattern
- Domain layer is pure Go (no GORM tags)
- Repository interface isolates domain from infrastructure
- Single file (`connection.go`) to change for SQLite → Oracle migration

### Frontend
- SvelteKit with TypeScript
- Svelte 5 Runes for reactive state
- HttpOnly cookies for JWT storage
- Server-side API proxy to Go backend
- Route protection via layout.server.ts

## API Endpoints

- `POST /api/register` - Register new user
- `POST /api/login` - Login user
- `POST /api/logout` - Logout user
- `GET /api/me` - Get current user (protected)
