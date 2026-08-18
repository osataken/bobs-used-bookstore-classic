package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type ReferenceDataService struct {
	repo *repository.ReferenceDataRepository
}

func NewReferenceDataService(repo *repository.ReferenceDataRepository) *ReferenceDataService {
	return &ReferenceDataService{repo: repo}
}

func (s *ReferenceDataService) GetItem(id int) (*model.ReferenceDataItem, error) {
	return s.repo.GetByID(id)
}

func (s *ReferenceDataService) ListAll() ([]model.ReferenceDataItem, error) {
	return s.repo.ListAll()
}

func (s *ReferenceDataService) ListByType(dataType model.ReferenceDataType) ([]model.ReferenceDataItem, error) {
	return s.repo.ListByType(dataType)
}

func (s *ReferenceDataService) ListPaginated(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListPaginated(filters, pageIndex, pageSize)
}

func (s *ReferenceDataService) Create(item *model.ReferenceDataItem) error {
	return s.repo.Create(item)
}

func (s *ReferenceDataService) Update(item *model.ReferenceDataItem) error {
	return s.repo.Update(item)
}
