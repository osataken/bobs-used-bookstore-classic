package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	cartRepo     *repository.ShoppingCartRepository
	customerRepo *repository.CustomerRepository
	bookRepo     *repository.BookRepository
	db           *gorm.DB
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	cartRepo *repository.ShoppingCartRepository,
	customerRepo *repository.CustomerRepository,
	bookRepo *repository.BookRepository,
	db *gorm.DB,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		customerRepo: customerRepo,
		bookRepo:     bookRepo,
		db:           db,
	}
}

func (s *OrderService) GetOrder(id int) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetOrdersByCustomer(sub string) ([]model.Order, error) {
	return s.orderRepo.ListBySub(sub)
}

func (s *OrderService) GetOrders(filters repository.OrderFilters, pageIndex, pageSize int) (*repository.PaginatedList, error) {
	return s.orderRepo.List(filters, pageIndex, pageSize)
}

func (s *OrderService) GetStatistics() (*repository.OrderStatistics, error) {
	return s.orderRepo.GetStatistics()
}

func (s *OrderService) CreateOrder(sub, correlationID string, addressID int) (int, error) {
	var orderID int

	// Get customer and cart outside transaction first
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return 0, err
	}
	if customer == nil {
		return 0, fmt.Errorf("customer not found")
	}

	cart, err := s.cartRepo.GetByCorrelationID(correlationID)
	if err != nil {
		return 0, err
	}
	if cart == nil {
		return 0, fmt.Errorf("shopping cart not found")
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		order := &model.Order{
			CustomerID:   customer.ID,
			AddressID:    addressID,
			DeliveryDate: time.Now().AddDate(0, 0, 7),
			OrderStatus:  model.OrderStatusPending,
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		var itemIDsToRemove []int
		for _, item := range cart.ShoppingCartItems {
			if item.WantToBuy && item.Book.Quantity > 0 {
				orderItem := &model.OrderItem{
					OrderID:  order.ID,
					BookID:   item.BookID,
					Quantity: item.Quantity,
				}
				if err := tx.Create(orderItem).Error; err != nil {
					return err
				}

				// Reduce stock
				newQty := item.Book.Quantity - item.Quantity
				if newQty < 0 {
					newQty = 0
				}
				if err := tx.Model(&model.Book{}).Where("id = ?", item.BookID).Update("quantity", newQty).Error; err != nil {
					return err
				}

				itemIDsToRemove = append(itemIDsToRemove, item.ID)
			}
		}

		// Remove purchased items from cart
		if len(itemIDsToRemove) > 0 {
			if err := tx.Where("id IN ? AND shopping_cart_id = ?", itemIDsToRemove, cart.ID).
				Delete(&model.ShoppingCartItem{}).Error; err != nil {
				return err
			}
		}

		orderID = order.ID
		return nil
	})

	return orderID, err
}

func (s *OrderService) UpdateOrderStatus(orderID int, status model.OrderStatus) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	order.OrderStatus = status
	order.UpdatedOn = time.Now().UTC()
	return s.orderRepo.Save(order)
}

func (s *OrderService) CancelOrder(sub string, orderID int) error {
	order, err := s.orderRepo.GetByIDAndSub(orderID, sub)
	if err != nil {
		return err
	}
	if order == nil {
		return nil
	}

	order.OrderStatus = model.OrderStatusCancelled
	return s.orderRepo.Save(order)
}
