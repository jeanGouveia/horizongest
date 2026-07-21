<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import { platformName } from '$lib/stores/brandStore';

	let email = '';
	let error = '';
	let success = false;
	let loading = false;

	async function handleSubmit() {
		error = '';
		loading = true;

		try {
			const res = await api.auth.requestPasswordReset({ email });
			if (res.error) {
				error = res.error;
			} else {
				success = true;
			}
		} catch (e: any) {
			error = e.message || 'Erro ao solicitar recuperação de senha.';
		} finally {
			loading = false;
		}
	}

	function goToLogin() {
		goto('/login');
	}
</script>

<svelte:head>
	<title>Recuperar Senha - {$platformName}</title>
</svelte:head>

<div class="auth-container">
	<div class="auth-card">
		<h1>Recuperar Senha</h1>
		<p>Digite seu e-mail para receber instruções de recuperação de senha.</p>

		{#if success}
			<Alert variant="success">
				✓ Se o e-mail estiver cadastrado, você receberá instruções para recuperar sua senha.
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
					<label for="email">E-mail</label>
					<Input
						id="email"
						type="email"
						bind:value={email}
						placeholder="seu@email.com"
						required
						disabled={loading}
					/>
				</div>

				<Button type="submit" variant="primary" fullWidth loading={loading}>
					Enviar Instruções
				</Button>
			</form>

			<div class="auth-footer">
				<button type="button" class="link-button" on:click={goToLogin}>
					Voltar para Login
				</button>
			</div>
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
