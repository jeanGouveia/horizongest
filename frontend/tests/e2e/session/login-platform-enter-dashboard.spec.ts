import { test, expect } from '@playwright/test';

test.describe('Cenário 1: Login Platform → Enter Company → Dashboard', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to login page
		await page.goto('/login');
	});

	test('should login to platform, enter company and reach dashboard', async ({ page }) => {
		// Login to platform
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');

		// Should be redirected to platform dashboard
		await expect(page).toHaveURL('/platform');
		await expect(page.locator('h1')).toContainText('Platform Dashboard');

		// Enter company
		await page.click('[data-testid="company-1"]');
		
		// Should be redirected to company dashboard
		await expect(page).toHaveURL('/dashboard');
		await expect(page.locator('h1')).toContainText('Company Dashboard');

		// Verify session cookies are set
		const cookies = await page.context().cookies();
		const platformCookie = cookies.find(c => c.name === 'platform_auth_token');
		const tenantCookie = cookies.find(c => c.name === 'auth_token');
		
		expect(platformCookie).toBeDefined();
		expect(tenantCookie).toBeDefined();
	});
});
