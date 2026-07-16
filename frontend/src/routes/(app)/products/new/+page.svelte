<script lang="ts">
  import { onMount } from 'svelte';
  import { createProduct } from '$lib/api/product';
  import { getCategories } from '$lib/api/category';
  import type { Product } from '$lib/types/product';
  import type { Category as CategoryType } from '$lib/types/category';
  import { Button, Input, Textarea, Select, Checkbox, Card, Alert, TabNavigation, PhotoUpload, PageContainer, Loading } from '$lib/components/ui';
  import { ArrowLeft, Save } from '@lucide/svelte';

  let loading = $state(false);
  let saving = $state(false);
  let error = $state('');
  let categories: CategoryType[] = $state([]);
  let activeTab = $state('info');

  let form = $state({
    name: '',
    description: '',
    price: 0,
    is_composto: false,
    active: true,
    photo_url: '',
    category_id: undefined as string | undefined,
    display_order: 0,
    preparation_time_minutes: 0,
    featured: false,
    is_new: false,
    promotion_price: undefined as number | undefined,
    promotion_start: '',
    promotion_end: '',
    available_from: '',
    available_until: '',
    sku: '',
    internal_notes: ''
  });

  let errors = $state<Record<string, string>>({});

  const tabs = [
    { id: 'info', label: 'Informações' },
    { id: 'sales', label: 'Venda' },
    { id: 'production', label: 'Produção' }
  ];

  onMount(async () => {
    await loadCategories();
  });

  async function loadCategories() {
    loading = true;
    try {
      categories = await getCategories();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar categorias.';
    } finally {
      loading = false;
    }
  }

  function validate() {
    errors = {};
    
    if (!form.name.trim()) {
      errors.name = 'Nome é obrigatório';
    }
    
    if (form.price <= 0) {
      errors.price = 'Preço deve ser maior que zero';
    }
    
    if (form.promotion_price && form.promotion_price >= form.price) {
      errors.promotion_price = 'Preço promocional deve ser menor que o preço normal';
    }
    
    return Object.keys(errors).length === 0;
  }

  async function handleSave() {
    if (!validate()) {
      return;
    }

    saving = true;
    error = '';
    try {
      const payload = {
        name: form.name,
        description: form.description,
        price: Number(form.price),
        is_composto: form.is_composto,
        active: form.active,
        photo_url: form.photo_url,
        category_id: form.category_id ? Number(form.category_id) : undefined,
        display_order: Number(form.display_order),
        preparation_time_minutes: Number(form.preparation_time_minutes),
        featured: form.featured,
        is_new: form.is_new,
        promotion_price: form.promotion_price ? Number(form.promotion_price) : undefined,
        promotion_start: form.promotion_start || undefined,
        promotion_end: form.promotion_end || undefined,
        available_from: form.available_from,
        available_until: form.available_until,
        sku: form.sku,
        internal_notes: form.internal_notes,
      };

      await createProduct(payload);
      window.location.href = '/products';
    } catch (e: any) {
      error = e?.message ?? 'Erro ao salvar produto.';
    } finally {
      saving = false;
    }
  }

  function handleCancel() {
    window.location.href = '/products';
  }

  function handlePhotoChange(detail: { file: File | null; previewUrl: string }) {
    form.photo_url = detail.previewUrl;
  }
</script>

