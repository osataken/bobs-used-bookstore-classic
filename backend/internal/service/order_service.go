package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"errors"
	"time"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
	bookRepo  *repository.BookRepository
	cartRepo  *repository.ShoppingCartRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, bookRepo *repository.BookRepository, cartRepo *repository.ShoppingCartRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		bookRepo:  bookRepo,
		cartRepo:  cartRepo,
	}
}

func (s *OrderService) GetByID(id int) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetByCustomerID(customerID int) ([]model.Order, error) {
	return s.orderRepo.GetByCustomerID(customerID)
}

func (s *OrderService) GetAll(pageIndex, pageSize int, statusFilter *model.OrderStatus, dateFrom, dateTo string) (repository.PaginatedList[model.Order], error) {
	return s.orderRepo.GetAll(pageIndex, pageSize, statusFilter, dateFrom, dateTo)
}

func (s *OrderService) CreateOrder(customerID, addressID int, cartID int) (*model.Order, error) {
	cartItems, err := s.cartRepo.GetCartItems(cartID)
	if err != nil {
		return nil, err
	}

	// Filter only in-stock items
	var validItems []model.ShoppingCartItem
	for _, item := range cartItems {
		if item.Book != nil && item.Book.IsInStock() {
			validItems = append(validItems, item)
		}
	}

	if len(validItems) == 0 {
		return nil, errors.New("no in-stock items in cart")
	}

	// Use transaction to ensure atomicity
	tx := s.orderRepo.DB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	order := &model.Order{
		CustomerID:   customerID,
		AddressID:    addressID,
		DeliveryDate: time.Now().AddDate(0, 0, 7),
		OrderStatus:  model.OrderStatusPending,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var orderedBookIDs []int
	for _, item := range validItems {
		orderItem := model.OrderItem{
			OrderID:  order.ID,
			BookID:   item.BookID,
			Quantity: item.Quantity,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Reduce stock
		var book model.Book
		if err := tx.First(&book, item.BookID).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		book.ReduceStockLevel(item.Quantity)
		if err := tx.Save(&book).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		orderedBookIDs = append(orderedBookIDs, item.BookID)
	}

	// Remove ordered items from cart
	if err := tx.Where("shopping_cart_id = ? AND book_id IN ? AND want_to_buy = ?", cartID, orderedBookIDs, true).Delete(&model.ShoppingCartItem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(order.ID)
}

func (s *OrderService) CancelOrder(id int) error {
	// Sets status to Cancelled only - stock is NOT restored
	return s.orderRepo.UpdateStatus(id, model.OrderStatusCancelled)
}

func (s *OrderService) UpdateStatus(id int, status model.OrderStatus) error {
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *OrderService) GetTotalCount() (int64, error) {
	return s.orderRepo.GetTotalCount()
}

func (s *OrderService) GetPendingCount() (int64, error) {
	return s.orderRepo.GetPendingCount()
}
