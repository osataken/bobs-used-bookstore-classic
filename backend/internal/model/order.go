package model

import "time"

// Order represents a customer order
type Order struct {
	Entity
	CustomerID   int         `gorm:"not null" json:"customerId"`
	AddressID    int         `gorm:"not null" json:"addressId"`
	DeliveryDate time.Time   `json:"deliveryDate"`
	OrderStatus  OrderStatus `gorm:"default:0" json:"orderStatus"`

	// Navigation properties
	Customer   Customer    `gorm:"foreignKey:CustomerID;constraint:OnDelete:NO ACTION" json:"customer,omitempty"`
	Address    Address     `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems,omitempty"`
}

// SubTotal calculates sum of book prices (intentionally does NOT multiply by quantity - preserving source bug)
func (o *Order) SubTotal() float64 {
	var total float64
	for _, item := range o.OrderItems {
		total += item.Book.Price
	}
	return total
}

// Tax calculates 10% tax on subtotal
func (o *Order) Tax() float64 {
	return o.SubTotal() * 0.1
}

// Total calculates subtotal + tax
func (o *Order) Total() float64 {
	return o.SubTotal() + o.Tax()
}

func (Order) TableName() string {
	return "Order"
}

// OrderItem represents a line item in an order
type OrderItem struct {
	Entity
	OrderID  int `gorm:"not null" json:"orderId"`
	BookID   int `gorm:"not null" json:"bookId"`
	Quantity int `gorm:"not null" json:"quantity"`

	// Navigation properties
	Book Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

func (OrderItem) TableName() string {
	return "OrderItem"
}
