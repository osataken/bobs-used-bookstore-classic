package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminInventoryHandler struct {
	bookService *service.BookService
}

func NewAdminInventoryHandler(bookService *service.BookService) *AdminInventoryHandler {
	return &AdminInventoryHandler{bookService: bookService}
}

func (h *AdminInventoryHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var filtersReq dto.BookFiltersRequest
	_ = c.ShouldBindQuery(&filtersReq)

	filters := repository.BookFilters{
		Name:        filtersReq.Name,
		Author:      filtersReq.Author,
		PublisherID: filtersReq.PublisherID,
		GenreID:     filtersReq.GenreID,
		BookTypeID:  filtersReq.BookTypeID,
		ConditionID: filtersReq.ConditionID,
		LowStock:    filtersReq.LowStock,
	}

	result, err := h.bookService.GetBooks(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get inventory"})
		return
	}

	var books []dto.BookResponse
	for _, b := range result.Items {
		books = append(books, mapBookToResponse(&b))
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items:       books,
		PageIndex:   result.PageIndex,
		TotalPages:  result.TotalPages,
		TotalCount:  result.TotalCount,
		HasNext:     result.HasNextPage(),
		HasPrevious: result.HasPreviousPage(),
	})
}

func (h *AdminInventoryHandler) Details(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := h.bookService.GetBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, mapBookToResponse(book))
}

func (h *AdminInventoryHandler) Create(c *gin.Context) {
	var req dto.BookCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	file, header, _ := c.Request.FormFile("coverImage")

	createDTO := service.CreateBookDTO{
		Name:        req.Name,
		Author:      req.Author,
		BookTypeID:  req.BookTypeID,
		ConditionID: req.ConditionID,
		GenreID:     req.GenreID,
		PublisherID: req.PublisherID,
		Year:        req.Year,
		ISBN:        req.ISBN,
		Summary:     req.Summary,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}

	if file != nil {
		createDTO.CoverImage = file
		createDTO.CoverFileName = header.Filename
		defer file.Close()
	}

	result, err := h.bookService.Create(createDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book created"})
}

func (h *AdminInventoryHandler) Update(c *gin.Context) {
	var req dto.BookUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	file, header, _ := c.Request.FormFile("coverImage")

	updateDTO := service.UpdateBookDTO{
		BookID:      req.ID,
		Name:        req.Name,
		Author:      req.Author,
		BookTypeID:  req.BookTypeID,
		ConditionID: req.ConditionID,
		GenreID:     req.GenreID,
		PublisherID: req.PublisherID,
		Year:        req.Year,
		ISBN:        req.ISBN,
		Summary:     req.Summary,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}

	if file != nil {
		updateDTO.CoverImage = file
		updateDTO.CoverFileName = header.Filename
		defer file.Close()
	}

	result, err := h.bookService.Update(updateDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book updated"})
}
