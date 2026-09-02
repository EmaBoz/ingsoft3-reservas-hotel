import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout } from '@/components/layout/Layout';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { useAuth } from '@/context/AuthContext';
import type { Room, RoomFilter, RoomType } from '@/types';
import { Wifi, Tv, Wind, Coffee, Users, MapPin, DollarSign, CheckCircle, RefreshCw } from 'lucide-react';
import { roomService } from '@/services/roomService';

export function Rooms() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [filters, setFilters] = useState<RoomFilter>({
    page: 1,
    limit: 12,
  });
  const [lastRefresh, setLastRefresh] = useState<number>(Date.now());

  // 🔄 Polling automático: Refrescar habitaciones cada 30 segundos
  useEffect(() => {
    const pollInterval = setInterval(() => {
      fetchRooms();
    }, 20000); // 30 segundos

    return () => clearInterval(pollInterval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]);

  // 👁️ Refrescar cuando el usuario regresa a la pestaña (después de hacer una reserva)
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        console.log('✅ Usuario regresó a la página de habitaciones, refrescando...');
        fetchRooms();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () =>
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]);

  // Refrescar cuando los filtros cambien
  useEffect(() => {
    fetchRooms();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]);

  const fetchRooms = async () => {
    setIsLoading(true);
    setError('');
    try {
      // 👉 Ahora pegamos contra rooms-api (getAvailableRooms)
      const response = await roomService.getAvailableRooms(filters);
      
      // ✅ FILTER CLIENT-SIDE: Asegurar que SOLO se muestren habitaciones disponibles
      // (El backend ya filtra, pero esto es medida de seguridad adicional)
      const availableRooms = (response.rooms || []).filter(
        (room) => room.status === 'available'
      );
      
      setRooms(availableRooms);
      setLastRefresh(Date.now()); // ⏰ Actualizar timestamp de última actualización
      
      // Si no hay resultados, mostrar mensaje amigable
      if (!availableRooms || availableRooms.length === 0) {
        setError('No se encontraron habitaciones con los filtros seleccionados');
      }
    } catch (err: any) {
      // Diferenciar tipos de errores
      const status = err?.response?.status;
      const errorMessage = err?.response?.data?.error || err?.message;

      if (status === 401 || status === 403) {
        setError('Tu sesión ha expirado. Por favor inicia sesión nuevamente.');
      } else if (status === 404) {
        setError('No se encontraron habitaciones disponibles.');
      } else if (status >= 500) {
        setError('Error del servidor. Por favor intenta más tarde.');
      } else if (status === 400) {
        setError(`Datos inválidos: ${errorMessage}`);
      } else if (err.message === 'Network Error') {
        setError('Error de conexión. Verifica tu conexión a internet.');
      } else {
        setError(errorMessage || 'Error desconocido al cargar las habitaciones');
      }
      
      setRooms([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleFilterChange = (key: keyof RoomFilter, value: any) => {
    setFilters((prev) => {
      const updated = {
        ...prev,
        [key]: value || undefined,
        page: 1,
      };

      // Validación: Si se cambia min_price, asegurar que no sea mayor que max_price
      if (key === 'min_price' && updated.min_price && updated.max_price) {
        if (updated.min_price > updated.max_price) {
          // Si el nuevo min es mayor que max, no actualizar
          setError('El precio mínimo no puede ser mayor que el precio máximo');
          return prev;
        }
      }

      // Validación: Si se cambia max_price, asegurar que no sea menor que min_price
      if (key === 'max_price' && updated.min_price && updated.max_price) {
        if (updated.max_price < updated.min_price) {
          // Si el nuevo max es menor que min, no actualizar
          setError('El precio máximo no puede ser menor que el precio mínimo');
          return prev;
        }
      }

      // Validación: El piso debe ser positivo
      if (key === 'floor' && value !== undefined && value !== null && value !== '') {
        if (parseInt(value, 10) < 1) {
          setError('El número de piso debe ser mayor a 0');
          return prev;
        }
      }

      setError(''); // Limpiar error si la validación pasa
      return updated;
    });
  };

  const handleReserve = (room: Room) => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    navigate(`/rooms/${room.id}/reserve`);
  };

  const getRoomTypeLabel = (type: RoomType) => {
    const labels: Record<RoomType, string> = {
      single: 'Individual',
      double: 'Doble',
      suite: 'Suite',
      deluxe: 'De Lujo',
      standard: 'Estándar',
    };

    return labels[type];
  };

  const handlePriceInput = (value: string, type: 'min' | 'max') => {
    if (value === '') {
      handleFilterChange(type === 'min' ? 'min_price' : 'max_price', undefined);
      return;
    }

    const numValue = parseFloat(value);
    if (isNaN(numValue) || numValue < 0) {
      setError('El precio debe ser un número positivo');
      return;
    }

    handleFilterChange(
      type === 'min' ? 'min_price' : 'max_price',
      numValue
    );
  };

  const handleFloorInput = (value: string) => {
    if (value === '') {
      handleFilterChange('floor', undefined);
      return;
    }

    const numValue = parseInt(value, 10);
    if (isNaN(numValue) || numValue < 1) {
      setError('El número de piso debe ser mayor a 0');
      return;
    }

    handleFilterChange('floor', numValue);
  };

  // 🔄 Función para refrescar manualmente
  const handleManualRefresh = async () => {
    setLastRefresh(Date.now());
    await fetchRooms();
  };

  // 📱 Función para obtener cuándo fue la última actualización
  const getLastUpdateText = () => {
    const now = Date.now();
    const diff = Math.floor((now - lastRefresh) / 1000); // segundos

    if (diff < 60) {
      return 'Hace unos segundos';
    }
    if (diff < 120) {
      return 'Hace 1 minuto';
    }
    if (diff < 3600) {
      const minutes = Math.floor(diff / 60);
      return `Hace ${minutes} minuto${minutes > 1 ? 's' : ''}`;
    }

    return 'Hace más de una hora';
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">
              Habitaciones Disponibles
            </h1>
            <p className="text-sm text-gray-500 mt-1">
              ✅ Actualizado {getLastUpdateText()}
            </p>
          </div>
          <button
            onClick={handleManualRefresh}
            disabled={isLoading}
            className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-all"
            title="Refrescar lista de habitaciones"
          >
            <RefreshCw
              className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`}
            />
            {isLoading ? 'Cargando...' : 'Refrescar'}
          </button>
        </div>

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>Filtros</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Tipo de Habitación
                </label>
                <select
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  value={filters.type || ''}
                  onChange={(e) => handleFilterChange('type', e.target.value)}
                >
                  <option value="">Todos</option>
                  <option value="single">Individual</option>
                  <option value="double">Doble</option>
                  <option value="suite">Suite</option>
                  <option value="deluxe">De Lujo</option>
                  <option value="standard">Estándar</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Precio Mínimo
                </label>
                <Input
                  type="number"
                  placeholder="$0"
                  min="0"
                  step="0.01"
                  value={filters.min_price || ''}
                  onChange={(e) => handlePriceInput(e.target.value, 'min')}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Precio Máximo
                </label>
                <Input
                  type="number"
                  placeholder="$999"
                  min="0"
                  step="0.01"
                  value={filters.max_price || ''}
                  onChange={(e) => handlePriceInput(e.target.value, 'max')}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Piso
                </label>
                <Input
                  type="number"
                  placeholder="Todos"
                  min="1"
                  value={filters.floor || ''}
                  onChange={(e) => handleFloorInput(e.target.value)}
                />
              </div>
            </div>

            <div className="mt-4 flex gap-4 flex-wrap">
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_wifi || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_wifi',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">WiFi</span>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_ac || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_ac',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">
                  Aire Acondicionado
                </span>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_tv || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_tv',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">TV</span>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_minibar || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_minibar',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">Minibar</span>
              </label>
            </div>

            <div className="mt-6 flex gap-3 justify-end">
              <Button
                variant="secondary"
                onClick={() => {
                  setFilters({ page: 1, limit: 12 });
                  setError('');
                }}
              >
                Limpiar Filtros
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Rooms Grid */}
        {error && !isLoading && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 flex items-start gap-3">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-red-400"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <div className="flex-1">
              <h3 className="text-sm font-medium text-red-800">
                {error.includes('No se encontraron')
                  ? 'Sin resultados'
                  : 'Error'}
              </h3>
              <p className="mt-1 text-sm text-red-700">{error}</p>
            </div>
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
            <p className="mt-4 text-gray-600">Cargando habitaciones...</p>
          </div>
        ) : rooms.length === 0 && !error ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-gray-600">
                No se encontraron habitaciones con los filtros seleccionados.
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {rooms.map((room) => (
              <Card
                key={room.id}
                className="hover:shadow-lg transition-shadow"
              >
                <CardContent className="p-0">
                  <div className="h-48 bg-gradient-to-br from-primary-400 to-primary-600 rounded-t-lg flex items-center justify-center">
                    <div className="text-center text-white">
                      <div className="text-5xl font-bold">#{room.number}</div>
                      <div className="text-lg mt-2">
                        {getRoomTypeLabel(room.type)}
                      </div>
                    </div>
                  </div>

                  <div className="p-6 space-y-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2 text-gray-600">
                        <Users className="w-4 h-4" />
                        <span className="text-sm">
                          {room.capacity} personas
                        </span>
                      </div>
                      <div className="flex items-center gap-2 text-green-600 bg-green-50 px-2 py-1 rounded-full">
                        <CheckCircle className="w-4 h-4" />
                        <span className="text-xs font-medium">Disponible</span>
                      </div>
                    </div>

                    <div className="flex items-center space-x-2 text-gray-600">
                      <MapPin className="w-4 h-4" />
                      <span className="text-sm">Piso {room.floor}</span>
                    </div>

                    <p className="text-gray-600 text-sm line-clamp-2">
                      {room.description || 'Habitación cómoda y acogedora.'}
                    </p>

                    <div className="flex gap-2 flex-wrap">
                      {room.has_wifi && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Wifi className="w-3 h-3" />
                          <span>WiFi</span>
                        </div>
                      )}
                      {room.has_tv && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Tv className="w-3 h-3" />
                          <span>TV</span>
                        </div>
                      )}
                      {room.has_ac && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Wind className="w-3 h-3" />
                          <span>A/C</span>
                        </div>
                      )}
                      {room.has_minibar && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Coffee className="w-3 h-3" />
                          <span>Minibar</span>
                        </div>
                      )}
                    </div>

                    <div className="border-t pt-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-1">
                          <DollarSign className="w-5 h-5 text-primary-600" />
                          <span className="text-2xl font-bold text-gray-900">
                            {room.price.toFixed(2)}
                          </span>
                          <span className="text-gray-600 text-sm">
                            /noche
                          </span>
                        </div>
                        <Button
                          variant="primary"
                          size="sm"
                          onClick={() => handleReserve(room)}
                        >
                          Reservar
                        </Button>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}
