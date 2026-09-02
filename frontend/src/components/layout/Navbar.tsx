import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/Button';
import { Hotel, User, LogOut, Home, Calendar, Settings } from 'lucide-react';

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="bg-white shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link to="/" className="flex items-center space-x-2 text-primary-600 hover:text-primary-700">
              <Hotel className="w-8 h-8" />
              <span className="text-xl font-bold">Grand Prix</span>
            </Link>

            {isAuthenticated && (
              <div className="ml-10 flex items-center space-x-4">
                <Link
                  to="/"
                  className="flex items-center space-x-1 text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium"
                >
                  <Home className="w-4 h-4" />
                  <span>Inicio</span>
                </Link>
                <Link
                  to="/rooms"
                  className="flex items-center space-x-1 text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium"
                >
                  <Hotel className="w-4 h-4" />
                  <span>Habitaciones</span>
                </Link>
                <Link
                  to="/my-reservations"
                  className="flex items-center space-x-1 text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium"
                >
                  <Calendar className="w-4 h-4" />
                  <span>Mis Reservas</span>
                </Link>
                {user?.role === 'admin' && (
                  <Link
                    to="/admin/rooms"
                    className="flex items-center space-x-1 bg-blue-600 text-white hover:bg-blue-700 px-3 py-2 rounded-md text-sm font-medium"
                  >
                    <Settings className="w-4 h-4" />
                    <span>Admin Panel</span>
                  </Link>
                )}
              </div>
            )}
          </div>

          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <div className="flex items-center space-x-2 text-gray-700">
                  <User className="w-5 h-5" />
                  <span className="text-sm font-medium">
                    {user?.first_name} {user?.last_name}
                  </span>
                  {user?.role === 'admin' && (
                    <span className="ml-2 px-2 py-1 text-xs font-semibold text-white bg-primary-600 rounded">
                      Admin
                    </span>
                  )}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleLogout}
                  className="flex items-center space-x-1"
                >
                  <LogOut className="w-4 h-4" />
                  <span>Salir</span>
                </Button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="outline" size="sm">
                    Iniciar Sesión
                  </Button>
                </Link>
                <Link to="/register">
                  <Button variant="primary" size="sm">
                    Registrarse
                  </Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
