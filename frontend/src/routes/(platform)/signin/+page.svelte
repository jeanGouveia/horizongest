<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui';
	import { Input } from '$lib/components/ui';
	import { Card } from '$lib/components/ui';
	import { showSuccess, showError } from '$lib/stores/toast';
	import { platformName } from '$lib/stores/brandStore';

	let email = '';
	let password = '';
	let loading = false;

	async function handleLogin() {
		if (!email || !password) {
			showError('Erro', 'Por favor, preencha todos os campos');
			return;
		}

		loading = true;
		try {
			const response = await fetch('/api/platform/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, password })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message || 'Falha no login');
			}

			// Store token in cookie
			document.cookie = `platform_auth_token=${data.token}; path=/; max-age=86400; secure; samesite=strict`;

			showSuccess('Sucesso', 'Login realizado com sucesso');
			goto('/platform/admin');
		} catch (error) {
			showError('Erro', error instanceof Error ? error.message : 'Falha no login');
		} finally {
			loading = false;
		}
	}
</script>

<div class="login-container">
	<Card class="login-card">
		<div class="login-header">
			<h1>{$platformName} Platform</h1>
			<p>Área Administrativa</p>
		</div>

		<form on:submit|preventDefault={handleLogin} class="login-form">
			<div class="form-group">
				<label for="email">E-mail</label>
				<Input
					id="email"
					type="email"
					placeholder="admin@example.com"
					bind:value={email}
					disabled={loading}
					required
				/>
			</div>

			<div class="form-group">
				<label for="password">Senha</label>
				<Input
					id="password"
					type="password"
					placeholder="••••••••"
					bind:value={password}
					disabled={loading}
					required
				/>
			</div>

			<Button type="submit" variant="primary" size="lg" loading={loading} class="login-button">
				Entrar
			</Button>
		</form>

		<div class="login-footer">
			<p>Acesso restrito a administradores da plataforma</p>
		</div>
	</Card>
</div>

<style>
	.login-container {
		width: 100%;
		max-width: 420px;
	}

	.login-card {
		background: white;
		border-radius: 16px;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
		padding: 2.5rem;
	}

	.login-header {
		text-align: center;
		margin-bottom: 2rem;
	}

	.login-header h1 {
		font-size: 1.75rem;
		font-weight: 700;
		color: #1e293b;
		margin: 0 0 0.5rem 0;
	}

	.login-header p {
		font-size: 0.875rem;
		color: #64748b;
		margin: 0;
	}

	.login-form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #475569;
	}

	.login-button {
		margin-top: 0.5rem;
		width: 100%;
	}

	.login-footer {
		text-align: center;
		margin-top: 2rem;
		padding-top: 1.5rem;
		border-top: 1px solid #e2e8f0;
	}

	.login-footer p {
		font-size: 0.75rem;
		color: #94a3b8;
		margin: 0;
	}
</style>
