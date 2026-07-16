<script lang="ts">
	interface Props {
		title?: string;
		subtitle?: string;
		breadcrumb?: string[];
	}

	let {
		title,
		subtitle,
		breadcrumb = [],
	}: Props = $props();
</script>

<div class="page-header">
	<div class="page-header-content">
		{#if breadcrumb.length > 0}
			<nav class="page-header-breadcrumb" aria-label="Breadcrumb">
				{#each breadcrumb as item, index}
					{#if index > 0}
						<span class="breadcrumb-separator">/</span>
					{/if}
					<span class="breadcrumb-item {index === breadcrumb.length - 1 ? 'breadcrumb-item-active' : ''}">
						{item}
					</span>
				{/each}
			</nav>
		{/if}
		{#if title}
			<h1 class="page-header-title">{title}</h1>
		{/if}
		{#if subtitle}
			<p class="page-header-subtitle">{subtitle}</p>
		{/if}
	</div>
	<div class="page-header-actions">
		<slot name="actions" />
	</div>
</div>

<style>
	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1.5rem;
		margin-bottom: 2rem;
		padding-bottom: 1.5rem;
		border-bottom: 1px solid #e2e8f0;
	}

	.page-header-content {
		flex: 1;
		min-width: 0;
	}

	.page-header-breadcrumb {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 0.75rem;
		font-size: 0.875rem;
		overflow: hidden;
		white-space: nowrap;
	}

	.breadcrumb-item {
		color: #64748b;
		transition: color 0.2s ease;
	}

	.breadcrumb-item:hover {
		color: #6366f1;
	}

	.breadcrumb-item-active {
		color: #0f172a;
		font-weight: 600;
	}

	.breadcrumb-separator {
		color: #cbd5e1;
		font-size: 0.75rem;
	}

	.page-header-title {
		font-size: 1.75rem;
		font-weight: 700;
		color: #0f172a;
		margin: 0 0 0.5rem 0;
		line-height: 1.2;
	}

	.page-header-subtitle {
		font-size: 0.875rem;
		color: #64748b;
		margin: 0;
		line-height: 1.5;
	}

	.page-header-actions {
		display: flex;
		gap: 0.75rem;
		align-items: flex-start;
		flex-shrink: 0;
	}

	@media (max-width: 768px) {
		.page-header {
			flex-direction: column;
			gap: 1rem;
		}

		.page-header-actions {
			width: 100%;
			justify-content: flex-start;
		}
	}
</style>
