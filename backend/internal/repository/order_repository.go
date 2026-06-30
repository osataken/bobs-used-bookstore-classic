package repository

import (
	"math"
	"time"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

// OrderFilters for filtering order queries
type OrderFilters struct {
	OrderStatus   *model.OrderStatus
	DateFrom      *time.Time
	DateTo        *time.Time
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
		Preload("OrderItems.Book.Genre").Preload("OrderItems.Book.Publisher").
		Preload("OrderItems.Book.BookType").Preload("OrderItems.Book.Condition").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByIDAndSub(id int, sub string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Customer").Preload("Address").
		Preload("OrderItems.Book.Genre").Preload("OrderItems.Book.Publisher").
		Preload("OrderItems.Book.BookType").Preload("OrderItems.Book.Condition").
		Joins("JOIN \"Customer\" ON \"Customer\".id = \"Order\".customer_id").
		Where("\"Order\".id = ? AND \"Customer\".sub = ?", id, sub).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListByCustomerSub(sub string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("OrderItems.Book").
		Joins("JOIN \"Customer\" ON \"Customer\".id = \"Order\".customer_id").
		Where("\"Customer\".sub = ?", sub).
		Order("\"Order\".id DESC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) List(filters OrderFilters, pageIndex, pageSize int) (*PaginatedResult[model.Order], error) {
	query := r.db.Model(&model.Order{}).Preload("Customer").Preload("OrderItems.Book")

	if filters.OrderStatus != nil {
		query = query.Where("order_status = ?", *filters.OrderStatus)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_on >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_on <= ?", *filters.DateTo)
	}

	var totalCount int64
	query.Count(&totalCount)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	offset := (pageIndex - 1) * pageSize

	var orders []model.Order
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[model.Order]{
		Items:      orders,
		PageIndex:  pageIndex,
		TotalPages: totalPages,
		TotalCount: int(totalCount),
	}, nil
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}

// ListBestSellingBookIDs returns book IDs ordered by number of order items
func (r *OrderRepository) ListBestSellingBookIDs(count int) ([]int, error) {
	var results []struct {
		BookID int
		Count  int
	}
	err := r.db.Model(&model.OrderItem{}).
		Select("book_id, COUNT(*) as count").
		Group("book_id").
		Order("count DESC").
		Limit(count).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(results))
	for i, r := range results {
		ids[i] = r.BookID
	}
	return ids, nil
}

// OrderStatistics holds order statistics
type OrderStatistics struct {
	PendingOrders  int `json:"pendingOrders"`
	PastDueOrders  int `json:"pastDueOrders"`
	OrdersThisMonth int `json:"ordersThisMonth"`
	OrdersTotal    int `json:"ordersTotal"`
}

func (r *OrderRepository) GetStatistics() (*OrderStatistics, error) {
	var stats OrderStatistics

	var pending int64
	r.db.Model(&model.Order{}).Where("order_status = ?", model.OrderStatusPending).Count(&pending)
	stats.PendingOrders = int(pending)

	var pastDue int64
	r.db.Model(&model.Order{}).Where("order_status = ? AND delivery_date < ?", model.OrderStatusOrdered, time.Now()).Count(&pastDue)
	stats.PastDueOrders = int(pastDue)

	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	var thisMonth int64
	r.db.Model(&model.Order{}).Where("created_on >= ?", startOfMonth).Count(&thisMonth)
	stats.OrdersThisMonth = int(thisMonth)

	var total int64
	r.db.Model(&model.Order{}).Count(&total)
	stats.OrdersTotal = int(total)

	return &stats, nil
}
