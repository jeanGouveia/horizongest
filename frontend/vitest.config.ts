import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { resolve } from 'path';

export default defineConfig({
	plugins: [svelte()],
	resolve: {
		alias: {
			$lib: resolve(__dirname, './src/lib'),
			'$lib/*': resolve(__dirname, './src/lib/*')
		}
	},
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./tests/unit/session/setup.ts'],
		include: ['tests/**/*.{test,spec}.{js,ts}'],
		coverage: {
			provider: 'v8',
			reporter: ['text', 'json', 'html'],
			exclude: ['node_modules/', 'tests/']
		}
	}
});
