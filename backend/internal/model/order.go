package model

import "time"

type Order struct {
	Entity
	CustomerID  uint        `json:"customerId" gorm:"not null"`
	Customer    Customer    `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	AddressID   uint        `json:"addressId" gorm:"not null"`
	Address     Address     `json:"address" gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE"`
	OrderItems  []OrderItem `json:"orderItems" gorm:"foreignKey:OrderID"`
	DeliveryDate time.Time  `json:"deliveryDate"`
	OrderStatus OrderStatus `json:"orderStatus" gorm:"default:0"`
}

func NewOrder(customerID uint, addressID uint) *Order {
	return &Order{
		CustomerID:   customerID,
		AddressID:    addressID,
		DeliveryDate: time.Now().AddDate(0, 0, 7),
		OrderStatus:  OrderStatusPending,
	}
}

// SubTotal preserves the source bug: Sum(Book.Price) per OrderItem, does NOT multiply by Quantity
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
	OrderID  uint `json:"orderId" gorm:"not null"`
	BookID   uint `json:"bookId" gorm:"not null"`
	Book     Book `json:"book" gorm:"foreignKey:BookID"`
	Quantity int  `json:"quantity" gorm:"not null"`
}
