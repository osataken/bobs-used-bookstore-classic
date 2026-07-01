package model

type ShoppingCart struct {
	Entity
	CorrelationID     string             `json:"correlationId" gorm:"uniqueIndex;not null"`
	ShoppingCartItems []ShoppingCartItem `json:"shoppingCartItems" gorm:"foreignKey:ShoppingCartID"`
}

type ShoppingCartItem struct {
	ID             uint `json:"id" gorm:"primaryKey;autoIncrement"`
	ShoppingCartID uint `json:"shoppingCartId" gorm:"primaryKey"`
	BookID         uint `json:"bookId" gorm:"not null"`
	Book           Book `json:"book" gorm:"foreignKey:BookID"`
	Quantity       int  `json:"quantity" gorm:"not null;default:1"`
	WantToBuy      bool `json:"wantToBuy" gorm:"not null"`
}

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

func (sc *ShoppingCart) GetWishlistItems() []ShoppingCartItem {
	var items []ShoppingCartItem
	for _, item := range sc.ShoppingCartItems {
		if !item.WantToBuy {
			items = append(items, item)
		}
	}
	return items
}

func (sc *ShoppingCart) AddItemToCart(bookID uint, quantity int) {
	for i, item := range sc.ShoppingCartItems {
		if item.BookID == bookID && item.WantToBuy {
			sc.ShoppingCartItems[i].Quantity += quantity
			return
		}
	}
	sc.ShoppingCartItems = append(sc.ShoppingCartItems, ShoppingCartItem{
		BookID:    bookID,
		Quantity:  quantity,
		WantToBuy: true,
	})
}

func (sc *ShoppingCart) AddItemToWishlist(bookID uint) {
	for _, item := range sc.ShoppingCartItems {
		if item.BookID == bookID && !item.WantToBuy {
			return
		}
	}
	sc.ShoppingCartItems = append(sc.ShoppingCartItems, ShoppingCartItem{
		BookID:    bookID,
		Quantity:  1,
		WantToBuy: false,
	})
}

func (sc *ShoppingCart) MoveWishlistItemToCart(itemID uint) {
	for i, item := range sc.ShoppingCartItems {
		if item.ID == itemID && !item.WantToBuy {
			sc.ShoppingCartItems[i].WantToBuy = true
			return
		}
	}
}

func (sc *ShoppingCart) RemoveItemByID(itemID uint) {
	for i, item := range sc.ShoppingCartItems {
		if item.ID == itemID {
			sc.ShoppingCartItems = append(sc.ShoppingCartItems[:i], sc.ShoppingCartItems[i+1:]...)
			return
		}
	}
}

// SubTotal preserves source bug: Sum(Book.Price), does NOT multiply by Quantity
func (sc *ShoppingCart) SubTotal(excludeOutOfStock bool) float64 {
	var total float64
	for _, item := range sc.GetCartItems(excludeOutOfStock) {
		total += item.Book.Price
	}
	return total
}
