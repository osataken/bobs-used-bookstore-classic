package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// GET /api/orders
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	orders, err := h.orderService.GetOrdersByCustomer(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GET /api/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"id":           order.ID,
		"customerId":   order.CustomerID,
		"customer":     order.Customer,
		"addressId":    order.AddressID,
		"address":      order.Address,
		"orderItems":   order.OrderItems,
		"deliveryDate": order.DeliveryDate,
		"orderStatus":  order.OrderStatus,
		"subTotal":     order.SubTotal(),
		"tax":          order.Tax(),
		"total":        order.Total(),
		"createdOn":    order.CreatedOn,
	})
}

// POST /api/orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub := middleware.GetUserSub(c)
	err = h.orderService.CancelOrder(sub, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled"})
}

// POST /api/checkout
func (h *OrderHandler) Checkout(c *gin.Context) {
	var req struct {
		AddressID int `json:"addressId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub := middleware.GetUserSub(c)
	cartID := middleware.GetShoppingCartID(c)

	orderID, err := h.orderService.CreateOrder(sub, cartID, req.AddressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"orderId": orderID})
}

// GET /api/admin/orders?orderStatus=&pageIndex=1&pageSize=10
func (h *OrderHandler) AdminList(c *gin.Context) {
	var filters repository.OrderFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.orderService.GetOrders(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/admin/orders/:id
func (h *OrderHandler) AdminGetOrder(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"id":           order.ID,
		"customerId":   order.CustomerID,
		"customer":     order.Customer,
		"addressId":    order.AddressID,
		"address":      order.Address,
		"orderItems":   order.OrderItems,
		"deliveryDate": order.DeliveryDate,
		"orderStatus":  order.OrderStatus,
		"subTotal":     order.SubTotal(),
		"tax":          order.Tax(),
		"total":        order.Total(),
		"createdOn":    order.CreatedOn,
	})
}

// PUT /api/admin/orders/:id/status
func (h *OrderHandler) AdminUpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		OrderStatus model.OrderStatus `json:"orderStatus"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orderService.UpdateOrderStatus(id, req.OrderStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

// GET /api/admin/orders/statistics
func (h *OrderHandler) AdminStatistics(c *gin.Context) {
	stats, err := h.orderService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
