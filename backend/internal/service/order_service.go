package service

import (
	"errors"
	"time"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
	bookRepo  *repository.BookRepository
	cartRepo  *repository.ShoppingCartRepository
	db        *gorm.DB
}

func NewOrderService(orderRepo *repository.OrderRepository, bookRepo *repository.BookRepository, cartRepo *repository.ShoppingCartRepository, db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		bookRepo:  bookRepo,
		cartRepo:  cartRepo,
		db:        db,
	}
}

func (s *OrderService) GetOrder(id int) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) ListOrdersByCustomer(customerID int) ([]model.Order, error) {
	return s.orderRepo.ListByCustomer(customerID)
}

func (s *OrderService) ListAllOrders(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.orderRepo.ListAll(filters, pageIndex, pageSize)
}

func (s *OrderService) CreateOrder(customerID, addressID int, correlationID string) (*model.Order, error) {
	// Use transaction to preserve atomicity (matching source EF6 shared unit of work)
	var order *model.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get cart with items
		var cart model.ShoppingCart
		if err := tx.Where("correlation_id = ?", correlationID).
			Preload("ShoppingCartItems", "want_to_buy = ?", true).
			Preload("ShoppingCartItems.Book").
			First(&cart).Error; err != nil {
			return err
		}

		cartItems := cart.ShoppingCartItems
		if len(cartItems) == 0 {
			return errors.New("shopping cart is empty")
		}

		// Create order
		order = &model.Order{
			CustomerID:   customerID,
			AddressID:    addressID,
			DeliveryDate: time.Now().AddDate(0, 0, 7),
			OrderStatus:  model.OrderStatusPending,
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Add order items and reduce stock
		for _, cartItem := range cartItems {
			orderItem := model.OrderItem{
				OrderID:  order.ID,
				BookID:   cartItem.BookID,
				Quantity: cartItem.Quantity,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Reduce stock level
			book := cartItem.Book
			book.ReduceStockLevel(cartItem.Quantity)
			if err := tx.Save(&book).Error; err != nil {
				return err
			}
		}

		// Clear cart (only WantToBuy items)
		if err := tx.Where("shopping_cart_id = ? AND want_to_buy = ?", cart.ID, true).
			Delete(&model.ShoppingCartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with all associations
	return s.orderRepo.GetByID(order.ID)
}

func (s *OrderService) UpdateOrderStatus(id int, status model.OrderStatus) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	order.OrderStatus = status
	return s.orderRepo.Update(order)
}

func (s *OrderService) CancelOrder(id, customerID int) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	if order.CustomerID != customerID {
		return errors.New("unauthorized")
	}
	// Preserve source behavior: sets status to Cancelled only — stock is NOT restored
	order.OrderStatus = model.OrderStatusCancelled
	return s.orderRepo.Update(order)
}

func (s *OrderService) GetStatistics() (*model.OrderStatistics, error) {
	return s.orderRepo.GetStatistics()
}
