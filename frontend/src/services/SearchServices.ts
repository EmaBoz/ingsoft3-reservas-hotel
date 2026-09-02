// src/services/searchService.ts
import api from '@/lib/axios';
import type { Room, RoomFilter } from '@/types';

export interface SearchRoomsParams extends RoomFilter {
  q?: string;
  status?: RoomFilter['status'];
}

export interface SearchRoomsResponse {
  success: boolean;
  total: number;
  page: number;
  limit: number;
  rooms: Room[];
}

export const searchService = {
  async searchRooms(params: SearchRoomsParams): Promise<SearchRoomsResponse> {
    const response = await api.get<SearchRoomsResponse>(
      '/search-api/api/search/rooms',
      { params }
    );
    return response.data;
  },
};
