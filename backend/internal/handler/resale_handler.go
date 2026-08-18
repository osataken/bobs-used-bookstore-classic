package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ResaleHandler struct {
	offerService   *service.OfferService
	refDataService *service.ReferenceDataService
	customerService *service.CustomerService
}

func NewResaleHandler(offerService *service.OfferService, refDataService *service.ReferenceDataService, customerService *service.CustomerService) *ResaleHandler {
	return &ResaleHandler{
		offerService:   offerService,
		refDataService: refDataService,
		customerService: customerService,
	}
}

func (h *ResaleHandler) ListOffers(c *gin.Context) {
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	offers, err := h.offerService.ListByCustomer(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch offers"})
		return
	}
	c.JSON(http.StatusOK, offers)
}

func (h *ResaleHandler) CreateOffer(c *gin.Context) {
	var req struct {
		BookName    string  `json:"bookName" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		ISBN        string  `json:"isbn" binding:"required"`
		BookTypeID  int     `json:"bookTypeId" binding:"required"`
		ConditionID int     `json:"conditionId" binding:"required"`
		GenreID     int     `json:"genreId" binding:"required"`
		PublisherID int     `json:"publisherId" binding:"required"`
		BookPrice   float64 `json:"bookPrice" binding:"required"`
		Summary     string  `json:"summary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	offer := &model.Offer{
		CustomerID:  customer.ID,
		BookName:    req.BookName,
		Author:      req.Author,
		ISBN:        req.ISBN,
		BookTypeID:  req.BookTypeID,
		ConditionID: req.ConditionID,
		GenreID:     req.GenreID,
		PublisherID: req.PublisherID,
		BookPrice:   req.BookPrice,
		Summary:     req.Summary,
	}

	if err := h.offerService.CreateOffer(offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create offer"})
		return
	}
	c.JSON(http.StatusCreated, offer)
}

func (h *ResaleHandler) GetReferenceData(c *gin.Context) {
	items, err := h.refDataService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reference data"})
		return
	}

	// Group by type
	publishers := filterByType(items, model.ReferenceDataTypePublisher)
	conditions := filterByType(items, model.ReferenceDataTypeCondition)
	bookTypes := filterByType(items, model.ReferenceDataTypeBookType)
	genres := filterByType(items, model.ReferenceDataTypeGenre)

	c.JSON(http.StatusOK, gin.H{
		"publishers": publishers,
		"conditions": conditions,
		"bookTypes":  bookTypes,
		"genres":     genres,
	})
}

func filterByType(items []model.ReferenceDataItem, dataType model.ReferenceDataType) []model.ReferenceDataItem {
	var result []model.ReferenceDataItem
	for _, item := range items {
		if item.DataType == dataType {
			result = append(result, item)
		}
	}
	return result
}
