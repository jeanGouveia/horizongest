<script lang="ts">
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import { userStore } from '$lib/stores/userStore.svelte';
  import { Button, Input, Alert, PageContainer } from '$lib/components/ui';

  let name     = $state('');
  let email    = $state('');
  let password = $state('');
  let error    = $state('');
  let loading  = $state(false);

  function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  async function handleSubmit() {
    error = '';

    if (!isValidEmail(email)) {
      error = 'Por favor, insira um e-mail válido';
      return;
    }

    if (password.length < 6) {
      error = 'A senha deve ter no mínimo 6 caracteres';
      return;
    }

    loading = true;

    const { data, error: err } = await api.auth.register({ name, email, password });
    loading = false;

    if (err || !data) {
      if (err?.includes('e-mail já cadastrado')) {
        error = 'Este e-mail já está cadastrado';
      } else if (err?.includes('dados inválidos')) {
        error = 'Verifique os campos marcados em vermelho';
      } else {
        error = err ?? 'Erro ao cadastrar';
      }
      return;
    }

    // Login automático após registro bem-sucedido
    const loginRes = await api.auth.login({ email, password });
    if (loginRes.data) {
      userStore.setUser({
        id:    loginRes.data.id,
        name:  loginRes.data.name,
        email: loginRes.data.email
      });
      goto('/dashboard');
    } else {
      goto('/login');
    }
  }
</script>

<PageContainer maxWidth="sm">
  <div class="auth-container">
    <div class="auth-card">
      <h1 class="auth-title">🍽️ Prato Online</h1>
      <h2 class="auth-subtitle">Criar conta</h2>

      {#if error}
        <Alert variant="error" dismissible onDismiss={() => error = ''}>
          {error}
        </Alert>
      {/if}

      <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
        <Input
          id="name"
          label="Nome"
          type="text"
          bind:value={name}
          placeholder="Seu nome completo"
          required
          autocomplete="name"
        />

        <Input
          id="email"
          label="E-mail"
          type="email"
          bind:value={email}
          placeholder="voce@email.com"
          required
          autocomplete="email"
        />

        <Input
          id="password"
          label="Senha"
          type="password"
          bind:value={password}
          placeholder="••••••"
          required
          minlength={6}
          autocomplete="new-password"
          helper="Mínimo 6 caracteres"
        />

        <Button type="submit" variant="primary" fullWidth loading={loading}>
          {loading ? 'Cadastrando...' : 'Criar conta'}
        </Button>
      </form>

      <p class="auth-link">
        Já tem conta? <a href="/login">Entrar</a>
      </p>
    </div>
  </div>
</PageContainer>

<style>
  .auth-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  }

  .auth-card {
    background: #ffffff;
    padding: 2.5rem;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(99, 102, 241, 0.1);
    width: 100%;
    max-width: 420px;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    border: 1px solid #e2e8f0;
  }

  .auth-title {
    font-size: 1.5rem;
    color: #6366f1;
    margin: 0;
    font-weight: 700;
    text-align: center;
  }

  .auth-subtitle {
    font-size: 1.125rem;
    color: #0f172a;
    margin: 0;
    font-weight: 600;
    text-align: center;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .auth-link {
    font-size: 0.875rem;
    text-align: center;
    color: #64748b;
  }

  .auth-link a {
    color: #6366f1;
    text-decoration: none;
    font-weight: 600;
  }

  .auth-link a:hover {
    text-decoration: underline;
  }
</style>
