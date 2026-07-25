# Session Testing Documentation

## Overview

This document describes the testing strategy for the HorizonGest Session Management system. The SessionManager and TenantSessionManager are critical components that require comprehensive test coverage to prevent regressions and ensure security.

## Testing Strategy

We use a two-tier testing approach:

1. **Unit Tests (Vitest)** - Fast, isolated tests for individual manager methods
2. **End-to-End Tests (Playwright)** - Real browser tests that validate the complete user flow

## Tools

### Unit Tests

- **Vitest** - Fast unit test framework with native TypeScript support
- **jsdom** - DOM environment for browser API simulation
- **vi** - Vitest's built-in mocking library

### E2E Tests

- **Playwright** - Cross-browser E2E testing framework
- **Chromium, Firefox, WebKit** - Multi-browser support

## How to Execute

### Unit Tests

```bash
# Run unit tests in watch mode
npm run test

# Run unit tests once
npm run test:run

# Run with coverage
npm run test:run -- --coverage
```

### E2E Tests

```bash
# Run E2E tests
npm run test:e2e

# Run E2E tests with UI
npm run test:e2e:ui

# Run specific test file
npx playwright test tests/e2e/session/login-platform-enter-dashboard.spec.ts
```

## Test Structure

```
frontend/
├── tests/
│   ├── unit/
│   │   └── session/
│   │       ├── setup.ts                 # Global test setup
│   │       ├── mocks/
│   │       │   └── stores.ts            # Mock implementations
│   │       ├── sessionManager.test.ts    # SessionManager unit tests
│   │       └── tenantSessionManager.test.ts # TenantSessionManager unit tests
│   └── e2e/
│       └── session/
│           ├── login-platform-enter-dashboard.spec.ts
│           ├── logout-login-required.spec.ts
│           ├── switch-company-verify-context.spec.ts
│           ├── backend-restarted-401-login.spec.ts
│           ├── close-browser-reopen-login.spec.ts
│           ├── 100-switches-stress.spec.ts
│           ├── logout-during-impersonation.spec.ts
│           ├── user-switch.spec.ts
│           └── platform-return-enter-new-company.spec.ts
├── vitest.config.ts
└── playwright.config.ts
```

## Unit Test Scenarios

### SessionManager Tests

#### validateSession
- ✅ No session when no tokens exist
- ✅ Invalid when tenant session exists without platform session
- ✅ Platform session when platform token is valid
- ✅ Invalid when platform token is invalid (401)
- ✅ Invalid when backend returns error
- ✅ Invalid when backend is unavailable
- ✅ Tenant session when both tokens are valid
- ✅ Platform session when tenant token is invalid
- ✅ Prevents concurrent validation

#### logout
- ✅ Destroys all sessions and redirects to login
- ✅ Calls end impersonation when tenant session exists
- ✅ Handles errors gracefully

#### destroyAllSessions
- ✅ Clears all stores and caches
- ✅ Clears cookies (platform and tenant)
- ✅ Clears localStorage and sessionStorage

#### destroyTenantSession
- ✅ Clears only tenant-specific data
- ✅ Preserves platform session

#### destroyPlatformSession
- ✅ Clears only platform-specific data
- ✅ Preserves tenant session

#### hasActiveSession
- ✅ Returns true when platform session exists
- ✅ Returns true when tenant session exists
- ✅ Returns false when no session exists
- ✅ hasPlatformSession works correctly
- ✅ hasTenantSession works correctly

### TenantSessionManager Tests

#### enterCompany
- ✅ Enters company successfully
- ✅ Prevents entry if already entering
- ✅ Prevents entry if already in company
- ✅ Clears stores before entering
- ✅ Returns error if token fetch fails
- ✅ Ends previous impersonation before entering

#### leaveCompany
- ✅ Leaves company successfully
- ✅ Prevents leaving if already leaving
- ✅ Clears stores when leaving
- ✅ Clears tenant cookie when leaving
- ✅ Clears localStorage when leaving

#### switchCompany
- ✅ Switches company successfully
- ✅ Fails if leaveCompany fails
- ✅ Fails if enterCompany fails
- ✅ Clears stores twice during switch

