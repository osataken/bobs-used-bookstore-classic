package service

import (
	"testing"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetBooks_Search(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 5}, DataType: model.ReferenceDataTypeCondition, Text: "Like New"})
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 13}, DataType: model.ReferenceDataTypeGenre, Text: "Sci-Fi"})
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 15}, DataType: model.ReferenceDataTypePublisher, Text: "Publisher"})

	db.Create(&model.Book{Name: "Alpha Book", Author: "John", Price: 10.00, Quantity: 5, PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5})
	db.Create(&model.Book{Name: "Beta Book", Author: "Jane", Price: 15.00, Quantity: 3, PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5})
	db.Create(&model.Book{Name: "Gamma Novel", Author: "Bob", Price: 8.00, Quantity: 0, PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5})

	bookRepo := repository.NewBookRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	bookService := NewBookService(bookRepo, orderRepo, &LocalFileService{basePath: "/tmp"}, &LocalImageValidationService{})

	// Search for "Book"
	result, err := bookService.GetBooks("Book", "Name", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.TotalCount)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, "Alpha Book", result.Items[0].Name) // sorted by name
}

func TestGetBooks_Pagination(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 5}, DataType: model.ReferenceDataTypeCondition, Text: "Like New"})
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 13}, DataType: model.ReferenceDataTypeGenre, Text: "Sci-Fi"})
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 15}, DataType: model.ReferenceDataTypePublisher, Text: "Publisher"})

	for i := 0; i < 15; i++ {
		db.Create(&model.Book{Name: "Book " + string(rune('A'+i)), Author: "Author", Price: 10.00, Quantity: 5, PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5})
	}

	bookRepo := repository.NewBookRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	bookService := NewBookService(bookRepo, orderRepo, &LocalFileService{basePath: "/tmp"}, &LocalImageValidationService{})

	result, err := bookService.GetBooks("", "Name", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), result.TotalCount)
	assert.Len(t, result.Items, 10)
	assert.Equal(t, 2, result.TotalPages)

	result2, err := bookService.GetBooks("", "Name", 2, 10)
	assert.NoError(t, err)
	assert.Len(t, result2.Items, 5)
}

func TestSubTotal_Bug_DoesNotMultiplyByQuantity(t *testing.T) {
	// This test verifies the intentional source bug is preserved
	order := &model.Order{
		OrderItems: []model.OrderItem{
			{Book: model.Book{Price: 10.00}, Quantity: 3},
			{Book: model.Book{Price: 5.00}, Quantity: 2},
		},
	}

	// SubTotal should be 10 + 5 = 15, NOT 10*3 + 5*2 = 40
	assert.Equal(t, 15.0, order.SubTotal())
	assert.Equal(t, 1.5, order.Tax())
	assert.Equal(t, 16.5, order.Total())
}
