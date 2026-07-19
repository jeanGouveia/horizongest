<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { Button, Input, Alert, Card, Skeleton } from '$lib/components/ui';
  import { Workspace } from '$lib/components/layout';
  import { Users, UserPlus, Shield, Trash2, Check, AlertTriangle, MoreVertical } from '@lucide/svelte';
  import { rbacStore } from '$lib/stores/rbacStore.svelte';

  let loading = $state(true);
  let users = $state<Array<{
    id: number;
    name: string;
    email: string;
    role: string | null;
    active: boolean;
    company_id: number | null;
  }>>([]);
  let error = $state('');
  let success = $state(false);
  
  // Add user modal
  let showAddModal = $state(false);
  let addingUser = $state(false);
  let addUserEmail = $state('');
  let addUserError = $state('');
  
  // Change role modal
  let showRoleModal = $state(false);
  let changingRole = $state(false);
  let selectedUser = $state<{
    id: number;
    name: string;
    email: string;
    role: string | null;
    active: boolean;
    company_id: number | null;
  } | null>(null);
  let newRole = $state('');
  let roleError = $state('');

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

  onMount(async () => {
    await rbacStore.load();
    await loadUsers();
  });

  async function loadUsers() {
    loading = true;
    error = '';
    try {
      const res = await api.companyUsers.list();
      if (res.error) throw new Error(res.error);
      if (res.data) {
        users = res.data;
      }
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar usuários.';
    } finally {
      loading = false;
    }
  }

  async function addUser() {
    if (!addUserEmail.trim()) {
      addUserError = 'E-mail é obrigatório.';
      return;
    }

    addingUser = true;
    addUserError = '';
    try {
      const res = await api.companyUsers.add({ email: addUserEmail.trim() });
      if (res.error) throw new Error(res.error);
      success = true;
      showAddModal = false;
      addUserEmail = '';
      await loadUsers();
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      addUserError = e?.message ?? 'Erro ao adicionar usuário.';
    } finally {
      addingUser = false;
    }
  }

  async function changeRole() {
    if (!newRole) {
      roleError = 'Cargo é obrigatório.';
      return;
    }

    if (!selectedUser) {
      roleError = 'Nenhum usuário selecionado.';
      return;
    }

    changingRole = true;
    roleError = '';
    try {
      const res = await api.companyUsers.changeRole(selectedUser.id, { role: newRole });
      if (res.error) throw new Error(res.error);
      success = true;
      showRoleModal = false;
      selectedUser = null;
      newRole = '';
      await loadUsers();
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      roleError = e?.message ?? 'Erro ao alterar cargo.';
    } finally {
      changingRole = false;
    }
  }

  async function removeUser(userId: number) {
    if (!confirm('Tem certeza que deseja remover este usuário da empresa?')) {
      return;
    }

    try {
      const res = await api.companyUsers.remove(userId);
      if (res.error) throw new Error(res.error);
      success = true;
      await loadUsers();
      setTimeout(() => (success = false), 3000);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao remover usuário.';
    }
  }

  function openRoleModal(user: any) {
    selectedUser = user;
    newRole = user.role || '';
    showRoleModal = true;
    roleError = '';
  }

  function getRoleBadge(role: string | null) {
    if (!role) return 'Sem cargo';
    return roleLabels[role] || role;
  }

  function getRoleBadgeColor(role: string | null) {
    if (!role) return 'bg-gray-100 text-gray-700';
    switch (role.toLowerCase()) {
      case 'owner': return 'bg-purple-100 text-purple-700';
      case 'admin': return 'bg-blue-100 text-blue-700';
      case 'manager': return 'bg-green-100 text-green-700';
      case 'cashier': return 'bg-yellow-100 text-yellow-700';
      case 'kitchen': return 'bg-orange-100 text-orange-700';
      case 'waiter': return 'bg-pink-100 text-pink-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  }
</script>

<Workspace
  breadcrumb={[
    { label: 'Configurações', href: '/settings' },
    { label: 'Usuários' }
  ]}
  title="Gerenciamento de Usuários"
  description="Gerencie os usuários e cargos da sua empresa"
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
    <div class="users-container">
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
        <Button onclick={() => showAddModal = true} variant="primary">
          <UserPlus size={16} />
          Adicionar Usuário
        </Button>
      </div>

      <!-- Users Table -->
      <Card class="users-table">
        {#if users.length === 0}
          <div class="empty-state">
            <Users size={48} />
            <p>Nenhum usuário na empresa.</p>
            <Button onclick={() => showAddModal = true} variant="secondary">
              <UserPlus size={16} />
              Adicionar Primeiro Usuário
            </Button>
          </div>
        {:else}
          <table class="table">
            <thead>
              <tr>
                <th>Nome</th>
                <th>E-mail</th>
                <th>Cargo</th>
                <th>Status</th>
                <th>Ações</th>
              </tr>
            </thead>
            <tbody>
              {#each users as user}
                <tr>
                  <td>{user.name}</td>
                  <td>{user.email}</td>
                  <td>
                    <span class="badge {getRoleBadgeColor(user.role)}">
                      {getRoleBadge(user.role)}
                    </span>
                  </td>
                  <td>
                    <span class="status-badge {user.active ? 'active' : 'inactive'}">
                      {user.active ? 'Ativo' : 'Inativo'}
                    </span>
                  </td>
                  <td>
                    <div class="actions">
                      <Button onclick={() => openRoleModal(user)} variant="ghost" size="sm">
                        <Shield size={16} />
                      </Button>
                      <Button onclick={() => removeUser(user.id)} variant="ghost" size="sm">
                        <Trash2 size={16} />
                      </Button>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </Card>
    </div>
  {/if}
</Workspace>

<!-- Add User Modal -->
{#if showAddModal}
  <div class="modal-overlay" onclick={() => showAddModal = false}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Adicionar Usuário</h3>
        <Button onclick={() => showAddModal = false} variant="ghost" size="sm">✕</Button>
      </div>
      <div class="modal-body">
        {#if addUserError}
          <Alert variant="error" dismissible onDismiss={() => addUserError = ''}>
            <AlertTriangle size={16} />
            {addUserError}
          </Alert>
        {/if}
        <Input
          id="add-user-email"
          label="E-mail do Usuário"
          bind:value={addUserEmail}
          placeholder="usuario@exemplo.com"
          disabled={addingUser}
        />
        <p class="help-text">
          O usuário deve já estar cadastrado no sistema. Será atribuído o cargo padrão de Gerente.
        </p>
      </div>
      <div class="modal-footer">
        <Button onclick={() => showAddModal = false} variant="secondary" disabled={addingUser}>
          Cancelar
        </Button>
        <Button onclick={addUser} variant="primary" disabled={addingUser} loading={addingUser}>
          Adicionar
        </Button>
      </div>
    </div>
  </div>
{/if}

<!-- Change Role Modal -->
{#if showRoleModal && selectedUser}
  <div class="modal-overlay" onclick={() => showRoleModal = false}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Alterar Cargo</h3>
        <Button onclick={() => showRoleModal = false} variant="ghost" size="sm">✕</Button>
      </div>
      <div class="modal-body">
        {#if roleError}
          <Alert variant="error" dismissible onDismiss={() => roleError = ''}>
            <AlertTriangle size={16} />
            {roleError}
          </Alert>
        {/if}
        <p class="user-info">
          Usuário: <strong>{selectedUser.name}</strong> ({selectedUser.email})
        </p>
        <div class="form-group">
          <label for="role-select" class="form-label">Novo Cargo</label>
          <select
            id="role-select"
            bind:value={newRole}
            disabled={changingRole}
            class="form-select"
          >
            {#each roles as role}
              <option value={role.value}>{role.label}</option>
            {/each}
          </select>
        </div>
      </div>
      <div class="modal-footer">
        <Button onclick={() => showRoleModal = false} variant="secondary" disabled={changingRole}>
          Cancelar
        </Button>
        <Button onclick={changeRole} variant="primary" disabled={changingRole} loading={changingRole}>
          Alterar Cargo
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

  /* Users Container */
  .users-container {
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

  /* Users Table */
  .users-table {
    padding: 1.5rem;
    overflow-x: auto;
  }

  .table {
    width: 100%;
    border-collapse: collapse;
  }

  .table thead {
    border-bottom: 2px solid #e5e7eb;
  }

  .table th {
    text-align: left;
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
  }

  .table td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .table tbody tr:last-child td {
    border-bottom: none;
  }

  /* Badge */
  .badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  /* Status Badge */
  .status-badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .status-badge.active {
    background: #dcfce7;
    color: #166534;
  }

  .status-badge.inactive {
    background: #fee2e2;
    color: #991b1b;
  }

  /* Actions */
  .actions {
    display: flex;
    gap: 0.5rem;
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

  .user-info {
    font-size: 0.875rem;
    color: #64748b;
    margin-bottom: 1rem;
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
    .table {
      font-size: 0.875rem;
    }

    .table th,
    .table td {
      padding: 0.5rem;
    }

    .actions {
      flex-direction: column;
    }
  }
</style>
