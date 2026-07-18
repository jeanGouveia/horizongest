<script lang="ts">
  import { onMount } from 'svelte';
  import { getIngredients, createIngredient, updateIngredient, deleteIngredient, updateIngredientStock } from '$lib/api/product';
  import { api } from '$lib/api/client';
  import type { Ingredient } from '$lib/types/ingredient';
  import type { DependencyCheck } from '$lib/types/dependency';
  import { Card, Button, Input, Select, Badge, Alert, Modal } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Search, Plus, AlertTriangle, TrendingUp, TrendingDown, ArrowUpDown, Package } from '@lucide/svelte';

  let ingredients: Ingredient[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let showLowStockOnly = $state(false);
  let ingredientSearch = $state('');
  let ingredientSortBy = $state<string>('name');
  let ingredientSortOrder = $state<'asc' | 'desc'>('asc');
  let ingredientCurrentPage = $state(1);
  const itemsPerPage = 12;

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

  // Modal de dependências
  let showDependencyModal = $state(false);
  let dependencyCheck = $state<DependencyCheck | null>(null);
  let deleteTargetId = $state<number | null>(null);

  onMount(async () => {
    await loadAll();
  });

  async function loadAll() {
    loading = true;
    error = '';
    try {
      ingredients = await getIngredients();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar ingredientes.';
    } finally {
      loading = false;
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
      // Melhorar tratamento de erro para mostrar mensagens específicas
      if (e?.message) {
        try {
          const errorData = JSON.parse(e.message);
          if (errorData.fields) {
            const fieldMessages = Object.entries(errorData.fields).map(([field, msg]) => {
              const fieldMap: Record<string, string> = {
                name: 'Nome',
                unit: 'Unidade',
                stock_quantity: 'Estoque inicial',
                min_stock: 'Estoque mínimo'
              };
              return `${fieldMap[field] || field}: ${msg}`;
            });
            ingError = fieldMessages.join('. ');
          } else {
            ingError = e.message;
          }
        } catch {
          ingError = e.message;
        }
      } else {
        ingError = 'Erro ao salvar ingrediente.';
      }
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
    try {
      const res = await api.canDeleteIngredient(id);
      if (res.error) {
        error = res.error;
        return;
      }

      const check = res.data as DependencyCheck;
      if (!check.canDelete) {
        dependencyCheck = check;
        deleteTargetId = id;
        showDependencyModal = true;
        return;
      }

      if (!confirm('Tem certeza que deseja excluir este ingrediente?')) return;
      await deleteIngredient(id);
      ingredients = ingredients.filter(i => i.ID !== id);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao excluir ingrediente.';
    }
  }

  async function confirmDeleteIngredient() {
    if (!deleteTargetId) return;
    try {
      await deleteIngredient(deleteTargetId);
      ingredients = ingredients.filter(i => i.ID !== deleteTargetId);
      showDependencyModal = false;
      dependencyCheck = null;
      deleteTargetId = null;
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

  function getStockStatus(ing: Ingredient) {
    if (ing.StockQuantity === 0) {
      return { label: 'Estoque Zerado', variant: 'danger' as const, icon: AlertTriangle };
    } else if (ing.StockQuantity <= ing.MinStock) {
      return { label: 'Estoque Baixo', variant: 'warning' as const, icon: TrendingDown };
    }
    return { label: 'OK', variant: 'success' as const, icon: TrendingUp };
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

  const ingredientTotalPages = $derived(Math.ceil(filteredIngredients.length / itemsPerPage));
  const paginatedIngredients = $derived(
    filteredIngredients.slice((ingredientCurrentPage - 1) * itemsPerPage, ingredientCurrentPage * itemsPerPage)
  );

  function goToIngredientPage(page: number) {
    if (page >= 1 && page <= ingredientTotalPages) {
      ingredientCurrentPage = page;
    }
  }
</script>

<Workspace
  breadcrumb={[{ label: 'Ingredientes' }]}
  title="Gerenciamento de Ingredientes"
  description="Controle de estoque e cadastro de ingredientes"
>
  <svelte:fragment slot="actions">
    <Button onclick={openIngredientCreate} variant="primary">
      <Plus size={16} />
      Novo Ingrediente
    </Button>
  </svelte:fragment>

  {#if loading}
    <div class="skeleton-grid">
      {#each Array(6) as _}
        <Card class="skeleton-card">
          <div class="skeleton-card-top">
            <div class="skeleton-content">
              <div class="skeleton-line"></div>
              <div class="skeleton-line short"></div>
            </div>
          </div>
          <div class="skeleton-card-bottom">
            <div class="skeleton-line"></div>
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
    <!-- Filtros -->
    <Card class="filters-card">
      <div class="filters-header">
        <div class="filters-title">
          <Search size={18} />
          <span>Filtros</span>
        </div>
        <div class="filters-stats">
          <span class="stat-item">{ingredients.length} ingredientes</span>
        </div>
      </div>
      <div class="filters-grid">
        <div class="filter-group">
          <label class="filter-label">Buscar ingredientes</label>
          <div class="search-wrapper">
            <Search size={16} class="search-icon" />
            <Input
              placeholder="Buscar por nome..."
              bind:value={ingredientSearch}
              class="search-input"
            />
          </div>
        </div>
        <div class="filter-group">
          <label class="filter-label">Ordenar por</label>
          <div class="sort-wrapper">
            <Select bind:value={ingredientSortBy} class="sort-select">
              <option value="name">Nome</option>
              <option value="stock">Estoque</option>
            </Select>
            <Button
              variant="ghost"
              size="sm"
              onclick={() => (ingredientSortOrder = ingredientSortOrder === 'asc' ? 'desc' : 'asc')}
              class="sort-toggle"
            >
              <ArrowUpDown size={16} />
            </Button>
          </div>
        </div>
        <div class="filter-group">
          <label class="filter-label">&nbsp;</label>
          <label class="filter-toggle">
            <input type="checkbox" bind:checked={showLowStockOnly} />
            <span>Apenas estoque baixo</span>
          </label>
        </div>
      </div>
    </Card>

    {#if filteredIngredients.length === 0}
      <div class="empty-state">
        <Package size={48} class="empty-icon" />
        <span class="empty-title">{showLowStockOnly ? 'Nenhum ingrediente com estoque baixo' : (ingredientSearch ? 'Nenhum ingrediente encontrado' : 'Nenhum ingrediente cadastrado')}</span>
        <span class="empty-subtitle">{showLowStockOnly ? 'Ótimo! Todos os ingredientes estão com estoque adequado' : (ingredientSearch ? 'Tente outro termo de busca' : 'Adicione ingredientes para começar')}</span>
        {#if !showLowStockOnly && !ingredientSearch}
          <Button onclick={openIngredientCreate} variant="primary">Adicionar Ingrediente</Button>
        {/if}
      </div>
    {:else}
      <div class="ingredients-grid">
        {#each paginatedIngredients as ing}
          {@const stockStatus = getStockStatus(ing)}
          <Card class={`ingredient-card ${stockStatus.variant === 'danger' ? 'ingredient-critical' : ''} ${stockStatus.variant === 'warning' ? 'ingredient-low' : ''}`}>
            <div class="ingredient-header">
              <div class="ingredient-name">{ing.Name}</div>
              <Badge variant={stockStatus.variant} size="sm" icon>
                <svelte:component this={stockStatus.icon} size={12} />
                {stockStatus.label}
              </Badge>
            </div>
            
            <div class="ingredient-stock">
              <div class="stock-value">
                <span class="stock-number">{ing.StockQuantity}</span>
                <span class="stock-unit">{ing.Unit}</span>
              </div>
              <div class="stock-min">Mín: {ing.MinStock}</div>
            </div>

            <div class="ingredient-actions">
              <Button variant="ghost" size="sm" onclick={() => openIngredientEdit(ing)}>
                Editar
              </Button>
              <Button variant="ghost" size="sm" onclick={() => openStockModal(ing)}>
                Ajustar
              </Button>
              <Button variant="ghost" size="sm" onclick={() => deleteIngredientById(ing.ID)} class="danger">
                Excluir
              </Button>
            </div>
          </Card>
        {/each}
      </div>

      {#if ingredientTotalPages > 1}
        <div class="pagination">
          <Button
            variant="ghost"
            size="sm"
            disabled={ingredientCurrentPage === 1}
            onclick={() => goToIngredientPage(ingredientCurrentPage - 1)}
          >
            Anterior
          </Button>
          <span class="pagination-info">
            Página {ingredientCurrentPage} de {ingredientTotalPages}
          </span>
          <Button
            variant="ghost"
            size="sm"
            disabled={ingredientCurrentPage === ingredientTotalPages}
            onclick={() => goToIngredientPage(ingredientCurrentPage + 1)}
          >
            Próxima
          </Button>
        </div>
      {/if}
    {/if}
  {/if}
</Workspace>

<!-- Modal: Novo Ingrediente -->
<Modal 
  open={showIngModal} 
  title={ingEditMode ? 'Editar Ingrediente' : 'Novo Ingrediente'}
  onClose={() => (showIngModal = false)}
>
  {#if ingError}
    <Alert variant="error">{ingError}</Alert>
  {/if}

  <Input
    id="i-name"
    label="Nome *"
    bind:value={ingForm.Name}
    placeholder="Ex: Feijão Preto"
  />
  <Select
    id="i-unit"
    label="Unidade *"
    bind:value={ingForm.Unit}
  >
    <option value="">Selecione...</option>
    <option value="kg">kg (quilograma)</option>
    <option value="g">g (grama)</option>
    <option value="L">L (litro)</option>
    <option value="ml">ml (mililitro)</option>
    <option value="un">un (unidade)</option>
  </Select>
  <Input
    id="i-stock"
    label="Estoque inicial"
    type="number"
    min="0"
    step="0.01"
    bind:value={ingForm.StockQuantity}
  />
  <Input
    id="i-minstock"
    label="Estoque mínimo"
    type="number"
    min="0"
    step="0.01"
    bind:value={ingForm.MinStock}
  />

  <div class="modal-actions">
    <Button variant="ghost" onclick={() => (showIngModal = false)}>Cancelar</Button>
    <Button 
      variant="primary" 
      onclick={saveIngredient} 
      disabled={ingSaving || !ingForm.Name || !ingForm.Unit}
      loading={ingSaving}
    >
      {ingEditMode ? 'Atualizar Ingrediente' : 'Criar Ingrediente'}
    </Button>
  </div>
</Modal>

<!-- Modal: Ajustar Estoque -->
<Modal 
  open={showStockModal} 
  title="Ajustar Estoque"
  onClose={() => (showStockModal = false)}
>
  {#if stockError}
    <Alert variant="error">{stockError}</Alert>
  {/if}

  <Input
    id="s-quantity"
    label="Quantidade *"
    type="number"
    min="0"
    step="0.01"
    bind:value={stockForm.Quantity}
  />

  <div class="modal-actions">
    <Button variant="ghost" onclick={() => (showStockModal = false)}>Cancelar</Button>
    <Button 
      variant="primary" 
      onclick={saveStock} 
      disabled={stockSaving || stockForm.Quantity < 0}
      loading={stockSaving}
    >
      Ajustar Estoque
    </Button>
  </div>
</Modal>

<style>
  /* Filters Card */
  .filters-card {
    margin-bottom: 1rem;
  }

  .filters-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
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
    gap: 0.75rem;
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

  .filter-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: #0f172a;
    cursor: pointer;
    user-select: none;
    padding: 0.5rem 0.75rem;
    border-radius: 8px;
    transition: background 0.15s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid #f1f5f9;
  }

  .filter-toggle:hover {
    background: #f8fafc;
    border-color: #e2e8f0;
  }

  .filter-toggle input {
    cursor: pointer;
  }

  /* Ingredients Grid */
  .ingredients-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .ingredient-card {
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .ingredient-card:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px 0 rgb(0 0 0 / 0.08);
  }

  .ingredient-critical {
    border-color: #fee2e2;
    background: #fef2f2;
  }

  .ingredient-low {
    border-color: #fef3c7;
    background: #fffbeb;
  }

  .ingredient-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
  }

  .ingredient-name {
    font-weight: 600;
    font-size: 0.875rem;
    color: #0f172a;
  }

  .ingredient-stock {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: #f8fafc;
    border-radius: 8px;
    margin-bottom: 0.75rem;
  }

  .stock-value {
    display: flex;
    align-items: baseline;
    gap: 0.25rem;
  }

  .stock-number {
    font-size: 1.5rem;
    font-weight: 700;
    color: #0f172a;
  }

  .stock-unit {
    font-size: 0.875rem;
    color: #64748b;
  }

  .stock-min {
    font-size: 0.75rem;
    color: #64748b;
  }

  .ingredient-actions {
    display: flex;
    gap: 0.5rem;
  }

  .ingredient-actions .danger {
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
    margin-top: 1rem;
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
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 0.75rem;
  }

  .skeleton-card {
    padding: 1rem;
  }

  .skeleton-card-top {
    margin-bottom: 0.75rem;
  }

  .skeleton-content {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .skeleton-line {
    height: 16px;
    background: #e2e8f0;
    border-radius: 4px;
    animation: pulse 1.5s ease-in-out infinite;
  }

  .skeleton-line.short {
    width: 60%;
  }

  .skeleton-card-bottom {
    margin-top: 0.75rem;
  }

  @keyframes pulse {
    0%, 100% { opacity: 0.4; }
    50% { opacity: 1; }
  }

  /* Responsive */
  @media (max-width: 768px) {
    .filters-grid {
      grid-template-columns: 1fr;
    }

    .ingredients-grid {
      grid-template-columns: 1fr;
    }
  }

  /* Dependency Modal */
  .dependency-modal {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .dependency-list {
    margin-top: 1rem;
  }

  .dependency-list h4 {
    font-size: 0.875rem;
    font-weight: 600;
    color: #0f172a;
    margin-bottom: 0.75rem;
  }

  .dependency-list ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .dependency-item {
    padding: 0.75rem;
    background: #fef3c7;
    border: 1px solid #fde68a;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .dependency-type {
    font-size: 0.6875rem;
    font-weight: 600;
    color: #92400e;
    text-transform: uppercase;
  }

  .dependency-name {
    font-size: 0.8125rem;
    font-weight: 500;
    color: #0f172a;
  }

  .dependency-desc {
    font-size: 0.6875rem;
    color: #64748b;
  }

  .dependency-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1rem;
  }
</style>

<!-- Modal de Dependências -->
<Modal
  open={showDependencyModal}
  onClose={() => {
    showDependencyModal = false;
    dependencyCheck = null;
    deleteTargetId = null;
  }}
  title="Não é possível excluir"
>
  <div class="dependency-modal">
    <Alert variant="warning" dismissible={false}>
      Este item possui dependências que impedem sua exclusão.
    </Alert>

    {#if dependencyCheck && dependencyCheck.reasons.length > 0}
      <div class="dependency-list">
        <h4>Dependências encontradas:</h4>
        <ul>
          {#each dependencyCheck.reasons as reason}
            <li class="dependency-item">
              <div class="dependency-type">{reason.type}</div>
              <div class="dependency-name">{reason.name}</div>
              <div class="dependency-desc">{reason.description}</div>
            </li>
          {/each}
        </ul>
      </div>
    {/if}

    <div class="dependency-actions">
      <Button variant="ghost" onclick={() => {
        showDependencyModal = false;
        dependencyCheck = null;
        deleteTargetId = null;
      }}>
        Fechar
      </Button>
    </div>
  </div>
</Modal>
