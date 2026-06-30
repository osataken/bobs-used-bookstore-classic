package service

import (
	"errors"
	"time"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"gorm.io/gorm"
)

// CreateOrderDTO holds data for creating an order
type CreateOrderDTO struct {
	CustomerSub   string
	CorrelationID string
	AddressID     int
}

// UpdateOrderStatusDTO holds data for updating order status
type UpdateOrderStatusDTO struct {
	OrderID     int
	OrderStatus model.OrderStatus
}

// CancelOrderDTO holds data for cancelling an order
type CancelOrderDTO struct {
	CustomerSub string
	OrderID     int
}

type OrderService struct {
	orderRepo    *repository.OrderRepository
	cartRepo     *repository.ShoppingCartRepository
	customerRepo *repository.CustomerRepository
	bookRepo     *repository.BookRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	cartRepo *repository.ShoppingCartRepository,
	customerRepo *repository.CustomerRepository,
	bookRepo *repository.BookRepository,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		customerRepo: customerRepo,
		bookRepo:     bookRepo,
	}
}

func (s *OrderService) GetOrder(id int) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetOrderByIDAndSub(id int, sub string) (*model.Order, error) {
	return s.orderRepo.GetByIDAndSub(id, sub)
}

func (s *OrderService) GetOrdersByCustomer(sub string) ([]model.Order, error) {
	return s.orderRepo.ListByCustomerSub(sub)
}

func (s *OrderService) GetOrders(filters repository.OrderFilters, pageIndex, pageSize int) (*repository.PaginatedResult[model.Order], error) {
	return s.orderRepo.List(filters, pageIndex, pageSize)
}

func (s *OrderService) GetStatistics() (*repository.OrderStatistics, error) {
	return s.orderRepo.GetStatistics()
}

func (s *OrderService) CreateOrder(dto CreateOrderDTO) (int, error) {
	// Get cart
	cart, err := s.cartRepo.GetByCorrelationID(dto.CorrelationID)
	if err != nil {
		return 0, errors.New("shopping cart not found")
	}

	// Get customer
	customer, err := s.customerRepo.GetBySub(dto.CustomerSub)
	if err != nil {
		return 0, errors.New("customer not found")
	}

	// Create order in a transaction
	var orderID int
	err = s.orderRepo.DB().Transaction(func(tx *gorm.DB) error {
		order := &model.Order{
			CustomerID:   customer.ID,
			AddressID:    dto.AddressID,
			DeliveryDate: time.Now().AddDate(0, 0, 7),
			OrderStatus:  model.OrderStatusPending,
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Process cart items (only in-stock)
		cartItems := cart.GetCartItems(true)
		for _, item := range cartItems {
			// Add order item
			orderItem := &model.OrderItem{
				OrderID:  order.ID,
				BookID:   item.BookID,
				Quantity: item.Quantity,
			}
			if err := tx.Create(orderItem).Error; err != nil {
				return err
			}

			// Reduce stock
			var book model.Book
			if err := tx.First(&book, item.BookID).Error; err != nil {
				return err
			}
			book.ReduceStockLevel(item.Quantity)
			if err := tx.Save(&book).Error; err != nil {
				return err
			}

			// Remove cart item
			if err := tx.Where("id = ? AND shopping_cart_id = ?", item.ID, item.ShoppingCartID).
				Delete(&model.ShoppingCartItem{}).Error; err != nil {
				return err
			}
		}

		orderID = order.ID
		return nil
	})

	return orderID, err
}

func (s *OrderService) UpdateOrderStatus(dto UpdateOrderStatusDTO) error {
	order, err := s.orderRepo.GetByID(dto.OrderID)
	if err != nil {
		return err
	}
	order.OrderStatus = dto.OrderStatus
	return s.orderRepo.Update(order)
}

func (s *OrderService) CancelOrder(dto CancelOrderDTO) error {
	order, err := s.orderRepo.GetByIDAndSub(dto.OrderID, dto.CustomerSub)
	if err != nil {
		return errors.New("order not found or does not belong to customer")
	}
	order.OrderStatus = model.OrderStatusCancelled
	return s.orderRepo.Update(order)
}
