import { test, expect } from '@playwright/test';

test.describe('Cenário 8: Troca de usuário (Usuário 1 → Logout → Usuário 2)', () => {
	test('should ensure no data from User 1 remains after User 2 login', async ({ page, context }) => {
		// Login as User 1
		await page.goto('/login');
		await page.fill('input[name="email"]', 'user1@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/platform');
		
		// Enter company as User 1
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
		
		// Verify User 1 data
		await expect(page.locator('[data-testid="user-email"]')).toHaveText('user1@example.com');
		
		// Logout User 1
		await page.click('[data-testid="logout-button"]');
		await expect(page).toHaveURL('/login');

		// Verify User 1 session is cleared
		const cookiesAfterLogout = await context.cookies();
		const platformCookieAfterLogout = cookiesAfterLogout.find(c => c.name === 'platform_auth_token');
		const tenantCookieAfterLogout = cookiesAfterLogout.find(c => c.name === 'auth_token');
		
		expect(platformCookieAfterLogout).toBeUndefined();
		expect(tenantCookieAfterLogout).toBeUndefined();

		const impersonationAfterLogout = await page.evaluate(() => localStorage.getItem('impersonation'));
		expect(impersonationAfterLogout).toBeNull();

		// Login as User 2
		await page.fill('input[name="email"]', 'user2@example.com');
		await page.fill('input[name="password"]', 'password456');
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/platform');

		// Verify User 2 data (no User 1 data)
		await expect(page.locator('[data-testid="user-email"]')).toHaveText('user2@example.com');
		
		// Verify no User 1 company data
		const companyData = await page.evaluate(() => {
			return window.__companyStore__;
		});
		expect(companyData).toBeNull(); // Should be null since User 2 hasn't entered a company yet
	});
});
