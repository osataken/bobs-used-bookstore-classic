package repository

import (
	"bobs-used-bookstore-api/internal/model"
	"math"
	"time"

	"gorm.io/gorm"
)

type OfferFilters struct {
	OfferStatus *int `form:"offerStatus"`
}

type OfferStatistics struct {
	PendingOffers  int64 `json:"pendingOffers"`
	OffersThisMonth int64 `json:"offersThisMonth"`
	OffersTotal    int64 `json:"offersTotal"`
}

type OfferRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) GetByID(id int) (*model.Offer, error) {
	var offer model.Offer
	err := r.db.Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").Preload("Customer").
		First(&offer, id).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepository) ListBySub(sub string) ([]model.Offer, error) {
	var offers []model.Offer
	err := r.db.Joins("JOIN customer ON customer.id = offer.customer_id").
		Where("customer.sub = ?", sub).
		Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").
		Order("offer.id ASC").
		Find(&offers).Error
	return offers, err
}


func (r *OfferRepository) List(filters OfferFilters, pageIndex, pageSize int) (*PaginatedList, error) {
	query := r.db.Model(&model.Offer{}).
		Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").Preload("Customer")

	if filters.OfferStatus != nil {
		query = query.Where("offer_status = ?", *filters.OfferStatus)
	}

	var totalCount int64
	query.Count(&totalCount)

	var offers []model.Offer
	offset := (pageIndex - 1) * pageSize
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&offers).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedList{
		Items:      offers,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}, nil
}

func (r *OfferRepository) Add(offer *model.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) Save(offer *model.Offer) error {
	return r.db.Save(offer).Error
}

func (r *OfferRepository) GetStatistics() (*OfferStatistics, error) {
	var stats OfferStatistics
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusPendingApproval).Count(&stats.PendingOffers)
	r.db.Model(&model.Offer{}).Where("created_on >= ?", startOfMonth).Count(&stats.OffersThisMonth)
	r.db.Model(&model.Offer{}).Count(&stats.OffersTotal)

	return &stats, nil
}
