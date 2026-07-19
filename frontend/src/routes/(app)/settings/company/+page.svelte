<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { Button, Input, Alert, Card, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Building2, Palette, Briefcase, Check, AlertTriangle } from '@lucide/svelte';
  import { themeStore } from '$lib/stores/themeStore.svelte';

  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');
  let success = $state(false);

  // General
  let name = $state('');
  let slug = $state('');
  let description = $state('');

  // Branding
  let logoUrl = $state('');
  let primaryColor = $state('');
  let secondaryColor = $state('');

  // Business
  let businessType = $state('');
  let locale = $state('');
  let currency = $state('');
  let timezone = $state('');

  // Preview state
  let previewPrimaryColor = $state('');
  let previewSecondaryColor = $state('');

  const businessTypes = [
    { value: 'restaurant', label: 'Restaurante' },
    { value: 'bakery', label: 'Padaria' },
    { value: 'cafe', label: 'Café' },
    { value: 'bar', label: 'Bar' },
    { value: 'food_truck', label: 'Food Truck' },
    { value: 'catering', label: 'Catering' },
    { value: 'other', label: 'Outro' }
  ];

  const locales = [
    { value: 'pt-BR', label: 'Português (Brasil)' },
    { value: 'en-US', label: 'English (United States)' },
    { value: 'es-ES', label: 'Español (España)' }
  ];

  const currencies = [
    { value: 'BRL', label: 'Real (BRL)' },
    { value: 'USD', label: 'Dólar (USD)' },
    { value: 'EUR', label: 'Euro (EUR)' }
  ];

  const timezones = [
    { value: 'America/Sao_Paulo', label: 'São Paulo (UTC-3)' },
    { value: 'America/New_York', label: 'New York (UTC-5)' },
    { value: 'Europe/London', label: 'London (UTC+0)' },
    { value: 'Asia/Tokyo', label: 'Tokyo (UTC+9)' }
  ];

  onMount(async () => {
    loading = true;
    try {
      const res = await api.companySettings.getSettings();
      if (res.error) throw new Error(res.error);
      if (res.data) {
        name = res.data.name || '';
        slug = res.data.slug || '';
        description = res.data.description || '';
        logoUrl = res.data.logo_url || '';
        primaryColor = res.data.primary_color || '';
        secondaryColor = res.data.secondary_color || '';
        businessType = res.data.business_type || '';
        locale = res.data.locale || '';
        currency = res.data.currency || '';
        timezone = res.data.timezone || '';
        
        // Initialize preview
        previewPrimaryColor = primaryColor;
        previewSecondaryColor = secondaryColor;
      }
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar configurações da empresa.';
    } finally {
      loading = false;
    }
  });

  async function saveSettings() {
    if (!name.trim()) {
      error = 'Nome da empresa é obrigatório.';
      return;
    }

    saving = true;
    error = '';
    success = false;
    try {
      const res = await api.companySettings.updateSettings({
        name: name.trim(),
        description: description.trim(),
        logo_url: logoUrl.trim(),
        primary_color: primaryColor.trim(),
        secondary_color: secondaryColor.trim(),
        business_type: businessType,
        locale: locale,
        currency: currency,
        timezone: timezone
      });
      if (res.error) throw new Error(res.error);
      success = true;
      
      // Reload theme to reflect changes
      await themeStore.loadTheme();
      
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao atualizar configurações da empresa.';
    } finally {
      saving = false;
    }
  }

  function updatePreview() {
    previewPrimaryColor = primaryColor;
    previewSecondaryColor = secondaryColor;
  }
</script>

<Workspace
  breadcrumb={[
    { label: 'Configurações', href: '/settings' },
    { label: 'Empresa' }
  ]}
  title="Configurações da Empresa"
  description="Personalize a identidade e configurações da sua empresa"
