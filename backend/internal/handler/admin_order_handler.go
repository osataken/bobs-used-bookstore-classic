package handler

import (
	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	orderService *service.OrderService
}

func NewAdminOrderHandler(orderService *service.OrderService) *AdminOrderHandler {
	return &AdminOrderHandler{orderService: orderService}
}

func (h *AdminOrderHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	dateFrom := c.Query("filters.OrderDateFromFilter")
	dateTo := c.Query("filters.OrderDateToFilter")

	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var statusFilter *model.OrderStatus
	if s := c.Query("filters.OrderStatusFilter"); s != "" {
		val, err := strconv.Atoi(s)
		if err == nil {
			status := model.OrderStatus(val)
			statusFilter = &status
		}
	}

	result, err := h.orderService.GetAll(pageIndex, pageSize, statusFilter, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminOrderHandler) Details(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.orderService.GetByID(id)
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

func (h *AdminOrderHandler) UpdateStatus(c *gin.Context) {
	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.UpdateStatus(req.OrderID, model.OrderStatus(req.OrderStatus)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
