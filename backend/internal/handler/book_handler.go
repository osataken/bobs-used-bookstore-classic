package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// GET /api/books/search?searchString=&sortBy=Name&pageIndex=1&pageSize=10
func (h *BookHandler) Search(c *gin.Context) {
	searchString := c.DefaultQuery("searchString", "")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/books/:id
func (h *BookHandler) GetByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, book)
}

// GET /api/books/best-selling?count=4
func (h *BookHandler) BestSelling(c *gin.Context) {
	count, _ := strconv.Atoi(c.DefaultQuery("count", "4"))
	books, err := h.bookService.ListBestSelling(count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// POST /api/books/:id/add-to-cart
func (h *BookHandler) AddToCart(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cartID := middleware.GetShoppingCartID(c)
	cartService := c.MustGet("cartService").(*service.ShoppingCartService)

	err = cartService.AddToCart(cartID, bookID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added to shopping cart"})
}

// POST /api/books/:id/add-to-wishlist
func (h *BookHandler) AddToWishlist(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cartID := middleware.GetShoppingCartID(c)
	cartService := c.MustGet("cartService").(*service.ShoppingCartService)

	err = cartService.AddToWishlist(cartID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added to wishlist"})
}

// GET /api/admin/inventory?searchString=&genreId=&bookTypeId=&conditionId=&publisherId=&pageIndex=1&pageSize=10
func (h *BookHandler) AdminList(c *gin.Context) {
	var filters repository.BookFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.bookService.GetBooks(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /api/admin/inventory
func (h *BookHandler) AdminCreate(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		ISBN        string  `json:"isbn"`
		PublisherID int     `json:"publisherId" binding:"required"`
		BookTypeID  int     `json:"bookTypeId" binding:"required"`
		GenreID     int     `json:"genreId" binding:"required"`
		ConditionID int     `json:"conditionId" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Quantity    int     `json:"quantity"`
		Year        *int    `json:"year"`
		Summary     string  `json:"summary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := service.CreateBookDTO{
		Name:        req.Name,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublisherID: req.PublisherID,
		BookTypeID:  req.BookTypeID,
		GenreID:     req.GenreID,
		ConditionID: req.ConditionID,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Year:        req.Year,
		Summary:     req.Summary,
	}

	book, err := h.bookService.Add(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// PUT /api/admin/inventory/:id
func (h *BookHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		ISBN        string  `json:"isbn"`
		PublisherID int     `json:"publisherId" binding:"required"`
		BookTypeID  int     `json:"bookTypeId" binding:"required"`
		GenreID     int     `json:"genreId" binding:"required"`
		ConditionID int     `json:"conditionId" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Quantity    int     `json:"quantity"`
		Year        *int    `json:"year"`
		Summary     string  `json:"summary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := service.UpdateBookDTO{
		BookID:      id,
		Name:        req.Name,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublisherID: req.PublisherID,
		BookTypeID:  req.BookTypeID,
		GenreID:     req.GenreID,
		ConditionID: req.ConditionID,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Year:        req.Year,
		Summary:     req.Summary,
	}

	book, err := h.bookService.Update(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// GET /api/admin/inventory/statistics
func (h *BookHandler) AdminStatistics(c *gin.Context) {
	stats, err := h.bookService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
