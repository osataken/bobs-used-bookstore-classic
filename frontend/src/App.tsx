import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Home from './pages/Home'
import Search from './pages/Search'
import BookDetail from './pages/BookDetail'
import ShoppingCart from './pages/ShoppingCart'
import Wishlist from './pages/Wishlist'
import Checkout from './pages/Checkout'
import Orders from './pages/Orders'
import OrderDetail from './pages/OrderDetail'
import Addresses from './pages/Addresses'
import Resale from './pages/Resale'
import AdminDashboard from './pages/admin/Dashboard'
import AdminInventory from './pages/admin/Inventory'
import AdminOrders from './pages/admin/Orders'
import AdminOffers from './pages/admin/Offers'
import AdminReferenceData from './pages/admin/ReferenceData'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="search" element={<Search />} />
        <Route path="books/:id" element={<BookDetail />} />
        <Route path="cart" element={<ShoppingCart />} />
        <Route path="wishlist" element={<Wishlist />} />
        <Route path="checkout" element={<Checkout />} />
        <Route path="orders" element={<Orders />} />
        <Route path="orders/:id" element={<OrderDetail />} />
        <Route path="addresses" element={<Addresses />} />
        <Route path="resale" element={<Resale />} />
        <Route path="admin" element={<AdminDashboard />} />
        <Route path="admin/inventory" element={<AdminInventory />} />
        <Route path="admin/orders" element={<AdminOrders />} />
        <Route path="admin/offers" element={<AdminOffers />} />
        <Route path="admin/reference-data" element={<AdminReferenceData />} />
      </Route>
    </Routes>
  )
}

export default App
