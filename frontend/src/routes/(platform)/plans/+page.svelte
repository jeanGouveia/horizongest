<script lang="ts">
	import { onMount } from 'svelte';
	import { Card } from '$lib/components/ui';
	import { Button } from '$lib/components/ui';
	import { Input } from '$lib/components/ui';
	import { Textarea } from '$lib/components/ui';
	import { Select } from '$lib/components/ui';
	import { Badge } from '$lib/components/ui';
	import { showSuccess, showError } from '$lib/stores/toast';

	interface Plan {
		id: number;
		name: string;
		slug: string;
		description: string;
		price: number;
		currency: string;
		interval: string;
		max_users: number;
		max_products: number;
		features: string;
		active: boolean;
		created_at: string;
		updated_at: string;
	}

	let plans: Plan[] = [];
	let loading = true;
	let showCreateModal = false;
	let showEditModal = false;
	let editingPlan: Plan | null = null;

	let createForm = {
		name: '',
		slug: '',
		description: '',
		price: 0,
		currency: 'BRL',
		interval: 'monthly',
		max_users: 1,
		max_products: 100,
		features: ''
	};

	const intervals = [
		{ value: 'monthly', label: 'Mensal' },
		{ value: 'yearly', label: 'Anual' }
	];

	onMount(async () => {
		await loadPlans();
	});

	async function loadPlans() {
		loading = true;
		try {
			const response = await fetch('/api/platform/plans');
			if (response.ok) {
				const data = await response.json();
				plans = data.plans || [];
			} else {
				showError('Erro', 'Erro ao carregar planos');
			}
		} catch (error) {
			showError('Erro', 'Erro ao carregar planos');
		} finally {
			loading = false;
		}
	}

	async function createPlan() {
		if (!createForm.name || !createForm.slug) {
			showError('Erro', 'Preencha todos os campos obrigatórios');
			return;
		}

		try {
			const response = await fetch('/api/platform/plans', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(createForm)
			});

			if (response.ok) {
				showSuccess('Sucesso', 'Plano criado com sucesso');
				showCreateModal = false;
				createForm = {
					name: '',
					slug: '',
					description: '',
					price: 0,
					currency: 'BRL',
					interval: 'monthly',
					max_users: 1,
					max_products: 100,
					features: ''
				};
				await loadPlans();
			} else {
				const data = await response.json();
				showError('Erro', data.message || 'Erro ao criar plano');
			}
		} catch (error) {
			showError('Erro', 'Erro ao criar plano');
		}
	}

	async function togglePlanStatus(plan: Plan) {
		try {
			const response = await fetch(`/api/platform/plans/${plan.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ active: !plan.active })
			});

			if (response.ok) {
				showSuccess('Sucesso', plan.active ? 'Plano desativado' : 'Plano ativado');
				await loadPlans();
			} else {
				showError('Erro', 'Erro ao alterar status do plano');
			}
		} catch (error) {
			showError('Erro', 'Erro ao alterar status do plano');
		}
	}

	async function deletePlan(plan: Plan) {
		if (!confirm(`Tem certeza que deseja deletar o plano "${plan.name}"?`)) {
			return;
		}

		try {
			const response = await fetch(`/api/platform/plans/${plan.id}`, {
				method: 'DELETE'
			});

			if (response.ok) {
				showSuccess('Sucesso', 'Plano deletado com sucesso');
				await loadPlans();
			} else {
				showError('Erro', 'Erro ao deletar plano');
			}
		} catch (error) {
			showError('Erro', 'Erro ao deletar plano');
		}
	}

	function openEditModal(plan: Plan) {
		editingPlan = plan;
		showEditModal = true;
	}

	async function updatePlan() {
		if (!editingPlan) return;

		try {
			const response = await fetch(`/api/platform/plans/${editingPlan.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					name: editingPlan.name,
					description: editingPlan.description,
					price: editingPlan.price,
					currency: editingPlan.currency,
					interval: editingPlan.interval,
					max_users: editingPlan.max_users,
					max_products: editingPlan.max_products,
					features: editingPlan.features,
					active: editingPlan.active
				})
			});

			if (response.ok) {
				showSuccess('Sucesso', 'Plano atualizado com sucesso');
				showEditModal = false;
				editingPlan = null;
				await loadPlans();
			} else {
				const data = await response.json();
				showError('Erro', data.message || 'Erro ao atualizar plano');
			}
		} catch (error) {
			showError('Erro', 'Erro ao atualizar plano');
		}
	}
