package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/service"

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

func (h *AdminDashboardHandler) GetDashboard(c *gin.Context) {
	orderStats, err := h.orderService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order statistics"})
		return
	}

	offerStats, err := h.offerService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get offer statistics"})
		return
	}

	bookStats, err := h.bookService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderStatistics": orderStats,
		"offerStatistics": offerStats,
		"bookStatistics":  bookStats,
	})
}
