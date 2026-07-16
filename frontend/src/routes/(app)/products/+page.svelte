<script lang="ts">
  import { onMount } from 'svelte';
  import { getProducts, createProduct, updateProduct, deleteProduct } from '$lib/api/product';
  import type { Product } from '$lib/types/product';
  import { Card, Button, Input, Textarea, Select, Checkbox, Badge, Alert, Loading, EmptyState, Modal, Table, Skeleton, ProductCard } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Search, Plus, Package, AlertTriangle, TrendingUp, TrendingDown, Activity, MoreHorizontal, Filter, ArrowUpDown } from '@lucide/svelte';

  let products: Product[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let productSearch = $state('');
  let productSortBy = $state<string>('name');
  let productSortOrder = $state<'asc' | 'desc'>('asc');
  let productCurrentPage = $state(1);
  let productFilter = $state<'all' | 'active' | 'archived' | 'promotion' | 'new' | 'featured' | 'composto'>('all');
  const itemsPerPage = 12;

  // Modal novo produto
  let showProductModal = $state(false);
  let productEditMode = $state(false);
  let productEditId = $state<number | null>(null);
  let productForm = $state({
    Name: '', Description: '', Price: 0, IsComposto: false, Active: true,
    PhotoURL: '', CategoryID: undefined as number | undefined, DisplayOrder: 0, PreparationTimeMinutes: 0,
    Featured: false, IsNew: false, PromotionPrice: undefined as number | undefined,
    PromotionStart: undefined as string | undefined, PromotionEnd: undefined as string | undefined,
    AvailableFrom: '', AvailableUntil: '', SKU: '', InternalNotes: ''
  });
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
        photo_url: productForm.PhotoURL,
        category_id: productForm.CategoryID,
        display_order: Number(productForm.DisplayOrder),
        preparation_time_minutes: Number(productForm.PreparationTimeMinutes),
        featured: productForm.Featured,
        is_new: productForm.IsNew,
        promotion_price: productForm.PromotionPrice ? Number(productForm.PromotionPrice) : undefined,
        promotion_start: productForm.PromotionStart,
        promotion_end: productForm.PromotionEnd,
        available_from: productForm.AvailableFrom,
        available_until: productForm.AvailableUntil,
        sku: productForm.SKU,
        internal_notes: productForm.InternalNotes,
      };

      if (productEditMode && productEditId) {
        const updated = await updateProduct(productEditId, payload);
        products = products.map(p => p.ID === productEditId ? updated : p);
      } else {
        const created = await createProduct(payload);
        products = [...products, created];
      }

      showProductModal = false;
      productForm = {
        Name: '', Description: '', Price: 0, IsComposto: false, Active: true,
        PhotoURL: '', CategoryID: undefined, DisplayOrder: 0, PreparationTimeMinutes: 0,
        Featured: false, IsNew: false, PromotionPrice: undefined, PromotionStart: undefined,
        PromotionEnd: undefined, AvailableFrom: '', AvailableUntil: '', SKU: '', InternalNotes: ''
      };
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
      PhotoURL: product.PhotoURL ?? '',
      CategoryID: product.CategoryID,
      DisplayOrder: product.DisplayOrder,
      PreparationTimeMinutes: product.PreparationTimeMinutes,
      Featured: product.Featured,
      IsNew: product.IsNew,
      PromotionPrice: product.PromotionPrice,
      PromotionStart: product.PromotionStart,
      PromotionEnd: product.PromotionEnd,
      AvailableFrom: product.AvailableFrom ?? '',
      AvailableUntil: product.AvailableUntil ?? '',
      SKU: product.SKU ?? '',
      InternalNotes: product.InternalNotes ?? '',
    };
    showProductModal = true;
  }

  function openProductCreate() {
    window.location.href = '/products/new';
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

  async function duplicateProduct(id: number) {
    if (!confirm('Tem certeza que deseja duplicar este produto?')) return;
    try {
      const product = products.find(p => p.ID === id);
      if (!product) return;
      
      const payload = {
        name: product.Name + ' (Cópia)',
        description: product.Description,
        price: product.Price,
        is_composto: product.IsComposto,
        active: false, // Cópia começa inativa
        photo_url: product.PhotoURL,
        category_id: product.CategoryID,
        display_order: product.DisplayOrder,
        preparation_time_minutes: product.PreparationTimeMinutes,
        featured: false, // Cópia não destacada
        is_new: false, // Cópia não é nova
        promotion_price: product.PromotionPrice,
        promotion_start: product.PromotionStart,
        promotion_end: product.PromotionEnd,
        available_from: product.AvailableFrom,
        available_until: product.AvailableUntil,
        sku: product.SKU,
        internal_notes: product.InternalNotes,
      };
      
      const created = await createProduct(payload);
      products = [...products, created];
    } catch (e: any) {
      error = e?.message ?? 'Erro ao duplicar produto.';
    }
  }

  async function archiveProduct(id: number) {
    if (!confirm('Tem certeza que deseja arquivar este produto?')) return;
    try {
      const product = products.find(p => p.ID === id);
      if (!product) return;
      
      const payload = {
        name: product.Name,
        description: product.Description,
        price: product.Price,
        is_composto: product.IsComposto,
        active: false, // Arquivar = desativar
        photo_url: product.PhotoURL,
        category_id: product.CategoryID,
        display_order: product.DisplayOrder,
        preparation_time_minutes: product.PreparationTimeMinutes,
        featured: product.Featured,
        is_new: product.IsNew,
        promotion_price: product.PromotionPrice,
        promotion_start: product.PromotionStart,
        promotion_end: product.PromotionEnd,
        available_from: product.AvailableFrom,
        available_until: product.AvailableUntil,
        sku: product.SKU,
        internal_notes: product.InternalNotes,
      };
      
      const updated = await updateProduct(id, payload);
      products = products.map(p => p.ID === id ? updated : p);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao arquivar produto.';
    }
  }

  async function toggleProductActive(id: number) {
    try {
      const product = products.find(p => p.ID === id);
      if (!product) return;
      
      const payload = {
        name: product.Name,
        description: product.Description,
        price: product.Price,
        is_composto: product.IsComposto,
        active: !product.Active,
        photo_url: product.PhotoURL,
        category_id: product.CategoryID,
        display_order: product.DisplayOrder,
        preparation_time_minutes: product.PreparationTimeMinutes,
        featured: product.Featured,
        is_new: product.IsNew,
        promotion_price: product.PromotionPrice,
        promotion_start: product.PromotionStart,
        promotion_end: product.PromotionEnd,
        available_from: product.AvailableFrom,
        available_until: product.AvailableUntil,
        sku: product.SKU,
        internal_notes: product.InternalNotes,
      };
      
      const updated = await updateProduct(id, payload);
      products = products.map(p => p.ID === id ? updated : p);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao alterar status do produto.';
    }
  }

  async function toggleProductFeatured(id: number) {
    try {
      const product = products.find(p => p.ID === id);
      if (!product) return;
      
      const payload = {
        name: product.Name,
        description: product.Description,
        price: product.Price,
        is_composto: product.IsComposto,
        active: product.Active,
        photo_url: product.PhotoURL,
        category_id: product.CategoryID,
        display_order: product.DisplayOrder,
        preparation_time_minutes: product.PreparationTimeMinutes,
        featured: !product.Featured,
        is_new: product.IsNew,
        promotion_price: product.PromotionPrice,
        promotion_start: product.PromotionStart,
        promotion_end: product.PromotionEnd,
        available_from: product.AvailableFrom,
        available_until: product.AvailableUntil,
        sku: product.SKU,
        internal_notes: product.InternalNotes,
      };
      
      const updated = await updateProduct(id, payload);
      products = products.map(p => p.ID === id ? updated : p);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao alterar destaque do produto.';
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
    products.filter(p => {
      // Filtro de busca
      if (productSearch && !p.Name.toLowerCase().includes(productSearch.toLowerCase())) {
        return false;
      }
      
      // Filtro de status
      switch (productFilter) {
        case 'active':
          return p.Active;
        case 'archived':
          return !p.Active;
        case 'promotion':
          return p.PromotionPrice !== undefined && p.PromotionPrice > 0;
        case 'new':
          return p.IsNew;
        case 'featured':
          return p.Featured;
        case 'composto':
          return p.IsComposto;
        default:
          return true;
      }
    })
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
          <label class="filter-label">Filtrar por</label>
          <Select bind:value={productFilter} class="filter-select">
            <option value="all">Todos</option>
            <option value="active">Ativos</option>
            <option value="archived">Arquivados</option>
            <option value="promotion">Em promoção</option>
            <option value="new">Novidades</option>
            <option value="featured">Destaques</option>
            <option value="composto">Compostos</option>
          </Select>
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
          <ProductCard
            product={product}
            onClick={() => window.location.href = `/products/${product.ID}/edit`}
            onEdit={() => window.location.href = `/products/${product.ID}/edit`}
            onDuplicate={() => duplicateProduct(product.ID)}
            onArchive={() => archiveProduct(product.ID)}
            onToggleActive={() => toggleProductActive(product.ID)}
            onToggleFeatured={() => toggleProductFeatured(product.ID)}
          />
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
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }

  @media (max-width: 768px) {
    .filters-grid {
      grid-template-columns: 1fr;
    }
  }

  .filter-select {
    width: 100%;
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
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }

  @media (max-width: 768px) {
    .products-grid {
      grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
      gap: 0.75rem;
    }
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
