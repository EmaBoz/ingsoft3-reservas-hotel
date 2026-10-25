import React, { useEffect, useState } from 'react';
import { Layout } from '@/components/layout/Layout';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Calendar, MapPin, X, RefreshCw } from 'lucide-react';
import { reservationService, traerJson } from '@/services/reservationService';
import { misReservas, etiquetaDeReserva } from '@/lib/reservas';
import { format } from 'date-fns';
import { authService } from '@/services/authService';
import type { Reservation } from '@/types';

export function MyReservations() {
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [cancelingId, setCancelingId] = useState<string | null>(null);

  const fetchReservations = async () => {
    try {
      setLoading(true);
      setError(null);
      const user = authService.getStoredUser();
      
      if (!user) {
        setReservations([]);
        return;
      }
      
      setReservations(await misReservas(String(user.id), traerJson));
    } catch (e: any) {
      setError(e?.response?.data?.error || 'Error cargando reservas');
      console.error('Error fetching reservations:', e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReservations();
  }, []);

  const handleCancel = async (id: string) => {
    if (!confirm('¿Estás seguro de que deseas cancelar esta reserva?')) {
      return;
    }

    try {
      setCancelingId(id);
      await reservationService.cancelReservation(id, 'Canceled by user');
      // Recargar las reservas después de cancelar
      await fetchReservations();
    } catch (e: any) {
      setError(e?.response?.data?.error || 'No se pudo cancelar la reserva');
      console.error('Error canceling reservation:', e);
    } finally {
      setCancelingId(null);
    }
  };

  const hoy = format(new Date(), 'yyyy-MM-dd');

  const getStatusBadge = (etiqueta: string) => {
    const colores: Record<string, string> = {
      'Próxima': 'bg-green-100 text-green-700',
      'En curso': 'bg-blue-100 text-blue-700',
      'Finalizada': 'bg-gray-100 text-gray-700',
      'Cancelada': 'bg-red-100 text-red-700',
    };

    return (
      <span className={`px-3 py-1 ${colores[etiqueta] || 'bg-gray-100 text-gray-700'} rounded-full text-sm font-medium`}>
        {etiqueta}
      </span>
    );
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-900">Mis Reservas</h1>
          <Button
            variant="outline"
            size="sm"
            onClick={fetchReservations}
            disabled={loading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
            Actualizar
          </Button>
        </div>

        {loading && (
          <Card>
            <CardContent className="py-12">
              <div className="text-center text-gray-600">
                <Calendar className="w-16 h-16 mx-auto mb-4 text-gray-400 animate-pulse" />
                <p className="text-lg font-medium mb-2">Cargando tus reservas...</p>
              </div>
            </CardContent>
          </Card>
        )}

        {!loading && error && (
          <Card>
            <CardContent className="py-6">
              <div className="text-center">
                <p className="text-red-600 text-sm mb-3">{error}</p>
                <Button variant="outline" size="sm" onClick={fetchReservations}>
                  Reintentar
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {!loading && !error && reservations.length === 0 && (
          <Card>
            <CardContent className="py-12">
              <div className="text-center text-gray-600">
                <Calendar className="w-16 h-16 mx-auto mb-4 text-gray-400" />
                <p className="text-lg font-medium mb-2">No tienes reservas activas</p>
                <p className="text-sm">Cuando realices una reserva, aparecerá aquí.</p>
              </div>
            </CardContent>
          </Card>
        )}

        {!loading && !error && reservations.length > 0 && (
          <div className="space-y-4">
            {reservations.map((reservation) => (
              <Card key={reservation.id}>
                <CardContent className="p-6">
                  <div className="flex justify-between items-start">
                    <div className="space-y-2">
                      <div className="flex items-center space-x-4">
                        <div className="text-2xl font-bold text-primary-600">
                          #{reservation.id.slice(-6)}
                        </div>
                        {getStatusBadge(etiquetaDeReserva(reservation, hoy))}
                      </div>

                      <div className="flex items-center space-x-6 text-sm text-gray-600">
                        <div className="flex items-center space-x-2">
                          <Calendar className="w-4 h-4" />
                          <span>{reservation.start_date} - {reservation.end_date}</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <MapPin className="w-4 h-4" />
                          <span>Habitación #{reservation.room_id}</span>
                        </div>
                      </div>
                    </div>

                    {reservation.status === 'active' && (
                      <Button
                        variant="danger"
                        size="sm"
                        disabled={cancelingId === reservation.id}
                        onClick={() => handleCancel(reservation.id)}
                        aria-label={`Cancelar la reserva de la habitación ${reservation.room_id}`}
                      >
                        <X className="w-4 h-4 mr-1" />
                        {cancelingId === reservation.id ? 'Cancelando...' : 'Cancelar'}
                      </Button>
                    )}
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