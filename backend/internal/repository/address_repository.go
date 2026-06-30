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

func (r *AddressRepository) GetByID(id int) (*model.Address, error) {
	var address model.Address
	err := r.db.Where("is_active = ?", true).First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) GetByCustomerID(customerID int) ([]model.Address, error) {
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

func (r *AddressRepository) SoftDelete(id int) error {
	return r.db.Model(&model.Address{}).Where("id = ?", id).Update("is_active", false).Error
}
