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
		locale: string;
		currency: string;
		timezone: string;
		created_at: string;
		updated_at: string;
	}

	let company: Company | null = null;
	let loading = true;
	let editing = false;
	let showEditModal = false;

	let editForm = {
		name: '',
		slug: '',
		description: '',
		business_type: '',
		locale: '',
		currency: '',
		timezone: ''
	};

	const businessTypes = [
		{ value: 'restaurant', label: 'Restaurante' },
		{ value: 'cafe', label: 'Café' },
		{ value: 'bakery', label: 'Padaria' },
		{ value: 'other', label: 'Outro' }
	];

	onMount(async () => {
		await loadCompany();
	});

	async function loadCompany() {
		loading = true;
		try {
			const id = window.location.pathname.split('/').pop();
			const response = await fetch(`/api/platform/companies/${id}`);
			if (response.ok) {
				company = await response.json();
			} else {
				showError('Erro', 'Erro ao carregar empresa');
				goto('/platform/companies');
			}
		} catch (error) {
			showError('Erro', 'Erro ao carregar empresa');
			goto('/platform/companies');
		} finally {
			loading = false;
		}
	}

	function openEditModal() {
		if (!company) return;
		editForm = {
			name: company.name,
			slug: company.slug,
			description: company.description,
			business_type: company.business_type,
			locale: company.locale,
			currency: company.currency,
			timezone: company.timezone
		};
		showEditModal = true;
	}

	async function saveCompany() {
		if (!company) return;

		try {
			const response = await fetch(`/api/platform/companies/${company.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(editForm)
			});

			if (response.ok) {
				showSuccess('Sucesso', 'Empresa atualizada com sucesso');
				showEditModal = false;
				await loadCompany();
			} else {
				const data = await response.json();
				showError('Erro', data.message || 'Erro ao atualizar empresa');
			}
		} catch (error) {
			showError('Erro', 'Erro ao atualizar empresa');
		}
	}

	async function toggleStatus() {
		if (!company) return;

		const action = company.active ? 'deactivate' : 'activate';
		try {
			const response = await fetch(`/api/platform/companies/${company.id}/${action}`, {
				method: 'POST'
			});

			if (response.ok) {
				showSuccess('Sucesso', company.active ? 'Empresa desativada' : 'Empresa ativada');
				await loadCompany();
			} else {
				showError('Erro', 'Erro ao alterar status da empresa');
			}
		} catch (error) {
			showError('Erro', 'Erro ao alterar status da empresa');
		}
	}

	async function loginAsCompany() {
		if (!company) return;

		try {
			const response = await fetch(`/api/platform/companies/${company.id}/login-as`, {
				method: 'POST'
			});

			if (response.ok) {
				const data = await response.json();
				showSuccess('Sucesso', `Login como empresa iniciado. Owner: ${data.owner_email}`);
				// In a real implementation, this would redirect to the company dashboard with a temporary token
			} else {
				showError('Erro', 'Erro ao fazer login como empresa');
			}
		} catch (error) {
			showError('Erro', 'Erro ao fazer login como empresa');
		}
	}
</script>

<div class="company-detail-container">
	<header class="company-header">
		<div class="header-content">
			<Button variant="secondary" on:click={() => goto('/platform/companies')}>← Voltar</Button>
			<div>
				<h1>{company?.name || 'Carregando...'}</h1>
				<p>Detalhes da empresa</p>
			</div>
		</div>
		<div class="header-actions">
			<Button variant="secondary" on:click={loginAsCompany}>Login como Empresa</Button>
			<Button variant="secondary" on:click={() => goto(`/platform/companies/${company?.id}/owner`)}>Gerenciar Owner</Button>
			<Button variant="secondary" on:click={openEditModal}>Editar</Button>
			<Button
				variant={company?.active ? 'secondary' : 'primary'}
				on:click={toggleStatus}
			>
				{company?.active ? 'Desativar' : 'Ativar'}
			</Button>
		</div>
	</header>

	{#if loading}
		<div class="loading-state">
			<p>Carregando...</p>
		</div>
	{:else if company}
		<div class="company-details">
			<Card class="detail-card">
				<h2>Informações Básicas</h2>
				<div class="detail-row">
					<strong>Nome:</strong>
					<span>{company.name}</span>
				</div>
				<div class="detail-row">
					<strong>Slug:</strong>
					<span>{company.slug}</span>
				</div>
				<div class="detail-row">
					<strong>Status:</strong>
					<Badge variant={company.active ? 'success' : 'error'}>
						{company.active ? 'Ativa' : 'Inativa'}
					</Badge>
				</div>
				<div class="detail-row">
					<strong>Descrição:</strong>
					<span>{company.description || 'Sem descrição'}</span>
				</div>
			</Card>

			<Card class="detail-card">
				<h2>Configurações de Negócio</h2>
				<div class="detail-row">
					<strong>Tipo de Negócio:</strong>
					<span>{company.business_type || 'Não definido'}</span>
				</div>
				<div class="detail-row">
					<strong>Locale:</strong>
					<span>{company.locale || 'Não definido'}</span>
				</div>
				<div class="detail-row">
					<strong>Moeda:</strong>
					<span>{company.currency || 'Não definido'}</span>
				</div>
				<div class="detail-row">
					<strong>Fuso Horário:</strong>
					<span>{company.timezone || 'Não definido'}</span>
				</div>
			</Card>

			<Card class="detail-card">
				<h2>Informações do Sistema</h2>
				<div class="detail-row">
					<strong>ID:</strong>
					<span>{company.id}</span>
				</div>
				<div class="detail-row">
					<strong>Criada em:</strong>
					<span>{new Date(company.created_at).toLocaleString('pt-BR')}</span>
				</div>
				<div class="detail-row">
					<strong>Atualizada em:</strong>
					<span>{new Date(company.updated_at).toLocaleString('pt-BR')}</span>
				</div>
			</Card>
		</div>
	{/if}

	{#if showEditModal}
		<div class="modal-overlay" on:click={() => (showEditModal = false)}>
			<Card class="modal" on:click={(e) => e.stopPropagation()}>
				<div class="modal-header">
					<h3>Editar Empresa</h3>
					<Button variant="ghost" size="sm" on:click={() => (showEditModal = false)}>✕</Button>
				</div>
				<div class="modal-body">
					<div class="form-group">
						<label>Nome</label>
						<Input type="text" bind:value={editForm.name} />
					</div>
					<div class="form-group">
						<label>Slug</label>
						<Input type="text" bind:value={editForm.slug} />
					</div>
					<div class="form-group">
						<label>Descrição</label>
						<Textarea bind:value={editForm.description} />
					</div>
					<div class="form-group">
						<label>Tipo de Negócio</label>
						<Select bind:value={editForm.business_type}>
							{#each businessTypes as type}
								<option value={type.value}>{type.label}</option>
							{/each}
						</Select>
					</div>
					<div class="form-group">
						<label>Locale</label>
						<Input type="text" bind:value={editForm.locale} placeholder="pt-BR" />
					</div>
					<div class="form-group">
						<label>Moeda</label>
						<Input type="text" bind:value={editForm.currency} placeholder="BRL" />
					</div>
					<div class="form-group">
						<label>Fuso Horário</label>
						<Input type="text" bind:value={editForm.timezone} placeholder="America/Sao_Paulo" />
					</div>
				</div>
				<div class="modal-footer">
					<Button variant="secondary" on:click={() => (showEditModal = false)}>Cancelar</Button>
					<Button variant="primary" on:click={saveCompany}>Salvar</Button>
				</div>
			</Card>
		</div>
	{/if}
</div>

<style>
	.company-detail-container {
		width: 100%;
		max-width: 1200px;
		padding: 2rem;
	}

	.company-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	.header-content {
		display: flex;
		align-items: center;
		gap: 1rem;
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

	.header-actions {
		display: flex;
		gap: 0.75rem;
	}

	.loading-state {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 4rem;
		color: white;
	}

	.company-details {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
		gap: 1.5rem;
	}

	.detail-card {
		background: rgba(255, 255, 255, 0.1);
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 12px;
		padding: 1.5rem;
	}

	.detail-card h2 {
		font-size: 1.125rem;
		font-weight: 600;
		color: white;
		margin: 0 0 1rem 0;
	}

	.detail-row {
		display: flex;
		justify-content: space-between;
		padding: 0.75rem 0;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.detail-row:last-child {
		border-bottom: none;
	}

	.detail-row strong {
		color: #cbd5e1;
		font-size: 0.875rem;
	}

	.detail-row span {
		color: #94a3b8;
		font-size: 0.875rem;
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
		max-width: 500px;
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
		gap: 1rem;
		margin-bottom: 1.5rem;
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
