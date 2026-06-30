package repository

import (
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

func (r *OrderRepository) GetByCustomerID(customerID int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("OrderItems").Preload("OrderItems.Book").
		Where("customer_id = ?", customerID).Order("id DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetAll(pageIndex, pageSize int, statusFilter *model.OrderStatus, dateFrom, dateTo string) (PaginatedList[model.Order], error) {
	var orders []model.Order
	var totalCount int64

	query := r.db.Model(&model.Order{})

	if statusFilter != nil {
		query = query.Where("order_status = ?", *statusFilter)
	}
	if dateFrom != "" {
		query = query.Where("created_on >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("created_on <= ?", dateTo)
	}

	query.Count(&totalCount)

	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Customer").Preload("Address").Preload("OrderItems").Preload("OrderItems.Book").
		Order("id ASC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return PaginatedList[model.Order]{}, err
	}

	return NewPaginatedList(orders, int(totalCount), pageIndex, pageSize), nil
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) UpdateStatus(id int, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("order_status", status).Error
}

func (r *OrderRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Order{}).Count(&count).Error
	return count, err
}

func (r *OrderRepository) GetPendingCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusPending).Count(&count).Error
	return count, err
}

func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}
