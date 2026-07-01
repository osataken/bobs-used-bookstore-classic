package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) Index(c *gin.Context) {
	cartID := middleware.GetCartID(c)
	cart, err := h.cartService.GetCart(cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}
	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "subTotal": 0})
		return
	}

	items := cart.GetCartItems(false)
	subTotal := cart.SubTotal(false)

	c.JSON(http.StatusOK, gin.H{
		"items":    items,
		"subTotal": subTotal,
	})
}

type AddToCartRequest struct {
	BookID   uint `json:"bookId" binding:"required"`
	Quantity int  `json:"quantity"`
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	cartID := middleware.GetCartID(c)

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try query param fallback
		bookIDStr := c.Query("bookId")
		if bookIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bookId is required"})
			return
		}
		bookID, _ := strconv.ParseUint(bookIDStr, 10, 32)
		req.BookID = uint(bookID)
		req.Quantity = 1
	}

	if req.Quantity < 1 {
		req.Quantity = 1
	}

	if err := h.cartService.AddToCart(cartID, req.BookID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to cart"})
}

func (h *CartHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
