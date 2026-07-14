<script lang="ts">
	import { onMount } from 'svelte';
	import { getPendingAdjustments, approveAdjustment, rejectAdjustment } from '$lib/api/stock-adjustment';
	import type { StockAdjustment } from '$lib/types/stock-adjustment';
	import { STOCK_ADJUSTMENT_STATUS_LABEL, STOCK_ADJUSTMENT_STATUS_COLOR } from '$lib/types/stock-adjustment';

	let adjustments: StockAdjustment[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let filter = $state<string>('all');

	onMount(loadAdjustments);

	async function loadAdjustments() {
		loading = true;
		error = '';
		try {
			adjustments = await getPendingAdjustments();
		} catch (e: any) {
			error = e?.message ?? 'Erro ao carregar ajustes.';
		} finally {
			loading = false;
		}
	}

	const filtered = $derived(
		filter === 'all' ? adjustments : adjustments.filter((a) => a.status === filter)
	);

	function formatDate(d?: string) {
		if (!d) return '—';
		return new Intl.DateTimeFormat('pt-BR', { dateStyle: 'short', timeStyle: 'short' }).format(new Date(d));
	}

	const STATUS_OPTIONS = [
		{ value: 'all', label: 'Todos' },
		{ value: 'pending', label: 'Pendentes' },
		{ value: 'approved', label: 'Aprovados' },
		{ value: 'rejected', label: 'Rejeitados' },
	];

	// Contagem por status para os pills
	const countByStatus = $derived(
		adjustments.reduce<Record<string, number>>((acc, a) => {
			acc[a.status] = (acc[a.status] ?? 0) + 1;
			return acc;
		}, {})
	);

	// Modal de ação
	let showActionModal = $state(false);
	let actionModalType: 'approve' | 'reject' | null = $state(null);
	let selectedAdjustment: StockAdjustment | null = $state(null);
	let actionNotes = $state('');
	let actionLoading = $state(false);

	function openApproveModal(adjustment: StockAdjustment) {
		selectedAdjustment = adjustment;
		actionModalType = 'approve';
		actionNotes = '';
		showActionModal = true;
	}

	function openRejectModal(adjustment: StockAdjustment) {
		selectedAdjustment = adjustment;
		actionModalType = 'reject';
		actionNotes = '';
		showActionModal = true;
	}

	function closeModal() {
		showActionModal = false;
		actionModalType = null;
		selectedAdjustment = null;
		actionNotes = '';
	}

	async function executeAction() {
		if (!selectedAdjustment || !actionModalType) return;

		actionLoading = true;
		error = '';

		try {
			if (actionModalType === 'approve') {
				await approveAdjustment(selectedAdjustment.id, { notes: actionNotes });
			} else {
				await rejectAdjustment(selectedAdjustment.id, { notes: actionNotes });
			}

			await loadAdjustments();
			closeModal();
		} catch (e: any) {
			error = e?.message ?? 'Erro ao executar ação.';
		} finally {
			actionLoading = false;
		}
	}</script>

<svelte:head>
	<title>Ajustes de Estoque - PratoOnline</title>
</svelte:head>

<div class="page-wrapper">
	<header class="page-header">
		<div>
			<h1 class="page-title">Ajustes de Estoque</h1>
			<p class="page-subtitle">{adjustments.length} ajuste{adjustments.length !== 1 ? 's' : ''} registrado{adjustments.length !== 1 ? 's' : ''}</p>
		</div>
	</header>

	<!-- Filtros -->
	<div class="filter-row">
		{#each STATUS_OPTIONS as opt}
			<button
				class="filter-pill"
				class:active={filter === opt.value}
				onclick={() => (filter = opt.value)}
			>
				{opt.label}
				{#if opt.value !== 'all' && countByStatus[opt.value]}
					<span class="pill-count">{countByStatus[opt.value]}</span>
				{/if}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="loading-state">
			<div class="spinner"></div>
			<span>Carregando ajustes…</span>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<span>⚠️ {error}</span>
			<button class="btn btn-sm" onclick={loadAdjustments}>Tentar novamente</button>
		</div>
	{:else if filtered.length === 0}
		<div class="empty-state">
			<p class="empty-icon">📦</p>
			<p class="empty-text">{filter === 'all' ? 'Nenhum ajuste ainda.' : 'Nenhum ajuste com esse status.'}</p>
		</div>
	{:else}
		<div class="adjustments-list">
			{#each filtered as adj}
				<div class="adjustment-card">
					<div class="adjustment-card-left">
						<span class="adjustment-id"># {adj.id}</span>
						<span class="adjustment-date">{formatDate(adj.created_at)}</span>
					</div>
					<div class="adjustment-card-center">
						<div class="adjustment-info">
							<span class="info-label">Pedido:</span>
							<span class="info-value">{adj.order_id}</span>
						</div>
						<div class="adjustment-info">
							<span class="info-label">Ingrediente:</span>
							<span class="info-value">{adj.ingredient_id}</span>
						</div>
						<div class="adjustment-info">
							<span class="info-label">Quantidade:</span>
							<span class="info-value">{adj.quantity?.toFixed?.(4) ?? '0.0000'}</span>
						</div>
						<div class="adjustment-info">
							<span class="info-label">Status Pedido:</span>
							<span class="info-value">{adj.order_status}</span>
						</div>
					</div>
					<div class="adjustment-card-right">
						<span class="status-badge {STOCK_ADJUSTMENT_STATUS_COLOR[adj.status]}">{STOCK_ADJUSTMENT_STATUS_LABEL[adj.status]}</span>
						{#if adj.status === 'pending'}
							<div class="card-actions">
								<button class="btn btn-sm btn-ghost" onclick={() => openApproveModal(adj)}>Aprovar</button>
								<button class="btn btn-sm btn-ghost" onclick={() => openRejectModal(adj)}>Rejeitar</button>
							</div>
						{:else}
							<span class="muted text-xs">
								{adj.processed_by ? `Processado por ID ${adj.processed_by}` : 'Processado'}
							</span>
						{/if}
					</div>
				</div>
				{#if adj.processed_by || adj.processing_notes}
					<div class="adjustment-details">
						{#if adj.processed_at}
							<strong>Processado em:</strong> {formatDate(adj.processed_at)}
							{#if adj.processed_by} - <strong>Por:</strong> ID {adj.processed_by}{/if}
						{/if}
						{#if adj.processing_notes}
							<br />
							<strong>Observações:</strong> {adj.processing_notes}
						{/if}
					</div>
				{/if}
			{/each}
		</div>
	{/if}
</div>

<!-- Modal de Ação -->
{#if showActionModal}
	<div class="modal-overlay" onclick={() => (showActionModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">
				{actionModalType === 'approve' ? 'Aprovar Ajuste' : 'Rejeitar Ajuste'}
			</h2>

			<div class="modal-info">
				<p class="info-item"><strong>ID:</strong> {selectedAdjustment?.id}</p>
				<p class="info-item"><strong>Pedido:</strong> {selectedAdjustment?.order_id}</p>
				<p class="info-item"><strong>Ingrediente:</strong> {selectedAdjustment?.ingredient_id}</p>
				<p class="info-item"><strong>Quantidade:</strong> {selectedAdjustment?.quantity?.toFixed?.(4) ?? '0.0000'}</p>
			</div>

			<div class="form-group">
				<label for="notes">Observações (opcional)</label>
				<textarea
					id="notes"
					bind:value={actionNotes}
					placeholder="Adicione observações sobre esta decisão..."
					rows="3"
				></textarea>
			</div>

			<div class="modal-actions">
				<button class="btn btn-ghost" onclick={closeModal} disabled={actionLoading}>Cancelar</button>
				<button
					class="btn {actionModalType === 'approve' ? 'btn-primary' : 'btn-danger'}"
					onclick={executeAction}
					disabled={actionLoading}
				>
					{actionLoading ? 'Processando...' : (actionModalType === 'approve' ? 'Aprovar' : 'Rejeitar')}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
  .page-wrapper   { max-width: 1000px; margin: 0 auto; padding: 2rem 1.5rem; }
  .page-header    { display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 1rem; margin-bottom: 1.75rem; }
  .page-title     { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .page-subtitle  { font-size: 0.875rem; color: var(--color-muted); margin: 0.25rem 0 0; }

  /* Filtros */
  .filter-row     { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .filter-pill    { border: 1px solid var(--color-border, #e5e7eb); background: var(--color-surface, #fff); color: var(--color-muted); font-size: 0.8rem; font-weight: 500; padding: 0.35rem 0.75rem; border-radius: 99px; cursor: pointer; transition: all 0.15s; display: flex; align-items: center; gap: 0.35rem; }
  .filter-pill:hover { border-color: var(--color-primary, #e85d04); color: var(--color-primary, #e85d04); }
  .filter-pill.active { background: var(--color-primary, #e85d04); border-color: var(--color-primary, #e85d04); color: #fff; }
  .pill-count     { background: rgba(255,255,255,0.3); padding: 0 0.3rem; border-radius: 99px; font-size: 0.7rem; font-weight: 700; }
  .filter-pill:not(.active) .pill-count { background: var(--color-surface-2, #f0f0f0); color: var(--color-text); }

  /* Lista de ajustes */
  .adjustments-list { display: flex; flex-direction: column; gap: 0.6rem; }
  .adjustment-card { display: flex; align-items: center; gap: 1rem; background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1rem 1.25rem; transition: box-shadow 0.15s, border-color 0.15s; }
  .adjustment-card:hover { box-shadow: 0 4px 14px rgba(0,0,0,0.07); border-color: var(--color-primary, #e85d04); }
  .adjustment-card-left { display: flex; flex-direction: column; gap: 0.2rem; min-width: 90px; }
  .adjustment-id { font-weight: 700; color: var(--color-text); font-size: 0.95rem; }
  .adjustment-date { font-size: 0.75rem; color: var(--color-muted); }
  .adjustment-card-center { flex: 1; display: flex; gap: 1rem; flex-wrap: wrap; }
  .adjustment-info { display: flex; gap: 0.5rem; font-size: 0.85rem; }
  .info-label { color: var(--color-muted); }
  .info-value { font-weight: 500; color: var(--color-text); }
  .adjustment-card-right { display: flex; flex-direction: column; align-items: flex-end; gap: 0.4rem; min-width: 110px; }
  .card-actions { display: flex; gap: 0.5rem; }
  .adjustment-details { padding: 0.5rem 1.25rem 1rem; background: var(--color-surface-2, #f9fafb); border-radius: 0 0 0.75rem 0.75rem; font-size: 0.85rem; color: var(--color-muted); margin-top: -0.6rem; margin-bottom: 0.6rem; border: 1px solid var(--color-border, #e5e7eb); border-top: none; }

  /* Status badges */
  .status-badge   { font-size: 0.72rem; font-weight: 600; padding: 0.2rem 0.55rem; border-radius: 99px; white-space: nowrap; }
  .badge-warning  { background: #fef9c3; color: #854d0e; }
  .badge-info     { background: #dbeafe; color: #1e40af; }
  .badge-success  { background: #dcfce7; color: #15803d; }
  .badge-neutral  { background: var(--color-surface-2, #f3f4f6); color: var(--color-muted); }
  .badge-error    { background: #fee2e2; color: #b91c1c; }

  /* Estados */
  .loading-state  { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: var(--color-muted); }
  .spinner        { width: 2rem; height: 2rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin  { to { transform: rotate(360deg); } }
  .empty-state    { text-align: center; padding: 4rem 1rem; }
  .empty-icon     { font-size: 2.5rem; margin: 0; }
  .empty-text     { color: var(--color-muted); margin: 0.5rem 0 1.25rem; }
  .alert-error    { display: flex; align-items: center; gap: 1rem; background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 1rem 1.25rem; border-radius: 0.5rem; margin-bottom: 1.5rem; }

  /* Modal */
  .modal-overlay  { position: fixed; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; z-index: 50; padding: 1rem; }
  .modal          { background: var(--color-surface, #fff); border-radius: 1rem; padding: 2rem; width: 100%; max-width: 480px; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }
  .modal-title    { font-size: 1.25rem; font-weight: 700; margin: 0 0 1.5rem; }
  .modal-info     { margin-bottom: 1.5rem; }
  .info-item      { font-size: 0.9rem; color: var(--color-muted); margin: 0.5rem 0; }
  .modal-actions  { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.75rem; }

  /* Formulário */
  .form-group     { display: flex; flex-direction: column; gap: 0.4rem; margin-bottom: 1rem; }
  .form-group label { font-size: 0.85rem; font-weight: 500; color: var(--color-text); }
  .form-group textarea { border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; padding: 0.6rem 0.75rem; font-size: 0.9rem; background: var(--color-surface, #fff); color: var(--color-text); transition: border-color 0.15s; width: 100%; box-sizing: border-box; }
  .form-group textarea:focus { outline: none; border-color: var(--color-primary, #e85d04); box-shadow: 0 0 0 3px rgba(232,93,4,0.12); }

  /* Botões */
  .btn            { display: inline-flex; align-items: center; padding: 0.55rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; text-decoration: none; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-secondary  { background: var(--color-surface-2, #f3f4f6); color: var(--color-text); border: 1px solid var(--color-border, #d1d5db); }
  .btn-secondary:hover:not(:disabled) { background: var(--color-border, #e5e7eb); }
  .btn-ghost      { background: transparent; color: var(--color-muted); border: 1px solid var(--color-border, #e5e7eb); }
  .btn-ghost:hover:not(:disabled) { background: var(--color-surface-2, #f3f4f6); }
  .btn-danger     { background: #fee2e2; color: #b91c1c; border: 1px solid #fca5a5; }
  .btn-danger:hover:not(:disabled) { background: #fecaca; }
  .btn-sm         { padding: 0.35rem 0.75rem; font-size: 0.8rem; }

  .muted          { color: var(--color-muted); }

  @media (max-width: 640px) {
    .adjustment-card         { flex-wrap: wrap; }
    .adjustment-card-center { width: 100%; }
    .adjustment-card-right   { flex-direction: row; min-width: unset; width: 100%; justify-content: space-between; align-items: center; }
  }
</style>
