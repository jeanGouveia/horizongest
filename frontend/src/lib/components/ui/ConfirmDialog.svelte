<script lang="ts">
	interface Props {
		open?: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		cancelText?: string;
		onConfirm?: () => void;
		onCancel?: () => void;
		variant?: 'danger' | 'warning' | 'info';
	}

	let {
		open = false,
		title = 'Confirmar',
		message,
		confirmText = 'Confirmar',
		cancelText = 'Cancelar',
		onConfirm,
		onCancel,
		variant = 'info',
	}: Props = $props();
</script>

{#if open}
	<div
		class="confirm-backdrop"
		role="dialog"
		aria-modal="true"
		onkeydown={(e) => e.key === 'Escape' && onCancel?.()}
	>
		<div class="confirm-content">
			<div class="confirm-header">
				<h3 class="confirm-title">{title}</h3>
			</div>
			<div class="confirm-body">
				<p class="confirm-message">{message}</p>
			</div>
			<div class="confirm-footer">
				<button
					class="btn btn-secondary"
					onclick={onCancel}
				>
					{cancelText}
				</button>
				<button
					class:btn-primary={variant === 'info'}
					class:btn-danger={variant === 'danger'}
					class:btn-warning={variant === 'warning'}
					class="btn"
					onclick={onConfirm}
				>
					{confirmText}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.confirm-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 24px;
	}

	.confirm-content {
		background-color: white;
		border-radius: 12px;
		box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
		max-width: 400px;
		width: 100%;
	}

	.confirm-header {
		padding: 16px 24px;
		border-bottom: 1px solid #e5e5e5;
	}

	.confirm-title {
		font-size: 1.125rem;
		font-weight: 600;
		color: #171717;
		margin: 0;
	}

	.confirm-body {
		padding: 24px;
	}

	.confirm-message {
		font-size: 0.875rem;
		color: #525252;
		margin: 0;
	}

	.confirm-footer {
		display: flex;
		gap: 12px;
		justify-content: flex-end;
		padding: 16px 24px;
		border-top: 1px solid #e5e5e5;
	}

	.btn {
		padding: 8px 16px;
		font-size: 0.875rem;
		font-weight: 500;
		border-radius: 6px;
		border: 1px solid transparent;
		cursor: pointer;
		transition: all 0.2s ease;
		font-family: inherit;
	}

	.btn-secondary {
		background-color: white;
		color: #171717;
		border-color: #e5e5e5;
	}

	.btn-secondary:hover {
		background-color: #f5f5f5;
		border-color: #d4d4d4;
	}

	.btn-primary {
		background-color: #3b82f6;
		color: white;
		border-color: #3b82f6;
	}

	.btn-primary:hover {
		background-color: #2563eb;
		border-color: #2563eb;
	}

	.btn-danger {
		background-color: #ef4444;
		color: white;
		border-color: #ef4444;
	}

	.btn-danger:hover {
		background-color: #dc2626;
		border-color: #dc2626;
	}

	.btn-warning {
		background-color: #eab308;
		color: white;
		border-color: #eab308;
	}

	.btn-warning:hover {
		background-color: #ca8a04;
		border-color: #ca8a04;
	}
</style>
