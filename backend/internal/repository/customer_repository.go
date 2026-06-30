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
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) CreateOrUpdate(sub, username, firstName, lastName string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("sub = ?", sub).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			customer = model.Customer{
				Sub:       sub,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
			}
			if err := r.db.Create(&customer).Error; err != nil {
				return nil, err
			}
			return &customer, nil
		}
		return nil, err
	}

	customer.Username = username
	customer.FirstName = firstName
	customer.LastName = lastName
	if err := r.db.Save(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}
