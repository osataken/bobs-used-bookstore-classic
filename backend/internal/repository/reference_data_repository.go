package repository

import (
	"bobs-used-bookstore-api/internal/model"
	"math"

	"gorm.io/gorm"
)

type ReferenceDataFilters struct {
	DataType *int `form:"dataType"`
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

func (r *ReferenceDataRepository) List(filters ReferenceDataFilters, pageIndex, pageSize int) (*PaginatedList, error) {
	query := r.db.Model(&model.ReferenceDataItem{})

	if filters.DataType != nil {
		query = query.Where("data_type = ?", *filters.DataType)
	}

	var totalCount int64
	query.Count(&totalCount)

	var items []model.ReferenceDataItem
	offset := (pageIndex - 1) * pageSize
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedList{
		Items:      items,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}, nil
}

func (r *ReferenceDataRepository) Add(item *model.ReferenceDataItem) error {
	return r.db.Create(item).Error
}

func (r *ReferenceDataRepository) Save(item *model.ReferenceDataItem) error {
	return r.db.Save(item).Error
}
