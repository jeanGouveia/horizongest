<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import type { Dashboard } from '$lib/types/dashboard';
  import { Card, Button, Badge, Alert, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { TrendingUp, TrendingDown, AlertTriangle, Package, ShoppingCart, DollarSign, ArrowUpRight, ArrowDownRight, CheckCircle, XCircle, Info, Clock, Users, Activity, MoreHorizontal } from '@lucide/svelte';

  // Dashboard data
  let dashboard = $state<Dashboard | null>(null);
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    try {
      const res = await api.dashboard();
      if (res.error) {
        error = res.error;
      } else {
        dashboard = res.data;
      }
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar dashboard.';
    } finally {
      loading = false;
    }
  });

  function formatCurrency(value: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
  }

  function formatNumber(value: number) {
    return new Intl.NumberFormat('pt-BR').format(value);
  }

  function getStatusVariant(status: string) {
    switch (status) {
      case 'delivered': return 'success';
      case 'pending': return 'warning';
      case 'confirmed': return 'primary';
      case 'cancelled': return 'danger';
      case 'preparing': return 'info';
      case 'ready': return 'success';
      default: return 'default';
    }
  }

  function getStatusLabel(status: string) {
    switch (status) {
      case 'delivered': return 'Entregue';
      case 'pending': return 'Pendente';
      case 'confirmed': return 'Confirmado';
      case 'cancelled': return 'Cancelado';
      case 'preparing': return 'Preparando';
      case 'ready': return 'Pronto';
      default: return status;
    }
  }

  function getAverageTicket() {
    if (!dashboard || dashboard.metrics.todayOrders === 0) return 0;
    return dashboard.metrics.todayRevenue / dashboard.metrics.todayOrders;
  }
</script>

<Workspace
  breadcrumb={[{ label: 'Dashboard' }]}
  title="Painel Executivo"
  description="Visão geral das operações do restaurante em tempo real"
