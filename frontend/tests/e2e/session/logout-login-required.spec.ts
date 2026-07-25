import { test, expect } from '@playwright/test';

test.describe('Cenário 2: Logout → Login obrigatório', () => {
	test.beforeEach(async ({ page }) => {
		// Login first
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/platform');
	});

	test('should logout and require login again', async ({ page }) => {
		// Click logout
		await page.click('[data-testid="logout-button"]');

		// Should be redirected to login page
		await expect(page).toHaveURL('/login');
		await expect(page.locator('h1')).toContainText('Login');

		// Verify session cookies are cleared
		const cookies = await page.context().cookies();
		const platformCookie = cookies.find(c => c.name === 'platform_auth_token');
		const tenantCookie = cookies.find(c => c.name === 'auth_token');
		
		expect(platformCookie).toBeUndefined();
		expect(tenantCookie).toBeUndefined();

		// Verify localStorage is cleared
		const impersonation = await page.evaluate(() => localStorage.getItem('impersonation'));
		expect(impersonation).toBeNull();
	});
});
