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

func (s *ReferenceDataService) GetByID(id int) (*model.ReferenceDataItem, error) {
	return s.repo.GetByID(id)
}

func (s *ReferenceDataService) GetByType(dataType model.ReferenceDataType) ([]model.ReferenceDataItem, error) {
	return s.repo.GetByType(dataType)
}

func (s *ReferenceDataService) GetAll(pageIndex, pageSize int, dataTypeFilter *model.ReferenceDataType) (repository.PaginatedList[model.ReferenceDataItem], error) {
	return s.repo.GetAll(pageIndex, pageSize, dataTypeFilter)
}

func (s *ReferenceDataService) Create(item *model.ReferenceDataItem) error {
	return s.repo.Create(item)
}

func (s *ReferenceDataService) Update(item *model.ReferenceDataItem) error {
	return s.repo.Update(item)
}
