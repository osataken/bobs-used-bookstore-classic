export interface Book {
  id: number;
  name: string;
  author: string;
  year: number | null;
  isbn: string;
  publisherId: number;
  bookTypeId: number;
  genreId: number;
  conditionId: number;
  coverImageUrl: string;
  summary: string;
  price: number;
  quantity: number;
  publisher?: ReferenceDataItem;
  bookType?: ReferenceDataItem;
  genre?: ReferenceDataItem;
  condition?: ReferenceDataItem;
}

export interface Customer {
  id: number;
  sub: string;
  username: string;
  firstName: string;
  lastName: string;
  email: string;
}

export interface Address {
  id: number;
  addressLine1: string;
  addressLine2: string;
  city: string;
  state: string;
  country: string;
  zipCode: string;
  customerId: number;
  isActive: boolean;
}

export interface Order {
  id: number;
  customerId: number;
  addressId: number;
  deliveryDate: string;
  orderStatus: number;
  statusText?: string;
  subTotal?: number;
  tax?: number;
  total?: number;
  createdOn: string;
  customer?: Customer;
  address?: Address;
  orderItems?: OrderItem[];
}

export interface OrderItem {
  id: number;
  orderId: number;
  bookId: number;
  quantity: number;
  book?: Book;
}

export interface Offer {
  id: number;
  author: string;
  isbn: string;
  bookName: string;
  frontUrl: string;
  genreId: number;
  conditionId: number;
  publisherId: number;
  bookTypeId: number;
  summary: string;
  offerStatus: number;
  comment: string;
  customerId: number;
  bookPrice: number;
}

export interface ShoppingCartItem {
  id: number;
  shoppingCartId: number;
  bookId: number;
  quantity: number;
  wantToBuy: boolean;
  book?: Book;
}

export interface ReferenceDataItem {
  id: number;
  dataType: number;
  text: string;
}

export interface PaginatedList<T> {
  items: T[];
  pageIndex: number;
  pageSize: number;
  totalCount: number;
  totalPages: number;
}

export interface DashboardData {
  totalOrders: number;
  pendingOrders: number;
  pendingOffers: number;
  totalBooks: number;
  lowStockBooks: number;
}

export interface AuthUser {
  authenticated: boolean;
  sub?: string;
  username?: string;
  firstName?: string;
  lastName?: string;
  role?: string;
  customerID?: number;
}

export enum OrderStatus {
  Pending = 0,
  Ordered = 1,
  Shipped = 2,
  Delivered = 3,
  Cancelled = 4,
}

export enum OfferStatus {
  PendingApproval = 0,
  Approved = 1,
  Received = 2,
  Paid = 3,
  Rejected = 4,
}

export enum ReferenceDataType {
  Publisher = 0,
  Condition = 1,
  BookType = 2,
  Genre = 3,
}
