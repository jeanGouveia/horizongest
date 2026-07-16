<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { getActiveProducts } from '$lib/api/product';
  import { createOrder } from '$lib/api/order';
  import type { Product } from '$lib/types/product';
  import type { OrderCreatePayload } from '$lib/types/order';
  import { Workspace } from '$lib/components/layout';
  import { Button, Input, Card, Alert } from '$lib/components/ui';
  import { Search, ShoppingCart, Plus, Minus, X, DollarSign, Activity, Utensils, Coffee, Cake, Beef } from '@lucide/svelte';

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
  let selectedCategory = $state<string>('all');

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

  const categories = $derived([
    { value: 'all', label: 'Todos', icon: Utensils },
    { value: 'bebidas', label: 'Bebidas', icon: Coffee },
    { value: 'pratos', label: 'Pratos', icon: Beef },
    { value: 'sobremesas', label: 'Sobremesas', icon: Cake },
  ]);

  const filteredProducts = $derived(
    products.filter((p) => {
      const matchesSearch = p.Name.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesCategory = selectedCategory === 'all' || 
        (selectedCategory === 'pratos' && !p.Name.toLowerCase().includes('suco') && !p.Name.toLowerCase().includes('refrigerante') && !p.Name.toLowerCase().includes('café') && !p.Name.toLowerCase().includes('bolo') && !p.Name.toLowerCase().includes('torta')) ||
        (selectedCategory === 'bebidas' && (p.Name.toLowerCase().includes('suco') || p.Name.toLowerCase().includes('refrigerante') || p.Name.toLowerCase().includes('café'))) ||
        (selectedCategory === 'sobremesas' && (p.Name.toLowerCase().includes('bolo') || p.Name.toLowerCase().includes('torta')));
      return matchesSearch && matchesCategory;
    })
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
    try {
      const order = await createOrder(payload);
      goto(`/orders/${order.ID}`);
    } catch (e: any) {
      const errorMsg = e?.message ?? 'Erro ao criar pedido.';
      if (errorMsg.includes('Ingredientes insuficientes')) {
        error = errorMsg;
      } else {
        error = errorMsg;
      }
      submitting = false;
    }
  }
</script>

<Workspace
  breadcrumb={[{ label: 'Pedidos', href: '/orders' }, { label: 'Novo Pedido' }]}
  title="Point of Sale"
  description="Selecione os produtos para criar um novo pedido"
