import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
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
