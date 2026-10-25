import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',                // los dos archivos: api.spec.js (integración) y reservas.spec.js (e2e)
  timeout: 60_000,                 // tope de cada test
  expect: { timeout: 15_000 },     // tope de cada aserción (el default es 5 s)
  use: {
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',   // el FRONT; la api usa API_BASE_URL
    actionTimeout: 15_000,           // un clic que no encuentra su elemento falla acá, y la limpieza del test alcanza a correr
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  retries: 1,                      // un reintento: si pasa recién ahí, el reporte lo marca como flaky
  workers: 1,                      // de a uno: los tests comparten la base de QA
  reporter: [['html', { open: 'never' }], ['list']],
})
