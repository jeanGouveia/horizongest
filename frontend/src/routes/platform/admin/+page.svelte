<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card } from '$lib/components/ui';
	import { Button } from '$lib/components/ui';
	import { showSuccess, showError } from '$lib/stores/toast';

	let stats = {
		totalCompanies: 0,
		totalOwners: 0,
		totalUsers: 0,
		blockedCompanies: 0,
		trialCompanies: 0,
		paidCompanies: 0
	};

	let loading = true;

	onMount(async () => {
		await loadDashboardStats();
	});

	async function loadDashboardStats() {
		loading = true;
		try {
			const response = await fetch('http://localhost:8080/api/platform/dashboard/stats');
			if (response.ok) {
				const data = await response.json();
				stats = data;
			}
		} catch (error) {
			if (error instanceof TypeError && error.message.includes('fetch')) {
				showError('Erro de conexão', 'Não foi possível conectar ao servidor. Verifique se o backend está rodando em http://localhost:8080');
			} else {
				showError('Erro', 'Erro ao carregar estatísticas');
			}
		} finally {
			loading = false;
		}
	}

	function logout() {
		document.cookie = 'platform_auth_token=; path=/; max-age=0';
		goto('/platform/signin');
	}
</script>

<div class="dashboard-container">
	<header class="dashboard-header">
		<div class="header-content">
			<h1>Dashboard da Plataforma</h1>
			<p>Visão geral do sistema</p>
		</div>
		<div class="header-actions">
			<Button variant="secondary" onclick={() => goto('/platform/plans')}>Planos</Button>
			<Button variant="secondary" onclick={logout}>Sair</Button>
		</div>
	</header>

	{#if loading}
		<div class="loading-state">
			<p>Carregando...</p>
		</div>
	{:else}
		<div class="stats-grid">
			<Card class="stat-card">
				<div class="stat-icon companies">🏢</div>
				<div class="stat-content">
					<p class="stat-label">Empresas</p>
					<p class="stat-value">{stats.totalCompanies}</p>
				</div>
			</Card>

			<Card class="stat-card">
				<div class="stat-icon owners">👤</div>
				<div class="stat-content">
					<p class="stat-label">Owners</p>
					<p class="stat-value">{stats.totalOwners}</p>
				</div>
			</Card>

			<Card class="stat-card">
				<div class="stat-icon users">👥</div>
				<div class="stat-content">
					<p class="stat-label">Usuários</p>
					<p class="stat-value">{stats.totalUsers}</p>
				</div>
			</Card>

			<Card class="stat-card">
				<div class="stat-icon blocked">🔒</div>
				<div class="stat-content">
					<p class="stat-label">Bloqueadas</p>
					<p class="stat-value">{stats.blockedCompanies}</p>
				</div>
			</Card>

			<Card class="stat-card">
				<div class="stat-icon trial">⏳</div>
				<div class="stat-content">
					<p class="stat-label">Trial</p>
					<p class="stat-value">{stats.trialCompanies}</p>
				</div>
			</Card>

			<Card class="stat-card">
				<div class="stat-icon paid">💳</div>
				<div class="stat-content">
					<p class="stat-label">Pagas</p>
					<p class="stat-value">{stats.paidCompanies}</p>
				</div>
			</Card>
		</div>

		<div class="actions-grid">
			<Card class="action-card" onclick={() => goto('/platform/companies')}>
				<h3>Gerenciar Empresas</h3>
				<p>Ver, criar, editar e desativar empresas</p>
				<Button variant="primary">Acessar</Button>
			</Card>

			<Card class="action-card">
				<h3>Gerenciar Owners</h3>
				<p>Gerenciar proprietários das empresas</p>
				<Button variant="secondary" disabled>Em breve</Button>
			</Card>

			<Card class="action-card">
				<h3>Planos</h3>
				<p>Gerenciar planos e assinaturas</p>
				<Button variant="secondary" disabled>Em breve</Button>
			</Card>

			<Card class="action-card">
				<h3>Auditoria</h3>
				<p>Visualizar logs de ações</p>
				<Button variant="secondary" disabled>Em breve</Button>
			</Card>
		</div>
	{/if}
</div>

<style>
	.dashboard-container {
		width: 100%;
		max-width: 1200px;
		padding: 2rem;
	}

	.dashboard-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	.header-content h1 {
		font-size: 1.75rem;
		font-weight: 700;
		color: white;
		margin: 0 0 0.5rem 0;
	}

	.header-content p {
		font-size: 0.875rem;
		color: #94a3b8;
		margin: 0;
	}

	.loading-state {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 4rem;
		color: white;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 1.5rem;
		margin-bottom: 2rem;
	}

	.stat-card {
		background: rgba(255, 255, 255, 0.1);
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 12px;
		padding: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.stat-icon {
		font-size: 2.5rem;
		width: 60px;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 12px;
	}

	.stat-content {
		flex: 1;
	}

	.stat-label {
		font-size: 0.875rem;
		color: #94a3b8;
		margin: 0 0 0.25rem 0;
	}

	.stat-value {
		font-size: 1.75rem;
		font-weight: 700;
		color: white;
		margin: 0;
	}

	.actions-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 1.5rem;
	}

	.action-card {
		background: rgba(255, 255, 255, 0.1);
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 12px;
		padding: 1.5rem;
		cursor: pointer;
		transition: transform 0.2s, background 0.2s;
	}

	.action-card:hover {
		transform: translateY(-2px);
		background: rgba(255, 255, 255, 0.15);
	}

	.action-card h3 {
		font-size: 1.125rem;
		font-weight: 600;
		color: white;
		margin: 0 0 0.5rem 0;
	}

	.action-card p {
		font-size: 0.875rem;
		color: #94a3b8;
		margin: 0 0 1rem 0;
	}

	.action-card button {
		width: 100%;
	}
</style>
