import { describe, it, expect, beforeEach, vi } from 'vitest';
import { userStore, companyStore, rbacStore, themeStore, brandStore, toast, goto } from './mocks/stores';

// Mock the imports BEFORE importing SessionManager
vi.mock('$lib/stores/userStore.svelte', () => ({ userStore }));
vi.mock('$lib/stores/companyStore.svelte', () => ({ companyStore }));
vi.mock('$lib/stores/rbacStore.svelte', () => ({ rbacStore }));
vi.mock('$lib/stores/themeStore.svelte', () => ({ themeStore }));
vi.mock('$lib/stores/brandStore.ts', () => ({ brandStore }));
vi.mock('$lib/stores/toast', () => ({ toast }));
vi.mock('$app/navigation', () => ({ goto }));

// Import SessionManager after mocks are set up
import { sessionManager } from '$lib/managers/sessionManager';

describe('SessionManager Unit Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		// Clear cookies and localStorage
		document.cookie.split(';').forEach(c => {
			document.cookie = c.replace(/^ +/, '').replace(/=.*/, '=;expires=' + new Date().toUTCString() + ';path=/');
		});
		localStorage.clear();
		sessionStorage.clear();
	});

	describe('validateSession', () => {
		it('should return no session when no tokens exist', async () => {
			const result = await sessionManager.validateSession();
			expect(result.valid).toBe(false);
			expect(result.sessionType).toBe('none');
		});

		it('should return invalid when tenant session exists without platform session', async () => {
			document.cookie = 'auth_token=tenant_token_value';
			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(false);
			expect(result.sessionType).toBe('none');
			expect(result.error).toContain('Tenant Session sem Platform Session');
		});

		it('should return platform session when platform token is valid', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			
			global.fetch = vi.fn(() =>
				Promise.resolve({
					ok: true,
					json: () => Promise.resolve({ id: 1, email: 'test@example.com' })
				})
			) as any;

			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(true);
			expect(result.sessionType).toBe('platform');
		});

		it('should return invalid when platform token is invalid (401)', async () => {
			document.cookie = 'platform_auth_token=invalid_token';
			
			global.fetch = vi.fn(() =>
				Promise.resolve({
					ok: false,
					status: 401
				})
			) as any;

			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(false);
			expect(result.sessionType).toBe('none');
		});

		it('should return invalid when backend returns error', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			
			global.fetch = vi.fn(() =>
				Promise.resolve({
					ok: false,
					status: 500
				})
			) as any;

			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(false);
			expect(result.sessionType).toBe('none');
		});

		it('should return invalid when backend is unavailable', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			
			global.fetch = vi.fn(() =>
				Promise.reject(new Error('Network error'))
			) as any;

			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(false);
			expect(result.sessionType).toBe('none');
		});

		it('should return tenant session when both tokens are valid', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			document.cookie = 'auth_token=tenant_token_value';
			
			global.fetch = vi.fn((url: string) => {
				if (url.includes('platform')) {
					return Promise.resolve({
						ok: true,
						json: () => Promise.resolve({ id: 1, email: 'test@example.com' })
					});
				}
				return Promise.resolve({
					ok: true,
					json: () => Promise.resolve({ id: 1, email: 'test@example.com' })
				});
			}) as any;

			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(true);
			expect(result.sessionType).toBe('tenant');
		});

		it('should return platform session when tenant token is invalid', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			document.cookie = 'auth_token=invalid_tenant_token';
			
			global.fetch = vi.fn((url: string) => {
				if (url.includes('platform')) {
					return Promise.resolve({
						ok: true,
						json: () => Promise.resolve({ id: 1, email: 'test@example.com' })
					});
				}
				return Promise.resolve({
					ok: false,
					status: 401
				});
			}) as any;

			const result = await sessionManager.validateSession();
			
			expect(result.valid).toBe(true);
			expect(result.sessionType).toBe('platform');
		});

		it('should prevent concurrent validation', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			
			global.fetch = vi.fn(() =>
				Promise.resolve({
					ok: true,
					json: () => Promise.resolve({ id: 1, email: 'test@example.com' })
				})
			) as any;

			const promise1 = sessionManager.validateSession();
			const promise2 = sessionManager.validateSession();
			
			const [result1, result2] = await Promise.all([promise1, promise2]);
			
			expect(result1.valid).toBe(true);
			expect(result2.valid).toBe(false);
			expect(result2.error).toContain('Já está validando sessão');
		});
	});

	describe('logout', () => {
		it('should destroy all sessions and redirect to login', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			document.cookie = 'auth_token=tenant_token_value';
			localStorage.setItem('impersonation', JSON.stringify({ isImpersonating: true }));
			
			global.fetch = vi.fn(() =>
				Promise.resolve({
					ok: true
				})
			) as any;

			const result = await sessionManager.logout();
			
			expect(result.success).toBe(true);
			expect(document.cookie).not.toContain('platform_auth_token');
			expect(document.cookie).not.toContain('auth_token');
			expect(localStorage.getItem('impersonation')).toBeNull();
		});

		it('should call end impersonation when tenant session exists', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			document.cookie = 'auth_token=tenant_token_value';
			
			const fetchMock = vi.fn(() =>
				Promise.resolve({
					ok: true
				})
			);
			global.fetch = fetchMock as any;

			await sessionManager.logout();
			
			expect(fetchMock).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({
					method: 'POST',
					headers: expect.objectContaining({
						'Authorization': expect.stringContaining('platform_auth_token')
					})
				})
			);
		});

		it('should handle errors gracefully', async () => {
			global.fetch = vi.fn(() =>
				Promise.reject(new Error('Network error'))
			) as any;

			const result = await sessionManager.logout();
			
			expect(result.success).toBe(false);
			expect(result.error).toBeDefined();
		});
	});

	describe('destroyAllSessions', () => {
		it('should clear all stores and caches', async () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			document.cookie = 'auth_token=tenant_token_value';
			localStorage.setItem('impersonation', JSON.stringify({ isImpersonating: true }));
			sessionStorage.setItem('test', 'value');

			await sessionManager.destroyAllSessions();

			expect(userStore.logout).toHaveBeenCalled();
			expect(companyStore.clear).toHaveBeenCalled();
			expect(rbacStore.reset).toHaveBeenCalled();
			expect(themeStore.clear).toHaveBeenCalled();
			expect(brandStore.clear).toHaveBeenCalled();
			expect(toast.clear).toHaveBeenCalled();
			expect(document.cookie).not.toContain('platform_auth_token');
			expect(document.cookie).not.toContain('auth_token');
			expect(localStorage.getItem('impersonation')).toBeNull();
			expect(sessionStorage.getItem('test')).toBeNull();
		});
	});

	describe('destroyTenantSession', () => {
		it('should clear only tenant-specific data', async () => {
			document.cookie = 'auth_token=tenant_token_value';
			document.cookie = 'platform_auth_token=platform_token_value';
			localStorage.setItem('impersonation', JSON.stringify({ isImpersonating: true }));

			await sessionManager.destroyTenantSession();

			expect(userStore.logout).toHaveBeenCalled();
			expect(companyStore.clear).toHaveBeenCalled();
			expect(rbacStore.reset).toHaveBeenCalled();
			expect(themeStore.clear).toHaveBeenCalled();
			expect(document.cookie).not.toContain('auth_token');
			expect(document.cookie).toContain('platform_auth_token');
			expect(localStorage.getItem('impersonation')).toBeNull();
		});
	});

	describe('destroyPlatformSession', () => {
		it('should clear only platform-specific data', () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			document.cookie = 'auth_token=tenant_token_value';

			sessionManager.destroyPlatformSession();

			expect(document.cookie).not.toContain('platform_auth_token');
			expect(document.cookie).toContain('auth_token');
		});
	});

	describe('hasActiveSession', () => {
		it('should return true when platform session exists', () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			expect(sessionManager.hasActiveSession()).toBe(true);
		});

		it('should return true when tenant session exists', () => {
			document.cookie = 'auth_token=tenant_token_value';
			expect(sessionManager.hasActiveSession()).toBe(true);
		});

		it('should return false when no session exists', () => {
			expect(sessionManager.hasActiveSession()).toBe(false);
		});

		it('should return true when hasPlatformSession', () => {
			document.cookie = 'platform_auth_token=platform_token_value';
			expect(sessionManager.hasPlatformSession()).toBe(true);
		});

		it('should return true when hasTenantSession', () => {
			document.cookie = 'auth_token=tenant_token_value';
			expect(sessionManager.hasTenantSession()).toBe(true);
		});
	});
});
