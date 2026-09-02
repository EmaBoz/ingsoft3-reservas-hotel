import api from '@/lib/axios';
import type { Room, RoomFilter, RoomListResponse, CreateRoomRequest } from '@/types';

export const roomService = {
  async getAllRooms(filters?: RoomFilter): Promise<RoomListResponse> {
    const params = new URLSearchParams();

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }

    const response = await api.get<RoomListResponse>(`/api/v1/rooms?${params.toString()}`);
    return response.data;
  },

  async getAvailableRooms(filters?: RoomFilter): Promise<RoomListResponse> {
    const params = new URLSearchParams();

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }

    const response = await api.get<RoomListResponse>(`/api/v1/rooms/available?${params.toString()}`);
    return response.data;
  },
  async getRoomsPrice(id: string): Promise<number> {
    const response = await api.get(`/api/v1/rooms/${id}`);
    return response.data.price;
  },
  async getRoomById(id: string): Promise<Room> {
    const response = await api.get<Room>(`/api/v1/rooms/${id}`);
    return response.data;
  },

  async getRoomByNumber(number: string): Promise<Room> {
    const response = await api.get<Room>(`/api/v1/rooms/number/${number}`);
    return response.data;
  },

  async createRoom(roomData: CreateRoomRequest): Promise<Room> {
    const response = await api.post<Room>('/api/v1/rooms', roomData);
    return response.data;
  },

  async updateRoom(id: string, roomData: Partial<CreateRoomRequest>): Promise<Room> {
    const response = await api.put<Room>(`/api/v1/rooms/${id}`, roomData);
    return response.data;
  },

  async updateRoomStatus(id: string, status: string): Promise<Room> {
    const response = await api.patch<Room>(`/api/v1/rooms/${id}/status`, { status });
    return response.data;
  },

  async deleteRoom(id: string): Promise<void> {
    await api.delete(`/api/v1/rooms/${id}`);
  },
};
