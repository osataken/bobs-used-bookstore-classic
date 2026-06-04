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

func (r *CustomerRepository) GetByID(id int) (*model.Customer, error) {
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Add(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) Save(customer *model.Customer) error {
	return r.db.Save(customer).Error
}
