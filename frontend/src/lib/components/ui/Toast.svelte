<script lang="ts">
	import { CheckCircle, AlertCircle, Info, AlertTriangle, X } from '@lucide/svelte';

	type ToastVariant = 'success' | 'error' | 'warning' | 'info';

	interface Toast {
		id: string;
		variant: ToastVariant;
		title: string;
		message?: string;
		duration?: number;
	}

	interface Props {
		toast: Toast;
		onDismiss: (id: string) => void;
	}

	let { toast, onDismiss }: Props = $props();

	let remaining = $state(toast.duration || 5000);
	let isPaused = $state(false);

	let interval: number;

	function getIcon() {
		switch (toast.variant) {
			case 'success': return CheckCircle;
			case 'error': return AlertCircle;
			case 'warning': return AlertTriangle;
			case 'info': return Info;
			default: return Info;
		}
	}

	function getVariantClass() {
		switch (toast.variant) {
			case 'success': return 'toast-success';
			case 'error': return 'toast-error';
			case 'warning': return 'toast-warning';
			case 'info': return 'toast-info';
			default: return 'toast-info';
		}
	}

	function startTimer() {
		if (toast.duration === 0) return;
		
		interval = window.setInterval(() => {
			if (!isPaused) {
				remaining -= 100;
				if (remaining <= 0) {
					onDismiss(toast.id);
				}
			}
		}, 100);
	}

	function pauseTimer() {
		isPaused = true;
	}

	function resumeTimer() {
		isPaused = false;
	}

	function dismiss() {
		onDismiss(toast.id);
	}

	$effect(() => {
		startTimer();
		return () => clearInterval(interval);
	});
</script>

<div
	class="toast {getVariantClass()}"
	onmouseenter={pauseTimer}
	onmouseleave={resumeTimer}
	role="alert"
	aria-live="polite"
>
	<div class="toast-icon">
		<svelte:component this={getIcon()} size={20} />
	</div>
	
	<div class="toast-content">
		<div class="toast-title">{toast.title}</div>
		{#if toast.message}
			<div class="toast-message">{toast.message}</div>
		{/if}
	</div>
	
	<button class="toast-dismiss" onclick={dismiss} aria-label="Fechar notificação">
		<X size={16} />
	</button>
	
	{#if toast.duration && toast.duration > 0}
		<div class="toast-progress" style="width: {(remaining / (toast.duration || 5000)) * 100}%"></div>
	{/if}
</div>

<style>
	.toast {
		display: flex;
		align-items: flex-start;
		gap: 0.75rem;
		padding: 1rem 1.25rem;
		border-radius: 12px;
		background: white;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		border: 1px solid #f1f5f9;
		min-width: 320px;
		max-width: 420px;
		position: relative;
		overflow: hidden;
		animation: slideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.toast:hover {
		box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
	}

	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateX(100%);
		}
		to {
			opacity: 1;
			transform: translateX(0);
		}
	}

	.toast-icon {
		flex-shrink: 0;
		padding-top: 0.125rem;
	}

	.toast-success .toast-icon {
		color: #10b981;
	}

	.toast-error .toast-icon {
		color: #ef4444;
	}

	.toast-warning .toast-icon {
		color: #f59e0b;
	}

	.toast-info .toast-icon {
		color: #6366f1;
	}

	.toast-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.toast-title {
		font-size: 0.875rem;
		font-weight: 600;
		color: #0f172a;
		line-height: 1.4;
		letter-spacing: -0.025em;
	}

	.toast-message {
		font-size: 0.8125rem;
		color: #64748b;
		line-height: 1.5;
	}

	.toast-dismiss {
		flex-shrink: 0;
		padding: 0.25rem;
		background: transparent;
		border: none;
		color: #94a3b8;
		cursor: pointer;
		border-radius: 4px;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.toast-dismiss:hover {
		background: #f1f5f9;
		color: #64748b;
	}

	.toast-progress {
		position: absolute;
		bottom: 0;
		left: 0;
		height: 2px;
		background: currentColor;
		transition: width 0.1s linear;
	}

	.toast-success .toast-progress {
		background: #10b981;
	}

	.toast-error .toast-progress {
		background: #ef4444;
	}

	.toast-warning .toast-progress {
		background: #f59e0b;
	}

	.toast-info .toast-progress {
		background: #6366f1;
	}
</style>
