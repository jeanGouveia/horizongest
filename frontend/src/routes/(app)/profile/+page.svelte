<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import { userStore } from '$lib/stores/userStore.svelte';
  import { Button, Input, Alert, Card, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { User, Lock, Settings, LogOut, Activity, Check, AlertTriangle } from '@lucide/svelte';

  let loading = $state(true);
  let saving = $state(false);
  let savingPassword = $state(false);
  let error = $state('');
  let passwordError = $state('');
  let success = $state(false);
  let passwordSuccess = $state(false);

  let name = $state('');
  let email = $state('');
  let originalEmail = $state('');

  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');

  let profilePassword = $state('');

  onMount(async () => {
    loading = true;
    try {
      const res = await api.auth.me();
      if (res.error) throw new Error(res.error);
      if (res.data) {
        name = res.data.name;
        email = res.data.email;
        originalEmail = res.data.email;
      }
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar perfil.';
    } finally {
      loading = false;
    }
  });

  async function saveProfile() {
    if (!name.trim() || !email.trim()) {
      error = 'Nome e e-mail são obrigatórios.';
      return;
    }

    // Se e-mail foi alterado, exigir confirmação de senha
    if (email.trim() !== originalEmail) {
      if (!profilePassword || profilePassword.length < 6) {
        error = 'A senha deve ter no mínimo 6 caracteres.';
        return;
      }
    }

    saving = true;
    error = '';
    success = false;
    try {
      const body: { name: string; email: string; current_password?: string } = {
        name: name.trim(),
        email: email.trim()
      };
      // Include current_password only when email is being changed
      if (email.trim() !== originalEmail) {
        body.current_password = profilePassword;
      }
      const res = await api.auth.updateProfile(body);
      if (res.error) throw new Error(res.error);
      success = true;
      originalEmail = email.trim();
      profilePassword = '';
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao atualizar perfil.';
    } finally {
      saving = false;
    }
  }

  async function changePassword() {
    if (!currentPassword || !newPassword || !confirmPassword) {
      passwordError = 'Todos os campos de senha são obrigatórios.';
      return;
    }

    if (newPassword.length < 6) {
      passwordError = 'A nova senha deve ter no mínimo 6 caracteres.';
      return;
    }

    if (newPassword !== confirmPassword) {
      passwordError = 'A nova senha e a confirmação não coincidem.';
      return;
    }

    savingPassword = true;
    passwordError = '';
    passwordSuccess = false;
    try {
      const res = await api.auth.changePassword({
        current_password: currentPassword,
        new_password: newPassword
      });
      if (res.error) throw new Error(res.error);
      passwordSuccess = true;
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
      setTimeout(() => (passwordSuccess = false), 3000);
    } catch (e: any) {
      passwordError = e?.message ?? 'Erro ao alterar senha.';
    } finally {
      savingPassword = false;
    }
  }

  async function logout() {
    try {
      await api.auth.logout();
      userStore.logout();
      goto('/login');
    } catch (e: any) {
      error = e?.message ?? 'Erro ao fazer logout.';
    }
  }
</script>

<Workspace
  breadcrumb={[{ label: 'Perfil' }]}
  title="Perfil do Usuário"
  description="Gerencie suas informações e configurações de conta"
>
  {#if loading}
    <div class="skeleton-grid">
      {#each Array(4) as _}
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
    <div class="profile-grid">
      <!-- Informações Pessoais -->
      <Card class="profile-section">
        <div class="section-header">
          <div class="section-title">
            <User size={20} />
            <span>Informações Pessoais</span>
          </div>
        </div>

        {#if error}
          <Alert variant="error" dismissible onDismiss={() => error = ''}>
            <AlertTriangle size={16} />
            {error}
          </Alert>
        {/if}

        {#if success}
          <Alert variant="success" dismissible onDismiss={() => success = false}>
            <Check size={16} />
            Perfil atualizado com sucesso!
          </Alert>
        {/if}

        <form onsubmit={(e) => { e.preventDefault(); saveProfile(); }}>
          <Input
            id="name"
            label="Nome"
            bind:value={name}
            placeholder="Seu nome completo"
            disabled={saving}
          />

          <Input
            id="email"
            label="E-mail"
            type="email"
            bind:value={email}
            placeholder="seu@email.com"
            disabled={saving}
          />

          {#if email !== originalEmail}
            <Input
              id="profile-password"
              label="Senha atual (obrigatório para alterar e-mail)"
              type="password"
              bind:value={profilePassword}
              placeholder="Confirme sua senha atual"
              disabled={saving}
            />
          {/if}

          <div class="form-actions">
            <Button
              type="submit"
              variant="primary"
              disabled={saving}
              loading={saving}
            >
              Salvar Alterações
            </Button>
          </div>
        </form>
      </Card>

      <!-- Segurança -->
      <Card class="profile-section">
        <div class="section-header">
          <div class="section-title">
            <Lock size={20} />
            <span>Segurança</span>
          </div>
        </div>

        {#if passwordError}
          <Alert variant="error" dismissible onDismiss={() => passwordError = ''}>
            <AlertTriangle size={16} />
            {passwordError}
          </Alert>
        {/if}

        {#if passwordSuccess}
          <Alert variant="success" dismissible onDismiss={() => passwordSuccess = false}>
            <Check size={16} />
            Senha alterada com sucesso!
          </Alert>
        {/if}

        <form onsubmit={(e) => { e.preventDefault(); changePassword(); }}>
          <Input
            id="current-password"
            label="Senha Atual"
            type="password"
            bind:value={currentPassword}
            placeholder="Digite sua senha atual"
            disabled={savingPassword}
          />

          <Input
            id="new-password"
            label="Nova Senha"
            type="password"
            bind:value={newPassword}
            placeholder="Mínimo 6 caracteres"
            disabled={savingPassword}
            minlength={6}
          />

          <Input
            id="confirm-password"
            label="Confirmar Nova Senha"
            type="password"
            bind:value={confirmPassword}
            placeholder="Repita a nova senha"
            disabled={savingPassword}
            minlength={6}
          />

          <div class="form-actions">
            <Button
              type="submit"
              variant="secondary"
              disabled={savingPassword}
              loading={savingPassword}
            >
              Alterar Senha
            </Button>
          </div>
        </form>
      </Card>

      <!-- Preferências -->
      <Card class="profile-section">
        <div class="section-header">
          <div class="section-title">
            <Settings size={20} />
            <span>Preferências</span>
          </div>
        </div>

        <div class="preferences-content">
          <p class="preferences-note">
            Configurações adicionais estarão disponíveis em breve.
          </p>
        </div>
      </Card>

      <!-- Sessão -->
      <Card class="profile-section">
        <div class="section-header">
          <div class="section-title">
            <LogOut size={20} />
            <span>Sessão</span>
          </div>
        </div>

        <div class="session-content">
          <p class="session-note">
            Encerre sua sessão atual para sair da aplicação.
          </p>
          <Button onclick={logout} variant="danger">
            <LogOut size={16} />
            Sair da Conta
          </Button>
        </div>
      </Card>
    </div>
  {/if}
</Workspace>

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

  /* Profile Grid */
  .profile-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: 1.5rem;
  }

  .profile-section {
    padding: 1.5rem;
    transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .profile-section:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px 0 rgb(0 0 0 / 0.08);
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

  .form-actions {
    margin-top: 1.5rem;
  }

  /* Skeleton Loading */
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

  /* Preferences & Session */
  .preferences-content,
  .session-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .preferences-note,
  .session-note {
    font-size: 0.875rem;
    color: #64748b;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .profile-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
