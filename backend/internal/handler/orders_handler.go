package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	orderService    *service.OrderService
	customerService *service.CustomerService
}

func NewOrdersHandler(orderService *service.OrderService, customerService *service.CustomerService) *OrdersHandler {
	return &OrdersHandler{orderService: orderService, customerService: customerService}
}

func (h *OrdersHandler) ListOrders(c *gin.Context) {
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	orders, err := h.orderService.ListOrdersByCustomer(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	// Build response with computed totals
	type orderResponse struct {
		model.Order
		SubTotal float64 `json:"subTotal"`
		Tax      float64 `json:"tax"`
		Total    float64 `json:"total"`
	}

	var result []orderResponse
	for _, o := range orders {
		result = append(result, orderResponse{
			Order:    o,
			SubTotal: o.SubTotal(),
			Tax:      o.Tax(),
			Total:    o.Total(),
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *OrdersHandler) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Verify ownership
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)
	if order.CustomerID != customer.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order":    order,
		"subTotal": order.SubTotal(),
		"tax":      order.Tax(),
		"total":    order.Total(),
	})
}

func (h *OrdersHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	if err := h.orderService.CancelOrder(id, customer.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}
