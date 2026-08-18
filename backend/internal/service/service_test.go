package service_test

import (
	"testing"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

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
	require.NoError(t, err)

	return db
}

func seedTestDB(t *testing.T, db *gorm.DB) {
	refData := []model.ReferenceDataItem{
		{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"},
		{Entity: model.Entity{ID: 5}, DataType: model.ReferenceDataTypeCondition, Text: "Like New"},
		{Entity: model.Entity{ID: 13}, DataType: model.ReferenceDataTypeGenre, Text: "Science Fiction & Fantasy"},
		{Entity: model.Entity{ID: 15}, DataType: model.ReferenceDataTypePublisher, Text: "Arcadia Books"},
	}
	require.NoError(t, db.Create(&refData).Error)

	books := []model.Book{
		{Entity: model.Entity{ID: 1}, Name: "2020: The Apocalypse", Author: "Li Juan", ISBN: "6556784356", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 10.95, Quantity: 25},
		{Entity: model.Entity{ID: 2}, Name: "Test Book", Author: "Author", ISBN: "1234567890", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 5.00, Quantity: 3},
	}
	require.NoError(t, db.Create(&books).Error)
}

func TestBookService_SearchBooks(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(t, db)

	bookRepo := repository.NewBookRepository(db)
	cfg := &config.Config{FileServiceMode: "local", ImageValidationMode: "local", LocalStoragePath: "/tmp"}
	fs := service.NewFileService(cfg)
	irs := service.NewImageResizeService()
	ivs := service.NewImageValidationService(cfg)
	bookService := service.NewBookService(bookRepo, fs, irs, ivs)

	result, err := bookService.SearchBooks("Apocalypse", "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.TotalCount)

	books := result.Items.([]model.Book)
	assert.Equal(t, "2020: The Apocalypse", books[0].Name)
}

func TestBookService_GetBook(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(t, db)

	bookRepo := repository.NewBookRepository(db)
	cfg := &config.Config{FileServiceMode: "local", ImageValidationMode: "local", LocalStoragePath: "/tmp"}
	fs := service.NewFileService(cfg)
	irs := service.NewImageResizeService()
	ivs := service.NewImageValidationService(cfg)
	bookService := service.NewBookService(bookRepo, fs, irs, ivs)

	book, err := bookService.GetBook(1)
	require.NoError(t, err)
	assert.Equal(t, "2020: The Apocalypse", book.Name)
	assert.Equal(t, 10.95, book.Price)
}

func TestOrderService_CreateOrder(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(t, db)

	// Create customer
	customer := &model.Customer{Sub: "test-sub", Username: "testuser"}
	require.NoError(t, db.Create(customer).Error)

	// Create address
	address := &model.Address{CustomerID: customer.ID, AddressLine1: "123 Main St", City: "Test", State: "TS", Country: "US", ZipCode: "12345", IsActive: true}
	require.NoError(t, db.Create(address).Error)

	// Create cart with items
	cart := &model.ShoppingCart{CorrelationID: "test-correlation"}
	require.NoError(t, db.Create(cart).Error)

	cartItem := &model.ShoppingCartItem{ShoppingCartID: cart.ID, BookID: 1, Quantity: 2, WantToBuy: true}
	require.NoError(t, db.Create(cartItem).Error)

	bookRepo := repository.NewBookRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	orderService := service.NewOrderService(orderRepo, bookRepo, cartRepo, db)

	order, err := orderService.CreateOrder(customer.ID, address.ID, "test-correlation")
	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, model.OrderStatusPending, order.OrderStatus)

	// Verify stock was reduced
	var book model.Book
	db.First(&book, 1)
	assert.Equal(t, 23, book.Quantity) // 25 - 2

	// Verify cart was cleared
	var remainingItems []model.ShoppingCartItem
	db.Where("\"ShoppingCartID\" = ? AND \"WantToBuy\" = ?", cart.ID, true).Find(&remainingItems)
	assert.Empty(t, remainingItems)
}

func TestOrderService_CancelOrder(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(t, db)

	customer := &model.Customer{Sub: "test-sub", Username: "testuser"}
	require.NoError(t, db.Create(customer).Error)

	address := &model.Address{CustomerID: customer.ID, AddressLine1: "123 Main St", City: "Test", State: "TS", Country: "US", ZipCode: "12345", IsActive: true}
	require.NoError(t, db.Create(address).Error)

	order := &model.Order{CustomerID: customer.ID, AddressID: address.ID, OrderStatus: model.OrderStatusPending}
	require.NoError(t, db.Create(order).Error)

	orderItem := &model.OrderItem{OrderID: order.ID, BookID: 1, Quantity: 2}
	require.NoError(t, db.Create(orderItem).Error)

	bookRepo := repository.NewBookRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	orderService := service.NewOrderService(orderRepo, bookRepo, cartRepo, db)

	err := orderService.CancelOrder(order.ID, customer.ID)
	require.NoError(t, err)

	// Verify status changed
	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusCancelled, updatedOrder.OrderStatus)

	// Verify stock NOT restored (preserves source behavior)
	var book model.Book
	db.First(&book, 1)
	assert.Equal(t, 25, book.Quantity) // unchanged
}

