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

func (r *OfferRepository) GetByID(id uint) (*model.Offer, error) {
	var offer model.Offer
	err := r.db.Preload("Customer").Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		First(&offer, id).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepository) ListBySub(sub string) ([]model.Offer, error) {
	var offers []model.Offer
	err := r.db.Joins("Customer").Where("\"Customer\".sub = ?", sub).
		Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Order("\"offers\".id ASC").Find(&offers).Error
	return offers, err
}

type OfferFilters struct {
	Status *model.OfferStatus
}

func (r *OfferRepository) ListFiltered(filters OfferFilters, pageIndex int, pageSize int) (*PaginatedResult[model.Offer], error) {
	query := r.db.Model(&model.Offer{}).Order("id ASC")

	if filters.Status != nil {
		query = query.Where("offer_status = ?", *filters.Status)
	}

	var totalCount int64
	query.Count(&totalCount)

	var offers []model.Offer
	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Customer").Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Offset(offset).Limit(pageSize).Find(&offers).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedResult[model.Offer]{
		Items:      offers,
		TotalCount: totalCount,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

type OfferStatistics struct {
	PendingOffers   int64
	OffersThisMonth int64
	OffersTotal     int64
}

func (r *OfferRepository) GetStatistics() (*OfferStatistics, error) {
	var stats OfferStatistics
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusPendingApproval).Count(&stats.PendingOffers)
	r.db.Model(&model.Offer{}).Count(&stats.OffersTotal)
	// Simplified: count all as this month for now
	r.db.Model(&model.Offer{}).Count(&stats.OffersThisMonth)
	return &stats, nil
}

func (r *OfferRepository) Create(offer *model.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) Save(offer *model.Offer) error {
	return r.db.Save(offer).Error
}
