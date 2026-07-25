import { test, expect } from '@playwright/test';

test.describe('Cenário 7: Logout durante impersonation', () => {
	test.beforeEach(async ({ page }) => {
		// Login and enter company
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('should completely end session when logging out during impersonation', async ({ page, context }) => {
		// Verify impersonation is active
		const impersonationBefore = await page.evaluate(() => localStorage.getItem('impersonation'));
		expect(impersonationBefore).not.toBeNull();

		// Logout
		await page.click('[data-testid="logout-button"]');

		// Should be redirected to login
		await expect(page).toHaveURL('/login');

		// Verify all session data is cleared
		const cookies = await context.cookies();
		const platformCookie = cookies.find(c => c.name === 'platform_auth_token');
		const tenantCookie = cookies.find(c => c.name === 'auth_token');
		
		expect(platformCookie).toBeUndefined();
		expect(tenantCookie).toBeUndefined();

		const impersonationAfter = await page.evaluate(() => localStorage.getItem('impersonation'));
		expect(impersonationAfter).toBeNull();

		// Verify sessionStorage is cleared
		const sessionData = await page.evaluate(() => sessionStorage.getItem('test'));
		expect(sessionData).toBeNull();
	});
});
