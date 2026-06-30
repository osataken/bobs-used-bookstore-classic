package model

// ShoppingCart represents a user's shopping cart
type ShoppingCart struct {
	Entity
	CorrelationID     string             `gorm:"type:varchar(450);uniqueIndex;not null" json:"correlationId"`
	ShoppingCartItems []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID" json:"shoppingCartItems,omitempty"`
}

func (ShoppingCart) TableName() string {
	return "ShoppingCart"
}

// GetCartItems returns items where WantToBuy is true
func (sc *ShoppingCart) GetCartItems(excludeOutOfStock bool) []ShoppingCartItem {
	var items []ShoppingCartItem
	for _, item := range sc.ShoppingCartItems {
		if item.WantToBuy {
			if excludeOutOfStock && !item.Book.IsInStock() {
				continue
			}
			items = append(items, item)
		}
	}
	return items
}

// GetWishlistItems returns items where WantToBuy is false
func (sc *ShoppingCart) GetWishlistItems() []ShoppingCartItem {
	var items []ShoppingCartItem
	for _, item := range sc.ShoppingCartItems {
		if !item.WantToBuy {
			items = append(items, item)
		}
	}
	return items
}

// GetSubTotal returns sum of book prices for cart items
func (sc *ShoppingCart) GetSubTotal(excludeOutOfStock bool) float64 {
	var total float64
	for _, item := range sc.GetCartItems(excludeOutOfStock) {
		total += item.Book.Price
	}
	return total
}

// ShoppingCartItem represents an item in a shopping cart
type ShoppingCartItem struct {
	ID             int  `gorm:"primaryKey;autoIncrement" json:"id"`
	ShoppingCartID int  `gorm:"primaryKey" json:"shoppingCartId"`
	BookID         int  `gorm:"not null" json:"bookId"`
	Quantity       int  `gorm:"not null" json:"quantity"`
	WantToBuy      bool `gorm:"not null" json:"wantToBuy"`

	// Navigation properties
	Book         Book         `gorm:"foreignKey:BookID" json:"book,omitempty"`
	ShoppingCart ShoppingCart `gorm:"foreignKey:ShoppingCartID" json:"-"`
}

func (ShoppingCartItem) TableName() string {
	return "ShoppingCartItem"
}
