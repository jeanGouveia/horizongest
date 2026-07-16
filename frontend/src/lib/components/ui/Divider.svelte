<script lang="ts">
	interface Props {
		vertical?: boolean;
		class?: string;
	}

	let {
		vertical = false,
		class: className = '',
	}: Props = $props();
</script>

<div
	class:divider={true}
	class:divider-vertical={vertical}
	class:divider-horizontal={!vertical}
	class={className}
></div>

<style>
	.divider {
		border-color: #e5e5e5;
	}

	.divider-horizontal {
		width: 100%;
		height: 1px;
		border-top: 1px solid;
		margin: 16px 0;
	}

	.divider-vertical {
		width: 1px;
		height: 100%;
		border-left: 1px solid;
		margin: 0 16px;
	}
</style>