#### destroy
- ✅ Clears all stores
- ✅ Clears tenant cookie
- ✅ Clears localStorage

#### getCurrentCompanyId
- ✅ Returns current company ID
- ✅ Returns null if not in company

#### isInCompany
- ✅ Returns true if in company
- ✅ Returns false if not in company

## E2E Test Scenarios

### Cenário 1: Login Platform → Enter Company → Dashboard
- Validates complete login flow
- Verifies session cookies are set
- Confirms dashboard navigation

### Cenário 2: Logout → Login obrigatório
- Validates logout clears all session data
- Confirms redirect to login page
- Verifies cookies and localStorage are cleared

### Cenário 3: Troca Empresa (A → B) com verificação de contexto
- Verifies branding changes correctly
- Confirms company name updates
- Validates user remains the same
- Confirms dashboard shows correct company data

### Cenário 4: Backend reiniciado → 401 → Tela Login
- Simulates backend restart
- Validates 401 handling
- Confirms automatic session destruction
- Verifies redirect to login

### Cenário 5: Fechar navegador → Abrir → Login obrigatório
- Validates session persistence policy
- Confirms login is required after browser close
- Verifies no session data persists

### Cenário 6: 100 trocas consecutivas
- Stress tests company switching
- Verifies no store contamination
- Confirms no cache survival
- Validates no orphaned impersonation sessions
- Checks for no console errors

### Cenário 7: Logout durante impersonation
- Validates complete session end during impersonation
- Confirms all session data is cleared
- Verifies impersonation data is removed

### Cenário 8: Troca de usuário (Usuário 1 → Logout → Usuário 2)
- Validates user isolation
- Confirms no data from User 1 remains
- Verifies User 2 starts with clean state

### Cenário 9: Voltar para Plataforma → Entrar novamente → Nova empresa
- Validates platform return flow
- Confirms previous context is destroyed
- Verifies new company context is clean

## Coverage Goals

- **Minimum Coverage**: 90% of SessionManager and TenantSessionManager
- **Critical Flows**: 100% coverage for login, logout, company switch, and 401 handling
- **E2E Coverage**: All 9 critical scenarios covered

## Mocks

Unit tests use comprehensive mocks to isolate the managers from external dependencies:

### Mocked Dependencies
- `fetch` - HTTP requests
- `localStorage` - Browser storage
- `sessionStorage` - Session storage
- `document.cookie` - Cookie management
- `userStore` - User state store
- `companyStore` - Company state store
- `rbacStore` - RBAC state store
- `themeStore` - Theme state store
- `brandStore` - Branding state store
- `toast` - Toast notification system
- `goto` - Navigation function

## Limitations

### Unit Tests
- **No Real Browser**: Tests run in jsdom, not a real browser
- **No Network**: All network calls are mocked
- **No Real Backend**: Backend responses are simulated

### E2E Tests
- **Requires Running Backend**: Tests need a running backend server
- **Slower Execution**: E2E tests are slower than unit tests
- **Flaky Potential**: Browser tests can be flaky due to timing issues
- **Test Data**: Requires test accounts and companies to exist

## Best Practices

### Writing Unit Tests
1. Mock all external dependencies
2. Test both success and failure scenarios
3. Verify side effects (cookies, localStorage, stores)
4. Use descriptive test names
5. Keep tests isolated and independent

### Writing E2E Tests
1. Use data-testid attributes for reliable element selection
2. Wait for elements to be visible before interacting
3. Verify both positive and negative outcomes
4. Clean up test data after tests
5. Use page objects for complex interactions

## Maintenance

### Adding New Tests
1. Add unit tests first for new manager methods
2. Add E2E tests for new user flows
3. Update this documentation with new scenarios
4. Ensure coverage goals are met

### Debugging Failed Tests
1. Run unit tests with `--reporter=verbose`
2. Run E2E tests with `--debug` flag
3. Check browser console for errors
4. Verify backend is running and accessible
5. Check test data exists in database

## Continuous Integration

Unit tests should run on every commit. E2E tests should run on pull requests and before releases.

```yaml
# Example CI configuration
- npm run test:run
- npm run build
- npm run test:e2e
```

## Future Improvements

