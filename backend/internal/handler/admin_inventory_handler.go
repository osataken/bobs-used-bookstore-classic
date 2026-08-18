package handler

import (
	"io"
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminInventoryHandler struct {
	bookService    *service.BookService
	refDataService *service.ReferenceDataService
}

func NewAdminInventoryHandler(bookService *service.BookService, refDataService *service.ReferenceDataService) *AdminInventoryHandler {
	return &AdminInventoryHandler{bookService: bookService, refDataService: refDataService}
}

func (h *AdminInventoryHandler) ListBooks(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if s := c.Query("searchString"); s != "" {
		filters["searchString"] = s
	}
	if v, err := strconv.Atoi(c.Query("genreId")); err == nil {
		filters["genreId"] = v
	}
	if v, err := strconv.Atoi(c.Query("bookTypeId")); err == nil {
		filters["bookTypeId"] = v
	}
	if v, err := strconv.Atoi(c.Query("conditionId")); err == nil {
		filters["conditionId"] = v
	}
	if v, err := strconv.Atoi(c.Query("publisherId")); err == nil {
		filters["publisherId"] = v
	}

	result, err := h.bookService.ListBooks(filters, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list books"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminInventoryHandler) GetBook(c *gin.Context) {
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

	// Also get reference data for form dropdowns
	refData, _ := h.refDataService.ListAll()

	c.JSON(http.StatusOK, gin.H{
		"book":          book,
		"referenceData": refData,
	})
}

func (h *AdminInventoryHandler) CreateBook(c *gin.Context) {
	name := c.PostForm("name")
	author := c.PostForm("author")
	isbn := c.PostForm("isbn")
	publisherID, _ := strconv.Atoi(c.PostForm("publisherId"))
	bookTypeID, _ := strconv.Atoi(c.PostForm("bookTypeId"))
	genreID, _ := strconv.Atoi(c.PostForm("genreId"))
	conditionID, _ := strconv.Atoi(c.PostForm("conditionId"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	quantity, _ := strconv.Atoi(c.PostForm("quantity"))
	summary := c.PostForm("summary")

	var imageData []byte
	var imageFilename string
	file, header, err := c.Request.FormFile("coverImage")
	if err == nil {
		defer file.Close()
		imageData, _ = io.ReadAll(file)
		imageFilename = header.Filename
	}

	book := &model.Book{
		Name:        name,
		Author:      author,
		ISBN:        isbn,
		PublisherID: publisherID,
		BookTypeID:  bookTypeID,
		GenreID:     genreID,
		ConditionID: conditionID,
		Price:       price,
		Quantity:    quantity,
		Summary:     summary,
	}

	if err := h.bookService.CreateBook(book, imageData, imageFilename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

func (h *AdminInventoryHandler) UpdateBook(c *gin.Context) {
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

	book.Name = c.PostForm("name")
	book.Author = c.PostForm("author")
	book.ISBN = c.PostForm("isbn")
	book.PublisherID, _ = strconv.Atoi(c.PostForm("publisherId"))
	book.BookTypeID, _ = strconv.Atoi(c.PostForm("bookTypeId"))
	book.GenreID, _ = strconv.Atoi(c.PostForm("genreId"))
	book.ConditionID, _ = strconv.Atoi(c.PostForm("conditionId"))
	book.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	book.Quantity, _ = strconv.Atoi(c.PostForm("quantity"))
	book.Summary = c.PostForm("summary")

	var imageData []byte
	var imageFilename string
	file, header, err := c.Request.FormFile("coverImage")
	if err == nil {
		defer file.Close()
		imageData, _ = io.ReadAll(file)
		imageFilename = header.Filename
	}

	if err := h.bookService.UpdateBook(book, imageData, imageFilename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}