<PageContainer>
  <div class="page-header">
    <div class="header-left">
      <Button variant="ghost" size="sm" onclick={handleCancel} class="back-button">
        <ArrowLeft size={16} />
        Produtos
      </Button>
      <div class="header-title">
        <h1>Cadastro de Produto</h1>
        <p class="header-subtitle">Crie um produto que poderá ser vendido no sistema, cardápio digital e marketplaces.</p>
      </div>
    </div>
    <div class="header-actions">
      <Button variant="ghost" onclick={handleCancel} disabled={saving}>
        Cancelar
      </Button>
      <Button variant="primary" onclick={handleSave} disabled={saving} loading={saving}>
        <Save size={16} />
        Salvar Produto
      </Button>
    </div>
  </div>

  {#if loading}
    <div class="loading-state">
      <Loading type="skeleton" />
    </div>
  {:else}
    {#if error}
      <Alert variant="error" dismissible onDismiss={() => error = ''}>
        {error}
      </Alert>
    {/if}

    <TabNavigation tabs={tabs} bind:activeTab={activeTab} />

    <div class="form-container">
      {#if activeTab === 'info'}
        <div class="info-grid">
          <div class="left-column">
            <Card class="section-card">
              <h3 class="section-title">Foto do Produto</h3>
              <PhotoUpload 
                photoUrl={form.photo_url} 
                onPhotoChange={handlePhotoChange}
                size="lg"
              />
            </Card>

            <Card class="section-card">
              <h3 class="section-title">Informações Básicas</h3>
              <Input
                label="Nome *"
                placeholder="Ex: Feijoada Completa"
                bind:value={form.name}
                error={errors.name}
              />
              <p class="helper-text">Nome do produto como aparecerá no cardápio</p>
              <Textarea
                label="Descrição"
                placeholder="Descrição detalhada do produto"
                bind:value={form.description}
                rows={3}
              />
              <p class="helper-text">Descreva os ingredientes e características do produto</p>
              <Select
                label="Categoria"
                placeholder="Selecione uma categoria"
                bind:value={form.category_id}
              >
                <option value="">Selecione uma categoria</option>
                {#each categories as cat}
                  <option value={cat.ID}>{cat.Name}</option>
                {/each}
              </Select>
              <p class="helper-text">Categoria para organização do cardápio</p>
              <Input
                label="SKU"
                placeholder="Ex: PROD-001"
                bind:value={form.sku}
              />
              <p class="helper-text">Código para integração com sistemas externos</p>
            </Card>
          </div>

          <div class="right-column">
            <Card class="section-card">
              <h3 class="section-title">Preço e Configuração</h3>
              <Input
                label="Preço (R$) *"
                type="number"
                min="0"
                step="0.01"
                bind:value={form.price}
                error={errors.price}
              />
              <p class="helper-text">Preço de venda do produto</p>
              <Input
                label="Tempo de Preparo (minutos)"
                type="number"
                min="0"
                bind:value={form.preparation_time_minutes}
              />
              <p class="helper-text">Tempo médio para preparar o produto</p>
              <Input
                label="Ordem de Exibição"
                type="number"
                min="0"
                bind:value={form.display_order}
              />
              <p class="helper-text">Ordem em que o produto aparecerá no cardápio</p>
              <div class="checkbox-group">
                <Checkbox
                  label="Produto Ativo"
                  bind:checked={form.active}
                />
                <p class="helper-text">Produto visível no cardápio e disponível para venda</p>
                <Checkbox
                  label="Produto Composto"
                  bind:checked={form.is_composto}
                />
                <p class="helper-text">Produto que requer ficha técnica com ingredientes</p>
              </div>
            </Card>
          </div>
        </div>
      {/if}

      {#if activeTab === 'sales'}
        <div class="sales-grid">
          <Card class="section-card">
            <h3 class="section-title">Destaque e Novidade</h3>
            <div class="checkbox-group">
              <Checkbox
                label="Produto em Destaque"
                bind:checked={form.featured}
              />
              <p class="helper-text">Destacar o produto na página inicial e seções especiais</p>
              <Checkbox
                label="Mostrar Selo NOVO"
                bind:checked={form.is_new}
              />
              <p class="helper-text">Exibir selo de produto novo por tempo limitado</p>
            </div>
          </Card>

          <Card class="section-card">
            <h3 class="section-title">Promoção</h3>
            <Input
              label="Preço Promocional (R$)"
              type="number"
              min="0"
              step="0.01"
              bind:value={form.promotion_price}
              error={errors.promotion_price}
            />
            <p class="helper-text">Preço especial válido apenas durante o período informado</p>
            <div class="date-row">
              <div class="date-field">
                <Input
                  label="Início da Promoção"
                  type="datetime-local"
                  bind:value={form.promotion_start}
                />
                <p class="helper-text">Data e hora de início</p>
              </div>
              <div class="date-field">
                <Input
                  label="Fim da Promoção"
                  type="datetime-local"
                  bind:value={form.promotion_end}
                />
                <p class="helper-text">Data e hora de término</p>
              </div>
            </div>
          </Card>

          <Card class="section-card">
            <h3 class="section-title">Disponibilidade</h3>
            <div class="date-row">
              <div class="date-field">
                <Input
                  label="Disponível Das"
                  type="time"
                  bind:value={form.available_from}
                />
                <p class="helper-text">Horário de início de disponibilidade</p>
              </div>
              <div class="date-field">
                <Input
                  label="Disponível Até"
                  type="time"
                  bind:value={form.available_until}
                />
                <p class="helper-text">Horário de fim de disponibilidade</p>
              </div>
            </div>
          </Card>
        </div>
      {/if}

      {#if activeTab === 'production'}
        <Card class="section-card">
          <h3 class="section-title">Produção</h3>
          <Textarea
            label="Observações Internas"
            placeholder="Notas para equipe de produção, instruções especiais, etc."
            bind:value={form.internal_notes}
            rows={6}
          />
          <p class="helper-text">Informações visíveis apenas para administradores e equipe de produção</p>
          {#if form.is_composto}
            <div class="production-note">
              <p>⚠️ Este produto é composto. Configure a ficha técnica após salvar o produto.</p>
            </div>
          {:else}
            <div class="production-note">
              <p>ℹ️ Este é um produto simples. Não requer ficha técnica.</p>
            </div>
          {/if}
        </Card>
      {/if}
    </div>
  {/if}
</PageContainer>

<style>
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
    gap: 2rem;
  }

  .header-left {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .back-button {
    width: fit-content;
  }

  .header-title h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #0f172a;
    margin: 0;
    letter-spacing: -0.025em;
  }

  .header-subtitle {
    font-size: 0.875rem;
    color: #64748b;
    margin: 0.5rem 0 0 0;
    line-height: 1.5;
  }

  .header-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .loading-state {
    padding: 4rem;
  }

  .form-container {
    margin-top: 2rem;
  }

  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }

  .left-column,
  .right-column {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .section-card {
    padding: 1.5rem;
  }

  .section-title {
    font-size: 1rem;
    font-weight: 600;
    color: #0f172a;
    margin: 0 0 1.5rem 0;
  }

  .checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .sales-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 1.5rem;
  }

  .date-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .date-field {
    display: flex;
    flex-direction: column;
  }

  .production-note {
    margin-top: 1.5rem;
    padding: 1rem;
    background: #f8fafc;
    border-left: 3px solid #6366f1;
    border-radius: 4px;
  }

  .production-note p {
    margin: 0;
    font-size: 0.875rem;
    color: #64748b;
    line-height: 1.5;
  }

  .helper-text {
    font-size: 0.75rem;
    color: #64748b;
    margin-top: 0.25rem;
    line-height: 1.4;
  }

  @media (max-width: 1024px) {
    .page-header {
      flex-direction: column;
      align-items: stretch;
    }

    .header-actions {
      justify-content: flex-end;
    }

    .info-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 768px) {
    .header-title h1 {
      font-size: 1.5rem;
    }

    .sales-grid {
      grid-template-columns: 1fr;
    }

    .date-row {
      grid-template-columns: 1fr;
    }
  }
</style>
