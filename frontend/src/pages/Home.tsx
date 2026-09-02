import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Layout } from '@/components/layout/Layout';
import { Button } from '@/components/ui/Button';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Hotel, Calendar, Star, Wifi, Tv, Wind, Coffee } from 'lucide-react';

export function Home() {
  const { user, isAuthenticated } = useAuth();

  return (
    <Layout>
      <div className="space-y-12">
        {/* Hero Section */}
        <div className="bg-gradient-to-r from-primary-600 to-primary-800 rounded-xl shadow-xl p-12 text-white">
          <div className="max-w-3xl">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Bienvenido a Grand Prix
            </h1>
            <p className="text-xl mb-8 text-primary-100">
              {isAuthenticated
                ? `Hola ${user?.first_name}, descubre nuestras mejores habitaciones y reserva tu estadía perfecta.`
                : 'Descubre la comodidad y el lujo en cada habitación. Reserva ahora y disfruta de una experiencia inolvidable.'}
            </p>
            <div className="flex gap-4">
              <Link to="/rooms">
                <Button variant="outline" size="lg" className="bg-white text-primary-600 hover:bg-gray-100">
                  Ver Habitaciones
                </Button>
              </Link>
              {!isAuthenticated && (
                <Link to="/register">
                  <Button size="lg" className="bg-primary-700 hover:bg-primary-800">
                    Registrarse
                  </Button>
                </Link>
              )}
            </div>
          </div>
        </div>

        {/* Features */}
        <div>
          <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">
            ¿Por qué elegir Grand Prix?
          </h2>
          <div className="grid md:grid-cols-3 gap-6">
            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center">
                  <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mb-4">
                    <Star className="w-8 h-8 text-primary-600" />
                  </div>
                  <h3 className="text-xl font-semibold mb-2">Calidad Premium</h3>
                  <p className="text-gray-600">
                    Habitaciones de alta calidad con las mejores comodidades para tu estadía.
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center">
                  <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mb-4">
                    <Calendar className="w-8 h-8 text-primary-600" />
                  </div>
                  <h3 className="text-xl font-semibold mb-2">Reservas Fáciles</h3>
                  <p className="text-gray-600">
                    Sistema de reservas intuitivo y rápido. Confirma tu habitación en minutos.
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center">
                  <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mb-4">
                    <Hotel className="w-8 h-8 text-primary-600" />
                  </div>
                  <h3 className="text-xl font-semibold mb-2">Excelente Ubicación</h3>
                  <p className="text-gray-600">
                    Ubicación privilegiada con acceso a las mejores atracciones de la ciudad.
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Amenities */}
        <div>
          <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">
            Nuestras Comodidades
          </h2>
          <Card>
            <CardContent className="pt-6">
              <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                <div className="flex flex-col items-center text-center">
                  <Wifi className="w-10 h-10 text-primary-600 mb-2" />
                  <span className="font-medium">WiFi Gratis</span>
                </div>
                <div className="flex flex-col items-center text-center">
                  <Tv className="w-10 h-10 text-primary-600 mb-2" />
                  <span className="font-medium">TV Premium</span>
                </div>
                <div className="flex flex-col items-center text-center">
                  <Wind className="w-10 h-10 text-primary-600 mb-2" />
                  <span className="font-medium">Aire Acondicionado</span>
                </div>
                <div className="flex flex-col items-center text-center">
                  <Coffee className="w-10 h-10 text-primary-600 mb-2" />
                  <span className="font-medium">Minibar</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* CTA Section */}
        <div className="bg-gray-100 rounded-xl p-8 text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            ¿Listo para reservar?
          </h2>
          <p className="text-gray-600 mb-6">
            Explora nuestras habitaciones disponibles y encuentra la perfecta para ti.
          </p>
          <Link to="/rooms">
            <Button size="lg">
              Explorar Habitaciones
            </Button>
          </Link>
        </div>
      </div>
    </Layout>
  );
}
