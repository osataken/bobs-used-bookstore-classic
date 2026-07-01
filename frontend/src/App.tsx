import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import HomePage from './pages/HomePage';
import SearchPage from './pages/SearchPage';
import BookDetailsPage from './pages/BookDetailsPage';
import CartPage from './pages/CartPage';
import WishlistPage from './pages/WishlistPage';
import CheckoutPage from './pages/CheckoutPage';
import OrdersPage from './pages/OrdersPage';
import OrderDetailsPage from './pages/OrderDetailsPage';
import ResalePage from './pages/ResalePage';
import AddressPage from './pages/AddressPage';
import DashboardPage from './pages/admin/DashboardPage';
import InventoryPage from './pages/admin/InventoryPage';
import AdminOrdersPage from './pages/admin/AdminOrdersPage';
import OffersPage from './pages/admin/OffersPage';
import ReferenceDataPage from './pages/admin/ReferenceDataPage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="search" element={<SearchPage />} />
        <Route path="search/:id" element={<BookDetailsPage />} />
        <Route path="cart" element={<CartPage />} />
        <Route path="wishlist" element={<WishlistPage />} />
        <Route path="checkout" element={<CheckoutPage />} />
        <Route path="orders" element={<OrdersPage />} />
        <Route path="orders/:id" element={<OrderDetailsPage />} />
        <Route path="resale" element={<ResalePage />} />
        <Route path="addresses" element={<AddressPage />} />
        <Route path="admin" element={<DashboardPage />} />
        <Route path="admin/inventory" element={<InventoryPage />} />
        <Route path="admin/orders" element={<AdminOrdersPage />} />
        <Route path="admin/offers" element={<OffersPage />} />
        <Route path="admin/reference-data" element={<ReferenceDataPage />} />
      </Route>
    </Routes>
  );
}

export default App;
