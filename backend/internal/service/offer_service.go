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

func (s *OfferService) GetByID(id int) (*model.Offer, error) {
	return s.repo.GetByID(id)
}

func (s *OfferService) GetByCustomerID(customerID int) ([]model.Offer, error) {
	return s.repo.GetByCustomerID(customerID)
}

func (s *OfferService) GetAll(pageIndex, pageSize int, bookName, author string, genreID, conditionID int, offerStatus *model.OfferStatus) (repository.PaginatedList[model.Offer], error) {
	return s.repo.GetAll(pageIndex, pageSize, bookName, author, genreID, conditionID, offerStatus)
}

func (s *OfferService) Create(offer *model.Offer) error {
	offer.OfferStatus = model.OfferStatusPendingApproval
	return s.repo.Create(offer)
}

func (s *OfferService) Approve(id int) error {
	return s.repo.UpdateStatus(id, model.OfferStatusApproved)
}

func (s *OfferService) Reject(id int) error {
	return s.repo.UpdateStatus(id, model.OfferStatusRejected)
}

func (s *OfferService) MarkReceived(id int) error {
	return s.repo.UpdateStatus(id, model.OfferStatusReceived)
}

func (s *OfferService) MarkPaid(id int) error {
	return s.repo.UpdateStatus(id, model.OfferStatusPaid)
}

func (s *OfferService) GetPendingCount() (int64, error) {
	return s.repo.GetPendingCount()
}
