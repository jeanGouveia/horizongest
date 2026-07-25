import { test, expect } from '@playwright/test';

test.describe('Cenário 4: Backend reiniciado → 401 → Tela Login', () => {
	test('should handle backend restart with 401 and redirect to login', async ({ page, context }) => {
		// Login first
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/platform');

		// Simulate backend restart by clearing server session
		// (This would typically be done via a test endpoint)
		await page.goto('/api/test/clear-server-session');

		// Try to navigate to a protected route
		await page.goto('/dashboard');

		// Should receive 401 and be redirected to login
		await expect(page).toHaveURL('/login');

		// Verify session cookies are cleared
		const cookies = await context.cookies();
		const platformCookie = cookies.find(c => c.name === 'platform_auth_token');
		const tenantCookie = cookies.find(c => c.name === 'auth_token');
		
		expect(platformCookie).toBeUndefined();
		expect(tenantCookie).toBeUndefined();

		// Verify localStorage is cleared
		const impersonation = await page.evaluate(() => localStorage.getItem('impersonation'));
		expect(impersonation).toBeNull();
	});
});
