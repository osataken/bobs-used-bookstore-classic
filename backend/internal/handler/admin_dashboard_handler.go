package handler

import (
	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct {
	orderService *service.OrderService
	offerService *service.OfferService
	bookService  *service.BookService
}

func NewAdminDashboardHandler(orderService *service.OrderService, offerService *service.OfferService, bookService *service.BookService) *AdminDashboardHandler {
	return &AdminDashboardHandler{
		orderService: orderService,
		offerService: offerService,
		bookService:  bookService,
	}
}

func (h *AdminDashboardHandler) Index(c *gin.Context) {
	totalOrders, _ := h.orderService.GetTotalCount()
	pendingOrders, _ := h.orderService.GetPendingCount()
	pendingOffers, _ := h.offerService.GetPendingCount()
	totalBooks, _ := h.bookService.GetTotalCount()
	lowStockBooks, _ := h.bookService.GetLowStockCount()

	c.JSON(http.StatusOK, dto.DashboardResponse{
		TotalOrders:   totalOrders,
		PendingOrders: pendingOrders,
		PendingOffers: pendingOffers,
		TotalBooks:    totalBooks,
		LowStockBooks: lowStockBooks,
	})
}
