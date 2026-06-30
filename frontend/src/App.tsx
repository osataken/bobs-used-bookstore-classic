import { Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
import ProtectedRoute from './components/ProtectedRoute'
import HomePage from './pages/HomePage'
import SearchPage from './pages/SearchPage'
import BookDetailsPage from './pages/BookDetailsPage'
import CartPage from './pages/CartPage'
import WishlistPage from './pages/WishlistPage'
import CheckoutPage from './pages/CheckoutPage'
import OrdersPage from './pages/OrdersPage'
import OrderDetailsPage from './pages/OrderDetailsPage'
import AddressPage from './pages/AddressPage'
import ResalePage from './pages/ResalePage'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/admin/DashboardPage'
import InventoryPage from './pages/admin/InventoryPage'
import AdminOrdersPage from './pages/admin/AdminOrdersPage'
import OffersPage from './pages/admin/OffersPage'
import ReferenceDataPage from './pages/admin/ReferenceDataPage'

function App() {
  return (
    <div>
      <Navbar />
      <main style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/search/:id" element={<BookDetailsPage />} />
          <Route path="/cart" element={<CartPage />} />
          <Route path="/wishlist" element={<WishlistPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/checkout" element={<ProtectedRoute><CheckoutPage /></ProtectedRoute>} />
          <Route path="/orders" element={<ProtectedRoute><OrdersPage /></ProtectedRoute>} />
          <Route path="/orders/:id" element={<ProtectedRoute><OrderDetailsPage /></ProtectedRoute>} />
          <Route path="/addresses" element={<ProtectedRoute><AddressPage /></ProtectedRoute>} />
          <Route path="/resale" element={<ProtectedRoute><ResalePage /></ProtectedRoute>} />
          <Route path="/admin/dashboard" element={<ProtectedRoute admin><DashboardPage /></ProtectedRoute>} />
          <Route path="/admin/inventory" element={<ProtectedRoute admin><InventoryPage /></ProtectedRoute>} />
          <Route path="/admin/orders" element={<ProtectedRoute admin><AdminOrdersPage /></ProtectedRoute>} />
          <Route path="/admin/offers" element={<ProtectedRoute admin><OffersPage /></ProtectedRoute>} />
          <Route path="/admin/reference-data" element={<ProtectedRoute admin><ReferenceDataPage /></ProtectedRoute>} />
        </Routes>
      </main>
    </div>
  )
}

export default App
