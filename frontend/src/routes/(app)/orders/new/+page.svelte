<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { getActiveProducts } from '$lib/api/product';
  import { createOrder } from '$lib/api/order';
  import type { Product } from '$lib/types/product';
  import type { OrderCreatePayload } from '$lib/types/order';

  interface CartItem {
    product: Product;
    quantity: number;
  }

  let products: Product[] = $state([]);
  let cart: CartItem[] = $state([]);
  let notes = $state('');
  let loading = $state(true);
  let submitting = $state(false);
  let error = $state('');
  let searchQuery = $state('');

  onMount(async () => {
    loading = true;
    try {
      products = await getActiveProducts();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar produtos.';
    } finally {
      loading = false;
    }
  });

  const filteredProducts = $derived(
    products.filter((p) =>
      p.Name.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  function addToCart(product: Product) {
    const existing = cart.find((c) => c.product.ID === product.ID);
    if (existing) {
      cart = cart.map((c) =>
        c.product.ID === product.ID ? { ...c, quantity: c.quantity + 1 } : c
      );
    } else {
      cart = [...cart, { product, quantity: 1 }];
    }
  }

  function removeFromCart(productId: number) {
    cart = cart.filter((c) => c.product.ID !== productId);
  }

  function updateQty(productId: number, qty: number) {
    if (qty <= 0) {
      removeFromCart(productId);
      return;
    }
    cart = cart.map((c) => (c.product.ID === productId ? { ...c, quantity: qty } : c));
  }

  const cartTotal = $derived(
    cart.reduce((sum, c) => sum + c.product.Price * c.quantity, 0)
  );

  const cartCount = $derived(cart.reduce((sum, c) => sum + c.quantity, 0));

  function formatPrice(v: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(v);
  }

  function getCartQty(productId: number) {
    return cart.find((c) => c.product.ID === productId)?.quantity ?? 0;
  }

  async function submitOrder() {
    if (cart.length === 0) return;
    submitting = true;
    error = '';
    const payload: OrderCreatePayload = {
      notes: notes.trim() || undefined,
      items: cart.map((c) => ({ product_id: c.product.ID, quantity: c.quantity })),
    };
    console.log('Frontend - Payload enviado:', payload);
    try {
      const order = await createOrder(payload);
      console.log('Frontend - Pedido criado:', order);
      console.log('Frontend - order.ID:', order.ID);
      console.log('Frontend - order.id:', order.id);
      console.log('Frontend - Tentando goto para:', `/orders/${order.ID}`);
      goto(`/orders/${order.ID}`);
    } catch (e: any) {
      console.error('Frontend - Erro ao criar pedido:', e);
      console.error('Frontend - Stack trace:', e.stack);
      
      // Formatar mensagem de erro de estoque insuficiente
      const errorMsg = e?.message ?? 'Erro ao criar pedido.';
      if (errorMsg.includes('Ingredientes insuficientes')) {
        // A mensagem já vem formatada do backend, apenas preservar as quebras de linha
        error = errorMsg;
      } else {
        error = errorMsg;
      }
      submitting = false;
    }
  }
</script>

<div class="page-wrapper">
  <a href="/orders" class="back-link">← Voltar para Pedidos</a>
  <h1 class="page-title">Novo Pedido</h1>

  {#if error}
    <div class="alert alert-error">{error}</div>
  {/if}

  <div class="order-layout">
    <!-- Painel esquerdo: seleção de produtos -->
    <div class="products-panel">
      <div class="search-box">
        <span class="search-icon">🔍</span>
        <input
          type="search"
          bind:value={searchQuery}
          placeholder="Buscar produto…"
          class="search-input"
        />
      </div>

      {#if loading}
        <div class="loading-state">
          <div class="spinner"></div>
          <span>Carregando produtos…</span>
        </div>
      {:else if filteredProducts.length === 0}
        <p class="empty-note">Nenhum produto encontrado.</p>
      {:else}
        <div class="product-grid">
          {#each filteredProducts as product}
            {@const qty = getCartQty(product.ID)}
            <div class="product-tile" class:in-cart={qty > 0}>
              <div class="tile-info">
                <span class="tile-name">{product.Name}</span>
                {#if product.IsComposto}
                  <span class="tile-tag">Composto</span>
                {/if}
                <span class="tile-price">{formatPrice(product.Price)}</span>
              </div>
              <div class="tile-actions">
                {#if qty > 0}
                  <div class="qty-control">
                    <button onclick={() => updateQty(product.ID, qty - 1)}>−</button>
                    <span>{qty}</span>
                    <button onclick={() => addToCart(product)}>+</button>
                  </div>
                {:else}
                  <button class="btn-add" onclick={() => addToCart(product)}>Adicionar</button>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Painel direito: resumo do pedido -->
    <aside class="cart-panel">
      <h2 class="cart-title">
        Resumo
        {#if cartCount > 0}
          <span class="cart-count">{cartCount}</span>
        {/if}
      </h2>

      {#if cart.length === 0}
        <div class="cart-empty">
          <p>🛒</p>
          <p>Selecione itens ao lado para montar o pedido.</p>
        </div>
      {:else}
        <ul class="cart-list">
          {#each cart as item}
            <li class="cart-item">
              <div class="cart-item-info">
                <span class="cart-item-name">{item.product.Name}</span>
                <span class="cart-item-unit">{formatPrice(item.product.Price)} × {item.quantity}</span>
              </div>
              <div class="cart-item-right">
                <span class="cart-item-sub">{formatPrice(item.product.Price * item.quantity)}</span>
                <button class="btn-remove" onclick={() => removeFromCart(item.product.ID)}>✕</button>
              </div>
            </li>
          {/each}
        </ul>

        <div class="cart-total">
          <span>Total</span>
          <span class="total-value">{formatPrice(cartTotal)}</span>
        </div>

        <div class="form-group">
          <label for="notes">Observações</label>
          <textarea id="notes" bind:value={notes} rows="3" placeholder="Ex: sem cebola, prato para viagem…"></textarea>
        </div>

        <button
          class="btn btn-primary btn-full"
          onclick={submitOrder}
          disabled={submitting || cart.length === 0}
        >
          {submitting ? 'Enviando pedido…' : `Confirmar Pedido · ${formatPrice(cartTotal)}`}
        </button>
      {/if}
    </aside>
  </div>
</div>

<style>
  .page-wrapper   { max-width: 1100px; margin: 0 auto; padding: 2rem 1.5rem; }
  .back-link      { display: inline-flex; align-items: center; gap: 0.3rem; color: var(--color-muted); font-size: 0.875rem; text-decoration: none; margin-bottom: 1.5rem; }
  .back-link:hover { color: var(--color-text); }
  .page-title     { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0 0 1.75rem; }

  .order-layout   { display: grid; grid-template-columns: 1fr 340px; gap: 1.5rem; align-items: start; }

  /* Painel produtos */
  .products-panel { min-width: 0; }
  .search-box     { position: relative; margin-bottom: 1.25rem; }
  .search-icon    { position: absolute; left: 0.8rem; top: 50%; transform: translateY(-50%); pointer-events: none; font-size: 0.85rem; }
  .search-input   { width: 100%; box-sizing: border-box; border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; padding: 0.6rem 0.75rem 0.6rem 2.2rem; font-size: 0.9rem; background: var(--color-surface, #fff); color: var(--color-text); }
  .search-input:focus { outline: none; border-color: var(--color-primary, #e85d04); box-shadow: 0 0 0 3px rgba(232,93,4,0.1); }
  .product-grid   { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; }
  .product-tile   { background: var(--color-surface, #fff); border: 1.5px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; transition: border-color 0.15s; }
  .product-tile.in-cart { border-color: var(--color-primary, #e85d04); }
  .tile-info      { display: flex; flex-direction: column; gap: 0.2rem; flex: 1; }
  .tile-name      { font-weight: 600; font-size: 0.9rem; color: var(--color-text); }
  .tile-tag       { font-size: 0.72rem; color: var(--color-muted); }
  .tile-price     { font-weight: 700; color: var(--color-primary, #e85d04); font-size: 0.9rem; margin-top: 0.3rem; }
  .tile-actions   { display: flex; }
  .btn-add        { width: 100%; padding: 0.4rem; border: 1.5px solid var(--color-primary, #e85d04); background: transparent; color: var(--color-primary, #e85d04); border-radius: 0.4rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: background 0.15s; }
  .btn-add:hover  { background: var(--color-primary, #e85d04); color: #fff; }
  .qty-control    { display: flex; align-items: center; gap: 0; width: 100%; border: 1.5px solid var(--color-primary, #e85d04); border-radius: 0.4rem; overflow: hidden; }
  .qty-control button { flex: 1; background: transparent; border: none; color: var(--color-primary, #e85d04); font-size: 1.1rem; font-weight: 700; cursor: pointer; padding: 0.3rem; transition: background 0.12s; }
  .qty-control button:hover { background: rgba(232,93,4,0.1); }
  .qty-control span { flex: 1; text-align: center; font-weight: 700; font-size: 0.9rem; color: var(--color-text); border-left: 1px solid var(--color-primary, #e85d04); border-right: 1px solid var(--color-primary, #e85d04); padding: 0.3rem 0; }

  /* Painel carrinho */
  .cart-panel     { background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 1rem; padding: 1.5rem; position: sticky; top: 1rem; }
  .cart-title     { font-size: 1.1rem; font-weight: 700; margin: 0 0 1.25rem; display: flex; align-items: center; gap: 0.5rem; }
  .cart-count     { background: var(--color-primary, #e85d04); color: #fff; font-size: 0.72rem; font-weight: 700; border-radius: 99px; padding: 0.1rem 0.5rem; }
  .cart-empty     { text-align: center; color: var(--color-muted); font-size: 0.875rem; padding: 2rem 0; line-height: 1.8; }
  .cart-list      { list-style: none; margin: 0 0 1rem; padding: 0; display: flex; flex-direction: column; gap: 0.5rem; }
  .cart-item      { display: flex; align-items: flex-start; justify-content: space-between; gap: 0.5rem; }
  .cart-item-info { flex: 1; }
  .cart-item-name { font-size: 0.875rem; font-weight: 500; display: block; }
  .cart-item-unit { font-size: 0.78rem; color: var(--color-muted); }
  .cart-item-right { display: flex; align-items: center; gap: 0.4rem; flex-shrink: 0; }
  .cart-item-sub  { font-size: 0.85rem; font-weight: 600; color: var(--color-text); }
  .btn-remove     { background: none; border: none; color: #9ca3af; cursor: pointer; font-size: 0.75rem; padding: 0.1rem 0.2rem; border-radius: 3px; }
  .btn-remove:hover { color: #ef4444; background: #fee2e2; }
  .cart-total     { display: flex; justify-content: space-between; align-items: center; border-top: 1px solid var(--color-border, #e5e7eb); padding: 0.85rem 0; margin-bottom: 1rem; font-weight: 600; }
  .total-value    { font-size: 1.1rem; color: var(--color-primary, #e85d04); }

  .form-group     { display: flex; flex-direction: column; gap: 0.4rem; margin-bottom: 1rem; }
  .form-group label { font-size: 0.8rem; font-weight: 500; color: var(--color-text); }
  .form-group textarea { border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; padding: 0.55rem 0.75rem; font-size: 0.875rem; resize: vertical; background: var(--color-surface, #fff); color: var(--color-text); width: 100%; box-sizing: border-box; }
  .form-group textarea:focus { outline: none; border-color: var(--color-primary, #e85d04); }

  .btn            { display: inline-flex; align-items: center; justify-content: center; padding: 0.6rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-full       { width: 100%; }

  .loading-state  { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 3rem; color: var(--color-muted); }
  .spinner        { width: 1.75rem; height: 1.75rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin  { to { transform: rotate(360deg); } }
  .empty-note     { color: var(--color-muted); font-size: 0.9rem; text-align: center; padding: 2rem; }
  .alert-error    { background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 0.85rem 1.1rem; border-radius: 0.5rem; margin-bottom: 1.5rem; font-size: 0.9rem; white-space: pre-line; line-height: 1.6; }

  @media (max-width: 768px) {
    .order-layout { grid-template-columns: 1fr; }
    .cart-panel   { position: static; }
  }
</style>
