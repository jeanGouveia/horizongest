/**
 * Sistema de tipografia do horizongest - Redesign Visual
 * Tokens de fontes para consistência visual
 * Inspirado em sistemas modernos como Stripe, Linear e Vercel
 */

export const typography = {
	// Fontes - Inter para UI, JetBrains Mono para código
	fontFamily: {
		sans: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
		mono: 'JetBrains Mono, "Fira Code", Consolas, Monaco, monospace',
	},

	// Tamanhos de fonte - Escala mais refinada
	fontSize: {
		xs: '0.75rem',     // 12px
		sm: '0.875rem',    // 14px
		base: '1rem',      // 16px
		lg: '1.125rem',    // 18px
		xl: '1.25rem',     // 20px
		'2xl': '1.5rem',   // 24px
		'3xl': '1.875rem', // 30px
		'4xl': '2.25rem',  // 36px
		'5xl': '3rem',     // 48px
		'6xl': '3.75rem',  // 60px
		'7xl': '4.5rem',   // 72px
		'8xl': '6rem',     // 96px
		'9xl': '8rem',     // 128px
	},

	// Altura de linha - Mais refinada para melhor legibilidade
	lineHeight: {
		none: '1',
		tight: '1.2',
		snug: '1.35',
		normal: '1.5',
		relaxed: '1.625',
		loose: '2',
	},

	// Espessura de fonte - Escala completa
	fontWeight: {
		thin: '100',
		extralight: '200',
		light: '300',
		normal: '400',
		medium: '500',
		semibold: '600',
		bold: '700',
		extrabold: '800',
		black: '900',
	},

	// Espaçamento de letra - Para títulos e UI
	letterSpacing: {
		tighter: '-0.05em',
		tight: '-0.025em',
		normal: '0em',
		wide: '0.025em',
		wider: '0.05em',
		widest: '0.1em',
	},
} as const;

export type FontSizeToken = keyof typeof typography.fontSize;
export type FontWeightToken = keyof typeof typography.fontWeight;

// Tipografia semântica - Hierarquia visual clara e profissional
export const semanticTypography = {
	// Títulos de página
	h1: {
		fontSize: typography.fontSize['3xl'],
		fontWeight: typography.fontWeight.bold,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.tight,
	},
	h2: {
		fontSize: typography.fontSize['2xl'],
		fontWeight: typography.fontWeight.semibold,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	h3: {
		fontSize: typography.fontSize.xl,
		fontWeight: typography.fontWeight.semibold,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	h4: {
		fontSize: typography.fontSize.lg,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	h5: {
		fontSize: typography.fontSize.base,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	h6: {
		fontSize: typography.fontSize.sm,
		fontWeight: typography.fontWeight.semibold,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	// Texto de corpo
	body: {
		fontSize: typography.fontSize.base,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
	},
	bodyLarge: {
		fontSize: typography.fontSize.lg,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
	},
	bodySmall: {
		fontSize: typography.fontSize.sm,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
	},
	// Texto auxiliar
	caption: {
		fontSize: typography.fontSize.xs,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
	},
	overline: {
		fontSize: typography.fontSize.xs,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.wide,
		textTransform: 'uppercase',
	},
	// Componentes
	button: {
		fontSize: typography.fontSize.sm,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	buttonLarge: {
		fontSize: typography.fontSize.base,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	buttonSmall: {
		fontSize: typography.fontSize.xs,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	label: {
		fontSize: typography.fontSize.sm,
		fontWeight: typography.fontWeight.medium,
		lineHeight: typography.lineHeight.tight,
		letterSpacing: typography.letterSpacing.normal,
	},
	input: {
		fontSize: typography.fontSize.base,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
	},
	// Código
	code: {
		fontSize: typography.fontSize.sm,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
		fontFamily: typography.fontFamily.mono,
	},
	codeInline: {
		fontSize: typography.fontSize.sm,
		fontWeight: typography.fontWeight.normal,
		lineHeight: typography.lineHeight.normal,
		letterSpacing: typography.letterSpacing.normal,
		fontFamily: typography.fontFamily.mono,
	},
} as const;

export type SemanticTypography = keyof typeof semanticTypography;
