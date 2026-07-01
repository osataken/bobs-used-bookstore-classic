package service

import (
	"testing"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbFile := t.TempDir() + "/test.db"
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&model.ReferenceDataItem{},
		&model.Customer{},
		&model.Address{},
		&model.Book{},
		&model.ShoppingCart{},
		&model.ShoppingCartItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.Offer{},
	)
	assert.NoError(t, err)

	return db
}

func TestCreateOrder_ReducesStock(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	customer := &model.Customer{Sub: "test-sub", Username: "testuser", FirstName: "Test", LastName: "User"}
	db.Create(customer)

	address := &model.Address{AddressLine1: "123 Main St", City: "Springfield", State: "IL", Country: "US", ZipCode: "62701", CustomerID: customer.ID, IsActive: true}
	db.Create(address)

	book := &model.Book{Name: "Test Book", Author: "Author", Price: 10.00, Quantity: 5, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)

	// Create reference data
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})

	cart := &model.ShoppingCart{CorrelationID: "test-sub"}
	db.Create(cart)

	cartItem := &model.ShoppingCartItem{ShoppingCartID: cart.ID, BookID: book.ID, Quantity: 2, WantToBuy: true}
	db.Create(cartItem)

	// Setup service
	repos := &repository.Repositories{
		Book:     repository.NewBookRepository(db),
		Order:    repository.NewOrderRepository(db),
		Customer: repository.NewCustomerRepository(db),
		Cart:     repository.NewCartRepository(db),
	}

	orderService := NewOrderService(repos.Order, repos.Cart, repos.Customer, repos.Book)

	// Execute
	orderID, err := orderService.CreateOrder(CreateOrderInput{
		CustomerSub:   "test-sub",
		CorrelationID: "test-sub",
		AddressID:     address.ID,
	})

	// Verify
	assert.NoError(t, err)
	assert.Greater(t, orderID, uint(0))

	// Verify stock was reduced
	var updatedBook model.Book
	db.First(&updatedBook, book.ID)
	assert.Equal(t, 3, updatedBook.Quantity) // 5 - 2 = 3

	// Verify cart item was removed
	var cartItems []model.ShoppingCartItem
	db.Where("shopping_cart_id = ?", cart.ID).Find(&cartItems)
	assert.Empty(t, cartItems)
}

func TestCreateOrder_ExcludesOutOfStock(t *testing.T) {
	db := setupTestDB(t)

	customer := &model.Customer{Sub: "test-sub-2", Username: "testuser2", FirstName: "Test", LastName: "User"}
	db.Create(customer)

	address := &model.Address{AddressLine1: "456 Oak Ave", City: "Portland", State: "OR", Country: "US", ZipCode: "97201", CustomerID: customer.ID, IsActive: true}
	db.Create(address)

	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})

	inStockBook := &model.Book{Name: "In Stock Book", Author: "A", Price: 10.00, Quantity: 5, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(inStockBook)

	outOfStockBook := &model.Book{Name: "Out of Stock Book", Author: "B", Price: 20.00, Quantity: 0, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(outOfStockBook)

	cart := &model.ShoppingCart{CorrelationID: "test-sub-2"}
	db.Create(cart)

	db.Create(&model.ShoppingCartItem{ShoppingCartID: cart.ID, BookID: inStockBook.ID, Quantity: 1, WantToBuy: true})
	db.Create(&model.ShoppingCartItem{ShoppingCartID: cart.ID, BookID: outOfStockBook.ID, Quantity: 1, WantToBuy: true})

	repos := &repository.Repositories{
		Book:     repository.NewBookRepository(db),
		Order:    repository.NewOrderRepository(db),
		Customer: repository.NewCustomerRepository(db),
		Cart:     repository.NewCartRepository(db),
	}

	orderService := NewOrderService(repos.Order, repos.Cart, repos.Customer, repos.Book)

	orderID, err := orderService.CreateOrder(CreateOrderInput{
		CustomerSub:   "test-sub-2",
		CorrelationID: "test-sub-2",
		AddressID:     address.ID,
	})

	assert.NoError(t, err)

	// Verify only 1 order item (the in-stock book)
	var orderItems []model.OrderItem
	db.Where("order_id = ?", orderID).Find(&orderItems)
	assert.Len(t, orderItems, 1)
	assert.Equal(t, inStockBook.ID, orderItems[0].BookID)
}

func TestCancelOrder_DoesNotRestoreStock(t *testing.T) {
	db := setupTestDB(t)

	customer := &model.Customer{Sub: "cancel-sub", Username: "canceluser", FirstName: "Cancel", LastName: "User"}
	db.Create(customer)

	address := &model.Address{AddressLine1: "789 Pine", City: "Seattle", State: "WA", Country: "US", ZipCode: "98101", CustomerID: customer.ID, IsActive: true}
	db.Create(address)

	book := &model.Book{Name: "Cancel Book", Author: "C", Price: 15.00, Quantity: 3, PublisherID: 1, BookTypeID: 1, GenreID: 1, ConditionID: 1}
	db.Create(book)
	db.Create(&model.ReferenceDataItem{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"})

	order := model.NewOrder(customer.ID, address.ID)
	db.Create(order)

	repos := &repository.Repositories{
		Book:     repository.NewBookRepository(db),
		Order:    repository.NewOrderRepository(db),
		Customer: repository.NewCustomerRepository(db),
		Cart:     repository.NewCartRepository(db),
	}

	orderService := NewOrderService(repos.Order, repos.Cart, repos.Customer, repos.Book)
	err := orderService.CancelOrder(order.ID, "cancel-sub")

	assert.NoError(t, err)

	// Verify order is cancelled
	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusCancelled, updatedOrder.OrderStatus)

	// Verify stock was NOT restored
	var updatedBook model.Book
	db.First(&updatedBook, book.ID)
	assert.Equal(t, 3, updatedBook.Quantity) // Unchanged
}
