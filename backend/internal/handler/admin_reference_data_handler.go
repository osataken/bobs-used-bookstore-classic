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

	var filtersReq dto.ReferenceDataFiltersRequest
	_ = c.ShouldBindQuery(&filtersReq)

	filters := repository.ReferenceDataFilters{}
	if filtersReq.DataType != nil {
		dt := model.ReferenceDataType(*filtersReq.DataType)
		filters.DataType = &dt
	}

	result, err := h.refDataService.List(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reference data"})
		return
	}

	var items []dto.ReferenceDataResponse
	for _, item := range result.Items {
		items = append(items, dto.ReferenceDataResponse{
			ID:       item.ID,
			DataType: int(item.DataType),
			TypeText: item.DataType.String(),
			Text:     item.Text,
		})
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items:       items,
		PageIndex:   result.PageIndex,
		TotalPages:  result.TotalPages,
		TotalCount:  result.TotalCount,
		HasNext:     result.HasNextPage(),
		HasPrevious: result.HasPreviousPage(),
	})
}

func (h *AdminReferenceDataHandler) GetAll(c *gin.Context) {
	items, err := h.refDataService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reference data"})
		return
	}

	var response []dto.ReferenceDataResponse
	for _, item := range items {
		response = append(response, dto.ReferenceDataResponse{
			ID:       item.ID,
			DataType: int(item.DataType),
			TypeText: item.DataType.String(),
			Text:     item.Text,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": response})
}

func (h *AdminReferenceDataHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.refDataService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reference data not found"})
		return
	}

	c.JSON(http.StatusOK, dto.ReferenceDataResponse{
		ID:       item.ID,
		DataType: int(item.DataType),
		TypeText: item.DataType.String(),
		Text:     item.Text,
	})
}

func (h *AdminReferenceDataHandler) Create(c *gin.Context) {
	var req dto.ReferenceDataCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.refDataService.Create(service.CreateReferenceDataDTO{
		DataType: model.ReferenceDataType(req.DataType),
		Text:     req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reference data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "reference data created"})
}

func (h *AdminReferenceDataHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.ReferenceDataUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.refDataService.Update(service.UpdateReferenceDataDTO{
		ID:       id,
		DataType: model.ReferenceDataType(req.DataType),
		Text:     req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update reference data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reference data updated"})
}
