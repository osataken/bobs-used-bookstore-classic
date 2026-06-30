package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

// CreateOfferDTO holds data for creating an offer
type CreateOfferDTO struct {
	CustomerSub string
	BookName    string
	Author      string
	ISBN        string
	BookTypeID  int
	ConditionID int
	GenreID     int
	PublisherID int
	BookPrice   float64
	Summary     string
}

// UpdateOfferStatusDTO holds data for updating offer status
type UpdateOfferStatusDTO struct {
	OfferID int
	Status  model.OfferStatus
}

type OfferService struct {
	offerRepo    *repository.OfferRepository
	customerRepo *repository.CustomerRepository
}

func NewOfferService(offerRepo *repository.OfferRepository, customerRepo *repository.CustomerRepository) *OfferService {
	return &OfferService{offerRepo: offerRepo, customerRepo: customerRepo}
}

func (s *OfferService) GetOffer(id int) (*model.Offer, error) {
	return s.offerRepo.GetByID(id)
}

func (s *OfferService) GetOffersByCustomer(sub string) ([]model.Offer, error) {
	return s.offerRepo.ListBySub(sub)
}

func (s *OfferService) GetOffers(filters repository.OfferFilters, pageIndex, pageSize int) (*repository.PaginatedResult[model.Offer], error) {
	return s.offerRepo.List(filters, pageIndex, pageSize)
}

func (s *OfferService) GetStatistics() (*repository.OfferStatistics, error) {
	return s.offerRepo.GetStatistics()
}

func (s *OfferService) Create(dto CreateOfferDTO) error {
	customer, err := s.customerRepo.GetBySub(dto.CustomerSub)
	if err != nil {
		return err
	}

	offer := &model.Offer{
		CustomerID:  customer.ID,
		BookName:    dto.BookName,
		Author:      dto.Author,
		ISBN:        dto.ISBN,
		BookTypeID:  dto.BookTypeID,
		ConditionID: dto.ConditionID,
		GenreID:     dto.GenreID,
		PublisherID: dto.PublisherID,
		BookPrice:   dto.BookPrice,
		Summary:     dto.Summary,
		OfferStatus: model.OfferStatusPendingApproval,
	}
	return s.offerRepo.Create(offer)
}

func (s *OfferService) UpdateStatus(dto UpdateOfferStatusDTO) error {
	offer, err := s.offerRepo.GetByID(dto.OfferID)
	if err != nil {
		return err
	}
	offer.OfferStatus = dto.Status
	return s.offerRepo.Update(offer)
}
