package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
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

	orders, err := h.orderService.GetOrdersByCustomer(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}

	var response []dto.OrderResponse
	for _, o := range orders {
		response = append(response, mapOrderToResponse(&o))
	}

	c.JSON(http.StatusOK, gin.H{"orders": response})
}

func (h *OrderHandler) Details(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.orderService.GetOrderByIDAndSub(id, sub)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, mapOrderToResponse(order))
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.orderService.CancelOrder(service.CancelOrderDTO{
		CustomerSub: sub,
		OrderID:     id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

func mapOrderToResponse(o *model.Order) dto.OrderResponse {
	var items []dto.OrderItemResponse
	for _, item := range o.OrderItems {
		items = append(items, dto.OrderItemResponse{
			ID:       item.ID,
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Book:     mapBookToResponse(&item.Book),
		})
	}

	var address *dto.AddressResponse
	if o.AddressID != 0 {
		address = &dto.AddressResponse{
			ID:           o.Address.ID,
			AddressLine1: o.Address.AddressLine1,
			AddressLine2: o.Address.AddressLine2,
			City:         o.Address.City,
			State:        o.Address.State,
			Country:      o.Address.Country,
			ZipCode:      o.Address.ZipCode,
		}
	}

	return dto.OrderResponse{
		ID:           o.ID,
		CustomerID:   o.CustomerID,
		CustomerName: o.Customer.FullName(),
		AddressID:    o.AddressID,
		DeliveryDate: o.DeliveryDate,
		OrderStatus:  int(o.OrderStatus),
		StatusText:   o.OrderStatus.String(),
		SubTotal:     o.SubTotal(),
		Tax:          o.Tax(),
		Total:        o.Total(),
		OrderItems:   items,
		Address:      address,
		CreatedOn:    o.CreatedOn,
	}
}
