<script lang="ts">
	interface Props {
		maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
		class?: string;
		padding?: 'sm' | 'md' | 'lg' | 'none';
	}

	let {
		maxWidth = 'xl',
		padding = 'md',
		class: className = '',
	}: Props = $props();
</script>

<div
	class:page-container={true}
	class:page-container-sm={maxWidth === 'sm'}
	class:page-container-md={maxWidth === 'md'}
	class:page-container-lg={maxWidth === 'lg'}
	class:page-container-xl={maxWidth === 'xl'}
	class:page-container-full={maxWidth === 'full'}
	class:page-container-padding-sm={padding === 'sm'}
	class:page-container-padding-md={padding === 'md'}
	class:page-container-padding-lg={padding === 'lg'}
	class:page-container-padding-none={padding === 'none'}
	class={className}
>
	<slot />
</div>

<style>
	.page-container {
		width: 100%;
		margin: 0 auto;
		padding: 0 1.5rem;
	}

	.page-container-padding-sm {
		padding: 0 1rem;
	}

	.page-container-padding-md {
		padding: 0 1.5rem;
	}

	.page-container-padding-lg {
		padding: 0 2rem;
	}

	.page-container-padding-none {
		padding: 0;
	}

	.page-container-sm {
		max-width: 640px;
	}

	.page-container-md {
		max-width: 768px;
	}

	.page-container-lg {
		max-width: 1024px;
	}

	.page-container-xl {
		max-width: 1280px;
	}

	.page-container-full {
		max-width: 100%;
	}

	@media (max-width: 768px) {
		.page-container {
			padding: 0 1rem;
		}

		.page-container-padding-sm {
			padding: 0 0.75rem;
		}

		.page-container-padding-md {
			padding: 0 1rem;
		}

		.page-container-padding-lg {
			padding: 0 1.25rem;
		}
	}
</style>
