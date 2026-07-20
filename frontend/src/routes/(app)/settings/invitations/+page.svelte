<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { Button, Input, Alert, Card, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Mail, Copy, Trash2, Check, AlertTriangle, Plus, Clock, X } from '@lucide/svelte';
  import { rbacStore } from '$lib/stores/rbacStore.svelte';

  let loading = $state(true);
  let invitations = $state<Array<{
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
  }>>([]);
  let error = $state('');
  let success = $state(false);
  
  // New invitation modal
  let showNewModal = $state(false);
  let creating = $state(false);
  let newEmail = $state('');
  let newRole = $state('manager');
  let newError = $state('');
  
  const roles = [
    { value: 'owner', label: 'Proprietário' },
    { value: 'admin', label: 'Administrador' },
    { value: 'manager', label: 'Gerente' },
    { value: 'cashier', label: 'Caixa' },
    { value: 'kitchen', label: 'Cozinha' },
    { value: 'waiter', label: 'Garçom' }
  ];

  const roleLabels: Record<string, string> = {
    owner: 'Proprietário',
    admin: 'Administrador',
    manager: 'Gerente',
    cashier: 'Caixa',
    kitchen: 'Cozinha',
    waiter: 'Garçom'
  };

  const statusLabels: Record<string, string> = {
    pending: 'Pendente',
    accepted: 'Aceito',
    expired: 'Expirado',
    revoked: 'Revogado'
  };

  onMount(async () => {
    await rbacStore.load();
    await loadInvitations();
  });

  async function loadInvitations() {
    loading = true;
    error = '';
    try {
      const res = await api.companyInvitations.list();
      if (res.error) throw new Error(res.error);
      if (res.data) {
        invitations = res.data;
      }
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar convites.';
    } finally {
      loading = false;
    }
  }

  async function createInvitation() {
    if (!newEmail.trim()) {
      newError = 'E-mail é obrigatório.';
      return;
    }

    creating = true;
    newError = '';
    try {
      const res = await api.companyInvitations.create({ email: newEmail.trim(), role: newRole });
      if (res.error) throw new Error(res.error);
      success = true;
      showNewModal = false;
      newEmail = '';
      newRole = 'manager';
      await loadInvitations();
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      newError = e?.message ?? 'Erro ao criar convite.';
    } finally {
      creating = false;
    }
  }

  async function revokeInvitation(id: number) {
    if (!confirm('Tem certeza que deseja revogar este convite?')) {
      return;
    }

    try {
      const res = await api.companyInvitations.revoke(id);
      if (res.error) throw new Error(res.error);
      success = true;
      await loadInvitations();
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao revogar convite.';
    }
  }

  async function copyLink(token: string) {
    const link = `${window.location.origin}/invite/${token}`;
    try {
      await navigator.clipboard.writeText(link);
      success = true;
      setTimeout(() => (success = false), 3000);
    } catch (e) {
      error = 'Não foi possível copiar o link.';
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

  function getStatusBadge(status: string) {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-700',
      accepted: 'bg-green-100 text-green-700',
      expired: 'bg-gray-100 text-gray-700',
      revoked: 'bg-red-100 text-red-700'
    };
    return colors[status] || 'bg-gray-100 text-gray-700';
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

  function isExpired(expiresAt: string): boolean {
    return new Date(expiresAt) < new Date();
  }
</script>

<Workspace
  breadcrumb={[
    { label: 'Configurações', href: '/settings' },
    { label: 'Convites' }
  ]}
  title="Convites da Empresa"
  description="Gerencie os convites para novos usuários ingressarem na sua empresa"
>
  {#if loading}
    <div class="skeleton-grid">
      {#each Array(5) as _}
        <Card class="skeleton-card">
          <div class="skeleton-card-body">
            <Skeleton variant="text" width="100%" height="16px" />
            <Skeleton variant="text" width="80%" height="16px" />
          </div>
        </Card>
      {/each}
    </div>
  {:else}
    <div class="invitations-container">
      {#if error}
        <Alert variant="error" dismissible onDismiss={() => error = ''} class="full-width">
          <AlertTriangle size={16} />
          {error}
        </Alert>
      {/if}

      {#if success}
        <Alert variant="success" dismissible onDismiss={() => success = false} class="full-width">
          <Check size={16} />
          Operação realizada com sucesso!
        </Alert>
      {/if}

      <!-- Actions Bar -->
      <div class="actions-bar">
        <Button onclick={() => showNewModal = true} variant="primary">
          <Plus size={16} />
          Novo Convite
        </Button>
      </div>

      <!-- Invitations List -->
      <Card class="invitations-list">
        {#if invitations.length === 0}
          <div class="empty-state">
            <Mail size={48} />
            <p>Nenhum convite enviado.</p>
            <Button onclick={() => showNewModal = true} variant="secondary">
              <Plus size={16} />
              Enviar Primeiro Convite
            </Button>
          </div>
        {:else}
          <div class="invitations-grid">
            {#each invitations as invitation}
              <div class="invitation-card {invitation.status}">
                <div class="invitation-header">
                  <div class="invitation-email">
                    <Mail size={16} />
                    {invitation.email}
                  </div>
                  <span class="badge {getStatusBadge(invitation.status)}">
                    {statusLabels[invitation.status] || invitation.status}
                  </span>
                </div>
                
                <div class="invitation-details">
                  <div class="detail-item">
                    <span class="detail-label">Cargo:</span>
                    <span class="badge {getRoleBadgeColor(invitation.role)}">
                      {roleLabels[invitation.role] || invitation.role}
                    </span>
                  </div>
                  
                  <div class="detail-item">
                    <Clock size={14} />
                    <span class="detail-label">Expira em:</span>
                    <span class="detail-value {isExpired(invitation.expires_at) ? 'expired' : ''}">
                      {formatDate(invitation.expires_at)}
                    </span>
                  </div>
                  
                  {#if invitation.accepted_at}
                    <div class="detail-item">
                      <Check size={14} />
                      <span class="detail-label">Aceito em:</span>
                      <span class="detail-value">{formatDate(invitation.accepted_at)}</span>
                    </div>
                  {/if}
                </div>

                <div class="invitation-actions">
                  {#if invitation.status === 'pending' && !isExpired(invitation.expires_at)}
                    <Button onclick={() => copyLink(invitation.token)} variant="ghost" size="sm">
                      <Copy size={16} />
                      Copiar Link
                    </Button>
                    <Button onclick={() => revokeInvitation(invitation.id)} variant="ghost" size="sm">
                      <X size={16} />
                      Revogar
                    </Button>
                  {:else}
                    <span class="no-actions">
                      {invitation.status === 'expired' ? 'Convite expirado' : 
                       invitation.status === 'revoked' ? 'Convite revogado' : 
                       invitation.status === 'accepted' ? 'Convite aceito' : ''}
                    </span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </Card>
    </div>
  {/if}
</Workspace>

<!-- New Invitation Modal -->
{#if showNewModal}
  <div class="modal-overlay" onclick={() => showNewModal = false}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Novo Convite</h3>
        <Button onclick={() => showNewModal = false} variant="ghost" size="sm">✕</Button>
      </div>
      <div class="modal-body">
        {#if newError}
          <Alert variant="error" dismissible onDismiss={() => newError = ''}>
            <AlertTriangle size={16} />
            {newError}
          </Alert>
        {/if}
        <Input
          id="new-invitation-email"
          label="E-mail do Convidado"
          bind:value={newEmail}
          placeholder="convidado@exemplo.com"
          disabled={creating}
        />
        <div class="form-group">
          <label for="role-select" class="form-label">Cargo</label>
          <select
            id="role-select"
            bind:value={newRole}
            disabled={creating}
            class="form-select"
          >
            {#each roles as role}
              <option value={role.value}>{role.label}</option>
            {/each}
          </select>
        </div>
        <p class="help-text">
          O convite será válido por 7 dias. O usuário precisará estar cadastrado no sistema para aceitar o convite.
        </p>
      </div>
      <div class="modal-footer">
        <Button onclick={() => showNewModal = false} variant="secondary" disabled={creating}>
          Cancelar
        </Button>
        <Button onclick={createInvitation} variant="primary" disabled={creating} loading={creating}>
          Criar Convite
        </Button>
      </div>
    </div>
  </div>
{/if}

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

  .skeleton-card-body {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  /* Invitations Container */
  .invitations-container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .full-width {
    width: 100%;
  }

  /* Actions Bar */
  .actions-bar {
    display: flex;
    justify-content: flex-end;
  }

  /* Invitations List */
  .invitations-list {
    padding: 1.5rem;
  }

  .invitations-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 1.5rem;
  }

  .invitation-card {
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    padding: 1.5rem;
    transition: all 0.2s;
  }

  .invitation-card:hover {
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .invitation-card.accepted {
    border-color: #10b981;
    background: #f0fdf4;
  }

  .invitation-card.expired {
    border-color: #9ca3af;
    background: #f9fafb;
    opacity: 0.7;
  }

  .invitation-card.revoked {
    border-color: #ef4444;
    background: #fef2f2;
    opacity: 0.7;
  }

  .invitation-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .invitation-email {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 500;
    color: #374151;
  }

  .badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .invitation-details {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .detail-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .detail-label {
    color: #64748b;
  }

  .detail-value {
    color: #374151;
    font-weight: 500;
  }

  .detail-value.expired {
    color: #ef4444;
  }

  .invitation-actions {
    display: flex;
    gap: 0.5rem;
    padding-top: 1rem;
    border-top: 1px solid #e5e7eb;
  }

  .no-actions {
    font-size: 0.875rem;
    color: #64748b;
    font-style: italic;
  }

  /* Empty State */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    padding: 4rem 2rem;
    color: #64748b;
  }

  .empty-state p {
    font-size: 1rem;
  }

  /* Modal */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: white;
    border-radius: 0.5rem;
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .modal-header h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: #0f172a;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1.5rem;
    border-top: 1px solid #e5e7eb;
  }

  .help-text {
    font-size: 0.875rem;
    color: #64748b;
    margin-top: 0.5rem;
  }

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
  }

  .form-select:disabled {
    background: #f3f4f6;
    cursor: not-allowed;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .invitations-grid {
      grid-template-columns: 1fr;
    }

    .invitation-actions {
      flex-direction: column;
    }
  }
</style>
