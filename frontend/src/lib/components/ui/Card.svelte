<script lang="ts">
	interface Props {
		title?: string;
		subtitle?: string;
		class?: string;
		elevated?: boolean;
		hoverable?: boolean;
		onclick?: (e: MouseEvent) => void;
	}

	let {
		title,
		subtitle,
		class: className = '',
		elevated = false,
		hoverable = false,
		onclick,
	}: Props = $props();
</script>

<div class="card {className} {elevated ? 'card-elevated' : ''} {hoverable ? 'card-hoverable' : ''}" on:click={onclick}>
	{#if title || subtitle}
		<div class="card-header">
			{#if title}
				<h3 class="card-title">{title}</h3>
			{/if}
			{#if subtitle}
				<p class="card-subtitle">{subtitle}</p>
			{/if}
		</div>
	{/if}
	<div class="card-body">
		<slot />
	</div>
</div>

<style>
	.card {
		background-color: #ffffff;
		border: 1px solid #e2e8f0;
		border-radius: 8px;
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.04);
		overflow: hidden;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.card:hover:not(.card-hoverable) {
		box-shadow: 0 2px 4px 0 rgb(0 0 0 / 0.06);
		border-color: #cbd5e1;
	}

	.card-elevated {
		box-shadow: 0 4px 8px 0 rgb(0 0 0 / 0.06);
		border-color: #e2e8f0;
	}

	.card-elevated:hover:not(.card-hoverable) {
		box-shadow: 0 6px 12px 0 rgb(0 0 0 / 0.08);
		border-color: #cbd5e1;
	}

	.card-hoverable {
		cursor: pointer;
	}

	.card-hoverable:hover {
		box-shadow: 0 6px 16px 0 rgb(0 0 0 / 0.1);
		transform: translateY(-1px);
		border-color: #cbd5e1;
	}

	.card-hoverable:active {
		transform: translateY(0);
		box-shadow: 0 2px 4px 0 rgb(0 0 0 / 0.06);
	}

	.card-header {
		padding: 0.875rem 1rem;
		border-bottom: 1px solid #e2e8f0;
		background-color: #ffffff;
	}

	.card-title {
		font-size: 0.9375rem;
		font-weight: 600;
		color: #0f172a;
		margin: 0;
		line-height: 1.3;
		letter-spacing: -0.025em;
	}

	.card-subtitle {
		font-size: 0.8125rem;
		color: #64748b;
		margin: 0.25rem 0 0 0;
		line-height: 1.4;
	}

	.card-body {
		padding: 1rem;
	}
</style>
