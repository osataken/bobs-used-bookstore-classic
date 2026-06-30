package dto

// CartResponse represents the shopping cart
type CartResponse struct {
	Items    []CartItemResponse `json:"items"`
	SubTotal float64            `json:"subTotal"`
}

// CartItemResponse represents a cart item
type CartItemResponse struct {
	ID             int          `json:"id"`
	ShoppingCartID int          `json:"shoppingCartId"`
	BookID         int          `json:"bookId"`
	Quantity       int          `json:"quantity"`
	WantToBuy      bool         `json:"wantToBuy"`
	Book           BookResponse `json:"book"`
}

// AddToCartRequest for adding item to cart
type AddToCartRequest struct {
	BookID   int `json:"bookId" binding:"required"`
	Quantity int `json:"quantity"`
}

// AddToWishlistRequest for adding item to wishlist
type AddToWishlistRequest struct {
	BookID int `json:"bookId" binding:"required"`
}

// DeleteCartItemRequest for removing a cart item
type DeleteCartItemRequest struct {
	ShoppingCartItemID int `json:"shoppingCartItemId" binding:"required"`
}

// MoveToCartRequest for moving a wishlist item to cart
type MoveToCartRequest struct {
	ShoppingCartItemID int `json:"shoppingCartItemId" binding:"required"`
}
