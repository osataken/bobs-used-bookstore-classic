package service

import (
	"time"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	cartRepo     *repository.CartRepository
	customerRepo *repository.CustomerRepository
	bookRepo     *repository.BookRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, cartRepo *repository.CartRepository, customerRepo *repository.CustomerRepository, bookRepo *repository.BookRepository) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		customerRepo: customerRepo,
		bookRepo:     bookRepo,
	}
}

func (s *OrderService) GetOrder(id uint) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetOrdersByCustomerSub(sub string) ([]model.Order, error) {
	return s.orderRepo.ListBySub(sub)
}

func (s *OrderService) GetOrdersFiltered(filters repository.OrderFilters, pageIndex int, pageSize int) (*repository.PaginatedResult[model.Order], error) {
	return s.orderRepo.ListFiltered(filters, pageIndex, pageSize)
}

func (s *OrderService) GetStatistics() (*repository.OrderStatistics, error) {
	return s.orderRepo.GetStatistics()
}

type CreateOrderInput struct {
	CustomerSub   string
	CorrelationID string
	AddressID     uint
}

func (s *OrderService) CreateOrder(input CreateOrderInput) (uint, error) {
	db := s.orderRepo.DB()

	var orderID uint
	err := db.Transaction(func(tx *gorm.DB) error {
		// Get cart
		cart, err := s.getCartInTx(tx, input.CorrelationID)
		if err != nil {
			return err
		}

		// Get customer
		customer, err := s.customerRepo.GetBySub(input.CustomerSub)
		if err != nil {
			return err
		}

		// Create order
		order := model.NewOrder(customer.ID, input.AddressID)
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Process in-stock cart items
		for _, item := range cart.GetCartItems(true) {
			orderItem := model.OrderItem{
				OrderID:  order.ID,
				BookID:   item.BookID,
				Quantity: item.Quantity,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Reduce stock
			item.Book.ReduceStockLevel(item.Quantity)
			if err := tx.Save(&item.Book).Error; err != nil {
				return err
			}

			// Remove cart item
			if err := tx.Where("id = ? AND shopping_cart_id = ?", item.ID, item.ShoppingCartID).Delete(&model.ShoppingCartItem{}).Error; err != nil {
				return err
			}
		}

		orderID = order.ID
		return nil
	})

	return orderID, err
}

func (s *OrderService) getCartInTx(tx *gorm.DB, correlationID string) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	err := tx.Where("correlation_id = ?", correlationID).
		Preload("ShoppingCartItems.Book").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status model.OrderStatus) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	order.OrderStatus = status
	order.UpdatedOn = time.Now().UTC()
	return s.orderRepo.Save(order)
}

func (s *OrderService) CancelOrder(orderID uint, customerSub string) error {
	order, err := s.orderRepo.GetByIDAndSub(orderID, customerSub)
	if err != nil {
		return err
	}
	order.OrderStatus = model.OrderStatusCancelled
	return s.orderRepo.Save(order)
}
