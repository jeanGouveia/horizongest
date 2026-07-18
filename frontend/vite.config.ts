import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    port: 3000,
    proxy: {
      // Todo /api/* é repassado ao Go — mesmo domínio, sem CORS
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      // Arquivos estáticos de upload (imagens)
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
});
