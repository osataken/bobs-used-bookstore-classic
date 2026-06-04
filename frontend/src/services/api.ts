import axios from 'axios';
import type {
  Book, Address, Order, Offer, ReferenceDataItem,
  CartResponse, WishlistResponse, PaginatedList, AuthStatus, DashboardStats,
  OrderStatus
} from '../types';

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
});

// Auth
export const getAuthStatus = () => api.get<AuthStatus>('/auth/status').then(r => r.data);
export const login = () => api.get('/auth/login').then(r => r.data);
export const logout = () => api.get('/auth/logout').then(r => r.data);

// Books
export const searchBooks = (params: { searchString?: string; sortBy?: string; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedList<Book>>('/books/search', { params }).then(r => r.data);
export const getBook = (id: number) => api.get<Book>(`/books/${id}`).then(r => r.data);
export const getBestSellingBooks = (count = 4) => api.get<Book[]>('/books/best-selling', { params: { count } }).then(r => r.data);
export const addToCart = (bookId: number) => api.post(`/books/${bookId}/add-to-cart`).then(r => r.data);
export const addToWishlist = (bookId: number) => api.post(`/books/${bookId}/add-to-wishlist`).then(r => r.data);

// Cart
export const getCart = () => api.get<CartResponse>('/cart').then(r => r.data);
export const deleteCartItem = (itemId: number) => api.delete(`/cart/${itemId}`).then(r => r.data);

// Wishlist
export const getWishlist = () => api.get<WishlistResponse>('/wishlist').then(r => r.data);
export const moveToCart = (itemId: number) => api.post(`/wishlist/${itemId}/move-to-cart`).then(r => r.data);
export const moveAllToCart = () => api.post('/wishlist/move-all-to-cart').then(r => r.data);
export const deleteWishlistItem = (itemId: number) => api.delete(`/wishlist/${itemId}`).then(r => r.data);

// Orders
export const getMyOrders = () => api.get<Order[]>('/orders').then(r => r.data);
export const getOrder = (id: number) => api.get<Order>(`/orders/${id}`).then(r => r.data);
export const cancelOrder = (id: number) => api.post(`/orders/${id}/cancel`).then(r => r.data);
export const checkout = (addressId: number) => api.post<{ orderId: number }>('/checkout', { addressId }).then(r => r.data);

// Addresses
export const getAddresses = () => api.get<Address[]>('/addresses').then(r => r.data);
export const getAddress = (id: number) => api.get<Address>(`/addresses/${id}`).then(r => r.data);
export const createAddress = (data: Omit<Address, 'id' | 'customerId' | 'isActive'>) => api.post('/addresses', data).then(r => r.data);
export const updateAddress = (id: number, data: Omit<Address, 'id' | 'customerId' | 'isActive'>) => api.put(`/addresses/${id}`, data).then(r => r.data);
export const deleteAddress = (id: number) => api.delete(`/addresses/${id}`).then(r => r.data);

// Offers
export const getMyOffers = () => api.get<Offer[]>('/offers').then(r => r.data);
export const createOffer = (data: { bookName: string; author: string; isbn: string; bookTypeId: number; conditionId: number; genreId: number; publisherId: number; bookPrice: number }) =>
  api.post('/offers', data).then(r => r.data);

// Reference Data
export const getReferenceData = () => api.get<ReferenceDataItem[]>('/reference-data').then(r => r.data);

// Admin
export const getAdminDashboard = () => api.get<DashboardStats>('/admin/dashboard').then(r => r.data);
export const getAdminInventory = (params: { searchString?: string; genreId?: number; bookTypeId?: number; conditionId?: number; publisherId?: number; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedList<Book>>('/admin/inventory', { params }).then(r => r.data);
export const createBook = (data: object) => api.post('/admin/inventory', data).then(r => r.data);
export const updateBook = (id: number, data: object) => api.put(`/admin/inventory/${id}`, data).then(r => r.data);
export const getAdminOrders = (params: { orderStatus?: number; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedList<Order>>('/admin/orders', { params }).then(r => r.data);
export const getAdminOrder = (id: number) => api.get<Order>(`/admin/orders/${id}`).then(r => r.data);
export const updateOrderStatus = (id: number, orderStatus: OrderStatus) => api.put(`/admin/orders/${id}/status`, { orderStatus }).then(r => r.data);
export const getAdminOffers = (params: { offerStatus?: number; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedList<Offer>>('/admin/offers', { params }).then(r => r.data);
export const approveOffer = (id: number) => api.post(`/admin/offers/${id}/approve`).then(r => r.data);
export const rejectOffer = (id: number) => api.post(`/admin/offers/${id}/reject`).then(r => r.data);
export const receivedOffer = (id: number) => api.post(`/admin/offers/${id}/received`).then(r => r.data);
export const paidOffer = (id: number) => api.post(`/admin/offers/${id}/paid`).then(r => r.data);
export const getAdminReferenceData = (params: { dataType?: number; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedList<ReferenceDataItem>>('/admin/reference-data', { params }).then(r => r.data);
export const createReferenceData = (data: { dataType: number; text: string }) => api.post('/admin/reference-data', data).then(r => r.data);
export const updateReferenceData = (id: number, data: { dataType: number; text: string }) => api.put(`/admin/reference-data/${id}`, data).then(r => r.data);
