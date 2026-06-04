package handler

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReferenceDataHandler struct {
	refDataService *service.ReferenceDataService
}

func NewReferenceDataHandler(refDataService *service.ReferenceDataService) *ReferenceDataHandler {
	return &ReferenceDataHandler{refDataService: refDataService}
}

// GET /api/reference-data
func (h *ReferenceDataHandler) GetAll(c *gin.Context) {
	items, err := h.refDataService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/admin/reference-data?dataType=&pageIndex=1&pageSize=10
func (h *ReferenceDataHandler) AdminList(c *gin.Context) {
	var filters repository.ReferenceDataFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.refDataService.List(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/admin/reference-data/:id
func (h *ReferenceDataHandler) AdminGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.refDataService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/admin/reference-data
func (h *ReferenceDataHandler) AdminCreate(c *gin.Context) {
	var req struct {
		DataType int    `json:"dataType"`
		Text     string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.refDataService.Create(model.ReferenceDataType(req.DataType), req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Reference data created"})
}

// PUT /api/admin/reference-data/:id
func (h *ReferenceDataHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		DataType int    `json:"dataType"`
		Text     string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.refDataService.Update(id, model.ReferenceDataType(req.DataType), req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reference data updated"})
}
