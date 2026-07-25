import { test, expect } from '@playwright/test';

test.describe('Cenário 3: Troca Empresa (A → B) com verificação de contexto', () => {
	test.beforeEach(async ({ page }) => {
		// Login and enter company A
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('should switch from company A to company B and verify context', async ({ page }) => {
		// Verify company A context
		await expect(page.locator('[data-testid="company-name"]')).toHaveText('Company A');
		await expect(page.locator('[data-testid="branding-logo"]')).toHaveAttribute('src', /company-a/);

		// Switch to company B
		await page.click('[data-testid="switch-company-button"]');
		await page.click('[data-testid="company-2"]');

		// Verify company B context
		await expect(page).toHaveURL('/dashboard');
		await expect(page.locator('[data-testid="company-name"]')).toHaveText('Company B');
		await expect(page.locator('[data-testid="branding-logo"]')).toHaveAttribute('src', /company-b/);

		// Verify user is still the same
		await expect(page.locator('[data-testid="user-email"]')).toHaveText('platform@example.com');

		// Verify dashboard is for company B
		await expect(page.locator('[data-testid="dashboard-content"]')).toContainText('Company B Dashboard');
	});
});
