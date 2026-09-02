import api from '@/lib/axios';
import type {
  Reservation,
  CreateReservationRequest,
  CreateReservationResponse,
  CancelReservationRequest
} from '@/types';
import type { get } from 'node_modules/axios/index.d.cts';

export const reservationService = {
  async getGuestName(reservationId: string): Promise<string> {
    const response = await api.get<{ guest_id: string }>(`/api/reservations/${reservationId}`);
    const user_id = response.data.guest_id;
    const userResponse = await api.get<{ name: string }>(`/api/users/${user_id}`);
    
    return userResponse.data.name;
  },
  async getGuestEmail(reservationId: string): Promise<string> {
    const response = await api.get<{ guest_id: string }>(`/api/reservations/${reservationId}`);
    const user_id = response.data.guest_id;
    const userResponse = await api.get<{ email: string }>(`/api/users/${user_id}`);
    return userResponse.data.email;
  },
  async getAllReservations(): Promise<Reservation[]> {
    const response = await api.get<Reservation[]>('/api/reservations');

    return response.data;
  },
  async createReservation(reservationData: CreateReservationRequest): Promise<CreateReservationResponse> {
    const response = await api.post<CreateReservationResponse>('/api/reservations', reservationData);
    return response.data;
  },

  async getReservationById(id: string): Promise<{ reservation: Reservation }> {
    const response = await api.get<{ reservation: Reservation }>(`/api/reservations/${id}`);
    return response.data;
  },

  async getMyReservations(userId: string): Promise<{ reservations: Reservation[]}> {
    const response = await api.get<{ reservations: Reservation[]}>(
      `/api/reservations/users/${userId}/myreservations`
    );
    return response.data;
  },

  async cancelReservation(id: string, reason: string): Promise<void> {
    const cancelData: CancelReservationRequest = { reason };
    await api.delete(`/api/reservations/${id}`, { data: cancelData });
  },
};