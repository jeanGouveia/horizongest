/**
 * Sistema de Dark Mode do PratoOnline - Sprint 9 PX
 * Tokens preparados para implementação futura
 * Cores adaptadas para baixo contraste e conforto visual
 */

export const darkMode = {
	// Cores de fundo - Adaptadas para dark mode
	background: {
		default: '#0f172a',      // Slate 900
		secondary: '#1e293b',    // Slate 800
		tertiary: '#334155',     // Slate 700
		elevated: '#1e293b',     // Slate 800
		surface: '#1e293b',      // Slate 800
	},

	// Cores de texto - Adaptadas para dark mode
	text: {
		primary: '#f8fafc',      // Slate 50
		secondary: '#cbd5e1',    // Slate 300
		tertiary: '#94a3b8',     // Slate 400
		inverse: '#0f172a',      // Slate 900
		disabled: '#64748b',     // Slate 500
	},

	// Cores de borda - Adaptadas para dark mode
	border: {
		default: '#334155',      // Slate 700
		light: '#475569',       // Slate 600
		dark: '#1e293b',        // Slate 800
		focus: '#818cf8',       // Indigo 400
	},

	// Sombras adaptadas para dark mode
	shadow: {
		xs: '0 1px 2px 0 rgb(0 0 0 / 0.3)',
		sm: '0 1px 3px 0 rgb(0 0 0 / 0.4), 0 1px 2px -1px rgb(0 0 0 / 0.4)',
		base: '0 4px 6px -1px rgb(0 0 0 / 0.4), 0 2px 4px -2px rgb(0 0 0 / 0.4)',
		md: '0 10px 15px -3px rgb(0 0 0 / 0.4), 0 4px 6px -4px rgb(0 0 0 / 0.4)',
		lg: '0 20px 25px -5px rgb(0 0 0 / 0.4), 0 8px 10px -6px rgb(0 0 0 / 0.4)',
		xl: '0 25px 50px -12px rgb(0 0 0 / 0.5)',
		'2xl': '0 50px 100px -20px rgb(0 0 0 / 0.5)',
		inner: 'inset 0 2px 4px 0 rgb(0 0 0 / 0.3)',
	},
} as const;

// Classe para ativar dark mode
export const darkModeClass = 'dark';

// Seletor para dark mode
export const darkModeSelector = '[data-theme="dark"]';

// Função para verificar preferência do sistema
export function prefersDarkMode(): boolean {
	if (typeof window === 'undefined') return false;
	return window.matchMedia('(prefers-color-scheme: dark)').matches;
}

// Função para aplicar dark mode
export function applyDarkMode(enabled: boolean): void {
	if (typeof document === 'undefined') return;
	document.documentElement.setAttribute('data-theme', enabled ? 'dark' : 'light');
}

// Função para alternar dark mode
export function toggleDarkMode(): void {
	if (typeof document === 'undefined') return;
	const current = document.documentElement.getAttribute('data-theme');
	const next = current === 'dark' ? 'light' : 'dark';
	applyDarkMode(next === 'dark');
}

// Hook para usar dark mode (para uso futuro em Svelte)
export function useDarkMode() {
	let isDark = $state(false);

	if (typeof window !== 'undefined') {
		isDark = prefersDarkMode();
		applyDarkMode(isDark);

		// Escutar mudanças de preferência do sistema
		const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
		const handler = (e: MediaQueryListEvent) => {
			isDark = e.matches;
			applyDarkMode(isDark);
		};

		mediaQuery.addEventListener('change', handler);

		// Cleanup (será implementado com onDestroy no componente)
		return {
			get isDark() { return isDark; },
			set isDark(value: boolean) {
				isDark = value;
				applyDarkMode(value);
			},
			toggle: () => { isDark = !isDark; },
		};
	}

	return {
		get isDark() { return isDark; },
		set isDark(value: boolean) { isDark = value; },
		toggle: () => { isDark = !isDark; },
	};
}
