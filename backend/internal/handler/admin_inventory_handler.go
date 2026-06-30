package handler

import (
	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminInventoryHandler struct {
	bookService  *service.BookService
	fileService  *service.FileService
	imageService *service.ImageService
}

func NewAdminInventoryHandler(bookService *service.BookService, fileService *service.FileService, imageService *service.ImageService) *AdminInventoryHandler {
	return &AdminInventoryHandler{
		bookService:  bookService,
		fileService:  fileService,
		imageService: imageService,
	}
}

func (h *AdminInventoryHandler) Index(c *gin.Context) {
	name := c.Query("filters.Name")
	author := c.Query("filters.Author")
	publisherID, _ := strconv.Atoi(c.Query("filters.PublisherId"))
	genreID, _ := strconv.Atoi(c.Query("filters.GenreId"))
	bookTypeID, _ := strconv.Atoi(c.Query("filters.BookTypeId"))
	conditionID, _ := strconv.Atoi(c.Query("filters.ConditionId"))
	lowStock := c.Query("filters.LowStock") == "true"
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	result, err := h.bookService.GetFiltered(name, author, publisherID, genreID, bookTypeID, conditionID, lowStock, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminInventoryHandler) Details(c *gin.Context) {
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

func (h *AdminInventoryHandler) Create(c *gin.Context) {
	var req dto.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := &model.Book{
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

	if err := h.bookService.Create(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *AdminInventoryHandler) Update(c *gin.Context) {
	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.GetByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	book.Name = req.Name
	book.Author = req.Author
	book.BookTypeID = req.BookTypeID
	book.ConditionID = req.ConditionID
	book.GenreID = req.GenreID
	book.PublisherID = req.PublisherID
	book.Year = req.Year
	book.ISBN = req.ISBN
	book.Summary = req.Summary
	book.Price = req.Price
	book.Quantity = req.Quantity

	if err := h.bookService.Update(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *AdminInventoryHandler) UploadImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	file, err := c.FormFile("coverImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// Validate image safety
	safe, err := h.imageService.ValidateImage(src)
	if err != nil || !safe {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image failed safety validation"})
		return
	}

	// Re-open for saving
	src.Close()
	src, _ = file.Open()
	defer src.Close()

	filename := "coverimages/" + file.Filename
	url, err := h.fileService.SaveFile(filename, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Delete old image if it exists
	if book.CoverImageUrl != "" {
		_ = h.fileService.DeleteFile(book.CoverImageUrl)
	}

	book.CoverImageUrl = url
	if err := h.bookService.Update(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coverImageUrl": url})
}
