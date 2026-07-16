<script lang="ts">
	import { onMount } from 'svelte';
	import { getPendingAdjustments, approveAdjustment, rejectAdjustment } from '$lib/api/stock-adjustment';
	import type { StockAdjustment } from '$lib/types/stock-adjustment';
	import { STOCK_ADJUSTMENT_STATUS_LABEL, STOCK_ADJUSTMENT_STATUS_COLOR } from '$lib/types/stock-adjustment';
	import { Button, Badge, Alert, Card, Modal, Textarea, Skeleton } from '$lib/components/ui';
	import { Workspace } from '$lib/components/layout';
	import { Check, X, Clock, Package, Filter, ArrowUpDown, Activity, Calendar, User } from '@lucide/svelte';

	let adjustments: StockAdjustment[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let filterStatus = $state<string>('all');
	let sortBy = $state<string>('date');
	let sortOrder = $state<'asc' | 'desc'>('desc');
	let currentPage = $state(1);
	const itemsPerPage = 12;

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
		adjustments.filter((a) => {
			if (filterStatus !== 'all' && a.status !== filterStatus) return false;
			return true;
		}).sort((a, b) => {
			let comparison = 0;
			if (sortBy === 'date') {
				comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
			} else if (sortBy === 'quantity') {
				comparison = a.quantity - b.quantity;
			} else if (sortBy === 'id') {
				comparison = a.id - b.id;
			}
			return sortOrder === 'asc' ? comparison : -comparison;
		})
	);

	const totalPages = $derived(Math.ceil(filtered.length / itemsPerPage));
	const paginated = $derived(
		filtered.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
	);

	function goToPage(page: number) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
		}
	}

	function formatDate(d?: string) {
		if (!d) return '—';
		return new Intl.DateTimeFormat('pt-BR', { dateStyle: 'short', timeStyle: 'short' }).format(new Date(d));
	}

	function getStatusVariant(status: string) {
		switch (status) {
			case 'pending': return 'warning';
			case 'approved': return 'success';
			case 'rejected': return 'danger';
			default: return 'default';
		}
	}

	const STATUS_OPTIONS = [
		{ value: 'all', label: 'Todos' },
		{ value: 'pending', label: 'Pendentes' },
		{ value: 'approved', label: 'Aprovados' },
		{ value: 'rejected', label: 'Rejeitados' },
	];

	const countByStatus = $derived(
		adjustments.reduce<Record<string, number>>((acc, a) => {
			acc[a.status] = (acc[a.status] ?? 0) + 1;
			return acc;
		}, {})
	);

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
	}
</script>

<Workspace
  breadcrumb={[{ label: 'Ajustes de Estoque' }]}
  title="Ajustes de Estoque"
  description="Gerencie as solicitações de ajuste de ingredientes"
