import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '@/context/AuthContext';
import { ProtectedRoute } from '@/components/auth/ProtectedRoute';
import AdminRoute from '@/components/auth/AdminRoute';

import { Login } from '@/pages/Login';
import { Register } from '@/pages/Register';
import { Home } from '@/pages/Home';
import { Rooms } from '@/pages/Rooms';
import { ReserveRoom } from '@/pages/ReserveRoom';
import { MyReservations } from '@/pages/MyReservations';
import AdminRoomList from '@/pages/admin/AdminRoomList';
import AdminRoomForm from '@/pages/admin/AdminRoomForm';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>

          {/* PUBLIC */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          {/* PRIVATE */}
          <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
          <Route path="/rooms" element={<ProtectedRoute><Rooms /></ProtectedRoute>} />
          <Route path="/rooms/:roomId/reserve" element={<ProtectedRoute><ReserveRoom /></ProtectedRoute>} />
          <Route path="/my-reservations" element={<ProtectedRoute><MyReservations /></ProtectedRoute>} />

          {/* ADMIN */}
          <Route path="/admin/rooms" element={<AdminRoute><AdminRoomList /></AdminRoute>} />
          <Route path="/admin/rooms/new" element={<AdminRoute><AdminRoomForm /></AdminRoute>} />
          <Route path="/admin/rooms/:id" element={<AdminRoute><AdminRoomForm /></AdminRoute>} />

          {/* 404 fallback */}
          <Route path="*" element={<Navigate to="/login" replace />} />

        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
