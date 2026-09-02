import api from '@/lib/axios';
import type { Room, CreateRoomRequest } from '@/types';

const API_URL = '/api/v1/admin/rooms';

export interface UpdateRoomRequest {
    number: string;
    type: string;
    price: number;
    description: string;
    capacity: number;
    floor: number;
    has_wifi: boolean;
    has_ac: boolean;
    has_tv: boolean;
    has_minibar: boolean;
}

export const adminRoomService = {
    getAllRooms: async (): Promise<Room[]> => {
        const response = await api.get(API_URL);
        return response.data.rooms || [];
    },

    getRoomById: async (id: string): Promise<Room> => {
        const response = await api.get(`${API_URL}/${id}`);
        return response.data;
    },

    createRoom: async (room: CreateRoomRequest): Promise<Room> => {
        const response = await api.post(API_URL, room);
        return response.data;
    },

    updateRoom: async (id: string, room: UpdateRoomRequest): Promise<Room> => {
        const response = await api.put(`${API_URL}/${id}`, room);
        return response.data;
    },

    updateRoomStatus: async (id: string, status: string): Promise<Room> => {
        const response = await api.patch(`${API_URL}/${id}/status`, { status });
        return response.data;
    },

    deleteRoom: async (id: string): Promise<void> => {
        await api.delete(`${API_URL}/${id}`);
    },
};
