package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

// CreateAddressDTO holds data for creating an address
type CreateAddressDTO struct {
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	ZipCode      string
	CustomerSub  string
}

// UpdateAddressDTO holds data for updating an address
type UpdateAddressDTO struct {
	AddressID    int
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	ZipCode      string
	CustomerSub  string
}

// DeleteAddressDTO holds data for deleting an address
type DeleteAddressDTO struct {
	AddressID   int
	CustomerSub string
}

type AddressService struct {
	addressRepo  *repository.AddressRepository
	customerRepo *repository.CustomerRepository
}

func NewAddressService(addressRepo *repository.AddressRepository, customerRepo *repository.CustomerRepository) *AddressService {
	return &AddressService{addressRepo: addressRepo, customerRepo: customerRepo}
}

func (s *AddressService) GetAddress(sub string, id int) (*model.Address, error) {
	return s.addressRepo.GetByIDAndSub(id, sub)
}

func (s *AddressService) GetAddresses(sub string) ([]model.Address, error) {
	return s.addressRepo.ListBySub(sub)
}

func (s *AddressService) Create(dto CreateAddressDTO) error {
	customer, err := s.customerRepo.GetBySub(dto.CustomerSub)
	if err != nil {
		return err
	}

	address := &model.Address{
		AddressLine1: dto.AddressLine1,
		AddressLine2: dto.AddressLine2,
		City:         dto.City,
		State:        dto.State,
		Country:      dto.Country,
		ZipCode:      dto.ZipCode,
		CustomerID:   customer.ID,
		IsActive:     true,
	}
	return s.addressRepo.Create(address)
}

func (s *AddressService) Update(dto UpdateAddressDTO) error {
	address, err := s.addressRepo.GetByIDAndSub(dto.AddressID, dto.CustomerSub)
	if err != nil {
		return err
	}

	address.AddressLine1 = dto.AddressLine1
	address.AddressLine2 = dto.AddressLine2
	address.City = dto.City
	address.State = dto.State
	address.Country = dto.Country
	address.ZipCode = dto.ZipCode
	return s.addressRepo.Update(address)
}

func (s *AddressService) Delete(dto DeleteAddressDTO) error {
	return s.addressRepo.SoftDelete(dto.AddressID, dto.CustomerSub)
}
