import axios from 'axios';
import type { Book, Address, Order, Offer, ShoppingCartItem, ReferenceDataItem, PaginatedList, DashboardData, AuthUser } from '../types';

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
});

// Auth
export const login = () => api.post<AuthUser>('/auth/login');
export const logout = () => api.post('/auth/logout');
export const getMe = () => api.get<AuthUser>('/auth/me');

// Home
export const getHome = () => api.get<{ bestSellers: Book[] }>('/home');

// Search
export const searchBooks = (params: { searchString?: string; sortBy?: string; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedList<Book>>('/search', { params });
export const getBookDetails = (id: number) => api.get<Book>(`/search/${id}`);
export const addToCart = (bookId: number) => api.post('/search/add-to-cart', null, { params: { bookId } });
export const addToWishlist = (bookId: number) => api.post('/search/add-to-wishlist', null, { params: { bookId } });

// Cart
export const getCart = () => api.get<{ items: ShoppingCartItem[]; subTotal: number }>('/cart');
export const deleteCartItem = (shoppingCartItemId: number) => api.post('/cart/delete', { shoppingCartItemId });

// Wishlist
export const getWishlist = () => api.get<{ items: ShoppingCartItem[] }>('/wishlist');
export const moveToCart = (shoppingCartItemId: number) => api.post('/wishlist/move-to-cart', { shoppingCartItemId });
export const moveAllToCart = () => api.post('/wishlist/move-all-to-cart');
export const deleteWishlistItem = (shoppingCartItemId: number) => api.post('/wishlist/delete', { shoppingCartItemId });

// Checkout
export const getCheckout = () => api.get<{ items: ShoppingCartItem[]; addresses: Address[]; subTotal: number; tax: number; total: number }>('/checkout');
export const placeOrder = (selectedAddressId: number) => api.post<{ message: string; orderId: number }>('/checkout', { selectedAddressId });

// Orders
export const getOrders = () => api.get<{ orders: Order[] }>('/orders');
export const getOrderDetails = (id: number) => api.get<{ order: Order; subTotal: number; tax: number; total: number }>(`/orders/${id}`);
export const cancelOrder = (id: number) => api.post('/orders/cancel', { id });

// Addresses
export const getAddresses = () => api.get<{ addresses: Address[] }>('/addresses');
export const createAddress = (data: Omit<Address, 'id' | 'customerId' | 'isActive'>) => api.post<Address>('/addresses', data);
export const updateAddress = (data: { id: number } & Omit<Address, 'customerId' | 'isActive'>) => api.put<Address>('/addresses', data);
export const deleteAddress = (id: number) => api.post('/addresses/delete', { id });

// Resale
export const getOffers = () => api.get<{ offers: Offer[] }>('/resale');
export const createOffer = (data: Partial<Offer>) => api.post<Offer>('/resale', data);

// Reference Data
export const getReferenceData = () => api.get<{ publishers: ReferenceDataItem[]; conditions: ReferenceDataItem[]; bookTypes: ReferenceDataItem[]; genres: ReferenceDataItem[] }>('/reference-data');

// Admin
export const getAdminDashboard = () => api.get<DashboardData>('/admin/dashboard');
export const getAdminInventory = (params: Record<string, string | number | boolean>) => api.get<PaginatedList<Book>>('/admin/inventory', { params });
export const getAdminBookDetails = (id: number) => api.get<Book>(`/admin/inventory/${id}`);
export const createBook = (data: Partial<Book>) => api.post<Book>('/admin/inventory', data);
export const updateBook = (data: Partial<Book>) => api.put<Book>('/admin/inventory', data);
export const uploadBookImage = (id: number, file: File) => {
  const formData = new FormData();
  formData.append('coverImage', file);
  return api.post(`/admin/inventory/${id}/image`, formData);
};

export const getAdminOrders = (params: Record<string, string | number>) => api.get<PaginatedList<Order>>('/admin/orders', { params });
export const getAdminOrderDetails = (id: number) => api.get<{ order: Order; subTotal: number; tax: number; total: number }>(`/admin/orders/${id}`);
export const updateOrderStatus = (orderId: number, orderStatus: number) => api.post('/admin/orders/update-status', { orderId, orderStatus });

export const getAdminOffers = (params: Record<string, string | number>) => api.get<PaginatedList<Offer>>('/admin/offers', { params });
export const approveOffer = (id: number) => api.post('/admin/offers/approve', { id });
export const rejectOffer = (id: number) => api.post('/admin/offers/reject', { id });
export const markOfferReceived = (id: number) => api.post('/admin/offers/received', { id });
export const markOfferPaid = (id: number) => api.post('/admin/offers/paid', { id });

export const getAdminReferenceData = (params: Record<string, string | number>) => api.get<PaginatedList<ReferenceDataItem>>('/admin/reference-data', { params });
export const createReferenceData = (data: { dataType: number; text: string }) => api.post<ReferenceDataItem>('/admin/reference-data', data);
export const updateReferenceData = (data: { id: number; dataType: number; text: string }) => api.put<ReferenceDataItem>('/admin/reference-data', data);

export default api;
