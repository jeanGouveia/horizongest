import { describe, it, expect, beforeEach, vi } from 'vitest';
import { userStore, companyStore, rbacStore, themeStore, brandStore, toast, goto } from './mocks/stores';

// Mock the imports BEFORE importing TenantSessionManager
vi.mock('$lib/stores/userStore.svelte', () => ({ userStore }));
vi.mock('$lib/stores/companyStore.svelte', () => ({ companyStore }));
vi.mock('$lib/stores/rbacStore.svelte', () => ({ rbacStore }));
vi.mock('$lib/stores/themeStore.svelte', () => ({ themeStore }));
vi.mock('$lib/stores/brandStore.ts', () => ({ brandStore }));
vi.mock('$lib/stores/toast', () => ({ toast }));
vi.mock('$app/navigation', () => ({ goto }));

// Import TenantSessionManager after mocks are set up
import { tenantSessionManager } from '$lib/managers/tenantSessionManager';

describe('TenantSessionManager Unit Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		// Reset manager state
		(tenantSessionManager as any).currentCompanyId = null;
		(tenantSessionManager as any).isEntering = false;
		(tenantSessionManager as any).isLeaving = false;
		// Clear cookies and localStorage
		document.cookie.split(';').forEach(c => {
			document.cookie = c.replace(/^ +/, '').replace(/=.*/, '=;expires=' + new Date().toUTCString() + ';path=/');
		});
		localStorage.clear();
		sessionStorage.clear();
	});

	describe('enterCompany', () => {
		it('should enter company successfully', async () => {
			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ token: 'test-token', companyName: 'Test Company' })
			}) as any;

			const result = await tenantSessionManager.enterCompany(1);

			expect(result.success).toBe(true);
			expect((tenantSessionManager as any).currentCompanyId).toBe(1);
		});

		it('should prevent entry if already entering', async () => {
			(tenantSessionManager as any).isEntering = true;

			const result = await tenantSessionManager.enterCompany(1);

			expect(result.success).toBe(false);
			expect(result.error).toBe('Já está entrando em uma empresa');
		});

		it('should prevent entry if already in company', async () => {
			(tenantSessionManager as any).currentCompanyId = 1;

			const result = await tenantSessionManager.enterCompany(1);

			expect(result.success).toBe(false);
			expect(result.error).toBe('Já está nesta empresa');
		});

		it('should clear stores before entering', async () => {
			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ token: 'test-token', companyName: 'Test Company' })
			}) as any;

			await tenantSessionManager.enterCompany(1);

			expect(userStore.logout).toHaveBeenCalled();
			expect(companyStore.clear).toHaveBeenCalled();
			expect(rbacStore.reset).toHaveBeenCalled();
			expect(themeStore.clear).toHaveBeenCalled();
			expect(brandStore.clear).toHaveBeenCalled();
			expect(toast.clear).toHaveBeenCalled();
		});

		it('should return error if token fetch fails', async () => {
			global.fetch = vi.fn().mockResolvedValue({
				ok: false,
				json: async () => ({ error: 'Token inválido' })
			}) as any;

			const result = await tenantSessionManager.enterCompany(1);

			expect(result.success).toBe(false);
			expect(result.error).toContain('Erro ao iniciar impersonation');
		});

		it('should end previous impersonation before entering', async () => {
			document.cookie = 'platform_auth_token=platform_token';
			localStorage.setItem('impersonation', JSON.stringify({ isImpersonating: true }));

			const fetchMock = vi.fn().mockImplementation((url: string) => {
				if (url.includes('end')) {
					return Promise.resolve({ ok: true });
				}
				return Promise.resolve({
					ok: true,
					json: async () => ({ token: 'test-token', companyName: 'Test Company' })
				});
			});
			global.fetch = fetchMock as any;

			await tenantSessionManager.enterCompany(1);

			expect(fetchMock).toHaveBeenCalledWith(
				expect.stringContaining('impersonation/end'),
				expect.objectContaining({
					method: 'POST'
				})
			);
		});
	});

	describe('leaveCompany', () => {
		it('should leave company successfully', async () => {
			(tenantSessionManager as any).currentCompanyId = 1;

			const result = await tenantSessionManager.leaveCompany();

			expect(result.success).toBe(true);
			expect((tenantSessionManager as any).currentCompanyId).toBeNull();
		});

		it('should prevent leaving if already leaving', async () => {
			(tenantSessionManager as any).isLeaving = true;

			const result = await tenantSessionManager.leaveCompany();

			expect(result.success).toBe(false);
			expect(result.error).toBe('Já está saindo da empresa');
		});

		it('should clear stores when leaving', async () => {
			await tenantSessionManager.leaveCompany();

			expect(userStore.logout).toHaveBeenCalled();
			expect(companyStore.clear).toHaveBeenCalled();
			expect(rbacStore.reset).toHaveBeenCalled();
			expect(themeStore.clear).toHaveBeenCalled();
			expect(brandStore.clear).toHaveBeenCalled();
			expect(toast.clear).toHaveBeenCalled();
		});

		it('should clear tenant cookie when leaving', async () => {
			document.cookie = 'auth_token=tenant_token';
			(tenantSessionManager as any).currentCompanyId = 1;

			await tenantSessionManager.leaveCompany();

			expect(document.cookie).not.toContain('auth_token');
		});

		it('should clear localStorage when leaving', async () => {
			localStorage.setItem('impersonation', JSON.stringify({ isImpersonating: true }));
			(tenantSessionManager as any).currentCompanyId = 1;

			await tenantSessionManager.leaveCompany();

			expect(localStorage.getItem('impersonation')).toBeNull();
		});
	});

	describe('switchCompany', () => {
		it('should switch company successfully', async () => {
			(tenantSessionManager as any).currentCompanyId = 1;

			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ token: 'test-token', companyName: 'Test Company' })
			}) as any;

			const result = await tenantSessionManager.switchCompany(2);

			expect(result.success).toBe(true);
			expect((tenantSessionManager as any).currentCompanyId).toBe(2);
		});

		it('should fail if leaveCompany fails', async () => {
			(tenantSessionManager as any).isLeaving = true; // Force failure

			const result = await tenantSessionManager.switchCompany(2);

			expect(result.success).toBe(false);
			expect(result.error).toContain('saindo da empresa');
		});

		it('should fail if enterCompany fails', async () => {
			(tenantSessionManager as any).isEntering = true; // Force failure

			const result = await tenantSessionManager.switchCompany(2);

			expect(result.success).toBe(false);
			expect(result.error).toContain('entrando em uma empresa');
		});

		it('should clear stores twice during switch', async () => {
			(tenantSessionManager as any).currentCompanyId = 1;

			global.fetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ token: 'test-token', companyName: 'Test Company' })
			}) as any;

			await tenantSessionManager.switchCompany(2);

			// Clear once for leave, once for enter
			expect(userStore.logout).toHaveBeenCalledTimes(2);
			expect(companyStore.clear).toHaveBeenCalledTimes(2);
			expect(rbacStore.reset).toHaveBeenCalledTimes(2);
			expect(themeStore.clear).toHaveBeenCalledTimes(2);
			expect(brandStore.clear).toHaveBeenCalledTimes(2);
			expect(toast.clear).toHaveBeenCalledTimes(2);
		});
	});

	describe('destroy', () => {
		it('should clear all stores', async () => {
			await tenantSessionManager.destroy();

			expect(userStore.logout).toHaveBeenCalled();
			expect(companyStore.clear).toHaveBeenCalled();
			expect(rbacStore.reset).toHaveBeenCalled();
			expect(themeStore.clear).toHaveBeenCalled();
			expect(brandStore.clear).toHaveBeenCalled();
			expect(toast.clear).toHaveBeenCalled();
		});

		it('should clear tenant cookie', async () => {
			document.cookie = 'auth_token=tenant_token';

			await tenantSessionManager.destroy();

			expect(document.cookie).not.toContain('auth_token');
		});

		it('should clear localStorage', async () => {
			localStorage.setItem('impersonation', 'test-data');

			await tenantSessionManager.destroy();

			expect(localStorage.getItem('impersonation')).toBeNull();
		});
	});

	describe('getCurrentCompanyId', () => {
		it('should return current company ID', () => {
			(tenantSessionManager as any).currentCompanyId = 1;

			expect(tenantSessionManager.getCurrentCompanyId()).toBe(1);
		});

		it('should return null if not in company', () => {
			(tenantSessionManager as any).currentCompanyId = null;

			expect(tenantSessionManager.getCurrentCompanyId()).toBeNull();
		});
	});

	describe('isInCompany', () => {
		it('should return true if in company', () => {
			(tenantSessionManager as any).currentCompanyId = 1;

			expect(tenantSessionManager.isInCompany()).toBe(true);
		});

		it('should return false if not in company', () => {
			(tenantSessionManager as any).currentCompanyId = null;

			expect(tenantSessionManager.isInCompany()).toBe(false);
		});
	});
});
