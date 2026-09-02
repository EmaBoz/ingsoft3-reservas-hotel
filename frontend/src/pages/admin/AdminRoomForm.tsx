import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { adminRoomService, type UpdateRoomRequest } from '@/services/adminRoomService';
import type { CreateRoomRequest } from '@/types';

export default function AdminRoomForm() {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const { user } = useAuth();
    const isEditMode = !!id;

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [formData, setFormData] = useState<CreateRoomRequest>({
        number: '',
        type: 'standard',
        price: 0,
        description: '',
        capacity: 1,
        floor: 1,
        has_wifi: false,
        has_ac: false,
        has_tv: false,
        has_minibar: false,
    });

    // Redirect if not admin
    useEffect(() => {
        if (!user || user.role !== 'admin') {
            navigate('/');
        }
    }, [user, navigate]);

    // Load room data if editing
    useEffect(() => {
        if (isEditMode && id) {
            loadRoom(id);
        }
    }, [id, isEditMode]);

    const loadRoom = async (roomId: string) => {
        try {
            setLoading(true);
            const room = await adminRoomService.getRoomById(roomId);
            setFormData({
                number: room.number,
                type: room.type,
                price: room.price,
                description: room.description,
                capacity: room.capacity,
                floor: room.floor,
                has_wifi: room.has_wifi,
                has_ac: room.has_ac,
                has_tv: room.has_tv,
                has_minibar: room.has_minibar,
            });
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to load room');
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        // Validation
        if (!formData.number.trim()) {
            setError('Room number is required');
            return;
        }
        if (formData.price <= 0) {
            setError('Price must be greater than 0');
            return;
        }
        if (formData.capacity < 1) {
            setError('Capacity must be at least 1');
            return;
        }
        if (formData.floor < 1) {
            setError('Floor must be at least 1');
            return;
        }

        try {
            setLoading(true);
            if (isEditMode && id) {
                await adminRoomService.updateRoom(id, formData as UpdateRoomRequest);
            } else {
                await adminRoomService.createRoom(formData);
            }
            navigate('/admin/rooms');
        } catch (err: any) {
            setError(err.response?.data?.error || err.response?.data?.message || 'Failed to save room');
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
    ) => {
        const { name, value, type } = e.target;

        if (type === 'checkbox') {
            const checked = (e.target as HTMLInputElement).checked;
            setFormData((prev) => ({ ...prev, [name]: checked }));
        } else if (type === 'number') {
            setFormData((prev) => ({ ...prev, [name]: parseFloat(value) || 0 }));
        } else {
            setFormData((prev) => ({ ...prev, [name]: value }));
        }
    };

    if (loading && isEditMode) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-xl">Cargando habitación...</div>
            </div>
        );
    }

    return (
        <div className="container mx-auto px-4 py-8 max-w-3xl">
            <div className="mb-6">
                <button
                    onClick={() => navigate('/admin/rooms')}
                    className="text-blue-600 hover:text-blue-800 flex items-center gap-2"
                >
                    ← Volver a Habitaciones
                </button>
            </div>

            <div className="bg-white rounded-lg shadow-lg p-8">
                <h1 className="text-3xl font-bold mb-6">
                    {isEditMode ? 'Edit Room' : 'Add New Room'}
                </h1>

                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                    {/* Room Number */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Número de Habitación *
                        </label>
                        <input
                            type="text"
                            name="number"
                            value={formData.number}
                            onChange={handleChange}
                            required
                            className="w-full border border-gray-300 rounded-md px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            placeholder="e.g., 101"
                        />
                    </div>

                    {/* Room Type */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Tipo de Habitación *
                        </label>
                        <select
                            name="type"
                            value={formData.type}
                            onChange={handleChange}
                            required
                            className="w-full border border-gray-300 rounded-md px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="standard">Estándar</option>
                            <option value="single">Individual</option>
                            <option value="double">Doble</option>
                            <option value="suite">Suite</option>
                            <option value="deluxe">De lujo</option>
                        </select>
                    </div>

                    {/* Price and Capacity */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Precio por Noche ($) *
                            </label>
                            <input
                                type="number"
                                name="price"
                                value={formData.price}
                                onChange={handleChange}
                                required
                                min="0"
                                step="0.01"
                                className="w-full border border-gray-300 rounded-md px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                placeholder="0.00"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Capacidad (Huéspedes) *
                            </label>
                            <input
                                type="number"
                                name="capacity"
                                value={formData.capacity}
                                onChange={handleChange}
                                required
                                min="1"
                                className="w-full border border-gray-300 rounded-md px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                placeholder="1"
                            />
                        </div>
                    </div>

                    {/* Floor */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Piso *
                        </label>
                        <input
                            type="number"
                            name="floor"
                            value={formData.floor}
                            onChange={handleChange}
                            required
                            min="1"
                            className="w-full border border-gray-300 rounded-md px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            placeholder="1"
                        />
                    </div>

                    {/* Description */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Descripción
                        </label>
                        <textarea
                            name="description"
                            value={formData.description}
                            onChange={handleChange}
                            rows={4}
                            className="w-full border border-gray-300 rounded-md px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            placeholder="Enter room description..."
                        />
                    </div>

                    {/* Amenities */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-3">
                            Amenidades
                        </label>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <label className="flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    name="has_wifi"
                                    checked={formData.has_wifi}
                                    onChange={handleChange}
                                    className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                                />
                                <span className="text-sm">WiFi</span>
                            </label>

                            <label className="flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    name="has_ac"
                                    checked={formData.has_ac}
                                    onChange={handleChange}
                                    className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                                />
                                <span className="text-sm">Aire Acondicionado</span>
                            </label>

                            <label className="flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    name="has_tv"
                                    checked={formData.has_tv}
                                    onChange={handleChange}
                                    className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                                />
                                <span className="text-sm">TV</span>
                            </label>

                            <label className="flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    name="has_minibar"
                                    checked={formData.has_minibar}
                                    onChange={handleChange}
                                    className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                                />
                                <span className="text-sm">Minibar</span>
                            </label>
                        </div>
                    </div>

                    {/* Submit Buttons */}
                    <div className="flex justify-end gap-4 pt-6">
                        <button
                            type="button"
                            onClick={() => navigate('/admin/rooms')}
                            className="px-6 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={loading}
                            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
                        >
                            {loading ? 'Saving...' : isEditMode ? 'Update Room' : 'Create Room'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
