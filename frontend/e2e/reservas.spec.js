import { expect, test } from '@playwright/test'

// e2e: un navegador de verdad contra el front de QA (E2E_BASE_URL), sin mocks.
// El usuario de prueba (admin) lo crea el job antes de correr; su contraseña
// viene del secret E2E_PASSWORD.
const USUARIO = 'e2e-admin'
const PASSWORD = process.env.E2E_PASSWORD

async function iniciarSesion (page) {
  // se entra por "/" como un usuario: la app redirige al login. (Abrir /login directo
  // no sirve: nginx manda /login a users-api, que es el POST del login.)
  await page.goto('/')
  await page.getByLabel('Usuario o Email').fill(USUARIO)
  await page.getByLabel('Contraseña').fill(PASSWORD)
  await page.getByRole('button', { name: 'Iniciar Sesión' }).click()
  await expect(page.getByRole('heading', { name: 'Bienvenido al Hotel Grand Prix' })).toBeVisible()
}

// Ayudas para preparar y limpiar datos por la api (a través del front, como el navegador)
async function tokenDeAdmin (request) {
  const r = await request.post('/login', { data: { username_or_email: USUARIO, password: PASSWORD } })
  expect(r.status()).toBe(200)
  return (await r.json()).login.token
}

async function crearHabitacion (request, token, numero) {
  const r = await request.post('/api/v1/admin/rooms', {
    headers: { Authorization: `Bearer ${token}` },
    data: { number: numero, type: 'single', price: 100, capacity: 1, floor: 1, description: 'creada por la e2e' },
  })
  expect(r.status()).toBe(201)
  return r.json()
}

async function eliminarHabitacionYComprobar (request, token, id) {
  const r = await request.delete(`/api/v1/admin/rooms/${id}`, { headers: { Authorization: `Bearer ${token}` } })
  expect(r.ok()).toBe(true)
  expect((await request.get(`/api/v1/rooms/${id}`)).status()).toBe(404)
}

test('un admin crea una habitación, la ve en el listado, y al eliminarla ya no está', async ({ page }) => {
  const numero = `E2E${Date.now()}`   // único: no choca con otra corrida
  await iniciarSesion(page)

  await page.goto('/admin/rooms/new')
  await page.getByLabel('Número de Habitación').fill(numero)
  await page.getByLabel('Precio por Noche').fill('120')
  await page.getByRole('button', { name: 'Create Room' }).click()
  await expect(page.getByRole('cell', { name: numero, exact: true })).toBeVisible()

  await page.getByRole('button', { name: `Eliminar la habitación ${numero}` }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Eliminar' }).click()
  await expect(page.getByRole('cell', { name: numero, exact: true })).toHaveCount(0)
})

test('crear una habitación con un número que ya existe muestra el error y no crea otra', async ({ page, request }) => {
  const numero = `E2E${Date.now()}`
  const token = await tokenDeAdmin(request)
  const existente = await crearHabitacion(request, token, numero)

  try {
    await iniciarSesion(page)
    await page.goto('/admin/rooms/new')
    await page.getByLabel('Número de Habitación').fill(numero)
    await page.getByLabel('Precio por Noche').fill('120')
    await page.getByRole('button', { name: 'Create Room' }).click()
    await expect(page.getByRole('alert')).toBeVisible()          // el usuario ve el error…

    await page.goto('/admin/rooms')
    await expect(page.getByRole('cell', { name: numero, exact: true })).toHaveCount(1)   // …y no se creó otra
  } finally {
    await eliminarHabitacionYComprobar(request, token, existente.id)
  }
})

test('un usuario reserva una habitación, la ve en Mis reservas y la cancela', async ({ page, request }) => {
  const token = await tokenDeAdmin(request)
  const habitacion = await crearHabitacion(request, token, `E2E${Date.now()}`)

  try {
    await iniciarSesion(page)
    await page.goto(`/rooms/${habitacion.id}/reserve`)
    await page.getByLabel('Fecha de Entrada').fill('2030-03-10')
    await page.getByLabel('Fecha de Salida').fill('2030-03-12')
    await page.getByRole('button', { name: 'Confirmar Reserva' }).click()

    // la app redirige a Mis reservas, y ahí está la reserva de ESTA habitación
    await expect(page).toHaveURL(/\/my-reservations/)
    const cancelar = page.getByRole('button', { name: `Cancelar la reserva de la habitación ${habitacion.id}` })
    await expect(cancelar).toBeVisible()

    page.once('dialog', (dialogo) => dialogo.accept())   // el "¿Estás seguro?" de la app
    await cancelar.click()
    await expect(cancelar).toHaveCount(0)                  // ya no está activa: no se puede volver a cancelar
  } finally {
    await eliminarHabitacionYComprobar(request, token, habitacion.id)
  }
})
