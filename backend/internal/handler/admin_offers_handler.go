package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminOffersHandler struct {
	offerService *service.OfferService
}

func NewAdminOffersHandler(offerService *service.OfferService) *AdminOffersHandler {
	return &AdminOffersHandler{offerService: offerService}
}

func (h *AdminOffersHandler) ListOffers(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if v, err := strconv.Atoi(c.Query("offerStatus")); err == nil {
		filters["offerStatus"] = v
	}

	result, err := h.offerService.ListAll(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list offers"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminOffersHandler) Approve(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.offerService.UpdateStatus(id, model.OfferStatusApproved); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer approved"})
}

func (h *AdminOffersHandler) Reject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.offerService.UpdateStatus(id, model.OfferStatusRejected); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer rejected"})
}

func (h *AdminOffersHandler) Received(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.offerService.UpdateStatus(id, model.OfferStatusReceived); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer marked received"})
}

func (h *AdminOffersHandler) Paid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.offerService.UpdateStatus(id, model.OfferStatusPaid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer marked paid"})
}
