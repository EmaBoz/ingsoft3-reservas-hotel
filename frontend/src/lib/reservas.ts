import { differenceInDays } from 'date-fns';
import type { Reservation } from '@/types';

export const ERROR_FECHAS_OBLIGATORIAS = 'Las fechas de entrada y salida son obligatorias';
export const ERROR_SALIDA_NO_POSTERIOR = 'La fecha de salida debe ser posterior a la fecha de entrada';

// Noches entre la fecha de entrada y la de salida (las dos en formato yyyy-MM-dd).
export function calcularNoches(entrada: string, salida: string): number {
  if (!entrada || !salida) return 0;
  return differenceInDays(new Date(salida), new Date(entrada));
}

// Total de la estadía. Si las fechas no forman una estadía válida, da 0.
export function calcularTotal(entrada: string, salida: string, precioPorNoche: number): number {
  const noches = calcularNoches(entrada, salida);
  return noches > 0 ? noches * precioPorNoche : 0;
}

// Devuelve el mensaje de error para mostrar en el formulario, o null si las fechas sirven.
export function validarFechas(entrada: string, salida: string): string | null {
  if (!entrada || !salida) return ERROR_FECHAS_OBLIGATORIAS;
  if (calcularNoches(entrada, salida) <= 0) return ERROR_SALIDA_NO_POSTERIOR;
  return null;
}

type TraerReservas = (url: string) => Promise<{ reservations: Reservation[] | null }>;

// Trae las reservas del usuario. El cliente HTTP entra por parámetro
// para poder probarla sin la API levantada.
export async function misReservas(userId: string, traer: TraerReservas): Promise<Reservation[]> {
  const data = await traer(`/api/reservations/users/${userId}/myreservations`);
  return data.reservations || [];
}

// Cómo se le muestra al usuario una reserva, según su estado y sus fechas.
// "hoy" entra por parámetro (yyyy-MM-dd) para no depender del reloj.
export function etiquetaDeReserva(
  reserva: Pick<Reservation, 'status' | 'start_date' | 'end_date'>,
  hoy: string
): string {
  if (reserva.status === 'canceled') {
    return 'Cancelada';
  }
  if (!reserva.start_date || !reserva.end_date) {
    return 'Sin fechas';
  }
  if (hoy < reserva.start_date) {
    return 'Próxima';
  }
  if (hoy < reserva.end_date) {
    return 'En curso';
  }
  return 'Finalizada';
}
