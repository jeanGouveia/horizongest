import { test, expect } from '@playwright/test';

test.describe('Cenário 6: 100 trocas consecutivas', () => {
	test.beforeEach(async ({ page }) => {
		// Login first
		await page.goto('/login');
		await page.fill('input[name="email"]', 'platform@example.com');
		await page.fill('input[name="password"]', 'password123');
		await page.click('button[type="submit"]');
		await page.click('[data-testid="company-1"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('should handle 100 consecutive company switches without errors', async ({ page }) => {
		// Perform 100 switches between company 1 and company 2
		for (let i = 1; i <= 100; i++) {
			const targetCompanyId = i % 2 === 0 ? 2 : 1;
			const targetCompanyName = targetCompanyId === 1 ? 'Company A' : 'Company B';
			
			// Switch company
			await page.click('[data-testid="switch-company-button"]');
			await page.click(`[data-testid="company-${targetCompanyId}"]`);
			
			// Verify correct company
			await expect(page).toHaveURL('/dashboard');
			await expect(page.locator('[data-testid="company-name"]')).toHaveText(targetCompanyName);
			
			// Verify no errors in console
			const errors: string[] = [];
			page.on('console', msg => {
				if (msg.type() === 'error') {
					errors.push(msg.text());
				}
			});
			
			expect(errors).toHaveLength(0);
		}

		// Verify final state
		await expect(page.locator('[data-testid="company-name"]')).toHaveText('Company A');
		
		// Verify no orphaned impersonation sessions
		const impersonation = await page.evaluate(() => localStorage.getItem('impersonation'));
		const impersonationData = impersonation ? JSON.parse(impersonation) : null;
		expect(impersonationData?.isImpersonating).toBe(true);
	});

	test('should verify no store contamination after 100 switches', async ({ page }) => {
		// Perform 100 switches
		for (let i = 1; i <= 100; i++) {
			const targetCompanyId = i % 2 === 0 ? 2 : 1;
			await page.click('[data-testid="switch-company-button"]');
			await page.click(`[data-testid="company-${targetCompanyId}"]`);
		}

		// Verify stores are clean (only current company data)
		const companyData = await page.evaluate(() => {
			// Access company store data
			return window.__companyStore__;
		});
		
		// Should only have data for current company
		expect(companyData?.companyId).toBe(1);
		expect(companyData?.companyName).toBe('Company A');
	});
});
