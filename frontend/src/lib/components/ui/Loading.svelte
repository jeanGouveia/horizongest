<script lang="ts">
	import { Loader2 } from '@lucide/svelte';

	interface Props {
		type?: 'spinner' | 'skeleton' | 'dots';
		size?: 'sm' | 'md' | 'lg';
		text?: string;
		message?: string;
		skeletonType?: 'card' | 'list' | 'table' | 'text' | 'avatar';
		lines?: number;
	}

	let {
		type = 'spinner',
		size = 'md',
		text,
		message,
		skeletonType = 'card',
		lines = 3,
	}: Props = $props();
</script>

<div class="loading-wrapper">
	{#if type === 'spinner'}
		<div class="spinner spinner-{size}">
			<Loader2 size={size === 'sm' ? 24 : size === 'lg' ? 48 : 32} class="spinner-icon" />
		</div>
	{:else if type === 'dots'}
		<div class="dots dots-{size}">
			<div class="dot"></div>
			<div class="dot"></div>
			<div class="dot"></div>
		</div>
	{:else if type === 'skeleton'}
		<div class="skeleton skeleton-{skeletonType}">
			{#if skeletonType === 'card'}
				<div class="skeleton-header"></div>
				<div class="skeleton-content">
					{#each Array(lines) as _}
						<div class="skeleton-line"></div>
					{/each}
				</div>
			{:else if skeletonType === 'list'}
				{#each Array(lines) as _}
					<div class="skeleton-list-item">
						<div class="skeleton-avatar"></div>
						<div class="skeleton-list-content">
							<div class="skeleton-line skeleton-line-short"></div>
							<div class="skeleton-line skeleton-line-medium"></div>
						</div>
					</div>
				{/each}
			{:else if skeletonType === 'table'}
				<div class="skeleton-table-header">
					{#each Array(4) as _}
						<div class="skeleton-cell"></div>
					{/each}
				</div>
				{#each Array(lines) as _}
					<div class="skeleton-table-row">
						{#each Array(4) as _}
							<div class="skeleton-cell"></div>
						{/each}
					</div>
				{/each}
			{:else if skeletonType === 'avatar'}
				<div class="skeleton-avatar-large"></div>
			{:else}
				<!-- text skeleton -->
				{#each Array(lines) as _}
					<div class="skeleton-line"></div>
				{/each}
			{/if}
		</div>
	{/if}
	{#if text || message}
		<span class="loading-text">{text || message}</span>
	{/if}
</div>

<style>
	.loading-wrapper {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		padding: 1.5rem;
	}

	/* Spinner */
	.spinner {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.spinner-icon {
		color: #6366f1;
		animation: spin 1s linear infinite;
	}

	.spinner-sm .spinner-icon {
		width: 24px;
		height: 24px;
	}

	.spinner-md .spinner-icon {
		width: 32px;
		height: 32px;
	}

	.spinner-lg .spinner-icon {
		width: 48px;
		height: 48px;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	/* Dots */
	.dots {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: #6366f1;
		animation: pulse 1.4s ease-in-out infinite both;
	}

	.dots-sm .dot {
		width: 6px;
		height: 6px;
	}

	.dots-lg .dot {
		width: 12px;
		height: 12px;
	}

	.dot:nth-child(1) {
		animation-delay: -0.32s;
	}

	.dot:nth-child(2) {
		animation-delay: -0.16s;
	}

	@keyframes pulse {
		0%, 80%, 100% {
			transform: scale(0);
			opacity: 0.5;
		}
		40% {
			transform: scale(1);
			opacity: 1;
		}
	}

	/* Skeleton */
	.skeleton {
		width: 100%;
	}

	.skeleton-line {
		height: 12px;
		background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
		background-size: 200% 100%;
		border-radius: 4px;
		animation: shimmer 1.5s infinite;
		margin-bottom: 0.5rem;
	}

	.skeleton-line:last-child {
		margin-bottom: 0;
	}

	.skeleton-line-short {
		width: 60%;
	}

	.skeleton-line-medium {
		width: 80%;
	}

	.skeleton-header {
		height: 24px;
		width: 40%;
		background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
		background-size: 200% 100%;
		border-radius: 4px;
		margin-bottom: 1rem;
		animation: shimmer 1.5s infinite;
	}

	.skeleton-content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.skeleton-avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
		background-size: 200% 100%;
		animation: shimmer 1.5s infinite;
		flex-shrink: 0;
	}

	.skeleton-avatar-large {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
		background-size: 200% 100%;
		animation: shimmer 1.5s infinite;
		margin: 0 auto;
	}

	.skeleton-list-item {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 0.75rem 0;
	}

	.skeleton-list-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.skeleton-table-header {
		display: flex;
		gap: 1rem;
		padding-bottom: 0.75rem;
		border-bottom: 1px solid #e2e8f0;
		margin-bottom: 0.75rem;
	}

	.skeleton-table-row {
		display: flex;
		gap: 1rem;
		padding: 0.75rem 0;
	}

	.skeleton-cell {
		height: 16px;
		flex: 1;
		background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
		background-size: 200% 100%;
		border-radius: 4px;
		animation: shimmer 1.5s infinite;
	}

	@keyframes shimmer {
		0% {
			background-position: 200% 0;
		}
		100% {
			background-position: -200% 0;
		}
	}

	.loading-text {
		font-size: 0.875rem;
		color: #64748b;
		font-weight: 500;
		letter-spacing: 0.025em;
	}
</style>
