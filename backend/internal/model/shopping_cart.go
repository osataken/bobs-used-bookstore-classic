package model

import "time"

// ShoppingCart entity
type ShoppingCart struct {
	Entity
	CorrelationID     string             `gorm:"not null;uniqueIndex" json:"correlationId"`
	ShoppingCartItems []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID" json:"shoppingCartItems,omitempty"`
}

func (ShoppingCart) TableName() string {
	return "ShoppingCart"
}

func (sc *ShoppingCart) GetShoppingCartItems(includeOutOfStock bool) []ShoppingCartItem {
	var items []ShoppingCartItem
	for _, item := range sc.ShoppingCartItems {
		if item.WantToBuy {
			if includeOutOfStock || item.Book.Quantity > 0 {
				items = append(items, item)
			}
		}
	}
	return items
}

func (sc *ShoppingCart) GetWishListItems() []ShoppingCartItem {
	var items []ShoppingCartItem
	for _, item := range sc.ShoppingCartItems {
		if !item.WantToBuy {
			items = append(items, item)
		}
	}
	return items
}

func (sc *ShoppingCart) GetSubTotal(includeOutOfStock bool) float64 {
	var total float64
	for _, item := range sc.GetShoppingCartItems(includeOutOfStock) {
		total += item.Book.Price
	}
	return total
}

// ShoppingCartItem entity - composite PK (ID, ShoppingCartID)
type ShoppingCartItem struct {
	ID             int          `gorm:"primaryKey;autoIncrement" json:"id"`
	ShoppingCartID int          `gorm:"primaryKey" json:"shoppingCartId"`
	ShoppingCart   ShoppingCart `gorm:"foreignKey:ShoppingCartID" json:"-"`
	BookID         int          `gorm:"not null" json:"bookId"`
	Book           Book         `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity       int          `gorm:"not null;default:1" json:"quantity"`
	WantToBuy      bool         `gorm:"not null" json:"wantToBuy"`
	CreatedBy      string       `gorm:"default:System" json:"createdBy"`
	CreatedOn      time.Time    `json:"createdOn"`
	UpdatedOn      time.Time    `json:"updatedOn"`
	Version        int          `gorm:"default:1" json:"-"`
}

func (ShoppingCartItem) TableName() string {
	return "ShoppingCartItem"
}
