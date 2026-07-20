import { browser } from '$app/environment';
import { api } from '$lib/api/client';

export interface Theme {
	primary_color: string;
	secondary_color: string;
	logo_url: string;
	font_family: string;
	border_radius: string;
	is_default: boolean;
}

const DEFAULT_THEME: Theme = {
	primary_color: '#6366f1',
	secondary_color: '#4f46e5',
	logo_url: '',
	font_family: 'Inter',
	border_radius: '8px',
	is_default: true
};

class ThemeStore {
	theme = $state<Theme>(DEFAULT_THEME);
	loading = $state(false);
	error = $state<string | null>(null);

	constructor() {
		if (browser) {
			this.loadTheme();
		}
	}

	async loadTheme() {
		if (!browser) return;

		this.loading = true;
		this.error = null;

		try {
			const response = await api.theme.getTheme();
			if (!response.error && response.data) {
				this.theme = response.data;
				this.applyThemeToDOM();
			}
		} catch (e) {
			console.error('Failed to load theme:', e);
			const errorMessage = e instanceof Error ? e.message : String(e);
			if (e instanceof TypeError && errorMessage.includes('fetch')) {
				this.error = 'Erro de conexão ao carregar tema. Usando tema padrão.';
			} else {
				this.error = `Erro ao carregar tema: ${errorMessage || 'Erro desconhecido'}. Usando tema padrão.`;
			}
			// Fall back to default theme
			this.theme = DEFAULT_THEME;
			this.applyThemeToDOM();
		} finally {
			this.loading = false;
		}
	}

	async loadDefaultTheme() {
		if (!browser) return;

		this.loading = true;
		this.error = null;

		try {
			const response = await api.theme.getDefaultTheme();
			if (!response.error && response.data) {
				this.theme = response.data;
				this.applyThemeToDOM();
			}
		} catch (e) {
			console.error('Failed to load default theme:', e);
			this.error = 'Failed to load default theme';
			this.theme = DEFAULT_THEME;
			this.applyThemeToDOM();
		} finally {
			this.loading = false;
		}
	}

	private applyThemeToDOM() {
		if (!browser) return;

		const root = document.documentElement;

		// Apply primary color to CSS variables
		root.style.setProperty('--color-primary-500', this.theme.primary_color);
		root.style.setProperty('--color-primary-600', this.theme.secondary_color);
		
		// Generate color palette from primary color (simplified)
		root.style.setProperty('--color-primary-400', this.lightenColor(this.theme.primary_color, 10));
		root.style.setProperty('--color-primary-700', this.darkenColor(this.theme.secondary_color, 10));
		
		// Apply logo URL if available
		if (this.theme.logo_url) {
			root.style.setProperty('--brand-logo-url', `url(${this.theme.logo_url})`);
		}

		// Apply font family
		if (this.theme.font_family) {
			root.style.setProperty('--font-family-sans', this.theme.font_family);
		}

		// Apply border radius
		if (this.theme.border_radius) {
			root.style.setProperty('--radius-base', this.theme.border_radius);
			root.style.setProperty('--radius-md', this.theme.border_radius);
		}
	}

	private lightenColor(color: string, percent: number): string {
		const num = parseInt(color.replace('#', ''), 16);
		const amt = Math.round(2.55 * percent);
		const R = (num >> 16) + amt;
		const G = (num >> 8 & 0x00FF) + amt;
		const B = (num & 0x0000FF) + amt;
		return '#' + (0x1000000 + 
			(R < 255 ? R < 1 ? 0 : R : 255) * 0x10000 + 
			(G < 255 ? G < 1 ? 0 : G : 255) * 0x100 + 
			(B < 255 ? B < 1 ? 0 : B : 255)
		).toString(16).slice(1);
	}

	private darkenColor(color: string, percent: number): string {
		const num = parseInt(color.replace('#', ''), 16);
		const amt = Math.round(2.55 * percent);
		const R = (num >> 16) - amt;
		const G = (num >> 8 & 0x00FF) - amt;
		const B = (num & 0x0000FF) - amt;
		return '#' + (0x1000000 + 
			(R > 0 ? R > 255 ? 255 : R : 0) * 0x10000 + 
			(G > 0 ? G > 255 ? 255 : G : 0) * 0x100 + 
			(B > 0 ? B > 255 ? 255 : B : 0)
		).toString(16).slice(1);
	}

	get primaryColor(): string {
		return this.theme.primary_color;
	}

	get secondaryColor(): string {
		return this.theme.secondary_color;
	}

	get logoUrl(): string {
		return this.theme.logo_url;
	}

	get isDefault(): boolean {
		return this.theme.is_default;
	}
}

export const themeStore = new ThemeStore();
