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

func (s *AddressService) GetAddress(sub string, id uint) (*model.Address, error) {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return nil, err
	}
	return s.addressRepo.GetByIDAndCustomerID(id, customer.ID)
}

func (s *AddressService) GetAddresses(sub string) ([]model.Address, error) {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return nil, err
	}
	return s.addressRepo.ListByCustomerID(customer.ID)
}

type CreateAddressInput struct {
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	ZipCode      string
}

func (s *AddressService) CreateAddress(sub string, input CreateAddressInput) error {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return err
	}

	address := &model.Address{
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		Country:      input.Country,
		ZipCode:      input.ZipCode,
		CustomerID:   customer.ID,
		IsActive:     true,
	}
	return s.addressRepo.Create(address)
}

type UpdateAddressInput struct {
	ID           uint
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	ZipCode      string
}

func (s *AddressService) UpdateAddress(sub string, input UpdateAddressInput) error {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return err
	}

	address, err := s.addressRepo.GetByIDAndCustomerID(input.ID, customer.ID)
	if err != nil {
		return err
	}

	address.AddressLine1 = input.AddressLine1
	address.AddressLine2 = input.AddressLine2
	address.City = input.City
	address.State = input.State
	address.Country = input.Country
	address.ZipCode = input.ZipCode
	return s.addressRepo.Update(address)
}

func (s *AddressService) DeleteAddress(sub string, id uint) error {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return err
	}
	return s.addressRepo.SoftDelete(id, customer.ID)
}
