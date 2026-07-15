<script lang="ts">
  import { onMount } from 'svelte';
  import { getOrders } from '$lib/api/order';
  import type { Order } from '$lib/types/order';
  import { ORDER_STATUS_LABEL, ORDER_STATUS_COLOR } from '$lib/types/order';

  let orders: Order[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let filterStatus = $state<string>('all');
  let searchQuery = $state<string>('');
  let dateFrom = $state<string>('');
  let dateTo = $state<string>('');
  let sortBy = $state<string>('date');
  let sortOrder = $state<'asc' | 'desc'>('desc');
  let currentPage = $state(1);
  const itemsPerPage = 20;

  onMount(loadOrders);

  async function loadOrders() {
    loading = true;
    error = '';
    try {
      orders = await getOrders();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar pedidos.';
    } finally {
      loading = false;
    }
  }

  const filtered = $derived(
    orders.filter((order) => {
      // Filtro por status
      if (filterStatus !== 'all' && order.Status !== filterStatus) return false;

      // Filtro por busca (ID ou nome de produtos)
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        const idMatch = order.ID.toString().includes(query);
        const productsMatch = order.Items.some(item =>
          item.ProductName.toLowerCase().includes(query)
        );
        if (!idMatch && !productsMatch) return false;
      }

      // Filtro por data
      if (dateFrom) {
        const orderDate = new Date(order.CreatedAt).setHours(0, 0, 0, 0);
        const fromDate = new Date(dateFrom).setHours(0, 0, 0, 0);
        if (orderDate < fromDate) return false;
      }
      if (dateTo) {
        const orderDate = new Date(order.CreatedAt).setHours(23, 59, 59, 999);
        const toDate = new Date(dateTo).setHours(23, 59, 59, 999);
        if (orderDate > toDate) return false;
      }

      return true;
    }).sort((a, b) => {
      // Ordenação
      let comparison = 0;
      if (sortBy === 'date') {
        comparison = new Date(a.CreatedAt).getTime() - new Date(b.CreatedAt).getTime();
      } else if (sortBy === 'total') {
        comparison = a.TotalPrice - b.TotalPrice;
      } else if (sortBy === 'id') {
        comparison = a.ID - b.ID;
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

  function nextPage() {
    goToPage(currentPage + 1);
  }

  function prevPage() {
    goToPage(currentPage - 1);
  }

  function formatDate(d?: string) {
    if (!d) return '—';
    return new Intl.DateTimeFormat('pt-BR', { dateStyle: 'short', timeStyle: 'short' }).format(new Date(d));
  }

  function formatTotal(v: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(v);
  }

  const STATUS_OPTIONS = [
    { value: 'all', label: 'Todos' },
    { value: 'pending', label: 'Pendentes' },
    { value: 'confirmed', label: 'Confirmados' },
    { value: 'preparing', label: 'Preparando' },
    { value: 'ready', label: 'Prontos' },
    { value: 'delivered', label: 'Entregues' },
    { value: 'cancelled', label: 'Cancelados' },
  ];

  // Contagem por status para os pills
  const countByStatus = $derived(
    orders.reduce<Record<string, number>>((acc, o) => {
      acc[o.Status] = (acc[o.Status] ?? 0) + 1;
      return acc;
    }, {})
  );
</script>

<div class="page-wrapper">
  <header class="page-header">
    <div>
      <h1 class="page-title">Pedidos</h1>
      <p class="page-subtitle">{orders.length} pedido{orders.length !== 1 ? 's' : ''} registrado{orders.length !== 1 ? 's' : ''}</p>
    </div>
    <a href="/orders/new" class="btn btn-primary">+ Novo Pedido</a>
  </header>

  <!-- Filtros -->
  <div class="filter-section" role="region" aria-label="Filtros de pedidos">
    <div class="filter-row" role="group" aria-label="Filtro por status">
      {#each STATUS_OPTIONS as opt}
        <button
          class="filter-pill"
          class:active={filterStatus === opt.value}
          onclick={() => (filterStatus = opt.value)}
          aria-pressed={filterStatus === opt.value}
          aria-label={`Filtrar por ${opt.label}`}
        >
          {opt.label}
          {#if opt.value !== 'all' && countByStatus[opt.value]}
            <span class="pill-count" aria-label={`${countByStatus[opt.value]} pedidos`}>{countByStatus[opt.value]}</span>
          {/if}
        </button>
      {/each}
    </div>

    <div class="filter-controls" role="search" aria-label="Busca e filtros adicionais">
      <div class="filter-control">
        <label for="order-search" class="sr-only">Buscar pedidos</label>
        <input
          id="order-search"
          type="text"
          placeholder="Buscar por ID ou produto..."
          bind:value={searchQuery}
          class="search-input"
          aria-label="Buscar pedidos por ID ou nome do produto"
        />
      </div>
      <div class="filter-control">
        <label for="date-from" class="sr-only">Data inicial</label>
        <input
          id="date-from"
          type="date"
          placeholder="Data inicial"
          bind:value={dateFrom}
          class="date-input"
          aria-label="Data inicial do filtro"
        />
      </div>
      <div class="filter-control">
        <label for="date-to" class="sr-only">Data final</label>
        <input
          id="date-to"
          type="date"
          placeholder="Data final"
          bind:value={dateTo}
          class="date-input"
          aria-label="Data final do filtro"
        />
      </div>
      <div class="filter-control sort-control">
        <label for="sort-by" class="sr-only">Ordenar por</label>
        <select bind:value={sortBy} id="sort-by" class="sort-select" aria-label="Critério de ordenação">
          <option value="date">Ordenar por Data</option>
          <option value="total">Ordenar por Valor</option>
          <option value="id">Ordenar por ID</option>
        </select>
        <button
          class="btn btn-sm btn-ghost"
          onclick={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
          title={sortOrder === 'asc' ? 'Crescente' : 'Decrescente'}
          aria-label={`Ordem ${sortOrder === 'asc' ? 'crescente' : 'decrescente'}, clique para inverter`}
        >
          {sortOrder === 'asc' ? '↑' : '↓'}
        </button>
      </div>
      {(searchQuery || dateFrom || dateTo) && (
        <button class="btn btn-sm btn-ghost" onclick={() => { searchQuery = ''; dateFrom = ''; dateTo = ''; }} aria-label="Limpar todos os filtros">
          Limpar filtros
        </button>
      )}
    </div>
  </div>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <span>Carregando pedidos…</span>
    </div>
  {:else if error}
    <div class="alert alert-error">
      <span>⚠️ {error}</span>
      <button class="btn btn-sm" onclick={loadOrders}>Tentar novamente</button>
    </div>
  {:else if filtered.length === 0}
    <div class="empty-state">
      <p class="empty-icon">🧾</p>
      <p class="empty-text">{filterStatus === 'all' ? 'Nenhum pedido ainda.' : 'Nenhum pedido com esse status.'}</p>
      {#if filterStatus === 'all'}
        <a href="/orders/new" class="btn btn-primary">Criar primeiro pedido</a>
      {:else}
        <button class="btn btn-ghost" onclick={() => (filterStatus = 'all')}>Ver todos</button>
      {/if}
    </div>
  {:else}
    <div class="orders-list">
      {#each paginated as order}
        <a href="/orders/{order.ID}" class="order-card">
          <div class="order-card-left">
            <span class="order-id"># {order.ID}</span>
            <span class="order-date">{formatDate(order.CreatedAt)}</span>
          </div>
          <div class="order-items-preview">
            {#each (order.Items ?? []).slice(0, 3) as item}
              <span class="item-chip">{item.Quantity}× Produto #{item.ProductID}</span>
            {/each}
            {#if (order.Items ?? []).length > 3}
              <span class="item-chip muted">+{order.Items.length - 3} mais</span>
            {/if}
          </div>
          <div class="order-card-right">
            <span class="order-total">{formatTotal(order.TotalPrice)}</span>
            <span class="status-badge {ORDER_STATUS_COLOR[order.Status]}">{ORDER_STATUS_LABEL[order.Status]}</span>
          </div>
        </a>
      {/each}
    </div>

    {#if totalPages > 1}
      <div class="pagination">
        <button
          class="pagination-btn"
          disabled={currentPage === 1}
          onclick={prevPage}
        >
          ← Anterior
        </button>
        <span class="pagination-info">
          Página {currentPage} de {totalPages} ({filtered.length} pedidos)
        </span>
        <button
          class="pagination-btn"
          disabled={currentPage === totalPages}
          onclick={nextPage}
        >
          Próxima →
        </button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .page-wrapper   { max-width: 1000px; margin: 0 auto; padding: 2rem 1.5rem; }
  .page-header    { display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 1rem; margin-bottom: 1.75rem; }
  .page-title     { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .page-subtitle  { font-size: 0.875rem; color: var(--color-muted); margin: 0.25rem 0 0; }

  /* Acessibilidade */
  .sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0; }

  /* Filtros */
  .filter-section { margin-bottom: 1.5rem; }
  .filter-row     { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1rem; }
  .filter-pill    { border: 1px solid var(--color-border, #e5e7eb); background: var(--color-surface, #fff); color: var(--color-muted); font-size: 0.8rem; font-weight: 500; padding: 0.4rem 0.85rem; border-radius: 99px; cursor: pointer; transition: all 0.2s ease; display: flex; align-items: center; gap: 0.4rem; box-shadow: 0 1px 2px rgba(0,0,0,0.05); }
  .filter-pill:hover { border-color: var(--color-primary, #e85d04); color: var(--color-primary, #e85d04); box-shadow: 0 2px 4px rgba(232,93,4,0.15); transform: translateY(-1px); }
  .filter-pill.active { background: var(--color-primary, #e85d04); border-color: var(--color-primary, #e85d04); color: #fff; box-shadow: 0 2px 8px rgba(232,93,4,0.3); }
  .pill-count     { background: rgba(255,255,255,0.25); padding: 0 0.35rem; border-radius: 99px; font-size: 0.7rem; font-weight: 700; min-width: 20px; text-align: center; }
  .filter-pill:not(.active) .pill-count { background: var(--color-surface-2, #f0f0f0); color: var(--color-text); }

  .filter-controls { display: flex; gap: 0.75rem; flex-wrap: wrap; align-items: center; }
  .filter-control { flex: 1; min-width: 180px; }
  .sort-control { display: flex; gap: 0.25rem; min-width: auto; }
  .search-input,
  .date-input,
  .sort-select { width: 100%; padding: 0.5rem 0.75rem; border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; font-size: 0.85rem; background: var(--color-surface, #fff); color: var(--color-text); transition: border-color 0.15s; box-sizing: border-box; }
  .search-input:focus,
  .date-input:focus,
  .sort-select:focus { outline: none; border-color: var(--color-primary, #e85d04); }
  .sort-select { flex: 1; cursor: pointer; }

  /* Paginação */
  .pagination { display: flex; justify-content: center; align-items: center; gap: 1rem; margin-top: 2rem; padding: 1rem; background: var(--color-surface-2, #f9fafb); border-radius: 0.5rem; }
  .pagination-btn { padding: 0.5rem 1rem; border: 1px solid var(--color-border, #d1d5db); background: var(--color-surface, #fff); color: var(--color-text); border-radius: 0.4rem; font-size: 0.85rem; font-weight: 500; cursor: pointer; transition: all 0.15s; }
  .pagination-btn:hover:not(:disabled) { background: var(--color-primary, #e85d04); color: #fff; border-color: var(--color-primary, #e85d04); }
  .pagination-btn:disabled { opacity: 0.4; cursor: not-allowed; }
  .pagination-info { font-size: 0.85rem; color: var(--color-muted); }

  /* Lista pedidos */
  .orders-list    { display: flex; flex-direction: column; gap: 0.6rem; }
  .order-card     { display: flex; align-items: center; gap: 1rem; background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1rem 1.25rem; text-decoration: none; color: inherit; transition: box-shadow 0.15s, border-color 0.15s; }
  .order-card:hover { box-shadow: 0 4px 14px rgba(0,0,0,0.07); border-color: var(--color-primary, #e85d04); }
  .order-card-left { display: flex; flex-direction: column; gap: 0.2rem; min-width: 90px; }
  .order-id       { font-weight: 700; color: var(--color-text); font-size: 0.95rem; }
  .order-date     { font-size: 0.75rem; color: var(--color-muted); }
  .order-items-preview { flex: 1; display: flex; gap: 0.4rem; flex-wrap: wrap; }
  .item-chip      { font-size: 0.78rem; background: var(--color-surface-2, #f5f5f5); padding: 0.2rem 0.5rem; border-radius: 4px; color: var(--color-text); }
  .item-chip.muted { color: var(--color-muted); }
  .order-card-right { display: flex; flex-direction: column; align-items: flex-end; gap: 0.4rem; min-width: 110px; }
  .order-total    { font-weight: 700; font-size: 0.95rem; color: var(--color-text); }

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

  .btn            { display: inline-flex; align-items: center; padding: 0.55rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; text-decoration: none; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-ghost      { background: transparent; color: var(--color-muted); border: 1px solid var(--color-border, #e5e7eb); }
  .btn-ghost:hover { background: var(--color-surface-2, #f3f4f6); }
  .btn-sm         { padding: 0.35rem 0.75rem; font-size: 0.8rem; }

  @media (max-width: 640px) {
    .order-card         { flex-wrap: wrap; }
    .order-items-preview { width: 100%; }
    .order-card-right   { flex-direction: row; min-width: unset; width: 100%; justify-content: space-between; align-items: center; }
  }
</style>
