<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';

  // Contadores para dar contexto real aos cards
  let counts = $state({ products: 0, orders: 0, pendingOrders: 0 });
  let loadingCounts = $state(true);

  onMount(async () => {
    try {
      const [productsRes, ordersRes] = await Promise.all([
        api.get('/products'),
        api.get('/orders'),
      ]);
      const products = await productsRes.json();
      const orders = await ordersRes.json();
      counts = {
        products: Array.isArray(products) ? products.length : 0,
        orders: Array.isArray(orders) ? orders.length : 0,
        pendingOrders: Array.isArray(orders)
          ? orders.filter((o: any) => o.status === 'pending' || o.status === 'confirmed').length
          : 0,
      };
    } catch {
      // silencia — contadores não são críticos
    } finally {
      loadingCounts = false;
    }
  });
</script>

<div class="page-wrapper">
  <header class="page-header">
    <div>
      <h1 class="page-title">Dashboard</h1>
      <p class="page-subtitle">Visão geral do sistema — clique em um módulo para começar</p>
    </div>
  </header>

  <!-- Cards de módulos -->
  <div class="module-grid">

    <!-- PRODUTOS -->
    <div class="module-card">
      <div class="module-icon">🍽️</div>
      <div class="module-body">
        <h2 class="module-title">Produtos</h2>
        <p class="module-desc">
          Cadastre e gerencie os itens do cardápio. Dentro de cada produto você define
          a composição com ingredientes e quantidades usadas.
        </p>
        <div class="module-note">
          💡 Ingredientes também são gerenciados aqui — acesse um produto e vá em "Composição do Prato".
        </div>
        {#if !loadingCounts}
          <p class="module-stat">{counts.products} produto{counts.products !== 1 ? 's' : ''} cadastrado{counts.products !== 1 ? 's' : ''}</p>
        {/if}
      </div>
      <div class="module-actions">
        <a href="/products" class="btn btn-primary">Ver Cardápio</a>
        <a href="/products" class="btn btn-ghost">+ Novo Produto</a>
      </div>
    </div>

    <!-- PEDIDOS -->
    <div class="module-card">
      <div class="module-icon">🧾</div>
      <div class="module-body">
        <h2 class="module-title">Pedidos</h2>
        <p class="module-desc">
          Acompanhe todos os pedidos realizados. Ao confirmar um pedido, o sistema
          realiza a baixa automática de estoque dos ingredientes utilizados.
        </p>
        {#if !loadingCounts && counts.pendingOrders > 0}
          <div class="module-alert">
            ⚡ {counts.pendingOrders} pedido{counts.pendingOrders !== 1 ? 's' : ''} aguardando atenção
          </div>
        {/if}
        {#if !loadingCounts}
          <p class="module-stat">{counts.orders} pedido{counts.orders !== 1 ? 's' : ''} no total</p>
        {/if}
      </div>
      <div class="module-actions">
        <a href="/orders/new" class="btn btn-primary">+ Novo Pedido</a>
        <a href="/orders" class="btn btn-ghost">Ver Todos</a>
      </div>
    </div>

    <!-- INGREDIENTES -->
    <div class="module-card module-card--secondary">
      <div class="module-icon">🥦</div>
      <div class="module-body">
        <h2 class="module-title">Ingredientes &amp; Estoque</h2>
        <p class="module-desc">
          Os ingredientes são gerenciados dentro do módulo de Produtos. Você pode
          cadastrar novos ingredientes e vinculá-los a cada prato com a quantidade exata usada.
        </p>
        <div class="module-note">
          📍 Para acessar: vá em <strong>Produtos</strong> e clique em qualquer produto,
          ou use o botão <strong>"+ Ingrediente"</strong> na tela de listagem.
        </div>
      </div>
      <div class="module-actions">
        <a href="/products" class="btn btn-secondary">Ir para Produtos</a>
      </div>
    </div>

    <!-- AJUSTES DE ESTOQUE -->
    <div class="module-card module-card--secondary">
      <div class="module-icon">📦</div>
      <div class="module-body">
        <h2 class="module-title">Ajustes de Estoque</h2>
        <p class="module-desc">
          Aprove ou rejeite ajustes de estoque pendentes de pedidos cancelados.
        </p>
      </div>
      <div class="module-actions">
        <a href="/stock-adjustments" class="btn btn-secondary">Ver Ajustes</a>
      </div>
    </div>

    <!-- FLUXO DE TESTES -->
    <div class="module-card module-card--guide">
      <div class="module-icon">🧪</div>
      <div class="module-body">
        <h2 class="module-title">Como testar o fluxo completo</h2>
        <ol class="guide-steps">
          <li>
            <span class="step-num">1</span>
            <span>Acesse <strong>Produtos → "+ Ingrediente"</strong> e cadastre ao menos um ingrediente com estoque</span>
          </li>
          <li>
            <span class="step-num">2</span>
            <span>Crie um produto em <strong>Produtos → "+ Produto"</strong></span>
          </li>
          <li>
            <span class="step-num">3</span>
            <span>Abra o produto criado e vincule ingredientes com quantidade em <strong>"Composição do Prato"</strong></span>
          </li>
          <li>
            <span class="step-num">4</span>
            <span>Vá em <strong>"+ Novo Pedido"</strong>, adicione o produto ao carrinho e confirme</span>
          </li>
          <li>
            <span class="step-num">5</span>
            <span>Verifique no detalhe do pedido o status e confira que o estoque foi baixado</span>
          </li>
        </ol>
      </div>
    </div>

  </div>
</div>

<style>
  .page-wrapper    { max-width: 1100px; margin: 0 auto; padding: 2rem 1.5rem; }
  .page-header     { margin-bottom: 2rem; }
  .page-title      { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .page-subtitle   { font-size: 0.9rem; color: var(--color-muted); margin: 0.25rem 0 0; }

  /* Grid de módulos — 2 colunas em desktop, 1 em mobile */
  .module-grid     { display: grid; grid-template-columns: repeat(auto-fill, minmax(420px, 1fr)); gap: 1.25rem; }

  /* Card base */
  .module-card {
    background: var(--color-surface, #fff);
    border: 1px solid var(--color-border, #e5e7eb);
    border-radius: 1rem;
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    transition: box-shadow 0.15s;
  }
  .module-card:hover { box-shadow: 0 4px 20px rgba(0,0,0,0.07); }

  /* Variantes */
  .module-card--secondary { border-color: var(--color-border, #e5e7eb); background: var(--color-surface-2, #fafafa); }
  .module-card--guide     { border-style: dashed; background: transparent; }

  .module-icon  { font-size: 2rem; line-height: 1; }
  .module-body  { flex: 1; display: flex; flex-direction: column; gap: 0.5rem; }
  .module-title { font-size: 1.1rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .module-desc  { font-size: 0.875rem; color: var(--color-muted); margin: 0; line-height: 1.55; }
  .module-stat  { font-size: 0.8rem; font-weight: 600; color: var(--color-muted); margin: 0; }

  .module-note  {
    font-size: 0.8rem;
    background: var(--color-surface-2, #f9fafb);
    border-left: 3px solid var(--color-primary, #e85d04);
    padding: 0.5rem 0.75rem;
    border-radius: 0 0.4rem 0.4rem 0;
    color: var(--color-text);
    line-height: 1.5;
  }

  .module-alert {
    font-size: 0.8rem;
    font-weight: 600;
    background: #fef9c3;
    color: #854d0e;
    padding: 0.4rem 0.7rem;
    border-radius: 0.4rem;
    border: 1px solid #fde68a;
  }

  /* Ações do card */
  .module-actions { display: flex; gap: 0.6rem; flex-wrap: wrap; }

  /* Guia de testes */
  .guide-steps  { list-style: none; margin: 0.25rem 0 0; padding: 0; display: flex; flex-direction: column; gap: 0.6rem; }
  .guide-steps li { display: flex; align-items: flex-start; gap: 0.6rem; font-size: 0.85rem; color: var(--color-muted); line-height: 1.5; }
  .step-num     {
    flex-shrink: 0;
    width: 1.4rem; height: 1.4rem;
    border-radius: 50%;
    background: var(--color-primary, #e85d04);
    color: #fff;
    font-size: 0.7rem;
    font-weight: 700;
    display: flex; align-items: center; justify-content: center;
  }

  /* Botões */
  .btn          { display: inline-flex; align-items: center; padding: 0.5rem 1rem; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; text-decoration: none; }
  .btn-primary  { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover { background: var(--color-primary-dark, #c84e00); }
  .btn-secondary { background: var(--color-surface-2, #f3f4f6); color: var(--color-text); border: 1px solid var(--color-border, #d1d5db); }
  .btn-secondary:hover { background: var(--color-border, #e5e7eb); }
  .btn-ghost    { background: transparent; color: var(--color-primary, #e85d04); border: 1.5px solid var(--color-primary, #e85d04); }
  .btn-ghost:hover { background: rgba(232,93,4,0.06); }

  @media (max-width: 640px) {
    .module-grid { grid-template-columns: 1fr; }
  }
</style>
