package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

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

func (s *ReferenceDataService) List(filters repository.ReferenceDataFilters, pageIndex, pageSize int) (*repository.PaginatedList, error) {
	return s.refDataRepo.List(filters, pageIndex, pageSize)
}

func (s *ReferenceDataService) Create(dataType model.ReferenceDataType, text string) error {
	item := &model.ReferenceDataItem{
		DataType: dataType,
		Text:     text,
	}
	return s.refDataRepo.Add(item)
}

func (s *ReferenceDataService) Update(id int, dataType model.ReferenceDataType, text string) error {
	item, err := s.refDataRepo.GetByID(id)
	if err != nil {
		return err
	}

	item.DataType = dataType
	item.Text = text
	return s.refDataRepo.Save(item)
}
