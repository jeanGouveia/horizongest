<script lang="ts">
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import { userStore } from '$lib/stores/userStore.svelte';

  let email    = $state('');
  let password = $state('');
  let error    = $state('');
  let loading  = $state(false);

  async function handleSubmit() {
    error   = '';
    loading = true;

    const { data, error: err } = await api.auth.login({ email, password });
    loading = false;

    if (err || !data) {
      // Melhorar mensagem de erro específica
      if (err?.includes('e-mail ou senha inválidos')) {
        error = 'E-mail ou senha incorretos';
      } else if (err?.includes('unauthorized')) {
        error = 'Não autorizado';
      } else {
        error = err ?? 'Erro ao fazer login';
      }
      return;
    }

    userStore.setUser({ id: data.id, name: data.name, email: data.email });
    goto('/dashboard');
  }
</script>

<div class="auth-container">
  <div class="auth-card">
    <h1>🍽️ Prato Online</h1>
    <h2>Entrar</h2>

    {#if error}
      <div class="alert-error">{error}</div>
    {/if}

    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      <label>
        E-mail
        <input
          type="email"
          bind:value={email}
          placeholder="voce@email.com"
          required
          autocomplete="email"
        />
      </label>

      <label>
        Senha
        <input
          type="password"
          bind:value={password}
          placeholder="••••••"
          required
          autocomplete="current-password"
        />
      </label>

      <button type="submit" disabled={loading}>
        {loading ? 'Entrando...' : 'Entrar'}
      </button>
    </form>

    <p class="auth-link">
      Não tem conta? <a href="/register">Cadastrar</a>
    </p>
  </div>
</div>

<style>
  .auth-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f5f5f0;
  }
  .auth-card {
    background: white;
    padding: 2.5rem;
    border-radius: 12px;
    box-shadow: 0 2px 16px rgba(0,0,0,0.08);
    width: 100%;
    max-width: 400px;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  h1 { font-size: 1.4rem; color: #e85d04; margin: 0; }
  h2 { font-size: 1.1rem; color: #333; margin: 0; font-weight: 500; }
  form { display: flex; flex-direction: column; gap: 1rem; }
  label {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #555;
  }
  input {
    padding: 0.65rem 0.875rem;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
    outline: none;
    transition: border-color 0.2s;
  }
  input:focus { border-color: #e85d04; }
  button {
    padding: 0.75rem;
    background: #e85d04;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    margin-top: 0.5rem;
    transition: background 0.2s;
  }
  button:hover:not(:disabled) { background: #c94f03; }
  button:disabled { opacity: 0.6; cursor: not-allowed; }
  .alert-error {
    background: #fff0f0;
    border: 1px solid #fca5a5;
    color: #dc2626;
    padding: 0.75rem;
    border-radius: 8px;
    font-size: 0.875rem;
  }
  .auth-link { font-size: 0.875rem; text-align: center; color: #666; }
  .auth-link a { color: #e85d04; text-decoration: none; font-weight: 500; }
</style>
