package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
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

func (h *OrderHandler) Index(c *gin.Context) {
	customerID := c.GetInt(middleware.ContextCustomerID)
	orders, err := h.orderService.GetByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type orderResponse struct {
		ID          int     `json:"id"`
		OrderStatus int     `json:"orderStatus"`
		StatusText  string  `json:"statusText"`
		SubTotal    float64 `json:"subTotal"`
		Tax         float64 `json:"tax"`
		Total       float64 `json:"total"`
		CreatedOn   string  `json:"createdOn"`
	}

	var response []orderResponse
	for _, o := range orders {
		response = append(response, orderResponse{
			ID:          o.ID,
			OrderStatus: int(o.OrderStatus),
			StatusText:  o.OrderStatus.String(),
			SubTotal:    o.SubTotal(),
			Tax:         o.Tax(),
			Total:       o.Total(),
			CreatedOn:   o.CreatedOn.Format("2006-01-02"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": response})
}

func (h *OrderHandler) Details(c *gin.Context) {
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

func (h *OrderHandler) Cancel(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.orderService.CancelOrder(body.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}
