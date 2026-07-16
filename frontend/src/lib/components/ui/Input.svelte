<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends Omit<HTMLInputAttributes, 'size'> {
		error?: string;
		label?: string;
		value?: string | number;
		helper?: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let {
		error,
		label,
		helper,
		size = 'md',
		class: className = '',
		value = $bindable(''),
		...restProps
	}: Props = $props();
</script>

<div class="input-wrapper {className} input-wrapper-{size}">
	{#if label}
		<label for={restProps.id} class="input-label">{label}</label>
	{/if}
	<input
		{...restProps}
		bind:value
		class:input={true}
		class:input-error={error}
		class:input-sm={size === 'sm'}
		class:input-md={size === 'md'}
		class:input-lg={size === 'lg'}
	/>
	{#if helper && !error}
		<span class="input-helper">{helper}</span>
	{/if}
	{#if error}
		<span class="input-error-text">{error}</span>
	{/if}
</div>

<style>
	.input-wrapper {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.input-wrapper-sm {
		gap: 0.375rem;
	}

	.input-wrapper-lg {
		gap: 0.625rem;
	}

	.input-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #0f172a;
		line-height: 1.4;
		letter-spacing: -0.025em;
	}

	.input {
		padding: 0.625rem 1rem;
		font-size: 0.875rem;
		line-height: 1.5;
		border: 1px solid #f1f5f9;
		border-radius: 8px;
		background-color: #ffffff;
		color: #0f172a;
		font-family: inherit;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
		letter-spacing: -0.025em;
	}

	.input:hover:not(:disabled) {
		border-color: #e2e8f0;
	}

	.input:disabled {
		background-color: #f8fafc;
		color: #94a3b8;
		cursor: not-allowed;
	}

	.input-sm {
		padding: 0.375rem 0.75rem;
		font-size: 0.75rem;
	}

	.input-lg {
		padding: 0.875rem 1.25rem;
		font-size: 1rem;
	}

	.input:focus {
		outline: none;
		border-color: #6366f1;
		box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.08);
		background-color: #ffffff;
	}

	.input:focus-visible {
		outline: 2px solid #6366f1;
		outline-offset: 2px;
	}

	.input::placeholder {
		color: #94a3b8;
	}

	.input-error {
		border-color: #ef4444;
	}

	.input-error:focus {
		border-color: #ef4444;
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.08);
	}

	.input-helper {
		font-size: 0.75rem;
		color: #64748b;
		line-height: 1.4;
	}

	.input-error-text {
		font-size: 0.75rem;
		color: #ef4444;
		line-height: 1.4;
	}
</style>
