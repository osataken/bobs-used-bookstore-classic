package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type CartService struct {
	cartRepo *repository.CartRepository
}

func NewCartService(cartRepo *repository.CartRepository) *CartService {
	return &CartService{cartRepo: cartRepo}
}

func (s *CartService) GetCart(correlationID string) (*model.ShoppingCart, error) {
	return s.cartRepo.GetByCorrelationID(correlationID)
}

func (s *CartService) AddToCart(correlationID string, bookID uint, quantity int) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = &model.ShoppingCart{CorrelationID: correlationID}
		if err := s.cartRepo.Create(cart); err != nil {
			return err
		}
		// Reload to get ID
		cart, err = s.cartRepo.GetByCorrelationID(correlationID)
		if err != nil {
			return err
		}
	}

	// Check if item already exists in cart
	for i, item := range cart.ShoppingCartItems {
		if item.BookID == bookID && item.WantToBuy {
			cart.ShoppingCartItems[i].Quantity += quantity
			return s.cartRepo.SaveItem(&cart.ShoppingCartItems[i])
		}
	}

	// Add new item
	newItem := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         bookID,
		Quantity:       quantity,
		WantToBuy:      true,
	}
	return s.cartRepo.DB().Create(newItem).Error
}

func (s *CartService) AddToWishlist(correlationID string, bookID uint) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = &model.ShoppingCart{CorrelationID: correlationID}
		if err := s.cartRepo.Create(cart); err != nil {
			return err
		}
		cart, err = s.cartRepo.GetByCorrelationID(correlationID)
		if err != nil {
			return err
		}
	}

	// Check if already in wishlist
	for _, item := range cart.ShoppingCartItems {
		if item.BookID == bookID && !item.WantToBuy {
			return nil
		}
	}

	newItem := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         bookID,
		Quantity:       1,
		WantToBuy:      false,
	}
	return s.cartRepo.DB().Create(newItem).Error
}

func (s *CartService) MoveWishlistItemToCart(correlationID string, itemID uint) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}

	for i, item := range cart.ShoppingCartItems {
		if item.ID == itemID && !item.WantToBuy {
			cart.ShoppingCartItems[i].WantToBuy = true
			return s.cartRepo.SaveItem(&cart.ShoppingCartItems[i])
		}
	}
	return nil
}

func (s *CartService) MoveAllWishlistItemsToCart(correlationID string) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}

	for i, item := range cart.ShoppingCartItems {
		if !item.WantToBuy {
			cart.ShoppingCartItems[i].WantToBuy = true
			if err := s.cartRepo.SaveItem(&cart.ShoppingCartItems[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *CartService) DeleteItem(correlationID string, itemID uint) error {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}

	for _, item := range cart.ShoppingCartItems {
		if item.ID == itemID {
			return s.cartRepo.DeleteItem(&item)
		}
	}
	return nil
}
