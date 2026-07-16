<script lang="ts">
  import { onMount } from 'svelte';
  import { getProducts, createProduct, updateProduct, deleteProduct } from '$lib/api/product';
  import type { Product } from '$lib/types/product';
  import { Card, Button, Input, Textarea, Select, Checkbox, Badge, Alert, Loading, EmptyState, Modal, Table, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Search, Plus, Package, AlertTriangle, TrendingUp, TrendingDown, Activity, MoreHorizontal, Filter, ArrowUpDown } from '@lucide/svelte';

  let products: Product[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let productSearch = $state('');
  let productSortBy = $state<string>('name');
  let productSortOrder = $state<'asc' | 'desc'>('asc');
  let productCurrentPage = $state(1);
  const itemsPerPage = 12;

  // Modal novo produto
  let showProductModal = $state(false);
  let productEditMode = $state(false);
  let productEditId = $state<number | null>(null);
  let productForm = $state({ Name: '', Description: '', Price: 0, IsComposto: false, Active: true });
  let productSaving = $state(false);
  let productError = $state('');


  onMount(async () => {
    await loadAll();
  });

  async function loadAll() {
    loading = true;
    error = '';
    try {
      products = await getProducts();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar produtos.';
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
      Description: product.Description ?? '',
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


  function formatPrice(value: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
  }


  function getProductLowStockStatus(product: Product) {
    if (!product.Ingredients || product.Ingredients.length === 0) return null;
    const lowStockIngredients = product.Ingredients.filter((i: any) => i.stock <= (i.min_stock || 5));
    if (lowStockIngredients.length === 0) return null;
    return { count: lowStockIngredients.length, total: product.Ingredients.length };
  }


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
  const paginatedProducts = $derived(
    filteredProducts.slice((productCurrentPage - 1) * itemsPerPage, productCurrentPage * itemsPerPage)
  );

  function goToProductPage(page: number) {
    if (page >= 1 && page <= productTotalPages) {
      productCurrentPage = page;
    }
  }

</script>

<Workspace
  breadcrumb={[{ label: 'Produtos' }]}
  title="Gerenciamento de Produtos"
  description="Catálogo de produtos"
>
  <svelte:fragment slot="actions">
    <Button onclick={openProductCreate} variant="primary">
      <Plus size={16} />
      Novo Produto
    </Button>
  </svelte:fragment>

  {#if loading}
    <div class="skeleton-grid">
      {#each Array(6) as _}
        <Card class="skeleton-card">
          <div class="skeleton-card-top">
            <Skeleton variant="rectangular" width="48px" height="48px" />
            <div class="skeleton-content">
              <Skeleton variant="text" width="120px" height="16px" />
              <Skeleton variant="text" width="80px" height="12px" />
            </div>
          </div>
          <div class="skeleton-card-bottom">
            <Skeleton variant="text" width="60px" height="14px" />
          </div>
        </Card>
      {/each}
    </div>
  {:else if error}
    <Alert variant="error" dismissible onDismiss={() => error = ''}>
      ⚠️ {error}
      <Button onclick={loadAll} size="sm">Tentar novamente</Button>
    </Alert>
  {:else}
    <!-- Filtros Rápidos -->
    <Card class="filters-card">
      <div class="filters-header">
        <div class="filters-title">
          <Filter size={18} />
          <span>Filtros</span>
        </div>
        <div class="filters-stats">
          <span class="stat-item">{products.length} produtos</span>
        </div>
      </div>
      <div class="filters-grid">
        <div class="filter-group">
          <label class="filter-label">Buscar produtos</label>
          <div class="search-wrapper">
            <Search size={16} class="search-icon" />
            <Input
              placeholder="Buscar por nome..."
              bind:value={productSearch}
              class="search-input"
            />
          </div>
        </div>
        <div class="filter-group">
          <label class="filter-label">Ordenar por</label>
          <div class="sort-wrapper">
            <Select bind:value={productSortBy} class="sort-select">
              <option value="name">Nome</option>
              <option value="price">Preço</option>
            </Select>
            <Button
              variant="ghost"
              size="sm"
              onclick={() => (productSortOrder = productSortOrder === 'asc' ? 'desc' : 'asc')}
              class="sort-toggle"
            >
              <ArrowUpDown size={16} />
            </Button>
          </div>
        </div>
      </div>
    </Card>

    <!-- Produtos -->
    <div class="section-header">
      <h2 class="section-title">Produtos</h2>
      <div class="section-actions">
        <Button onclick={openProductCreate} variant="secondary" size="sm">
          <Plus size={14} />
          Adicionar
        </Button>
      </div>
    </div>

    {#if filteredProducts.length === 0}
      <div class="empty-state">
        <Package size={48} class="empty-icon" />
        <span class="empty-title">{productSearch ? 'Nenhum produto encontrado' : 'Nenhum produto cadastrado'}</span>
        <span class="empty-subtitle">{productSearch ? 'Tente outro termo de busca' : 'Comece adicionando seu primeiro produto'}</span>
        {#if !productSearch}
          <Button onclick={openProductCreate} variant="primary">Criar Produto</Button>
        {/if}
      </div>
    {:else}
      <div class="products-grid">
        {#each paginatedProducts as product}
          {@const lowStock = getProductLowStockStatus(product)}
          <Card class="product-card">
            <div class="product-card-top">
              <div class="product-info">
                <a href="/products/{product.ID}" class="product-name">{product.Name}</a>
                {#if product.Description}
                  <p class="product-desc">{product.Description}</p>
                {/if}
              </div>
              <div class="product-price">{formatPrice(product.Price)}</div>
            </div>
            
            <div class="product-meta">
              {#if product.IsComposto}
                <Badge variant="primary" size="sm">Composto</Badge>
              {/if}
              <Badge variant={product.Active ? 'success' : 'default'} size="sm">
                {product.Active ? 'Ativo' : 'Inativo'}
              </Badge>
              {#if lowStock}
                <Badge variant="warning" size="sm" icon>
                  <AlertTriangle size={12} />
                  {lowStock.count}/{lowStock.total} ingredientes baixos
                </Badge>
              {/if}
            </div>

            <div class="product-actions">
              <Button variant="ghost" size="sm" onclick={() => openProductEdit(product)}>
                Editar
              </Button>
              <Button variant="ghost" size="sm" onclick={() => deleteProductById(product.ID)} class="danger">
                Excluir
              </Button>
            </div>
          </Card>
        {/each}
      </div>

      {#if productTotalPages > 1}
        <div class="pagination">
          <Button
            variant="ghost"
            size="sm"
            disabled={productCurrentPage === 1}
            onclick={() => goToProductPage(productCurrentPage - 1)}
          >
            Anterior
          </Button>
          <span class="pagination-info">
            Página {productCurrentPage} de {productTotalPages}
          </span>
          <Button
            variant="ghost"
            size="sm"
            disabled={productCurrentPage === productTotalPages}
            onclick={() => goToProductPage(productCurrentPage + 1)}
          >
            Próxima
          </Button>
        </div>
      {/if}
    {/if}

  {/if}
</Workspace>

<!-- Modal: Novo Produto -->
<Modal 
  open={showProductModal} 
  title={productEditMode ? 'Editar Produto' : 'Novo Produto'}
  onClose={() => (showProductModal = false)}
>
  {#if productError}
    <Alert variant="error">{productError}</Alert>
  {/if}

  <Input
    id="p-name"
    label="Nome *"
    bind:value={productForm.Name}
    placeholder="Ex: Feijoada Completa"
  />
  <Textarea
    id="p-desc"
    label="Descrição"
    bind:value={productForm.Description}
    placeholder="Descrição opcional"
    rows={2}
  />
  <Input
    id="p-price"
    label="Preço (R$) *"
    type="number"
    min="0"
    step="0.01"
    bind:value={productForm.Price}
  />
  <Checkbox
    label="Produto composto (requer ficha técnica)"
    bind:checked={productForm.IsComposto}
  />
  <Checkbox
    label="Produto ativo (visível no cardápio)"
    bind:checked={productForm.Active}
  />

  <div class="modal-actions">
    <Button variant="ghost" onclick={() => (showProductModal = false)}>Cancelar</Button>
    <Button 
      variant="primary" 
      onclick={saveProduct} 
      disabled={productSaving || !productForm.Name}
      loading={productSaving}
    >
      {productEditMode ? 'Atualizar Produto' : 'Criar Produto'}
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

  .filters-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1rem;
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .filter-label {
    font-size: 0.75rem;
    font-weight: 500;
    color: #64748b;
  }

  .search-wrapper {
    position: relative;
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

  .sort-wrapper {
    display: flex;
    gap: 0.5rem;
  }

  .sort-select {
    flex: 1;
  }

  .sort-toggle {
    padding: 0.5rem;
  }

  /* Section Header */
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #0f172a;
    margin: 0;
    letter-spacing: -0.025em;
  }

  .section-actions {
    display: flex;
    gap: 0.5rem;
  }

  /* Products Grid */
  .products-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
  }

  .product-card {
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .product-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px 0 rgb(0 0 0 / 0.08);
  }

  .product-card-top {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .product-info {
    flex: 1;
  }

  .product-name {
    font-weight: 600;
    font-size: 1rem;
    color: #0f172a;
    text-decoration: none;
    display: block;
    margin-bottom: 0.5rem;
  }

  .product-name:hover {
    color: #6366f1;
  }

  .product-desc {
    font-size: 0.875rem;
    color: #64748b;
    margin: 0;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .product-price {
    font-weight: 700;
    color: #6366f1;
    font-size: 1.25rem;
    letter-spacing: -0.025em;
  }

  .product-meta {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .product-actions {
    display: flex;
    gap: 0.5rem;
    padding-top: 1rem;
    border-top: 1px solid #f1f5f9;
  }

  .product-actions .danger {
    color: #ef4444;
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

  /* Modal Actions */
  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1.75rem;
  }

  /* Skeleton Loading */
  .skeleton-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
  }

  .skeleton-card {
    padding: 1.5rem;
  }

  .skeleton-card-top {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .skeleton-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .skeleton-card-bottom {
    margin-top: 1rem;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .filters-grid {
      grid-template-columns: 1fr;
    }

    .products-grid {
      grid-template-columns: 1fr;
    }


    .section-header {
      flex-direction: column;
      align-items: stretch;
      gap: 1rem;
    }

    .section-actions {
      justify-content: stretch;
    }

    .section-actions button {
      flex: 1;
    }
  }
</style>
