package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type ShoppingCartRepository struct {
	db *gorm.DB
}

func NewShoppingCartRepository(db *gorm.DB) *ShoppingCartRepository {
	return &ShoppingCartRepository{db: db}
}

func (r *ShoppingCartRepository) GetByCorrelationID(correlationID string) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	err := r.db.Preload("ShoppingCartItems.Book.Genre").
		Preload("ShoppingCartItems.Book.Publisher").
		Preload("ShoppingCartItems.Book.BookType").
		Preload("ShoppingCartItems.Book.Condition").
		Where("correlation_id = ?", correlationID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *ShoppingCartRepository) Create(cart *model.ShoppingCart) error {
	return r.db.Create(cart).Error
}

func (r *ShoppingCartRepository) Save(cart *model.ShoppingCart) error {
	return r.db.Save(cart).Error
}

func (r *ShoppingCartRepository) AddItem(item *model.ShoppingCartItem) error {
	return r.db.Create(item).Error
}

func (r *ShoppingCartRepository) DeleteItem(itemID, cartID int) error {
	return r.db.Where("id = ? AND shopping_cart_id = ?", itemID, cartID).
		Delete(&model.ShoppingCartItem{}).Error
}

func (r *ShoppingCartRepository) UpdateItem(item *model.ShoppingCartItem) error {
	return r.db.Save(item).Error
}

func (r *ShoppingCartRepository) DeleteItemsByCartID(cartID int, onlyWantToBuy bool) error {
	query := r.db.Where("shopping_cart_id = ?", cartID)
	if onlyWantToBuy {
		query = query.Where("want_to_buy = ?", true)
	}
	return query.Delete(&model.ShoppingCartItem{}).Error
}

func (r *ShoppingCartRepository) DB() *gorm.DB {
	return r.db
}
