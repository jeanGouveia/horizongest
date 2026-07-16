<script lang="ts">
	interface Props {
		class?: string;
		zebra?: boolean;
		hover?: boolean;
		compact?: boolean;
		bordered?: boolean;
	}

	let {
		class: className = '',
		zebra = false,
		hover = true,
		compact = false,
		bordered = true,
	}: Props = $props();
</script>

<div class="table-wrapper {className}">
	<table class:table-zebra={zebra} class:table-hover={hover} class:table-compact={compact} class:table-bordered={bordered}>
		<slot />
	</table>
</div>

<style>
	.table-wrapper {
		overflow-x: auto;
		border: 1px solid #f1f5f9;
		border-radius: 16px;
		background-color: #ffffff;
		box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.05);
	}

	.table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.875rem;
	}

	.table th {
		padding: 1rem 1.25rem;
		text-align: left;
		font-weight: 600;
		color: #64748b;
		background-color: #f8fafc;
		border-bottom: 1px solid #f1f5f9;
		white-space: nowrap;
		font-size: 0.6875rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.table td {
		padding: 1rem 1.25rem;
		color: #0f172a;
		border-bottom: 1px solid #f1f5f9;
		line-height: 1.5;
	}

	.table-compact th {
		padding: 0.75rem 1rem;
	}

	.table-compact td {
		padding: 0.75rem 1rem;
	}

	.table-bordered tbody tr td {
		border-bottom: 1px solid #f1f5f9;
	}

	.table-bordered tbody tr:last-child td {
		border-bottom: none;
	}

	.table-zebra tbody tr:nth-child(even) {
		background-color: #fafbfc;
	}

	.table-hover tbody tr {
		transition: background-color 0.15s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.table-hover tbody tr:hover {
		background-color: #f8fafc;
	}

	.table-hover tbody tr:hover td {
		color: #0f172a;
	}
</style>