>
  {#if loading}
    <div class="skeleton-grid">
      {#each Array(6) as _}
        <Card class="skeleton-card">
          <div class="skeleton-card-header">
            <Skeleton variant="text" width="80px" height="16px" />
            <Skeleton variant="rectangular" width="60px" height="20px" />
          </div>
          <div class="skeleton-card-meta">
            <Skeleton variant="text" width="120px" height="12px" />
            <Skeleton variant="text" width="100px" height="12px" />
          </div>
          <div class="skeleton-card-details">
            <Skeleton variant="text" width="100%" height="12px" />
            <Skeleton variant="text" width="80%" height="12px" />
          </div>
        </Card>
      {/each}
    </div>
  {:else if error}
    <Alert variant="error" dismissible onDismiss={() => error = ''}>
      ⚠️ {error}
      <Button onclick={loadAdjustments} size="sm">Tentar novamente</Button>
    </Alert>
  {:else}
    <!-- Filtros -->
    <Card class="filters-card">
      <div class="filters-header">
        <div class="filters-title">
          <Filter size={18} />
          <span>Filtros</span>
        </div>
        <div class="filters-stats">
          <span class="stat-item">{filtered.length} ajustes</span>
        </div>
      </div>
      
      <!-- Status Pills -->
      <div class="status-pills">
        {#each STATUS_OPTIONS as opt}
          <button
            class="status-pill"
            class:active={filterStatus === opt.value}
            onclick={() => (filterStatus = opt.value)}
          >
            {opt.label}
            {#if opt.value !== 'all' && countByStatus[opt.value]}
              <span class="pill-count">{countByStatus[opt.value]}</span>
            {/if}
          </button>
        {/each}
      </div>

      <!-- Sort -->
      <div class="filters-row">
        <div class="sort-wrapper">
          <select bind:value={sortBy} class="sort-select">
            <option value="date">Data</option>
            <option value="quantity">Quantidade</option>
            <option value="id">ID</option>
          </select>
          <Button
            variant="ghost"
            size="sm"
            onclick={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
            class="sort-toggle"
          >
            <ArrowUpDown size={16} />
          </Button>
        </div>
      </div>
    </Card>

    {#if filtered.length === 0}
      <div class="empty-state">
        <Package size={48} class="empty-icon" />
        <span class="empty-title">{filterStatus === 'all' ? 'Nenhum ajuste encontrado' : 'Nenhum ajuste com esse status'}</span>
        <span class="empty-subtitle">Não há solicitações de ajuste de estoque</span>
      </div>
    {:else}
      <div class="adjustments-timeline">
        {#each paginated as adj}
          <Card class={`adjustment-card ${adj.status === 'pending' ? 'pending' : ''} ${adj.status === 'approved' ? 'approved' : ''} ${adj.status === 'rejected' ? 'rejected' : ''}`}>
            <div class="adjustment-header">
              <div class="adjustment-id">#{adj.id}</div>
              <Badge variant={getStatusVariant(adj.status)} size="sm">
                {STOCK_ADJUSTMENT_STATUS_LABEL[adj.status]}
              </Badge>
            </div>

            <div class="adjustment-meta">
              <div class="meta-item">
                <Calendar size={14} />
                <span>{formatDate(adj.created_at)}</span>
              </div>
              {#if adj.processed_at}
                <div class="meta-item">
                  <Clock size={14} />
                  <span>Processado em {formatDate(adj.processed_at)}</span>
                </div>
              {/if}
            </div>

            <div class="adjustment-details">
              <div class="detail-row">
                <span class="detail-label">Pedido</span>
                <span class="detail-value">#{adj.order_id}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Ingrediente</span>
                <span class="detail-value">#{adj.ingredient_id}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Quantidade</span>
                <span class="detail-value quantity">{adj.quantity?.toFixed?.(4) ?? '0.0000'}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Status Pedido</span>
                <span class="detail-value">{adj.order_status}</span>
              </div>
            </div>

            {#if adj.status === 'pending'}
              <div class="adjustment-actions">
                <Button onclick={() => openApproveModal(adj)} variant="success" size="sm">
                  <Check size={14} />
                  Aprovar
                </Button>
                <Button onclick={() => openRejectModal(adj)} variant="danger" size="sm">
                  <X size={14} />
                  Rejeitar
                </Button>
              </div>
            {:else}
              <div class="adjustment-processed">
                {#if adj.processed_by}
                  <div class="processed-info">
                    <User size={14} />
                    <span>Processado por ID {adj.processed_by}</span>
                  </div>
                {/if}
                {#if adj.processing_notes}
                  <div class="processed-notes">
                    <span>{adj.processing_notes}</span>
                  </div>
                {/if}
              </div>
            {/if}
          </Card>
        {/each}
      </div>

      {#if totalPages > 1}
        <div class="pagination">
          <Button
            variant="ghost"
            size="sm"
            disabled={currentPage === 1}
            onclick={() => goToPage(currentPage - 1)}
          >
            Anterior
          </Button>
          <span class="pagination-info">
            Página {currentPage} de {totalPages}
          </span>
          <Button
            variant="ghost"
            size="sm"
            disabled={currentPage === totalPages}
            onclick={() => goToPage(currentPage + 1)}
          >
            Próxima
          </Button>
        </div>
      {/if}
    {/if}
  {/if}
</Workspace>

<!-- Modal de Ação -->
<Modal 
  open={showActionModal} 
  title={actionModalType === 'approve' ? 'Aprovar Ajuste' : 'Rejeitar Ajuste'}
  onClose={closeModal}
>
  <div class="modal-info">
    <div class="info-row">
      <span class="info-label">ID:</span>
      <span class="info-value">{selectedAdjustment?.id}</span>
    </div>
    <div class="info-row">
      <span class="info-label">Pedido:</span>
      <span class="info-value">#{selectedAdjustment?.order_id}</span>
    </div>
    <div class="info-row">
      <span class="info-label">Ingrediente:</span>
      <span class="info-value">#{selectedAdjustment?.ingredient_id}</span>
    </div>
    <div class="info-row">
      <span class="info-label">Quantidade:</span>
      <span class="info-value">{selectedAdjustment?.quantity?.toFixed?.(4) ?? '0.0000'}</span>
    </div>
  </div>

  <Textarea
    id="notes"
    label="Observações (opcional)"
    bind:value={actionNotes}
    placeholder="Adicione observações sobre esta decisão..."
    rows={3}
  />

  <div class="modal-actions">
    <Button variant="ghost" onclick={closeModal} disabled={actionLoading}>Cancelar</Button>
    <Button
      variant={actionModalType === 'approve' ? 'primary' : 'danger'}
      onclick={executeAction}
      disabled={actionLoading}
      loading={actionLoading}
    >
      {actionModalType === 'approve' ? 'Aprovar' : 'Rejeitar'}
    </Button>
  </div>
</Modal>

<style>
  /* Loading State */
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    padding: 4rem;
    color: #64748b;
    font-size: 0.875rem;
  }

  .spinner {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Filters Card */
  .filters-card {
    margin-bottom: 2rem;
  }

  .filters-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .filters-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: #0f172a;
  }

  .filters-stats {
    display: flex;
    gap: 1rem;
  }

  .stat-item {
    font-size: 0.75rem;
    color: #64748b;
    padding: 0.25rem 0.5rem;
    background: #f8fafc;
    border-radius: 4px;
  }

  /* Status Pills */
  .status-pills {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .status-pill {
    border: 1px solid #f1f5f9;
    background: #ffffff;
    color: #64748b;
    font-size: 0.875rem;
    font-weight: 500;
    padding: 0.5rem 1rem;
    border-radius: 999px;
    cursor: pointer;
    transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .status-pill:hover {
    border-color: #6366f1;
    color: #6366f1;
    transform: translateY(-1px);
  }

  .status-pill.active {
    background: #6366f1;
    border-color: #6366f1;
    color: #ffffff;
  }

  .pill-count {
    background: rgba(255, 255, 255, 0.2);
    padding: 0 0.375rem;
    border-radius: 999px;
    font-size: 0.75rem;
    font-weight: 600;
    min-width: 20px;
    text-align: center;
  }

  .status-pill:not(.active) .pill-count {
    background: #f1f5f9;
    color: #0f172a;
  }

  /* Filters Row */
  .filters-row {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .sort-wrapper {
    display: flex;
    gap: 0.5rem;
  }

  .sort-select {
    padding: 0.5rem 0.75rem;
    border: 1px solid #f1f5f9;
    border-radius: 8px;
    font-size: 0.875rem;
    background: #ffffff;
    color: #0f172a;
    cursor: pointer;
  }

  .sort-toggle {
    padding: 0.5rem;
  }

  /* Adjustments Timeline */
  .adjustments-timeline {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .adjustment-card {
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .adjustment-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px 0 rgb(0 0 0 / 0.08);
  }

  .adjustment-card.pending {
    border-left: 4px solid #f59e0b;
  }

  .adjustment-card.approved {
    border-left: 4px solid #10b981;
  }

  .adjustment-card.rejected {
    border-left: 4px solid #ef4444;
  }

  .adjustment-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .adjustment-id {
    font-size: 1rem;
    font-weight: 700;
    color: #0f172a;
  }

  .adjustment-meta {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
    font-size: 0.75rem;
    color: #64748b;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .adjustment-details {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 0.75rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 8px;
    margin-bottom: 1rem;
  }

  .detail-row {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .detail-label {
    font-size: 0.75rem;
    color: #64748b;
  }

  .detail-value {
    font-size: 0.875rem;
    font-weight: 500;
    color: #0f172a;
  }

  .detail-value.quantity {
    font-weight: 700;
    color: #6366f1;
  }

  .adjustment-actions {
    display: flex;
    gap: 0.5rem;
  }

  .adjustment-processed {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .processed-info {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.75rem;
    color: #64748b;
  }

  .processed-notes {
    font-size: 0.875rem;
    color: #64748b;
    padding: 0.75rem;
    background: #f8fafc;
    border-radius: 8px;
  }

  /* Empty State */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    padding: 4rem;
    text-align: center;
  }

  .empty-icon {
    color: #cbd5e1;
  }

  .empty-title {
    font-size: 1rem;
    font-weight: 600;
    color: #0f172a;
  }

  .empty-subtitle {
    font-size: 0.875rem;
    color: #64748b;
  }

  /* Pagination */
  .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 2rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 12px;
    border: 1px solid #f1f5f9;
  }

  .pagination-info {
    font-size: 0.875rem;
    color: #64748b;
  }

  /* Modal */
  .modal-info {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    font-size: 0.875rem;
  }

  .info-label {
    color: #64748b;
  }

  .info-value {
    font-weight: 500;
    color: #0f172a;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1.75rem;
  }

  /* Skeleton Loading */
  .skeleton-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
  }

  .skeleton-card {
    padding: 1.5rem;
  }

  .skeleton-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .skeleton-card-meta {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .skeleton-card-details {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .adjustment-details {
      grid-template-columns: 1fr;
    }

    .status-pills {
      overflow-x: auto;
      padding-bottom: 0.5rem;
    }

    .status-pill {
      flex-shrink: 0;
    }

    .adjustment-actions {
      flex-direction: column;
    }

    .adjustment-actions button {
      width: 100%;
    }
  }
</style>
