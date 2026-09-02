import React, { createContext, useContext, useState, useEffect } from 'react';
import { authService } from '@/services/authService';
import type { User, LoginRequest, RegisterRequest } from '@/types';
import type {ReactNode} from 'react';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const storedUser = authService.getStoredUser();
    const storedToken = authService.getStoredToken();

    if (storedUser && storedToken) {
      setUser(storedUser);
    }

    setIsLoading(false);
  }, []);

  const login = async (credentials: LoginRequest) => {
    try {
      const response = await authService.login(credentials);

      if (!response || !response.login) {
        throw new Error("Credenciales inválidas");
      }

      const { token, user: userData } = response.login;

      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(userData));
      setUser(userData);

    } catch (err: any) {
      // se reenvía el error para que Login.tsx pueda mostrarlo en pantalla
      throw err.response?.data?.error || "Error al iniciar sesión. Verifica tus credenciales.";
    }
  };

  const register = async (userData: RegisterRequest) => {
    const response = await authService.register(userData);
    await login({ username_or_email: userData.username, password: userData.password });
  };

  const logout = () => {
    authService.logout();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within an AuthProvider');
  return context;
}
