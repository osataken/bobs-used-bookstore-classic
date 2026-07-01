package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) GetByIDAndCustomerID(id uint, customerID uint) (*model.Address, error) {
	var address model.Address
	err := r.db.Where("id = ? AND customer_id = ? AND is_active = ?", id, customerID, true).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) ListByCustomerID(customerID uint) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Where("customer_id = ? AND is_active = ?", customerID, true).Order("id ASC").Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *AddressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *AddressRepository) SoftDelete(id uint, customerID uint) error {
	return r.db.Model(&model.Address{}).Where("id = ? AND customer_id = ?", id, customerID).
		Update("is_active", false).Error
}
