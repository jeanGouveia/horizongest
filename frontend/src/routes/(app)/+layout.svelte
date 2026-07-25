<script lang="ts">
	import { Header } from '$lib/components/layout';
	import { Sidebar } from '$lib/components/layout';
	import { Footer } from '$lib/components/layout';
	import ImpersonationBanner from '$lib/components/ImpersonationBanner.svelte';
	import { page } from '$app/stores';
	import { userStore } from '$lib/stores/userStore.svelte';
	import { themeStore } from '$lib/stores/themeStore.svelte';
	import { companyStore } from '$lib/stores/companyStore.svelte';
	import { api } from '$lib/api/client';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';

	let { children } = $props();
	let sidebarCollapsed = $state(false);
	let sidebarOpen = $state(false);
	let showMenuButton = $state(false);

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
		if (path === '/settings/company') return 'Empresa';
		if (path === '/settings/users') return 'Usuários';
		if (path === '/settings/invitations') return 'Convites';
		if (path === '/profile') return 'Perfil';
		return '';
	});

	onMount(async () => {
		if (browser) {
			showMenuButton = window.innerWidth < 768;
			// Load theme on app initialization
			themeStore.loadTheme();

			// Buscar configurações da empresa usando o proxy (como no perfil)
			try {
				const response = await fetch('/api/company/settings', {
					credentials: 'include'
				});

				if (response.ok) {
					const data = await response.json();
					const companyName = data.Name || data.name;
					if (companyName) {
						companyStore.setCompany({ name: companyName });
					}
				}
			} catch (e) {
				// Error silently ignored - company will use default name
			}
		}
	});

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}
</script>

<div class="app-layout">
	<ImpersonationBanner />
	<Header
		breadcrumb={breadcrumb()}
		showMenuButton={showMenuButton}
		onMenuToggle={toggleSidebar}
	/>
	<div class="main-content">
		<Sidebar
			currentPath={currentPath}
			collapsed={sidebarCollapsed}
			onToggle={() => sidebarCollapsed = !sidebarCollapsed}
			userName={userStore.user?.name}
			userAvatar={userStore.user?.avatar}
			open={sidebarOpen}
			onClose={() => sidebarOpen = false}
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
