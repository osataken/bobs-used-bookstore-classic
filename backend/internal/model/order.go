package model

import "time"

type Order struct {
	Entity
	CustomerID   int         `gorm:"not null" json:"customerId"`
	AddressID    int         `gorm:"not null" json:"addressId"`
	DeliveryDate time.Time   `gorm:"not null" json:"deliveryDate"`
	OrderStatus  OrderStatus `gorm:"not null;default:0" json:"orderStatus"`

	Customer   *Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Address    *Address     `gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE" json:"address,omitempty"`
	OrderItems []OrderItem  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"orderItems,omitempty"`
}

// SubTotal: sum of Book.Price per OrderItem (intentionally does NOT multiply by Quantity - source bug preserved)
func (o *Order) SubTotal() float64 {
	var total float64
	for _, item := range o.OrderItems {
		if item.Book != nil {
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
