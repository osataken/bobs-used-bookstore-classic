package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type ShoppingCartService struct {
	cartRepo *repository.ShoppingCartRepository
}

func NewShoppingCartService(cartRepo *repository.ShoppingCartRepository) *ShoppingCartService {
	return &ShoppingCartService{cartRepo: cartRepo}
}

func (s *ShoppingCartService) GetCart(correlationID string) (*model.ShoppingCart, error) {
	return s.cartRepo.GetByCorrelationID(correlationID)
}

func (s *ShoppingCartService) AddToCart(correlationID string, bookID, quantity int) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = &model.ShoppingCart{CorrelationID: correlationID}
		err = s.cartRepo.Add(cart)
		if err != nil {
			return err
		}
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
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = &model.ShoppingCart{CorrelationID: correlationID}
		err = s.cartRepo.Add(cart)
		if err != nil {
			return err
		}
	}

	item := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         bookID,
		Quantity:       1,
		WantToBuy:      false,
	}
	return s.cartRepo.AddItem(item)
}

func (s *ShoppingCartService) MoveToCart(correlationID string, itemID int) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}

	for i := range cart.ShoppingCartItems {
		if cart.ShoppingCartItems[i].ID == itemID {
			cart.ShoppingCartItems[i].WantToBuy = true
			return s.cartRepo.UpdateItem(&cart.ShoppingCartItems[i])
		}
	}
	return nil
}

func (s *ShoppingCartService) MoveAllToCart(correlationID string) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}

	for i := range cart.ShoppingCartItems {
		if !cart.ShoppingCartItems[i].WantToBuy {
			cart.ShoppingCartItems[i].WantToBuy = true
			err = s.cartRepo.UpdateItem(&cart.ShoppingCartItems[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ShoppingCartService) DeleteItem(correlationID string, itemID int) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}
	return s.cartRepo.DeleteItem(itemID, cart.ID)
}
