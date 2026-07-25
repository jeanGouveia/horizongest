import { vi, beforeEach } from 'vitest';

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
	writable: true,
	value: vi.fn().mockImplementation(query => ({
		matches: false,
		media: query,
		onchange: null,
		addListener: vi.fn(),
		removeListener: vi.fn(),
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		dispatchEvent: vi.fn()
	}))
});

// Mock localStorage
const localStorageMock = (() => {
	let store: Record<string, string> = {};

	return {
		getItem: (key: string) => store[key] || null,
		setItem: (key: string, value: string) => {
			store[key] = value.toString();
		},
		removeItem: (key: string) => {
			delete store[key];
		},
		clear: () => {
			store = {};
		}
	};
})();

Object.defineProperty(window, 'localStorage', {
	value: localStorageMock
});

// Mock sessionStorage
const sessionStorageMock = (() => {
	let store: Record<string, string> = {};

	return {
		getItem: (key: string) => store[key] || null,
		setItem: (key: string, value: string) => {
			store[key] = value.toString();
		},
		removeItem: (key: string) => {
			delete store[key];
		},
		clear: () => {
			store = {};
		}
	};
})();

Object.defineProperty(window, 'sessionStorage', {
	value: sessionStorageMock
});

// Mock document.cookie
const cookieMock = (() => {
	let cookies: string = '';

	return {
		get: () => cookies,
		set: (cookie: string) => {
			cookies = cookie;
		},
		clear: () => {
			cookies = '';
		}
	};
})();

Object.defineProperty(document, 'cookie', {
	get: () => cookieMock.get(),
	set: (value: string) => {
		cookieMock.set(value);
	},
	configurable: true
});

// Reset before each test
beforeEach(() => {
	cookieMock.clear();
	localStorageMock.clear();
	sessionStorageMock.clear();
	vi.clearAllMocks();
});
