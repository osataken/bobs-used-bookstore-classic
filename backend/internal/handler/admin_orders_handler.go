package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminOrdersHandler struct {
	orderService *service.OrderService
}

func NewAdminOrdersHandler(orderService *service.OrderService) *AdminOrdersHandler {
	return &AdminOrdersHandler{orderService: orderService}
}

func (h *AdminOrdersHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var filters repository.OrderFilters
	if v := c.Query("status"); v != "" {
		status, _ := strconv.Atoi(v)
		s := model.OrderStatus(status)
		filters.Status = &s
	}

	result, err := h.orderService.GetOrdersFiltered(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items:      result.Items,
		TotalCount: result.TotalCount,
		PageIndex:  result.PageIndex,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *AdminOrdersHandler) Details(c *gin.Context) {
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
		"customer":     order.Customer,
		"address":      order.Address,
		"orderItems":   order.OrderItems,
	})
}

func (h *AdminOrdersHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		OrderID     uint `json:"orderId" binding:"required"`
		OrderStatus int  `json:"orderStatus" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.UpdateOrderStatus(req.OrderID, model.OrderStatus(req.OrderStatus)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}
