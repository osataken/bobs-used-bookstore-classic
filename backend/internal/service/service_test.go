package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use a unique temp file for each test to avoid shared state
	tmpFile := fmt.Sprintf("/tmp/test_%s.db", t.Name())
	os.Remove(tmpFile)
	t.Cleanup(func() { os.Remove(tmpFile) })

	db, err := gorm.Open(sqlite.Open(tmpFile), &gorm.Config{})
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

	return db
}

func seedTestDB(db *gorm.DB) {
	refData := []model.ReferenceDataItem{
		{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"},
		{Entity: model.Entity{ID: 5}, DataType: model.ReferenceDataTypeCondition, Text: "Like New"},
		{Entity: model.Entity{ID: 13}, DataType: model.ReferenceDataTypeGenre, Text: "Science Fiction & Fantasy"},
		{Entity: model.Entity{ID: 15}, DataType: model.ReferenceDataTypePublisher, Text: "Arcadia Books"},
	}
	db.Create(&refData)

	books := []model.Book{
		{Entity: model.Entity{ID: 1}, Name: "Test Book 1", Author: "Author 1", ISBN: "111", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 10.95, Quantity: 25},
		{Entity: model.Entity{ID: 2}, Name: "Test Book 2", Author: "Author 2", ISBN: "222", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 5.00, Quantity: 0},
	}
	db.Create(&books)
}

func TestOrderService_CreateOrder(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(db)

	// Create customer
	customer := &model.Customer{Sub: "test-sub", Username: "testuser", FirstName: "Test", LastName: "User"}
	err := db.Create(customer).Error
	assert.NoError(t, err)

	// Create address
	address := &model.Address{AddressLine1: "123 Main St", City: "Test", State: "TS", Country: "US", ZipCode: "12345", CustomerID: customer.ID, IsActive: true}
	err = db.Create(address).Error
	assert.NoError(t, err)

	// Create cart with items
	cart := &model.ShoppingCart{CorrelationID: "test-sub"}
	err = db.Create(cart).Error
	assert.NoError(t, err)

	cartItems := []model.ShoppingCartItem{
		{ShoppingCartID: cart.ID, BookID: 1, Quantity: 2, WantToBuy: true},
		{ShoppingCartID: cart.ID, BookID: 2, Quantity: 1, WantToBuy: true}, // out of stock - should be excluded
	}
	err = db.Create(&cartItems).Error
	assert.NoError(t, err)

	// Setup service
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	bookRepo := repository.NewBookRepository(db)
	orderService := NewOrderService(orderRepo, cartRepo, customerRepo, bookRepo, db)

	// Create order
	orderID, err := orderService.CreateOrder("test-sub", "test-sub", address.ID)
	assert.NoError(t, err)
	assert.Greater(t, orderID, 0)

	// Verify order
	order, err := orderService.GetOrder(orderID)
	assert.NoError(t, err)
	if order == nil {
		t.Fatal("order is nil")
	}
	assert.Equal(t, 1, len(order.OrderItems)) // only in-stock item
	assert.Equal(t, 1, order.OrderItems[0].BookID)
	assert.Equal(t, 2, order.OrderItems[0].Quantity)

	// Verify stock reduced
	var book model.Book
	db.First(&book, 1)
	assert.Equal(t, 23, book.Quantity) // 25 - 2

	// Verify cart item removed
	var remainingItems []model.ShoppingCartItem
	db.Where("shopping_cart_id = ?", cart.ID).Find(&remainingItems)
	assert.Equal(t, 1, len(remainingItems)) // out-of-stock item stays
}

func TestOrderService_SubTotal_DoesNotMultiplyByQuantity(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(db)

	customer := &model.Customer{Sub: "sub1", Username: "user1", FirstName: "F", LastName: "L"}
	db.Create(customer)
	address := &model.Address{AddressLine1: "1 St", City: "C", State: "S", Country: "US", ZipCode: "11111", CustomerID: customer.ID, IsActive: true}
	db.Create(address)
	cart := &model.ShoppingCart{CorrelationID: "sub1"}
	db.Create(cart)
	db.Create(&model.ShoppingCartItem{ShoppingCartID: cart.ID, BookID: 1, Quantity: 3, WantToBuy: true})

	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	bookRepo := repository.NewBookRepository(db)
	svc := NewOrderService(orderRepo, cartRepo, customerRepo, bookRepo, db)

	orderID, err := svc.CreateOrder("sub1", "sub1", address.ID)
	assert.NoError(t, err)

	order, err := svc.GetOrder(orderID)
	assert.NoError(t, err)

	// SubTotal is just Book.Price (10.95), NOT Price * Quantity
	assert.InDelta(t, 10.95, order.SubTotal(), 0.01)
	assert.InDelta(t, 1.095, order.Tax(), 0.01)
}

func TestOrderService_CancelOrder(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(db)

	customer := &model.Customer{Sub: "sub2", Username: "user2", FirstName: "F", LastName: "L"}
	db.Create(customer)

	order := &model.Order{CustomerID: customer.ID, AddressID: 1, OrderStatus: model.OrderStatusPending}
	db.Create(order)

	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	bookRepo := repository.NewBookRepository(db)
	svc := NewOrderService(orderRepo, cartRepo, customerRepo, bookRepo, db)

	// Cancel order
	err := svc.CancelOrder("sub2", order.ID)
	assert.NoError(t, err)

	// Verify cancelled
	var updated model.Order
	db.First(&updated, order.ID)
	assert.Equal(t, model.OrderStatusCancelled, updated.OrderStatus)

	// Stock NOT restored (intentional behavior)
}

func TestCustomerService_CreateOrUpdate(t *testing.T) {
	db := setupTestDB(t)
	customerRepo := repository.NewCustomerRepository(db)
	svc := NewCustomerService(customerRepo)

	// Create new
	err := svc.CreateOrUpdate("new-sub", "newuser", "New", "User")
	assert.NoError(t, err)

	customer, err := svc.GetBySub("new-sub")
	assert.NoError(t, err)
	assert.Equal(t, "newuser", customer.Username)

	// Update existing
	err = svc.CreateOrUpdate("new-sub", "updateduser", "Updated", "User")
	assert.NoError(t, err)

	customer, err = svc.GetBySub("new-sub")
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", customer.Username)
	assert.Equal(t, "Updated", customer.FirstName)
}

func TestShoppingCartService_AddAndDelete(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(db)

	cartRepo := repository.NewShoppingCartRepository(db)
	svc := NewShoppingCartService(cartRepo)

	// Add to cart
	err := svc.AddToCart("cart-1", 1, 2)
	assert.NoError(t, err)

	cart, err := svc.GetCart("cart-1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cart.ShoppingCartItems))
	assert.True(t, cart.ShoppingCartItems[0].WantToBuy)
	assert.Equal(t, 2, cart.ShoppingCartItems[0].Quantity)

	// Add to wishlist
	err = svc.AddToWishlist("cart-1", 2)
	assert.NoError(t, err)

	cart, err = svc.GetCart("cart-1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cart.ShoppingCartItems))

	// Move wishlist to cart
	var wishlistItem model.ShoppingCartItem
	for _, item := range cart.ShoppingCartItems {
		if !item.WantToBuy {
			wishlistItem = item
		}
	}
	err = svc.MoveToCart("cart-1", wishlistItem.ID)
	assert.NoError(t, err)

	cart, err = svc.GetCart("cart-1")
	assert.NoError(t, err)
	for _, item := range cart.ShoppingCartItems {
		assert.True(t, item.WantToBuy)
	}
}

func TestAddressService_SoftDelete(t *testing.T) {
	db := setupTestDB(t)

	customer := &model.Customer{Sub: "addr-sub", Username: "addruser", FirstName: "A", LastName: "B"}
	db.Create(customer)

	addressRepo := repository.NewAddressRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	svc := NewAddressService(addressRepo, customerRepo)

	// Create address
	err := svc.Create("addr-sub", CreateAddressDTO{
		AddressLine1: "123 Main", City: "City", State: "ST", Country: "US", ZipCode: "00000",
	})
	assert.NoError(t, err)

	addresses, err := svc.GetAddresses("addr-sub")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(addresses))

	// Soft delete
	err = svc.Delete("addr-sub", addresses[0].ID)
	assert.NoError(t, err)

	addresses, err = svc.GetAddresses("addr-sub")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(addresses))

	// Verify still in DB
	var count int64
	db.Model(&model.Address{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestBookService_StockLevels(t *testing.T) {
	// Test IsInStock and IsLowInStock
	book := &model.Book{Quantity: 10}
	assert.True(t, book.IsInStock())
	assert.True(t, book.IsLowInStock()) // 10 > 5

	book.Quantity = 3
	assert.True(t, book.IsInStock())
	assert.False(t, book.IsLowInStock()) // 3 <= 5

	book.Quantity = 0
	assert.False(t, book.IsInStock())

	// ReduceStockLevel
	book.Quantity = 5
	book.ReduceStockLevel(3)
	assert.Equal(t, 2, book.Quantity)

	book.ReduceStockLevel(10) // Should not go below 0
	assert.Equal(t, 0, book.Quantity)
}

func TestOfferService_CreateAndUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	seedTestDB(db)

	customer := &model.Customer{Sub: "offer-sub", Username: "offeruser", FirstName: "O", LastName: "U"}
	db.Create(customer)

	offerRepo := repository.NewOfferRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	svc := NewOfferService(offerRepo, customerRepo)

	err := svc.CreateOffer("offer-sub", CreateOfferDTO{
		BookName:    "My Book",
		Author:      "Me",
		ISBN:        "999",
		BookTypeID:  1,
		ConditionID: 5,
		GenreID:     13,
		PublisherID: 15,
		BookPrice:   12.50,
	})
	assert.NoError(t, err)

	offers, err := svc.GetOffersByCustomer("offer-sub")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(offers))
	assert.Equal(t, model.OfferStatusPendingApproval, offers[0].OfferStatus)

	// Approve
	err = svc.UpdateOfferStatus(offers[0].ID, model.OfferStatusApproved)
	assert.NoError(t, err)

	offer, err := svc.GetOffer(offers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, model.OfferStatusApproved, offer.OfferStatus)
}

func TestReferenceDataService_CRUD(t *testing.T) {
	db := setupTestDB(t)

	refDataRepo := repository.NewReferenceDataRepository(db)
	svc := NewReferenceDataService(refDataRepo)

	// Create
	err := svc.Create(model.ReferenceDataTypeGenre, "Horror")
	assert.NoError(t, err)

	items, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "Horror", items[0].Text)

	// Update
	err = svc.Update(items[0].ID, model.ReferenceDataTypeGenre, "Horror & Thriller")
	assert.NoError(t, err)

	item, err := svc.GetByID(items[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, "Horror & Thriller", item.Text)
}
