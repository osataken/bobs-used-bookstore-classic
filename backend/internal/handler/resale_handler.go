package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ResaleHandler struct {
	offerService     *service.OfferService
	refDataService   *service.ReferenceDataService
}

func NewResaleHandler(offerService *service.OfferService, refDataService *service.ReferenceDataService) *ResaleHandler {
	return &ResaleHandler{offerService: offerService, refDataService: refDataService}
}

func (h *ResaleHandler) Index(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	offers, err := h.offerService.GetOffersBySub(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get offers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

func (h *ResaleHandler) Create(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req struct {
		BookName    string  `json:"bookName" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		ISBN        string  `json:"isbn"`
		GenreID     uint    `json:"genreId" binding:"required"`
		ConditionID uint    `json:"conditionId" binding:"required"`
		PublisherID uint    `json:"publisherId" binding:"required"`
		BookTypeID  uint    `json:"bookTypeId" binding:"required"`
		Summary     string  `json:"summary"`
		BookPrice   float64 `json:"bookPrice" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.offerService.CreateOffer(sub, service.CreateOfferInput{
		BookName:    req.BookName,
		Author:      req.Author,
		ISBN:        req.ISBN,
		GenreID:     req.GenreID,
		ConditionID: req.ConditionID,
		PublisherID: req.PublisherID,
		BookTypeID:  req.BookTypeID,
		Summary:     req.Summary,
		BookPrice:   req.BookPrice,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create offer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "offer submitted"})
}

func (h *ResaleHandler) GetReferenceData(c *gin.Context) {
	items, err := h.refDataService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reference data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"referenceData": items})
}
