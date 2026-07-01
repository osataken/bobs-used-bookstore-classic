package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	bookService *service.BookService
}

func NewHomeHandler(bookService *service.BookService) *HomeHandler {
	return &HomeHandler{bookService: bookService}
}

func (h *HomeHandler) Index(c *gin.Context) {
	books, err := h.bookService.ListBestSellingBooks(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get best sellers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bestSellers": books})
}
