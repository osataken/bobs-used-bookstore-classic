package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminOfferHandler struct {
	offerService *service.OfferService
}

func NewAdminOfferHandler(offerService *service.OfferService) *AdminOfferHandler {
	return &AdminOfferHandler{offerService: offerService}
}

func (h *AdminOfferHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var filtersReq dto.OfferFiltersRequest
	_ = c.ShouldBindQuery(&filtersReq)

	filters := repository.OfferFilters{
		BookName: filtersReq.BookName,
		Author:   filtersReq.Author,
		GenreID:  filtersReq.GenreID,
		ConditionID: filtersReq.ConditionID,
	}
	if filtersReq.OfferStatus != nil {
		status := model.OfferStatus(*filtersReq.OfferStatus)
		filters.OfferStatus = &status
	}

	result, err := h.offerService.GetOffers(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get offers"})
		return
	}

	var offers []dto.OfferResponse
	for _, o := range result.Items {
		offers = append(offers, mapOfferToResponse(&o))
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items:       offers,
		PageIndex:   result.PageIndex,
		TotalPages:  result.TotalPages,
		TotalCount:  result.TotalCount,
		HasNext:     result.HasNextPage(),
		HasPrevious: result.HasPreviousPage(),
	})
}

func (h *AdminOfferHandler) Approve(c *gin.Context) {
	h.updateOfferStatus(c, model.OfferStatusApproved)
}

func (h *AdminOfferHandler) Reject(c *gin.Context) {
	h.updateOfferStatus(c, model.OfferStatusRejected)
}

func (h *AdminOfferHandler) Received(c *gin.Context) {
	h.updateOfferStatus(c, model.OfferStatusReceived)
}

func (h *AdminOfferHandler) Paid(c *gin.Context) {
	h.updateOfferStatus(c, model.OfferStatusPaid)
}

func (h *AdminOfferHandler) updateOfferStatus(c *gin.Context, status model.OfferStatus) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.offerService.UpdateStatus(service.UpdateOfferStatusDTO{
		OfferID: id,
		Status:  status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update offer status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "offer status updated"})
}
