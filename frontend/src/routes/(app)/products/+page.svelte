<script lang="ts">
  import { onMount } from 'svelte';
  import { getProducts, createProduct, updateProduct, deleteProduct, getIngredients, createIngredient, updateIngredient, deleteIngredient, updateIngredientStock } from '$lib/api/product';
  import type { Product } from '$lib/types/product';
  import type { Ingredient } from '$lib/types/ingredient';

  let products: Product[] = $state([]);
  let ingredients: Ingredient[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let showLowStockOnly = $state(false);
  let productSearch = $state('');
  let ingredientSearch = $state('');
  let productSortBy = $state<string>('name');
  let productSortOrder = $state<'asc' | 'desc'>('asc');
  let ingredientSortBy = $state<string>('name');
  let ingredientSortOrder = $state<'asc' | 'desc'>('asc');
  let productCurrentPage = $state(1);
  let ingredientCurrentPage = $state(1);
  const itemsPerPage = 20;

  // Modal novo produto
  let showProductModal = $state(false);
  let productEditMode = $state(false);
  let productEditId = $state<number | null>(null);
  let productForm = $state({ Name: '', Description: '', Price: 0, IsComposto: false, Active: true });
  let productSaving = $state(false);
  let productError = $state('');

  // Modal novo ingrediente
  let showIngModal = $state(false);
  let ingEditMode = $state(false);
  let ingEditId = $state<number | null>(null);
  let ingForm = $state({ Name: '', Unit: '', StockQuantity: 0, MinStock: 0 });
  let ingSaving = $state(false);
  let ingError = $state('');

  // Modal ajustar estoque
  let showStockModal = $state(false);
  let stockEditId = $state<number | null>(null);
  let stockForm = $state({ Quantity: 0 });
  let stockSaving = $state(false);
  let stockError = $state('');

  onMount(async () => {
    await loadAll();
  });

  async function loadAll() {
    loading = true;
    error = '';
    try {
      [products, ingredients] = await Promise.all([getProducts(), getIngredients()]);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar dados.';
    } finally {
      loading = false;
    }
  }

  async function saveProduct() {
    productSaving = true;
    productError = '';
    try {
      const payload = {
        name: productForm.Name,
        description: productForm.Description,
        price: Number(productForm.Price),
        is_composto: productForm.IsComposto,
        active: productForm.Active,
      };

      if (productEditMode && productEditId) {
        const updated = await updateProduct(productEditId, payload);
        products = products.map(p => p.ID === productEditId ? updated : p);
      } else {
        const created = await createProduct(payload);
        products = [...products, created];
      }

      showProductModal = false;
      productForm = { Name: '', Description: '', Price: 0, IsComposto: false, Active: true };
      productEditMode = false;
      productEditId = null;
    } catch (e: any) {
      productError = e?.message ?? 'Erro ao salvar produto.';
    } finally {
      productSaving = false;
    }
  }

  function openProductEdit(product: Product) {
    productEditMode = true;
    productEditId = product.ID;
    productForm = {
      Name: product.Name,
      Description: product.Description,
      Price: product.Price,
      IsComposto: product.IsComposto,
      Active: product.Active,
    };
    showProductModal = true;
  }

  function openProductCreate() {
    productEditMode = false;
    productEditId = null;
    productForm = { Name: '', Description: '', Price: 0, IsComposto: false, Active: true };
    showProductModal = true;
  }

  async function deleteProductById(id: number) {
    if (!confirm('Tem certeza que deseja excluir este produto?')) return;
    try {
      await deleteProduct(id);
      products = products.filter(p => p.ID !== id);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao excluir produto.';
    }
  }

  async function saveIngredient() {
    ingSaving = true;
    ingError = '';
    try {
      const payload = {
        name: ingForm.Name,
        unit: ingForm.Unit,
        stock_quantity: Number(ingForm.StockQuantity),
        min_stock: Number(ingForm.MinStock),
      };

      if (ingEditMode && ingEditId) {
        const updated = await updateIngredient(ingEditId, payload);
        ingredients = ingredients.map(i => i.ID === ingEditId ? updated : i);
      } else {
        const created = await createIngredient(payload);
        ingredients = [...ingredients, created];
      }

      showIngModal = false;
      ingForm = { Name: '', Unit: '', StockQuantity: 0, MinStock: 0 };
      ingEditMode = false;
      ingEditId = null;
    } catch (e: any) {
      ingError = e?.message ?? 'Erro ao salvar ingrediente.';
    } finally {
      ingSaving = false;
    }
  }

  function openIngredientEdit(ingredient: Ingredient) {
    ingEditMode = true;
    ingEditId = ingredient.ID;
    ingForm = {
      Name: ingredient.Name,
      Unit: ingredient.Unit,
      StockQuantity: ingredient.StockQuantity,
      MinStock: ingredient.MinStock,
    };
    showIngModal = true;
  }

  function openIngredientCreate() {
    ingEditMode = false;
    ingEditId = null;
    ingForm = { Name: '', Unit: '', StockQuantity: 0, MinStock: 0 };
    showIngModal = true;
  }

  async function deleteIngredientById(id: number) {
    if (!confirm('Tem certeza que deseja excluir este ingrediente?')) return;
    try {
      await deleteIngredient(id);
      ingredients = ingredients.filter(i => i.ID !== id);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao excluir ingrediente.';
    }
  }

  function openStockModal(ingredient: Ingredient) {
    stockEditId = ingredient.ID;
    stockForm = { Quantity: ingredient.StockQuantity };
    showStockModal = true;
  }

  async function saveStock() {
    stockSaving = true;
    stockError = '';
    try {
      if (stockEditId) {
        const updated = await updateIngredientStock(stockEditId, Number(stockForm.Quantity));
        ingredients = ingredients.map(i => i.ID === stockEditId ? updated : i);
      }
      showStockModal = false;
      stockForm = { Quantity: 0 };
      stockEditId = null;
    } catch (e: any) {
      stockError = e?.message ?? 'Erro ao ajustar estoque.';
    } finally {
      stockSaving = false;
    }
  }

  function formatPrice(value: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
  }

  function getStockStatus(ing: Ingredient) {
    if (ing.StockQuantity === 0) {
      return { label: '🔴 Estoque Zerado', class: 'stock-zero' };
    } else if (ing.StockQuantity <= ing.MinStock) {
      return { label: '🟡 Estoque Baixo', class: 'stock-low' };
    }
    return null;
  }

  const filteredIngredients = $derived(
    (showLowStockOnly
      ? ingredients.filter(ing => ing.StockQuantity <= ing.MinStock)
      : ingredients.filter(ing => !ingredientSearch || ing.Name.toLowerCase().includes(ingredientSearch.toLowerCase()))
    ).sort((a, b) => {
      let comparison = 0;
      if (ingredientSortBy === 'name') {
        comparison = a.Name.localeCompare(b.Name);
      } else if (ingredientSortBy === 'stock') {
        comparison = a.StockQuantity - b.StockQuantity;
      }
      return ingredientSortOrder === 'asc' ? comparison : -comparison;
    })
  );

  const filteredProducts = $derived(
    products.filter(p => !productSearch || p.Name.toLowerCase().includes(productSearch.toLowerCase()))
    .sort((a, b) => {
      let comparison = 0;
      if (productSortBy === 'name') {
        comparison = a.Name.localeCompare(b.Name);
      } else if (productSortBy === 'price') {
        comparison = a.Price - b.Price;
      }
      return productSortOrder === 'asc' ? comparison : -comparison;
    })
  );

  const productTotalPages = $derived(Math.ceil(filteredProducts.length / itemsPerPage));
  const ingredientTotalPages = $derived(Math.ceil(filteredIngredients.length / itemsPerPage));
  const paginatedProducts = $derived(
    filteredProducts.slice((productCurrentPage - 1) * itemsPerPage, productCurrentPage * itemsPerPage)
  );
  const paginatedIngredients = $derived(
    filteredIngredients.slice((ingredientCurrentPage - 1) * itemsPerPage, ingredientCurrentPage * itemsPerPage)
  );

  function goToProductPage(page: number) {
    if (page >= 1 && page <= productTotalPages) {
      productCurrentPage = page;
    }
  }

  function goToIngredientPage(page: number) {
    if (page >= 1 && page <= ingredientTotalPages) {
      ingredientCurrentPage = page;
    }
  }
</script>

<div class="page-wrapper">
  <header class="page-header">
    <div>
      <h1 class="page-title">Produtos</h1>
      <p class="page-subtitle">Gerencie produtos e ingredientes do seu restaurante</p>
    </div>
    <div class="header-actions">
      <button class="btn btn-secondary" onclick={openIngredientCreate}>
        + Ingrediente
      </button>
      <button class="btn btn-primary" onclick={openProductCreate}>
        + Produto
      </button>
    </div>
  </header>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <span>Carregando cardápio…</span>
    </div>
  {:else if error}
    <div class="alert alert-error">
      <span>⚠️ {error}</span>
      <button class="btn btn-sm" onclick={loadAll}>Tentar novamente</button>
    </div>
  {:else}
    <!-- Produtos -->
    <section class="section" aria-labelledby="products-heading">
      <div class="section-header">
        <h2 id="products-heading" class="section-title">Produtos <span class="badge">{products.length}</span></h2>
        <div class="section-filters" role="search" aria-label="Filtros de produtos">
          <label for="product-search" class="sr-only">Buscar produtos</label>
          <input
            id="product-search"
            type="text"
            placeholder="Buscar produtos..."
            bind:value={productSearch}
            class="search-input"
            aria-label="Buscar produtos por nome"
          />
          <div class="sort-control">
            <label for="product-sort" class="sr-only">Ordenar produtos</label>
            <select bind:value={productSortBy} id="product-sort" class="sort-select" aria-label="Critério de ordenação de produtos">
              <option value="name">Nome</option>
              <option value="price">Preço</option>
            </select>
            <button
              class="btn btn-sm btn-ghost"
              onclick={() => (productSortOrder = productSortOrder === 'asc' ? 'desc' : 'asc')}
              title={productSortOrder === 'asc' ? 'Crescente' : 'Decrescente'}
              aria-label={`Ordem ${productSortOrder === 'asc' ? 'crescente' : 'decrescente'}, clique para inverter`}
            >
              {productSortOrder === 'asc' ? '↑' : '↓'}
            </button>
          </div>
        </div>
      </div>

      {#if filteredProducts.length === 0}
        <div class="empty-state">
          <p class="empty-icon">🍽️</p>
          <p class="empty-text">{productSearch ? 'Nenhum produto encontrado.' : 'Nenhum produto cadastrado ainda.'}</p>
          <button class="btn btn-primary" onclick={() => (showProductModal = true)}>Criar primeiro produto</button>
        </div>
      {:else}
        <div class="card-grid">
          {#each paginatedProducts as product}
            <div class="product-card">
              <div class="product-card-header">
                <a href="/products/{product.ID}" class="product-name">{product.Name}</a>
                <span class="product-price">{formatPrice(product.Price)}</span>
              </div>
              {#if product.Description}
                <p class="product-desc">{product.Description}</p>
              {/if}
              <div class="product-footer">
                {#if product.IsComposto}
                  <span class="tag">Composto</span>
                {/if}
                <span class="status-dot" class:active={product.Active}>
                  {product.Active ? 'Ativo' : 'Inativo'}
                </span>
              </div>
              <div class="card-actions">
                <button class="btn btn-sm btn-ghost" onclick={() => openProductEdit(product)}>Editar</button>
                <button class="btn btn-sm btn-danger" onclick={() => deleteProductById(product.ID)}>Excluir</button>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      {#if productTotalPages > 1}
        <div class="pagination">
          <button
            class="pagination-btn"
            disabled={productCurrentPage === 1}
            onclick={() => goToProductPage(productCurrentPage - 1)}
          >
            ← Anterior
          </button>
          <span class="pagination-info">
            Página {productCurrentPage} de {productTotalPages} ({filteredProducts.length} produtos)
          </span>
          <button
            class="pagination-btn"
            disabled={productCurrentPage === productTotalPages}
            onclick={() => goToProductPage(productCurrentPage + 1)}
          >
            Próxima →
          </button>
        </div>
      {/if}
    </section>

    <!-- Ingredientes -->
    <section class="section" aria-labelledby="ingredients-heading">
      <div class="section-header">
        <h2 id="ingredients-heading" class="section-title">Ingredientes <span class="badge">{ingredients.length}</span></h2>
        <div class="section-filters" role="search" aria-label="Filtros de ingredientes">
          <label for="ingredient-search" class="sr-only">Buscar ingredientes</label>
          <input
            id="ingredient-search"
            type="text"
            placeholder="Buscar ingredientes..."
            bind:value={ingredientSearch}
            class="search-input"
            aria-label="Buscar ingredientes por nome"
          />
          <div class="sort-control">
            <label for="ingredient-sort" class="sr-only">Ordenar ingredientes</label>
            <select bind:value={ingredientSortBy} id="ingredient-sort" class="sort-select" aria-label="Critério de ordenação de ingredientes">
              <option value="name">Nome</option>
              <option value="stock">Estoque</option>
            </select>
            <button
              class="btn btn-sm btn-ghost"
              onclick={() => (ingredientSortOrder = ingredientSortOrder === 'asc' ? 'desc' : 'asc')}
              title={ingredientSortOrder === 'asc' ? 'Crescente' : 'Decrescente'}
              aria-label={`Ordem ${ingredientSortOrder === 'asc' ? 'crescente' : 'decrescente'}, clique para inverter`}
            >
              {ingredientSortOrder === 'asc' ? '↑' : '↓'}
            </button>
          </div>
          <label class="filter-toggle">
            <input type="checkbox" bind:checked={showLowStockOnly} aria-label="Mostrar apenas ingredientes com estoque baixo" />
            <span>Apenas estoque baixo</span>
          </label>
        </div>
      </div>

      {#if filteredIngredients.length === 0}
        <div class="empty-state">
          <p class="empty-icon">🥦</p>
          <p class="empty-text">{showLowStockOnly ? 'Nenhum ingrediente com estoque baixo.' : (ingredientSearch ? 'Nenhum ingrediente encontrado.' : 'Nenhum ingrediente cadastrado.')}</p>
          <button class="btn btn-secondary" onclick={() => (showIngModal = true)}>Adicionar ingrediente</button>
        </div>
      {:else}
        <div class="ing-table-wrapper">
          <table class="ing-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Nome</th>
                <th>Unidade</th>
                <th>Estoque</th>
                <th>Estoque Mínimo</th>
                <th>Status</th>
                <th>Ações</th>
              </tr>
            </thead>
            <tbody>
              {#each paginatedIngredients as ing}
                {@const stockStatus = getStockStatus(ing)}
                <tr class:stock-zero={stockStatus?.class === 'stock-zero'} class:stock-low={stockStatus?.class === 'stock-low'}>
                  <td class="muted">{ing.ID}</td>
                  <td class="bold">{ing.Name}</td>
                  <td>{ing.Unit}</td>
                  <td>{ing.StockQuantity}</td>
                  <td>{ing.MinStock}</td>
                  <td>
                    {#if stockStatus}
                      <span class="stock-badge {stockStatus.class}">{stockStatus.label}</span>
                    {:else}
                      <span class="stock-badge stock-ok">✅ OK</span>
                    {/if}
                  </td>
                  <td>
                    <button class="btn btn-sm btn-ghost" onclick={() => openIngredientEdit(ing)}>Editar</button>
                    <button class="btn btn-sm btn-ghost" onclick={() => openStockModal(ing)}>Ajustar</button>
                    <button class="btn btn-sm btn-danger" onclick={() => deleteIngredientById(ing.ID)}>Excluir</button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}

      {#if ingredientTotalPages > 1}
        <div class="pagination">
          <button
            class="pagination-btn"
            disabled={ingredientCurrentPage === 1}
            onclick={() => goToIngredientPage(ingredientCurrentPage - 1)}
          >
            ← Anterior
          </button>
          <span class="pagination-info">
            Página {ingredientCurrentPage} de {ingredientTotalPages} ({filteredIngredients.length} ingredientes)
          </span>
          <button
            class="pagination-btn"
            disabled={ingredientCurrentPage === ingredientTotalPages}
            onclick={() => goToIngredientPage(ingredientCurrentPage + 1)}
          >
            Próxima →
          </button>
        </div>
      {/if}
    </section>
  {/if}
</div>

<!-- Modal: Novo Produto -->
{#if showProductModal}
  <div class="modal-overlay" onclick={() => (showProductModal = false)}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <h2 class="modal-title">{productEditMode ? 'Editar Produto' : 'Novo Produto'}</h2>

      {#if productError}
        <p class="form-error">{productError}</p>
      {/if}

      <div class="form-group">
        <label for="p-name">Nome *</label>
        <input id="p-name" type="text" bind:value={productForm.Name} placeholder="Ex: Feijoada Completa" />
      </div>
      <div class="form-group">
        <label for="p-desc">Descrição</label>
        <textarea id="p-desc" bind:value={productForm.Description} rows="2" placeholder="Descrição opcional"></textarea>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label for="p-price">Preço (R$) *</label>
          <input id="p-price" type="number" min="0" step="0.01" bind:value={productForm.Price} />
        </div>
      </div>
      <div class="form-group checkbox-group">
        <label>
          <input type="checkbox" bind:checked={productForm.IsComposto} />
          Produto composto (requer ficha técnica)
        </label>
      </div>
      <div class="form-group checkbox-group">
        <label>
          <input type="checkbox" bind:checked={productForm.Active} />
          Produto ativo (visível no cardápio)
        </label>
      </div>

      <div class="modal-actions">
        <button class="btn btn-ghost" onclick={() => (showProductModal = false)}>Cancelar</button>
        <button class="btn btn-primary" onclick={saveProduct} disabled={productSaving || !productForm.Name}>
          {productSaving ? 'Salvando…' : (productEditMode ? 'Atualizar Produto' : 'Criar Produto')}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal: Novo Ingrediente -->
{#if showIngModal}
  <div class="modal-overlay" onclick={() => (showIngModal = false)}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <h2 class="modal-title">{ingEditMode ? 'Editar Ingrediente' : 'Novo Ingrediente'}</h2>

      {#if ingError}
        <p class="form-error">{ingError}</p>
      {/if}

      <div class="form-group">
        <label for="i-name">Nome *</label>
        <input id="i-name" type="text" bind:value={ingForm.Name} placeholder="Ex: Feijão Preto" />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label for="i-unit">Unidade *</label>
          <input id="i-unit" type="text" bind:value={ingForm.Unit} placeholder="kg, g, L, un…" />
        </div>
        <div class="form-group">
          <label for="i-stock">Estoque inicial</label>
          <input id="i-stock" type="number" min="0" step="0.01" bind:value={ingForm.StockQuantity} />
        </div>
      </div>
      <div class="form-group">
        <label for="i-minstock">Estoque mínimo</label>
        <input id="i-minstock" type="number" min="0" step="0.01" bind:value={ingForm.MinStock} />
      </div>

      <div class="modal-actions">
        <button class="btn btn-ghost" onclick={() => (showIngModal = false)}>Cancelar</button>
        <button class="btn btn-primary" onclick={saveIngredient} disabled={ingSaving || !ingForm.Name || !ingForm.Unit}>
          {ingSaving ? 'Salvando…' : (ingEditMode ? 'Atualizar Ingrediente' : 'Criar Ingrediente')}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal: Ajustar Estoque -->
{#if showStockModal}
  <div class="modal-overlay" onclick={() => (showStockModal = false)}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <h2 class="modal-title">Ajustar Estoque</h2>

      {#if stockError}
        <p class="form-error">{stockError}</p>
      {/if}

      <div class="form-group">
        <label for="s-quantity">Quantidade *</label>
        <input id="s-quantity" type="number" min="0" step="0.01" bind:value={stockForm.Quantity} />
      </div>

      <div class="modal-actions">
        <button class="btn btn-ghost" onclick={() => (showStockModal = false)}>Cancelar</button>
        <button class="btn btn-primary" onclick={saveStock} disabled={stockSaving || stockForm.Quantity < 0}>
          {stockSaving ? 'Salvando…' : 'Ajustar Estoque'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Acessibilidade */
  .sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0; }

  .page-wrapper   { max-width: 1100px; margin: 0 auto; padding: 2rem 1.5rem; }
  .page-header    { display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 1rem; margin-bottom: 2.5rem; }
  .page-title     { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .page-subtitle  { font-size: 0.9rem; color: var(--color-muted); margin: 0.25rem 0 0; }
  .header-actions { display: flex; gap: 0.75rem; }

  .section        { margin-bottom: 3rem; }
  .section-title  { font-size: 1.1rem; font-weight: 600; color: var(--color-text); margin: 0 0 1.25rem; display: flex; align-items: center; gap: 0.5rem; }
  .badge          { background: var(--color-surface-2, #f0f0f0); color: var(--color-muted); font-size: 0.75rem; font-weight: 600; padding: 0.1rem 0.5rem; border-radius: 99px; }

  /* Cards de produto */
  .card-grid      { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 1rem; }
  .product-card   { background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1.25rem; transition: box-shadow 0.15s, border-color 0.15s; }
  .product-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); border-color: var(--color-primary, #e85d04); }
  .product-card-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 0.5rem; margin-bottom: 0.5rem; }
  .product-name   { font-weight: 600; font-size: 1rem; color: var(--color-text); text-decoration: none; }
  .product-name:hover { color: var(--color-primary, #e85d04); text-decoration: underline; }
  .product-price  { font-weight: 700; color: var(--color-primary, #e85d04); white-space: nowrap; }
  .product-desc   { font-size: 0.85rem; color: var(--color-muted); margin: 0 0 0.75rem; line-height: 1.45; }
  .product-footer { display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
  .card-actions   { display: flex; gap: 0.5rem; margin-top: 0.75rem; padding-top: 0.75rem; border-top: 1px solid var(--color-border, #e5e7eb); }
  .tag            { font-size: 0.75rem; background: var(--color-surface-2, #f5f5f5); padding: 0.15rem 0.5rem; border-radius: 4px; color: var(--color-muted); }
  .status-dot     { font-size: 0.75rem; color: #9ca3af; margin-left: auto; }
  .status-dot.active { color: #16a34a; }

  /* Tabela ingredientes */
  .ing-table-wrapper { overflow-x: auto; border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; }
  .ing-table      { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
  .ing-table th   { text-align: left; padding: 0.75rem 1rem; background: var(--color-surface-2, #f9fafb); font-weight: 600; color: var(--color-muted); font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.04em; }
  .ing-table td   { padding: 0.75rem 1rem; border-top: 1px solid var(--color-border, #e5e7eb); }
  .ing-table tr:hover td { background: var(--color-surface-2, #f9fafb); }
  .muted          { color: var(--color-muted); }
  .bold           { font-weight: 500; }

  /* Estados */
  .loading-state  { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: var(--color-muted); }
  .spinner        { width: 2rem; height: 2rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin  { to { transform: rotate(360deg); } }
  .empty-state    { text-align: center; padding: 3rem 1rem; }
  .empty-icon     { font-size: 2.5rem; margin: 0; }
  .empty-text     { color: var(--color-muted); margin: 0.5rem 0 1.25rem; }
  .alert-error    { display: flex; align-items: center; gap: 1rem; background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 1rem 1.25rem; border-radius: 0.5rem; margin-bottom: 1.5rem; }

  /* Section header with filter */
  .section-header { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; }
  .section-filters { display: flex; gap: 0.75rem; align-items: center; flex-wrap: wrap; }
  .filter-toggle { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; color: var(--color-text); cursor: pointer; user-select: none; white-space: nowrap; padding: 0.4rem 0.6rem; border-radius: 0.4rem; transition: background 0.15s; }
  .filter-toggle:hover { background: var(--color-surface-2, #f5f5f5); }
  .filter-toggle input { cursor: pointer; }
  .search-input { padding: 0.5rem 0.75rem; border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; font-size: 0.85rem; background: var(--color-surface, #fff); color: var(--color-text); transition: all 0.2s ease; min-width: 180px; box-shadow: 0 1px 2px rgba(0,0,0,0.05); }
  .search-input:focus { outline: none; border-color: var(--color-primary, #e85d04); box-shadow: 0 0 0 3px rgba(232,93,4,0.1); }
  .sort-control { display: flex; gap: 0.25rem; }
  .sort-select { padding: 0.5rem 0.75rem; border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; font-size: 0.85rem; background: var(--color-surface, #fff); color: var(--color-text); transition: all 0.2s ease; cursor: pointer; box-shadow: 0 1px 2px rgba(0,0,0,0.05); }
  .sort-select:focus { outline: none; border-color: var(--color-primary, #e85d04); box-shadow: 0 0 0 3px rgba(232,93,4,0.1); }

  /* Paginação */
  .pagination { display: flex; justify-content: center; align-items: center; gap: 1rem; margin-top: 1.5rem; padding: 0.75rem; background: var(--color-surface-2, #f9fafb); border-radius: 0.4rem; }
  .pagination-btn { padding: 0.4rem 0.85rem; border: 1px solid var(--color-border, #d1d5db); background: var(--color-surface, #fff); color: var(--color-text); border-radius: 0.35rem; font-size: 0.8rem; font-weight: 500; cursor: pointer; transition: all 0.15s; }
  .pagination-btn:hover:not(:disabled) { background: var(--color-primary, #e85d04); color: #fff; border-color: var(--color-primary, #e85d04); }
  .pagination-btn:disabled { opacity: 0.4; cursor: not-allowed; }
  .pagination-info { font-size: 0.8rem; color: var(--color-muted); }

  /* Stock status badges */
  .stock-badge    { display: inline-flex; align-items: center; padding: 0.25rem 0.6rem; border-radius: 99px; font-size: 0.75rem; font-weight: 600; }
  .stock-ok       { background: #dcfce7; color: #15803d; }
  .stock-low      { background: #fef9c3; color: #854d0e; }
  .stock-zero     { background: #fee2e2; color: #b91c1c; }
  .ing-table tr.stock-zero { background: #fef2f2; }
  .ing-table tr.stock-low { background: #fffbeb; }

  /* Modal */
  .modal-overlay  { position: fixed; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; z-index: 50; padding: 1rem; }
  .modal          { background: var(--color-surface, #fff); border-radius: 1rem; padding: 2rem; width: 100%; max-width: 480px; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }
  .modal-title    { font-size: 1.25rem; font-weight: 700; margin: 0 0 1.5rem; }
  .modal-actions  { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.75rem; }

  /* Formulário */
  .form-group     { display: flex; flex-direction: column; gap: 0.4rem; margin-bottom: 1rem; }
  .form-group label { font-size: 0.85rem; font-weight: 500; color: var(--color-text); }
  .form-group input,
  .form-group textarea { border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; padding: 0.6rem 0.75rem; font-size: 0.9rem; background: var(--color-surface, #fff); color: var(--color-text); transition: border-color 0.15s; width: 100%; box-sizing: border-box; }
  .form-group input:focus,
  .form-group textarea:focus { outline: none; border-color: var(--color-primary, #e85d04); box-shadow: 0 0 0 3px rgba(232,93,4,0.12); }
  .form-row       { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; }
  .checkbox-group label { flex-direction: row; align-items: center; gap: 0.5rem; cursor: pointer; font-weight: 400; }
  .form-error     { color: #b91c1c; font-size: 0.85rem; margin-bottom: 1rem; background: #fef2f2; padding: 0.6rem 0.75rem; border-radius: 0.4rem; border: 1px solid #fca5a5; }

  /* Botões */
  .btn            { display: inline-flex; align-items: center; justify-content: center; gap: 0.4rem; padding: 0.55rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s, opacity 0.15s; text-decoration: none; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-secondary  { background: var(--color-surface-2, #f3f4f6); color: var(--color-text); border: 1px solid var(--color-border, #d1d5db); }
  .btn-secondary:hover:not(:disabled) { background: var(--color-border, #e5e7eb); }
  .btn-ghost      { background: transparent; color: var(--color-muted); }
  .btn-ghost:hover:not(:disabled) { background: var(--color-surface-2, #f3f4f6); }
  .btn-danger     { background: #fee2e2; color: #b91c1c; border: 1px solid #fca5a5; }
  .btn-danger:hover:not(:disabled) { background: #fecaca; }
  .btn-sm         { padding: 0.35rem 0.75rem; font-size: 0.8rem; }

  @media (max-width: 600px) {
    .page-header  { flex-direction: column; }
    .form-row     { grid-template-columns: 1fr; }
    .header-actions { width: 100%; }
    .header-actions .btn { flex: 1; }
  }
</style>
