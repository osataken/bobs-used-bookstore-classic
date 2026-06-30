package repository

import "gorm.io/gorm"

// Repositories holds all repository instances
type Repositories struct {
	Book          *BookRepository
	Order         *OrderRepository
	Customer      *CustomerRepository
	Address       *AddressRepository
	Offer         *OfferRepository
	ShoppingCart  *ShoppingCartRepository
	ReferenceData *ReferenceDataRepository
}

// NewRepositories creates all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Book:          NewBookRepository(db),
		Order:         NewOrderRepository(db),
		Customer:      NewCustomerRepository(db),
		Address:       NewAddressRepository(db),
		Offer:         NewOfferRepository(db),
		ShoppingCart:  NewShoppingCartRepository(db),
		ReferenceData: NewReferenceDataRepository(db),
	}
}

// PaginatedResult holds paginated query results
type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	PageIndex  int `json:"pageIndex"`
	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"`
}

func (p *PaginatedResult[T]) HasNextPage() bool {
	return p.PageIndex < p.TotalPages
}

func (p *PaginatedResult[T]) HasPreviousPage() bool {
	return p.PageIndex > 1
}
