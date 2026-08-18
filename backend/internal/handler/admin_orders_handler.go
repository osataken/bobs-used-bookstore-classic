package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminOrdersHandler struct {
	orderService *service.OrderService
}

func NewAdminOrdersHandler(orderService *service.OrderService) *AdminOrdersHandler {
	return &AdminOrdersHandler{orderService: orderService}
}

func (h *AdminOrdersHandler) ListOrders(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if v, err := strconv.Atoi(c.Query("orderStatus")); err == nil {
		filters["orderStatus"] = v
	}

	result, err := h.orderService.ListAllOrders(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminOrdersHandler) GetOrder(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"order":    order,
		"subTotal": order.SubTotal(),
		"tax":      order.Tax(),
		"total":    order.Total(),
	})
}

func (h *AdminOrdersHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req struct {
		OrderStatus int `json:"orderStatus" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.orderService.UpdateOrderStatus(id, model.OrderStatus(req.OrderStatus)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}
