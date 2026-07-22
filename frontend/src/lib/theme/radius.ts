/**
 * Sistema de bordas arredondadas do horizongest - Redesign Visual
 * Tokens de border-radius para consistência visual
 * Inspirado em sistemas modernos como Linear e Vercel
 */

export const radius = {
	none: '0px',
	xs: '2px',
	sm: '4px',
	base: '6px',
	md: '8px',
	lg: '12px',
	xl: '16px',
	'2xl': '20px',
	'3xl': '24px',
	full: '9999px',
} as const;

export type RadiusToken = keyof typeof radius;

// Bordas arredondadas semânticas - Mais refinadas para look profissional
export const semanticRadius = {
	button: radius.md,
	input: radius.md,
	card: radius.lg,
	modal: radius.xl,
	badge: radius.full,
	avatar: radius.full,
	dropdown: radius.lg,
	tooltip: radius.md,
} as const;

export type SemanticRadius = keyof typeof semanticRadius;
