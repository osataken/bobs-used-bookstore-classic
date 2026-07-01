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

func (h *AdminDashboardHandler) Index(c *gin.Context) {
	orderStats, err := h.orderService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order stats"})
		return
	}

	offerStats, err := h.offerService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get offer stats"})
		return
	}

	bookStats, err := h.bookService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": gin.H{
			"pastDueOrders":   orderStats.PastDueOrders,
			"pendingOrders":   orderStats.PendingOrders,
			"ordersThisMonth": orderStats.OrdersThisMonth,
			"ordersTotal":     orderStats.OrdersTotal,
		},
		"offers": gin.H{
			"pendingOffers":   offerStats.PendingOffers,
			"offersThisMonth": offerStats.OffersThisMonth,
			"offersTotal":     offerStats.OffersTotal,
		},
		"inventory": gin.H{
			"lowStock":   bookStats.LowStock,
			"outOfStock": bookStats.OutOfStock,
			"stockTotal": bookStats.StockTotal,
		},
	})
}