</script>

<div class="plans-container">
	<header class="plans-header">
		<h1>Planos</h1>
		<Button variant="primary" on:click={() => showCreateModal = true}>+ Novo Plano</Button>
	</header>

	{#if loading}
		<div class="loading-state">
			<p>Carregando...</p>
		</div>
	{:else}
		<div class="plans-grid">
			{#each plans as plan}
				<Card class="plan-card">
					<div class="plan-header">
						<h3>{plan.name}</h3>
						<Badge variant={plan.active ? 'success' : 'error'}>
							{plan.active ? 'Ativo' : 'Inativo'}
						</Badge>
					</div>
					<div class="plan-price">
						<span class="price">{plan.currency} {plan.price.toFixed(2)}</span>
						<span class="interval">/{plan.interval === 'monthly' ? 'mês' : 'ano'}</span>
					</div>
					<div class="plan-details">
						<div class="detail-item">
							<strong>Slug:</strong>
							<span>{plan.slug}</span>
						</div>
						<div class="detail-item">
							<strong>Max Usuários:</strong>
							<span>{plan.max_users}</span>
						</div>
						<div class="detail-item">
							<strong>Max Produtos:</strong>
							<span>{plan.max_products}</span>
						</div>
						{#if plan.description}
							<p class="description">{plan.description}</p>
						{/if}
					</div>
					<div class="plan-actions">
						<Button variant="secondary" size="sm" on:click={() => openEditModal(plan)}>Editar</Button>
						<Button variant="secondary" size="sm" on:click={() => togglePlanStatus(plan)}>
							{plan.active ? 'Desativar' : 'Ativar'}
						</Button>
						<Button variant="secondary" size="sm" on:click={() => deletePlan(plan)}>Deletar</Button>
					</div>
				</Card>
			{/each}
		</div>
	{/if}

	{#if showCreateModal}
		<div class="modal-overlay" on:click={() => showCreateModal = false}>
			<Card class="modal-card">
				<h2>Novo Plano</h2>
				<div class="form-group">
					<label>Nome *</label>
					<Input bind:value={createForm.name} placeholder="Nome do plano" />
				</div>
				<div class="form-group">
					<label>Slug *</label>
					<Input bind:value={createForm.slug} placeholder="slug-do-plano" />
				</div>
				<div class="form-group">
					<label>Descrição</label>
					<Textarea bind:value={createForm.description} placeholder="Descrição do plano" />
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Preço</label>
						<Input type="number" bind:value={createForm.price} placeholder="0.00" />
					</div>
					<div class="form-group">
						<label>Moeda</label>
						<Input bind:value={createForm.currency} placeholder="BRL" />
					</div>
				</div>
				<div class="form-group">
					<label>Intervalo</label>
					<Select bind:value={createForm.interval}>
						{#each intervals as interval}
							<option value={interval.value}>{interval.label}</option>
						{/each}
					</Select>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Max Usuários</label>
						<Input type="number" bind:value={createForm.max_users} />
					</div>
					<div class="form-group">
						<label>Max Produtos</label>
						<Input type="number" bind:value={createForm.max_products} />
					</div>
				</div>
				<div class="form-group">
					<label>Features (JSON)</label>
					<Textarea bind:value={createForm.features} placeholder='feature1: true, feature2: false' />
				</div>
				<div class="modal-actions">
					<Button variant="secondary" on:click={() => showCreateModal = false}>Cancelar</Button>
					<Button variant="primary" on:click={createPlan}>Criar</Button>
				</div>
			</Card>
		</div>
	{/if}

	{#if showEditModal && editingPlan}
		<div class="modal-overlay" on:click={() => showEditModal = false}>
			<Card class="modal-card">
				<h2>Editar Plano</h2>
				<div class="form-group">
					<label>Nome *</label>
					<Input bind:value={editingPlan.name} placeholder="Nome do plano" />
				</div>
				<div class="form-group">
					<label>Descrição</label>
					<Textarea bind:value={editingPlan.description} placeholder="Descrição do plano" />
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Preço</label>
						<Input type="number" bind:value={editingPlan.price} placeholder="0.00" />
					</div>
					<div class="form-group">
						<label>Moeda</label>
						<Input bind:value={editingPlan.currency} placeholder="BRL" />
					</div>
				</div>
				<div class="form-group">
					<label>Intervalo</label>
					<Select bind:value={editingPlan.interval}>
						{#each intervals as interval}
							<option value={interval.value}>{interval.label}</option>
						{/each}
					</Select>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Max Usuários</label>
						<Input type="number" bind:value={editingPlan.max_users} />
					</div>
					<div class="form-group">
						<label>Max Produtos</label>
						<Input type="number" bind:value={editingPlan.max_products} />
					</div>
				</div>
				<div class="form-group">
					<label>Features (JSON)</label>
					<Textarea bind:value={editingPlan.features} placeholder='feature1: true, feature2: false' />
				</div>
				<div class="form-group">
					<label>Status</label>
					<div class="radio-group">
						<label>
							<input type="radio" bind:group={editingPlan.active} value={true} />
							Ativo
						</label>
						<label>
							<input type="radio" bind:group={editingPlan.active} value={false} />
							Inativo
						</label>
					</div>
				</div>
				<div class="modal-actions">
					<Button variant="secondary" on:click={() => showEditModal = false}>Cancelar</Button>
					<Button variant="primary" on:click={updatePlan}>Salvar</Button>
				</div>
			</Card>
		</div>
	{/if}
</div>

<style>
	.plans-container {
		padding: 2rem;
		max-width: 1400px;
		margin: 0 auto;
	}

	.plans-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	.plans-header h1 {
		font-size: 1.5rem;
		font-weight: 600;
		margin: 0;
	}

	.loading-state {
		text-align: center;
		padding: 2rem;
		color: #64748b;
	}

	.plans-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.plan-card {
		background: rgba(255, 255, 255, 0.1);
		padding: 1.5rem;
	}

	.plan-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	.plan-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
	}

	.plan-price {
		display: flex;
		align-items: baseline;
		gap: 0.25rem;
		margin-bottom: 1rem;
	}

	.plan-price .price {
		font-size: 2rem;
		font-weight: 700;
	}

	.plan-price .interval {
		color: #64748b;
	}

	.plan-details {
		margin-bottom: 1.5rem;
	}

	.detail-item {
		display: flex;
		justify-content: space-between;
		padding: 0.5rem 0;
		border-bottom: 1px solid rgba(0, 0, 0, 0.1);
	}

	.detail-item:last-child {
		border-bottom: none;
	}

	.detail-item strong {
		font-weight: 500;
	}

	.description {
		margin-top: 1rem;
		color: #64748b;
		line-height: 1.5;
	}

	.plan-actions {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.plan-actions button {
		flex: 1;
		min-width: 80px;
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal-card {
		background: #1e293b;
		padding: 2rem;
		min-width: 500px;
		max-width: 600px;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-card h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin-bottom: 1.5rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 500;
		color: #e2e8f0;
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.modal-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1.5rem;
	}

	.modal-actions button {
		flex: 1;
	}

	.radio-group {
		display: flex;
		gap: 1rem;
	}

	.radio-group label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
	}

	.radio-group input {
		cursor: pointer;
	}
</style>
