package repository

import (
	"bobs-used-bookstore-api/internal/model"
	"math"
	"time"

	"gorm.io/gorm"
)

type OrderFilters struct {
	OrderStatus *int `form:"orderStatus"`
}

type OrderStatistics struct {
	PastDueOrders  int64 `json:"pastDueOrders"`
	PendingOrders  int64 `json:"pendingOrders"`
	OrdersThisMonth int64 `json:"ordersThisMonth"`
	OrdersTotal    int64 `json:"ordersTotal"`
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByID(id int) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Customer").Preload("Address").
		Preload("OrderItems").Preload("OrderItems.Book").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByIDAndSub(id int, sub string) (*model.Order, error) {
	var order model.Order
	err := r.db.Joins("JOIN customer ON customer.id = `order`.customer_id").
		Where("`order`.id = ? AND customer.sub = ?", id, sub).
		Preload("OrderItems").Preload("OrderItems.Book").Preload("Address").
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListBySub(sub string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Joins("JOIN customer ON customer.id = `order`.customer_id").
		Where("customer.sub = ?", sub).
		Preload("OrderItems").Preload("OrderItems.Book").Preload("Address").
		Order("`order`.id ASC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) List(filters OrderFilters, pageIndex, pageSize int) (*PaginatedList, error) {
	query := r.db.Model(&model.Order{}).Preload("Customer").Preload("Address").
		Preload("OrderItems").Preload("OrderItems.Book")

	if filters.OrderStatus != nil {
		query = query.Where("order_status = ?", *filters.OrderStatus)
	}

	var totalCount int64
	query.Count(&totalCount)

	var orders []model.Order
	offset := (pageIndex - 1) * pageSize
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedList{
		Items:      orders,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}, nil
}

func (r *OrderRepository) Add(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) Save(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) GetStatistics() (*OrderStatistics, error) {
	var stats OrderStatistics
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	r.db.Model(&model.Order{}).Where("order_status = ? AND delivery_date < ?", model.OrderStatusOrdered, now).Count(&stats.PastDueOrders)
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusPending).Count(&stats.PendingOrders)
	r.db.Model(&model.Order{}).Where("created_on >= ?", startOfMonth).Count(&stats.OrdersThisMonth)
	r.db.Model(&model.Order{}).Count(&stats.OrdersTotal)

	return &stats, nil
}
