<script lang="ts">
  import { onMount } from 'svelte';
  import { request } from '$lib/api/client';
  import { Card, Button, Badge, Alert, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { TrendingUp, TrendingDown, AlertTriangle, Package, ShoppingCart, DollarSign, ArrowUpRight, ArrowDownRight, CheckCircle, XCircle, Info, Clock, Users, Activity, MoreHorizontal } from '@lucide/svelte';

  // Métricas executivas
  let metrics = $state({
    products: 0,
    orders: 0,
    pendingOrders: 0,
    todayOrders: 0,
    todayRevenue: 0,
    averageTicket: 0,
    lowStock: 0,
    pendingAdjustments: 0,
  });
  let loadingMetrics = $state(true);
  let error = $state('');

  // Últimos pedidos
  interface RecentOrder {
    id: number;
    table: string;
    total: number;
    status: string;
    time: string;
  }
  let recentOrders = $state<RecentOrder[]>([]);
  let loadingOrders = $state(true);

  // Ingredientes críticos
  interface CriticalIngredient {
    name: string;
    stock: number;
    unit: string;
    minStock: number;
  }
  let criticalIngredients = $state<CriticalIngredient[]>([]);
  let loadingIngredients = $state(true);

  // Atividades recentes
  let recentActivities = $state([
    { type: 'order', message: 'Novo pedido #1234 criado', time: '15min atrás', icon: ShoppingCart },
    { type: 'stock', message: 'Estoque de Tomate ajustado', time: '1h atrás', icon: Package },
    { type: 'approval', message: 'Ajuste #56 aprovado', time: '2h atrás', icon: CheckCircle },
    { type: 'alert', message: 'Estoque baixo: Queijo', time: '3h atrás', icon: AlertTriangle },
  ]);

  onMount(async () => {
    try {
      const [productsRes, ordersRes, adjustmentsRes] = await Promise.all([
        request('/products'),
        request('/orders'),
        request('/stock-adjustments'),
      ]);
      const products = productsRes.data;
      const orders = ordersRes.data;
      const adjustments = adjustmentsRes.data;

      // Calcular métricas
      const today = new Date().toISOString().split('T')[0];
      const todayOrdersList = Array.isArray(orders)
        ? orders.filter((o: any) => o.created_at?.startsWith(today))
        : [];

      const avgTicket = todayOrdersList.length > 0
        ? todayOrdersList.reduce((sum: number, o: any) => sum + (o.total_price || 0), 0) / todayOrdersList.length
        : 0;

      metrics = {
        products: Array.isArray(products) ? products.length : 0,
        orders: Array.isArray(orders) ? orders.length : 0,
        pendingOrders: Array.isArray(orders)
          ? orders.filter((o: any) => o.status === 'pending' || o.status === 'confirmed').length
          : 0,
        todayOrders: todayOrdersList.length,
        todayRevenue: todayOrdersList.reduce((sum: number, o: any) => sum + (o.total_price || 0), 0),
        averageTicket: avgTicket,
        lowStock: Array.isArray(products)
          ? products.filter((p: any) => (p.ingredients || []).some((i: any) => i.stock < 10)).length
          : 0,
        pendingAdjustments: Array.isArray(adjustments)
          ? adjustments.filter((a: any) => a.status === 'pending').length
          : 0,
      };

      // Últimos pedidos (últimos 5)
      recentOrders = Array.isArray(orders)
        ? orders.slice(0, 5).map((o: any) => ({
            id: o.id,
            table: o.table_number || 'N/A',
            total: o.total_price || 0,
            status: o.status,
            time: o.created_at ? new Date(o.created_at).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' }) : '—',
          }))
        : [];
      loadingOrders = false;

      // Ingredientes críticos (estoque < 10)
      if (Array.isArray(products)) {
        const allIngredients = products.flatMap((p: any) => p.ingredients || []);
        criticalIngredients = allIngredients
          .filter((i: any) => i.stock < 10)
          .map((i: any) => ({
            name: i.name,
            stock: i.stock,
            unit: i.unit || 'un',
            minStock: i.min_stock || 5,
          }))
          .slice(0, 5);
      }
      loadingIngredients = false;

    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar métricas.';
      loadingOrders = false;
      loadingIngredients = false;
    } finally {
      loadingMetrics = false;
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
      case 'completed': return 'success';
      case 'pending': return 'warning';
      case 'confirmed': return 'primary';
      case 'cancelled': return 'danger';
      default: return 'default';
    }
  }

  function getStatusLabel(status: string) {
    switch (status) {
      case 'completed': return 'Concluído';
      case 'pending': return 'Pendente';
      case 'confirmed': return 'Confirmado';
      case 'cancelled': return 'Cancelado';
      default: return status;
    }
  }
</script>

<Workspace
  breadcrumb={[{ label: 'Dashboard' }]}
  title="Painel Executivo"
  description="Visão geral das operações do restaurante em tempo real"
>
  {#if loadingMetrics}
    <div class="skeleton-metrics-grid">
      {#each Array(4) as _}
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
  {:else}
    <!-- KPIs Executivos -->
  <div class="kpi-grid">
    <Card class="kpi-card kpi-primary">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <ShoppingCart size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Pedidos Hoje</span>
      </div>
      <div class="kpi-value">{loadingMetrics ? '—' : formatNumber(metrics.todayOrders)}</div>
      <div class="kpi-footer">
        <div class="kpi-change positive">
          <ArrowUpRight size={14} />
          <span>+12% vs ontem</span>
        </div>
        <div class="kpi-sparkline positive"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-success">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <DollarSign size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Faturamento Hoje</span>
      </div>
      <div class="kpi-value">{loadingMetrics ? '—' : formatCurrency(metrics.todayRevenue)}</div>
      <div class="kpi-footer">
        <div class="kpi-change positive">
          <ArrowUpRight size={14} />
          <span>+8% vs ontem</span>
        </div>
        <div class="kpi-sparkline positive"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-info">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <Users size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Ticket Médio</span>
      </div>
      <div class="kpi-value">{loadingMetrics ? '—' : formatCurrency(metrics.averageTicket)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Estável</span>
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
      <div class="kpi-value">{loadingMetrics ? '—' : formatNumber(metrics.lowStock)}</div>
      <div class="kpi-footer">
        <div class="kpi-change negative">
          <ArrowDownRight size={14} />
          <span>+2 novos</span>
        </div>
        <div class="kpi-sparkline negative"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-danger">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <Clock size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Ajustes Pendentes</span>
      </div>
      <div class="kpi-value">{loadingMetrics ? '—' : formatNumber(metrics.pendingAdjustments)}</div>
      <div class="kpi-footer">
        <div class="kpi-change neutral">
          <span>Aguardando</span>
        </div>
        <div class="kpi-sparkline neutral"></div>
      </div>
    </Card>

    <Card class="kpi-card kpi-secondary">
      <div class="kpi-header">
        <div class="kpi-icon-wrapper">
          <Package size={20} class="kpi-icon" />
        </div>
        <span class="kpi-label">Produtos Ativos</span>
      </div>
      <div class="kpi-value">{loadingMetrics ? '—' : formatNumber(metrics.products)}</div>
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
        {#if loadingOrders}
          <div class="loading-state">
            <Activity class="spinner" size={24} />
            <span>Carregando pedidos...</span>
          </div>
        {:else if recentOrders.length === 0}
          <div class="empty-state">
            <ShoppingCart size={32} class="empty-icon" />
            <span>Nenhum pedido encontrado</span>
          </div>
        {:else}
          {#each recentOrders as order}
            <div class="order-item">
              <div class="order-info">
                <div class="order-id">#{order.id}</div>
                <div class="order-table">Mesa {order.table}</div>
              </div>
              <div class="order-status">
                <Badge variant={getStatusVariant(order.status)} size="sm">
                  {getStatusLabel(order.status)}
                </Badge>
              </div>
              <div class="order-time">{order.time}</div>
              <div class="order-total">{formatCurrency(order.total)}</div>
            </div>
          {/each}
        {/if}
      </div>
    </Card>

    <!-- Ingredientes Críticos -->
    <Card class="critical-ingredients-card">
      <div class="card-header">
        <h3 class="card-title">Ingredientes Críticos</h3>
        <Badge variant="danger" size="sm">{criticalIngredients.length} itens</Badge>
      </div>
      <div class="ingredients-list">
        {#if loadingIngredients}
          <div class="loading-state">
            <Activity class="spinner" size={24} />
            <span>Carregando ingredientes...</span>
          </div>
        {:else if criticalIngredients.length === 0}
          <div class="empty-state">
            <Package size={32} class="empty-icon" />
            <span>Estoque em dia</span>
          </div>
        {:else}
          {#each criticalIngredients as ingredient}
            <div class="ingredient-item">
              <div class="ingredient-info">
                <div class="ingredient-name">{ingredient.name}</div>
                <div class="ingredient-stock">
                  <span class="stock-value">{ingredient.stock} {ingredient.unit}</span>
                  <span class="stock-min">mín: {ingredient.minStock}</span>
                </div>
              </div>
              <div class="ingredient-status">
                <Badge variant="danger" size="sm">Crítico</Badge>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </Card>

    <!-- Atividades Recentes -->
    <Card class="activities-card">
      <div class="card-header">
        <h3 class="card-title">Atividades Recentes</h3>
        <Button variant="ghost" size="sm">
          <MoreHorizontal size={16} />
        </Button>
      </div>
      <div class="activities-list">
        {#each recentActivities as activity}
          <div class="activity-item activity-{activity.type}">
            <div class="activity-icon-wrapper">
              <svelte:component this={activity.icon} size={16} class="activity-icon" />
            </div>
            <div class="activity-content">
              <div class="activity-message">{activity.message}</div>
              <div class="activity-time">{activity.time}</div>
            </div>
          </div>
        {/each}
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
