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

func (r *AddressRepository) GetByIDAndSub(id int, sub string) (*model.Address, error) {
	var address model.Address
	err := r.db.Joins("JOIN \"Customer\" ON \"Customer\".id = \"Address\".customer_id").
		Where("\"Address\".id = ? AND \"Customer\".sub = ? AND \"Address\".is_active = ?", id, sub, true).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) ListBySub(sub string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Joins("JOIN \"Customer\" ON \"Customer\".id = \"Address\".customer_id").
		Where("\"Customer\".sub = ? AND \"Address\".is_active = ?", sub, true).
		Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *AddressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *AddressRepository) SoftDelete(id int, sub string) error {
	return r.db.Model(&model.Address{}).
		Joins("JOIN \"Customer\" ON \"Customer\".id = \"Address\".customer_id").
		Where("\"Address\".id = ? AND \"Customer\".sub = ?", id, sub).
		Update("is_active", false).Error
}
