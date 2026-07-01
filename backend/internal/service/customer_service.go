package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) GetBySub(sub string) (*model.Customer, error) {
	return s.customerRepo.GetBySub(sub)
}

func (s *CustomerService) GetByID(id uint) (*model.Customer, error) {
	return s.customerRepo.GetByID(id)
}

func (s *CustomerService) CreateOrUpdate(sub, username, firstName, lastName string) error {
	customer := &model.Customer{
		Sub:       sub,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	}
	return s.customerRepo.CreateOrUpdate(customer)
}
