import api from './api';
import { Book, PaginatedList, Order, Address, Offer, ReferenceDataItem, ShoppingCartItem, Customer } from '../types';

// Auth
export const authService = {
  login: () => api.get('/auth/login'),
  logout: () => api.get('/auth/logout'),
  me: () => api.get<Customer>('/auth/me'),
};

// Home
export const homeService = {
  getFeatured: () => api.get<Book[]>('/home/featured'),
};

// Search
export const searchService = {
  search: (params: { searchString?: string; sortBy?: string; pageIndex?: number; pageSize?: number }) =>
    api.get<PaginatedList<Book>>('/search', { params }),
  getBook: (id: number) => api.get<Book>(`/search/${id}`),
};

// Cart
export const cartService = {
  getCart: () => api.get<{ items: ShoppingCartItem[]; subTotal: number }>('/cart'),
  addItem: (bookId: number, quantity: number = 1) => api.post('/cart/add', { bookId, quantity }),
  deleteItem: (shoppingCartItemId: number) => api.post('/cart/delete', { shoppingCartItemId }),
};

// Wishlist
export const wishlistService = {
  getWishlist: () => api.get<{ items: ShoppingCartItem[] }>('/wishlist'),
  addItem: (bookId: number) => api.post('/wishlist/add', { bookId }),
  moveToCart: (shoppingCartItemId: number) => api.post('/wishlist/move', { shoppingCartItemId }),
  moveAllToCart: () => api.post('/wishlist/move-all'),
  deleteItem: (shoppingCartItemId: number) => api.post('/wishlist/delete', { shoppingCartItemId }),
};

// Checkout
export const checkoutService = {
  getCheckout: () => api.get('/checkout'),
  submitOrder: (addressId: number) => api.post('/checkout', { addressId }),
};

// Orders
export const ordersService = {
  list: () => api.get<Order[]>('/orders'),
  get: (id: number) => api.get(`/orders/${id}`),
  cancel: (id: number) => api.post(`/orders/${id}/cancel`),
};

// Addresses
export const addressService = {
  list: () => api.get<Address[]>('/addresses'),
  create: (data: Partial<Address>) => api.post<Address>('/addresses', data),
  update: (id: number, data: Partial<Address>) => api.put<Address>(`/addresses/${id}`, data),
  delete: (id: number) => api.delete(`/addresses/${id}`),
};

// Resale
export const resaleService = {
  list: () => api.get<Offer[]>('/resale'),
  create: (data: Partial<Offer>) => api.post<Offer>('/resale', data),
  getReferenceData: () => api.get<{
    publishers: ReferenceDataItem[];
    conditions: ReferenceDataItem[];
    bookTypes: ReferenceDataItem[];
    genres: ReferenceDataItem[];
  }>('/resale/reference-data'),
};

// Admin
export const adminService = {
  getDashboard: () => api.get('/admin/dashboard'),
  listBooks: (params: Record<string, string | number>) => api.get<PaginatedList<Book>>('/admin/inventory', { params }),
  getBook: (id: number) => api.get(`/admin/inventory/${id}`),
  createBook: (data: FormData) => api.post('/admin/inventory', data, { headers: { 'Content-Type': 'multipart/form-data' } }),
  updateBook: (id: number, data: FormData) => api.put(`/admin/inventory/${id}`, data, { headers: { 'Content-Type': 'multipart/form-data' } }),
  listOffers: (params: Record<string, string | number>) => api.get<PaginatedList<Offer>>('/admin/offers', { params }),
  approveOffer: (id: number) => api.post(`/admin/offers/${id}/approve`),
  rejectOffer: (id: number) => api.post(`/admin/offers/${id}/reject`),
  receivedOffer: (id: number) => api.post(`/admin/offers/${id}/received`),
  paidOffer: (id: number) => api.post(`/admin/offers/${id}/paid`),
  listOrders: (params: Record<string, string | number>) => api.get<PaginatedList<Order>>('/admin/orders', { params }),
  getOrder: (id: number) => api.get(`/admin/orders/${id}`),
  updateOrderStatus: (id: number, orderStatus: number) => api.put(`/admin/orders/${id}/status`, { orderStatus }),
  listReferenceData: (params: Record<string, string | number>) => api.get<PaginatedList<ReferenceDataItem>>('/admin/reference-data', { params }),
  createReferenceData: (data: { dataType: number; text: string }) => api.post('/admin/reference-data', data),
  updateReferenceData: (id: number, data: { text: string }) => api.put(`/admin/reference-data/${id}`, data),
};
