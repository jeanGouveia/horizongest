<script lang="ts">
  import { onMount } from 'svelte';
  import { getCategories, createCategory, updateCategory, deleteCategory } from '$lib/api/category';
  import { api } from '$lib/api/client';
  import type { Category } from '$lib/types/category';
  import type { DependencyCheck } from '$lib/types/dependency';
  import { Card, Button, Input, Textarea, Checkbox, Badge, Alert, Loading, EmptyState, Modal } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Plus, Folder, AlertTriangle, MoreHorizontal } from '@lucide/svelte';

  let categories: Category[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let categorySearch = $state('');

  // Modal nova categoria
  let showCategoryModal = $state(false);
  let categoryEditMode = $state(false);
  let categoryEditId = $state<number | null>(null);
  let categoryForm = $state({ Name: '', Description: '', DisplayOrder: 0, Active: true });
  let categorySaving = $state(false);
  let categoryError = $state('');

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
      categories = await getCategories();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar categorias.';
    } finally {
      loading = false;
    }
  }

  async function saveCategory() {
    categorySaving = true;
    categoryError = '';
    try {
      const payload = {
        name: categoryForm.Name,
        description: categoryForm.Description,
        display_order: Number(categoryForm.DisplayOrder),
        active: categoryForm.Active,
      };

      if (categoryEditMode && categoryEditId) {
        const updated = await updateCategory(categoryEditId, payload);
        categories = categories.map(c => c.ID === categoryEditId ? updated : c);
      } else {
        const created = await createCategory(payload);
        categories = [...categories, created];
      }

      showCategoryModal = false;
      categoryForm = { Name: '', Description: '', DisplayOrder: 0, Active: true };
      categoryEditMode = false;
      categoryEditId = null;
    } catch (e: any) {
      categoryError = e?.message ?? 'Erro ao salvar categoria.';
    } finally {
      categorySaving = false;
    }
  }

  function openCategoryEdit(category: Category) {
    categoryEditMode = true;
    categoryEditId = category.ID;
    categoryForm = {
      Name: category.Name,
      Description: category.Description ?? '',
      DisplayOrder: category.DisplayOrder,
      Active: category.Active,
    };
    showCategoryModal = true;
  }

  function openCategoryCreate() {
    categoryEditMode = false;
    categoryEditId = null;
    categoryForm = { Name: '', Description: '', DisplayOrder: 0, Active: true };
    showCategoryModal = true;
  }

  async function deleteCategoryById(id: number) {
    try {
      const res = await api.canDeleteCategory(id);
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

      if (!confirm('Tem certeza que deseja excluir esta categoria?')) return;
      await deleteCategory(id);
      categories = categories.filter(c => c.ID !== id);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao excluir categoria.';
    }
  }

  async function confirmDeleteCategory() {
    if (!deleteTargetId) return;
    try {
      await deleteCategory(deleteTargetId);
      categories = categories.filter(c => c.ID !== deleteTargetId);
      showDependencyModal = false;
      dependencyCheck = null;
      deleteTargetId = null;
    } catch (e: any) {
      error = e?.message ?? 'Erro ao excluir categoria.';
    }
  }

  const filteredCategories = $derived(
    categories.filter(c => !categorySearch || c.Name.toLowerCase().includes(categorySearch.toLowerCase()))
    .sort((a, b) => a.DisplayOrder - b.DisplayOrder || a.Name.localeCompare(b.Name))
  );
</script>

<Workspace
  breadcrumb={[{ label: 'Categorias' }]}
  title="Gerenciamento de Categorias"
  description="Organização do cardápio"