>
  {#if error}
    <Alert variant="error" dismissible onDismiss={() => error = ''}>
      ⚠️ {error}
    </Alert>
  {/if}

  <div class="pos-layout">
    <!-- Left Panel: Products -->
    <div class="pos-products">
      <!-- Search -->
      <div class="pos-search">
        <div class="search-wrapper">
          <Search size={18} class="search-icon" />
          <Input
            placeholder="Buscar produtos..."
            bind:value={searchQuery}
            class="search-input"
          />
        </div>
      </div>

      <!-- Categories -->
      <div class="pos-categories">
        {#each categories as cat}
          <button
            class="category-btn"
            class:active={selectedCategory === cat.value}
            onclick={() => selectedCategory = cat.value}
          >
            <svelte:component this={cat.icon} size={18} />
            <span>{cat.label}</span>
          </button>
        {/each}
      </div>

      <!-- Products Grid -->
      {#if loading}
        <div class="loading-state">
          <Activity class="spinner" size={32} />
          <span>Carregando produtos...</span>
        </div>
      {:else if filteredProducts.length === 0}
        <div class="empty-state">
          <Utensils size={48} class="empty-icon" />
          <span class="empty-title">Nenhum produto encontrado</span>
          <span class="empty-subtitle">Tente outro termo de busca ou categoria</span>
        </div>
      {:else}
        <div class="products-grid">
          {#each filteredProducts as product}
            {@const qty = getCartQty(product.ID)}
            <Card class={`product-card ${qty > 0 ? 'in-cart' : ''}`}>
              <div class="product-image-placeholder">
                <Utensils size={32} />
              </div>
              <div class="product-info">
                <div class="product-name">{product.Name}</div>
                {#if product.IsComposto}
                  <span class="product-tag">Composto</span>
                {/if}
                <div class="product-price">{formatPrice(product.Price)}</div>
              </div>
              <div class="product-actions">
                {#if qty > 0}
                  <div class="qty-control">
                    <button onclick={() => updateQty(product.ID, qty - 1)}>
                      <Minus size={16} />
                    </button>
                    <span>{qty}</span>
                    <button onclick={() => addToCart(product)}>
                      <Plus size={16} />
                    </button>
                  </div>
                {:else}
                  <Button onclick={() => addToCart(product)} variant="primary" class="add-btn">
                    <Plus size={16} />
                    Adicionar
                  </Button>
                {/if}
              </div>
            </Card>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Right Panel: Cart -->
    <div class="pos-cart">
      <Card class="cart-card">
        <div class="cart-header">
          <div class="cart-title">
            <ShoppingCart size={20} />
            <span>Carrinho</span>
          </div>
          {#if cartCount > 0}
            <span class="cart-count">{cartCount}</span>
          {/if}
        </div>

        {#if cart.length === 0}
          <div class="cart-empty">
            <ShoppingCart size={48} class="empty-icon" />
            <span class="empty-title">Carrinho vazio</span>
            <span class="empty-subtitle">Selecione produtos ao lado</span>
          </div>
        {:else}
          <div class="cart-items">
            {#each cart as item}
              <div class="cart-item">
                <div class="cart-item-info">
                  <div class="cart-item-name">{item.product.Name}</div>
                  <div class="cart-item-unit">{formatPrice(item.product.Price)} × {item.quantity}</div>
                </div>
                <div class="cart-item-right">
                  <div class="cart-item-sub">{formatPrice(item.product.Price * item.quantity)}</div>
                  <button class="btn-remove" onclick={() => removeFromCart(item.product.ID)}>
                    <X size={16} />
                  </button>
                </div>
              </div>
            {/each}
          </div>

          <div class="cart-total">
            <span>Total</span>
            <div class="total-value">
              <DollarSign size={18} />
              <span>{formatPrice(cartTotal)}</span>
            </div>
          </div>

          <div class="cart-notes">
            <label class="notes-label">Observações</label>
            <Input
              bind:value={notes}
              placeholder="Ex: sem cebola, para viagem..."
              class="notes-input"
            />
          </div>

          <Button
            onclick={submitOrder}
            disabled={submitting || cart.length === 0}
            variant="primary"
            class="checkout-btn"
            loading={submitting}
          >
            <span>Confirmar Pedido</span>
            <span class="checkout-total">{formatPrice(cartTotal)}</span>
          </Button>
        {/if}
      </Card>
    </div>
  </div>
</Workspace>

<style>
  /* POS Layout */
  .pos-layout {
    display: grid;
    grid-template-columns: 1fr 380px;
    gap: 2rem;
    align-items: start;
  }

  /* Products Panel */
  .pos-products {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .pos-search {
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .search-wrapper {
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 0.875rem;
    top: 50%;
    transform: translateY(-50%);
    color: #94a3b8;
  }

  .search-input {
    padding-left: 2.75rem;
  }

  /* Categories */
  .pos-categories {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .category-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1rem;
    border: 1px solid #f1f5f9;
    background: #ffffff;
    color: #64748b;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .category-btn:hover {
    border-color: #6366f1;
    color: #6366f1;
  }

  .category-btn.active {
    background: #6366f1;
    border-color: #6366f1;
    color: #ffffff;
  }

  /* Products Grid */
  .products-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }

  .product-card {
    padding: 1.25rem;
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .product-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px 0 rgb(0 0 0 / 0.08);
  }

  .product-card.in-cart {
    border-color: #6366f1;
    background: #eef2ff;
  }

  .product-image-placeholder {
    width: 100%;
    aspect-ratio: 16/10;
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #cbd5e1;
    margin-bottom: 1rem;
  }

  .product-info {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
    margin-bottom: 1rem;
  }

  .product-name {
    font-weight: 600;
    font-size: 0.9375rem;
    color: #0f172a;
  }

  .product-tag {
    font-size: 0.75rem;
    color: #64748b;
  }

  .product-price {
    font-weight: 700;
    font-size: 1.125rem;
    color: #6366f1;
  }

  .product-actions {
    display: flex;
  }

  .add-btn {
    width: 100%;
  }

  .qty-control {
    display: flex;
    align-items: center;
    width: 100%;
    border: 1px solid #6366f1;
    border-radius: 8px;
    overflow: hidden;
  }

  .qty-control button {
    flex: 1;
    background: transparent;
    border: none;
    color: #6366f1;
    padding: 0.5rem;
    cursor: pointer;
    transition: background 0.15s;
  }

  .qty-control button:hover {
    background: rgba(99, 102, 241, 0.1);
  }

  .qty-control span {
    flex: 1;
    text-align: center;
    font-weight: 700;
    font-size: 0.875rem;
    color: #0f172a;
    border-left: 1px solid #6366f1;
    border-right: 1px solid #6366f1;
  }

  /* Cart Panel */
  .pos-cart {
    position: sticky;
    top: 0;
  }

  .cart-card {
    padding: 1.5rem;
  }

  .cart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .cart-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: #0f172a;
  }

  .cart-count {
    background: #6366f1;
    color: #ffffff;
    font-size: 0.75rem;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.25rem 0.625rem;
  }

  .cart-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    padding: 3rem 0;
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

  .cart-items {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
    max-height: 400px;
    overflow-y: auto;
  }

  .cart-item {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.75rem;
    background: #f8fafc;
    border-radius: 8px;
  }

  .cart-item-info {
    flex: 1;
  }

  .cart-item-name {
    font-size: 0.875rem;
    font-weight: 500;
    color: #0f172a;
  }

  .cart-item-unit {
    font-size: 0.75rem;
    color: #64748b;
    margin-top: 0.125rem;
  }

  .cart-item-right {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  .cart-item-sub {
    font-size: 0.875rem;
    font-weight: 600;
    color: #0f172a;
  }

  .btn-remove {
    background: none;
    border: none;
    color: #94a3b8;
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 4px;
    transition: all 0.15s;
  }

  .btn-remove:hover {
    color: #ef4444;
    background: #fef2f2;
  }

  .cart-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 0;
    border-top: 1px solid #f1f5f9;
    border-bottom: 1px solid #f1f5f9;
    margin-bottom: 1.5rem;
  }

  .total-value {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 1.25rem;
    font-weight: 700;
    color: #6366f1;
  }

  .cart-notes {
    margin-bottom: 1.5rem;
  }

  .notes-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #0f172a;
    margin-bottom: 0.5rem;
    display: block;
  }

  .checkout-btn {
    width: 100%;
    padding: 1rem;
    font-size: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .checkout-total {
    font-weight: 700;
  }

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

  /* Responsive */
  @media (max-width: 1024px) {
    .pos-layout {
      grid-template-columns: 1fr;
    }

    .pos-cart {
      position: static;
    }
  }
</style>
