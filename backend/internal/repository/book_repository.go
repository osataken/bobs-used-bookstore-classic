package repository

import (
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

func (r *BookRepository) GetAll(pageIndex, pageSize int, searchString, sortBy string) (PaginatedList[model.Book], error) {
	var books []model.Book
	var totalCount int64

	query := r.db.Model(&model.Book{})

	if searchString != "" {
		like := "%" + searchString + "%"
		query = query.Where("name LIKE ? OR author LIKE ? OR isbn LIKE ?", like, like, like)
	}

	query.Count(&totalCount)

	switch sortBy {
	case "Author":
		query = query.Order("author ASC")
	case "Price":
		query = query.Order("price ASC")
	default:
		query = query.Order("id ASC")
	}

	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return PaginatedList[model.Book]{}, err
	}

	return NewPaginatedList(books, int(totalCount), pageIndex, pageSize), nil
}

func (r *BookRepository) GetFiltered(name, author string, publisherID, genreID, bookTypeID, conditionID int, lowStock bool, pageIndex, pageSize int) (PaginatedList[model.Book], error) {
	var books []model.Book
	var totalCount int64

	query := r.db.Model(&model.Book{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}
	if publisherID > 0 {
		query = query.Where("publisher_id = ?", publisherID)
	}
	if genreID > 0 {
		query = query.Where("genre_id = ?", genreID)
	}
	if bookTypeID > 0 {
		query = query.Where("book_type_id = ?", bookTypeID)
	}
	if conditionID > 0 {
		query = query.Where("condition_id = ?", conditionID)
	}
	if lowStock {
		query = query.Where("quantity <= ?", model.LowBookThreshold)
	}

	query.Count(&totalCount)

	offset := (pageIndex - 1) * pageSize
	err := query.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Order("id ASC").Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		return PaginatedList[model.Book]{}, err
	}

	return NewPaginatedList(books, int(totalCount), pageIndex, pageSize), nil
}

func (r *BookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) GetBestSellers(limit int) ([]model.Book, error) {
	var books []model.Book
	err := r.db.Preload("Publisher").Preload("BookType").Preload("Genre").Preload("Condition").
		Where("quantity > 0").Order("id ASC").Limit(limit).Find(&books).Error
	return books, err
}

func (r *BookRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Book{}).Count(&count).Error
	return count, err
}

func (r *BookRepository) GetLowStockCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Book{}).Where("quantity <= ?", model.LowBookThreshold).Count(&count).Error
	return count, err
}
