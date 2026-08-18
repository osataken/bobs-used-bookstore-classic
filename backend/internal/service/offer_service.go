package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type OfferService struct {
	repo *repository.OfferRepository
}

func NewOfferService(repo *repository.OfferRepository) *OfferService {
	return &OfferService{repo: repo}
}

func (s *OfferService) GetOffer(id int) (*model.Offer, error) {
	return s.repo.GetByID(id)
}

func (s *OfferService) ListByCustomer(customerID int) ([]model.Offer, error) {
	return s.repo.ListByCustomer(customerID)
}

func (s *OfferService) ListAll(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListAll(filters, pageIndex, pageSize)
}

func (s *OfferService) CreateOffer(offer *model.Offer) error {
	offer.OfferStatus = model.OfferStatusPendingApproval
	return s.repo.Create(offer)
}

func (s *OfferService) UpdateStatus(id int, status model.OfferStatus) error {
	offer, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	offer.OfferStatus = status
	return s.repo.Update(offer)
}

func (s *OfferService) GetStatistics() (*model.OfferStatistics, error) {
	return s.repo.GetStatistics()
}
