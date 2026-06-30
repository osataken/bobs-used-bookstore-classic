package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.ShoppingCartService
}

func NewCartHandler(cartService *service.ShoppingCartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) Index(c *gin.Context) {
	correlationID := getCartCorrelationID(c)

	cart, err := h.cartService.GetShoppingCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	cartItems := cart.GetCartItems(false)
	var items []dto.CartItemResponse
	for _, item := range cartItems {
		items = append(items, dto.CartItemResponse{
			ID:             item.ID,
			ShoppingCartID: item.ShoppingCartID,
			BookID:         item.BookID,
			Quantity:       item.Quantity,
			WantToBuy:      item.WantToBuy,
			Book:           mapBookToResponse(&item.Book),
		})
	}

	c.JSON(http.StatusOK, dto.CartResponse{
		Items:    items,
		SubTotal: cart.GetSubTotal(false),
	})
}

func (h *CartHandler) Delete(c *gin.Context) {
	var req dto.DeleteCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)

	err := h.cartService.DeleteShoppingCartItem(service.DeleteShoppingCartItemDTO{
		CorrelationID:      correlationID,
		ShoppingCartItemID: req.ShoppingCartItemID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	err := h.cartService.AddToShoppingCart(service.AddToShoppingCartDTO{
		CorrelationID: correlationID,
		BookID:        req.BookID,
		Quantity:      quantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to cart"})
}

// getCartCorrelationID gets or sets the cart cookie
func getCartCorrelationID(c *gin.Context) string {
	return middleware.GetCartCorrelationID(c)
}
