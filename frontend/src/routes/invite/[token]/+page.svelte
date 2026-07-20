<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api } from '$lib/api/client';
  import { Button, Alert, Card, Skeleton } from '$lib/components/ui';
  import { Mail, Shield, Check, AlertTriangle, X, Calendar, Building } from '@lucide/svelte';

  let loading = $state(true);
  let invitation = $state<{
    id: number;
    company_id: number;
    email: string;
    role: string;
    token: string;
    status: string;
    expires_at: string;
    accepted_at: string | null;
    created_by: number;
    created_at: string;
  } | null>(null);
  let error = $state('');
  let success = $state(false);
  let accepting = $state(false);
  let companyName = $state('');

  const roleLabels: Record<string, string> = {
    owner: 'Proprietário',
    admin: 'Administrador',
    manager: 'Gerente',
    cashier: 'Caixa',
    kitchen: 'Cozinha',
    waiter: 'Garçom'
  };

  onMount(async () => {
    const token = $page.params.token;
    if (!token) {
      error = 'Token não fornecido.';
      loading = false;
      return;
    }
    await loadInvitation(token);
  });

  async function loadInvitation(token: string) {
    loading = true;
    error = '';
    try {
      const res = await api.invitations.getByToken(token);
      if (res.error) throw new Error(res.error);
      if (res.data) {
        invitation = res.data;
        // Try to get company name (this might fail if we don't have a company endpoint, but we'll try)
        // For now, we'll just use the company ID
      }
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar convite.';
    } finally {
      loading = false;
    }
  }

  async function acceptInvitation() {
    if (!invitation) return;

    accepting = true;
    error = '';
    try {
      const res = await api.invitations.accept({ token: invitation.token });
      if (res.error) throw new Error(res.error);
      success = true;
      // Reload invitation to show updated status
      await loadInvitation(invitation.token);
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      if (e?.message?.includes('usuário não encontrado')) {
        error = 'Usuário não encontrado. Por favor, realize o cadastro primeiro.';
      } else {
        error = e?.message ?? 'Erro ao aceitar convite.';
      }
    } finally {
      accepting = false;
    }
  }

  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function isExpired(): boolean {
    if (!invitation) return true;
    return new Date(invitation.expires_at) < new Date();
  }

  function getRoleBadgeColor(role: string) {
    const colors: Record<string, string> = {
      owner: 'bg-purple-100 text-purple-700',
      admin: 'bg-blue-100 text-blue-700',
      manager: 'bg-green-100 text-green-700',
      cashier: 'bg-yellow-100 text-yellow-700',
      kitchen: 'bg-orange-100 text-orange-700',
      waiter: 'bg-pink-100 text-pink-700'
    };
    return colors[role] || 'bg-gray-100 text-gray-700';
  }
</script>

