<script lang="ts">
	import { Search } from '@lucide/svelte';

	interface Props {
		breadcrumb?: string;
	}

	let {
		breadcrumb = '',
	}: Props = $props();

	let searchQuery = $state('');

	function handleSearch(e: KeyboardEvent) {
		if (e.key === 'Enter' && searchQuery.trim()) {
			// Implementar busca
			console.log('Searching for:', searchQuery);
		}
	}
</script>

<header class="header">
	<div class="header-left">
		{#if breadcrumb}
			<span class="breadcrumb">{breadcrumb}</span>
		{/if}
	</div>

	<div class="header-center">
		<div class="search-container">
			<Search size={14} class="search-icon" />
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Buscar..."
				class="search-input"
				onkeydown={handleSearch}
			/>
		</div>
	</div>

	<div class="header-right">
		<!-- Spacer for balance -->
	</div>
</header>

<style>
	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 0.5rem 1rem;
		background-color: #ffffff;
		border-bottom: 1px solid #e2e8f0;
		position: sticky;
		top: 0;
		z-index: 100;
		height: 40px;
	}

	.header-left {
		display: flex;
		align-items: center;
		flex: 0 0 auto;
		min-width: 0;
	}

	.breadcrumb {
		font-size: 0.75rem;
		font-weight: 500;
		color: #64748b;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.header-center {
		flex: 1;
		display: flex;
		justify-content: center;
		max-width: 400px;
	}

	.search-container {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
		padding: 0.375rem 0.75rem;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 6px;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.search-container:focus-within {
		background: #ffffff;
		border-color: #6366f1;
		box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.1);
	}

	.search-icon {
		color: #94a3b8;
		flex-shrink: 0;
	}

	.search-input {
		flex: 1;
		border: none;
		background: transparent;
		font-size: 0.8125rem;
		color: #0f172a;
		outline: none;
	}

	.search-input::placeholder {
		color: #94a3b8;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex: 0 0 auto;
		width: 100px;
	}

	@media (max-width: 768px) {
		.header {
			padding: 0.5rem 0.75rem;
		}

		.header-center {
			max-width: 180px;
		}

		.search-input::placeholder {
			font-size: 0.75rem;
		}

		.header-right {
			width: 50px;
		}
	}
</style>
