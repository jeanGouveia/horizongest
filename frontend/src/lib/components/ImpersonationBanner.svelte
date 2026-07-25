<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from './ui';
	import { tenantSessionManager } from '$lib/managers/tenantSessionManager';

	interface ImpersonationInfo {
		isImpersonating: boolean;
		companyName: string;
		ownerEmail: string;
		companyId: number;
	}

	let impersonationInfo: ImpersonationInfo | null = $state(null);
	let visible = $state(false);

	onMount(() => {
		const stored = localStorage.getItem('impersonation');
		if (stored) {
			try {
				impersonationInfo = JSON.parse(stored);
				if (impersonationInfo?.isImpersonating) {
					visible = true;
				}
			} catch (e) {
				console.error('Error parsing impersonation info:', e);
				localStorage.removeItem('impersonation');
			}
		}
	});

	async function endImpersonation() {
		try {
			const result = await tenantSessionManager.leaveCompany();
			
			if (result.success) {
				impersonationInfo = null;
				visible = false;
			} else {
				console.error('Error ending impersonation:', result.error);
			}
		} catch (error) {
			console.error('Error ending impersonation:', error);
		}
	}
</script>

{#if visible && impersonationInfo}
	<div class="impersonation-banner">
		<div class="banner-content">
			<div class="banner-icon">⚠️</div>
			<div class="banner-text">
				<strong>Você está operando como {impersonationInfo.companyName}</strong>
				<p>através da Plataforma. Todas as ações estão sendo auditadas.</p>
			</div>
			<Button variant="primary" onclick={endImpersonation}>Voltar para Plataforma</Button>
		</div>
	</div>
{/if}

<style>
	.impersonation-banner {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
		border-bottom: 3px solid #b45309;
		z-index: 9999;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
	}

	.banner-content {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 2rem;
		max-width: 1400px;
		margin: 0 auto;
	}

	.banner-icon {
		font-size: 1.5rem;
		flex-shrink: 0;
	}

	.banner-text {
		flex: 1;
		color: white;
	}

	.banner-text strong {
		display: block;
		font-size: 1rem;
		margin-bottom: 0.25rem;
	}

	.banner-text p {
		margin: 0;
		font-size: 0.875rem;
		opacity: 0.9;
	}

	.banner-content :global(button) {
		flex-shrink: 0;
		background: white;
		color: #d97706;
		font-weight: 600;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.banner-content :global(button:hover) {
		background: #fef3c7;
		transform: translateY(-1px);
	}
</style>
