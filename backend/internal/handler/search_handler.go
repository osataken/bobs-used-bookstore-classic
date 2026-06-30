package handler

import (
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	bookService         *service.BookService
	shoppingCartService *service.ShoppingCartService
}

func NewSearchHandler(bookService *service.BookService, cartService *service.ShoppingCartService) *SearchHandler {
	return &SearchHandler{bookService: bookService, shoppingCartService: cartService}
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

	result, err := h.bookService.Search(pageIndex, pageSize, searchString, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SearchHandler) Details(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := h.bookService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *SearchHandler) AddToCart(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}

	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.shoppingCartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.shoppingCartService.AddToCart(cart.ID, bookID, 1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to cart"})
}

func (h *SearchHandler) AddToWishlist(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}

	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.shoppingCartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.shoppingCartService.AddToWishlist(cart.ID, bookID, 1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to wishlist"})
}
