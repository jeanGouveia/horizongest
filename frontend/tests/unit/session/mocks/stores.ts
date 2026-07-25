import { vi } from 'vitest';

// Mock userStore
export const userStore = {
	logout: vi.fn(),
	setUser: vi.fn(),
	setLoading: vi.fn()
};

// Mock companyStore
export const companyStore = {
	clear: vi.fn(),
	setCompany: vi.fn()
};

// Mock rbacStore
export const rbacStore = {
	reset: vi.fn(),
	load: vi.fn()
};

// Mock themeStore
export const themeStore = {
	clear: vi.fn(),
	load: vi.fn()
};

// Mock brandStore
export const brandStore = {
	clear: vi.fn(),
	load: vi.fn()
};

// Mock toast
export const toast = {
	clear: vi.fn()
};

// Mock navigation
export const goto = vi.fn();
