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
	bookService    *service.BookService
	refDataService *service.ReferenceDataService
}

func NewAdminInventoryHandler(bookService *service.BookService, refDataService *service.ReferenceDataService) *AdminInventoryHandler {
	return &AdminInventoryHandler{
		bookService:    bookService,
		refDataService: refDataService,
	}
}

func (h *AdminInventoryHandler) Index(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := repository.BookFilters{
		SearchString: c.Query("searchString"),
	}

	if v := c.Query("publisherId"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 32)
		uid := uint(id)
		filters.PublisherID = &uid
	}
	if v := c.Query("bookTypeId"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 32)
		uid := uint(id)
		filters.BookTypeID = &uid
	}
	if v := c.Query("genreId"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 32)
		uid := uint(id)
		filters.GenreID = &uid
	}
	if v := c.Query("conditionId"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 32)
		uid := uint(id)
		filters.ConditionID = &uid
	}

	result, err := h.bookService.GetBooksFiltered(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get books"})
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

func (h *AdminInventoryHandler) Details(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	book, err := h.bookService.GetBook(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *AdminInventoryHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		BookTypeID  uint    `json:"bookTypeId" binding:"required"`
		ConditionID uint    `json:"conditionId" binding:"required"`
		GenreID     uint    `json:"genreId" binding:"required"`
		PublisherID uint    `json:"publisherId" binding:"required"`
		Year        *int    `json:"year"`
		ISBN        string  `json:"isbn"`
		Summary     string  `json:"summary"`
		Price       float64 `json:"price" binding:"required"`
		Quantity    int     `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.bookService.Add(service.CreateBookInput{
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
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.ErrorMessage})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book added to inventory"})
}

func (h *AdminInventoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		BookTypeID  uint    `json:"bookTypeId" binding:"required"`
		ConditionID uint    `json:"conditionId" binding:"required"`
		GenreID     uint    `json:"genreId" binding:"required"`
		PublisherID uint    `json:"publisherId" binding:"required"`
		Year        *int    `json:"year"`
		ISBN        string  `json:"isbn"`
		Summary     string  `json:"summary"`
		Price       float64 `json:"price" binding:"required"`
		Quantity    int     `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.bookService.Update(service.UpdateBookInput{
		BookID:      uint(id),
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
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.ErrorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book updated"})
}
