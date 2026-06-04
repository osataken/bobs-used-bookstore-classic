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

func (r *AddressRepository) GetBySubAndID(sub string, id int) (*model.Address, error) {
	var address model.Address
	err := r.db.Joins("JOIN customer ON customer.id = address.customer_id").
		Where("customer.sub = ? AND address.id = ? AND address.is_active = ?", sub, id, true).
		First(&address).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) ListBySub(sub string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Joins("JOIN customer ON customer.id = address.customer_id").
		Where("customer.sub = ? AND address.is_active = ?", sub, true).
		Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) Add(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *AddressRepository) Save(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *AddressRepository) SoftDelete(sub string, id int) error {
	var address model.Address
	err := r.db.Joins("JOIN customer ON customer.id = address.customer_id").
		Where("customer.sub = ? AND address.id = ?", sub, id).
		First(&address).Error
	if err != nil {
		return err
	}
	address.IsActive = false
	return r.db.Save(&address).Error
}
