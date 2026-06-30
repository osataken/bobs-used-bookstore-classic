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
	err := r.db.Preload("Items").Preload("Items.Book").Preload("Items.Book.Publisher").
		Preload("Items.Book.BookType").Preload("Items.Book.Genre").Preload("Items.Book.Condition").
		Where("correlation_id = ?", correlationID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newCart := model.ShoppingCart{CorrelationID: correlationID}
			if err := r.db.Create(&newCart).Error; err != nil {
				return nil, err
			}
			return &newCart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *ShoppingCartRepository) AddItem(cartID, bookID, quantity int, wantToBuy bool) error {
	var existing model.ShoppingCartItem
	err := r.db.Where("shopping_cart_id = ? AND book_id = ? AND want_to_buy = ?", cartID, bookID, wantToBuy).First(&existing).Error
	if err == nil {
		existing.Quantity += quantity
		return r.db.Save(&existing).Error
	}
	item := model.ShoppingCartItem{
		ShoppingCartID: cartID,
		BookID:         bookID,
		Quantity:       quantity,
		WantToBuy:      wantToBuy,
	}
	return r.db.Create(&item).Error
}

func (r *ShoppingCartRepository) RemoveItem(itemID, cartID int) error {
	return r.db.Where("id = ? AND shopping_cart_id = ?", itemID, cartID).Delete(&model.ShoppingCartItem{}).Error
}

func (r *ShoppingCartRepository) MoveToCart(itemID, cartID int) error {
	return r.db.Model(&model.ShoppingCartItem{}).Where("id = ? AND shopping_cart_id = ?", itemID, cartID).Update("want_to_buy", true).Error
}

func (r *ShoppingCartRepository) MoveAllToCart(cartID int) error {
	return r.db.Model(&model.ShoppingCartItem{}).Where("shopping_cart_id = ? AND want_to_buy = ?", cartID, false).Update("want_to_buy", true).Error
}

func (r *ShoppingCartRepository) GetCartItems(cartID int) ([]model.ShoppingCartItem, error) {
	var items []model.ShoppingCartItem
	err := r.db.Preload("Book").Where("shopping_cart_id = ? AND want_to_buy = ?", cartID, true).Find(&items).Error
	return items, err
}

func (r *ShoppingCartRepository) GetWishlistItems(cartID int) ([]model.ShoppingCartItem, error) {
	var items []model.ShoppingCartItem
	err := r.db.Preload("Book").Where("shopping_cart_id = ? AND want_to_buy = ?", cartID, false).Find(&items).Error
	return items, err
}

func (r *ShoppingCartRepository) RemoveCartItems(cartID int, bookIDs []int) error {
	return r.db.Where("shopping_cart_id = ? AND book_id IN ? AND want_to_buy = ?", cartID, bookIDs, true).Delete(&model.ShoppingCartItem{}).Error
}

func (r *ShoppingCartRepository) DB() *gorm.DB {
	return r.db
}
