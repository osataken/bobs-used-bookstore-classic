package repository

import (
	"math"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByID(id int) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Customer").Preload("Address").Preload("OrderItems").Preload("OrderItems.Book").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListByCustomer(customerID int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("customer_id = ?", customerID).
		Preload("OrderItems").Preload("OrderItems.Book").Preload("Address").
		Order("id DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) ListAll(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	var orders []model.Order
	query := r.db.Preload("Customer").Preload("Address").Preload("OrderItems").Preload("OrderItems.Book")

	if v, ok := filters["orderStatus"]; ok {
		if status, valid := v.(int); valid && status >= 0 {
			query = query.Where("order_status = ?", status)
		}
	}
	if v, ok := filters["searchString"]; ok && v != "" {
		// Search by order ID or customer name
		query = query.Where("id = ?", v)
	}

	var totalCount int64
	query.Model(&model.Order{}).Count(&totalCount)
	query = query.Order("id ASC")

	offset := (pageIndex - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return &model.PaginatedList{
		Items:      orders,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) GetStatistics() (*model.OrderStatistics, error) {
	var stats model.OrderStatistics
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusPending).Count(&stats.Pending)
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusOrdered).Count(&stats.Ordered)
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusShipped).Count(&stats.Shipped)
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusDelivered).Count(&stats.Delivered)
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusCancelled).Count(&stats.Cancelled)
	return &stats, nil
}

func (r *OrderRepository) ListBestSellingBookIDs(limit int) ([]int, error) {
	var bookIDs []int
	err := r.db.Model(&model.OrderItem{}).
		Select("book_id").
		Group("book_id").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("book_id", &bookIDs).Error
	return bookIDs, err
}
