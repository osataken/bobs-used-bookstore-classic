package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"time"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) GetByID(id int) (*model.Customer, error) {
	return s.customerRepo.GetByID(id)
}

func (s *CustomerService) GetBySub(sub string) (*model.Customer, error) {
	return s.customerRepo.GetBySub(sub)
}

func (s *CustomerService) CreateOrUpdate(sub, username, firstName, lastName string) error {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return err
	}

	if customer == nil {
		customer = &model.Customer{}
	}

	customer.Sub = sub
	customer.Username = username
	customer.FirstName = firstName
	customer.LastName = lastName
	customer.UpdatedOn = time.Now().UTC()

	if customer.ID == 0 {
		return s.customerRepo.Add(customer)
	}
	return s.customerRepo.Save(customer)
}
