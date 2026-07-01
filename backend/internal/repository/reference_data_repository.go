package repository

import (
	"math"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type ReferenceDataRepository struct {
	db *gorm.DB
}

func NewReferenceDataRepository(db *gorm.DB) *ReferenceDataRepository {
	return &ReferenceDataRepository{db: db}
}

func (r *ReferenceDataRepository) GetByID(id uint) (*model.ReferenceDataItem, error) {
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

type ReferenceDataFilters struct {
	DataType *model.ReferenceDataType
}

func (r *ReferenceDataRepository) ListFiltered(filters ReferenceDataFilters, pageIndex int, pageSize int) (*PaginatedResult[model.ReferenceDataItem], error) {
	query := r.db.Model(&model.ReferenceDataItem{}).Order("id ASC")

	if filters.DataType != nil {
		query = query.Where("data_type = ?", *filters.DataType)
	}

	var totalCount int64
	query.Count(&totalCount)

	var items []model.ReferenceDataItem
	offset := (pageIndex - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedResult[model.ReferenceDataItem]{
		Items:      items,
		TotalCount: totalCount,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *ReferenceDataRepository) Create(item *model.ReferenceDataItem) error {
	return r.db.Create(item).Error
}

func (r *ReferenceDataRepository) Update(item *model.ReferenceDataItem) error {
	return r.db.Save(item).Error
}
