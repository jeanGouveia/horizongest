/**
 * Sistema de espaçamentos do horizongest - Redesign Visual
 * Tokens de espaçamento para consistência visual
 * Inspirado em sistemas modernos como Tailwind CSS
 */

export const spacing = {
	0: '0px',
	1: '4px',
	2: '8px',
	3: '12px',
	4: '16px',
	5: '20px',
	6: '24px',
	7: '28px',
	8: '32px',
	9: '36px',
	10: '40px',
	11: '44px',
	12: '48px',
	14: '56px',
	16: '64px',
	20: '80px',
	24: '96px',
	28: '112px',
	32: '128px',
	36: '144px',
	40: '160px',
	44: '176px',
	48: '192px',
	52: '208px',
	56: '224px',
	60: '240px',
	64: '256px',
	72: '288px',
	80: '320px',
	96: '384px',
} as const;

export type SpacingToken = keyof typeof spacing;

// Espaçamentos semânticos - Mais granulares para melhor controle
export const semanticSpacing = {
	none: spacing[0],
	xs: spacing[1],
	sm: spacing[2],
	base: spacing[3],
	md: spacing[4],
	lg: spacing[6],
	xl: spacing[8],
	'2xl': spacing[12],
	'3xl': spacing[16],
	'4xl': spacing[20],
	'5xl': spacing[24],
} as const;

export type SemanticSpacing = keyof typeof semanticSpacing;
