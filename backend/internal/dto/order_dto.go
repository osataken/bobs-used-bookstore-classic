package dto

import "time"

// OrderResponse represents an order in API responses
type OrderResponse struct {
	ID           int               `json:"id"`
	CustomerID   int               `json:"customerId"`
	CustomerName string            `json:"customerName"`
	AddressID    int               `json:"addressId"`
	DeliveryDate time.Time         `json:"deliveryDate"`
	OrderStatus  int               `json:"orderStatus"`
	StatusText   string            `json:"statusText"`
	SubTotal     float64           `json:"subTotal"`
	Tax          float64           `json:"tax"`
	Total        float64           `json:"total"`
	OrderItems   []OrderItemResponse `json:"orderItems"`
	Address      *AddressResponse  `json:"address,omitempty"`
	CreatedOn    time.Time         `json:"createdOn"`
}

// OrderItemResponse represents an order item in API responses
type OrderItemResponse struct {
	ID       int          `json:"id"`
	BookID   int          `json:"bookId"`
	Quantity int          `json:"quantity"`
	Book     BookResponse `json:"book"`
}

// CheckoutRequest for creating an order
type CheckoutRequest struct {
	AddressID int `json:"addressId" binding:"required"`
}

// UpdateOrderStatusRequest for admin order status update
type UpdateOrderStatusRequest struct {
	OrderStatus int `json:"orderStatus" binding:"required"`
}

// OrderFiltersRequest for filtering orders
type OrderFiltersRequest struct {
	OrderStatus *int    `form:"orderStatus"`
	DateFrom    *string `form:"dateFrom"`
	DateTo      *string `form:"dateTo"`
}

// OrderStatisticsResponse for dashboard stats
type OrderStatisticsResponse struct {
	PendingOrders   int `json:"pendingOrders"`
	PastDueOrders   int `json:"pastDueOrders"`
	OrdersThisMonth int `json:"ordersThisMonth"`
	OrdersTotal     int `json:"ordersTotal"`
}
