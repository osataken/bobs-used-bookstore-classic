package model

type ShoppingCart struct {
	Entity
	CorrelationID     string             `gorm:"column:correlation_id" json:"correlationId"`
	ShoppingCartItems []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID" json:"shoppingCartItems,omitempty"`
}

func (ShoppingCart) TableName() string {
	return "shopping_cart"
}

type ShoppingCartItem struct {
	ID             int          `gorm:"primaryKey;autoIncrement" json:"id"`
	ShoppingCartID int          `gorm:"primaryKey;column:shopping_cart_id" json:"shoppingCartId"`
	ShoppingCart   ShoppingCart `gorm:"foreignKey:ShoppingCartID" json:"-"`
	BookID         int          `gorm:"column:book_id" json:"bookId"`
	Book           Book         `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity       int          `gorm:"column:quantity" json:"quantity"`
	WantToBuy      bool         `gorm:"column:want_to_buy" json:"wantToBuy"`
}

func (ShoppingCartItem) TableName() string {
	return "shopping_cart_item"
}
