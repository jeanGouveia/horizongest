# Coding Standards

**HorizonGest Platform - Coding Standards**

---

## Go Code

### Naming Conventions

- **Structs:** PascalCase (User, Product)
- **Methods:** PascalCase (GetUser, CreateProduct)
- **Variables:** camelCase (userID, productName)
- **Constants:** PascalCase or UPPER_CASE (MaxRetries, MAX_RETRIES)
- **Packages:** lowercase (service, repository)

### Code Style

- Use `gofmt` for formatting
- Use `golint` for linting
- Add godoc comments for public functions
- Keep functions focused and small

### Error Handling

- Return errors, don't panic
- Use specific errors (ErrUserNotFound)
- Propagate errors to Handler
- Log errors in Service Layer

---

## TypeScript Code

### Naming Conventions

- **Classes/Interfaces:** PascalCase
- **Functions/Variables:** camelCase
- **Constants:** UPPER_CASE
- **Files:** kebab-case

### Code Style

- Use `eslint` for linting
- Use `prettier` for formatting
- Add JSDoc comments for public functions
- Keep functions focused and small

---

## General Standards

- Write clear, readable code
- Add comments for complex logic
- Follow existing patterns
- Write tests for new code
- Update documentation

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
