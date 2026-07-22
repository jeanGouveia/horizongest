<script lang="ts">
	type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'link' | 'success';
	type ButtonSize = 'sm' | 'md' | 'lg' | 'xl';

	interface Props {
		variant?: ButtonVariant;
		size?: ButtonSize;
		fullWidth?: boolean;
		loading?: boolean;
		href?: string;
		disabled?: boolean;
		class?: string;
		icon?: string;
		iconPosition?: 'left' | 'right';
		onclick?: (e: MouseEvent) => void;
		[key: string]: any;
	}

	let {
		variant = 'primary',
		size = 'md',
		fullWidth = false,
		loading = false,
		disabled = false,
		href,
		icon,
		iconPosition = 'left',
		class: className = '',
		onclick,
		...restProps
	}: Props = $props();

	const disabledAttr = $derived(disabled || loading);
</script>

{#if href}
	<a
		{...restProps}
		{href}
		onclick={onclick}
		class:btn={true}
		class:btn-primary={variant === 'primary'}
		class:btn-secondary={variant === 'secondary'}
		class:btn-danger={variant === 'danger'}
		class:btn-ghost={variant === 'ghost'}
		class:btn-link={variant === 'link'}
		class:btn-success={variant === 'success'}
		class:btn-sm={size === 'sm'}
		class:btn-md={size === 'md'}
		class:btn-lg={size === 'lg'}
		class:btn-xl={size === 'xl'}
		class:btn-full={fullWidth}
		class:btn-loading={loading}
		class:btn-icon={icon && !loading}
		class:disabled={disabledAttr}
		class={className}
	>
		{#if loading}
			<span class="btn-spinner"></span>
		{:else if icon && iconPosition === 'left'}
			<span class="btn-icon">{icon}</span>
		{/if}
		<slot />
		{#if icon && iconPosition === 'right' && !loading}
			<span class="btn-icon btn-icon-right">{icon}</span>
		{/if}
	</a>
{:else}
	<button
		{...restProps}
		onclick={onclick}
		class:btn={true}
		class:btn-primary={variant === 'primary'}
		class:btn-secondary={variant === 'secondary'}
		class:btn-danger={variant === 'danger'}
		class:btn-ghost={variant === 'ghost'}
		class:btn-link={variant === 'link'}
		class:btn-success={variant === 'success'}
		class:btn-sm={size === 'sm'}
		class:btn-md={size === 'md'}
		class:btn-lg={size === 'lg'}
		class:btn-xl={size === 'xl'}
		class:btn-full={fullWidth}
		class:btn-loading={loading}
		class:btn-icon={icon && !loading}
		class:disabled={disabledAttr}
		class={className}
		disabled={disabledAttr}
	>
		{#if loading}
			<span class="btn-spinner"></span>
		{:else if icon && iconPosition === 'left'}
			<span class="btn-icon">{icon}</span>
		{/if}
		<slot />
		{#if icon && iconPosition === 'right' && !loading}
			<span class="btn-icon btn-icon-right">{icon}</span>
		{/if}
	</button>
{/if}

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.625rem 1rem;
		font-size: 0.875rem;
		font-weight: 500;
		line-height: 1.25;
		border-radius: 8px;
		border: 1px solid transparent;
		cursor: pointer;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		font-family: inherit;
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
		letter-spacing: -0.025em;
		position: relative;
	}

	.btn:focus-visible {
		outline: 2px solid var(--color-primary-500);
		outline-offset: 2px;
	}

	.btn:active:not(:disabled) {
		transform: translateY(0);
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-primary {
		background-color: var(--color-primary-500);
		color: white;
		border-color: var(--color-primary-500);
	}

	.btn-primary:hover:not(:disabled) {
		background-color: var(--color-primary-600);
		border-color: var(--color-primary-600);
		box-shadow: 0 2px 8px rgba(99, 102, 241, 0.2);
		transform: translateY(-1px);
	}

	.btn-primary:active:not(:disabled) {
		background-color: var(--color-primary-700);
		border-color: var(--color-primary-700);
	}

	.btn-secondary {
		background-color: #ffffff;
		color: #0f172a;
		border-color: #f1f5f9;
	}

	.btn-secondary:hover:not(:disabled) {
		background-color: #f8fafc;
		border-color: #e2e8f0;
	}

	.btn-secondary:active:not(:disabled) {
		background-color: #f1f5f9;
		border-color: #cbd5e1;
	}

	.btn-danger {
		background-color: #ef4444;
		color: white;
		border-color: #ef4444;
	}

	.btn-danger:hover:not(:disabled) {
		background-color: #dc2626;
		border-color: #dc2626;
		box-shadow: 0 2px 8px rgba(239, 68, 68, 0.2);
		transform: translateY(-1px);
	}

	.btn-danger:active:not(:disabled) {
		background-color: #b91c1c;
		border-color: #b91c1c;
	}

	.btn-ghost {
		background-color: transparent;
		color: #0f172a;
		border-color: #f1f5f9;
	}

	.btn-ghost:hover:not(:disabled) {
		background-color: #f8fafc;
		border-color: #e2e8f0;
	}

	.btn-ghost:active:not(:disabled) {
		background-color: #f1f5f9;
		border-color: #cbd5e1;
	}

	.btn-link {
		background-color: transparent;
		color: var(--color-primary-500);
		border-color: transparent;
		padding: 0;
		box-shadow: none;
	}

	.btn-link:hover:not(:disabled) {
		text-decoration: underline;
	}

	.btn-link:active:not(:disabled) {
		text-decoration: underline;
		color: var(--color-primary-600);
	}

	.btn-success {
		background-color: #10b981;
		color: white;
		border-color: #10b981;
	}

	.btn-success:hover:not(:disabled) {
		background-color: #059669;
		border-color: #059669;
		box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
		transform: translateY(-1px);
	}

	.btn-success:active:not(:disabled) {
		background-color: #047857;
		border-color: #047857;
	}

	.btn-sm {
		padding: 0.375rem 0.75rem;
		font-size: 0.75rem;
		gap: 0.375rem;
	}

	.btn-md {
		padding: 0.625rem 1rem;
		font-size: 0.875rem;
		gap: 0.5rem;
	}

	.btn-lg {
		padding: 0.75rem 1.5rem;
		font-size: 1rem;
		gap: 0.625rem;
	}

	.btn-xl {
		padding: 1rem 2rem;
		font-size: 1.125rem;
		gap: 0.75rem;
	}

	.btn-full {
		width: 100%;
	}

	.btn-loading {
		position: relative;
		color: transparent !important;
		pointer-events: none;
	}

	.btn-spinner {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		display: inline-block;
		width: 16px;
		height: 16px;
		border: 2px solid currentColor;
		border-radius: 50%;
		border-top-color: transparent;
		animation: spin 0.6s linear infinite;
	}

	.btn-icon {
		font-size: 1rem;
		line-height: 1;
	}

	.btn-icon-right {
		margin-left: 0.25rem;
	}

	@keyframes spin {
		to {
			transform: translate(-50%, -50%) rotate(360deg);
		}
	}
</style>
