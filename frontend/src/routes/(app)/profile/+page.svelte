<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';

  let loading = $state(true);
  let saving = $state(false);
  let savingPassword = $state(false);
  let error = $state('');
  let passwordError = $state('');
  let success = $state(false);
  let passwordSuccess = $state(false);

  let name = $state('');
  let email = $state('');

  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');

  onMount(async () => {
    loading = true;
    try {
      const res = await api.auth.me();
      if (res.error) throw new Error(res.error);
      if (res.data) {
        name = res.data.name;
        email = res.data.email;
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

    saving = true;
    error = '';
    success = false;
    try {
      const res = await api.auth.updateProfile({ name: name.trim(), email: email.trim() });
      if (res.error) throw new Error(res.error);
      success = true;
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
</script>

<div class="page-wrapper">
  <a href="/dashboard" class="back-link">← Voltar ao Dashboard</a>

  <div class="profile-container">
    <header class="profile-header">
      <h1 class="profile-title">Meu Perfil</h1>
      <p class="profile-subtitle">Gerencie suas informações de conta</p>
    </header>

    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
        <span>Carregando perfil…</span>
      </div>

    {:else}
      <div class="profile-card">
        {#if error}
          <div class="alert alert-error">
            <span>⚠️ {error}</span>
          </div>
        {/if}

        {#if success}
          <div class="alert alert-success">
            <span>✓ Perfil atualizado com sucesso!</span>
          </div>
        {/if}

        <form onsubmit={(e) => { e.preventDefault(); saveProfile(); }}>
          <div class="form-group">
            <label for="name">Nome</label>
            <input
              id="name"
              type="text"
              bind:value={name}
              placeholder="Seu nome completo"
              disabled={saving}
              required
            />
          </div>

          <div class="form-group">
            <label for="email">E-mail</label>
            <input
              id="email"
              type="email"
              bind:value={email}
              placeholder="seu@email.com"
              disabled={saving}
              required
            />
          </div>

          <div class="form-actions">
            <button
              type="submit"
              class="btn btn-primary"
              disabled={saving}
            >
              {saving ? 'Salvando…' : 'Salvar Alterações'}
            </button>
          </div>
        </form>
      </div>

      <div class="password-card">
        <h2 class="card-title">Alterar Senha</h2>

        {#if passwordError}
          <div class="alert alert-error">
            <span>⚠️ {passwordError}</span>
          </div>
        {/if}

        {#if passwordSuccess}
          <div class="alert alert-success">
            <span>✓ Senha alterada com sucesso!</span>
          </div>
        {/if}

        <form onsubmit={(e) => { e.preventDefault(); changePassword(); }}>
          <div class="form-group">
            <label for="current-password">Senha Atual</label>
            <input
              id="current-password"
              type="password"
              bind:value={currentPassword}
              placeholder="Digite sua senha atual"
              disabled={savingPassword}
              required
            />
          </div>

          <div class="form-group">
            <label for="new-password">Nova Senha</label>
            <input
              id="new-password"
              type="password"
              bind:value={newPassword}
              placeholder="Mínimo 6 caracteres"
              disabled={savingPassword}
              required
              minlength="6"
            />
          </div>

          <div class="form-group">
            <label for="confirm-password">Confirmar Nova Senha</label>
            <input
              id="confirm-password"
              type="password"
              bind:value={confirmPassword}
              placeholder="Repita a nova senha"
              disabled={savingPassword}
              required
              minlength="6"
            />
          </div>

          <div class="form-actions">
            <button
              type="submit"
              class="btn btn-secondary"
              disabled={savingPassword}
            >
              {savingPassword ? 'Alterando…' : 'Alterar Senha'}
            </button>
          </div>
        </form>
      </div>
    {/if}
  </div>
</div>

<style>
  .page-wrapper {
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem 1.5rem;
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    color: var(--color-muted);
    font-size: 0.875rem;
    text-decoration: none;
    margin-bottom: 1.75rem;
  }

  .back-link:hover {
    color: var(--color-text);
  }

  .profile-container {
    background: var(--color-surface, #fff);
    border-radius: 1rem;
    padding: 2rem;
    border: 1px solid var(--color-border, #e5e7eb);
  }

  .profile-header {
    margin-bottom: 2rem;
  }

  .profile-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 0.5rem;
  }

  .profile-subtitle {
    color: var(--color-muted);
    font-size: 0.9rem;
    margin: 0;
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 3rem;
    color: var(--color-muted);
  }

  .spinner {
    width: 2rem;
    height: 2rem;
    border: 3px solid var(--color-border, #e5e7eb);
    border-top-color: var(--color-primary, #e85d04);
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .profile-card,
  .password-card {
    max-width: 400px;
    margin-bottom: 2rem;
  }

  .card-title {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 1.25rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    margin-bottom: 1.25rem;
  }

  .form-group label {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--color-text);
  }

  .form-group input {
    border: 1px solid var(--color-border, #d1d5db);
    border-radius: 0.5rem;
    padding: 0.6rem 0.75rem;
    font-size: 0.9rem;
    background: var(--color-surface, #fff);
    color: var(--color-text);
    transition: border-color 0.15s;
    width: 100%;
    box-sizing: border-box;
  }

  .form-group input:focus {
    outline: none;
    border-color: var(--color-primary, #e85d04);
    box-shadow: 0 0 0 3px rgba(232,93,4,0.12);
  }

  .form-group input:disabled {
    background: var(--color-surface-2, #f3f4f6);
    cursor: not-allowed;
  }

  .form-actions {
    margin-top: 1.5rem;
  }

  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.6rem 1.25rem;
    border-radius: 0.5rem;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    border: none;
    transition: background 0.15s;
    width: 100%;
  }

  .btn:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

  .btn-primary {
    background: var(--color-primary, #e85d04);
    color: #fff;
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--color-primary-dark, #c84e00);
  }

  .alert {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.85rem 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }

  .alert-error {
    background: #fef2f2;
    border: 1px solid #fca5a5;
    color: #b91c1c;
  }

  .alert-success {
    background: #f0fdf4;
    border: 1px solid #86efac;
    color: #15803d;
  }
</style>
