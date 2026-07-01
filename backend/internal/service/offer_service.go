package service

import (
	"time"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type OfferService struct {
	offerRepo    *repository.OfferRepository
	customerRepo *repository.CustomerRepository
}

func NewOfferService(offerRepo *repository.OfferRepository, customerRepo *repository.CustomerRepository) *OfferService {
	return &OfferService{offerRepo: offerRepo, customerRepo: customerRepo}
}

func (s *OfferService) GetOffer(id uint) (*model.Offer, error) {
	return s.offerRepo.GetByID(id)
}

func (s *OfferService) GetOffersBySub(sub string) ([]model.Offer, error) {
	return s.offerRepo.ListBySub(sub)
}

func (s *OfferService) GetOffersFiltered(filters repository.OfferFilters, pageIndex int, pageSize int) (*repository.PaginatedResult[model.Offer], error) {
	return s.offerRepo.ListFiltered(filters, pageIndex, pageSize)
}

func (s *OfferService) GetStatistics() (*repository.OfferStatistics, error) {
	return s.offerRepo.GetStatistics()
}

type CreateOfferInput struct {
	BookName    string
	Author      string
	ISBN        string
	GenreID     uint
	ConditionID uint
	PublisherID uint
	BookTypeID  uint
	Summary     string
	BookPrice   float64
}

func (s *OfferService) CreateOffer(sub string, input CreateOfferInput) error {
	customer, err := s.customerRepo.GetBySub(sub)
	if err != nil {
		return err
	}

	offer := &model.Offer{
		CustomerID:  customer.ID,
		BookName:    input.BookName,
		Author:      input.Author,
		ISBN:        input.ISBN,
		GenreID:     input.GenreID,
		ConditionID: input.ConditionID,
		PublisherID: input.PublisherID,
		BookTypeID:  input.BookTypeID,
		Summary:     input.Summary,
		BookPrice:   input.BookPrice,
		OfferStatus: model.OfferStatusPendingApproval,
	}
	return s.offerRepo.Create(offer)
}

func (s *OfferService) UpdateOfferStatus(id uint, status model.OfferStatus) error {
	offer, err := s.offerRepo.GetByID(id)
	if err != nil {
		return err
	}
	offer.OfferStatus = status
	offer.UpdatedOn = time.Now().UTC()
	return s.offerRepo.Save(offer)
}
