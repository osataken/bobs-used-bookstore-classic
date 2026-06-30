package service

import (
	"bobs-used-bookstore-api/internal/database"
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *repository.BookRepository {
	cfg := &config.Config{Database: "local", DatabaseDSN: ":memory:"}
	db, err := database.Initialize(cfg)
	require.NoError(t, err)
	return repository.NewBookRepository(db)
}

func TestBookService_GetBestSellers(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewBookService(repo)

	books, err := svc.GetBestSellers(4)
	assert.NoError(t, err)
	assert.Len(t, books, 4)
	// All returned books should be in stock
	for _, b := range books {
		assert.True(t, b.Quantity > 0)
	}
}

func TestBookService_Search(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewBookService(repo)

	result, err := svc.Search(1, 10, "Apocalypse", "Name")
	assert.NoError(t, err)
	assert.Equal(t, 1, result.TotalCount)
	assert.Equal(t, "2020: The Apocalypse", result.Items[0].Name)
}

func TestBookService_GetByID(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewBookService(repo)

	book, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "2020: The Apocalypse", book.Name)
	assert.Equal(t, 10.95, book.Price)
	assert.NotNil(t, book.Publisher)
	assert.Equal(t, "Arcadia Books", book.Publisher.Text)
}

func TestBook_ReduceStockLevel(t *testing.T) {
	book := &model.Book{Quantity: 10}
	book.ReduceStockLevel(3)
	assert.Equal(t, 7, book.Quantity)

	book.ReduceStockLevel(20)
	assert.Equal(t, 0, book.Quantity)
}

func TestBook_IsLowInStock(t *testing.T) {
	book := &model.Book{Quantity: 5}
	assert.True(t, book.IsLowInStock())

	book.Quantity = 6
	assert.False(t, book.IsLowInStock())
}
