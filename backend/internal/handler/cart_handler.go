package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.ShoppingCartService
}

func NewCartHandler(cartService *service.ShoppingCartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// GET /api/cart
func (h *CartHandler) GetCart(c *gin.Context) {
	cartID := middleware.GetShoppingCartID(c)
	cart, err := h.cartService.GetCart(cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "subTotal": 0})
		return
	}

	// Filter cart items (WantToBuy=true)
	var cartItems []interface{}
	var subTotal float64
	for _, item := range cart.ShoppingCartItems {
		if item.WantToBuy {
			cartItems = append(cartItems, gin.H{
				"id":       item.ID,
				"bookId":   item.BookID,
				"book":     item.Book,
				"quantity": item.Quantity,
				"inStock":  item.Book.Quantity > 0,
			})
			subTotal += item.Book.Price
		}
	}
	if cartItems == nil {
		cartItems = []interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    cartItems,
		"subTotal": subTotal,
	})
}

// DELETE /api/cart/:itemId
func (h *CartHandler) DeleteItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	cartID := middleware.GetShoppingCartID(c)
	err = h.cartService.DeleteItem(cartID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from shopping cart"})
}

// GET /api/wishlist
func (h *CartHandler) GetWishlist(c *gin.Context) {
	cartID := middleware.GetShoppingCartID(c)
	cart, err := h.cartService.GetCart(cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
		return
	}

	var wishlistItems []interface{}
	for _, item := range cart.ShoppingCartItems {
		if !item.WantToBuy {
			wishlistItems = append(wishlistItems, gin.H{
				"id":       item.ID,
				"bookId":   item.BookID,
				"book":     item.Book,
				"quantity": item.Quantity,
			})
		}
	}
	if wishlistItems == nil {
		wishlistItems = []interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"items": wishlistItems})
}

// POST /api/wishlist/:itemId/move-to-cart
func (h *CartHandler) MoveToCart(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	cartID := middleware.GetShoppingCartID(c)
	err = h.cartService.MoveToCart(cartID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item moved to shopping cart"})
}

// POST /api/wishlist/move-all-to-cart
func (h *CartHandler) MoveAllToCart(c *gin.Context) {
	cartID := middleware.GetShoppingCartID(c)
	err := h.cartService.MoveAllToCart(cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All items moved to shopping cart"})
}

// DELETE /api/wishlist/:itemId
func (h *CartHandler) DeleteWishlistItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	cartID := middleware.GetShoppingCartID(c)
	err = h.cartService.DeleteItem(cartID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from wishlist"})
}
