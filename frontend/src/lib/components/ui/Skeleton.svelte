<script lang="ts">
	interface Props {
		class?: string;
		variant?: 'text' | 'circular' | 'rectangular';
		width?: string;
		height?: string;
	}

	let {
		class: className = '',
		variant = 'text',
		width = '100%',
		height = '1rem',
	}: Props = $props();
</script>

<div
	class="skeleton skeleton-{variant} {className}"
	style="width: {width}; height: {height};"
	aria-hidden="true"
></div>

<style>
	.skeleton {
		background-color: #f1f5f9;
		border-radius: 4px;
		animation: shimmer 1.5s infinite linear;
		background-image: linear-gradient(
			90deg,
			#f1f5f9 0px,
			#e2e8f0 40px,
			#f1f5f9 80px
		);
		background-size: 200% 100%;
	}

	.skeleton-text {
		border-radius: 4px;
	}

	.skeleton-circular {
		border-radius: 50%;
	}

	.skeleton-rectangular {
		border-radius: 8px;
	}

	@keyframes shimmer {
		0% {
			background-position: 200% 0;
		}
		100% {
			background-position: -200% 0;
		}
	}
</style>
