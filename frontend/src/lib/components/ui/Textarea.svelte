<script lang="ts">
	import type { HTMLTextareaAttributes } from 'svelte/elements';

	interface Props extends HTMLTextareaAttributes {
		error?: string;
		label?: string;
		value?: string;
	}

	let {
		error,
		label,
		class: className = '',
		value = $bindable(''),
		...restProps
	}: Props = $props();
</script>

<div class="textarea-wrapper {className}">
	{#if label}
		<label for={restProps.id} class="textarea-label">{label}</label>
	{/if}
	<textarea
		{...restProps}
		bind:value
		class:textarea={true}
		class:textarea-error={error}
	></textarea>
	{#if error}
		<span class="textarea-error-text">{error}</span>
	{/if}
</div>

<style>
	.textarea-wrapper {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.textarea-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #525252;
	}

	.textarea {
		padding: 8px 12px;
		font-size: 0.875rem;
		line-height: 1.5;
		border: 1px solid #e5e5e5;
		border-radius: 6px;
		background-color: white;
		color: #171717;
		font-family: inherit;
		transition: border-color 0.2s ease;
		resize: vertical;
		min-height: 80px;
	}

	.textarea:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.textarea::placeholder {
		color: #a3a3a3;
	}

	.textarea-error {
		border-color: #ef4444;
	}

	.textarea-error:focus {
		border-color: #ef4444;
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
	}

	.textarea-error-text {
		font-size: 0.75rem;
		color: #ef4444;
	}
</style>
