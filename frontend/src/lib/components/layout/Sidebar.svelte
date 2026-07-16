<script lang="ts">
	import { LayoutDashboard, Plus, ShoppingCart, Utensils, Leaf, Scale, Users, User, Settings, LogOut, ChevronLeft, ChevronRight, Bell } from '@lucide/svelte';

	interface NavItem {
		label: string;
		href: string;
		icon: any;
		badge?: number;
	}

	interface NavGroup {
		title: string;
		items: NavItem[];
	}

	interface Props {
		currentPath?: string;
		collapsed?: boolean;
		onToggle?: () => void;
		userName?: string;
		userAvatar?: string;
	}

	let {
		currentPath = '/',
		collapsed = false,
		onToggle,
		userName = 'Usuário',
		userAvatar,
	}: Props = $props();

	const navGroups: NavGroup[] = [
		{
			title: 'OPERAÇÃO',
			items: [
				{ label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
				{ label: 'Novo Pedido', href: '/orders/new', icon: Plus },
				{ label: 'Pedidos', href: '/orders', icon: ShoppingCart, badge: 3 },
			],
		},
		{
			title: 'CATÁLOGO',
			items: [
				{ label: 'Produtos', href: '/products', icon: Utensils },
				{ label: 'Ingredientes', href: '/ingredients', icon: Leaf },
			],
		},
		{
			title: 'ESTOQUE',
			items: [
				{ label: 'Ajustes', href: '/stock-adjustments', icon: Scale, badge: 2 },
			],
		},
		{
			title: 'ADMINISTRAÇÃO',
			items: [
				{ label: 'Usuários', href: '/users', icon: Users },
				{ label: 'Perfil', href: '/profile', icon: User },
				{ label: 'Configurações', href: '/settings', icon: Settings },
			],
		},
	];

	function isActive(href: string): boolean {
		return currentPath === href || currentPath.startsWith(href + '/');
	}
</script>

<aside class="sidebar {collapsed ? 'sidebar-collapsed' : ''}">
	<div class="sidebar-header">
		<button class="sidebar-toggle" onclick={onToggle} aria-label="Toggle sidebar">
			{#if collapsed}
				<ChevronRight size={14} class="toggle-icon" />
			{:else}
				<ChevronLeft size={14} class="toggle-icon" />
			{/if}
		</button>
		{#if !collapsed}
			<div class="sidebar-brand">
				<span class="brand-text">PratoOnline</span>
			</div>
		{/if}
	</div>

	<nav class="sidebar-nav" aria-label="Navegação principal">
		{#each navGroups as group}
			{#if !collapsed}
				<div class="nav-group">
					<span class="nav-group-title">{group.title}</span>
					<ul class="nav-group-list">
						{#each group.items as item}
							<li class="nav-item">
								<a
									href={item.href}
									class="nav-link {isActive(item.href) ? 'nav-link-active' : ''}"
									aria-current={isActive(item.href) ? 'page' : undefined}
								>
									<div class="nav-link-content">
										<svelte:component this={item.icon} size={18} class="nav-icon" />
										<span class="nav-label">{item.label}</span>
									</div>
									{#if item.badge}
										<span class="nav-badge">{item.badge}</span>
									{/if}
								</a>
							</li>
						{/each}
					</ul>
				</div>
			{:else}
				<ul class="nav-group-list nav-group-list-collapsed">
					{#each group.items as item}
						<li class="nav-item">
							<a
								href={item.href}
								class="nav-link {isActive(item.href) ? 'nav-link-active' : ''}"
								aria-current={isActive(item.href) ? 'page' : undefined}
								title={item.label}
							>
								<svelte:component this={item.icon} size={18} class="nav-icon" />
								{#if item.badge}
									<span class="nav-badge nav-badge-collapsed">{item.badge}</span>
								{/if}
							</a>
						</li>
					{/each}
				</ul>
			{/if}
		{/each}
	</nav>

	<div class="sidebar-footer">
		<div class="sidebar-footer-content">
			<div class="user-section">
				{#if userAvatar}
					<img src={userAvatar} alt={userName} class="user-avatar" />
				{:else}
					<div class="user-avatar-placeholder">{userName.charAt(0).toUpperCase()}</div>
				{/if}
				{#if !collapsed}
					<div class="user-info">
						<div class="user-name">{userName}</div>
					</div>
				{/if}
			</div>
			<div class="sidebar-actions">
				<button class="sidebar-action-btn" title="Notificações">
					<Bell size={14} />
				</button>
				<a href="/settings" class="sidebar-action-btn" title="Configurações">
					<Settings size={14} />
				</a>
				<a href="/logout" class="sidebar-action-btn danger" title="Sair">
					<LogOut size={14} />
				</a>
			</div>
		</div>
	</div>
</aside>

<style>
	.sidebar {
		display: flex;
		flex-direction: column;
		width: 240px;
		background-color: #ffffff;
		border-right: 1px solid #e2e8f0;
		position: sticky;
		top: 0;
		height: 100vh;
		overflow-y: auto;
		transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.sidebar-collapsed {
		width: 64px;
	}

	.sidebar-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		border-bottom: 1px solid #e2e8f0;
	}

	.sidebar-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		background: transparent;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		color: #64748b;
	}

	.sidebar-toggle:hover {
		background: #f1f5f9;
		color: #0f172a;
	}

	.toggle-icon {
		flex-shrink: 0;
	}

	.sidebar-brand {
		flex: 1;
	}

	.brand-text {
		font-size: 0.875rem;
		font-weight: 600;
		color: #0f172a;
		letter-spacing: -0.025em;
	}

	.sidebar-nav {
		flex: 1;
		padding: 0.5rem 0;
		overflow-y: auto;
	}

	.nav-group {
		margin-bottom: 0.75rem;
	}

	.nav-group:last-child {
		margin-bottom: 0;
	}

	.nav-group-title {
		display: block;
		padding: 0 0.75rem;
		margin-bottom: 0.375rem;
		font-size: 0.625rem;
		font-weight: 600;
		color: #64748b;
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.nav-group-list {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.nav-group-list-collapsed {
		padding: 0 0.75rem;
	}

	.nav-item {
		margin: 0;
		padding: 0;
	}

	.nav-link {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.375rem 0.75rem;
		margin: 0 0.5rem 0.125rem 0.5rem;
		color: #64748b;
		text-decoration: none;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		border-radius: 6px;
		font-size: 0.8125rem;
		font-weight: 500;
		position: relative;
	}

	.nav-link:hover {
		background: #f1f5f9;
		color: #0f172a;
	}

	.nav-link-active {
		background: #eef2ff;
		color: #6366f1;
	}

	.nav-link-active::before {
		content: '';
		position: absolute;
		left: 0.25rem;
		top: 50%;
		transform: translateY(-50%);
		width: 2px;
		height: 16px;
		background: #6366f1;
		border-radius: 1px;
	}

	.nav-link-content {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.nav-icon {
		flex-shrink: 0;
	}

	.nav-label {
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.nav-badge {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 18px;
		height: 18px;
		padding: 0 0.25rem;
		background: #6366f1;
		color: white;
		font-size: 0.625rem;
		font-weight: 600;
		border-radius: 9px;
	}

	.nav-badge-collapsed {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
	}

	.sidebar-collapsed .nav-link {
		justify-content: center;
		padding: 0.5rem;
		margin: 0 0.5rem 0.125rem 0.5rem;
	}

	.sidebar-collapsed .nav-link {
		position: relative;
	}

	.sidebar-footer {
		padding: 0.5rem 0.75rem;
		border-top: 1px solid #e2e8f0;
	}

	.sidebar-footer-content {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
	}

	.user-section {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.5rem;
		border-radius: 6px;
		transition: background 0.15s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.user-section:hover {
		background: #f1f5f9;
	}

	.user-avatar,
	.user-avatar-placeholder {
		width: 24px;
		height: 24px;
		border-radius: 50%;
		object-fit: cover;
		border: 1px solid #e2e8f0;
		flex-shrink: 0;
	}

	.user-avatar-placeholder {
		background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.6875rem;
		font-weight: 600;
	}

	.user-info {
		flex: 1;
		min-width: 0;
	}

	.user-name {
		font-size: 0.75rem;
		font-weight: 500;
		color: #0f172a;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.sidebar-actions {
		display: flex;
		gap: 0.25rem;
	}

	.sidebar-action-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		background: transparent;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		color: #64748b;
	}

	.sidebar-action-btn:hover {
		background: #f1f5f9;
		color: #0f172a;
	}

	.sidebar-action-btn.danger:hover {
		background: #fef2f2;
		color: #ef4444;
	}

	.logout-link {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.5rem;
		color: #64748b;
		text-decoration: none;
		border-radius: 6px;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		font-size: 0.8125rem;
		font-weight: 500;
	}

	.logout-link:hover {
		background: #fef2f2;
		color: #ef4444;
	}

	.sidebar-collapsed .sidebar-footer-content {
		align-items: center;
	}

	.sidebar-collapsed .user-section {
		padding: 0.375rem;
	}

	.sidebar-collapsed .user-info,
	.sidebar-collapsed .logout-label {
		display: none;
	}

	.sidebar-collapsed .sidebar-actions {
		flex-direction: column;
		gap: 0.375rem;
	}

	.logout-icon {
		flex-shrink: 0;
	}

	.sidebar-collapsed .logout-link {
		justify-content: center;
		padding: 0.375rem;
	}

	@media (max-width: 1024px) {
		.sidebar {
			width: 64px;
		}

		.sidebar-brand,
		.nav-group-title,
		.nav-label,
		.user-info,
		.logout-label {
			display: none;
		}

		.nav-link {
			justify-content: center;
			padding: 0.5rem;
			margin: 0 0.5rem 0.125rem 0.5rem;
		}

		.nav-group-list {
			padding: 0 0.5rem;
		}

		.sidebar-footer-content {
			align-items: center;
		}

		.sidebar-actions {
			flex-direction: column;
			gap: 0.375rem;
		}
	}

	@media (max-width: 768px) {
		.sidebar {
			display: none;
		}
	}
</style>
