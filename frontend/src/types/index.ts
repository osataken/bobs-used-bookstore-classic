export interface Book {
  id: number
  name: string
  author: string
  year: number | null
  isbn: string
  publisherId: number
  publisher: string
  bookTypeId: number
  bookType: string
  genreId: number
  genre: string
  conditionId: number
  condition: string
  coverImageUrl: string
  summary: string
  price: number
  quantity: number
  isInStock: boolean
  isLowInStock: boolean
}

export interface Order {
  id: number
  customerId: number
  customerName: string
  addressId: number
  deliveryDate: string
  orderStatus: number
  statusText: string
  subTotal: number
  tax: number
  total: number
  orderItems: OrderItem[]
  address?: Address
  createdOn: string
}

export interface OrderItem {
  id: number
  bookId: number
  quantity: number
  book: Book
}

export interface CartItem {
  id: number
  shoppingCartId: number
  bookId: number
  quantity: number
  wantToBuy: boolean
  book: Book
}

export interface Cart {
  items: CartItem[]
  subTotal: number
}

export interface Address {
  id: number
  addressLine1: string
  addressLine2: string
  city: string
  state: string
  country: string
  zipCode: string
}

export interface Offer {
  id: number
  bookName: string
  author: string
  isbn: string
  genreId: number
  genre: string
  conditionId: number
  condition: string
  publisherId: number
  publisher: string
  bookTypeId: number
  bookType: string
  summary: string
  offerStatus: number
  statusText: string
  comment: string
  customerId: number
  bookPrice: number
}

export interface ReferenceData {
  id: number
  dataType: number
  typeText: string
  text: string
}

export interface AuthUser {
  sub: string
  username: string
  firstName: string
  lastName: string
  isAdmin: boolean
}

export interface PaginatedResponse<T> {
  items: T[]
  pageIndex: number
  totalPages: number
  totalCount: number
  hasNext: boolean
  hasPrevious: boolean
}
