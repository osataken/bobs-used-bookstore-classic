package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type AddressService struct {
	repo *repository.AddressRepository
}

func NewAddressService(repo *repository.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) GetByID(id int) (*model.Address, error) {
	return s.repo.GetByID(id)
}

func (s *AddressService) GetByCustomerID(customerID int) ([]model.Address, error) {
	return s.repo.GetByCustomerID(customerID)
}

func (s *AddressService) Create(address *model.Address) error {
	return s.repo.Create(address)
}

func (s *AddressService) Update(address *model.Address) error {
	return s.repo.Update(address)
}

func (s *AddressService) Delete(id int) error {
	return s.repo.SoftDelete(id)
}
