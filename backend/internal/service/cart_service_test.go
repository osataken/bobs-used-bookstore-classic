package service

import (
	"testing"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestAddToCart(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	book := &model.Book{Name: "Cart Book", Author: "A", Price: 10.00, Quantity: 5, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)

	cartRepo := repository.NewCartRepository(db)
	cartService := NewCartService(cartRepo)

	// Add to cart
	err := cartService.AddToCart("test-correlation", book.ID, 1)
	assert.NoError(t, err)

	// Verify
	cart, err := cartService.GetCart("test-correlation")
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Len(t, cart.GetCartItems(false), 1)
	assert.Equal(t, book.ID, cart.GetCartItems(false)[0].BookID)
	assert.Equal(t, 1, cart.GetCartItems(false)[0].Quantity)
}

func TestAddToCart_IncreasesQuantity(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	book := &model.Book{Name: "Cart Book 2", Author: "B", Price: 5.00, Quantity: 10, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)

	cartRepo := repository.NewCartRepository(db)
	cartService := NewCartService(cartRepo)

	err := cartService.AddToCart("test-correlation-2", book.ID, 1)
	assert.NoError(t, err)

	err = cartService.AddToCart("test-correlation-2", book.ID, 2)
	assert.NoError(t, err)

	cart, err := cartService.GetCart("test-correlation-2")
	assert.NoError(t, err)
	items := cart.GetCartItems(false)
	assert.Len(t, items, 1)
	assert.Equal(t, 3, items[0].Quantity) // 1 + 2
}

func TestAddToWishlist(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	book := &model.Book{Name: "Wishlist Book", Author: "C", Price: 8.00, Quantity: 3, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)

	cartRepo := repository.NewCartRepository(db)
	cartService := NewCartService(cartRepo)

	err := cartService.AddToWishlist("wish-correlation", book.ID)
	assert.NoError(t, err)

	cart, err := cartService.GetCart("wish-correlation")
	assert.NoError(t, err)
	assert.Len(t, cart.GetWishlistItems(), 1)
	assert.Empty(t, cart.GetCartItems(false))
}

func TestMoveWishlistItemToCart(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	book := &model.Book{Name: "Move Book", Author: "D", Price: 12.00, Quantity: 7, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)

	cartRepo := repository.NewCartRepository(db)
	cartService := NewCartService(cartRepo)

	err := cartService.AddToWishlist("move-correlation", book.ID)
	assert.NoError(t, err)

	cart, err := cartService.GetCart("move-correlation")
	assert.NoError(t, err)
	wishItems := cart.GetWishlistItems()
	assert.Len(t, wishItems, 1)

	err = cartService.MoveWishlistItemToCart("move-correlation", wishItems[0].ID)
	assert.NoError(t, err)

	cart, err = cartService.GetCart("move-correlation")
	assert.NoError(t, err)
	assert.Empty(t, cart.GetWishlistItems())
	assert.Len(t, cart.GetCartItems(false), 1)
}

func TestDeleteItem(t *testing.T) {
	db := setupTestDB(t)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})
	book := &model.Book{Name: "Delete Book", Author: "E", Price: 6.00, Quantity: 2, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)

	cartRepo := repository.NewCartRepository(db)
	cartService := NewCartService(cartRepo)

	err := cartService.AddToCart("delete-correlation", book.ID, 1)
	assert.NoError(t, err)

	cart, err := cartService.GetCart("delete-correlation")
	assert.NoError(t, err)
	items := cart.GetCartItems(false)
	assert.Len(t, items, 1)

	err = cartService.DeleteItem("delete-correlation", items[0].ID)
	assert.NoError(t, err)

	cart, err = cartService.GetCart("delete-correlation")
	assert.NoError(t, err)
	assert.Empty(t, cart.GetCartItems(false))
}
