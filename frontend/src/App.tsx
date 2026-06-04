import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import Layout from './components/Layout';
import HomePage from './pages/HomePage';
import SearchPage from './pages/SearchPage';
import BookDetailsPage from './pages/BookDetailsPage';
import CartPage from './pages/CartPage';
import WishlistPage from './pages/WishlistPage';
import CheckoutPage from './pages/CheckoutPage';
import OrdersPage from './pages/OrdersPage';
import OrderDetailsPage from './pages/OrderDetailsPage';
import AddressesPage from './pages/AddressesPage';
import AddressFormPage from './pages/AddressFormPage';
import ResalePage from './pages/ResalePage';
import ResaleCreatePage from './pages/ResaleCreatePage';
import AdminDashboard from './pages/admin/AdminDashboard';
import AdminInventory from './pages/admin/AdminInventory';
import AdminInventoryForm from './pages/admin/AdminInventoryForm';
import AdminOrders from './pages/admin/AdminOrders';
import AdminOrderDetails from './pages/admin/AdminOrderDetails';
import AdminOffers from './pages/admin/AdminOffers';
import AdminReferenceData from './pages/admin/AdminReferenceData';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { auth, loading } = useAuth();
  if (loading) return <div>Loading...</div>;
  if (!auth.authenticated) return <Navigate to="/" />;
  return <>{children}</>;
}

function AdminRoute({ children }: { children: React.ReactNode }) {
  const { auth, loading } = useAuth();
  if (loading) return <div>Loading...</div>;
  if (!auth.authenticated || auth.role !== 'Administrators') return <Navigate to="/" />;
  return <>{children}</>;
}

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/search" element={<SearchPage />} />
        <Route path="/books/:id" element={<BookDetailsPage />} />
        <Route path="/cart" element={<CartPage />} />
        <Route path="/wishlist" element={<WishlistPage />} />
        <Route path="/checkout" element={<ProtectedRoute><CheckoutPage /></ProtectedRoute>} />
        <Route path="/orders" element={<ProtectedRoute><OrdersPage /></ProtectedRoute>} />
        <Route path="/orders/:id" element={<ProtectedRoute><OrderDetailsPage /></ProtectedRoute>} />
        <Route path="/addresses" element={<ProtectedRoute><AddressesPage /></ProtectedRoute>} />
        <Route path="/addresses/new" element={<ProtectedRoute><AddressFormPage /></ProtectedRoute>} />
        <Route path="/addresses/:id/edit" element={<ProtectedRoute><AddressFormPage /></ProtectedRoute>} />
        <Route path="/resale" element={<ProtectedRoute><ResalePage /></ProtectedRoute>} />
        <Route path="/resale/new" element={<ProtectedRoute><ResaleCreatePage /></ProtectedRoute>} />
        <Route path="/admin" element={<AdminRoute><AdminDashboard /></AdminRoute>} />
        <Route path="/admin/inventory" element={<AdminRoute><AdminInventory /></AdminRoute>} />
        <Route path="/admin/inventory/new" element={<AdminRoute><AdminInventoryForm /></AdminRoute>} />
        <Route path="/admin/inventory/:id/edit" element={<AdminRoute><AdminInventoryForm /></AdminRoute>} />
        <Route path="/admin/orders" element={<AdminRoute><AdminOrders /></AdminRoute>} />
        <Route path="/admin/orders/:id" element={<AdminRoute><AdminOrderDetails /></AdminRoute>} />
        <Route path="/admin/offers" element={<AdminRoute><AdminOffers /></AdminRoute>} />
        <Route path="/admin/reference-data" element={<AdminRoute><AdminReferenceData /></AdminRoute>} />
      </Routes>
    </Layout>
  );
}
