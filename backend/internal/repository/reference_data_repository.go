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

func (r *ReferenceDataRepository) GetByID(id int) (*model.ReferenceDataItem, error) {
	var item model.ReferenceDataItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ReferenceDataRepository) ListAll() ([]model.ReferenceDataItem, error) {
	var items []model.ReferenceDataItem
	err := r.db.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *ReferenceDataRepository) ListByType(dataType model.ReferenceDataType) ([]model.ReferenceDataItem, error) {
	var items []model.ReferenceDataItem
	err := r.db.Where("data_type = ?", dataType).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *ReferenceDataRepository) ListPaginated(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	var items []model.ReferenceDataItem
	query := r.db.Model(&model.ReferenceDataItem{})

	if v, ok := filters["dataType"]; ok {
		if dt, valid := v.(int); valid && dt >= 0 {
			query = query.Where("data_type = ?", dt)
		}
	}

	var totalCount int64
	query.Count(&totalCount)
	query = query.Order("id ASC")

	offset := (pageIndex - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&items).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return &model.PaginatedList{
		Items:      items,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (r *ReferenceDataRepository) Create(item *model.ReferenceDataItem) error {
	return r.db.Create(item).Error
}

func (r *ReferenceDataRepository) Update(item *model.ReferenceDataItem) error {
	return r.db.Save(item).Error
}
