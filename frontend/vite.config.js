import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:7070',
      '/media': 'http://localhost:7070',
      '/healthz': {
        target: 'http://localhost:7070',
        rewrite: () => '/api/healthz'
      }
    }
  },
  build: {
    sourcemap: false
  }
});
