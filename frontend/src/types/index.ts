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
}

export interface Order {
  id: number;
  orderStatus: number;
  deliveryDate: string;
  subTotal: number;
  tax: number;
  total: number;
  itemCount?: number;
  address?: Address;
  orderItems?: OrderItem[];
  customer?: Customer;
}

export interface OrderItem {
  id: number;
  orderId: number;
  bookId: number;
  book: Book;
  quantity: number;
}

export interface Address {
  id: number;
  addressLine1: string;
  addressLine2: string;
  city: string;
  state: string;
  country: string;
  zipCode: string;
}

export interface ShoppingCartItem {
  id: number;
  shoppingCartId: number;
  bookId: number;
  book: Book;
  quantity: number;
  wantToBuy: boolean;
}

export interface Offer {
  id: number;
  bookName: string;
  author: string;
  isbn: string;
  offerStatus: number;
  bookPrice: number;
  genre: ReferenceDataItem;
  condition: ReferenceDataItem;
  publisher: ReferenceDataItem;
  bookType: ReferenceDataItem;
}

export interface ReferenceDataItem {
  id: number;
  dataType: number;
  text: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  totalCount: number;
  pageIndex: number;
  pageSize: number;
  totalPages: number;
}

export interface AuthUser {
  authenticated: boolean;
  sub?: string;
  username?: string;
  firstName?: string;
  lastName?: string;
  roles?: string[];
}

export const OrderStatusLabels: Record<number, string> = {
  0: 'Pending',
  1: 'Ordered',
  2: 'Shipped',
  3: 'Delivered',
  4: 'Cancelled',
};

export const OfferStatusLabels: Record<number, string> = {
  0: 'Pending Approval',
  1: 'Approved',
  2: 'Received',
  3: 'Paid',
  4: 'Rejected',
};

export const ReferenceDataTypeLabels: Record<number, string> = {
  0: 'Publisher',
  1: 'Condition',
  2: 'Book Type',
  3: 'Genre',
};
