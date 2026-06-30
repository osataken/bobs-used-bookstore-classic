package model

type ShoppingCartItem struct {
	ID             int  `gorm:"primaryKey;autoIncrement" json:"id"`
	ShoppingCartID int  `gorm:"primaryKey;not null" json:"shoppingCartId"`
	BookID         int  `gorm:"not null" json:"bookId"`
	Quantity       int  `gorm:"not null" json:"quantity"`
	WantToBuy      bool `gorm:"not null" json:"wantToBuy"`

	ShoppingCart *ShoppingCart `gorm:"foreignKey:ShoppingCartID;constraint:OnDelete:CASCADE" json:"-"`
	Book         *Book        `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE" json:"book,omitempty"`
}
