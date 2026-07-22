/**
 * Sistema de sombras do horizongest - Redesign Visual
 * Tokens de shadow para consistência visual
 * Inspirado em sistemas modernos como Linear e Vercel
 * Sombras mais sutis e elegantes para look profissional
 */

export const shadow = {
	none: 'none',
	xs: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
	sm: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
	base: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
	md: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
	lg: '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)',
	xl: '0 25px 50px -12px rgb(0 0 0 / 0.25)',
	'2xl': '0 50px 100px -20px rgb(0 0 0 / 0.25)',
	inner: 'inset 0 2px 4px 0 rgb(0 0 0 / 0.05)',
} as const;

export type ShadowToken = keyof typeof shadow;

// Sombras semânticas - Mais refinadas para look profissional
export const semanticShadow = {
	button: shadow.sm,
	buttonHover: shadow.base,
	input: shadow.sm,
	inputFocus: shadow.xs,
	card: shadow.base,
	cardHover: shadow.md,
	modal: shadow.xl,
	dropdown: shadow.lg,
	tooltip: shadow.md,
	floating: shadow.lg,
} as const;

export type SemanticShadow = keyof typeof semanticShadow;
