import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { adminRoomService } from '@/services/adminRoomService';
import type { Room } from '@/types';
import { useAuth } from '@/context/AuthContext';

export default function AdminRoomList() {
    const navigate = useNavigate();
    const { user } = useAuth();
    const [rooms, setRooms] = useState<Room[]>([]);
    const [filteredRooms, setFilteredRooms] = useState<Room[]>([]);
    const [filters, setFilters] = useState({
        type: '',
        status: '',
        floor: '',
        maxPrice: '',
    });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [updatingRoomId, setUpdatingRoomId] = useState<string | null>(null);
    const [deleteModal, setDeleteModal] = useState<{ show: boolean; roomId: string | null }>({
        show: false,
        roomId: null,
    });

    useEffect(() => {
        if (!user || user.role !== 'admin') {
            navigate('/');
        }
    }, [user, navigate]);

    useEffect(() => {
        loadRooms();
    }, []);

    const loadRooms = async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await adminRoomService.getAllRooms();
            setRooms(data);
            setFilteredRooms(data);
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to load rooms');
        } finally {
            setLoading(false);
        }
    };

    const filterRooms = () => {
        let result = [...rooms];

        if (filters.type) {
            result = result.filter((room) => room.type === filters.type);
        }

        if (filters.status) {
            result = result.filter((room) => room.status === filters.status);
        }

        if (filters.floor) {
            result = result.filter((room) => room.floor === parseInt(filters.floor));
        }

        if (filters.maxPrice) {
            result = result.filter((room) => room.price <= parseFloat(filters.maxPrice));
        }

        setFilteredRooms(result);
    };

    useEffect(() => {
        filterRooms();
    }, [filters, rooms]);

    const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement | HTMLInputElement>) => {
        const { name, value } = e.target;
        setFilters((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const resetFilters = () => {
        setFilters({
            type: '',
            status: '',
            floor: '',
            maxPrice: '',
        });
    };

    const handleDelete = async () => {
        if (!deleteModal.roomId) return;

        try {
            await adminRoomService.deleteRoom(deleteModal.roomId);
            setDeleteModal({ show: false, roomId: null });
            loadRooms();
        } catch (err: any) {
            alert(err.response?.data?.error || 'Failed to delete room');
        }
    };

    const handleStatusChange = async (roomId: string, newStatus: string) => {
        if (updatingRoomId) return; // Prevenir múltiples clics

        setUpdatingRoomId(roomId);
        setError(null);
        setSuccess(null);

        try {
            const response = await adminRoomService.updateRoomStatus(roomId, newStatus);
            setSuccess(`✅ Estado de la habitación actualizado a: ${newStatus}`);

            // Actualizar la habitación en la lista localmente
            setRooms((prevRooms) =>
                prevRooms.map((room) =>
                    room.id === roomId ? { ...room, status: newStatus as any } : room
                )
            );

            // Limpiar mensaje de éxito después de 3 segundos
            setTimeout(() => {
                setSuccess(null);
            }, 3000);

            // Recargar después de un pequeño delay para sincronizar con el backend
            setTimeout(() => {
                loadRooms();
            }, 1000);
        } catch (err: any) {
            const errorMessage =
                err?.response?.status === 401
                    ? 'No tienes permisos para realizar esta acción'
                    : err?.response?.status === 404
                        ? 'La habitación no fue encontrada'
                        : err?.response?.data?.error || 'Error al actualizar el estado de la habitación';

            setError(errorMessage);
            console.error('Error updating room status:', err);
        } finally {
            setUpdatingRoomId(null);
        }
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="text-xl">Loading rooms...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 py-8">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                {/* Header */}
                <div className="flex justify-between items-center mb-8">
                    <div>
                        <h1 className="text-3xl font-bold text-gray-900">Manejo de habitaciones</h1>
                        <p className="mt-2 text-sm text-gray-600">Administra todas las habitaciones del hotel</p>
                    </div>
                    <div className="flex gap-4">
                        <Link
                            to="/"
                            className="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors"
                        >
                            ← Volver al Inicio
                        </Link>
                        <Link
                            to="/admin/rooms/new"
                            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                        >
                            + agregar Habitación
                        </Link>
                    </div>
                </div>

                {error && (
                    <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
                        <p className="text-red-800">{error}</p>
                    </div>
                )}

                {success && (
                    <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-md">
                        <p className="text-green-800">{success}</p>
                    </div>
                )}

                {/* Filters */}
                <div className="bg-white rounded-lg shadow p-6 mb-6">
                    <h2 className="text-lg font-semibold mb-4">Filtros</h2>
                    <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Tipo</label>
                            <select
                                name="type"
                                value={filters.type}
                                onChange={handleFilterChange}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            >
                                <option value="">Todo</option>
                                <option value="single">Individual</option>
                                <option value="double">Doble</option>
                                <option value="suite">Suite</option>
                                <option value="deluxe">De lujo</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Estado</label>
                            <select
                                name="status"
                                value={filters.status}
                                onChange={handleFilterChange}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            >
                                <option value="">Todos los estados</option>
                                <option value="available">Disponible</option>
                                <option value="occupied">Ocupado</option>
                                <option value="maintenance">Mantenimiento</option>
                                <option value="cleaning">Limpieza</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Piso</label>
                            <select
                                name="floor"
                                value={filters.floor}
                                onChange={handleFilterChange}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            >
                                <option value="">Cualquier piso</option>
                                <option value="1">Piso 1</option>
                                <option value="2">Piso 2</option>
                                <option value="3">Piso 3</option>
                                <option value="4">Piso 4</option>
                                <option value="5">Piso 5</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Precio Máximo</label>
                            <input
                                type="number"
                                name="maxPrice"
                                value={filters.maxPrice}
                                onChange={handleFilterChange}
                                placeholder="Cualquier precio"
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>

                        <div className="flex items-end">
                            <button
                                onClick={resetFilters}
                                className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors"
                            >
                                Restablecer filtros
                            </button>
                        </div>
                    </div>
                </div>

                {/* Results count */}
                <div className="mb-4 text-sm text-gray-600">
                    Mostrando {filteredRooms.length} de {rooms.length} habitaciones
                </div>

                {/* Rooms Table */}
                <div className="bg-white rounded-lg shadow overflow-hidden">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Habitación #
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Tipo
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Piso
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Precio
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Capacidad
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Estado
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Amenities
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {filteredRooms.length === 0 ? (
                                <tr>
                                    <td colSpan={8} className="px-6 py-12 text-center text-gray-500">
                                        No se encontraron habitaciones {filters.type || filters.status || filters.floor || filters.maxPrice ? 'Intenta ajustar tus filtros.' : 'Crea tu primera habitación para comenzar.'}
                                    </td>
                                </tr>
                            ) : (
                                filteredRooms.map((room) => (
                                    <tr key={room.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {room.number}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize">
                                            {room.type}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            Piso {room.floor}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            ${room.price}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {room.capacity} persona{room.capacity > 1 ? 's' : ''}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <div className="relative">
                                                <select
                                                    value={room.status}
                                                    onChange={(e) => handleStatusChange(room.id, e.target.value)}
                                                    disabled={updatingRoomId === room.id}
                                                    className={`text-sm rounded-full px-3 py-1 font-semibold cursor-pointer transition-all ${updatingRoomId === room.id
                                                        ? 'opacity-50 cursor-not-allowed bg-gray-200 text-gray-600'
                                                        : room.status === 'available'
                                                            ? 'bg-green-100 text-green-800 hover:bg-green-200'
                                                            : room.status === 'occupied'
                                                                ? 'bg-red-100 text-red-800 hover:bg-red-200'
                                                                : room.status === 'maintenance'
                                                                    ? 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200'
                                                                    : 'bg-gray-100 text-gray-800 hover:bg-gray-200'
                                                        }`}
                                                >
                                                    <option value="available">Disponible</option>
                                                    <option value="occupied">Ocupado</option>
                                                    <option value="maintenance">Mantenimiento</option>
                                                    <option value="cleaning">Limpieza</option>
                                                </select>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            <div className="flex flex-wrap gap-1">
                                                {room.has_wifi && (
                                                    <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">WiFi</span>
                                                )}
                                                {room.has_ac && (
                                                    <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">AC</span>
                                                )}
                                                {room.has_tv && (
                                                    <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">TV</span>
                                                )}
                                                {room.has_minibar && (
                                                    <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">Minibar</span>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <Link
                                                to={`/admin/rooms/${room.id}`}
                                                className="text-blue-600 hover:text-blue-900 mr-4"
                                            >
                                                Edit
                                            </Link>
                                            <button
                                                onClick={() => setDeleteModal({ show: true, roomId: room.id })}
                                                className="text-red-600 hover:text-red-900"
                                            >
                                                Eliminar
                                            </button>
                                        </td>
                                    </tr>

                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {/* Delete Confirmation Modal */}
                {deleteModal.show && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                        <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
                            <h3 className="text-lg font-semibold mb-4">Confirmar</h3>
                            <p className="text-gray-600 mb-6">
                                ¿Estás seguro de que deseas eliminar esta habitación? Esta acción no se puede deshacer.
                            </p>
                            <div className="flex justify-end gap-3">
                                <button
                                    onClick={() => setDeleteModal({ show: false, roomId: null })}
                                    className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
                                >
                                    Cancelar
                                </button>
                                <button
                                    onClick={handleDelete}
                                    className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
                                >
                                    Eliminar
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
