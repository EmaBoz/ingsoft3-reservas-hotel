import { describe, expect, it, vi } from 'vitest';
import {
  calcularTotal,
  validarFechas,
  misReservas,
  ERROR_FECHAS_OBLIGATORIAS,
  ERROR_SALIDA_NO_POSTERIOR,
} from './reservas';

describe('calcularTotal', () => {
  it.each([
    ['una noche', '2026-11-10', '2026-11-11', 100, 100],
    ['tres noches', '2026-11-10', '2026-11-13', 100, 300],
    ['salida el mismo día', '2026-11-10', '2026-11-10', 100, 0],
    ['salida antes que la entrada', '2026-11-10', '2026-11-08', 100, 0],
  ])('con %s da el total correcto', (_caso, entrada, salida, precio, esperado) => {
    expect(calcularTotal(entrada, salida, precio)).toBe(esperado);
  });

  it('da 0 y no NaN si el usuario borró la fecha de entrada', () => {
    expect(calcularTotal('', '2026-11-12', 100)).toBe(0);
  });
});

describe('validarFechas', () => {
  it('acepta una salida posterior a la entrada', () => {
    expect(validarFechas('2026-11-10', '2026-11-12')).toBeNull();
  });

  it('rechaza una salida igual a la entrada con el mensaje del formulario', () => {
    expect(validarFechas('2026-11-10', '2026-11-10')).toBe(ERROR_SALIDA_NO_POSTERIOR);
  });

  it.each([
    ['sin entrada', '', '2026-11-12'],
    ['sin salida', '2026-11-10', ''],
  ])('rechaza las fechas %s', (_caso, entrada, salida) => {
    expect(validarFechas(entrada, salida)).toBe(ERROR_FECHAS_OBLIGATORIAS);
  });
});

describe('misReservas', () => {
  it('devuelve las reservas que contesta la API', async () => {
    const traer = vi.fn().mockResolvedValue({
      reservations: [{ id: 'a1', status: 'active' }, { id: 'b2', status: 'canceled' }],
    });

    const reservas = await misReservas('7', traer);

    expect(reservas.map((r) => r.id)).toEqual(['a1', 'b2']);
  });

  it('le pide a la API la ruta de las reservas del usuario', async () => {
    const traer = vi.fn().mockResolvedValue({ reservations: [] });

    await misReservas('7', traer);

    expect(traer).toHaveBeenCalledWith('/api/reservations/users/7/myreservations');
  });

  it('devuelve una lista vacía si la API contesta null', async () => {
    const traer = vi.fn().mockResolvedValue({ reservations: null });

    expect(await misReservas('7', traer)).toEqual([]);
  });
});
