package service

import (
	"io"
	"mime/multipart"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

// CreateBookDTO holds data for creating a book
type CreateBookDTO struct {
	Name          string
	Author        string
	BookTypeID    int
	ConditionID   int
	GenreID       int
	PublisherID   int
	Year          *int
	ISBN          string
	Summary       string
	Price         float64
	Quantity      int
	CoverImage    multipart.File
	CoverFileName string
}

// UpdateBookDTO holds data for updating a book
type UpdateBookDTO struct {
	BookID        int
	Name          string
	Author        string
	BookTypeID    int
	ConditionID   int
	GenreID       int
	PublisherID   int
	Year          *int
	ISBN          string
	Summary       string
	Price         float64
	Quantity      int
	CoverImage    multipart.File
	CoverFileName string
}

// BookResult holds the result of a book operation
type BookResult struct {
	Success bool
	Error   string
}

type BookService struct {
	bookRepo        *repository.BookRepository
	orderRepo       *repository.OrderRepository
	imageResize     *ImageResizeService
	imageValidation ImageValidationService
	fileService     FileService
}

func NewBookService(
	bookRepo *repository.BookRepository,
	orderRepo *repository.OrderRepository,
	imageResize *ImageResizeService,
	imageValidation ImageValidationService,
	fileService FileService,
) *BookService {
	return &BookService{
		bookRepo:        bookRepo,
		orderRepo:       orderRepo,
		imageResize:     imageResize,
		imageValidation: imageValidation,
		fileService:     fileService,
	}
}

func (s *BookService) GetBook(id int) (*model.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookService) GetBooks(filters repository.BookFilters, pageIndex, pageSize int) (*repository.PaginatedResult[model.Book], error) {
	return s.bookRepo.List(filters, pageIndex, pageSize)
}

func (s *BookService) SearchBooks(searchString, sortBy string, pageIndex, pageSize int) (*repository.PaginatedResult[model.Book], error) {
	return s.bookRepo.Search(searchString, sortBy, pageIndex, pageSize)
}

func (s *BookService) ListBestSellingBooks(count int) ([]model.Book, error) {
	bookIDs, err := s.orderRepo.ListBestSellingBookIDs(count)
	if err != nil {
		return nil, err
	}

	if len(bookIDs) == 0 {
		// Fallback: return first N books
		result, err := s.bookRepo.List(repository.BookFilters{}, 1, count)
		if err != nil {
			return nil, err
		}
		return result.Items, nil
	}

	var books []model.Book
	for _, id := range bookIDs {
		book, err := s.bookRepo.GetByID(id)
		if err == nil {
			books = append(books, *book)
		}
	}
	return books, nil
}

func (s *BookService) GetStatistics() (*repository.BookStatistics, error) {
	return s.bookRepo.GetStatistics()
}

func (s *BookService) Create(dto CreateBookDTO) (*BookResult, error) {
	book := &model.Book{
		Name:        dto.Name,
		Author:      dto.Author,
		BookTypeID:  dto.BookTypeID,
		ConditionID: dto.ConditionID,
		GenreID:     dto.GenreID,
		PublisherID: dto.PublisherID,
		Year:        dto.Year,
		ISBN:        dto.ISBN,
		Summary:     dto.Summary,
		Price:       dto.Price,
		Quantity:    dto.Quantity,
	}

	if dto.CoverImage != nil {
		url, err := s.processImage(dto.CoverImage, dto.CoverFileName)
		if err != nil {
			return &BookResult{Success: false, Error: err.Error()}, nil
		}
		if url == "" {
			return &BookResult{Success: false, Error: "Image failed safety validation"}, nil
		}
		book.CoverImageURL = url
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return &BookResult{Success: true}, nil
}

func (s *BookService) Update(dto UpdateBookDTO) (*BookResult, error) {
	book, err := s.bookRepo.GetByID(dto.BookID)
	if err != nil {
		return nil, err
	}

	book.Name = dto.Name
	book.Author = dto.Author
	book.BookTypeID = dto.BookTypeID
	book.ConditionID = dto.ConditionID
	book.GenreID = dto.GenreID
	book.PublisherID = dto.PublisherID
	book.Year = dto.Year
	book.ISBN = dto.ISBN
	book.Summary = dto.Summary
	book.Price = dto.Price
	book.Quantity = dto.Quantity

	if dto.CoverImage != nil {
		url, err := s.processImage(dto.CoverImage, dto.CoverFileName)
		if err != nil {
			return &BookResult{Success: false, Error: err.Error()}, nil
		}
		if url == "" {
			return &BookResult{Success: false, Error: "Image failed safety validation"}, nil
		}
		// Delete old image
		if book.CoverImageURL != "" {
			_ = s.fileService.Delete(book.CoverImageURL)
		}
		book.CoverImageURL = url
	}

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return &BookResult{Success: true}, nil
}

func (s *BookService) processImage(file multipart.File, filename string) (string, error) {
	// Resize
	resized, err := s.imageResize.ResizeImage(file)
	if err != nil {
		return "", err
	}

	// Validate safety
	safe, err := s.imageValidation.IsSafe(resized)
	if err != nil {
		return "", err
	}
	if !safe {
		return "", nil
	}

	// Reset reader for saving
	if seeker, ok := resized.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	// Save
	url, err := s.fileService.Save(resized, filename)
	if err != nil {
		return "", err
	}

	return url, nil
}
