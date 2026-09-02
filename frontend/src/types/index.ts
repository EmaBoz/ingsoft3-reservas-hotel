// User Types
export interface User {
  id: number;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: 'normal' | 'admin';
}

export interface LoginRequest {
  username_or_email: string;
  password: string;
}

export interface LoginResponse {
  login: {
    token: string;
    user: User;
  };
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface RegisterResponse {
  user: User;
}

// Room Types
export type RoomType = 'single' | 'double' | 'suite' | 'deluxe' | 'standard';
export type RoomStatus = 'available' | 'occupied' | 'maintenance' | 'reserved';

export interface Room {
  id: string;
  number: string;
  type: RoomType;
  status: RoomStatus;
  price: number;
  description: string;
  capacity: number;
  floor: number;
  has_wifi: boolean;
  has_ac: boolean;
  has_tv: boolean;
  has_minibar: boolean;
  created_at: string;
  updated_at: string;
}

export interface RoomFilter {
  type?: RoomType;
  status?: RoomStatus;
  floor?: number;
  min_price?: number;
  max_price?: number;
  has_wifi?: boolean;
  has_ac?: boolean;
  has_tv?: boolean;
  has_minibar?: boolean;
  page?: number;
  limit?: number;
}

export interface RoomListResponse {
  rooms: Room[];
  total: number;
  page: number;
  limit: number;
}

export interface CreateRoomRequest {
  number: string;
  type: RoomType;
  price: number;
  description?: string;
  capacity: number;
  floor: number;
  has_wifi?: boolean;
  has_ac?: boolean;
  has_tv?: boolean;
  has_minibar?: boolean;
}

// Reservation Types
export type ReservationStatus = 'active' | 'canceled';

export interface Reservation {
  id: string;
  user_id: number;
  room_id: number;
  start_date: string;
  end_date: string;
  status: ReservationStatus;
}

export interface CreateReservationRequest {
  user_id: number;
  room_id: number;
  start_date: string;
  end_date: string;
}

export interface CreateReservationResponse {
  reservation: Reservation;
}

export interface CancelReservationRequest {
  reason: string;
}

// API Error Response
export interface ApiError {
  error: string;
  details?: string;
}
