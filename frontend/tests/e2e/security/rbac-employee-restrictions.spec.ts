import { test, expect } from '@playwright/test';

test.describe('RBAC: Employee Role Restrictions', () => {
	test.beforeEach(async ({ page }) => {
		// Login as Employee
		await page.goto('/login');
		await page.fill('input[name="email"]', 'employee@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		
		// Enter company
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('Employee should NOT be able to access company settings', async ({ page }) => {
		await page.goto('/company/settings');
		
		// Should be redirected or show 403
		await expect(page).toHaveURL(/.*403.*/);
	});

	test('Employee should NOT be able to manage users', async ({ page }) => {
		await page.goto('/company/users');
		
		// Should be redirected or show 403
		await expect(page).toHaveURL(/.*403.*/);
	});

	test('Employee should NOT be able to change user roles', async ({ page }) => {
		await page.goto('/company/users');
		
		// Should be redirected or show 403
		await expect(page).toHaveURL(/.*403.*/);
	});

	test('Employee should NOT be able to delete company', async ({ page }) => {
		await page.goto('/company/settings');
		
		// Should be redirected or show 403
		await expect(page).toHaveURL(/.*403.*/);
	});

	test('Employee should NOT be able to approve stock adjustments', async ({ page }) => {
		await page.goto('/stock-adjustments/pending');
		
		// Should be redirected or show 403
		await expect(page).toHaveURL(/.*403.*/);
	});

	test('Employee should be able to create orders', async ({ page }) => {
		await page.goto('/orders');
		await expect(page).toHaveURL('/orders');
		
		// Should see create order button
		await expect(page.locator('[data-testid="create-order-button"]')).toBeVisible();
	});

	test('Employee should be able to view products', async ({ page }) => {
		await page.goto('/products');
		await expect(page).toHaveURL('/products');
		
		// Should see products list
		await expect(page.locator('[data-testid="products-list"]')).toBeVisible();
	});

	test('Employee should NOT be able to create products', async ({ page }) => {
		await page.goto('/products');
		
		// Should NOT see create product button
		await expect(page.locator('[data-testid="create-product-button"]')).not.toBeVisible();
	});
});
