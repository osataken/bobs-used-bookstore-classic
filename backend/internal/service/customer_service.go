package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"

	"gorm.io/gorm"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetByID(id int) (*model.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) GetBySub(sub string) (*model.Customer, error) {
	return s.repo.GetBySub(sub)
}

func (s *CustomerService) CreateOrUpdate(sub, username, firstName, lastName, email string) (*model.Customer, error) {
	customer, err := s.repo.GetBySub(sub)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new customer
			customer = &model.Customer{
				Sub:       sub,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
			}
			if err := s.repo.Create(customer); err != nil {
				return nil, err
			}
			return customer, nil
		}
		return nil, err
	}

	// Update existing
	customer.Username = username
	customer.FirstName = firstName
	customer.LastName = lastName
	customer.Email = email
	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}
	return customer, nil
}
