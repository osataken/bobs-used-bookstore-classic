import axios from 'axios'
import type { Book, Cart, Order, Address, Offer, ReferenceData, AuthUser, PaginatedResponse } from '../types'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

// Auth
export const login = (redirectUri?: string) =>
  api.get<{ redirect: string }>('/auth/login', { params: { redirectUri } })

export const logout = () =>
  api.get<{ redirect: string }>('/auth/logout')

export const getMe = () =>
  api.get<AuthUser>('/auth/me')

// Home
export const getHome = () =>
  api.get<{ bestSellingBooks: Book[] }>('/home')

// Search
export const searchBooks = (params: { searchString?: string; sortBy?: string; pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedResponse<Book>>('/search', { params })

export const getBookDetails = (id: number) =>
  api.get<Book>(`/search/${id}`)

export const addToCart = (bookId: number, quantity = 1) =>
  api.post('/search/add-to-cart', null, { params: { bookId, quantity } })

export const addToWishlist = (bookId: number) =>
  api.post('/search/add-to-wishlist', null, { params: { bookId } })

// Cart
export const getCart = () =>
  api.get<Cart>('/cart')

export const addItemToCart = (bookId: number, quantity = 1) =>
  api.post('/cart/add', { bookId, quantity })

export const deleteCartItem = (shoppingCartItemId: number) =>
  api.post('/cart/delete', { shoppingCartItemId })

// Wishlist
export const getWishlist = () =>
  api.get<{ items: import('../types').CartItem[] }>('/wishlist')

export const moveToCart = (shoppingCartItemId: number) =>
  api.post('/wishlist/move-to-cart', { shoppingCartItemId })

export const moveAllToCart = () =>
  api.post('/wishlist/move-all-to-cart')

export const deleteWishlistItem = (shoppingCartItemId: number) =>
  api.post('/wishlist/delete', { shoppingCartItemId })

// Checkout
export const getCheckout = () =>
  api.get<{ cart: Cart; addresses: Address[] }>('/checkout')

export const submitCheckout = (addressId: number) =>
  api.post<{ orderId: number; message: string }>('/checkout', { addressId })

// Orders
export const getOrders = () =>
  api.get<{ orders: Order[] }>('/orders')

export const getOrderDetails = (id: number) =>
  api.get<Order>(`/orders/${id}`)

export const cancelOrder = (id: number) =>
  api.post(`/orders/${id}/cancel`)

// Addresses
export const getAddresses = () =>
  api.get<{ addresses: Address[] }>('/addresses')

export const getAddress = (id: number) =>
  api.get<Address>(`/addresses/${id}`)

export const createAddress = (data: Omit<Address, 'id'>) =>
  api.post('/addresses', data)

export const updateAddress = (id: number, data: Omit<Address, 'id'>) =>
  api.put(`/addresses/${id}`, data)

export const deleteAddress = (id: number) =>
  api.delete(`/addresses/${id}`)

// Resale
export const getResaleOffers = () =>
  api.get<{ offers: Offer[] }>('/resale')

export const createResaleOffer = (data: {
  bookName: string; author: string; isbn: string;
  bookTypeId: number; conditionId: number; genreId: number;
  publisherId: number; bookPrice: number; summary?: string
}) => api.post('/resale', data)

// Reference Data
export const getReferenceData = () =>
  api.get<{ items: ReferenceData[] }>('/reference-data')

// Admin
export const getAdminDashboard = () =>
  api.get('/admin/dashboard')

export const getAdminInventory = (params: { pageIndex?: number; pageSize?: number; name?: string; author?: string }) =>
  api.get<PaginatedResponse<Book>>('/admin/inventory', { params })

export const getAdminBookDetails = (id: number) =>
  api.get<Book>(`/admin/inventory/${id}`)

export const createBook = (data: FormData) =>
  api.post('/admin/inventory', data)

export const updateBook = (data: FormData) =>
  api.put('/admin/inventory', data)

export const getAdminOrders = (params: { pageIndex?: number; pageSize?: number; orderStatus?: number }) =>
  api.get<PaginatedResponse<Order>>('/admin/orders', { params })

export const getAdminOrderDetails = (id: number) =>
  api.get<Order>(`/admin/orders/${id}`)

export const updateOrderStatus = (id: number, orderStatus: number) =>
  api.post(`/admin/orders/${id}/status`, { orderStatus })

export const getAdminOffers = (params: { pageIndex?: number; pageSize?: number }) =>
  api.get<PaginatedResponse<Offer>>('/admin/offers', { params })

export const approveOffer = (id: number) =>
  api.post(`/admin/offers/${id}/approve`)

export const rejectOffer = (id: number) =>
  api.post(`/admin/offers/${id}/reject`)

export const receivedOffer = (id: number) =>
  api.post(`/admin/offers/${id}/received`)

export const paidOffer = (id: number) =>
  api.post(`/admin/offers/${id}/paid`)

export const getAdminReferenceData = (params: { pageIndex?: number; pageSize?: number; dataType?: number }) =>
  api.get<PaginatedResponse<ReferenceData>>('/admin/reference-data', { params })

export const createReferenceData = (data: { dataType: number; text: string }) =>
  api.post('/admin/reference-data', data)

export const updateReferenceData = (id: number, data: { dataType: number; text: string }) =>
  api.put(`/admin/reference-data/${id}`, data)

export default api
