package handler

import (
	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResaleHandler struct {
	offerService         *service.OfferService
	referenceDataService *service.ReferenceDataService
}

func NewResaleHandler(offerService *service.OfferService, refService *service.ReferenceDataService) *ResaleHandler {
	return &ResaleHandler{offerService: offerService, referenceDataService: refService}
}

func (h *ResaleHandler) Index(c *gin.Context) {
	customerID := c.GetInt(middleware.ContextCustomerID)
	offers, err := h.offerService.GetByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

func (h *ResaleHandler) Create(c *gin.Context) {
	var req dto.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := c.GetInt(middleware.ContextCustomerID)
	offer := &model.Offer{
		BookName:    req.BookName,
		Author:      req.Author,
		ISBN:        req.ISBN,
		BookTypeID:  req.BookTypeID,
		ConditionID: req.ConditionID,
		GenreID:     req.GenreID,
		PublisherID: req.PublisherID,
		BookPrice:   req.BookPrice,
		Summary:     req.Summary,
		CustomerID:  customerID,
	}

	if err := h.offerService.Create(offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

func (h *ResaleHandler) GetReferenceData(c *gin.Context) {
	publishers, _ := h.referenceDataService.GetByType(model.ReferenceDataTypePublisher)
	conditions, _ := h.referenceDataService.GetByType(model.ReferenceDataTypeCondition)
	bookTypes, _ := h.referenceDataService.GetByType(model.ReferenceDataTypeBookType)
	genres, _ := h.referenceDataService.GetByType(model.ReferenceDataTypeGenre)

	c.JSON(http.StatusOK, gin.H{
		"publishers": publishers,
		"conditions": conditions,
		"bookTypes":  bookTypes,
		"genres":     genres,
	})
}
