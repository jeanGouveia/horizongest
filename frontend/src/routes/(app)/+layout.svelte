<script lang="ts">
	import { Header } from '$lib/components/layout';
	import { Sidebar } from '$lib/components/layout';
	import { Footer } from '$lib/components/layout';
	import { page } from '$app/stores';
	import { userStore } from '$lib/stores/userStore.svelte';

	let { children } = $props();
	let sidebarCollapsed = $state(false);

	// Get current path for sidebar active state
	let currentPath = $derived($page.url.pathname);

	// Generate breadcrumb from current path
	let breadcrumb = $derived(() => {
		const path = $page.url.pathname;
		if (path === '/dashboard') return 'Dashboard';
		if (path === '/orders') return 'Pedidos';
		if (path === '/orders/new') return 'Novo Pedido';
		if (path === '/products') return 'Produtos';
		if (path === '/ingredients') return 'Ingredientes';
		if (path === '/stock-adjustments') return 'Ajustes de Estoque';
		if (path === '/profile') return 'Perfil';
		if (path === '/settings') return 'Configurações';
		return '';
	});
</script>

<div class="app-layout">
	<Header breadcrumb={breadcrumb()} />
	<div class="main-content">
		<Sidebar 
			currentPath={currentPath} 
			collapsed={sidebarCollapsed} 
			onToggle={() => sidebarCollapsed = !sidebarCollapsed}
			userName={userStore.user?.name}
			userAvatar={userStore.user?.avatar}
		/>
		<main class="content">
			{@render children()}
		</main>
	</div>
	<Footer />
</div>

<style>
	.app-layout {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
		background-color: #f8fafc;
	}

	.main-content {
		display: flex;
		flex: 1;
		margin-top: 0px; /* Header height */
	}

	.content {
		flex: 1;
		padding: 1rem 1.25rem;
		overflow-y: auto;
		background-color: #f8fafc;
	}

	@media (max-width: 768px) {
		.content {
			padding: 0.75rem 1rem;
		}
	}
</style>
