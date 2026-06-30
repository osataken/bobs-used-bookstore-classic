package repository

import (
	"math"
	"time"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

// OfferFilters for filtering offer queries
type OfferFilters struct {
	BookName    string
	Author      string
	GenreID     *int
	ConditionID *int
	OfferStatus *model.OfferStatus
}

type OfferRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) GetByID(id int) (*model.Offer, error) {
	var offer model.Offer
	err := r.db.Preload("Customer").Preload("Genre").Preload("Condition").
		Preload("Publisher").Preload("BookType").
		First(&offer, id).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepository) ListBySub(sub string) ([]model.Offer, error) {
	var offers []model.Offer
	err := r.db.Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType").
		Joins("JOIN \"Customer\" ON \"Customer\".id = \"Offer\".customer_id").
		Where("\"Customer\".sub = ?", sub).
		Order("\"Offer\".id DESC").
		Find(&offers).Error
	return offers, err
}

func (r *OfferRepository) List(filters OfferFilters, pageIndex, pageSize int) (*PaginatedResult[model.Offer], error) {
	query := r.db.Model(&model.Offer{}).Preload("Customer").Preload("Genre").Preload("Condition").Preload("Publisher").Preload("BookType")

	if filters.BookName != "" {
		query = query.Where("book_name LIKE ?", "%"+filters.BookName+"%")
	}
	if filters.Author != "" {
		query = query.Where("author LIKE ?", "%"+filters.Author+"%")
	}
	if filters.GenreID != nil {
		query = query.Where("genre_id = ?", *filters.GenreID)
	}
	if filters.ConditionID != nil {
		query = query.Where("condition_id = ?", *filters.ConditionID)
	}
	if filters.OfferStatus != nil {
		query = query.Where("offer_status = ?", *filters.OfferStatus)
	}

	var totalCount int64
	query.Count(&totalCount)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	offset := (pageIndex - 1) * pageSize

	var offers []model.Offer
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&offers).Error
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[model.Offer]{
		Items:      offers,
		PageIndex:  pageIndex,
		TotalPages: totalPages,
		TotalCount: int(totalCount),
	}, nil
}

func (r *OfferRepository) Create(offer *model.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) Update(offer *model.Offer) error {
	return r.db.Save(offer).Error
}

// OfferStatistics holds offer statistics
type OfferStatistics struct {
	PendingOffers   int `json:"pendingOffers"`
	OffersThisMonth int `json:"offersThisMonth"`
	OffersTotal     int `json:"offersTotal"`
}

func (r *OfferRepository) GetStatistics() (*OfferStatistics, error) {
	var stats OfferStatistics

	var pending int64
	r.db.Model(&model.Offer{}).Where("offer_status = ?", model.OfferStatusPendingApproval).Count(&pending)
	stats.PendingOffers = int(pending)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	var thisMonth int64
	r.db.Model(&model.Offer{}).Where("created_on >= ?", startOfMonth).Count(&thisMonth)
	stats.OffersThisMonth = int(thisMonth)

	var total int64
	r.db.Model(&model.Offer{}).Count(&total)
	stats.OffersTotal = int(total)

	return &stats, nil
}
