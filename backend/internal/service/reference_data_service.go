package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

// CreateReferenceDataDTO holds data for creating a reference data item
type CreateReferenceDataDTO struct {
	DataType model.ReferenceDataType
	Text     string
}

// UpdateReferenceDataDTO holds data for updating a reference data item
type UpdateReferenceDataDTO struct {
	ID       int
	DataType model.ReferenceDataType
	Text     string
}

type ReferenceDataService struct {
	refDataRepo *repository.ReferenceDataRepository
}

func NewReferenceDataService(refDataRepo *repository.ReferenceDataRepository) *ReferenceDataService {
	return &ReferenceDataService{refDataRepo: refDataRepo}
}

func (s *ReferenceDataService) GetAll() ([]model.ReferenceDataItem, error) {
	return s.refDataRepo.GetAll()
}

func (s *ReferenceDataService) GetByID(id int) (*model.ReferenceDataItem, error) {
	return s.refDataRepo.GetByID(id)
}

func (s *ReferenceDataService) List(filters repository.ReferenceDataFilters, pageIndex, pageSize int) (*repository.PaginatedResult[model.ReferenceDataItem], error) {
	return s.refDataRepo.List(filters, pageIndex, pageSize)
}

func (s *ReferenceDataService) Create(dto CreateReferenceDataDTO) error {
	item := &model.ReferenceDataItem{
		DataType: dto.DataType,
		Text:     dto.Text,
	}
	return s.refDataRepo.Create(item)
}

func (s *ReferenceDataService) Update(dto UpdateReferenceDataDTO) error {
	item, err := s.refDataRepo.GetByID(dto.ID)
	if err != nil {
		return err
	}
	item.DataType = dto.DataType
	item.Text = dto.Text
	return s.refDataRepo.Update(item)
}
