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
        changeOrigin: true,
        secure: false,
        configure: (proxy, options) => {
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // Repassar cookies
            if (req.headers.cookie) {
              proxyReq.setHeader('Cookie', req.headers.cookie);
            }
            console.log('Proxy request:', req.method, req.url, 'Cookie:', req.headers.cookie ? 'present' : 'missing');
          });
          proxy.on('proxyRes', (proxyRes, req, res) => {
            console.log('Proxy response:', proxyRes.statusCode, req.url);
          });
          proxy.on('error', (err, req, res) => {
            console.log('proxy error', err);
          });
        }
      },
      // Arquivos estáticos de upload (imagens)
      '/uploads': {
        target: backendUrl,
        changeOrigin: true
      }
    }
  }
});