<div class="invite-page">
  <div class="invite-container">
    {#if loading}
      <Card class="invite-card">
        <div class="skeleton-content">
          <Skeleton variant="text" width="60%" height="24px" />
          <Skeleton variant="text" width="40%" height="16px" />
          <Skeleton variant="text" width="100%" height="16px" />
          <Skeleton variant="text" width="80%" height="16px" />
        </div>
      </Card>
    {:else if error}
      <Card class="invite-card error">
        <div class="error-content">
          <X size={48} class="error-icon" />
          <h2>Erro ao Carregar Convite</h2>
          <p>{error}</p>
          <Button onclick={() => window.location.href = '/login'} variant="primary">
            Ir para Login
          </Button>
        </div>
      </Card>
    {:else if invitation}
      <Card class="invite-card">
        <div class="invite-header">
          <Mail size={48} class="invite-icon" />
          <h1>Convite para Ingressar na Empresa</h1>
        </div>

        <div class="invite-details">
          <div class="detail-row">
            <Building size={20} />
            <span class="detail-label">Empresa ID:</span>
            <span class="detail-value">#{invitation.company_id}</span>
          </div>

          <div class="detail-row">
            <Mail size={20} />
            <span class="detail-label">E-mail:</span>
            <span class="detail-value">{invitation.email}</span>
          </div>

          <div class="detail-row">
            <Shield size={20} />
            <span class="detail-label">Cargo:</span>
            <span class="badge {getRoleBadgeColor(invitation.role)}">
              {roleLabels[invitation.role] || invitation.role}
            </span>
          </div>

          <div class="detail-row">
            <Calendar size={20} />
            <span class="detail-label">Expira em:</span>
            <span class="detail-value {isExpired() ? 'expired' : ''}">
              {formatDate(invitation.expires_at)}
            </span>
          </div>

          {#if invitation.accepted_at}
            <div class="detail-row">
              <Check size={20} />
              <span class="detail-label">Aceito em:</span>
              <span class="detail-value">{formatDate(invitation.accepted_at)}</span>
            </div>
          {/if}
        </div>

        <div class="invite-status">
          {#if invitation.status === 'accepted'}
            <Alert variant="success">
              <Check size={16} />
              Este convite já foi aceito! Você já faz parte da empresa.
            </Alert>
          {:else if invitation.status === 'revoked'}
            <Alert variant="error">
              <X size={16} />
              Este convite foi revogado pelo administrador da empresa.
            </Alert>
          {:else if isExpired()}
            <Alert variant="error">
              <AlertTriangle size={16} />
              Este convite expirou. Entre em contato com o administrador da empresa para solicitar um novo convite.
            </Alert>
          {:else}
            <Alert variant="info">
              <Mail size={16} />
              Você foi convidado a ingressar na empresa. Aceite o convite para começar a colaborar.
            </Alert>
          {/if}
        </div>

        <div class="invite-actions">
          {#if invitation.status === 'pending' && !isExpired()}
            <Button onclick={acceptInvitation} variant="primary" disabled={accepting} loading={accepting}>
              <Check size={16} />
              Aceitar Convite
            </Button>
          {/if}

          <Button onclick={() => window.location.href = '/login'} variant="secondary">
            Ir para Login
          </Button>

          {#if error && error.includes('usuário não encontrado')}
            <Button onclick={() => window.location.href = '/register'} variant="ghost">
              Realizar Cadastro
            </Button>
          {/if}
        </div>

        {#if success}
          <div class="success-message">
            <Check size={20} />
            <span>Convite aceito com sucesso! Você já faz parte da empresa.</span>
          </div>
        {/if}
      </Card>
    {/if}
  </div>
</div>

<style>
  .invite-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .invite-container {
    width: 100%;
    max-width: 500px;
  }

  .invite-card {
    padding: 2rem;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  }

  .invite-card.error {
    background: #fef2f2;
    border: 1px solid #fecaca;
  }

  .skeleton-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .error-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    text-align: center;
  }

  .error-icon {
    color: #ef4444;
  }

  .error-content h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #374151;
  }

  .error-content p {
    color: #64748b;
  }

  .invite-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    margin-bottom: 2rem;
    text-align: center;
  }

  .invite-icon {
    color: #6366f1;
  }

  .invite-header h1 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #0f172a;
  }

  .invite-details {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .detail-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    background: #f9fafb;
    border-radius: 0.5rem;
  }

  .detail-row svg {
    color: #64748b;
    flex-shrink: 0;
  }

  .detail-label {
    color: #64748b;
    font-size: 0.875rem;
  }

  .detail-value {
    color: #374151;
    font-weight: 500;
    margin-left: auto;
  }

  .detail-value.expired {
    color: #ef4444;
  }

  .badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .invite-status {
    margin-bottom: 1.5rem;
  }

  .invite-actions {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .success-message {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 1rem;
    background: #f0fdf4;
    border: 1px solid #bbf7d0;
    border-radius: 0.5rem;
    color: #166534;
    font-weight: 500;
    margin-top: 1rem;
  }

  @media (max-width: 640px) {
    .invite-page {
      padding: 1rem;
    }

    .invite-card {
      padding: 1.5rem;
    }

    .invite-header h1 {
      font-size: 1.25rem;
    }
  }
</style>
