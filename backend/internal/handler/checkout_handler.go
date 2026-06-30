package handler

import (
	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"
	"net/http"

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
	customerID := c.GetInt(middleware.ContextCustomerID)
	correlationID := c.GetString("cartCorrelationID")

	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	addresses, err := h.addressService.GetByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cartItems []interface{}
	var subTotal float64
	for _, item := range cart.Items {
		if item.WantToBuy && item.Book != nil && item.Book.IsInStock() {
			cartItems = append(cartItems, item)
			subTotal += item.Book.Price
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     cartItems,
		"addresses": addresses,
		"subTotal":  subTotal,
		"tax":       subTotal * 0.1,
		"total":     subTotal + subTotal*0.1,
	})
}

func (h *CheckoutHandler) PlaceOrder(c *gin.Context) {
	var req dto.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := c.GetInt(middleware.ContextCustomerID)
	correlationID := c.GetString("cartCorrelationID")

	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(customerID, req.SelectedAddressID, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order placed",
		"orderId": order.ID,
	})
}
