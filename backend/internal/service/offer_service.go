package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"time"
)

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

func (s *OfferService) GetOffers(filters repository.OfferFilters, pageIndex, pageSize int) (*repository.PaginatedList, error) {
	return s.offerRepo.List(filters, pageIndex, pageSize)
}

func (s *OfferService) GetStatistics() (*repository.OfferStatistics, error) {
	return s.offerRepo.GetStatistics()
}

func (s *OfferService) CreateOffer(sub string, dto CreateOfferDTO) error {
	customer, err := s.customerRepo.GetBySub(sub)
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
		OfferStatus: model.OfferStatusPendingApproval,
	}

	return s.offerRepo.Add(offer)
}

func (s *OfferService) UpdateOfferStatus(offerID int, status model.OfferStatus) error {
	offer, err := s.offerRepo.GetByID(offerID)
	if err != nil {
		return err
	}

	offer.OfferStatus = status
	offer.UpdatedOn = time.Now().UTC()
	return s.offerRepo.Save(offer)
}

type CreateOfferDTO struct {
	BookName    string
	Author      string
	ISBN        string
	BookTypeID  int
	ConditionID int
	GenreID     int
	PublisherID int
	BookPrice   float64
}
