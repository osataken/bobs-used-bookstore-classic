package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
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

	offers, err := h.offerService.GetOffersByCustomer(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get offers"})
		return
	}

	var response []dto.OfferResponse
	for _, o := range offers {
		response = append(response, mapOfferToResponse(&o))
	}

	c.JSON(http.StatusOK, gin.H{"offers": response})
}

func (h *ResaleHandler) Create(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req dto.OfferCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.offerService.Create(service.CreateOfferDTO{
		CustomerSub: sub,
		BookName:    req.BookName,
		Author:      req.Author,
		ISBN:        req.ISBN,
		BookTypeID:  req.BookTypeID,
		ConditionID: req.ConditionID,
		GenreID:     req.GenreID,
		PublisherID: req.PublisherID,
		BookPrice:   req.BookPrice,
		Summary:     req.Summary,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create offer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "offer submitted"})
}

func mapOfferToResponse(o *model.Offer) dto.OfferResponse {
	return dto.OfferResponse{
		ID:          o.ID,
		BookName:    o.BookName,
		Author:      o.Author,
		ISBN:        o.ISBN,
		GenreID:     o.GenreID,
		Genre:       o.Genre.Text,
		ConditionID: o.ConditionID,
		Condition:   o.Condition.Text,
		PublisherID: o.PublisherID,
		Publisher:   o.Publisher.Text,
		BookTypeID:  o.BookTypeID,
		BookType:    o.BookType.Text,
		Summary:     o.Summary,
		OfferStatus: int(o.OfferStatus),
		StatusText:  o.OfferStatus.String(),
		Comment:     o.Comment,
		CustomerID:  o.CustomerID,
		BookPrice:   o.BookPrice,
	}
}
