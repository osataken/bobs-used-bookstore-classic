package handler

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

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
	bookName := c.Query("filters.BookName")
	author := c.Query("filters.Author")
	genreID, _ := strconv.Atoi(c.Query("filters.GenreId"))
	conditionID, _ := strconv.Atoi(c.Query("filters.ConditionId"))

	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var offerStatus *model.OfferStatus
	if s := c.Query("filters.OfferStatus"); s != "" {
		val, err := strconv.Atoi(s)
		if err == nil {
			status := model.OfferStatus(val)
			offerStatus = &status
		}
	}

	result, err := h.offerService.GetAll(pageIndex, pageSize, bookName, author, genreID, conditionID, offerStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminOfferHandler) Approve(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.offerService.Approve(body.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer approved"})
}

func (h *AdminOfferHandler) Reject(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.offerService.Reject(body.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer rejected"})
}

func (h *AdminOfferHandler) Received(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.offerService.MarkReceived(body.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer marked as received"})
}

func (h *AdminOfferHandler) Paid(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.offerService.MarkPaid(body.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer marked as paid"})
}
