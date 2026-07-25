import { writable, derived } from 'svelte/store';
import { userStore } from './userStore.svelte';

export interface BrandConfig {
	platformName: string;
	platformShortName: string;
	website: string;
	logoPath: string;
	faviconPath: string;
	logoLight: string;
	logoDark: string;
	icon: string;
	loginBackground: string;
	loginIllustration: string;
	copyright: string;
	primaryColor: string;
	secondaryColor: string
}

const defaultBrand: BrandConfig = {
	platformName: 'HorizonGest',
	platformShortName: 'Horizon',
	website: 'https://horizongest.com',
	logoPath: '/assets/platform/logo.svg',
	faviconPath: '/assets/platform/favicon.ico',
	logoLight: '',
	logoDark: '',
	icon: '',
	loginBackground: '',
	loginIllustration: '',
	copyright: '© 2024 HorizonGest Inc. All rights reserved.',
	primaryColor: '#0f172a',
	secondaryColor: '#6366f1'
};

function createBrandStore() {
	const { subscribe, set, update } = writable<BrandConfig>(defaultBrand);

	return {
		subscribe,
		load: async () => {
			try {
				const response = await fetch('/api/public/brand');
				if (response.ok) {
					const data = await response.json();
					set(data);
				}
			} catch (error) {
				console.error('Failed to load brand config:', error);
				// Keep default brand on error
			}
		},
		reload: async () => {
			await brandStore.load();
		},
		setCompanyName: (companyName: string) => {
			update(brand => ({
				...brand,
				platformName: companyName,
				platformShortName: companyName
			}));
		},
		clear: () => {
			set(defaultBrand);
		}
	};
}

export const brandStore = createBrandStore();

// Derived stores for commonly used values
export const platformName = derived(brandStore, ($brand) => $brand.platformName);
export const platformShortName = derived(brandStore, ($brand) => $brand.platformShortName);
export const copyright = derived(brandStore, ($brand) => $brand.copyright);
export const primaryColor = derived(brandStore, ($brand) => $brand.primaryColor);
export const secondaryColor = derived(brandStore, ($brand) => $brand.secondaryColor);
export const logoPath = derived(brandStore, ($brand) => $brand.logoPath || $brand.logoLight);
export const faviconPath = derived(brandStore, ($brand) => $brand.faviconPath);
