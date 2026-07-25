<script lang="ts">
  import { onMount } from 'svelte';
  import { getOrders } from '$lib/api/order';
  import type { Order } from '$lib/types/order';
  import { ORDER_STATUS_LABEL, ORDER_STATUS_COLOR } from '$lib/types/order';
  import { Button, Input, Badge, Alert, Card, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Search, Calendar, DollarSign, Clock, Filter, ArrowUpDown, Plus, ShoppingBag, Activity, MoreHorizontal } from '@lucide/svelte';

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
  const itemsPerPage = 12;

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
      if (filterStatus !== 'all' && order.Status !== filterStatus) return false;

      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        const idMatch = order.OrderNumber.toString().includes(query);
        const productsMatch = order.Items.some(item =>
          item.Product?.Name?.toLowerCase().includes(query)
        );
        if (!idMatch && !productsMatch) return false;
      }

      if (dateFrom) {
        const orderDate = new Date(order.CreatedAt || '').setHours(0, 0, 0, 0);
        const fromDate = new Date(dateFrom).setHours(0, 0, 0, 0);
        if (orderDate < fromDate) return false;
      }
      if (dateTo) {
        const orderDate = new Date(order.CreatedAt || '').setHours(23, 59, 59, 999);
        const toDate = new Date(dateTo).setHours(23, 59, 59, 999);
        if (orderDate > toDate) return false;
      }

      return true;
    }).sort((a, b) => {
      let comparison = 0;
      if (sortBy === 'date') {
        comparison = new Date(a.CreatedAt || '').getTime() - new Date(b.CreatedAt || '').getTime();
      } else if (sortBy === 'total') {
        comparison = a.TotalPrice - b.TotalPrice;
      } else if (sortBy === 'id') {
        comparison = a.OrderNumber - b.OrderNumber;
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

  function formatTotal(v: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(v);
  }

  function getStatusVariant(status: string) {
    switch (status) {
      case 'pending': return 'warning';
      case 'confirmed': return 'primary';
      case 'preparing': return 'info';
      case 'ready': return 'success';
      case 'delivered': return 'success';
      case 'cancelled': return 'danger';
      default: return 'default';
    }
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

  const countByStatus = $derived(
    orders.reduce<Record<string, number>>((acc, o) => {
      acc[o.Status] = (acc[o.Status] ?? 0) + 1;
      return acc;
    }, {})
  );
</script>

<Workspace
  breadcrumb={[{ label: 'Pedidos' }]}
  title="Gerenciamento de Pedidos"
  description="Acompanhe e gerencie todos os pedidos do restaurante"
>
  <svelte:fragment slot="actions">
    <a href="/orders/new">
      <Button variant="primary">
        <Plus size={16} />
        Novo Pedido
      </Button>
    </a>
  </svelte:fragment>

  {#if loading}
    <div class="skeleton-grid">
      {#each Array(6) as _}
        <Card class="skeleton-card">
          <div class="skeleton-card-header">
            <Skeleton variant="text" width="60px" height="16px" />
            <Skeleton variant="rectangular" width="80px" height="24px" />
          </div>
          <div class="skeleton-card-meta">
            <Skeleton variant="text" width="100px" height="12px" />
            <Skeleton variant="text" width="80px" height="12px" />
          </div>
          <div class="skeleton-card-items">
            <Skeleton variant="text" width="120px" height="12px" />
            <Skeleton variant="text" width="100px" height="12px" />
          </div>
          <div class="skeleton-card-footer">
            <Skeleton variant="text" width="80px" height="16px" />
          </div>
        </Card>
      {/each}
    </div>
  {:else if error}
    <Alert variant="error" dismissible onDismiss={() => error = ''}>
      ⚠️ {error}
      <Button onclick={loadOrders} size="sm">Tentar novamente</Button>
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
          <span class="stat-item">{filtered.length} pedidos</span>
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

      <!-- Search and Sort -->
      <div class="filters-row">
        <div class="search-wrapper">
          <Search size={16} class="search-icon" />
          <Input
            placeholder="Buscar por ID ou produto..."
            bind:value={searchQuery}
            class="search-input"
          />
        </div>
        <div class="date-filters">
          <div class="date-wrapper">
            <Calendar size={16} class="date-icon" />
            <Input
              type="date"
              placeholder="Data inicial"
              bind:value={dateFrom}
              class="date-input"
            />
          </div>
          <div class="date-wrapper">
            <Calendar size={16} class="date-icon" />
            <Input
              type="date"
              placeholder="Data final"
              bind:value={dateTo}
              class="date-input"
            />
          </div>
        </div>
        <div class="sort-wrapper">
          <select bind:value={sortBy} class="sort-select">
            <option value="date">Data</option>
            <option value="total">Valor</option>
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
        {#if searchQuery || dateFrom || dateTo}
          <Button variant="ghost" size="sm" onclick={() => { searchQuery = ''; dateFrom = ''; dateTo = ''; }}>
            Limpar
          </Button>
        {/if}
      </div>
    </Card>

    {#if filtered.length === 0}
      <div class="empty-state">
        <ShoppingBag size={48} class="empty-icon" />
        <span class="empty-title">{filterStatus === 'all' ? 'Nenhum pedido encontrado' : 'Nenhum pedido com esse status'}</span>
        <span class="empty-subtitle">{filterStatus === 'all' ? 'Comece criando um novo pedido' : 'Tente outro filtro'}</span>
        {#if filterStatus === 'all'}
          <a href="/orders/new">
            <Button variant="primary">Criar Pedido</Button>
          </a>
        {/if}
      </div>
    {:else}
      <div class="orders-grid">
        {#each paginated as order}
          <Card class="order-card">
            <div class="order-header">
              <div class="order-id">#{order.OrderNumber}</div>
              <Badge variant={getStatusVariant(order.Status)} size="sm">
                {ORDER_STATUS_LABEL[order.Status]}
              </Badge>
            </div>
            
            <div class="order-meta">
              <div class="order-date">
                <Clock size={14} />
                <span>{formatDate(order.CreatedAt)}</span>
              </div>
              <div class="order-table">
                <span>Mesa {order.TableNumber || 'N/A'}</span>
              </div>
            </div>

            <div class="order-items">
              {#each (order.Items ?? []).slice(0, 3) as item}
                <div class="order-item">
                  <span class="item-quantity">{item.Quantity}×</span>
                  <span class="item-name">{item.Product?.Name || `Produto #${item.ProductID}`}</span>
                </div>
              {/each}
              {#if (order.Items ?? []).length > 3}
                <div class="order-item more-items">
                  <span>+{order.Items.length - 3} itens</span>
                </div>
              {/if}
            </div>

            <div class="order-footer">
              <div class="order-total">
                <DollarSign size={16} />
                <span>{formatTotal(order.TotalPrice)}</span>
              </div>
              <Button href="/orders/{order.ID}" variant="ghost" size="sm">
                Ver detalhes
              </Button>
            </div>
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
    flex-wrap: wrap;
  }

  .search-wrapper {
    position: relative;
    flex: 1;
    min-width: 250px;
  }

  .search-icon {
    position: absolute;
    left: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: #94a3b8;
  }

  .search-input {
    padding-left: 2.5rem;
  }

  .date-filters {
    display: flex;
    gap: 0.5rem;
  }

  .date-wrapper {
    position: relative;
  }

  .date-icon {
    position: absolute;
    left: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: #94a3b8;
  }

  .date-input {
    padding-left: 2.5rem;
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

  /* Orders Grid */
  .orders-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
  }

  .order-card {
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .order-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px 0 rgb(0 0 0 / 0.08);
  }

  .order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .order-id {
    font-size: 1rem;
    font-weight: 700;
    color: #0f172a;
  }

  .order-meta {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
    font-size: 0.75rem;
    color: #64748b;
  }

  .order-date,
  .order-table {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .order-items {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
    padding: 0.75rem;
    background: #f8fafc;
    border-radius: 8px;
  }

  .order-item {
    display: flex;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .item-quantity {
    font-weight: 600;
    color: #6366f1;
  }

  .item-name {
    color: #0f172a;
  }

  .more-items {
    color: #64748b;
    font-size: 0.75rem;
  }

  .order-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 1rem;
    border-top: 1px solid #f1f5f9;
  }

  .order-total {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.125rem;
    font-weight: 700;
    color: #6366f1;
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

  .skeleton-card-items {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .skeleton-card-footer {
    display: flex;
    justify-content: flex-end;
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

  /* Responsive */
  @media (max-width: 768px) {
    .filters-row {
      flex-direction: column;
      align-items: stretch;
    }

    .search-wrapper,
    .date-filters,
    .sort-wrapper {
      width: 100%;
    }

    .date-filters {
      flex-direction: column;
    }

    .orders-grid {
      grid-template-columns: 1fr;
    }

    .status-pills {
      overflow-x: auto;
      padding-bottom: 0.5rem;
    }

    .status-pill {
      flex-shrink: 0;
    }
  }
</style>
