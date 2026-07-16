<script lang="ts">
	interface Props {
		version?: string;
		year?: number;
	}

	let {
		version = 'MVP',
		year = new Date().getFullYear(),
	}: Props = $props();
</script>

<footer class="footer">
	<div class="footer-content">
		<div class="footer-brand">
			<span class="footer-logo">🍽️ PratoOnline</span>
			<span class="footer-version">v{version}</span>
		</div>
		<div class="footer-copyright">
			© {year} PratoOnline. Todos os direitos reservados.
		</div>
	</div>
</footer>

<style>
	.footer {
		background-color: white;
		border-top: 1px solid #e5e5e5;
		padding: 16px 24px;
		margin-top: auto;
	}

	.footer-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
	}

	.footer-brand {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 0.875rem;
		color: #171717;
	}

	.footer-logo {
		font-weight: 600;
	}

	.footer-version {
		color: #525252;
		font-size: 0.75rem;
	}

	.footer-copyright {
		font-size: 0.75rem;
		color: #525252;
	}

	@media (max-width: 768px) {
		.footer-content {
			flex-direction: column;
			text-align: center;
		}
	}
</style>
