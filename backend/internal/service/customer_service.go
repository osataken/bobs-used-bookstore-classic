package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
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

func (s *CustomerService) CreateOrUpdate(sub, username, firstName, lastName string) (*model.Customer, error) {
	return s.repo.CreateOrUpdate(sub, username, firstName, lastName)
}
