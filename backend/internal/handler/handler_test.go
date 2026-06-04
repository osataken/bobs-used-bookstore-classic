package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&model.ReferenceDataItem{},
		&model.Book{},
		&model.Customer{},
		&model.Address{},
		&model.Order{},
		&model.OrderItem{},
		&model.ShoppingCart{},
		&model.ShoppingCartItem{},
		&model.Offer{},
	)
	assert.NoError(t, err)

	// Seed data
	refData := []model.ReferenceDataItem{
		{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"},
		{Entity: model.Entity{ID: 5}, DataType: model.ReferenceDataTypeCondition, Text: "Like New"},
		{Entity: model.Entity{ID: 13}, DataType: model.ReferenceDataTypeGenre, Text: "Science Fiction & Fantasy"},
		{Entity: model.Entity{ID: 15}, DataType: model.ReferenceDataTypePublisher, Text: "Arcadia Books"},
	}
	db.Create(&refData)

	books := []model.Book{
		{Entity: model.Entity{ID: 1}, Name: "Test Book", Author: "Test Author", ISBN: "1234567890", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 10.95, Quantity: 25},
	}
	db.Create(&books)

	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	bookRepo := repository.NewBookRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	fileService := service.NewLocalFileService("/tmp/test-uploads")
	imageService := service.NewLocalImageValidationService()
	bookService := service.NewBookService(bookRepo, fileService, imageService, "/tmp/test-uploads")
	cartService := service.NewShoppingCartService(cartRepo)
	refDataRepo := repository.NewReferenceDataRepository(db)
	refDataService := service.NewReferenceDataService(refDataRepo)

	bookHandler := NewBookHandler(bookService)
	cartHandler := NewCartHandler(cartService)
	refDataHandler := NewReferenceDataHandler(refDataService)

	r.Use(func(c *gin.Context) {
		c.Set("cartService", cartService)
		c.Set(middleware.ContextKeySub, "test-sub")
		c.Set(middleware.ContextKeyRole, "Administrators")
		c.Set("shoppingCartId", "test-cart-id")
		c.Next()
	})

	api := r.Group("/api")
	api.GET("/books/search", bookHandler.Search)
	api.GET("/books/best-selling", bookHandler.BestSelling)
	api.GET("/books/:id", bookHandler.GetByID)
	api.POST("/books/:id/add-to-cart", bookHandler.AddToCart)
	api.POST("/books/:id/add-to-wishlist", bookHandler.AddToWishlist)
	api.GET("/cart", cartHandler.GetCart)
	api.DELETE("/cart/:itemId", cartHandler.DeleteItem)
	api.GET("/wishlist", cartHandler.GetWishlist)
	api.GET("/reference-data", refDataHandler.GetAll)

	return r
}

func TestBookHandler_Search(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books/search?searchString=test&pageIndex=1&pageSize=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), result["totalCount"])
}

func TestBookHandler_GetByID(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var book map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, "Test Book", book["name"])
}

func TestBookHandler_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookHandler_AddToCart(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/books/1/add-to-cart", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check cart
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/cart", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var cartResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cartResp)
	items := cartResp["items"].([]interface{})
	assert.Equal(t, 1, len(items))
}

func TestBookHandler_AddToWishlist(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/books/1/add-to-wishlist", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check wishlist
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/wishlist", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	assert.Equal(t, 1, len(items))
}

func TestCartHandler_EmptyCart(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cart", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	assert.Equal(t, 0, len(items))
}

func TestReferenceDataHandler_GetAll(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/reference-data", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var items []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &items)
	assert.Equal(t, 4, len(items))
}

func TestBookHandler_BestSelling(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books/best-selling?count=4", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Suppress unused import warning
var _ = strings.NewReader
