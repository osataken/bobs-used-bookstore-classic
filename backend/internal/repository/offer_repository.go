package repository

import (
	"math"

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
	err := r.db.Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").Preload("Customer").First(&offer, id).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepository) ListByCustomer(customerID int) ([]model.Offer, error) {
	var offers []model.Offer
	err := r.db.Where("customer_id = ?", customerID).
		Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").
		Order("id DESC").Find(&offers).Error
	return offers, err
}

func (r *OfferRepository) ListAll(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	var offers []model.Offer
	query := r.db.Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").Preload("Customer")

	if v, ok := filters["offerStatus"]; ok {
		if status, valid := v.(int); valid && status >= 0 {
			query = query.Where("offer_status = ?", status)
		}
	}

	var totalCount int64
	query.Model(&model.Offer{}).Count(&totalCount)
	query = query.Order("id ASC")

	offset := (pageIndex - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&offers).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return &model.PaginatedList{
		Items:      offers,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (r *OfferRepository) Create(offer *model.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) Update(offer *model.Offer) error {
	return r.db.Save(offer).Error
}

func (r *OfferRepository) GetStatistics() (*model.OfferStatistics, error) {
	var stats model.OfferStatistics
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusPendingApproval).Count(&stats.PendingApproval)
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusApproved).Count(&stats.Approved)
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusReceived).Count(&stats.Received)
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusPaid).Count(&stats.Paid)
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusRejected).Count(&stats.Rejected)
	return &stats, nil
}
