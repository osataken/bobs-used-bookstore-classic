package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	bookService     *service.BookService
	cartService     *service.ShoppingCartService
}

func NewSearchHandler(bookService *service.BookService, cartService *service.ShoppingCartService) *SearchHandler {
	return &SearchHandler{bookService: bookService, cartService: cartService}
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

	result, err := h.bookService.SearchBooks(searchString, sortBy, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
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

func (h *SearchHandler) Details(c *gin.Context) {
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

func (h *SearchHandler) AddItemToShoppingCart(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}

	correlationID := getCartCorrelationID(c)

	err = h.cartService.AddToShoppingCart(service.AddToShoppingCartDTO{
		CorrelationID: correlationID,
		BookID:        bookID,
		Quantity:      1,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to cart"})
}

func (h *SearchHandler) AddItemToWishlist(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}

	correlationID := getCartCorrelationID(c)

	err = h.cartService.AddToWishlist(service.AddToWishlistDTO{
		CorrelationID: correlationID,
		BookID:        bookID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item to wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to wishlist"})
}

// mapBookToResponse converts a Book model to a BookResponse DTO
func mapBookToResponse(b *model.Book) dto.BookResponse {
	return dto.BookResponse{
		ID:            b.ID,
		Name:          b.Name,
		Author:        b.Author,
		Year:          b.Year,
		ISBN:          b.ISBN,
		PublisherID:   b.PublisherID,
		Publisher:     b.Publisher.Text,
		BookTypeID:    b.BookTypeID,
		BookType:      b.BookType.Text,
		GenreID:       b.GenreID,
		Genre:         b.Genre.Text,
		ConditionID:   b.ConditionID,
		Condition:     b.Condition.Text,
		CoverImageURL: b.CoverImageURL,
		Summary:       b.Summary,
		Price:         b.Price,
		Quantity:      b.Quantity,
		IsInStock:     b.IsInStock(),
		IsLowInStock:  b.IsLowInStock(),
	}
}