>
  {#if loading}
    <div class="skeleton-grid">
      {#each Array(3) as _}
        <Card class="skeleton-card">
          <div class="skeleton-card-header">
            <Skeleton variant="circular" width="24px" height="24px" />
            <Skeleton variant="text" width="120px" height="16px" />
          </div>
          <div class="skeleton-card-body">
            <Skeleton variant="text" width="100%" height="12px" />
            <Skeleton variant="text" width="80%" height="12px" />
          </div>
        </Card>
      {/each}
    </div>
  {:else}
    <div class="settings-grid">
      {#if error}
        <Alert variant="error" dismissible onDismiss={() => error = ''} class="full-width">
          <AlertTriangle size={16} />
          {error}
        </Alert>
      {/if}

      {#if success}
        <Alert variant="success" dismissible onDismiss={() => success = false} class="full-width">
          <Check size={16} />
          Configurações atualizadas com sucesso!
        </Alert>
      {/if}

      <!-- Dados Gerais -->
      <Card class="settings-section">
        <div class="section-header">
          <div class="section-title">
            <Building2 size={20} />
            <span>Dados Gerais</span>
          </div>
        </div>

        <form onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
          <Input
            id="name"
            label="Nome da Empresa"
            bind:value={name}
            placeholder="Nome da sua empresa"
            disabled={saving}
          />

          <Input
            id="slug"
            label="Slug"
            bind:value={slug}
            placeholder="identificador-unico"
            disabled={saving}
            readonly
          />

          <Input
            id="description"
            label="Descrição"
            bind:value={description}
            placeholder="Descrição da sua empresa"
            disabled={saving}
            multiline
            rows={3}
          />
        </form>
      </Card>

      <!-- Branding -->
      <Card class="settings-section">
        <div class="section-header">
          <div class="section-title">
            <Palette size={20} />
            <span>Branding</span>
          </div>
        </div>

        <form onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
          <Input
            id="logo-url"
            label="Logo URL"
            bind:value={logoUrl}
            placeholder="https://example.com/logo.png"
            disabled={saving}
          />

          <Input
            id="primary-color"
            label="Cor Primária"
            bind:value={primaryColor}
            placeholder="#3b82f6"
            disabled={saving}
            type="color"
            oninput={updatePreview}
          />

          <Input
            id="secondary-color"
            label="Cor Secundária"
            bind:value={secondaryColor}
            placeholder="#1e40af"
            disabled={saving}
            type="color"
            oninput={updatePreview}
          />

          <!-- Color Preview -->
          <div class="color-preview">
            <div class="preview-label">Preview:</div>
            <div class="preview-buttons">
              <Button variant="primary" style={`--primary: ${previewPrimaryColor}; --primary-hover: ${previewSecondaryColor};`}>
                Botão Primário
              </Button>
              <Button variant="secondary" style={`--secondary: ${previewSecondaryColor};`}>
                Botão Secundário
              </Button>
            </div>
          </div>
        </form>
      </Card>

      <!-- Negócio -->
      <Card class="settings-section">
        <div class="section-header">
          <div class="section-title">
            <Briefcase size={20} />
            <span>Negócio</span>
          </div>
        </div>

        <form onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
          <div class="form-group">
            <label for="business-type" class="form-label">Tipo do Negócio</label>
            <select
              id="business-type"
              bind:value={businessType}
              disabled={saving}
              class="form-select"
            >
              <option value="">Selecione...</option>
              {#each businessTypes as type}
                <option value={type.value}>{type.label}</option>
              {/each}
            </select>
          </div>

          <div class="form-group">
            <label for="locale" class="form-label">Idioma</label>
            <select
              id="locale"
              bind:value={locale}
              disabled={saving}
              class="form-select"
            >
              <option value="">Selecione...</option>
              {#each locales as loc}
                <option value={loc.value}>{loc.label}</option>
              {/each}
            </select>
          </div>

          <div class="form-group">
            <label for="currency" class="form-label">Moeda</label>
            <select
              id="currency"
              bind:value={currency}
              disabled={saving}
              class="form-select"
            >
              <option value="">Selecione...</option>
              {#each currencies as curr}
                <option value={curr.value}>{curr.label}</option>
              {/each}
            </select>
          </div>

          <div class="form-group">
            <label for="timezone" class="form-label">Timezone</label>
            <select
              id="timezone"
              bind:value={timezone}
              disabled={saving}
              class="form-select"
            >
              <option value="">Selecione...</option>
              {#each timezones as tz}
                <option value={tz.value}>{tz.label}</option>
              {/each}
            </select>
          </div>
        </form>
      </Card>

      <!-- Save Button -->
      <div class="save-section">
        <Button
          onclick={saveSettings}
          variant="primary"
          disabled={saving}
          loading={saving}
          size="lg"
        >
          Salvar Alterações
        </Button>
      </div>
    </div>
  {/if}
</Workspace>

<style>
  /* Loading State */
  .skeleton-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
  }

  .skeleton-card {
    padding: 1.5rem;
  }

  .skeleton-card-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .skeleton-card-body {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  /* Settings Grid */
  .settings-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: 1.5rem;
  }

  .full-width {
    grid-column: 1 / -1;
  }

  .settings-section {
    padding: 1.5rem;
  }

  .section-header {
    margin-bottom: 1.5rem;
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: #0f172a;
  }

  /* Color Preview */
  .color-preview {
    margin-top: 1.5rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 0.5rem;
    border: 1px solid #e2e8f0;
  }

  .preview-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #64748b;
    margin-bottom: 0.75rem;
  }

  .preview-buttons {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  /* Form Elements */
  .form-group {
    margin-bottom: 1rem;
  }

  .form-label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: #374151;
    margin-bottom: 0.5rem;
  }

  .form-select {
    width: 100%;
    padding: 0.625rem 0.875rem;
    border: 1px solid #d1d5db;
    border-radius: 0.375rem;
    font-size: 0.875rem;
    color: #374151;
    background: white;
    transition: border-color 0.15s, box-shadow 0.15s;
  }

  .form-select:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .form-select:disabled {
    background: #f3f4f6;
    cursor: not-allowed;
  }

  /* Save Section */
  .save-section {
    grid-column: 1 / -1;
    display: flex;
    justify-content: flex-end;
    padding-top: 1rem;
    border-top: 1px solid #e5e7eb;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .settings-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
