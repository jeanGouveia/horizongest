<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
</script>

<div class="platform-layout">
	<slot />
</div>

<style>
	.platform-layout {
		min-height: 100vh;
		background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
	}
</style>
