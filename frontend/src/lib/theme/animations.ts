/**
 * Sistema de animações do PratoOnline - Sprint 9 PX
 * Keyframes e animações para microinterações
 * Animações sutis e profissionais
 */

export const animations = {
	// Spin (loading)
	spin: 'spin 1s linear infinite',

	// Pulse (attention)
	pulse: 'pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',

	// Bounce (feedback)
	bounce: 'bounce 1s infinite',

	// Ping (notification)
	ping: 'ping 1s cubic-bezier(0, 0, 0.2, 1) infinite',

	// Fade
	fadeIn: 'fadeIn 200ms ease-out',
	fadeOut: 'fadeOut 200ms ease-in',

	// Slide
	slideInUp: 'slideInUp 200ms ease-out',
	slideInDown: 'slideInDown 200ms ease-out',
	slideInLeft: 'slideInLeft 200ms ease-out',
	slideInRight: 'slideInRight 200ms ease-out',

	// Scale
	scaleIn: 'scaleIn 200ms ease-out',
	scaleOut: 'scaleOut 200ms ease-in',

	// Shake (error)
	shake: 'shake 500ms ease-in-out',
} as const;

export type AnimationName = keyof typeof animations;

// Keyframes para uso em CSS
export const keyframes = {
	spin: {
		'0%': { transform: 'rotate(0deg)' },
		'100%': { transform: 'rotate(360deg)' },
	},

	pulse: {
		'0%, 100%': { opacity: '1' },
		'50%': { opacity: '0.5' },
	},

	bounce: {
		'0%, 100%': {
			transform: 'translateY(-5%)',
			animationTimingFunction: 'cubic-bezier(0.8, 0, 1, 1)',
		},
		'50%': {
			transform: 'translateY(0)',
			animationTimingFunction: 'cubic-bezier(0, 0, 0.2, 1)',
		},
	},

	ping: {
		'75%, 100%': {
			transform: 'scale(2)',
			opacity: '0',
		},
	},

	fadeIn: {
		'0%': { opacity: '0' },
		'100%': { opacity: '1' },
	},

	fadeOut: {
		'0%': { opacity: '1' },
		'100%': { opacity: '0' },
	},

	slideInUp: {
		'0%': {
			transform: 'translateY(10px)',
			opacity: '0',
		},
		'100%': {
			transform: 'translateY(0)',
			opacity: '1',
		},
	},

	slideInDown: {
		'0%': {
			transform: 'translateY(-10px)',
			opacity: '0',
		},
		'100%': {
			transform: 'translateY(0)',
			opacity: '1',
		},
	},

	slideInLeft: {
		'0%': {
			transform: 'translateX(-10px)',
			opacity: '0',
		},
		'100%': {
			transform: 'translateX(0)',
			opacity: '1',
		},
	},

	slideInRight: {
		'0%': {
			transform: 'translateX(10px)',
			opacity: '0',
		},
		'100%': {
			transform: 'translateX(0)',
			opacity: '1',
		},
	},

	scaleIn: {
		'0%': {
			transform: 'scale(0.95)',
			opacity: '0',
		},
		'100%': {
			transform: 'scale(1)',
			opacity: '1',
		},
	},

	scaleOut: {
		'0%': {
			transform: 'scale(1)',
			opacity: '1',
		},
		'100%': {
			transform: 'scale(0.95)',
			opacity: '0',
		},
	},

	shake: {
		'0%, 100%': { transform: 'translateX(0)' },
		'10%, 30%, 50%, 70%, 90%': { transform: 'translateX(-4px)' },
		'20%, 40%, 60%, 80%': { transform: 'translateX(4px)' },
	},
} as const;

// Animações semânticas por contexto
export const semanticAnimations = {
	// Loading states
	loading: animations.spin,
	loadingSmall: animations.pulse,

	// Attention
	notification: animations.ping,
	attention: animations.pulse,

	// Entry
	entry: animations.fadeIn,
	entrySlide: animations.slideInUp,

	// Exit
	exit: animations.fadeOut,

	// Error feedback
	error: animations.shake,

	// Success feedback
	success: animations.scaleIn,
} as const;

export type SemanticAnimation = keyof typeof semanticAnimations;
