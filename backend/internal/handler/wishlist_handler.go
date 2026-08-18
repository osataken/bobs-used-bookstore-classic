package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	cartService *service.ShoppingCartService
}

func NewWishlistHandler(cartService *service.ShoppingCartService) *WishlistHandler {
	return &WishlistHandler{cartService: cartService}
}

func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	correlationID := getCartCorrelationID(c)

	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
		return
	}

	items := cart.GetWishListItems()
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *WishlistHandler) AddItem(c *gin.Context) {
	var req struct {
		BookID int `json:"bookId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)
	if err := h.cartService.AddToWishlist(correlationID, req.BookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item added to wishlist"})
}

func (h *WishlistHandler) MoveToCart(c *gin.Context) {
	var req struct {
		ShoppingCartItemID int `json:"shoppingCartItemId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)
	if err := h.cartService.MoveWishlistItemToCart(correlationID, req.ShoppingCartItemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item moved to cart"})
}

func (h *WishlistHandler) MoveAllToCart(c *gin.Context) {
	correlationID := getCartCorrelationID(c)
	if err := h.cartService.MoveAllWishlistItemsToCart(correlationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move items"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all items moved to cart"})
}

func (h *WishlistHandler) DeleteItem(c *gin.Context) {
	var req struct {
		ShoppingCartItemID int `json:"shoppingCartItemId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)
	if err := h.cartService.DeleteItem(correlationID, req.ShoppingCartItemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item removed from wishlist"})
}
