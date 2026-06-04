package repository

import (
	"bobs-used-bookstore-api/internal/model"
	"math"
	"strings"

	"gorm.io/gorm"
)

type PaginatedList struct {
	Items      interface{} `json:"items"`
	TotalPages int         `json:"totalPages"`
	PageIndex  int         `json:"pageIndex"`
	PageSize   int         `json:"pageSize"`
	TotalCount int64       `json:"totalCount"`
}

type BookFilters struct {
	SearchString string `form:"searchString"`
	GenreID      int    `form:"genreId"`
	BookTypeID   int    `form:"bookTypeId"`
	ConditionID  int    `form:"conditionId"`
	PublisherID  int    `form:"publisherId"`
}

type BookStatistics struct {
	LowStock   int64 `json:"lowStock"`
	OutOfStock int64 `json:"outOfStock"`
	StockTotal int64 `json:"stockTotal"`
}

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetByID(id int) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) List(filters BookFilters, pageIndex, pageSize int) (*PaginatedList, error) {
	query := r.db.Model(&model.Book{}).Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition")

	if filters.SearchString != "" {
		search := "%" + strings.ToLower(filters.SearchString) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(author) LIKE ? OR LOWER(isbn) LIKE ?", search, search, search)
	}
	if filters.GenreID > 0 {
		query = query.Where("genre_id = ?", filters.GenreID)
	}
	if filters.BookTypeID > 0 {
		query = query.Where("book_type_id = ?", filters.BookTypeID)
	}
	if filters.ConditionID > 0 {
		query = query.Where("condition_id = ?", filters.ConditionID)
	}
	if filters.PublisherID > 0 {
		query = query.Where("publisher_id = ?", filters.PublisherID)
	}

	var totalCount int64
	query.Count(&totalCount)

	var books []model.Book
	offset := (pageIndex - 1) * pageSize
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedList{
		Items:      books,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}, nil
}

func (r *BookRepository) Search(searchString, sortBy string, pageIndex, pageSize int) (*PaginatedList, error) {
	query := r.db.Model(&model.Book{}).Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition")

	if searchString != "" {
		search := "%" + strings.ToLower(searchString) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(author) LIKE ? OR LOWER(isbn) LIKE ?", search, search, search)
	}

	var totalCount int64
	query.Count(&totalCount)

	switch strings.ToLower(sortBy) {
	case "price":
		query = query.Order("price ASC")
	case "author":
		query = query.Order("author ASC")
	default:
		query = query.Order("name ASC")
	}

	var books []model.Book
	offset := (pageIndex - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedList{
		Items:      books,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}, nil
}

func (r *BookRepository) Add(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) GetStatistics() (*BookStatistics, error) {
	var stats BookStatistics
	r.db.Model(&model.Book{}).Where("quantity = 0").Count(&stats.OutOfStock)
	r.db.Model(&model.Book{}).Where("quantity > 0 AND quantity <= ?", model.LowBookThreshold).Count(&stats.LowStock)
	r.db.Model(&model.Book{}).Count(&stats.StockTotal)
	return &stats, nil
}

func (r *BookRepository) ListBestSelling(count int) ([]model.Book, error) {
	var bookIDs []int
	err := r.db.Model(&model.OrderItem{}).
		Select("book_id").
		Group("book_id").
		Order("COUNT(*) DESC").
		Limit(count).
		Pluck("book_id", &bookIDs).Error
	if err != nil {
		return nil, err
	}

	if len(bookIDs) == 0 {
		var books []model.Book
		err = r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
			Order("id ASC").Limit(count).Find(&books).Error
		return books, err
	}

	var books []model.Book
	err = r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Where("id IN ?", bookIDs).Find(&books).Error
	return books, err
}
