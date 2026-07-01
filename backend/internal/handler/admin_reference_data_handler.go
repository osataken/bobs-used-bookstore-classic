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

type AdminReferenceDataHandler struct {
	refDataService *service.ReferenceDataService
}

func NewAdminReferenceDataHandler(refDataService *service.ReferenceDataService) *AdminReferenceDataHandler {
	return &AdminReferenceDataHandler{refDataService: refDataService}
}

func (h *AdminReferenceDataHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var filters repository.ReferenceDataFilters
	if v := c.Query("dataType"); v != "" {
		dt, _ := strconv.Atoi(v)
		t := model.ReferenceDataType(dt)
		filters.DataType = &t
	}

	result, err := h.refDataService.GetFiltered(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reference data"})
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

func (h *AdminReferenceDataHandler) GetAll(c *gin.Context) {
	items, err := h.refDataService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reference data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"referenceData": items})
}

func (h *AdminReferenceDataHandler) Create(c *gin.Context) {
	var req struct {
		DataType int    `json:"dataType" binding:"required"`
		Text     string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.refDataService.Create(model.ReferenceDataType(req.DataType), req.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reference data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "reference data created"})
}

func (h *AdminReferenceDataHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		DataType int    `json:"dataType" binding:"required"`
		Text     string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.refDataService.Update(uint(id), model.ReferenceDataType(req.DataType), req.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update reference data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reference data updated"})
}
