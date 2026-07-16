<script lang="ts">
	import type { HTMLSelectAttributes } from 'svelte/elements';

	interface Props extends HTMLSelectAttributes {
		error?: string;
		label?: string;
		value?: string;
	}

	let {
		error,
		label,
		value = $bindable(),
		class: className = '',
		...restProps
	}: Props = $props();
</script>

<div class="select-wrapper {className}">
	{#if label}
		<label for={restProps.id} class="select-label">{label}</label>
	{/if}
	<select
		bind:value
		{...restProps}
		class:select={true}
		class:select-error={error}
	>
		<slot />
	</select>
	{#if error}
		<span class="select-error-text">{error}</span>
	{/if}
</div>

<style>
	.select-wrapper {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.select-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #0f172a;
		line-height: 1.4;
		letter-spacing: -0.025em;
	}

	.select {
		padding: 0.625rem 1rem;
		font-size: 0.875rem;
		line-height: 1.5;
		border: 1px solid #f1f5f9;
		border-radius: 8px;
		background-color: white;
		color: #0f172a;
		font-family: inherit;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		cursor: pointer;
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
		letter-spacing: -0.025em;
	}

	.select:hover:not(:disabled) {
		border-color: #e2e8f0;
	}

	.select:disabled {
		background-color: #f8fafc;
		color: #94a3b8;
		cursor: not-allowed;
	}

	.select:focus {
		outline: none;
		border-color: #6366f1;
		box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.08);
	}

	.select:focus-visible {
		outline: 2px solid #6366f1;
		outline-offset: 2px;
	}

	.select-error {
		border-color: #ef4444;
	}

	.select-error:focus {
		border-color: #ef4444;
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.08);
	}

	.select-error-text {
		font-size: 0.75rem;
		color: #ef4444;
		line-height: 1.4;
	}
</style>
