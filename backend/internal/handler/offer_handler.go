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

type OfferHandler struct {
	offerService *service.OfferService
}

func NewOfferHandler(offerService *service.OfferService) *OfferHandler {
	return &OfferHandler{offerService: offerService}
}

// GET /api/offers
func (h *OfferHandler) GetMyOffers(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	offers, err := h.offerService.GetOffersByCustomer(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offers)
}

// POST /api/offers
func (h *OfferHandler) Create(c *gin.Context) {
	var req struct {
		BookName    string  `json:"bookName" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		ISBN        string  `json:"isbn"`
		BookTypeID  int     `json:"bookTypeId" binding:"required"`
		ConditionID int     `json:"conditionId" binding:"required"`
		GenreID     int     `json:"genreId" binding:"required"`
		PublisherID int     `json:"publisherId" binding:"required"`
		BookPrice   float64 `json:"bookPrice" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub := middleware.GetUserSub(c)
	dto := service.CreateOfferDTO{
		BookName:    req.BookName,
		Author:      req.Author,
		ISBN:        req.ISBN,
		BookTypeID:  req.BookTypeID,
		ConditionID: req.ConditionID,
		GenreID:     req.GenreID,
		PublisherID: req.PublisherID,
		BookPrice:   req.BookPrice,
	}

	err := h.offerService.CreateOffer(sub, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Offer submitted"})
}

// GET /api/admin/offers?offerStatus=&pageIndex=1&pageSize=10
func (h *OfferHandler) AdminList(c *gin.Context) {
	var filters repository.OfferFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.offerService.GetOffers(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /api/admin/offers/:id/approve
func (h *OfferHandler) AdminApprove(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.offerService.UpdateOfferStatus(id, model.OfferStatusApproved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Offer approved"})
}

// POST /api/admin/offers/:id/reject
func (h *OfferHandler) AdminReject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.offerService.UpdateOfferStatus(id, model.OfferStatusRejected)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Offer rejected"})
}

// POST /api/admin/offers/:id/received
func (h *OfferHandler) AdminReceived(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.offerService.UpdateOfferStatus(id, model.OfferStatusReceived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Offer marked as received"})
}

// POST /api/admin/offers/:id/paid
func (h *OfferHandler) AdminPaid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.offerService.UpdateOfferStatus(id, model.OfferStatusPaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Offer marked as paid"})
}

// GET /api/admin/offers/statistics
func (h *OfferHandler) AdminStatistics(c *gin.Context) {
	stats, err := h.offerService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
