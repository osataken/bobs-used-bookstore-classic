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

func (s *AddressService) GetAddress(id int) (*model.Address, error) {
	return s.repo.GetByID(id)
}

func (s *AddressService) ListByCustomer(customerID int) ([]model.Address, error) {
	return s.repo.ListByCustomer(customerID)
}

func (s *AddressService) CreateAddress(address *model.Address) error {
	return s.repo.Create(address)
}

func (s *AddressService) UpdateAddress(address *model.Address) error {
	return s.repo.Update(address)
}

func (s *AddressService) DeleteAddress(id int) error {
	return s.repo.SoftDelete(id)
}
