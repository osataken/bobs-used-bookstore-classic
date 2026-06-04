package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type AddressService struct {
	addressRepo  *repository.AddressRepository
	customerRepo *repository.CustomerRepository
}

func NewAddressService(addressRepo *repository.AddressRepository, customerRepo *repository.CustomerRepository) *AddressService {
	return &AddressService{addressRepo: addressRepo, customerRepo: customerRepo}
}

func (s *AddressService) GetAddress(sub string, id int) (*model.Address, error) {
	return s.addressRepo.GetBySubAndID(sub, id)
}

func (s *AddressService) GetAddresses(sub string) ([]model.Address, error) {
	return s.addressRepo.ListBySub(sub)
}

func (s *AddressService) Create(sub string, dto CreateAddressDTO) error {
	customer, err := s.customerRepo.GetBySub(sub)
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

	return s.addressRepo.Add(address)
}

func (s *AddressService) Update(sub string, dto UpdateAddressDTO) error {
	address, err := s.addressRepo.GetBySubAndID(sub, dto.ID)
	if err != nil {
		return err
	}
	if address == nil {
		return nil
	}

	address.AddressLine1 = dto.AddressLine1
	address.AddressLine2 = dto.AddressLine2
	address.City = dto.City
	address.State = dto.State
	address.Country = dto.Country
	address.ZipCode = dto.ZipCode

	return s.addressRepo.Save(address)
}

func (s *AddressService) Delete(sub string, id int) error {
	return s.addressRepo.SoftDelete(sub, id)
}

type CreateAddressDTO struct {
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	ZipCode      string
}

type UpdateAddressDTO struct {
	ID           int
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	ZipCode      string
}
