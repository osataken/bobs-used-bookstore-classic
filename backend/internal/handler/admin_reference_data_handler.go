package handler

import (
	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminReferenceDataHandler struct {
	refService *service.ReferenceDataService
}

func NewAdminReferenceDataHandler(refService *service.ReferenceDataService) *AdminReferenceDataHandler {
	return &AdminReferenceDataHandler{refService: refService}
}

func (h *AdminReferenceDataHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var dataTypeFilter *model.ReferenceDataType
	if s := c.Query("filters.ReferenceDataType"); s != "" {
		val, err := strconv.Atoi(s)
		if err == nil {
			dt := model.ReferenceDataType(val)
			dataTypeFilter = &dt
		}
	}

	result, err := h.refService.GetAll(pageIndex, pageSize, dataTypeFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminReferenceDataHandler) Create(c *gin.Context) {
	var req dto.CreateReferenceDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &model.ReferenceDataItem{
		DataType: model.ReferenceDataType(req.DataType),
		Text:     req.Text,
	}

	if err := h.refService.Create(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *AdminReferenceDataHandler) Update(c *gin.Context) {
	var req dto.UpdateReferenceDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.refService.GetByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	item.DataType = model.ReferenceDataType(req.DataType)
	item.Text = req.Text

	if err := h.refService.Update(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *AdminReferenceDataHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.refService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}
