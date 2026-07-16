<script lang="ts">
	import { CheckCircle, Clock, XCircle, AlertTriangle } from '@lucide/svelte';
	import { semanticTransitions } from '$lib/theme/transitions';

	type BadgeVariant = 'default' | 'primary' | 'success' | 'error' | 'danger' | 'warning' | 'info' | 'active' | 'inactive' | 'paid' | 'pending' | 'low-stock' | 'no-stock';
	type BadgeSize = 'sm' | 'md' | 'lg';

	interface Props {
		variant?: BadgeVariant;
		size?: BadgeSize;
		class?: string;
		dot?: boolean;
		icon?: boolean;
	}

	let {
		variant = 'default',
		size = 'md',
		class: className = '',
		dot = false,
		icon = false,
	}: Props = $props();

	function getIconComponent() {
		switch (variant) {
			case 'success':
			case 'active':
			case 'paid':
				return CheckCircle;
			case 'error':
			case 'danger':
			case 'inactive':
			case 'no-stock':
				return XCircle;
			case 'warning':
			case 'pending':
			case 'low-stock':
				return AlertTriangle;
			default:
				return null;
		}
	}

	let IconComponent = getIconComponent();
</script>

<span
	class:badge={true}
	class:badge-default={variant === 'default'}
	class:badge-primary={variant === 'primary'}
	class:badge-success={variant === 'success'}
	class:badge-error={variant === 'error'}
	class:badge-danger={variant === 'danger'}
	class:badge-warning={variant === 'warning'}
	class:badge-info={variant === 'info'}
	class:badge-active={variant === 'active'}
	class:badge-inactive={variant === 'inactive'}
	class:badge-paid={variant === 'paid'}
	class:badge-pending={variant === 'pending'}
	class:badge-low-stock={variant === 'low-stock'}
	class:badge-no-stock={variant === 'no-stock'}
	class:badge-sm={size === 'sm'}
	class:badge-md={size === 'md'}
	class:badge-lg={size === 'lg'}
	class:badge-dot={dot}
	class:badge-icon={icon}
	class={className}
>
	{#if dot}
		<span class="badge-dot-indicator"></span>
	{/if}
	{#if icon && IconComponent}
		<svelte:component this={IconComponent} size={size === 'sm' ? 12 : size === 'lg' ? 16 : 14} class="badge-icon" />
	{/if}
	<slot />
</span>

<style>
	.badge {
		display: inline-flex;
		align-items: center;
		gap: var(--spacing-3);
		padding: var(--spacing-1) var(--spacing-2-5);
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
		line-height: var(--line-height-tight);
		border-radius: var(--radius-full);
		white-space: nowrap;
		transition: all var(--transition-duration-base) var(--transition-easing-base);
		letter-spacing: var(--letter-spacing-wide);
	}

	.badge-sm {
		padding: var(--spacing-0-5) var(--spacing-2);
		font-size: 0.6875rem;
		gap: var(--spacing-1);
	}

	.badge-md {
		padding: var(--spacing-1) var(--spacing-2-5);
		font-size: var(--font-size-xs);
		gap: var(--spacing-1-5);
	}

	.badge-lg {
		padding: var(--spacing-1-5) var(--spacing-3);
		font-size: 0.8125rem;
		gap: var(--spacing-2);
	}

	.badge-dot {
		padding: 0.25rem;
	}

	.badge-dot-indicator {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: currentColor;
	}

	.badge-default {
		background-color: #f1f5f9;
		color: #475569;
		border: 1px solid #e2e8f0;
	}

	.badge-default:hover {
		background-color: #e2e8f0;
	}

	.badge-primary {
		background-color: #eef2ff;
		color: #4f46e5;
		border: 1px solid #e0e7ff;
	}

	.badge-primary:hover {
		background-color: #e0e7ff;
	}

	.badge-success {
		background-color: #ecfdf5;
		color: #059669;
		border: 1px solid #d1fae5;
	}

	.badge-success:hover {
		background-color: #d1fae5;
	}

	.badge-error {
		background-color: #fef2f2;
		color: #dc2626;
		border: 1px solid #fee2e2;
	}

	.badge-error:hover {
		background-color: #fee2e2;
	}

	.badge-danger {
		background-color: #fef2f2;
		color: #dc2626;
		border: 1px solid #fee2e2;
	}

	.badge-danger:hover {
		background-color: #fee2e2;
	}

	.badge-warning {
		background-color: #fffbeb;
		color: #d97706;
		border: 1px solid #fef3c7;
	}

	.badge-warning:hover {
		background-color: #fef3c7;
	}

	.badge-info {
		background-color: #f0f9ff;
		color: #0284c7;
		border: 1px solid #e0f2fe;
	}

	.badge-info:hover {
		background-color: #e0f2fe;
	}

	.badge-active {
		background-color: #ecfdf5;
		color: #059669;
		border: 1px solid #a7f3d0;
	}

	.badge-active:hover {
		background-color: #d1fae5;
	}

	.badge-inactive {
		background-color: #f1f5f9;
		color: #64748b;
		border: 1px solid #e2e8f0;
	}

	.badge-inactive:hover {
		background-color: #e2e8f0;
	}

	.badge-paid {
		background-color: #ecfdf5;
		color: #059669;
		border: 1px solid #a7f3d0;
	}

	.badge-paid:hover {
		background-color: #d1fae5;
	}

	.badge-pending {
		background-color: #fffbeb;
		color: #d97706;
		border: 1px solid #fde68a;
	}

	.badge-pending:hover {
		background-color: #fef3c7;
	}

	.badge-low-stock {
		background-color: #fffbeb;
		color: #d97706;
		border: 1px solid #fde68a;
	}

	.badge-low-stock:hover {
		background-color: #fef3c7;
	}

	.badge-no-stock {
		background-color: #fef2f2;
		color: #dc2626;
		border: 1px solid #fecaca;
	}

	.badge-no-stock:hover {
		background-color: #fee2e2;
	}

	.badge-icon {
		flex-shrink: 0;
	}
</style>
