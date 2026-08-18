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

func (r *BookRepository) GetByID(id int) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) List(searchString string, sortBy string, pageIndex, pageSize int) (*model.PaginatedList, error) {
	var books []model.Book
	query := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition")

	if searchString != "" {
		search := "%" + searchString + "%"
		query = query.Where("Name LIKE ? OR Author LIKE ? OR ISBN LIKE ?", search, search, search)
	}

	var totalCount int64
	query.Model(&model.Book{}).Count(&totalCount)

	switch sortBy {
	case "price_asc":
		query = query.Order("Price ASC")
	case "price_desc":
		query = query.Order("Price DESC")
	case "name_asc":
		query = query.Order("Name ASC")
	case "name_desc":
		query = query.Order("Name DESC")
	default:
		query = query.Order("id ASC")
	}

	offset := (pageIndex - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&books).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &model.PaginatedList{
		Items:      books,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (r *BookRepository) ListWithFilters(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	var books []model.Book
	query := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition")

	if v, ok := filters["searchString"]; ok && v != "" {
		search := "%" + v.(string) + "%"
		query = query.Where("Name LIKE ? OR Author LIKE ? OR ISBN LIKE ?", search, search, search)
	}
	if v, ok := filters["genreId"]; ok && v.(int) > 0 {
		query = query.Where("genre_id = ?", v)
	}
	if v, ok := filters["bookTypeId"]; ok && v.(int) > 0 {
		query = query.Where("book_type_id = ?", v)
	}
	if v, ok := filters["conditionId"]; ok && v.(int) > 0 {
		query = query.Where("condition_id = ?", v)
	}
	if v, ok := filters["publisherId"]; ok && v.(int) > 0 {
		query = query.Where("publisher_id = ?", v)
	}

	var totalCount int64
	query.Model(&model.Book{}).Count(&totalCount)
	query = query.Order("id ASC")

	offset := (pageIndex - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&books).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &model.PaginatedList{
		Items:      books,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (r *BookRepository) ListBestSelling(limit int) ([]model.Book, error) {
	var books []model.Book
	// Best selling = most ordered. Subquery on order_items.
	err := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Order("id ASC").Limit(limit).Find(&books).Error
	return books, err
}

func (r *BookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) GetStatistics() (*model.BookStatistics, error) {
	var stats model.BookStatistics
	r.db.Model(&model.Book{}).Where("quantity = 0").Count(&stats.OutOfStock)
	r.db.Model(&model.Book{}).Where("quantity > 0 AND quantity <= ?", model.LowBookThreshold).Count(&stats.LowStock)
	r.db.Model(&model.Book{}).Count(&stats.StockTotal)
	return &stats, nil
}
