package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	cartService *service.CartService
}

func NewWishlistHandler(cartService *service.CartService) *WishlistHandler {
	return &WishlistHandler{cartService: cartService}
}

func (h *WishlistHandler) Index(c *gin.Context) {
	cartID := middleware.GetCartID(c)
	cart, err := h.cartService.GetCart(cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get wishlist"})
		return
	}
	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
		return
	}

	items := cart.GetWishlistItems()
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	cartID := middleware.GetCartID(c)

	var req struct {
		BookID uint `json:"bookId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bookId is required"})
		return
	}

	if err := h.cartService.AddToWishlist(cartID, req.BookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to wishlist"})
}

func (h *WishlistHandler) MoveToCart(c *gin.Context) {
	cartID := middleware.GetCartID(c)

	var req struct {
		ShoppingCartItemID uint `json:"shoppingCartItemId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shoppingCartItemId is required"})
		return
	}

	if err := h.cartService.MoveWishlistItemToCart(cartID, req.ShoppingCartItemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item moved to cart"})
}

func (h *WishlistHandler) MoveAllToCart(c *gin.Context) {
	cartID := middleware.GetCartID(c)

	if err := h.cartService.MoveAllWishlistItemsToCart(cartID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all items moved to cart"})
}

func (h *WishlistHandler) Delete(c *gin.Context) {
	cartID := middleware.GetCartID(c)

	var req struct {
		ShoppingCartItemID uint `json:"shoppingCartItemId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shoppingCartItemId is required"})
		return
	}

	if err := h.cartService.DeleteItem(cartID, req.ShoppingCartItemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from wishlist"})
}
