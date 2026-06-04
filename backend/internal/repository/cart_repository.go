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
	err := r.db.Preload("ShoppingCartItems").Preload("ShoppingCartItems.Book").
		Where("correlation_id = ?", correlationID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *ShoppingCartRepository) Add(cart *model.ShoppingCart) error {
	return r.db.Create(cart).Error
}

func (r *ShoppingCartRepository) Save(cart *model.ShoppingCart) error {
	return r.db.Save(cart).Error
}

func (r *ShoppingCartRepository) AddItem(item *model.ShoppingCartItem) error {
	return r.db.Create(item).Error
}

func (r *ShoppingCartRepository) DeleteItem(id int, shoppingCartID int) error {
	return r.db.Where("id = ? AND shopping_cart_id = ?", id, shoppingCartID).
		Delete(&model.ShoppingCartItem{}).Error
}

func (r *ShoppingCartRepository) UpdateItem(item *model.ShoppingCartItem) error {
	return r.db.Save(item).Error
}

func (r *ShoppingCartRepository) DeleteItemsByCartID(shoppingCartID int, itemIDs []int) error {
	return r.db.Where("shopping_cart_id = ? AND id IN ?", shoppingCartID, itemIDs).
		Delete(&model.ShoppingCartItem{}).Error
}