- [ ] Add visual regression tests for branding
- [ ] Add performance tests for session operations
- [ ] Add accessibility tests for login flows
- [ ] Add mobile-specific E2E tests
- [ ] Add API integration tests
- [ ] Add contract tests for backend endpoints

## Architectural Decisions (Sprint 2026-07-24)

### Storage Keys Centralization
**Decision:** Centralized all storage keys in `frontend/src/lib/constants/storage-keys.ts`

**Rationale:**
- Prevents typos in storage key strings
- Facilitates refactoring when keys need to change
- Documents the purpose of each storage key
- Provides type safety through TypeScript constants

**Implementation:**
```typescript
export const CookieKeys = {
  PLATFORM_TOKEN: 'platform_auth_token',
  TENANT_TOKEN: 'auth_token'
} as const;

export const TenantLocalStorageKeys = {
  IMPERSONATION: 'impersonation',
  TENANT_CONTEXT: 'tenant_context',
  // ...
} as const;
```

**Rule:** All code must import and use these constants instead of string literals.

### Granular Cache Clearing
**Decision:** Replaced `sessionStorage.clear()` with granular key removal

**Rationale:**
- `sessionStorage.clear()` removes ALL session data, including platform-specific data
- Granular clearing ensures only tenant-specific data is removed
- Prevents accidental removal of platform session data
- More predictable and maintainable

**Implementation:**
```typescript
private clearTenantSessionStorage(): void {
  sessionStorage.removeItem(TenantSessionStorageKeys.TENANT_NAVIGATION);
  sessionStorage.removeItem(TenantSessionStorageKeys.TENANT_FORMS);
  sessionStorage.removeItem(TenantSessionStorageKeys.TENANT_FILTERS);
}
```

**Rule:** Never use `sessionStorage.clear()` or `localStorage.clear()` globally.

### Browser Compatibility
**Decision:** Removed Navigation API Experimental usage

**Rationale:**
- `window.navigation.entries` is not supported by Firefox, Safari, and other browsers
- Experimental APIs can change or be removed without notice
- Using standard APIs ensures cross-browser compatibility
- Reduces risk of runtime errors in production

**Removed Code:**
 Previously attempted to use:
```typescript
// REMOVED - Not supported by all browsers
window.navigation?.entries?.forEach((entry: any) => {
  // Clear navigation cache
});
```

**Solution:** Granular clearing of localStorage and sessionStorage is sufficient for cache invalidation.

**Rule:** Only use APIs supported by all target browsers. Check caniuse.com before using new APIs.

### Differentiated Error Handling
**Decision:** Implemented typed error classes for better error handling

**Rationale:**
- Different error types require different user messages
- Infrastructure errors (network) need different handling than backend errors
- Session validation errors should be handled differently from UI errors
- Improves debugging and logging

**Implementation:**
```typescript
class SessionError extends Error {
  constructor(message: string, public type: 'session' | 'infrastructure' | 'backend' | 'ui') {
    super(message);
  }
}

class InfrastructureError extends SessionError { /* ... */ }
class SessionValidationError extends SessionError { /* ... */ }
class BackendError extends SessionError { /* ... */ }
class UIError extends SessionError { /* ... */ }
```

**Rule:** Use appropriate error types for different scenarios. Provide user-friendly messages based on error type.

### No @ts-ignore
**Decision:** Removed all `@ts-ignore` comments from the codebase

**Rationale:**
- `@ts-ignore` silences TypeScript errors without fixing the root cause
- Hides potential bugs and type safety issues
- Makes the codebase less maintainable
- Prevents TypeScript from doing its job

**Rule:** Fix TypeScript errors properly instead of using `@ts-ignore`.

### Cache Documentation
**Decision:** Added comprehensive documentation for cache lifecycle

**Rationale:**
- Clear documentation of what caches exist, who creates them, who destroys them
- Prevents confusion about cache ownership
- Helps future developers understand the system
- Reduces risk of cache-related bugs

**Implementation:** Added detailed JSDoc comments to cache-related methods explaining:
- Responsibilities
- What is cleared
- What is NOT cleared
- When the method is called

**Rule:** Document cache operations clearly in code comments.
