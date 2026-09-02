import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Layout } from '@/components/layout/Layout';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { roomService } from '@/services/roomService';
import { reservationService } from '@/services/reservationService';
import { useAuth } from '@/context/AuthContext';
import type { Room } from '@/types';
import { format, addDays, differenceInDays } from 'date-fns';
import { Calendar, Users, MapPin, DollarSign } from 'lucide-react';

export function ReserveRoom() {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [room, setRoom] = useState<Room | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const today = format(new Date(), 'yyyy-MM-dd');
  const tomorrow = format(addDays(new Date(), 1), 'yyyy-MM-dd');

  const [formData, setFormData] = useState({
    start_date: today,
    end_date: tomorrow,
  });

  useEffect(() => {
    if (roomId) {
      fetchRoom();
    }
  }, [roomId]);

  const fetchRoom = async () => {
    if (!roomId) return;

    setIsLoading(true);
    setError('');
    try {
      const data = await roomService.getRoomById(roomId);
      setRoom(data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al cargar la habitación');
    } finally {
      setIsLoading(false);
    }
  };

  const calculateTotal = () => {
    if (!room || !formData.start_date || !formData.end_date) return 0;

    const nights = differenceInDays(
      new Date(formData.end_date),
      new Date(formData.start_date)
    );

    return nights > 0 ? nights * room.price : 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!user || !room) return;

    const nights = differenceInDays(
      new Date(formData.end_date),
      new Date(formData.start_date)
    );

    if (nights <= 0) {
      setError('La fecha de salida debe ser posterior a la fecha de entrada');
      return;
    }

    setIsSubmitting(true);
    setError('');

    try {
      await reservationService.createReservation({
        user_id: user.id,
        room_id: parseInt(room.id),
        start_date: formData.start_date,
        end_date: formData.end_date,
      });

      setSuccess(true);
      setTimeout(() => {
        navigate('/my-reservations');
      }, 2000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al crear la reserva');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <Layout>
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <p className="mt-4 text-gray-600">Cargando información...</p>
        </div>
      </Layout>
    );
  }

  if (error && !room) {
    return (
      <Layout>
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      </Layout>
    );
  }

  if (!room) {
    return <Layout><div>Habitación no encontrada</div></Layout>;
  }

  const nights = differenceInDays(
    new Date(formData.end_date),
    new Date(formData.start_date)
  );
  const total = calculateTotal();

  return (
    <Layout>
      <div className="max-w-4xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold text-gray-900">Reservar Habitación</h1>

        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg">
            Reserva creada exitosamente. Redirigiendo a tus reservas...
          </div>
        )}

        <div className="grid md:grid-cols-2 gap-6">
          {/* Room Info */}
          <Card>
            <CardHeader>
              <CardTitle>Detalles de la Habitación</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="h-48 bg-gradient-to-br from-primary-400 to-primary-600 rounded-lg flex items-center justify-center">
                  <div className="text-center text-white">
                    <div className="text-5xl font-bold">#{room.number}</div>
                    <div className="text-lg mt-2">{room.type.toUpperCase()}</div>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center space-x-2 text-gray-600">
                      <Users className="w-4 h-4" />
                      <span>Capacidad</span>
                    </div>
                    <span className="font-medium">{room.capacity} personas</span>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center space-x-2 text-gray-600">
                      <MapPin className="w-4 h-4" />
                      <span>Piso</span>
                    </div>
                    <span className="font-medium">{room.floor}</span>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center space-x-2 text-gray-600">
                      <DollarSign className="w-4 h-4" />
                      <span>Precio por noche</span>
                    </div>
                    <span className="font-medium">${room.price.toFixed(2)}</span>
                  </div>
                </div>

                <div className="border-t pt-4">
                  <p className="text-sm text-gray-600">
                    {room.description || 'Habitación cómoda y acogedora con todas las comodidades necesarias.'}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Reservation Form */}
          <Card>
            <CardHeader>
              <CardTitle>Información de Reserva</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
                    {error}
                  </div>
                )}

                <Input
                  label="Fecha de Entrada"
                  type="date"
                  value={formData.start_date}
                  onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                  min={today}
                  required
                />

                <Input
                  label="Fecha de Salida"
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  min={formData.start_date || tomorrow}
                  required
                />

                <div className="bg-gray-50 rounded-lg p-4 space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Noches:</span>
                    <span className="font-medium">{nights > 0 ? nights : 0}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Precio por noche:</span>
                    <span className="font-medium">${room.price.toFixed(2)}</span>
                  </div>
                  <div className="border-t pt-2 flex justify-between">
                    <span className="font-semibold">Total:</span>
                    <span className="text-xl font-bold text-primary-600">
                      ${total.toFixed(2)}
                    </span>
                  </div>
                </div>

                <div className="space-y-2">
                  <Button
                    type="submit"
                    fullWidth
                    disabled={isSubmitting || nights <= 0 || success}
                  >
                    {isSubmitting ? 'Procesando...' : 'Confirmar Reserva'}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    fullWidth
                    onClick={() => navigate('/rooms')}
                    disabled={isSubmitting}
                  >
                    Cancelar
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </Layout>
  );
}
