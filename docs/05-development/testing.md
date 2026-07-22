# Testing Documentation

**HorizonGest Platform - Testing Guide**

---

## Backend Testing

### Unit Tests

**Location:** `backend/internal/service/*_test.go`, `backend/internal/domain/*_test.go`

**Run:**
```bash
cd backend
go test ./internal/service/...
go test ./internal/domain/...
```

**Coverage Goals:**
- Services: 80%
- Domain: 95%

### Integration Tests

**Location:** `backend/internal/handler/*_test.go`

**Run:**
```bash
cd backend
go test ./internal/handler/...
```

---

## Frontend Testing

### Unit Tests

**Location:** `frontend/src/lib/components/*.test.ts`

**Run:**
```bash
cd frontend
npm test
```

---

## Test Guidelines

- Write tests for all new code
- Use mocks for external dependencies
- Test edge cases
- Keep tests fast and isolated
- Update tests when code changes

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
