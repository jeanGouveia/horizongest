<script lang="ts">
	import { ChevronRight } from '@lucide/svelte';

	interface BreadcrumbItem {
		label: string;
		href?: string;
	}

	interface Props {
		breadcrumb?: BreadcrumbItem[];
		title?: string;
		description?: string;
	}

	let {
		breadcrumb = [],
		title = '',
		description = '',
	}: Props = $props();
</script>

<div class="workspace">
	{#if breadcrumb.length > 0}
		<nav class="workspace-breadcrumb" aria-label="Breadcrumb">
			{#each breadcrumb as item, index}
				{#if index > 0}
					<ChevronRight size={14} class="breadcrumb-separator" />
				{/if}
				{#if item.href}
					<a href={item.href} class="breadcrumb-item">
						{item.label}
					</a>
				{:else}
					<span class="breadcrumb-item breadcrumb-item-current">
						{item.label}
					</span>
				{/if}
			{/each}
		</nav>
	{/if}

	<div class="workspace-header">
		<div class="workspace-header-content">
			{#if title}
				<h1 class="workspace-title">{title}</h1>
			{/if}
			{#if description}
				<p class="workspace-description">{description}</p>
			{/if}
		</div>
		<slot name="actions" />
	</div>

	<div class="workspace-content">
		<slot />
	</div>
</div>

<style>
	.workspace {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.workspace-breadcrumb {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		font-size: 0.8125rem;
	}

	.breadcrumb-item {
		color: #64748b;
		text-decoration: none;
		transition: color 0.15s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.breadcrumb-item:hover {
		color: #6366f1;
	}

	.breadcrumb-item-current {
		color: #0f172a;
		font-weight: 500;
	}

	.breadcrumb-separator {
		color: #cbd5e1;
		flex-shrink: 0;
	}

	.workspace-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 0.75rem;
	}

	.workspace-header-content {
		flex: 1;
		min-width: 0;
	}

	.workspace-title {
		font-size: 1.5rem;
		font-weight: 600;
		color: #0f172a;
		letter-spacing: -0.025em;
		line-height: 1.2;
		margin: 0 0 0.25rem 0;
	}

	.workspace-description {
		font-size: 0.8125rem;
		color: #64748b;
		line-height: 1.4;
		margin: 0;
	}

	.workspace-actions {
		flex: 0 0 auto;
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.workspace-content {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	@media (max-width: 768px) {
		.workspace {
			gap: 0.5rem;
		}

		.workspace-header {
			flex-direction: column;
			gap: 0.5rem;
		}

		.workspace-title {
			font-size: 1.25rem;
		}

		.workspace-actions {
			width: 100%;
			justify-content: flex-start;
		}
	}
</style>
