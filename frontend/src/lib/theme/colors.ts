/**
 * Paleta de cores do PratoOnline - Redesign Visual
 * Inspirada em produtos profissionais como Stripe, Linear e Vercel
 * Sistema de tokens de cores para consistência visual
 */

export const colors = {
	// Cor primária - Índigo sofisticado (inspirado em Stripe)
	primary: {
		50: '#eef2ff',
		100: '#e0e7ff',
		200: '#c7d2fe',
		300: '#a5b4fc',
		400: '#818cf8',
		500: '#6366f1',
		600: '#4f46e5',
		700: '#4338ca',
		800: '#3730a3',
		900: '#312e81',
		950: '#1e1b4b',
	},

	// Cores de sucesso - Verde esmeralda (mais sofisticado)
	success: {
		50: '#ecfdf5',
		100: '#d1fae5',
		200: '#a7f3d0',
		300: '#6ee7b7',
		400: '#34d399',
		500: '#10b981',
		600: '#059669',
		700: '#047857',
		800: '#065f46',
		900: '#064e3b',
		950: '#022c22',
	},

	// Cores de erro - Vermelho coral (mais moderno)
	error: {
		50: '#fef2f2',
		100: '#fee2e2',
		200: '#fecaca',
		300: '#fca5a5',
		400: '#f87171',
		500: '#ef4444',
		600: '#dc2626',
		700: '#b91c1c',
		800: '#991b1b',
		900: '#7f1d1d',
		950: '#450a0a',
	},

	// Cores de aviso - Âmbar dourado (mais elegante)
	warning: {
		50: '#fffbeb',
		100: '#fef3c7',
		200: '#fde68a',
		300: '#fcd34d',
		400: '#fbbf24',
		500: '#f59e0b',
		600: '#d97706',
		700: '#b45309',
		800: '#92400e',
		900: '#78350f',
		950: '#451a03',
	},

	// Cores neutras - Cinza ardesia (inspirado em Linear)
	neutral: {
		50: '#f8fafc',
		100: '#f1f5f9',
		200: '#e2e8f0',
		300: '#cbd5e1',
		400: '#94a3b8',
		500: '#64748b',
		600: '#475569',
		700: '#334155',
		800: '#1e293b',
		900: '#0f172a',
		950: '#020617',
	},

	// Cores de fundo - Tons mais suaves e modernos
	background: {
		default: '#ffffff',
		secondary: '#f8fafc',
		tertiary: '#f1f5f9',
		elevated: '#ffffff',
		surface: '#ffffff',
	},

	// Cores de texto - Hierarquia visual clara
	text: {
		primary: '#0f172a',
		secondary: '#475569',
		tertiary: '#64748b',
		inverse: '#ffffff',
		disabled: '#94a3b8',
	},

	// Cores de borda - Mais sutis e elegantes
	border: {
		default: '#e2e8f0',
		light: '#f1f5f9',
		dark: '#1e293b',
		focus: '#6366f1',
	},

	// Cores de acento - Para destacar elementos importantes
	accent: {
		blue: '#3b82f6',
		purple: '#8b5cf6',
		pink: '#ec4899',
		cyan: '#06b6d4',
	},
} as const;

export type ColorName = keyof typeof colors;
export type ColorShade = 50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900 | 950;
