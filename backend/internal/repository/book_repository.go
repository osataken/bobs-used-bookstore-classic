package repository

import (
	"math"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetByID(id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	TotalCount int64 `json:"totalCount"`
	PageIndex  int   `json:"pageIndex"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

func (r *BookRepository) List(searchString string, sortBy string, pageIndex int, pageSize int) (*PaginatedResult[model.Book], error) {
	query := r.db.Model(&model.Book{})

	if searchString != "" {
		search := "%" + searchString + "%"
		query = query.Where("name LIKE ? OR author LIKE ? OR isbn LIKE ?", search, search, search)
	}

	switch sortBy {
	case "Author":
		query = query.Order("author ASC")
	case "Price":
		query = query.Order("price ASC")
	default:
		query = query.Order("name ASC")
	}

	var totalCount int64
	query.Count(&totalCount)

	var books []model.Book
	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedResult[model.Book]{
		Items:      books,
		TotalCount: totalCount,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

type BookFilters struct {
	SearchString string
	PublisherID  *uint
	BookTypeID   *uint
	GenreID      *uint
	ConditionID  *uint
	InStock      *bool
}

func (r *BookRepository) ListFiltered(filters BookFilters, pageIndex int, pageSize int) (*PaginatedResult[model.Book], error) {
	query := r.db.Model(&model.Book{}).Order("id ASC")

	if filters.SearchString != "" {
		search := "%" + filters.SearchString + "%"
		query = query.Where("name LIKE ? OR author LIKE ?", search, search)
	}
	if filters.PublisherID != nil {
		query = query.Where("publisher_id = ?", *filters.PublisherID)
	}
	if filters.BookTypeID != nil {
		query = query.Where("book_type_id = ?", *filters.BookTypeID)
	}
	if filters.GenreID != nil {
		query = query.Where("genre_id = ?", *filters.GenreID)
	}
	if filters.ConditionID != nil {
		query = query.Where("condition_id = ?", *filters.ConditionID)
	}
	if filters.InStock != nil {
		if *filters.InStock {
			query = query.Where("quantity > 0")
		} else {
			query = query.Where("quantity = 0")
		}
	}

	var totalCount int64
	query.Count(&totalCount)

	var books []model.Book
	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedResult[model.Book]{
		Items:      books,
		TotalCount: totalCount,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

type BookStatistics struct {
	LowStock   int64
	OutOfStock int64
	StockTotal int64
}

func (r *BookRepository) GetStatistics() (*BookStatistics, error) {
	var stats BookStatistics
	r.db.Model(&model.Book{}).Where("quantity <= ? AND quantity > 0", model.LowBookThreshold).Count(&stats.LowStock)
	r.db.Model(&model.Book{}).Where("quantity = 0").Count(&stats.OutOfStock)
	r.db.Model(&model.Book{}).Count(&stats.StockTotal)
	return &stats, nil
}

func (r *BookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) DB() *gorm.DB {
	return r.db
}
