<script lang="ts">
	interface Props {
		open?: boolean;
		title?: string;
		onClose?: () => void;
	}

	let {
		open = false,
		title,
		onClose,
	}: Props = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onClose?.();
		}
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onClose?.();
		}
	}
</script>

{#if open}
	<div
		class="modal-backdrop"
		onclick={handleBackdropClick}
		onkeydown={handleKeyDown}
		role="dialog"
		aria-modal="true"
	>
		<div class="modal-content">
			{#if title}
				<div class="modal-header">
					<h2 class="modal-title">{title}</h2>
					<button class="modal-close" onclick={onClose} aria-label="Fechar">
						×
					</button>
				</div>
			{/if}
			<div class="modal-body">
				<slot />
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
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

	.modal-content {
		background-color: white;
		border-radius: 12px;
		box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
		max-width: 500px;
		width: 100%;
		max-height: 90vh;
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 16px 24px;
		border-bottom: 1px solid #e5e5e5;
	}

	.modal-title {
		font-size: 1.25rem;
		font-weight: 600;
		color: #171717;
		margin: 0;
	}

	.modal-close {
		background: none;
		border: none;
		font-size: 1.5rem;
		line-height: 1;
		color: #525252;
		cursor: pointer;
		padding: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
	}

	.modal-close:hover {
		background-color: rgba(0, 0, 0, 0.05);
	}

	.modal-body {
		padding: 24px;
		overflow-y: auto;
	}
</style>
