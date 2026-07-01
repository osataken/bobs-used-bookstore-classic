package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Index(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	orders, err := h.orderService.GetOrdersByCustomerSub(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}

	type orderResponse struct {
		ID           uint    `json:"id"`
		OrderStatus  int     `json:"orderStatus"`
		DeliveryDate string  `json:"deliveryDate"`
		SubTotal     float64 `json:"subTotal"`
		Tax          float64 `json:"tax"`
		Total        float64 `json:"total"`
		ItemCount    int     `json:"itemCount"`
	}

	var response []orderResponse
	for _, o := range orders {
		response = append(response, orderResponse{
			ID:           o.ID,
			OrderStatus:  int(o.OrderStatus),
			DeliveryDate: o.DeliveryDate.Format("2006-01-02"),
			SubTotal:     o.SubTotal(),
			Tax:          o.Tax(),
			Total:        o.Total(),
			ItemCount:    len(o.OrderItems),
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": response})
}

func (h *OrderHandler) Details(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderService.GetOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           order.ID,
		"orderStatus":  order.OrderStatus,
		"deliveryDate": order.DeliveryDate,
		"subTotal":     order.SubTotal(),
		"tax":          order.Tax(),
		"total":        order.Total(),
		"address":      order.Address,
		"orderItems":   order.OrderItems,
	})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req struct {
		OrderID uint `json:"orderId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orderId is required"})
		return
	}

	if err := h.orderService.CancelOrder(req.OrderID, sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}
