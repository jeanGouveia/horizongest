<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api } from '$lib/api/client';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import { platformName } from '$lib/stores/brandStore';

	let newPassword = '';
	let confirmPassword = '';
	let error = '';
	let success = false;
	let loading = false;

	// Get token from URL query parameter
	let token = '';
	$: {
		const urlToken = $page.url.searchParams.get('token');
		if (urlToken && !token) {
			token = urlToken;
		}
	}

	async function handleSubmit() {
		error = '';

		if (newPassword.length < 6) {
			error = 'A senha deve ter no mínimo 6 caracteres.';
			return;
		}

		if (newPassword !== confirmPassword) {
			error = 'As senhas não coincidem.';
			return;
		}

		loading = true;

		try {
			const res = await api.auth.resetPassword({ token, new_password: newPassword });
			if (res.error) {
				error = res.error;
			} else {
				success = true;
			}
		} catch (e: any) {
			error = e.message || 'Erro ao redefinir senha.';
		} finally {
			loading = false;
		}
	}

	function goToLogin() {
		goto('/login');
	}
</script>

<svelte:head>
	<title>Redefinir Senha - {$platformName}</title>
</svelte:head>

<div class="auth-container">
	<div class="auth-card">
		<h1>Redefinir Senha</h1>
		<p>Digite sua nova senha abaixo.</p>

		{#if success}
			<Alert variant="success">
				✓ Senha redefinida com sucesso.
			</Alert>
			<Button onclick={goToLogin} variant="primary" fullWidth>
				Fazer Login
			</Button>
		{:else}
			{#if !token}
				<Alert variant="error">
					Token inválido. Solicite uma nova recuperação de senha.
				</Alert>
				<Button onclick={goToLogin} variant="primary" fullWidth>
					Voltar para Login
				</Button>
			{:else}
				{#if error}
					<Alert variant="error" dismissible onDismiss={() => error = ''}>
						{error}
					</Alert>
				{/if}

				<form on:submit|preventDefault={handleSubmit}>
					<div class="form-group">
						<label for="newPassword">Nova Senha</label>
						<Input
							id="newPassword"
							type="password"
							bind:value={newPassword}
							placeholder="••••••"
							required
							disabled={loading}
						/>
					</div>

					<div class="form-group">
						<label for="confirmPassword">Confirmar Senha</label>
						<Input
							id="confirmPassword"
							type="password"
							bind:value={confirmPassword}
							placeholder="••••••"
							required
							disabled={loading}
						/>
					</div>

					<Button type="submit" variant="primary" fullWidth loading={loading}>
						Redefinir Senha
					</Button>
				</form>

				<div class="auth-footer">
					<button type="button" class="link-button" on:click={goToLogin}>
						Voltar para Login
					</button>
				</div>
			{/if}
		{/if}
	</div>
</div>

<style>
	.auth-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		padding: 1rem;
	}

	.auth-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		width: 100%;
		max-width: 400px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
	}

	.auth-card h1 {
		margin: 0 0 0.5rem 0;
		font-size: 1.5rem;
		color: #1a202c;
	}

	.auth-card > p {
		margin: 0 0 1.5rem 0;
		color: #718096;
		font-size: 0.875rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 500;
		color: #4a5568;
		font-size: 0.875rem;
	}

	.auth-footer {
		margin-top: 1.5rem;
		text-align: center;
	}

	.link-button {
		background: none;
		border: none;
		color: #667eea;
		cursor: pointer;
		font-size: 0.875rem;
		padding: 0;
		text-decoration: underline;
	}

	.link-button:hover {
		color: #764ba2;
	}
</style>
