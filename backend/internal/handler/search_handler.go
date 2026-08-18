package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	bookService *service.BookService
	cartService *service.ShoppingCartService
}

func NewSearchHandler(bookService *service.BookService, cartService *service.ShoppingCartService) *SearchHandler {
	return &SearchHandler{bookService: bookService, cartService: cartService}
}

func (h *SearchHandler) Search(c *gin.Context) {
	searchString := c.Query("searchString")
	sortBy := c.Query("sortBy")
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.bookService.SearchBooks(searchString, sortBy, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SearchHandler) Details(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	book, err := h.bookService.GetBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}
