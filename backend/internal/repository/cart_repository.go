package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetByCorrelationID(correlationID string) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	err := r.db.Where("correlation_id = ?", correlationID).
		Preload("ShoppingCartItems.Book.Publisher").
		Preload("ShoppingCartItems.Book.BookType").
		Preload("ShoppingCartItems.Book.Genre").
		Preload("ShoppingCartItems.Book.Condition").
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) Create(cart *model.ShoppingCart) error {
	return r.db.Create(cart).Error
}

func (r *CartRepository) Save(cart *model.ShoppingCart) error {
	return r.db.Save(cart).Error
}

func (r *CartRepository) SaveItem(item *model.ShoppingCartItem) error {
	return r.db.Save(item).Error
}

func (r *CartRepository) DeleteItem(item *model.ShoppingCartItem) error {
	return r.db.Where("id = ? AND shopping_cart_id = ?", item.ID, item.ShoppingCartID).Delete(&model.ShoppingCartItem{}).Error
}

func (r *CartRepository) DB() *gorm.DB {
	return r.db
}
