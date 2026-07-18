<script lang="ts">
	type AlertVariant = 'info' | 'success' | 'error' | 'warning';

	interface Props {
		variant?: AlertVariant;
		dismissible?: boolean;
		onDismiss?: () => void;
		class?: string;
	}

	let {
		variant = 'info',
		dismissible = false,
		onDismiss,
		class: className = '',
	}: Props = $props();

	let visible = $state(true);

	function handleDismiss() {
		visible = false;
		onDismiss?.();
	}
</script>

{#if visible}
	<div
		class:alert={true}
		class:alert-info={variant === 'info'}
		class:alert-success={variant === 'success'}
		class:alert-error={variant === 'error'}
		class:alert-warning={variant === 'warning'}
		class={className}
	>
		<div class="alert-content">
			<slot />
		</div>
		{#if dismissible}
			<button class="alert-dismiss" onmousedown={handleDismiss} aria-label="Fechar">
				×
			</button>
		{/if}
	</div>
{/if}

<style>
	.alert {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 12px 16px;
		border-radius: 6px;
		border: 1px solid;
		background-color: white;
	}

	.alert-info {
		border-color: #3b82f6;
		background-color: #eff6ff;
	}

	.alert-success {
		border-color: #22c55e;
		background-color: #f0fdf4;
	}

	.alert-error {
		border-color: #ef4444;
		background-color: #fef2f2;
	}

	.alert-warning {
		border-color: #eab308;
		background-color: #fefce8;
	}

	.alert-content {
		flex: 1;
		font-size: 0.875rem;
		color: #171717;
	}

	.alert-dismiss {
		background: none;
		border: none;
		font-size: 1.5rem;
		line-height: 1;
		color: #525252;
		cursor: pointer;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
	}

	.alert-dismiss:hover {
		background-color: rgba(0, 0, 0, 0.05);
	}
</style>
