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

type AdminOffersHandler struct {
	offerService *service.OfferService
}

func NewAdminOffersHandler(offerService *service.OfferService) *AdminOffersHandler {
	return &AdminOffersHandler{offerService: offerService}
}

func (h *AdminOffersHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var filters repository.OfferFilters
	if v := c.Query("status"); v != "" {
		status, _ := strconv.Atoi(v)
		s := model.OfferStatus(status)
		filters.Status = &s
	}

	result, err := h.offerService.GetOffersFiltered(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get offers"})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items:      result.Items,
		TotalCount: result.TotalCount,
		PageIndex:  result.PageIndex,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *AdminOffersHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		ID     uint `json:"id" binding:"required"`
		Status int  `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.offerService.UpdateOfferStatus(req.ID, model.OfferStatus(req.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update offer status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "offer status updated"})
}
