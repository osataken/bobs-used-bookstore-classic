package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	bookService *service.BookService
}

func NewSearchHandler(bookService *service.BookService) *SearchHandler {
	return &SearchHandler{bookService: bookService}
}

func (h *SearchHandler) Index(c *gin.Context) {
	searchString := c.Query("searchString")
	sortBy := c.DefaultQuery("sortBy", "Name")
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	result, err := h.bookService.GetBooks(searchString, sortBy, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search books"})
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

func (h *SearchHandler) Details(c *gin.Context) {
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
