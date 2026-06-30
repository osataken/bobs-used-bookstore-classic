package service

import (
	"time"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

// CreateOrUpdateCustomerDTO holds data for creating/updating a customer
type CreateOrUpdateCustomerDTO struct {
	Sub       string
	Username  string
	FirstName string
	LastName  string
}

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

func (s *CustomerService) CreateOrUpdate(dto CreateOrUpdateCustomerDTO) error {
	customer, err := s.customerRepo.GetBySub(dto.Sub)
	if err != nil {
		// Customer doesn't exist, create
		customer = &model.Customer{
			Sub:       dto.Sub,
			Username:  dto.Username,
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
		}
		return s.customerRepo.Create(customer)
	}

	// Update existing
	customer.Sub = dto.Sub
	customer.Username = dto.Username
	customer.FirstName = dto.FirstName
	customer.LastName = dto.LastName
	customer.UpdatedOn = time.Now().UTC()
	return s.customerRepo.Update(customer)
}
