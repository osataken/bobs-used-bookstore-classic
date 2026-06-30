package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type ReferenceDataRepository struct {
	db *gorm.DB
}

func NewReferenceDataRepository(db *gorm.DB) *ReferenceDataRepository {
	return &ReferenceDataRepository{db: db}
}

func (r *ReferenceDataRepository) GetByID(id int) (*model.ReferenceDataItem, error) {
	var item model.ReferenceDataItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ReferenceDataRepository) GetByType(dataType model.ReferenceDataType) ([]model.ReferenceDataItem, error) {
	var items []model.ReferenceDataItem
	err := r.db.Where("data_type = ?", dataType).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *ReferenceDataRepository) GetAll(pageIndex, pageSize int, dataTypeFilter *model.ReferenceDataType) (PaginatedList[model.ReferenceDataItem], error) {
	var items []model.ReferenceDataItem
	var totalCount int64

	query := r.db.Model(&model.ReferenceDataItem{})

	if dataTypeFilter != nil {
		query = query.Where("data_type = ?", *dataTypeFilter)
	}

	query.Count(&totalCount)

	offset := (pageIndex - 1) * pageSize
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return PaginatedList[model.ReferenceDataItem]{}, err
	}

	return NewPaginatedList(items, int(totalCount), pageIndex, pageSize), nil
}

func (r *ReferenceDataRepository) Create(item *model.ReferenceDataItem) error {
	return r.db.Create(item).Error
}

func (r *ReferenceDataRepository) Update(item *model.ReferenceDataItem) error {
	return r.db.Save(item).Error
}
