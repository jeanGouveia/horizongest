import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const backendUrl = process.env.VITE_BACKEND_URL || 'http://localhost:8080';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    port: 3000,
    proxy: {
      // Todo /api/* é repassado ao Go — mesmo domínio, sem CORS
      '/api': {
        target: backendUrl,
        changeOrigin: true
      },
      // Arquivos estáticos de upload (imagens)
      '/uploads': {
        target: backendUrl,
        changeOrigin: true
      }
    }
  }
});
