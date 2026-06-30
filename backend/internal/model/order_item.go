package model

type OrderItem struct {
	Entity
	OrderID  int  `gorm:"not null" json:"orderId"`
	BookID   int  `gorm:"not null" json:"bookId"`
	Quantity int  `gorm:"not null" json:"quantity"`

	Book *Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE" json:"book,omitempty"`
}
