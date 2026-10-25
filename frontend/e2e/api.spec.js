import { expect, test } from '@playwright/test'

// Integración: le hablamos directo a reservations-api de QA, con su Mongo de verdad
// atrás. Sin navegador y sin dobles. La api tiene su propia dirección (no es la del front).
const API = process.env.API_BASE_URL || 'http://localhost:8082'

// Cada prueba se fabrica un usuario y una habitación que no se repiten (con la hora
// adentro), así no choca con otra corrida ni con datos reales. reservations-api no
// valida que la habitación exista en rooms-api. user_id < 2^32 porque la ruta de
// "mis reservas" lo parsea como uint32.
function datosUnicos () {
  const n = Date.now() % 1_000_000_000
  return { user_id: n, room_id: n + 1 }
}

async function reservar (request, datos) {
  return request.post(`${API}/api/reservations`, { data: datos })
}

async function cancelarYComprobar (request, id) {
  const baja = await request.delete(`${API}/api/reservations/${id}`, { data: { reason: 'limpieza de la prueba' } })
  expect(baja.status()).toBe(204)
  const despues = await request.get(`${API}/api/reservations/${id}`)
  expect((await despues.json()).reservation.status).toBe('canceled')
}

test('el alta guarda la reserva en la base, y cancelarla la deja cancelada', async ({ request }) => {
  const { user_id, room_id } = datosUnicos()

  const alta = await reservar(request, { user_id, room_id, start_date: '2030-01-10', end_date: '2030-01-12' })
  expect(alta.status()).toBe(201)
  const id = (await alta.json()).reservation.id

  // la devuelve la BASE, no un doble
  const mias = await request.get(`${API}/api/reservations/users/${user_id}/myreservations`)
  expect(mias.status()).toBe(200)
  const reserva = (await mias.json()).reservations.find((r) => r.id === id)
  expect(reserva.status).toBe('active')

  // cancelar es el "borrado" de la app (soft delete): deja de estar activa
  await cancelarYComprobar(request, id)
})

test('una reserva sin fecha de entrada la rechaza la api, y no se crea nada', async ({ request }) => {
  const { user_id, room_id } = datosUnicos()

  const alta = await reservar(request, { user_id, room_id, end_date: '2030-01-12' })
  expect(alta.status()).toBe(400)

  const mias = await request.get(`${API}/api/reservations/users/${user_id}/myreservations`)
  expect(mias.status()).toBe(200)
  expect((await mias.json()).reservations ?? []).toHaveLength(0)
})

test('no deja reservar una habitación que ya está reservada en esas fechas', async ({ request }) => {
  const { user_id, room_id } = datosUnicos()

  const primera = await reservar(request, { user_id, room_id, start_date: '2030-02-10', end_date: '2030-02-15' })
  expect(primera.status()).toBe(201)
  const id = (await primera.json()).reservation.id

  // otro usuario, la misma habitación, fechas que se pisan: la consulta de
  // superposición contra Mongo tiene que encontrar la primera
  const pisada = await reservar(request, { user_id: user_id + 1, room_id, start_date: '2030-02-12', end_date: '2030-02-14' })
  expect(pisada.status()).toBe(409)

  await cancelarYComprobar(request, id)
})
