<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card } from '$lib/components/ui';
	import { Button } from '$lib/components/ui';
	import { Input } from '$lib/components/ui';
	import { Textarea } from '$lib/components/ui';
	import { Select } from '$lib/components/ui';
	import { Badge } from '$lib/components/ui';
	import { showSuccess, showError } from '$lib/stores/toast';

	interface Company {
		id: number;
		name: string;
		slug: string;
		description: string;
		active: boolean;
		business_type: string;
		created_at: string;
	}

	let companies: Company[] = [];
	let loading = true;
	let searchQuery = '';
	let showCreateModal = false;

	let createForm = {
		name: '',
		slug: '',
		description: '',
		business_type: 'restaurant',
		locale: 'pt-BR',
		currency: 'BRL',
		timezone: 'America/Sao_Paulo',
		owner_email: '',
		owner_name: '',
		password: ''
	};

	const businessTypes = [
		{ value: 'restaurant', label: 'Restaurante' },
		{ value: 'cafe', label: 'Café' },
		{ value: 'bakery', label: 'Padaria' },
		{ value: 'other', label: 'Outro' }
	];

	onMount(async () => {
		await loadCompanies();
	});

	async function loadCompanies() {
		loading = true;
		try {
			const response = await fetch('/api/platform/companies');
			if (response.ok) {
				const data = await response.json();
				companies = data.companies || [];
			} else {
				showError('Erro', 'Erro ao carregar empresas');
			}
		} catch (error) {
			showError('Erro', 'Erro ao carregar empresas');
		} finally {
			loading = false;
		}
	}

	async function toggleCompanyStatus(company: Company) {
		const action = company.active ? 'deactivate' : 'activate';
		try {
			const response = await fetch(`/api/platform/companies/${company.id}/${action}`, {
				method: 'POST'
			});

			if (response.ok) {
				showSuccess('Sucesso', company.active ? 'Empresa desativada' : 'Empresa ativada');
				await loadCompanies();
			} else {
				showError('Erro', 'Erro ao alterar status da empresa');
			}
		} catch (error) {
			showError('Erro', 'Erro ao alterar status da empresa');
		}
	}

	function viewCompany(company: Company) {
		goto(`/platform/companies/${company.id}`);
	}

	async function createCompany() {
		if (!createForm.name || !createForm.slug || !createForm.owner_email || !createForm.owner_name || !createForm.password) {
			showError('Erro', 'Preencha todos os campos obrigatórios');
			return;
		}

		try {
			const response = await fetch('/api/platform/companies', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					company: {
						name: createForm.name,
						slug: createForm.slug,
						description: createForm.description,
						business_type: createForm.business_type,
						locale: createForm.locale,
						currency: createForm.currency,
						timezone: createForm.timezone
					},
					owner_email: createForm.owner_email,
					owner_name: createForm.owner_name,
					password: createForm.password
				})
			});

			if (response.ok) {
				showSuccess('Sucesso', 'Empresa criada com sucesso');
				showCreateModal = false;
				createForm = {
					name: '',
					slug: '',
					description: '',
					business_type: 'restaurant',
					locale: 'pt-BR',
					currency: 'BRL',
					timezone: 'America/Sao_Paulo',
					owner_email: '',
					owner_name: '',
					password: ''
				};
				await loadCompanies();
			} else {
				const data = await response.json();
				showError('Erro', data.message || 'Erro ao criar empresa');
			}
		} catch (error) {
			showError('Erro', 'Erro ao criar empresa');
		}
	}

	$: filteredCompanies = companies.filter(
		(c) =>
			c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			c.slug.toLowerCase().includes(searchQuery.toLowerCase())
	);
</script>

