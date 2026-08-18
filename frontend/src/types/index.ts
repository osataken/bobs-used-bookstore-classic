export interface Book {
  id: number;
  name: string;
  author: string;
  year?: number;
  isbn: string;
  publisherId: number;
  publisher?: ReferenceDataItem;
  bookTypeId: number;
  bookType?: ReferenceDataItem;
  genreId: number;
  genre?: ReferenceDataItem;
  conditionId: number;
  condition?: ReferenceDataItem;
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
  isAdmin: boolean;
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
  customer?: Customer;
  addressId: number;
  address?: Address;
  orderItems?: OrderItem[];
  deliveryDate: string;
  orderStatus: number;
  subTotal?: number;
  tax?: number;
  total?: number;
}

export interface OrderItem {
  id: number;
  orderId: number;
  bookId: number;
  book?: Book;
  quantity: number;
}

export interface ShoppingCartItem {
  id: number;
  shoppingCartId: number;
  bookId: number;
  book?: Book;
  quantity: number;
  wantToBuy: boolean;
}

export interface Offer {
  id: number;
  author: string;
  isbn: string;
  bookName: string;
  frontUrl: string;
  genreId: number;
  genre?: ReferenceDataItem;
  conditionId: number;
  condition?: ReferenceDataItem;
  publisherId: number;
  publisher?: ReferenceDataItem;
  bookTypeId: number;
  bookType?: ReferenceDataItem;
  summary: string;
  offerStatus: number;
  comment: string;
  customerId: number;
  customer?: Customer;
  bookPrice: number;
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

export interface BookStatistics {
  lowStock: number;
  outOfStock: number;
  stockTotal: number;
}

export interface OrderStatistics {
  pending: number;
  ordered: number;
  shipped: number;
  delivered: number;
  cancelled: number;
}

export interface OfferStatistics {
  pendingApproval: number;
  approved: number;
  received: number;
  paid: number;
  rejected: number;
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
  2: 'BookType',
  3: 'Genre',
};
