package service

import (
	"testing"

	"bobs-used-bookstore-api/internal/model"
)

func TestOrderServiceCreateOrder(t *testing.T) {
	repos := setupTestDB(t)
	orderSvc := NewOrderService(repos.Order, repos.ShoppingCart, repos.Customer, repos.Book)
	customerSvc := NewCustomerService(repos.Customer)
	cartSvc := NewShoppingCartService(repos.ShoppingCart)

	// Create customer
	err := customerSvc.CreateOrUpdate(CreateOrUpdateCustomerDTO{
		Sub:       "test-sub-123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	})
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	// Add items to cart
	err = cartSvc.AddToShoppingCart(AddToShoppingCartDTO{
		CorrelationID: "test-sub-123",
		BookID:        1,
		Quantity:      2,
	})
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Create address
	addressSvc := NewAddressService(repos.Address, repos.Customer)
	err = addressSvc.Create(CreateAddressDTO{
		AddressLine1: "123 Main St",
		City:         "Springfield",
		State:        "IL",
		Country:      "US",
		ZipCode:      "62701",
		CustomerSub:  "test-sub-123",
	})
	if err != nil {
		t.Fatalf("failed to create address: %v", err)
	}

	addresses, _ := addressSvc.GetAddresses("test-sub-123")
	if len(addresses) == 0 {
		t.Fatal("no addresses found")
	}

	// Create order
	orderID, err := orderSvc.CreateOrder(CreateOrderDTO{
		CustomerSub:   "test-sub-123",
		CorrelationID: "test-sub-123",
		AddressID:     addresses[0].ID,
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	if orderID == 0 {
		t.Error("expected non-zero order ID")
	}

	// Verify stock was reduced
	book, _ := repos.Book.GetByID(1)
	if book.Quantity != 23 { // was 25, reduced by 2
		t.Errorf("expected quantity 23, got %d", book.Quantity)
	}

	// Verify cart is empty
	cart, _ := cartSvc.GetShoppingCart("test-sub-123")
	cartItems := cart.GetCartItems(false)
	if len(cartItems) != 0 {
		t.Errorf("expected empty cart, got %d items", len(cartItems))
	}
}

func TestOrderServiceCancelOrder(t *testing.T) {
	repos := setupTestDB(t)
	orderSvc := NewOrderService(repos.Order, repos.ShoppingCart, repos.Customer, repos.Book)
	customerSvc := NewCustomerService(repos.Customer)
	cartSvc := NewShoppingCartService(repos.ShoppingCart)
	addressSvc := NewAddressService(repos.Address, repos.Customer)

	// Setup
	_ = customerSvc.CreateOrUpdate(CreateOrUpdateCustomerDTO{Sub: "cancel-test", Username: "user", FirstName: "F", LastName: "L"})
	_ = cartSvc.AddToShoppingCart(AddToShoppingCartDTO{CorrelationID: "cancel-test", BookID: 2, Quantity: 1})
	_ = addressSvc.Create(CreateAddressDTO{AddressLine1: "1 St", City: "C", State: "S", Country: "US", ZipCode: "12345", CustomerSub: "cancel-test"})
	addrs, _ := addressSvc.GetAddresses("cancel-test")
	orderID, _ := orderSvc.CreateOrder(CreateOrderDTO{CustomerSub: "cancel-test", CorrelationID: "cancel-test", AddressID: addrs[0].ID})

	// Cancel
	err := orderSvc.CancelOrder(CancelOrderDTO{CustomerSub: "cancel-test", OrderID: orderID})
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}

	// Verify status
	order, _ := orderSvc.GetOrder(orderID)
	if order.OrderStatus != model.OrderStatusCancelled {
		t.Errorf("expected Cancelled status, got %s", order.OrderStatus.String())
	}

	// Stock should NOT be restored (per source behavior)
	book, _ := repos.Book.GetByID(2)
	if book.Quantity != 2 { // was 3, reduced by 1, NOT restored
		t.Errorf("expected quantity 2, got %d", book.Quantity)
	}
}
