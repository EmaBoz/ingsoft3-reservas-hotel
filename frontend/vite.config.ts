/// <reference types="vitest/config" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import { configDefaults } from 'vitest/config'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  test: {
    // vitest no toca los specs de Playwright (e2e/), que corren contra QA
    exclude: [...configDefaults.exclude, 'e2e/**'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'lcov', 'json-summary'],
      // Entra en la cuenta la lógica de src/lib (ver decisiones.md, TP5)
      include: ['src/lib/**'],
      exclude: ['src/lib/axios.ts'],
      thresholds: { lines: 80, branches: 80 },
    },
  },
  server: {
    port: 3000,
    proxy: {
      // Users API - login, register, users
      '/login': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/users': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // Rooms API (including admin routes)
      '/api/v1': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      // Reservations API
      '/api/reservations': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/search-api': {
        target: 'http://localhost:8083',   // o la URL interna del container search-api
        changeOrigin: true,
      },
    },
  },
})