>
  <svelte:fragment slot="actions">
    <Button onclick={openCategoryCreate} variant="primary">
      <Plus size={16} />
      Nova Categoria
    </Button>
  </svelte:fragment>

  {#if loading}
    <div class="loading-state">
      <Loading />
      <span>Carregando categorias…</span>
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
          <span>Filtros</span>
        </div>
        <div class="filters-stats">
          <span class="stat-item">{categories.length} categorias</span>
        </div>
      </div>
      <div class="filters-grid">
        <div class="filter-group">
          <label class="filter-label">Buscar categorias</label>
          <Input
            placeholder="Buscar por nome..."
            bind:value={categorySearch}
          />
        </div>
      </div>
    </Card>

    <!-- Categorias -->
    <div class="section-header">
      <h2 class="section-title">Categorias</h2>
    </div>

    {#if filteredCategories.length === 0}
      <div class="empty-state">
        <Folder size={48} class="empty-icon" />
        <span class="empty-title">{categorySearch ? 'Nenhuma categoria encontrada' : 'Nenhuma categoria cadastrada'}</span>
        <span class="empty-subtitle">{categorySearch ? 'Tente outro termo de busca' : 'Comece adicionando sua primeira categoria'}</span>
        {#if !categorySearch}
          <Button onclick={openCategoryCreate} variant="primary">Criar Categoria</Button>
        {/if}
      </div>
    {:else}
      <div class="categories-grid">
        {#each filteredCategories as category}
          <Card class="category-card">
            <div class="category-card-top">
              <div class="category-info">
                <span class="category-name">{category.Name}</span>
                {#if category.Description}
                  <p class="category-desc">{category.Description}</p>
                {/if}
              </div>
              <div class="category-order">#{category.DisplayOrder}</div>
            </div>
            
            <div class="category-meta">
              <Badge variant={category.Active ? 'success' : 'default'} size="sm">
                {category.Active ? 'Ativo' : 'Inativo'}
              </Badge>
            </div>

            <div class="category-actions">
              <Button variant="ghost" size="sm" onclick={() => openCategoryEdit(category)}>
                Editar
              </Button>
              <Button variant="ghost" size="sm" onclick={() => deleteCategoryById(category.ID)} class="danger">
                Excluir
              </Button>
            </div>
          </Card>
        {/each}
      </div>
    {/if}

  {/if}
</Workspace>

<!-- Modal: Nova Categoria -->
<Modal 
  open={showCategoryModal} 
  title={categoryEditMode ? 'Editar Categoria' : 'Nova Categoria'}
  onClose={() => (showCategoryModal = false)}
>
  {#if categoryError}
    <Alert variant="error">{categoryError}</Alert>
  {/if}

  <Input
    id="c-name"
    label="Nome *"
    bind:value={categoryForm.Name}
    placeholder="Ex: Lanches"
  />
  <Textarea
    id="c-desc"
    label="Descrição"
    bind:value={categoryForm.Description}
    placeholder="Descrição opcional"
    rows={2}
  />
  <Input
    id="c-order"
    label="Ordem de Exibição *"
    type="number"
    min="0"
    bind:value={categoryForm.DisplayOrder}
  />
  <Checkbox
    label="Categoria ativa (visível no cardápio)"
    bind:checked={categoryForm.Active}
  />

  <div class="modal-actions">
    <Button variant="ghost" onclick={() => (showCategoryModal = false)}>Cancelar</Button>
    <Button 
      variant="primary" 
      onclick={saveCategory} 
      disabled={categorySaving || !categoryForm.Name}
      loading={categorySaving}
    >
      {categoryEditMode ? 'Atualizar Categoria' : 'Criar Categoria'}
    </Button>
  </div>
</Modal>

<style>
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 4rem;
    color: #64748b;
    font-size: 0.875rem;
  }

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

  .categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  .category-card {
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .category-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px 0 rgb(0 0 0 / 0.08);
  }

  .category-card-top {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .category-info {
    flex: 1;
  }

  .category-name {
    font-weight: 600;
    font-size: 1rem;
    color: #0f172a;
    display: block;
    margin-bottom: 0.5rem;
  }

  .category-desc {
    font-size: 0.875rem;
    color: #64748b;
    margin: 0;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .category-order {
    font-weight: 700;
    color: #6366f1;
    font-size: 1.25rem;
    letter-spacing: -0.025em;
  }

  .category-meta {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .category-actions {
    display: flex;
    gap: 0.5rem;
    padding-top: 1rem;
    border-top: 1px solid #f1f5f9;
  }

  .category-actions .danger {
    color: #ef4444;
  }

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

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1.75rem;
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

  /* Responsive */
  @media (max-width: 768px) {
    .categories-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }

    .filters-grid {
      grid-template-columns: 1fr;
    }

    .category-actions {
      flex-direction: column;
    }

    .category-actions .danger {
      color: #ef4444;
    }
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
