package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) GetByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) GetBySub(sub string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("sub = ?", sub).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) CreateOrUpdate(customer *model.Customer) error {
	var existing model.Customer
	result := r.db.Where("sub = ?", customer.Sub).First(&existing)
	if result.Error == gorm.ErrRecordNotFound {
		return r.db.Create(customer).Error
	}
	existing.Username = customer.Username
	existing.FirstName = customer.FirstName
	existing.LastName = customer.LastName
	return r.db.Save(&existing).Error
}
