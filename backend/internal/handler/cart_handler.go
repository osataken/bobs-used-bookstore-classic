package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShoppingCartHandler struct {
	cartService *service.ShoppingCartService
}

func NewShoppingCartHandler(cartService *service.ShoppingCartService) *ShoppingCartHandler {
	return &ShoppingCartHandler{cartService: cartService}
}

func getCartCorrelationID(c *gin.Context) string {
	// Check for existing cookie
	correlationID, err := c.Cookie("ShoppingCartId")
	if err == nil && correlationID != "" {
		return correlationID
	}
	// Generate new UUID
	correlationID = uuid.New().String()
	c.SetCookie("ShoppingCartId", correlationID, 86400*30, "/", "", false, false)
	return correlationID
}

func (h *ShoppingCartHandler) GetCart(c *gin.Context) {
	correlationID := getCartCorrelationID(c)

	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "subTotal": 0})
		return
	}

	items := cart.GetShoppingCartItems(true)
	subTotal := cart.GetSubTotal(false)

	c.JSON(http.StatusOK, gin.H{
		"items":    items,
		"subTotal": subTotal,
	})
}

func (h *ShoppingCartHandler) AddItem(c *gin.Context) {
	var req struct {
		BookID   int `json:"bookId" binding:"required"`
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	correlationID := getCartCorrelationID(c)
	if err := h.cartService.AddToCart(correlationID, req.BookID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item added to cart"})
}

func (h *ShoppingCartHandler) DeleteItem(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
