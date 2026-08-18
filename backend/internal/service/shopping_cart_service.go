package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type ShoppingCartService struct {
	cartRepo *repository.ShoppingCartRepository
	bookRepo *repository.BookRepository
}

func NewShoppingCartService(cartRepo *repository.ShoppingCartRepository, bookRepo *repository.BookRepository) *ShoppingCartService {
	return &ShoppingCartService{cartRepo: cartRepo, bookRepo: bookRepo}
}

func (s *ShoppingCartService) GetOrCreateCart(correlationID string) (*model.ShoppingCart, error) {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err == nil {
		return cart, nil
	}
	// Create new cart
	cart = &model.ShoppingCart{CorrelationID: correlationID}
	if err := s.cartRepo.Create(cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *ShoppingCartService) AddToCart(correlationID string, bookID int, quantity int) error {
	cart, err := s.GetOrCreateCart(correlationID)
	if err != nil {
		return err
	}

	item := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         bookID,
		Quantity:       quantity,
		WantToBuy:      true,
	}
	return s.cartRepo.AddItem(item)
}

func (s *ShoppingCartService) AddToWishlist(correlationID string, bookID int) error {
	cart, err := s.GetOrCreateCart(correlationID)
	if err != nil {
		return err
	}

	item := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         bookID,
		Quantity:       1,
		WantToBuy:      false,
	}
	return s.cartRepo.AddItem(item)
}

func (s *ShoppingCartService) MoveWishlistItemToCart(correlationID string, itemID int) error {
	cart, err := s.GetOrCreateCart(correlationID)
	if err != nil {
		return err
	}

	// Find the item
	for i := range cart.ShoppingCartItems {
		if cart.ShoppingCartItems[i].ID == itemID {
			cart.ShoppingCartItems[i].WantToBuy = true
			return s.cartRepo.UpdateItem(&cart.ShoppingCartItems[i])
		}
	}
	return nil
}

func (s *ShoppingCartService) MoveAllWishlistItemsToCart(correlationID string) error {
	cart, err := s.GetOrCreateCart(correlationID)
	if err != nil {
		return err
	}

	for i := range cart.ShoppingCartItems {
		if !cart.ShoppingCartItems[i].WantToBuy {
			cart.ShoppingCartItems[i].WantToBuy = true
			if err := s.cartRepo.UpdateItem(&cart.ShoppingCartItems[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ShoppingCartService) DeleteItem(correlationID string, itemID int) error {
	cart, err := s.GetOrCreateCart(correlationID)
	if err != nil {
		return err
	}
	return s.cartRepo.DeleteItem(itemID, cart.ID)
}