>
  {#if loading}
    <div class="skeleton-metrics-grid">
      {#each Array(6) as _}
        <Card class="skeleton-metric-card">
          <div class="skeleton-metric-header">
            <Skeleton variant="circular" width="32px" height="32px" />
            <Skeleton variant="text" width="80px" height="12px" />
          </div>
          <Skeleton variant="text" width="60px" height="24px" />
        </Card>
      {/each}
    </div>
  {:else if error}
    <Alert variant="error" dismissible onDismiss={() => error = ''}>
      ⚠️ {error}
    </Alert>
  {:else if dashboard}
    <!-- KPIs Executivos -->
  <div class="kpi-grid">
    <Card class="kpi-card kpi-primary">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <ShoppingCart size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Pedidos Hoje</span>
      </div>
      <div class="kpi-value">{formatNumber(dashboard.metrics.todayOrders)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Hoje</span>
        </div>
        <div class="kpi-sparkline neutral"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-success">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <DollarSign size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Faturamento Hoje</span>
      </div>
      <div class="kpi-value">{formatCurrency(dashboard.metrics.todayRevenue)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Hoje</span>
        </div>
        <div class="kpi-sparkline neutral"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-info">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <Users size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Ticket Médio</span>
      </div>
      <div class="kpi-value">{formatCurrency(getAverageTicket())}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Hoje</span>
        </div>
        <div class="kpi-sparkline neutral"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-warning">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <AlertTriangle size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Estoque Baixo</span>
      </div>
      <div class="kpi-value">{formatNumber(dashboard.metrics.lowStockCount)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Crítico</span>
        </div>
        <div class="kpi-sparkline negative"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-danger">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <Clock size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Pedidos Pendentes</span>
      </div>
      <div class="kpi-value">{formatNumber(dashboard.metrics.pendingOrders)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Aguardando</span>
        </div>
        <div class="kpi-sparkline warning"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-secondary">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <Package size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Produtos Ativos</span>
      </div>
      <div class="kpi-value">{formatNumber(dashboard.metrics.activeProducts)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Catálogo</span>
        </div>
        <div class="kpi-sparkline neutral"></div>
      </div>
    </Card>
  </div>

  <!-- Grid Principal -->
  <div class="main-grid">
    <!-- Últimos Pedidos -->
    <Card class="recent-orders-card">
      <div class="card-header">
        <h3 class="card-title">Últimos Pedidos</h3>
        <Button href="/orders" variant="ghost" size="sm">Ver Todos</Button>
      </div>
      <div class="orders-list">
        {#if dashboard.recentOrders.length === 0}
          <div class="empty-state">
            <ShoppingCart size={32} class="empty-icon" />
            <span>Nenhum pedido encontrado</span>
          </div>
        {:else}
          {#each dashboard.recentOrders as order}
            <div class="order-item">
              <div class="order-info">
                <div class="order-id">#{order.orderNumber}</div>
                <div class="order-table">{order.itemsCount} itens</div>
              </div>
              <div class="order-status">
                <Badge variant={getStatusVariant(order.status)} size="sm">
                  {getStatusLabel(order.status)}
                </Badge>
              </div>
              <div class="order-time">{order.createdAt}</div>
              <div class="order-total">{formatCurrency(order.totalPrice)}</div>
            </div>
          {/each}
        {/if}
      </div>
    </Card>

    <!-- Ingredientes Críticos -->
    <Card class="critical-ingredients-card">
      <div class="card-header">
        <h3 class="card-title">Ingredientes Críticos</h3>
        <Badge variant="danger" size="sm">{dashboard.lowStock.length} itens</Badge>
      </div>
      <div class="ingredients-list">
        {#if dashboard.lowStock.length === 0}
          <div class="empty-state">
            <Package size={32} class="empty-icon" />
            <span>Estoque em dia</span>
          </div>
        {:else}
          {#each dashboard.lowStock as ingredient}
            <div class="ingredient-item">
              <div class="ingredient-info">
                <div class="ingredient-name">{ingredient.name}</div>
                <div class="ingredient-stock">
                  <span class="stock-value">{ingredient.stockQuantity} {ingredient.unit}</span>
                  <span class="stock-min">mín: {ingredient.minStock}</span>
                </div>
              </div>
              <div class="ingredient-status">
                {#if ingredient.stockQuantity === 0}
                  <Badge variant="danger" size="sm">Zerado</Badge>
                {:else}
                  <Badge variant="danger" size="sm">Crítico</Badge>
                {/if}
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </Card>

    <!-- Atividades Recentes - Removido (não implementado no backend) -->
    <Card class="activities-card">
      <div class="card-header">
        <h3 class="card-title">Totais</h3>
      </div>
      <div class="activities-list">
        <div class="activity-item activity-order">
          <div class="activity-icon-wrapper">
            <Package size={16} class="activity-icon" />
          </div>
          <div class="activity-content">
            <div class="activity-message">Total de Produtos</div>
            <div class="activity-time">{formatNumber(dashboard.totalProducts)}</div>
          </div>
        </div>
        <div class="activity-item activity-stock">
          <div class="activity-icon-wrapper">
            <Package size={16} class="activity-icon" />
          </div>
          <div class="activity-content">
            <div class="activity-message">Total de Categorias</div>
            <div class="activity-time">{formatNumber(dashboard.totalCategories)}</div>
          </div>
        </div>
        <div class="activity-item activity-approval">
          <div class="activity-icon-wrapper">
            <Package size={16} class="activity-icon" />
          </div>
          <div class="activity-content">
            <div class="activity-message">Total de Ingredientes</div>
            <div class="activity-time">{formatNumber(dashboard.totalIngredients)}</div>
          </div>
        </div>
      </div>
    </Card>
  </div>

  {/if}
</Workspace>

<style>
  /* KPI Grid */
  .kpi-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .kpi-card {
    padding: 0.875rem;
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .kpi-card:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px 0 rgb(0 0 0 / 0.08);
  }

  .kpi-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .kpi-icon-wrapper {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f8fafc;
  }

  .kpi-primary .kpi-icon-wrapper {
    background: #eef2ff;
    color: #6366f1;
  }

  .kpi-success .kpi-icon-wrapper {
    background: #ecfdf5;
    color: #10b981;
  }

  .kpi-info .kpi-icon-wrapper {
    background: #f0f9ff;
    color: #0284c7;
  }

  .kpi-warning .kpi-icon-wrapper {
    background: #fffbeb;
    color: #d97706;
  }

  .kpi-danger .kpi-icon-wrapper {
    background: #fef2f2;
    color: #dc2626;
  }

  .kpi-secondary .kpi-icon-wrapper {
    background: #f8fafc;
    color: #64748b;
  }

  .kpi-label {
    font-size: 0.6875rem;
    font-weight: 500;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .kpi-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: #0f172a;
    margin-bottom: 0.5rem;
    letter-spacing: -0.025em;
  }

  .kpi-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .kpi-change {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.6875rem;
    font-weight: 500;
  }

  .kpi-change.positive {
    color: #10b981;
  }

  .kpi-change.negative {
    color: #ef4444;
  }

  .kpi-change.neutral {
    color: #64748b;
  }

  .kpi-sparkline {
    width: 40px;
    height: 16px;
    border-radius: 3px;
    background: linear-gradient(90deg, transparent 0%, currentColor 100%);
    opacity: 0.15;
  }

  .kpi-sparkline.positive {
    background: linear-gradient(90deg, transparent 0%, #10b981 100%);
  }

  .kpi-sparkline.negative {
    background: linear-gradient(90deg, transparent 0%, #ef4444 100%);
  }

  .kpi-sparkline.neutral {
    background: linear-gradient(90deg, transparent 0%, #64748b 100%);
  }

  /* Main Grid */
  .main-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  .card-title {
    font-size: 0.9375rem;
    font-weight: 600;
    color: #0f172a;
    margin: 0;
    letter-spacing: -0.025em;
  }

  /* Orders List */
  .orders-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .order-item {
    display: grid;
    grid-template-columns: auto auto auto 1fr;
    gap: 0.75rem;
    align-items: center;
    padding: 0.625rem;
    border-radius: 8px;
    background: #f8fafc;
    transition: background 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .order-item:hover {
    background: #f1f5f9;
  }

  .order-info {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .order-id {
    font-size: 0.8125rem;
    font-weight: 600;
    color: #0f172a;
  }

  .order-table {
    font-size: 0.6875rem;
    color: #64748b;
  }

  .order-time {
    font-size: 0.6875rem;
    color: #64748b;
  }

  .order-total {
    font-size: 0.8125rem;
    font-weight: 600;
    color: #0f172a;
    text-align: right;
  }

  /* Ingredients List */
  .ingredients-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .ingredient-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.625rem;
    border-radius: 8px;
    background: #fef2f2;
    border: 1px solid #fee2e2;
  }

  .ingredient-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .ingredient-name {
    font-size: 0.8125rem;
    font-weight: 500;
    color: #0f172a;
  }

  .ingredient-stock {
    display: flex;
    gap: 0.375rem;
    font-size: 0.6875rem;
  }

  .stock-value {
    font-weight: 600;
    color: #dc2626;
  }

  .stock-min {
    color: #64748b;
  }

  /* Activities List */
  .activities-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .activity-item {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    padding: 0.625rem;
    border-radius: 8px;
    background: #f8fafc;
    transition: background 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .activity-item:hover {
    background: #f1f5f9;
  }

  .activity-icon-wrapper {
    width: 28px;
    height: 28px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .activity-item.activity-order .activity-icon-wrapper {
    background: #eef2ff;
    color: #6366f1;
  }

  .activity-item.activity-stock .activity-icon-wrapper {
    background: #f0f9ff;
    color: #0284c7;
  }

  .activity-item.activity-approval .activity-icon-wrapper {
    background: #ecfdf5;
    color: #10b981;
  }

  .activity-item.activity-alert .activity-icon-wrapper {
    background: #fffbeb;
    color: #d97706;
  }

  .activity-content {
    flex: 1;
    min-width: 0;
  }

  .activity-message {
    font-size: 0.8125rem;
    font-weight: 500;
    color: #0f172a;
  }

  .activity-time {
    font-size: 0.6875rem;
    color: #64748b;
    margin-top: 0.125rem;
  }

  /* Loading & Empty States */
  .loading-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1.5rem;
    color: #64748b;
    font-size: 0.8125rem;
  }

  .spinner {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .empty-icon {
    color: #cbd5e1;
  }


  /* Skeleton Loading */
  .skeleton-metrics-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 0.75rem;
  }

  .skeleton-metric-card {
    padding: 0.875rem;
  }

  .skeleton-metric-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
  }

  /* Responsive */
  @media (max-width: 1200px) {
    .kpi-grid {
      grid-template-columns: repeat(3, 1fr);
    }

    .main-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 768px) {
    .kpi-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 0.5rem;
    }

    .main-grid {
      grid-template-columns: 1fr;
    }

    .order-item {
      grid-template-columns: 1fr auto;
      gap: 0.5rem;
    }

    .order-total {
      grid-column: 2;
    }
  }

  @media (max-width: 480px) {
    .kpi-grid {
      grid-template-columns: 1fr;
    }

    .kpi-value {
      font-size: 1.5rem;
    }
  }
</style>
