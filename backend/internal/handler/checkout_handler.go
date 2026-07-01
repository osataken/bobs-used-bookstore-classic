package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	cartService    *service.CartService
	orderService   *service.OrderService
	addressService *service.AddressService
}

func NewCheckoutHandler(cartService *service.CartService, orderService *service.OrderService, addressService *service.AddressService) *CheckoutHandler {
	return &CheckoutHandler{
		cartService:    cartService,
		orderService:   orderService,
		addressService: addressService,
	}
}

func (h *CheckoutHandler) Index(c *gin.Context) {
	cartID := middleware.GetCartID(c)
	sub := middleware.GetUserSub(c)

	cart, err := h.cartService.GetCart(cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	addresses, err := h.addressService.GetAddresses(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get addresses"})
		return
	}

	var items interface{}
	var subTotal float64
	if cart != nil {
		items = cart.GetCartItems(true)
		subTotal = cart.SubTotal(true)
	} else {
		items = []interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"subTotal":  subTotal,
		"tax":       subTotal * 0.1,
		"total":     subTotal + (subTotal * 0.1),
		"addresses": addresses,
	})
}

func (h *CheckoutHandler) CreateOrder(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	cartID := middleware.GetCartID(c)

	var req struct {
		SelectedAddressID uint `json:"selectedAddressId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "selectedAddressId is required"})
		return
	}

	orderID, err := h.orderService.CreateOrder(service.CreateOrderInput{
		CustomerSub:   sub,
		CorrelationID: cartID,
		AddressID:     req.SelectedAddressID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"orderId": orderID})
}

func (h *CheckoutHandler) Finished(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Query("orderId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orderId"})
		return
	}

	order, err := h.orderService.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":      order.ID,
		"orderStatus":  order.OrderStatus,
		"deliveryDate": order.DeliveryDate,
		"subTotal":     order.SubTotal(),
		"tax":          order.Tax(),
		"total":        order.Total(),
	})
}