func TestOrder_SubTotal_Bug(t *testing.T) {
	// Test that SubTotal preserves the source bug: Sum(Book.Price) NOT multiplied by Quantity
	order := &model.Order{
		OrderItems: []model.OrderItem{
			{Book: model.Book{Entity: model.Entity{ID: 1}, Price: 10.00}, Quantity: 3},
			{Book: model.Book{Entity: model.Entity{ID: 2}, Price: 5.00}, Quantity: 2},
		},
	}

	// Bug: SubTotal = 10.00 + 5.00 = 15.00 (NOT 10*3 + 5*2 = 40.00)
	assert.Equal(t, 15.0, order.SubTotal())
	assert.Equal(t, 1.5, order.Tax())
	assert.Equal(t, 16.5, order.Total())
}

func TestCustomerService_CreateOrUpdate(t *testing.T) {
	db := setupTestDB(t)

	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)

	// Create
	customer, err := customerService.CreateOrUpdate("sub-123", "user1", "John", "Doe", "john@example.com")
	require.NoError(t, err)
	assert.Equal(t, "sub-123", customer.Sub)
	assert.Equal(t, "John", customer.FirstName)

	// Update
	customer2, err := customerService.CreateOrUpdate("sub-123", "user1", "Jane", "Doe", "jane@example.com")
	require.NoError(t, err)
	assert.Equal(t, customer.ID, customer2.ID)
	assert.Equal(t, "Jane", customer2.FirstName)
}

func TestShoppingCartService(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(t, db)

	bookRepo := repository.NewBookRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	cartService := service.NewShoppingCartService(cartRepo, bookRepo)

	// Add to cart
	err := cartService.AddToCart("cart-1", 1, 1)
	require.NoError(t, err)

	// Get cart
	cart, err := cartService.GetOrCreateCart("cart-1")
	require.NoError(t, err)
	assert.Len(t, cart.ShoppingCartItems, 1)
	assert.True(t, cart.ShoppingCartItems[0].WantToBuy)

	// Add to wishlist
	err = cartService.AddToWishlist("cart-1", 2)
	require.NoError(t, err)

	cart, err = cartService.GetOrCreateCart("cart-1")
	require.NoError(t, err)
	assert.Len(t, cart.ShoppingCartItems, 2)

	wishlistItems := cart.GetWishListItems()
	assert.Len(t, wishlistItems, 1)
	assert.Equal(t, 2, wishlistItems[0].BookID)
}

func TestAddressService_SoftDelete(t *testing.T) {
	db := setupTestDB(t)

	customer := &model.Customer{Sub: "sub-1", Username: "user1"}
	require.NoError(t, db.Create(customer).Error)

	addressRepo := repository.NewAddressRepository(db)
	addressService := service.NewAddressService(addressRepo)

	address := &model.Address{
		AddressLine1: "123 Main St",
		City:         "Test City",
		State:        "TS",
		Country:      "US",
		ZipCode:      "12345",
		CustomerID:   customer.ID,
		IsActive:     true,
	}
	err := addressService.CreateAddress(address)
	require.NoError(t, err)

	// Delete (soft)
	err = addressService.DeleteAddress(address.ID)
	require.NoError(t, err)

	// Should not be found via active-only query
	_, err = addressService.GetAddress(address.ID)
	assert.Error(t, err)

	// But still exists in DB
	var rawAddress model.Address
	db.Unscoped().First(&rawAddress, address.ID)
	assert.False(t, rawAddress.IsActive)
}

func TestBookModel_ReduceStockLevel(t *testing.T) {
	book := &model.Book{Quantity: 10}

	book.ReduceStockLevel(3)
	assert.Equal(t, 7, book.Quantity)

	book.ReduceStockLevel(100)
	assert.Equal(t, 0, book.Quantity) // Clamps to 0
}

func TestEnums(t *testing.T) {
	assert.Equal(t, 0, int(model.OrderStatusPending))
	assert.Equal(t, 1, int(model.OrderStatusOrdered))
	assert.Equal(t, 2, int(model.OrderStatusShipped))
	assert.Equal(t, 3, int(model.OrderStatusDelivered))
	assert.Equal(t, 4, int(model.OrderStatusCancelled))

	assert.Equal(t, 0, int(model.OfferStatusPendingApproval))
	assert.Equal(t, 1, int(model.OfferStatusApproved))
	assert.Equal(t, 2, int(model.OfferStatusReceived))
	assert.Equal(t, 3, int(model.OfferStatusPaid))
	assert.Equal(t, 4, int(model.OfferStatusRejected))

	assert.Equal(t, 0, int(model.ReferenceDataTypePublisher))
	assert.Equal(t, 1, int(model.ReferenceDataTypeCondition))
	assert.Equal(t, 2, int(model.ReferenceDataTypeBookType))
	assert.Equal(t, 3, int(model.ReferenceDataTypeGenre))
}