<div class="companies-container">
	<header class="companies-header">
		<div class="header-content">
			<h1>Empresas</h1>
			<p>Gerencie todas as empresas da plataforma</p>
		</div>
		<Button variant="primary" on:click={() => (showCreateModal = true)}>Nova Empresa</Button>
	</header>

	<div class="search-bar">
		<Input
			type="text"
			placeholder="Buscar por nome ou slug..."
			bind:value={searchQuery}
		/>
	</div>

	{#if loading}
		<div class="loading-state">
			<p>Carregando empresas...</p>
		</div>
	{:else if filteredCompanies.length === 0}
		<div class="empty-state">
			<p>Nenhuma empresa encontrada</p>
		</div>
	{:else}
		<div class="companies-grid">
			{#each filteredCompanies as company}
				<Card class="company-card">
					<div class="company-header">
						<h3>{company.name}</h3>
						<Badge variant={company.active ? 'success' : 'error'}>
							{company.active ? 'Ativa' : 'Inativa'}
						</Badge>
					</div>

					<div class="company-info">
						<p><strong>Slug:</strong> {company.slug}</p>
						<p><strong>Tipo:</strong> {company.business_type || 'Não definido'}</p>
						<p><strong>Criada em:</strong> {new Date(company.created_at).toLocaleDateString('pt-BR')}</p>
					</div>

					<div class="company-actions">
						<Button variant="secondary" size="sm" on:click={() => viewCompany(company)}>
							Ver Detalhes
						</Button>
						<Button
							variant={company.active ? 'secondary' : 'primary'}
							size="sm"
							on:click={() => toggleCompanyStatus(company)}
						>
							{company.active ? 'Desativar' : 'Ativar'}
						</Button>
					</div>
				</Card>
			{/each}
		</div>
	{/if}

	{#if showCreateModal}
		<div class="modal-overlay" on:click={() => (showCreateModal = false)}>
			<Card class="modal" on:click={(e) => e.stopPropagation()}>
				<div class="modal-header">
					<h3>Nova Empresa</h3>
					<Button variant="ghost" size="sm" on:click={() => (showCreateModal = false)}>✕</Button>
				</div>
				<div class="modal-body">
					<div class="form-section">
						<h4>Dados da Empresa</h4>
						<div class="form-group">
							<label>Nome *</label>
							<Input type="text" bind:value={createForm.name} placeholder="Nome da empresa" />
						</div>
						<div class="form-group">
							<label>Slug *</label>
							<Input type="text" bind:value={createForm.slug} placeholder="identificador-unico" />
						</div>
						<div class="form-group">
							<label>Descrição</label>
							<Textarea bind:value={createForm.description} placeholder="Descrição da empresa" />
						</div>
						<div class="form-group">
							<label>Tipo de Negócio</label>
							<Select bind:value={createForm.business_type}>
								{#each businessTypes as type}
									<option value={type.value}>{type.label}</option>
								{/each}
							</Select>
						</div>
						<div class="form-group">
							<label>Locale</label>
							<Input type="text" bind:value={createForm.locale} placeholder="pt-BR" />
						</div>
						<div class="form-group">
							<label>Moeda</label>
							<Input type="text" bind:value={createForm.currency} placeholder="BRL" />
						</div>
						<div class="form-group">
							<label>Fuso Horário</label>
							<Input type="text" bind:value={createForm.timezone} placeholder="America/Sao_Paulo" />
						</div>
					</div>

					<div class="form-section">
						<h4>Dados do Owner</h4>
						<div class="form-group">
							<label>Nome do Owner *</label>
							<Input type="text" bind:value={createForm.owner_name} placeholder="Nome completo" />
						</div>
						<div class="form-group">
							<label>E-mail do Owner *</label>
							<Input type="email" bind:value={createForm.owner_email} placeholder="owner@empresa.com" />
						</div>
						<div class="form-group">
							<label>Senha Temporária *</label>
							<Input type="password" bind:value={createForm.password} placeholder="Mínimo 8 caracteres" />
						</div>
					</div>
				</div>
				<div class="modal-footer">
					<Button variant="secondary" on:click={() => (showCreateModal = false)}>Cancelar</Button>
					<Button variant="primary" on:click={createCompany}>Criar Empresa</Button>
				</div>
			</Card>
		</div>
	{/if}
</div>

<style>
	.companies-container {
		width: 100%;
		max-width: 1200px;
		padding: 2rem;
	}

	.companies-header {
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

	.search-bar {
		margin-bottom: 2rem;
	}

	.loading-state,
	.empty-state {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 4rem;
		color: white;
	}

	.companies-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
		gap: 1.5rem;
	}

	.company-card {
		background: rgba(255, 255, 255, 0.1);
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 12px;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.company-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.company-header h3 {
		font-size: 1.125rem;
		font-weight: 600;
		color: white;
		margin: 0;
	}

	.company-info {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.company-info p {
		font-size: 0.875rem;
		color: #94a3b8;
		margin: 0;
	}

	.company-info strong {
		color: #cbd5e1;
	}

	.company-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: auto;
	}

	.company-actions button {
		flex: 1;
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: #1e293b;
		border-radius: 12px;
		padding: 2rem;
		width: 100%;
		max-width: 600px;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.modal-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: white;
		margin: 0;
	}

	.modal-body {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		margin-bottom: 1.5rem;
	}

	.form-section h4 {
		font-size: 1rem;
		font-weight: 600;
		color: white;
		margin: 0 0 1rem 0;
		padding-bottom: 0.5rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #cbd5e1;
	}

	.modal-footer {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
	}
</style>
