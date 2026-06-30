package service

import (
	"testing"
)

func TestShoppingCartServiceAddAndGet(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewShoppingCartService(repos.ShoppingCart)

	// Add to cart
	err := svc.AddToShoppingCart(AddToShoppingCartDTO{
		CorrelationID: "cart-test-1",
		BookID:        1,
		Quantity:      2,
	})
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Get cart
	cart, err := svc.GetShoppingCart("cart-test-1")
	if err != nil {
		t.Fatalf("failed to get cart: %v", err)
	}

	items := cart.GetCartItems(false)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].BookID != 1 {
		t.Errorf("expected book ID 1, got %d", items[0].BookID)
	}
	if items[0].Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", items[0].Quantity)
	}
}

func TestShoppingCartServiceWishlist(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewShoppingCartService(repos.ShoppingCart)

	// Add to wishlist
	err := svc.AddToWishlist(AddToWishlistDTO{
		CorrelationID: "wishlist-test-1",
		BookID:        3,
	})
	if err != nil {
		t.Fatalf("failed to add to wishlist: %v", err)
	}

	// Get cart - should have no cart items
	cart, err := svc.GetShoppingCart("wishlist-test-1")
	if err != nil {
		t.Fatalf("failed to get cart: %v", err)
	}

	cartItems := cart.GetCartItems(false)
	if len(cartItems) != 0 {
		t.Errorf("expected 0 cart items, got %d", len(cartItems))
	}

	wishlistItems := cart.GetWishlistItems()
	if len(wishlistItems) != 1 {
		t.Fatalf("expected 1 wishlist item, got %d", len(wishlistItems))
	}

	// Move to cart
	err = svc.MoveWishlistItemToShoppingCart(MoveWishlistItemDTO{
		CorrelationID:      "wishlist-test-1",
		ShoppingCartItemID: wishlistItems[0].ID,
	})
	if err != nil {
		t.Fatalf("failed to move to cart: %v", err)
	}

	// Verify moved
	cart, _ = svc.GetShoppingCart("wishlist-test-1")
	cartItems = cart.GetCartItems(false)
	if len(cartItems) != 1 {
		t.Errorf("expected 1 cart item after move, got %d", len(cartItems))
	}
	wishlistItems = cart.GetWishlistItems()
	if len(wishlistItems) != 0 {
		t.Errorf("expected 0 wishlist items after move, got %d", len(wishlistItems))
	}
}

func TestShoppingCartServiceDeleteItem(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewShoppingCartService(repos.ShoppingCart)

	// Add item
	_ = svc.AddToShoppingCart(AddToShoppingCartDTO{CorrelationID: "del-test", BookID: 1, Quantity: 1})

	cart, _ := svc.GetShoppingCart("del-test")
	items := cart.GetCartItems(false)
	if len(items) != 1 {
		t.Fatal("expected 1 item")
	}

	// Delete
	err := svc.DeleteShoppingCartItem(DeleteShoppingCartItemDTO{
		CorrelationID:      "del-test",
		ShoppingCartItemID: items[0].ID,
	})
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	// Verify empty
	cart, _ = svc.GetShoppingCart("del-test")
	items = cart.GetCartItems(false)
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestShoppingCartServiceEmptyCart(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewShoppingCartService(repos.ShoppingCart)

	// Get non-existent cart
	cart, err := svc.GetShoppingCart("non-existent")
	if err != nil {
		t.Fatalf("expected no error for non-existent cart, got: %v", err)
	}
	if cart.CorrelationID != "non-existent" {
		t.Errorf("expected correlation ID 'non-existent', got '%s'", cart.CorrelationID)
	}
}
