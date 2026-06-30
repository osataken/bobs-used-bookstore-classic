package service

import (
	"errors"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"gorm.io/gorm"
)

// AddToShoppingCartDTO holds data for adding an item to cart
type AddToShoppingCartDTO struct {
	CorrelationID string
	BookID        int
	Quantity      int
}

// AddToWishlistDTO holds data for adding to wishlist
type AddToWishlistDTO struct {
	CorrelationID string
	BookID        int
}

// MoveWishlistItemDTO holds data for moving wishlist item to cart
type MoveWishlistItemDTO struct {
	CorrelationID      string
	ShoppingCartItemID int
}

// MoveAllWishlistItemsDTO holds data for moving all wishlist items to cart
type MoveAllWishlistItemsDTO struct {
	CorrelationID string
}

// DeleteShoppingCartItemDTO holds data for removing a cart item
type DeleteShoppingCartItemDTO struct {
	CorrelationID      string
	ShoppingCartItemID int
}

type ShoppingCartService struct {
	cartRepo *repository.ShoppingCartRepository
}

func NewShoppingCartService(cartRepo *repository.ShoppingCartRepository) *ShoppingCartService {
	return &ShoppingCartService{cartRepo: cartRepo}
}

func (s *ShoppingCartService) GetShoppingCart(correlationID string) (*model.ShoppingCart, error) {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty cart
			return &model.ShoppingCart{
				CorrelationID: correlationID,
			}, nil
		}
		return nil, err
	}
	return cart, nil
}

func (s *ShoppingCartService) AddToShoppingCart(dto AddToShoppingCartDTO) error {
	cart, err := s.getOrCreateCart(dto.CorrelationID)
	if err != nil {
		return err
	}

	item := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         dto.BookID,
		Quantity:       dto.Quantity,
		WantToBuy:      true,
	}
	return s.cartRepo.AddItem(item)
}

func (s *ShoppingCartService) AddToWishlist(dto AddToWishlistDTO) error {
	cart, err := s.getOrCreateCart(dto.CorrelationID)
	if err != nil {
		return err
	}

	item := &model.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		BookID:         dto.BookID,
		Quantity:       1,
		WantToBuy:      false,
	}
	return s.cartRepo.AddItem(item)
}

func (s *ShoppingCartService) MoveWishlistItemToShoppingCart(dto MoveWishlistItemDTO) error {
	cart, err := s.cartRepo.GetByCorrelationID(dto.CorrelationID)
	if err != nil {
		return err
	}

	for i := range cart.ShoppingCartItems {
		if cart.ShoppingCartItems[i].ID == dto.ShoppingCartItemID && !cart.ShoppingCartItems[i].WantToBuy {
			cart.ShoppingCartItems[i].WantToBuy = true
			return s.cartRepo.UpdateItem(&cart.ShoppingCartItems[i])
		}
	}
	return errors.New("wishlist item not found")
}

func (s *ShoppingCartService) MoveAllWishlistItemsToShoppingCart(dto MoveAllWishlistItemsDTO) error {
	cart, err := s.cartRepo.GetByCorrelationID(dto.CorrelationID)
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

func (s *ShoppingCartService) DeleteShoppingCartItem(dto DeleteShoppingCartItemDTO) error {
	cart, err := s.cartRepo.GetByCorrelationID(dto.CorrelationID)
	if err != nil {
		return err
	}
	return s.cartRepo.DeleteItem(dto.ShoppingCartItemID, cart.ID)
}

func (s *ShoppingCartService) getOrCreateCart(correlationID string) (*model.ShoppingCart, error) {
	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &model.ShoppingCart{
				CorrelationID: correlationID,
			}
			if err := s.cartRepo.Create(cart); err != nil {
				return nil, err
			}
			return cart, nil
		}
		return nil, err
	}
	return cart, nil
}
