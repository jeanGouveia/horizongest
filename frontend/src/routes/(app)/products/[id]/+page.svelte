<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { getProduct, getProductIngredients, updateProductIngredients, getIngredients } from '$lib/api/product';
  import type { Product } from '$lib/types/product';
  import type { Ingredient } from '$lib/types/ingredient';

  const productId = Number($page.params.id);

  let product: Product | null = $state(null);
  let allIngredients: Ingredient[] = $state([]);
  let productIngredients: (Ingredient & { quantity: number })[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let saving = $state(false);
  let saveSuccess = $state(false);
  let saveError = $state('');

  // Estado de edição de ingredientes
  let selectedIngId = $state('');
  let selectedQty = $state(1);

  onMount(async () => {
    loading = true;
    error = '';
    try {
      const [prod, ings, prodIngs] = await Promise.all([
        getProduct(productId),
        getIngredients(),
        getProductIngredients(productId),
      ]);
      product = prod;
      allIngredients = ings;
      productIngredients = (prodIngs as any[]).map((i) => ({ ...i, quantity: i.Quantity ?? 1 }));
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar produto.';
    } finally {
      loading = false;
    }
  });

  function addIngredient() {
    const id = Number(selectedIngId);
    if (!id) return;
    if (productIngredients.find((i) => i.ID === id)) return;
    const ing = allIngredients.find((i) => i.ID === id);
    if (!ing) return;
    productIngredients = [...productIngredients, { ...ing, quantity: Number(selectedQty) }];
    selectedIngId = '';
    selectedQty = 1;
  }

  function removeIngredient(id: number) {
    productIngredients = productIngredients.filter((i) => i.ID !== id);
  }

  function updateQty(id: number, qty: number) {
    productIngredients = productIngredients.map((i) => (i.ID === id ? { ...i, quantity: qty } : i));
  }

  async function saveIngredients() {
    saving = true;
    saveSuccess = false;
    saveError = '';
    try {
      await updateProductIngredients(
        productId,
        productIngredients.map((i) => ({ ingredient_id: i.ID, quantity: i.quantity }))
      );
      saveSuccess = true;
      setTimeout(() => (saveSuccess = false), 3000);
    } catch (e: any) {
      saveError = e?.message ?? 'Erro ao salvar ingredientes.';
    } finally {
      saving = false;
    }
  }

  function formatPrice(v: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(v);
  }

  const availableToAdd = $derived(
    allIngredients.filter((i) => !productIngredients.find((pi) => pi.ID === i.ID))
  );
</script>

<div class="page-wrapper">
  <a href="/products" class="back-link">← Voltar para Cardápio</a>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <span>Carregando produto…</span>
    </div>
  {:else if error}
    <div class="alert alert-error">{error}</div>
  {:else if product}
    <header class="product-header">
      <div class="product-meta">
        <h1 class="product-title">{product.Name}</h1>
        {#if product.IsComposto}
          <span class="tag">Composto</span>
        {/if}
        <span class="status-badge" class:active={product.Active}>
          {product.Active ? 'Ativo' : 'Inativo'}
        </span>
      </div>
      <div class="product-price">{formatPrice(product.Price)}</div>
    </header>

    {#if product.Description}
      <p class="product-desc">{product.Description}</p>
    {/if}

    <div class="divider"></div>

    {#if product.IsComposto}
      <!-- Ingredientes -->
      <section class="section">
        <h2 class="section-title">Composição do Prato</h2>

        <!-- Lista atual -->
        {#if productIngredients.length > 0}
          <div class="ing-list">
            {#each productIngredients as ing}
              <div class="ing-row">
                <span class="ing-name">{ing.Name}</span>
                <div class="ing-controls">
                  <input
                    type="number"
                    min="0.01"
                    step="0.01"
                    value={ing.quantity}
                    oninput={(e) => updateQty(ing.ID, Number((e.target as HTMLInputElement).value))}
                    class="qty-input"
                  />
                  <span class="ing-unit">{ing.Unit}</span>
                  <button class="btn-remove" onclick={() => removeIngredient(ing.ID)} title="Remover">✕</button>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="empty-note">Nenhum ingrediente vinculado ainda. Adicione abaixo.</p>
        {/if}

        <!-- Adicionar ingrediente -->
        <div class="add-ing-form">
          <select bind:value={selectedIngId} class="select-input">
            <option value="">Selecione um ingrediente…</option>
            {#each availableToAdd as ing}
              <option value={ing.ID}>{ing.Name} ({ing.Unit})</option>
            {/each}
          </select>
          <input
            type="number"
            min="0.01"
            step="0.01"
            bind:value={selectedQty}
            class="qty-input"
            placeholder="Qtd"
          />
          <button class="btn btn-secondary" onclick={addIngredient} disabled={!selectedIngId}>
            Adicionar
          </button>
        </div>

        <!-- Salvar -->
        <div class="save-row">
          {#if saveError}
            <p class="form-error">{saveError}</p>
          {/if}
          {#if saveSuccess}
            <p class="form-success">✓ Ingredientes atualizados!</p>
          {/if}
          <button class="btn btn-primary" onclick={saveIngredients} disabled={saving}>
            {saving ? 'Salvando…' : 'Salvar Composição'}
          </button>
        </div>
      </section>
    {:else}
      <div class="info-box">
        <p>Este é um produto simples e não possui ficha técnica.</p>
      </div>
    {/if}
  {/if}
</div>

<style>
  .page-wrapper   { max-width: 780px; margin: 0 auto; padding: 2rem 1.5rem; }
  .back-link      { display: inline-flex; align-items: center; gap: 0.3rem; color: var(--color-muted); font-size: 0.875rem; text-decoration: none; margin-bottom: 1.75rem; }
  .back-link:hover { color: var(--color-text); }

  .product-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 1rem; flex-wrap: wrap; margin-bottom: 0.75rem; }
  .product-meta   { display: flex; align-items: center; gap: 0.6rem; flex-wrap: wrap; }
  .product-title  { font-size: 1.6rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .product-price  { font-size: 1.5rem; font-weight: 700; color: var(--color-primary, #e85d04); }
  .product-desc   { color: var(--color-muted); font-size: 0.95rem; margin: 0 0 1.5rem; }
  .tag            { font-size: 0.75rem; background: var(--color-surface-2, #f5f5f5); padding: 0.2rem 0.55rem; border-radius: 4px; color: var(--color-muted); }
  .status-badge   { font-size: 0.75rem; font-weight: 600; padding: 0.2rem 0.55rem; border-radius: 4px; background: #fee2e2; color: #b91c1c; }
  .status-badge.active { background: #dcfce7; color: #15803d; }

  .divider        { border: none; border-top: 1px solid var(--color-border, #e5e7eb); margin: 1.5rem 0; }
  .section        { margin-bottom: 2rem; }
  .section-title  { font-size: 1.1rem; font-weight: 600; margin: 0 0 1.25rem; }

  .ing-list       { display: flex; flex-direction: column; gap: 0.5rem; margin-bottom: 1.25rem; }
  .ing-row        { display: flex; align-items: center; gap: 0.75rem; background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.5rem; padding: 0.6rem 0.9rem; }
  .ing-name       { flex: 1; font-weight: 500; }
  .ing-controls   { display: flex; align-items: center; gap: 0.5rem; }
  .ing-unit       { color: var(--color-muted); font-size: 0.85rem; min-width: 28px; }
  .qty-input      { width: 72px; border: 1px solid var(--color-border, #d1d5db); border-radius: 0.4rem; padding: 0.35rem 0.5rem; font-size: 0.85rem; text-align: right; background: var(--color-surface, #fff); color: var(--color-text); }
  .qty-input:focus { outline: none; border-color: var(--color-primary, #e85d04); }
  .btn-remove     { background: none; border: none; color: #ef4444; cursor: pointer; font-size: 0.9rem; padding: 0.2rem 0.3rem; border-radius: 4px; line-height: 1; }
  .btn-remove:hover { background: #fee2e2; }

  .empty-note     { color: var(--color-muted); font-size: 0.9rem; margin-bottom: 1.25rem; }

  .add-ing-form   { display: flex; gap: 0.5rem; align-items: center; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .select-input   { flex: 1; min-width: 180px; border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; padding: 0.6rem 0.75rem; font-size: 0.9rem; background: var(--color-surface, #fff); color: var(--color-text); }
  .select-input:focus { outline: none; border-color: var(--color-primary, #e85d04); }

  .save-row       { display: flex; align-items: center; justify-content: flex-end; gap: 1rem; flex-wrap: wrap; }
  .form-error     { color: #b91c1c; font-size: 0.85rem; background: #fef2f2; padding: 0.5rem 0.75rem; border-radius: 0.4rem; border: 1px solid #fca5a5; margin: 0; }
  .form-success   { color: #15803d; font-size: 0.85rem; background: #dcfce7; padding: 0.5rem 0.75rem; border-radius: 0.4rem; margin: 0; }

  .info-box       { background: var(--color-surface-2, #f9fafb); padding: 1rem 1.25rem; border-radius: 0.5rem; color: var(--color-muted); font-size: 0.9rem; }

  .loading-state  { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: var(--color-muted); }
  .spinner        { width: 2rem; height: 2rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin  { to { transform: rotate(360deg); } }
  .alert-error    { background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 1rem 1.25rem; border-radius: 0.5rem; }

  .btn            { display: inline-flex; align-items: center; padding: 0.55rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-secondary  { background: var(--color-surface-2, #f3f4f6); color: var(--color-text); border: 1px solid var(--color-border, #d1d5db); }
  .btn-secondary:hover:not(:disabled) { background: var(--color-border, #e5e7eb); }
</style>
