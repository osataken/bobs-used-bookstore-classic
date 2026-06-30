package handler

import (
	"net/http"
	"strconv"
	"time"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

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

	var filtersReq dto.OrderFiltersRequest
	_ = c.ShouldBindQuery(&filtersReq)

	filters := repository.OrderFilters{}
	if filtersReq.OrderStatus != nil {
		status := model.OrderStatus(*filtersReq.OrderStatus)
		filters.OrderStatus = &status
	}
	if filtersReq.DateFrom != nil {
		t, err := time.Parse("2006-01-02", *filtersReq.DateFrom)
		if err == nil {
			filters.DateFrom = &t
		}
	}
	if filtersReq.DateTo != nil {
		t, err := time.Parse("2006-01-02", *filtersReq.DateTo)
		if err == nil {
			filters.DateTo = &t
		}
	}

	result, err := h.orderService.GetOrders(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}

	var orders []dto.OrderResponse
	for _, o := range result.Items {
		orders = append(orders, mapOrderToResponse(&o))
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items:       orders,
		PageIndex:   result.PageIndex,
		TotalPages:  result.TotalPages,
		TotalCount:  result.TotalCount,
		HasNext:     result.HasNextPage(),
		HasPrevious: result.HasPreviousPage(),
	})
}

func (h *AdminOrderHandler) Details(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, mapOrderToResponse(order))
}

func (h *AdminOrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.orderService.UpdateOrderStatus(service.UpdateOrderStatusDTO{
		OrderID:     id,
		OrderStatus: model.OrderStatus(req.OrderStatus),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}
