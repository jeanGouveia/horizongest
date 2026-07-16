/**
 * Sistema de transições do PratoOnline - Sprint 9 PX
 * Tokens de transição para consistência visual
 * Transições suaves e profissionais
 */

export const transitions = {
	// Durações
	duration: {
		fast: '150ms',
		base: '200ms',
		slow: '300ms',
		slower: '500ms',
	},

	// Easing functions
	easing: {
		linear: 'linear',
		easeIn: 'ease-in',
		easeOut: 'ease-out',
		easeInOut: 'ease-in-out',
		// Cubic beziers para transições mais naturais
		easeOutCubic: 'cubic-bezier(0.33, 1, 0.68, 1)',
		easeInOutCubic: 'cubic-bezier(0.65, 0, 0.35, 1)',
		easeOutQuart: 'cubic-bezier(0.25, 1, 0.5, 1)',
		easeInOutQuart: 'cubic-bezier(0.76, 0, 0.24, 1)',
		easeOutExpo: 'cubic-bezier(0.16, 1, 0.3, 1)',
		easeInOutExpo: 'cubic-bezier(0.87, 0, 0.13, 1)',
	},
} as const;

export type TransitionDuration = keyof typeof transitions.duration;
export type TransitionEasing = keyof typeof transitions.easing;

// Transições semânticas - Para uso consistente
export const semanticTransitions = {
	// Microinterações rápidas
	fast: {
		duration: transitions.duration.fast,
		easing: transitions.easing.easeOutCubic,
	},

	// Transições padrão
	base: {
		duration: transitions.duration.base,
		easing: transitions.easing.easeOutCubic,
	},

	// Transições mais complexas
	slow: {
		duration: transitions.duration.slow,
		easing: transitions.easing.easeInOutCubic,
	},

	// Transições muito suaves
	slower: {
		duration: transitions.duration.slower,
		easing: transitions.easing.easeInOutCubic,
	},
} as const;

// Propriedades que podem ser animadas
export const animatableProperties = {
	opacity: 'opacity',
	transform: 'transform',
	color: 'color, background-color, border-color',
	spacing: 'padding, margin, gap',
	size: 'width, height',
	all: 'all',
} as const;

// Transições completas para uso direto
export const transitionPresets = {
	// Fade
	fadeIn: 'opacity 150ms ease-out',
	fadeOut: 'opacity 150ms ease-in',
	fade: 'opacity 200ms ease-in-out',

	// Slide
	slideUp: 'transform 200ms ease-out',
	slideDown: 'transform 200ms ease-out',
	slideLeft: 'transform 200ms ease-out',
	slideRight: 'transform 200ms ease-out',

	// Scale
	scaleIn: 'transform 200ms ease-out',
	scaleOut: 'transform 200ms ease-in',
	scale: 'transform 200ms ease-in-out',

	// Color
	color: 'color 200ms ease-out, background-color 200ms ease-out, border-color 200ms ease-out',

	// Spacing
	spacing: 'padding 200ms ease-out, margin 200ms ease-out, gap 200ms ease-out',

	// Default para componentes UI
	default: 'all 200ms ease-out',

	// Hover
	hover: 'all 150ms ease-out',

	// Focus
	focus: 'box-shadow 150ms ease-out, border-color 150ms ease-out',

	// Modal
	modal: 'opacity 300ms ease-in-out, transform 300ms ease-in-out',

	// Drawer
	drawer: 'transform 300ms ease-out',

	// Dropdown
	dropdown: 'opacity 200ms ease-out, transform 200ms ease-out',
} as const;

export type TransitionPreset = keyof typeof transitionPresets;
