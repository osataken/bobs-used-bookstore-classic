package repository

import (
	"math"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

// ReferenceDataFilters for filtering reference data queries
type ReferenceDataFilters struct {
	DataType *model.ReferenceDataType
}

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

func (r *ReferenceDataRepository) GetAll() ([]model.ReferenceDataItem, error) {
	var items []model.ReferenceDataItem
	err := r.db.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *ReferenceDataRepository) List(filters ReferenceDataFilters, pageIndex, pageSize int) (*PaginatedResult[model.ReferenceDataItem], error) {
	query := r.db.Model(&model.ReferenceDataItem{})

	if filters.DataType != nil {
		query = query.Where("data_type = ?", *filters.DataType)
	}

	var totalCount int64
	query.Count(&totalCount)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	offset := (pageIndex - 1) * pageSize

	var items []model.ReferenceDataItem
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[model.ReferenceDataItem]{
		Items:      items,
		PageIndex:  pageIndex,
		TotalPages: totalPages,
		TotalCount: int(totalCount),
	}, nil
}

func (r *ReferenceDataRepository) Create(item *model.ReferenceDataItem) error {
	return r.db.Create(item).Error
}

func (r *ReferenceDataRepository) Update(item *model.ReferenceDataItem) error {
	return r.db.Save(item).Error
}
