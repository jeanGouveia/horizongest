<script lang="ts">
	import { Inbox, AlertCircle, CheckCircle, Info, Search, FileText, ShoppingCart, Package } from '@lucide/svelte';

	interface Props {
		title?: string;
		description?: string;
		icon?: string;
		action?: string;
		actionText?: string;
		onAction?: () => void;
		variant?: 'default' | 'error' | 'success' | 'info';
		size?: 'sm' | 'md' | 'lg';
		iconType?: 'inbox' | 'search' | 'file' | 'cart' | 'package';
	}

	let {
		title = 'Nenhum dado encontrado',
		description,
		icon,
		action,
		actionText,
		onAction,
		variant = 'default',
		size = 'md',
		iconType = 'inbox',
	}: Props = $props();

	function getIconComponent() {
		if (icon) return null; // Use custom emoji if provided
		switch (iconType) {
			case 'search':
				return Search;
			case 'file':
				return FileText;
			case 'cart':
				return ShoppingCart;
			case 'package':
				return Package;
			default:
				return Inbox;
		}
	}

	function getVariantIcon() {
		switch (variant) {
			case 'error':
				return AlertCircle;
			case 'success':
				return CheckCircle;
			case 'info':
				return Info;
			default:
				return null;
		}
	}

	let IconComponent = $derived(icon ? null : getIconComponent());
	let VariantIcon = $derived(getVariantIcon());
</script>

<div class="empty-state empty-state-{variant} empty-state-{size}">
	<div class="empty-state-icon-wrapper">
		{#if icon}
			<div class="empty-state-icon">{icon}</div>
		{:else if VariantIcon}
			<VariantIcon size={size === 'sm' ? 32 : size === 'lg' ? 64 : 48} class="empty-state-icon" />
		{:else if IconComponent}
			<IconComponent size={size === 'sm' ? 32 : size === 'lg' ? 64 : 48} class="empty-state-icon" />
		{/if}
	</div>
	<h3 class="empty-state-title">{title}</h3>
	{#if description}
		<p class="empty-state-description">{description}</p>
	{/if}
	{#if (action || actionText) && onAction}
		<button class="empty-state-action" onclick={onAction}>{action || actionText}</button>
	{/if}
</div>

<style>
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1.5rem;
		padding: 4rem 2rem;
		text-align: center;
		border-radius: 16px;
		background: #ffffff;
		border: 1px solid #e2e8f0;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.empty-state-sm {
		padding: 2.5rem 1.5rem;
		gap: 1rem;
	}

	.empty-state-lg {
		padding: 6rem 3rem;
		gap: 2rem;
	}

	.empty-state-icon-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 80px;
		height: 80px;
		border-radius: 50%;
		background: #f8fafc;
		color: #64748b;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.empty-state-sm .empty-state-icon-wrapper {
		width: 56px;
		height: 56px;
	}

	.empty-state-lg .empty-state-icon-wrapper {
		width: 96px;
		height: 96px;
	}

	.empty-state-icon {
		flex-shrink: 0;
		color: inherit;
	}

	.empty-state-icon-wrapper {
		font-size: 2rem;
		line-height: 1;
	}

	.empty-state-title {
		font-size: 1.125rem;
		font-weight: 600;
		color: #0f172a;
		margin: 0;
		line-height: 1.4;
		letter-spacing: -0.025em;
	}

	.empty-state-sm .empty-state-title {
		font-size: 1rem;
	}

	.empty-state-lg .empty-state-title {
		font-size: 1.25rem;
	}

	.empty-state-description {
		font-size: 0.9375rem;
		color: #64748b;
		max-width: 480px;
		margin: 0;
		line-height: 1.6;
	}

	.empty-state-sm .empty-state-description {
		font-size: 0.875rem;
		max-width: 360px;
	}

	.empty-state-lg .empty-state-description {
		font-size: 1rem;
		max-width: 600px;
	}

	.empty-state-action {
		margin-top: 0.5rem;
		padding: 0.75rem 1.5rem;
		background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
		letter-spacing: 0.025em;
	}

	.empty-state-action:hover {
		background: linear-gradient(135deg, #4f46e5 0%, #4338ca 100%);
		box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
		transform: translateY(-2px);
	}

	.empty-state-error {
		background: #fef2f2;
		border-color: #fecaca;
	}

	.empty-state-error .empty-state-icon-wrapper {
		background: #fee2e2;
		color: #dc2626;
	}

	.empty-state-success {
		background: #ecfdf5;
		border-color: #a7f3d0;
	}

	.empty-state-success .empty-state-icon-wrapper {
		background: #d1fae5;
		color: #059669;
	}

	.empty-state-info {
		background: #eff6ff;
		border-color: #bfdbfe;
	}

	.empty-state-info .empty-state-icon-wrapper {
		background: #dbeafe;
		color: #2563eb;
	}

	.empty-state:hover {
		border-color: #cbd5e1;
	}

	.empty-state-error:hover {
		border-color: #fca5a5;
	}

	.empty-state-success:hover {
		border-color: #6ee7b7;
	}

	.empty-state-info:hover {
		border-color: #93c5fd;
	}
</style>
