export interface Book {
  id: number;
  name: string;
  author: string;
  year: number | null;
  isbn: string;
  publisherId: number;
  publisher: ReferenceDataItem;
  bookTypeId: number;
  bookType: ReferenceDataItem;
  genreId: number;
  genre: ReferenceDataItem;
  conditionId: number;
  condition: ReferenceDataItem;
  coverImageUrl: string;
  summary: string;
  price: number;
  quantity: number;
}

export interface Customer {
  id: number;
  sub: string;
  username: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
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

export enum OrderStatus {
  Pending = 0,
  Ordered = 1,
  Shipped = 2,
  Delivered = 3,
  Cancelled = 4,
}

export interface Order {
  id: number;
  customerId: number;
  customer: Customer;
  addressId: number;
  address: Address;
  orderItems: OrderItem[];
  deliveryDate: string;
  orderStatus: OrderStatus;
  subTotal: number;
  tax: number;
  total: number;
  createdOn: string;
}

export interface OrderItem {
  id: number;
  orderId: number;
  bookId: number;
  book: Book;
  quantity: number;
}

export interface CartItem {
  id: number;
  bookId: number;
  book: Book;
  quantity: number;
  inStock: boolean;
}

export interface CartResponse {
  items: CartItem[];
  subTotal: number;
}

export interface WishlistResponse {
  items: CartItem[];
}

export enum OfferStatus {
  PendingApproval = 0,
  Approved = 1,
  Received = 2,
  Paid = 3,
  Rejected = 4,
}

export interface Offer {
  id: number;
  bookName: string;
  author: string;
  isbn: string;
  genreId: number;
  genre: ReferenceDataItem;
  conditionId: number;
  condition: ReferenceDataItem;
  publisherId: number;
  publisher: ReferenceDataItem;
  bookTypeId: number;
  bookType: ReferenceDataItem;
  offerStatus: OfferStatus;
  bookPrice: number;
  customerId: number;
}

export enum ReferenceDataType {
  Publisher = 0,
  Condition = 1,
  BookType = 2,
  Genre = 3,
}

export interface ReferenceDataItem {
  id: number;
  dataType: ReferenceDataType;
  text: string;
}

export interface PaginatedList<T> {
  items: T[];
  totalPages: number;
  pageIndex: number;
  pageSize: number;
  totalCount: number;
}

export interface AuthStatus {
  authenticated: boolean;
  sub?: string;
  role?: string;
  username?: string;
}

export interface DashboardStats {
  orders: {
    pastDueOrders: number;
    pendingOrders: number;
    ordersThisMonth: number;
    ordersTotal: number;
  };
  offers: {
    pendingOffers: number;
    offersThisMonth: number;
    offersTotal: number;
  };
  inventory: {
    lowStock: number;
    outOfStock: number;
    stockTotal: number;
  };
}
