import api from '@/lib/axios';
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User } from '@/types';

export const authService = {

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await api.post<LoginResponse>('/login', credentials);
      return response.data;
    } catch (err: any) {
      throw err.response?.data?.error || "Credenciales inválidas";
    }
  },

  async register(userData: RegisterRequest): Promise<RegisterResponse> {
    try {
      const response = await api.post<RegisterResponse>('/users', userData);
      return response.data;
    } catch (err: any) {
      throw err.response?.data?.error || "Error al registrarte";
    }
  },

  async getUser(id: number): Promise<{ user: User }> {
    const response = await api.get<{ user: User }>(`/users/${id}`);
    return response.data;
  },

  async getAllUsers(): Promise<{ users: User[] }> {
    const response = await api.get<{ users: User[] }>('/users');
    return response.data;
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getStoredUser(): User | null {
    const userStr = localStorage.getItem('user');
    if (!userStr) return null;
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  },

  getStoredToken(): string | null {
    return localStorage.getItem('token');
  },

  isAuthenticated(): boolean {
    return !!this.getStoredToken();
  },
};
