package model

import "time"

type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusOrdered   OrderStatus = 1
	OrderStatusShipped   OrderStatus = 2
	OrderStatusDelivered OrderStatus = 3
	OrderStatusCancelled OrderStatus = 4
)

type Order struct {
	Entity
	CustomerID   int         `gorm:"column:customer_id" json:"customerId"`
	Customer     Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	AddressID    int         `gorm:"column:address_id" json:"addressId"`
	Address      Address     `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems,omitempty"`
	DeliveryDate time.Time   `gorm:"column:delivery_date" json:"deliveryDate"`
	OrderStatus  OrderStatus `gorm:"column:order_status;default:0" json:"orderStatus"`
}

func (Order) TableName() string {
	return "order"
}

// SubTotal sums Book.Price per OrderItem (intentionally does NOT multiply by Quantity)
func (o *Order) SubTotal() float64 {
	var total float64
	for _, item := range o.OrderItems {
		total += item.Book.Price
	}
	return total
}

func (o *Order) Tax() float64 {
	return o.SubTotal() * 0.1
}

func (o *Order) Total() float64 {
	return o.SubTotal() + o.Tax()
}

type OrderItem struct {
	Entity
	OrderID  int  `gorm:"column:order_id" json:"orderId"`
	Order    Order `gorm:"foreignKey:OrderID" json:"-"`
	BookID   int  `gorm:"column:book_id" json:"bookId"`
	Book     Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity int  `gorm:"column:quantity" json:"quantity"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
