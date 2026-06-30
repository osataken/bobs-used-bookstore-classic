package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	orderService   *service.OrderService
	cartService    *service.ShoppingCartService
	addressService *service.AddressService
}

func NewCheckoutHandler(orderService *service.OrderService, cartService *service.ShoppingCartService, addressService *service.AddressService) *CheckoutHandler {
	return &CheckoutHandler{
		orderService:   orderService,
		cartService:    cartService,
		addressService: addressService,
	}
}

func (h *CheckoutHandler) Index(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	correlationID := getCartCorrelationID(c)

	// Get cart
	cart, err := h.cartService.GetShoppingCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	// Get addresses
	addresses, err := h.addressService.GetAddresses(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get addresses"})
		return
	}

	cartItems := cart.GetCartItems(true)
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

	var addressResponses []dto.AddressResponse
	for _, a := range addresses {
		addressResponses = append(addressResponses, dto.AddressResponse{
			ID:           a.ID,
			AddressLine1: a.AddressLine1,
			AddressLine2: a.AddressLine2,
			City:         a.City,
			State:        a.State,
			Country:      a.Country,
			ZipCode:      a.ZipCode,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"cart": dto.CartResponse{
			Items:    items,
			SubTotal: cart.GetSubTotal(true),
		},
		"addresses": addressResponses,
	})
}

func (h *CheckoutHandler) Submit(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	correlationID := getCartCorrelationID(c)

	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	orderID, err := h.orderService.CreateOrder(service.CreateOrderDTO{
		CustomerSub:   sub,
		CorrelationID: correlationID,
		AddressID:     req.AddressID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orderId": orderID, "message": "order placed successfully"})
}
