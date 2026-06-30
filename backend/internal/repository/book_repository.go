package repository

import (
	"math"

	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

// BookFilters for filtering book queries
type BookFilters struct {
	Name        string
	Author      string
	PublisherID *int
	GenreID     *int
	BookTypeID  *int
	ConditionID *int
	LowStock    bool
}

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetByID(id int) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("Genre").Preload("Publisher").Preload("BookType").Preload("Condition").
		First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) List(filters BookFilters, pageIndex, pageSize int) (*PaginatedResult[model.Book], error) {
	query := r.db.Model(&model.Book{}).Preload("Genre").Preload("Publisher").Preload("BookType").Preload("Condition")

	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.Author != "" {
		query = query.Where("author LIKE ?", "%"+filters.Author+"%")
	}
	if filters.PublisherID != nil {
		query = query.Where("publisher_id = ?", *filters.PublisherID)
	}
	if filters.GenreID != nil {
		query = query.Where("genre_id = ?", *filters.GenreID)
	}
	if filters.BookTypeID != nil {
		query = query.Where("book_type_id = ?", *filters.BookTypeID)
	}
	if filters.ConditionID != nil {
		query = query.Where("condition_id = ?", *filters.ConditionID)
	}
	if filters.LowStock {
		query = query.Where("quantity > 0 AND quantity <= ?", model.LowBookThreshold)
	}

	var totalCount int64
	query.Count(&totalCount)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	offset := (pageIndex - 1) * pageSize

	var books []model.Book
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[model.Book]{
		Items:      books,
		PageIndex:  pageIndex,
		TotalPages: totalPages,
		TotalCount: int(totalCount),
	}, nil
}

func (r *BookRepository) Search(searchString, sortBy string, pageIndex, pageSize int) (*PaginatedResult[model.Book], error) {
	query := r.db.Model(&model.Book{}).
		Preload("Genre").Preload("Publisher").Preload("BookType").Preload("Condition").
		Joins("LEFT JOIN \"ReferenceData\" AS genre ON genre.id = \"Book\".genre_id").
		Joins("LEFT JOIN \"ReferenceData\" AS book_type ON book_type.id = \"Book\".book_type_id").
		Joins("LEFT JOIN \"ReferenceData\" AS publisher ON publisher.id = \"Book\".publisher_id")

	if searchString != "" {
		like := "%" + searchString + "%"
		query = query.Where(
			"\"Book\".name LIKE ? OR genre.text LIKE ? OR book_type.text LIKE ? OR \"Book\".isbn LIKE ? OR publisher.text LIKE ?",
			like, like, like, like, like,
		)
	}

	switch sortBy {
	case "PriceAsc":
		query = query.Order("\"Book\".price ASC")
	case "PriceDesc":
		query = query.Order("\"Book\".price DESC")
	default:
		query = query.Order("\"Book\".name ASC")
	}

	var totalCount int64
	query.Count(&totalCount)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	offset := (pageIndex - 1) * pageSize

	var books []model.Book
	err := query.Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[model.Book]{
		Items:      books,
		PageIndex:  pageIndex,
		TotalPages: totalPages,
		TotalCount: int(totalCount),
	}, nil
}

func (r *BookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Save() error {
	return nil // GORM saves immediately
}

// BookStatistics holds inventory statistics
type BookStatistics struct {
	OutOfStock int `json:"outOfStock"`
	LowStock   int `json:"lowStock"`
	StockTotal int `json:"stockTotal"`
}

func (r *BookRepository) GetStatistics() (*BookStatistics, error) {
	var stats BookStatistics
	r.db.Model(&model.Book{}).Where("quantity = 0").Count(new(int64))

	var outOfStock int64
	r.db.Model(&model.Book{}).Where("quantity = 0").Count(&outOfStock)

	var lowStock int64
	r.db.Model(&model.Book{}).Where("quantity > 0 AND quantity <= ?", model.LowBookThreshold).Count(&lowStock)

	var stockTotal int64
	r.db.Model(&model.Book{}).Count(&stockTotal)

	stats.OutOfStock = int(outOfStock)
	stats.LowStock = int(lowStock)
	stats.StockTotal = int(stockTotal)

	return &stats, nil
}
