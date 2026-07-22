# Deployment Documentation

**HorizonGest Platform - Deployment Guide**

---

## Environment Variables

### Required

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=pratodb
JWT_PLATFORM_SECRET=
JWT_TENANT_SECRET=
APP_VERSION=1.0.0
PORT=8080
```

### Optional

```bash
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
EMAIL_ENABLED=false
```

---

## Database Migration

### Development

```bash
cd backend
goose -dir migrations sqlite "app.db" up
```

### Production

```bash
cd backend
goose -dir migrations oracle "$ORACLE_CONNECTION_STRING" up
```

---

## Build

### Backend

```bash
cd backend
go build -o horizon-gest cmd/server/main.go
```

### Frontend

```bash
cd frontend
npm run build
```

---

## Run

### Backend

```bash
./horizon-gest
```

### Frontend

```bash
cd frontend
npm run preview
```

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
