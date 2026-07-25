import { test, expect } from '@playwright/test';

test.describe('Cenário 5: Fechar navegador → Abrir → Login obrigatório', () => {
	test('should require login after closing and reopening browser', async ({ page, context }) => {
		// Login first
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/platform');

		// Close browser context (simulates closing browser)
		await context.close();

		// Create new context and page (simulates reopening browser)
		const newContext = await page.context().browser()?.newContext();
		const newPage = await newContext?.newPage();

		if (newPage) {
			// Try to navigate to a protected route
			await newPage.goto('/dashboard');

			// Should be redirected to login (per HorizonGest policy)
			await expect(newPage).toHaveURL('/login');
			await expect(newPage.locator('h1')).toContainText('Login');

			// Verify no session cookies exist
			const cookies = await newContext.cookies();
			const platformCookie = cookies.find(c => c.name === 'platform_auth_token');
			const tenantCookie = cookies.find(c => c.name === 'auth_token');
			
			expect(platformCookie).toBeUndefined();
			expect(tenantCookie).toBeUndefined();

			await newContext.close();
		}
	});
});
