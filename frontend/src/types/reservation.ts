export interface Reservation {
  id: string;
  user_id: string;
  room_id: string;
  check_in_date: string;
  check_out_date: string;
  total_price: number;
  status: 'pending' | 'confirmed' | 'cancelled' | 'completed';
  guest_name: string;
  guest_email: string;
  guest_phone?: string;
  special_requests?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateReservationRequest {
  user_id: string;
  room_id: string;
  check_in_date: string;
  check_out_date: string;
  total_price: number;
  guest_name: string;
  guest_email: string;
  guest_phone?: string;
  special_requests?: string;
}

export interface CreateReservationResponse {
  message: string;
  reservation: Reservation;
}

export interface CancelReservationRequest {
  reason: string;
}

export interface GetMyReservationsResponse {
  reservations: Reservation[];
  total: number;
}