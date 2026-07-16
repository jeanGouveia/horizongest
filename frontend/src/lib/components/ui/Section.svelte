<script lang="ts">
	interface Props {
		title?: string;
		class?: string;
	}

	let {
		title,
		class: className = '',
	}: Props = $props();
</script>

<div class="section {className}">
	{#if title}
		<h2 class="section-title">{title}</h2>
	{/if}
	<div class="section-content">
		<slot />
	</div>
</div>

<style>
	.section {
		margin-bottom: 32px;
	}

	.section-title {
		font-size: 1.25rem;
		font-weight: 600;
		color: #171717;
		margin: 0 0 16px 0;
	}

	.section-content {
		padding: 16px;
		background-color: white;
		border: 1px solid #e5e5e5;
		border-radius: 8px;
	}
</style>
