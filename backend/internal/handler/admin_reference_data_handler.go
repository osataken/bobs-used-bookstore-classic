package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminReferenceDataHandler struct {
	refDataService *service.ReferenceDataService
}

func NewAdminReferenceDataHandler(refDataService *service.ReferenceDataService) *AdminReferenceDataHandler {
	return &AdminReferenceDataHandler{refDataService: refDataService}
}

func (h *AdminReferenceDataHandler) List(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if v, err := strconv.Atoi(c.Query("dataType")); err == nil {
		filters["dataType"] = v
	}

	result, err := h.refDataService.ListPaginated(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list reference data"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminReferenceDataHandler) Create(c *gin.Context) {
	var req struct {
		DataType int    `json:"dataType" binding:"required"`
		Text     string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item := &model.ReferenceDataItem{
		DataType: model.ReferenceDataType(req.DataType),
		Text:     req.Text,
	}

	if err := h.refDataService.Create(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reference data"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *AdminReferenceDataHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.refDataService.GetItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	item.Text = req.Text
	if err := h.refDataService.Update(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update reference data"})
		return
	}
	c.JSON(http.StatusOK, item)
}
