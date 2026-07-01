package repository

import (
	"math"
	"time"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Customer").Preload("Address").Preload("OrderItems.Book").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByIDAndSub(id uint, sub string) (*model.Order, error) {
	var order model.Order
	err := r.db.Joins("Customer").Where("\"Customer\".sub = ?", sub).
		Preload("Address").Preload("OrderItems.Book").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListBySub(sub string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Joins("Customer").Where("\"Customer\".sub = ?", sub).
		Preload("Address").Preload("OrderItems.Book").
		Order("\"orders\".id ASC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

type OrderFilters struct {
	Status *model.OrderStatus
}

func (r *OrderRepository) ListFiltered(filters OrderFilters, pageIndex int, pageSize int) (*PaginatedResult[model.Order], error) {
	query := r.db.Model(&model.Order{}).Order("id ASC")

	if filters.Status != nil {
		query = query.Where("order_status = ?", *filters.Status)
	}

	var totalCount int64
	query.Count(&totalCount)

	var orders []model.Order
	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Customer").Preload("Address").Preload("OrderItems.Book").
		Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedResult[model.Order]{
		Items:      orders,
		TotalCount: totalCount,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) Save(order *model.Order) error {
	return r.db.Save(order).Error
}

type OrderStatistics struct {
	PastDueOrders   int64
	PendingOrders   int64
	OrdersThisMonth int64
	OrdersTotal     int64
}

func (r *OrderRepository) GetStatistics() (*OrderStatistics, error) {
	var stats OrderStatistics
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	r.db.Model(&model.Order{}).Where("delivery_date < ? AND order_status NOT IN (?, ?)",
		now, model.OrderStatusDelivered, model.OrderStatusCancelled).Count(&stats.PastDueOrders)
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusPending).Count(&stats.PendingOrders)
	r.db.Model(&model.Order{}).Where("created_on >= ?", startOfMonth).Count(&stats.OrdersThisMonth)
	r.db.Model(&model.Order{}).Count(&stats.OrdersTotal)

	return &stats, nil
}

func (r *OrderRepository) ListBestSellingBooks(count int) ([]model.Book, error) {
	var books []model.Book
	err := r.db.Raw(`
		SELECT b.* FROM books b
		INNER JOIN (
			SELECT book_id, SUM(quantity) as total_sold
			FROM order_items
			GROUP BY book_id
			ORDER BY total_sold DESC
			LIMIT ?
		) oi ON b.id = oi.book_id
		ORDER BY oi.total_sold DESC
	`, count).Scan(&books).Error

	if err != nil || len(books) == 0 {
		// Fallback: return first N books if no orders exist
		r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
			Order("id ASC").Limit(count).Find(&books)
	}
	return books, nil
}

func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}
