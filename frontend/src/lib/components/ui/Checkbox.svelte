<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends Omit<HTMLInputAttributes, 'size'> {
		label?: string;
		checked?: boolean;
	}

	let {
		label,
		class: className = '',
		checked = $bindable(false),
		...restProps
	}: Props = $props();
</script>

<div class="checkbox-wrapper {className}">
	<label class="checkbox-label">
		<input
			type="checkbox"
			{...restProps}
			bind:checked
			class:checkbox={true}
		/>
		{#if label}
			<span class="checkbox-text">{label}</span>
		{/if}
	</label>
</div>

<style>
	.checkbox-wrapper {
		display: inline-flex;
	}

	.checkbox-label {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		cursor: pointer;
		user-select: none;
	}

	.checkbox {
		width: 16px;
		height: 16px;
		border: 1px solid #e5e5e5;
		border-radius: 4px;
		background-color: white;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.checkbox:checked {
		background-color: #3b82f6;
		border-color: #3b82f6;
		background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 16 16' fill='white' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M12.207 4.793a1 1 0 010 1.414l-5 5a1 1 0 01-1.414 0l-2-2a1 1 0 011.414-1.414L6.5 9.086l4.293-4.293a1 1 0 011.414 0z'/%3E%3C/svg%3E");
		background-position: center;
		background-repeat: no-repeat;
		background-size: 12px;
	}

	.checkbox:focus {
		outline: none;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.checkbox:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.checkbox-text {
		font-size: 0.875rem;
		color: #171717;
	}
</style>
