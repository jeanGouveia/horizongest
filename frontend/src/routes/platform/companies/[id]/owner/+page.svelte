<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card } from '$lib/components/ui';
	import { Button } from '$lib/components/ui';
	import { Input } from '$lib/components/ui';
	import { Badge } from '$lib/components/ui';
	import { showSuccess, showError } from '$lib/stores/toast';
	import { CookieKeys } from '$lib/constants/storage-keys';

	interface Owner {
		ID: number;
		Name: string;
		Email: string;
		Role: string;
		Active: boolean;
		CreatedAt: string;
	}

	let companyID: number;
	let owner: Owner | null = null;
	let loading = true;
	let showResetPasswordModal = false;
	let newPassword = '';
	let newPasswordConfirm = '';

	onMount(async () => {
		const pathParts = window.location.pathname.split('/');
		companyID = parseInt(pathParts[3]);
		if (isNaN(companyID)) {
			showError('Erro', 'ID da empresa inválido');
			goto('/platform/companies');
			return;
		}
		await loadOwner();
	});

	async function loadOwner() {
		loading = true;
		try {
			const token = document.cookie.split('; ').find(row => row.startsWith(`${CookieKeys.PLATFORM_TOKEN}=`))?.split('=')[1];
			const headers: Record<string, string> = {
				'Content-Type': 'application/json'
			};
			if (token) {
				headers['Authorization'] = `Bearer ${token}`;
			}

			const response = await fetch(`http://localhost:8080/api/platform/companies/${companyID}/owner`, { headers });
			if (response.ok) {
				owner = await response.json();
			} else {
				showError('Erro', 'Erro ao carregar owner');
				goto(`/platform/companies/${companyID}`);
			}
		} catch (error) {
			showError('Erro', 'Erro ao carregar owner');
			goto(`/platform/companies/${companyID}`);
		} finally {
			loading = false;
		}
	}

	async function resetPassword() {
		if (!newPassword || newPassword.length < 8) {
			showError('Erro', 'Senha deve ter no mínimo 8 caracteres');
			return;
		}
		if (newPassword !== newPasswordConfirm) {
			showError('Erro', 'Senhas não conferem');
			return;
		}

		try {
			const token = document.cookie.split('; ').find(row => row.startsWith(`${CookieKeys.PLATFORM_TOKEN}=`))?.split('=')[1];
			const headers: Record<string, string> = {
				'Content-Type': 'application/json'
			};
			if (token) {
				headers['Authorization'] = `Bearer ${token}`;
			}

			const response = await fetch(`http://localhost:8080/api/platform/companies/${companyID}/owner/reset-password`, {
				method: 'POST',
				headers,
				body: JSON.stringify({ password: newPassword })
			});

			if (response.ok) {
				showSuccess('Sucesso', 'Senha redefinida com sucesso');
				showResetPasswordModal = false;
				newPassword = '';
				newPasswordConfirm = '';
			} else {
				showError('Erro', 'Erro ao redefinir senha');
			}
		} catch (error) {
			if (error instanceof TypeError && error.message.includes('fetch')) {
				showError('Erro de conexão', 'Não foi possível conectar ao servidor. Verifique se o backend está rodando em http://localhost:8080');
			} else {
				showError('Erro', 'Erro ao redefinir senha');
			}
		}
	}

	async function toggleUserBlock() {
		if (!owner) return;

		const action = owner.Active ? 'block' : 'unblock';
		try {
			const token = document.cookie.split('; ').find(row => row.startsWith(`${CookieKeys.PLATFORM_TOKEN}=`))?.split('=')[1];
			const headers: Record<string, string> = {};
			if (token) {
				headers['Authorization'] = `Bearer ${token}`;
			}

			const response = await fetch(`http://localhost:8080/api/platform/users/${owner.ID}/${action}`, {
				method: 'POST',
				headers
			});

			if (response.ok) {
				showSuccess('Sucesso', owner.Active ? 'Usuário bloqueado' : 'Usuário desbloqueado');
				await loadOwner();
			} else {
				showError('Erro', 'Erro ao alterar status do usuário');
			}
		} catch (error) {
			if (error instanceof TypeError && error.message.includes('fetch')) {
				showError('Erro de conexão', 'Não foi possível conectar ao servidor. Verifique se o backend está rodando em http://localhost:8080');
			} else {
				showError('Erro', 'Erro ao alterar status do usuário');
			}
		}
	}
