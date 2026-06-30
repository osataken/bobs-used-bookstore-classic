package service

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/database"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOrderTestDB(t *testing.T) (*repository.OrderRepository, *repository.BookRepository, *repository.ShoppingCartRepository, *repository.CustomerRepository, *repository.AddressRepository) {
	cfg := &config.Config{Database: "local", DatabaseDSN: ":memory:"}
	db, err := database.Initialize(cfg)
	require.NoError(t, err)
	return repository.NewOrderRepository(db),
		repository.NewBookRepository(db),
		repository.NewShoppingCartRepository(db),
		repository.NewCustomerRepository(db),
		repository.NewAddressRepository(db)
}

func TestOrderService_CreateOrder(t *testing.T) {
	orderRepo, bookRepo, cartRepo, customerRepo, addressRepo := setupOrderTestDB(t)
	svc := NewOrderService(orderRepo, bookRepo, cartRepo)

	// Create customer
	customer, err := customerRepo.CreateOrUpdate("test-sub", "testuser", "Test", "User")
	require.NoError(t, err)

	// Create address
	addr := &model.Address{
		AddressLine1: "123 Main St",
		City:         "TestCity",
		State:        "TS",
		Country:      "US",
		ZipCode:      "12345",
		CustomerID:   customer.ID,
		IsActive:     true,
	}
	require.NoError(t, addressRepo.Create(addr))

	// Add items to cart
	cart, err := cartRepo.GetByCorrelationID("test-sub")
	require.NoError(t, err)
	require.NoError(t, cartRepo.AddItem(cart.ID, 1, 2, true)) // Book ID 1, qty 2

	// Place order
	order, err := svc.CreateOrder(customer.ID, addr.ID, cart.ID)
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, model.OrderStatusPending, order.OrderStatus)
	assert.Equal(t, customer.ID, order.CustomerID)

	// Verify stock was reduced
	book, err := bookRepo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 23, book.Quantity) // Was 25, reduced by 2

	// Verify cart items were removed
	items, err := cartRepo.GetCartItems(cart.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestOrderService_CreateOrder_ExcludesOutOfStock(t *testing.T) {
	orderRepo, bookRepo, cartRepo, customerRepo, addressRepo := setupOrderTestDB(t)
	svc := NewOrderService(orderRepo, bookRepo, cartRepo)

	// Create customer and address
	customer, err := customerRepo.CreateOrUpdate("test-sub2", "testuser2", "Test", "User")
	require.NoError(t, err)
	addr := &model.Address{
		AddressLine1: "456 Oak Ave",
		City:         "City",
		State:        "ST",
		Country:      "US",
		ZipCode:      "00000",
		CustomerID:   customer.ID,
		IsActive:     true,
	}
	require.NoError(t, addressRepo.Create(addr))

	// Set book 4 to out of stock
	book4, _ := bookRepo.GetByID(4)
	book4.Quantity = 0
	bookRepo.Update(book4)

	// Add out-of-stock book to cart
	cart, err := cartRepo.GetByCorrelationID("test-sub2")
	require.NoError(t, err)
	require.NoError(t, cartRepo.AddItem(cart.ID, 4, 1, true))
	// Also add in-stock book
	require.NoError(t, cartRepo.AddItem(cart.ID, 1, 1, true))

	order, err := svc.CreateOrder(customer.ID, addr.ID, cart.ID)
	assert.NoError(t, err)
	// Only the in-stock book should be in the order
	assert.Len(t, order.OrderItems, 1)
	assert.Equal(t, 1, order.OrderItems[0].BookID)
}

func TestOrderService_CancelOrder(t *testing.T) {
	orderRepo, bookRepo, cartRepo, customerRepo, addressRepo := setupOrderTestDB(t)
	svc := NewOrderService(orderRepo, bookRepo, cartRepo)

	customer, err := customerRepo.CreateOrUpdate("test-sub3", "testuser3", "Test", "User")
	require.NoError(t, err)
	addr := &model.Address{
		AddressLine1: "789 Pine",
		City:         "City",
		State:        "ST",
		Country:      "US",
		ZipCode:      "11111",
		CustomerID:   customer.ID,
		IsActive:     true,
	}
	require.NoError(t, addressRepo.Create(addr))

	cart, _ := cartRepo.GetByCorrelationID("test-sub3")
	cartRepo.AddItem(cart.ID, 1, 1, true)

	order, err := svc.CreateOrder(customer.ID, addr.ID, cart.ID)
	require.NoError(t, err)

	// Cancel order - stock NOT restored
	bookBefore, _ := bookRepo.GetByID(1)
	qtyBefore := bookBefore.Quantity

	err = svc.CancelOrder(order.ID)
	assert.NoError(t, err)

	// Verify status is cancelled
	cancelledOrder, _ := svc.GetByID(order.ID)
	assert.Equal(t, model.OrderStatusCancelled, cancelledOrder.OrderStatus)

	// Verify stock NOT restored
	bookAfter, _ := bookRepo.GetByID(1)
	assert.Equal(t, qtyBefore, bookAfter.Quantity)
}

func TestOrder_SubTotal_Bug_Preserved(t *testing.T) {
	// SubTotal = Sum(Book.Price) per OrderItem - does NOT multiply by Quantity
	order := &model.Order{
		OrderItems: []model.OrderItem{
			{Book: &model.Book{Entity: model.Entity{ID: 1}, Price: 10.00}, Quantity: 3},
			{Book: &model.Book{Entity: model.Entity{ID: 2}, Price: 5.00}, Quantity: 2},
		},
	}
	// Bug: should be 10*3 + 5*2 = 40, but source just sums prices = 15
	assert.Equal(t, 15.0, order.SubTotal())
	assert.Equal(t, 1.5, order.Tax())
	assert.Equal(t, 16.5, order.Total())
}
