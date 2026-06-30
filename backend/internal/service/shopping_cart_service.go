package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type ShoppingCartService struct {
	repo *repository.ShoppingCartRepository
}

func NewShoppingCartService(repo *repository.ShoppingCartRepository) *ShoppingCartService {
	return &ShoppingCartService{repo: repo}
}

func (s *ShoppingCartService) GetOrCreateCart(correlationID string) (*model.ShoppingCart, error) {
	return s.repo.GetByCorrelationID(correlationID)
}

func (s *ShoppingCartService) AddToCart(cartID, bookID, quantity int) error {
	return s.repo.AddItem(cartID, bookID, quantity, true)
}

func (s *ShoppingCartService) AddToWishlist(cartID, bookID, quantity int) error {
	return s.repo.AddItem(cartID, bookID, quantity, false)
}

func (s *ShoppingCartService) RemoveItem(itemID, cartID int) error {
	return s.repo.RemoveItem(itemID, cartID)
}

func (s *ShoppingCartService) MoveToCart(itemID, cartID int) error {
	return s.repo.MoveToCart(itemID, cartID)
}

func (s *ShoppingCartService) MoveAllToCart(cartID int) error {
	return s.repo.MoveAllToCart(cartID)
}

func (s *ShoppingCartService) GetCartItems(cartID int) ([]model.ShoppingCartItem, error) {
	return s.repo.GetCartItems(cartID)
}

func (s *ShoppingCartService) GetWishlistItems(cartID int) ([]model.ShoppingCartItem, error) {
	return s.repo.GetWishlistItems(cartID)
}
