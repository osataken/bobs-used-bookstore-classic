package service

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/database"
	"bobs-used-bookstore-api/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCartTestDB(t *testing.T) *repository.ShoppingCartRepository {
	cfg := &config.Config{Database: "local", DatabaseDSN: ":memory:"}
	db, err := database.Initialize(cfg)
	require.NoError(t, err)
	return repository.NewShoppingCartRepository(db)
}

func TestShoppingCartService_AddToCart(t *testing.T) {
	repo := setupCartTestDB(t)
	svc := NewShoppingCartService(repo)

	cart, err := svc.GetOrCreateCart("test-correlation")
	require.NoError(t, err)

	err = svc.AddToCart(cart.ID, 1, 1)
	assert.NoError(t, err)

	items, err := svc.GetCartItems(cart.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, 1, items[0].BookID)
	assert.True(t, items[0].WantToBuy)
}

func TestShoppingCartService_AddToWishlist(t *testing.T) {
	repo := setupCartTestDB(t)
	svc := NewShoppingCartService(repo)

	cart, err := svc.GetOrCreateCart("test-wishlist")
	require.NoError(t, err)

	err = svc.AddToWishlist(cart.ID, 2, 1)
	assert.NoError(t, err)

	items, err := svc.GetWishlistItems(cart.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.False(t, items[0].WantToBuy)
}

func TestShoppingCartService_MoveToCart(t *testing.T) {
	repo := setupCartTestDB(t)
	svc := NewShoppingCartService(repo)

	cart, err := svc.GetOrCreateCart("test-move")
	require.NoError(t, err)

	// Add to wishlist
	err = svc.AddToWishlist(cart.ID, 3, 1)
	require.NoError(t, err)

	wishlistItems, _ := svc.GetWishlistItems(cart.ID)
	require.Len(t, wishlistItems, 1)

	// Move to cart
	err = svc.MoveToCart(wishlistItems[0].ID, cart.ID)
	assert.NoError(t, err)

	// Should now be in cart, not wishlist
	cartItems, _ := svc.GetCartItems(cart.ID)
	assert.Len(t, cartItems, 1)
	newWishlist, _ := svc.GetWishlistItems(cart.ID)
	assert.Len(t, newWishlist, 0)
}

func TestShoppingCartService_MoveAllToCart(t *testing.T) {
	repo := setupCartTestDB(t)
	svc := NewShoppingCartService(repo)

	cart, err := svc.GetOrCreateCart("test-moveall")
	require.NoError(t, err)

	svc.AddToWishlist(cart.ID, 1, 1)
	svc.AddToWishlist(cart.ID, 2, 1)

	err = svc.MoveAllToCart(cart.ID)
	assert.NoError(t, err)

	cartItems, _ := svc.GetCartItems(cart.ID)
	assert.Len(t, cartItems, 2)
	wishlistItems, _ := svc.GetWishlistItems(cart.ID)
	assert.Len(t, wishlistItems, 0)
}

func TestShoppingCartService_RemoveItem(t *testing.T) {
	repo := setupCartTestDB(t)
	svc := NewShoppingCartService(repo)

	cart, _ := svc.GetOrCreateCart("test-remove")
	svc.AddToCart(cart.ID, 1, 1)

	items, _ := svc.GetCartItems(cart.ID)
	require.Len(t, items, 1)

	err := svc.RemoveItem(items[0].ID, cart.ID)
	assert.NoError(t, err)

	items, _ = svc.GetCartItems(cart.ID)
	assert.Len(t, items, 0)
}

func TestShoppingCartService_DuplicateAdd_IncreasesQuantity(t *testing.T) {
	repo := setupCartTestDB(t)
	svc := NewShoppingCartService(repo)

	cart, _ := svc.GetOrCreateCart("test-dup")
	svc.AddToCart(cart.ID, 1, 1)
	svc.AddToCart(cart.ID, 1, 2)

	items, _ := svc.GetCartItems(cart.ID)
	assert.Len(t, items, 1)
	assert.Equal(t, 3, items[0].Quantity)
}
