import { test, expect } from '@playwright/test';

test.describe('RBAC: Owner Role Restrictions', () => {
	test.beforeEach(async ({ page }) => {
		// Login as Owner
		await page.goto('/login');
		await page.fill('input[name="email"]', 'owner@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		
		// Enter company
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('Owner should be able to access company settings', async ({ page }) => {
		await page.goto('/company/settings');
		await expect(page).toHaveURL('/company/settings');
		await expect(page.locator('h1')).toContainText('Company Settings');
	});

	test('Owner should be able to manage users', async ({ page }) => {
		await page.goto('/company/users');
		await expect(page).toHaveURL('/company/users');
		await expect(page.locator('h1')).toContainText('User Management');
		
		// Should see add user button
		await expect(page.locator('[data-testid="add-user-button"]')).toBeVisible();
	});

	test('Owner should be able to change user roles', async ({ page }) => {
		await page.goto('/company/users');
		await page.click('[data-testid="user-1"]');
		await page.click('[data-testid="change-role-button"]');
		
		// Should see role change modal
		await expect(page.locator('[data-testid="role-change-modal"]')).toBeVisible();
	});

	test('Owner should be able to delete company', async ({ page }) => {
		await page.goto('/company/settings');
		await page.click('[data-testid="delete-company-button"]');
		
		// Should see confirmation modal
		await expect(page.locator('[data-testid="delete-confirmation-modal"]')).toBeVisible();
	});

	test('Owner should be able to approve stock adjustments', async ({ page }) => {
		await page.goto('/stock-adjustments/pending');
		await expect(page).toHaveURL('/stock-adjustments/pending');
		
		// Should see approve button
		await expect(page.locator('[data-testid="approve-button"]')).toBeVisible();
	});
});
