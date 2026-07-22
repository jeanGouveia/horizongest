# Getting Started

**HorizonGest Platform - Setup Instructions**

---

## Prerequisites

- Go 1.21+
- Node.js 18+
- npm

---

## Backend Setup

### Install Go Dependencies

```bash
cd backend
go mod download
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

---

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

---

## Verification

1. Backend: Visit `http://localhost:8080/api/public/brand`
2. Frontend: Visit `http://localhost:5173`

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
