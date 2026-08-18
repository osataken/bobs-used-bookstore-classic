package model

import "time"

// Order entity
type Order struct {
	Entity
	CustomerID   int         `gorm:"not null" json:"customerId"`
	Customer     Customer    `gorm:"foreignKey:CustomerID;constraint:OnDelete:RESTRICT" json:"customer,omitempty"`
	AddressID    int         `gorm:"not null" json:"addressId"`
	Address      Address     `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems,omitempty"`
	DeliveryDate time.Time   `json:"deliveryDate"`
	OrderStatus  OrderStatus `gorm:"default:0" json:"orderStatus"`
}

func (Order) TableName() string {
	return "Order"
}

// SubTotal - preserves the source bug: Sum(Book.Price) per OrderItem, NOT multiplied by Quantity
func (o *Order) SubTotal() float64 {
	var total float64
	for _, item := range o.OrderItems {
		if item.Book.ID != 0 {
			total += item.Book.Price
		}
	}
	return total
}

func (o *Order) Tax() float64 {
	return o.SubTotal() * 0.1
}

func (o *Order) Total() float64 {
	return o.SubTotal() + o.Tax()
}

// OrderItem entity
type OrderItem struct {
	Entity
	OrderID  int  `gorm:"not null" json:"orderId"`
	Order    Order `gorm:"foreignKey:OrderID" json:"-"`
	BookID   int  `gorm:"not null" json:"bookId"`
	Book     Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity int  `gorm:"not null" json:"quantity"`
}

func (OrderItem) TableName() string {
	return "OrderItem"
}
