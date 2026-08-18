package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	orderService   *service.OrderService
	cartService    *service.ShoppingCartService
	addressService *service.AddressService
	customerService *service.CustomerService
}

func NewCheckoutHandler(orderService *service.OrderService, cartService *service.ShoppingCartService, addressService *service.AddressService, customerService *service.CustomerService) *CheckoutHandler {
	return &CheckoutHandler{
		orderService:   orderService,
		cartService:    cartService,
		addressService: addressService,
		customerService: customerService,
	}
}

func (h *CheckoutHandler) GetCheckout(c *gin.Context) {
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	addresses, err := h.addressService.ListByCustomer(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch addresses"})
		return
	}

	correlationID := getCartCorrelationID(c)
	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cart"})
		return
	}

	items := cart.GetShoppingCartItems(false)
	subTotal := cart.GetSubTotal(false)

	c.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
		"cartItems": items,
		"subTotal":  subTotal,
		"tax":       subTotal * 0.1,
		"total":     subTotal + (subTotal * 0.1),
	})
}

func (h *CheckoutHandler) SubmitOrder(c *gin.Context) {
	var req struct {
		AddressID int `json:"addressId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)
	correlationID := getCartCorrelationID(c)

	order, err := h.orderService.CreateOrder(customer.ID, req.AddressID, correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
		"orderId": order.ID,
	})
}

func (h *CheckoutHandler) Finished(c *gin.Context) {
	orderIDStr := c.Query("orderId")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