</script>

<div class="owner-container">
	<header class="owner-header">
		<div class="header-content">
			<Button variant="secondary" onclick={() => goto(`/platform/companies/${companyID}`)}>← Voltar</Button>
			<div>
				<h1>Gerenciar Owner</h1>
				<p>Configurações do dono da empresa</p>
			</div>
		</div>
	</header>

	{#if loading}
		<div class="loading-state">
			<p>Carregando...</p>
		</div>
	{:else if owner}
		<div class="owner-content">
			<Card class="owner-card">
				<h2>Informações do Owner</h2>
				<div class="detail-row">
					<strong>Nome:</strong>
					<span>{owner.Name}</span>
				</div>
				<div class="detail-row">
					<strong>E-mail:</strong>
					<span>{owner.Email}</span>
				</div>
				<div class="detail-row">
					<strong>Função:</strong>
					<Badge variant={owner.Role === 'owner' ? 'primary' : 'default'}>{owner.Role}</Badge>
				</div>
				<div class="detail-row">
					<strong>Status:</strong>
					<Badge variant={owner.Active ? 'success' : 'error'}>{owner.Active ? 'Ativo' : 'Bloqueado'}</Badge>
				</div>
				<div class="detail-row">
					<strong>Criado em:</strong>
					<span>{new Date(owner.CreatedAt).toLocaleString('pt-BR')}</span>
				</div>

				<div class="owner-actions">
					<Button variant="secondary" onclick={() => showResetPasswordModal = true}>
						Redefinir Senha
					</Button>
					<Button
						variant={owner.Active ? 'secondary' : 'primary'}
						onclick={toggleUserBlock}
					>
						{owner.Active ? 'Bloquear Usuário' : 'Desbloquear Usuário'}
					</Button>
				</div>
			</Card>
		</div>
	{/if}

	{#if showResetPasswordModal}
		<div class="modal-overlay" onclick={() => showResetPasswordModal = false}>
			<Card class="modal-card">
				<h2>Redefinir Senha</h2>
				<div class="form-group">
					<label>Nova Senha *</label>
					<Input type="password" bind:value={newPassword} placeholder="Mínimo 8 caracteres" />
				</div>
				<div class="form-group">
					<label>Confirmar Senha *</label>
					<Input type="password" bind:value={newPasswordConfirm} placeholder="Confirme a senha" />
				</div>
				<div class="modal-actions">
					<Button variant="secondary" onclick={() => showResetPasswordModal = false}>Cancelar</Button>
					<Button variant="primary" onclick={resetPassword}>Redefinir</Button>
				</div>
			</Card>
		</div>
	{/if}
</div>

<style>
	.owner-container {
		padding: 2rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	.owner-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	.header-content {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.owner-header h1 {
		font-size: 1.5rem;
		font-weight: 600;
		margin: 0;
	}

	.owner-header p {
		margin: 0;
		color: #64748b;
	}

	.loading-state {
		text-align: center;
		padding: 2rem;
		color: #64748b;
	}

	.owner-content {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.owner-card {
		background: rgba(255, 255, 255, 0.1);
		padding: 1.5rem;
	}

	.owner-card h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin-bottom: 1.5rem;
	}

	.detail-row {
		display: flex;
		justify-content: space-between;
		padding: 0.75rem 0;
		border-bottom: 1px solid rgba(0, 0, 0, 0.1);
	}

	.detail-row:last-child {
		border-bottom: none;
	}

	.detail-row strong {
		font-weight: 500;
		color: #0f172a;
	}

	.owner-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1.5rem;
	}

	.owner-actions button {
		flex: 1;
	}

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

	.modal-card {
		background: #1e293b;
		padding: 2rem;
		min-width: 400px;
		max-width: 500px;
	}

	.modal-card h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin-bottom: 1.5rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 500;
		color: #e2e8f0;
	}

	.modal-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1.5rem;
	}

	.modal-actions button {
		flex: 1;
	}
</style>
