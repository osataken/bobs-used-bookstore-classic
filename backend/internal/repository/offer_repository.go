package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type OfferRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) GetByID(id int) (*model.Offer, error) {
	var offer model.Offer
	err := r.db.Preload("Customer").Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").First(&offer, id).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepository) GetByCustomerID(customerID int) ([]model.Offer, error) {
	var offers []model.Offer
	err := r.db.Where("customer_id = ?", customerID).Order("id DESC").Find(&offers).Error
	return offers, err
}

func (r *OfferRepository) GetAll(pageIndex, pageSize int, bookName, author string, genreID, conditionID int, offerStatus *model.OfferStatus) (PaginatedList[model.Offer], error) {
	var offers []model.Offer
	var totalCount int64

	query := r.db.Model(&model.Offer{})

	if bookName != "" {
		query = query.Where("book_name LIKE ?", "%"+bookName+"%")
	}
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}
	if genreID > 0 {
		query = query.Where("genre_id = ?", genreID)
	}
	if conditionID > 0 {
		query = query.Where("condition_id = ?", conditionID)
	}
	if offerStatus != nil {
		query = query.Where("offer_status = ?", *offerStatus)
	}

	query.Count(&totalCount)

	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Customer").Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").
		Order("id ASC").Offset(offset).Limit(pageSize).Find(&offers).Error
	if err != nil {
		return PaginatedList[model.Offer]{}, err
	}

	return NewPaginatedList(offers, int(totalCount), pageIndex, pageSize), nil
}

func (r *OfferRepository) Create(offer *model.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) UpdateStatus(id int, status model.OfferStatus) error {
	return r.db.Model(&model.Offer{}).Where("id = ?", id).Update("offer_status", status).Error
}

func (r *OfferRepository) GetPendingCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusPendingApproval).Count(&count).Error
	return count, err
}
