import { test, expect } from '@playwright/test';

test.describe('Cenário 9: Voltar para Plataforma → Entrar novamente → Nova empresa', () => {
	test.beforeEach(async ({ page }) => {
		// Login and enter company
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('should destroy all previous context when returning to platform and entering new company', async ({ page }) => {
		// Verify company 1 context
		await expect(page.locator('[data-testid="company-name"]')).toHaveText('Company A');
		
		// Leave company (return to platform)
		await page.click('[data-testid="leave-company-button"]');
		await expect(page).toHaveURL('/platform');
		
		// Verify tenant session is cleared
		const cookies = await page.context().cookies();
		const tenantCookie = cookies.find(c => c.name === 'auth_token');
		expect(tenantCookie).toBeUndefined();
		
		const impersonation = await page.evaluate(() => localStorage.getItem('impersonation'));
		expect(impersonation).toBeNull();

		// Enter new company (company 2)
		await page.click('[data-testid="company-2"]');
		await expect(page).toHaveURL('/dashboard');
		
		// Verify company 2 context
		await expect(page.locator('[data-testid="company-name"]')).toHaveText('Company B');
		await expect(page.locator('[data-testid="branding-logo"]')).toHaveAttribute('src', /company-b/);
		
		// Verify no company 1 data remains
		const companyData = await page.evaluate(() => {
			return window.__companyStore__;
		});
		expect(companyData?.companyId).toBe(2);
		expect(companyData?.companyName).toBe('Company B');
		
		// Verify stores are clean (only company 2 data)
		await expect(page.locator('[data-testid="dashboard-content"]')).toContainText('Company B Dashboard');
		await expect(page.locator('[data-testid="dashboard-content"]')).not.toContainText('Company A');
	});
});
