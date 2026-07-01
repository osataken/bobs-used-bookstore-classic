package service

import (
	"io"
	"time"

	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type BookService struct {
	bookRepo     *repository.BookRepository
	orderRepo    *repository.OrderRepository
	fileService  FileService
	imageService ImageValidationService
}

func NewBookService(bookRepo *repository.BookRepository, orderRepo *repository.OrderRepository, fileService FileService, imageService ImageValidationService) *BookService {
	return &BookService{
		bookRepo:     bookRepo,
		orderRepo:    orderRepo,
		fileService:  fileService,
		imageService: imageService,
	}
}

func (s *BookService) GetBook(id uint) (*model.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookService) GetBooks(searchString string, sortBy string, pageIndex int, pageSize int) (*repository.PaginatedResult[model.Book], error) {
	return s.bookRepo.List(searchString, sortBy, pageIndex, pageSize)
}

func (s *BookService) GetBooksFiltered(filters repository.BookFilters, pageIndex int, pageSize int) (*repository.PaginatedResult[model.Book], error) {
	return s.bookRepo.ListFiltered(filters, pageIndex, pageSize)
}

func (s *BookService) ListBestSellingBooks(count int) ([]model.Book, error) {
	return s.orderRepo.ListBestSellingBooks(count)
}

func (s *BookService) GetStatistics() (*repository.BookStatistics, error) {
	return s.bookRepo.GetStatistics()
}

type CreateBookInput struct {
	Name         string
	Author       string
	BookTypeID   uint
	ConditionID  uint
	GenreID      uint
	PublisherID  uint
	Year         *int
	ISBN         string
	Summary      string
	Price        float64
	Quantity     int
	CoverImage   io.Reader
	CoverImageFileName string
}

type BookResult struct {
	Success      bool
	ErrorMessage string
}

func (s *BookService) Add(input CreateBookInput) (*BookResult, error) {
	book := &model.Book{
		Name:        input.Name,
		Author:      input.Author,
		ISBN:        input.ISBN,
		PublisherID: input.PublisherID,
		BookTypeID:  input.BookTypeID,
		GenreID:     input.GenreID,
		ConditionID: input.ConditionID,
		Price:       input.Price,
		Quantity:    input.Quantity,
		Year:        input.Year,
		Summary:     input.Summary,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return s.saveImage(book, input.CoverImage, input.CoverImageFileName)
}

type UpdateBookInput struct {
	BookID       uint
	Name         string
	Author       string
	BookTypeID   uint
	ConditionID  uint
	GenreID      uint
	PublisherID  uint
	Year         *int
	ISBN         string
	Summary      string
	Price        float64
	Quantity     int
	CoverImage   io.Reader
	CoverImageFileName string
}

func (s *BookService) Update(input UpdateBookInput) (*BookResult, error) {
	book, err := s.bookRepo.GetByID(input.BookID)
	if err != nil {
		return nil, err
	}

	book.Name = input.Name
	book.Author = input.Author
	book.ISBN = input.ISBN
	book.PublisherID = input.PublisherID
	book.BookTypeID = input.BookTypeID
	book.GenreID = input.GenreID
	book.ConditionID = input.ConditionID
	book.Price = input.Price
	book.Quantity = input.Quantity
	book.Year = input.Year
	book.Summary = input.Summary
	book.UpdatedOn = time.Now().UTC()

	return s.saveImage(book, input.CoverImage, input.CoverImageFileName)
}

func (s *BookService) saveImage(book *model.Book, coverImage io.Reader, fileName string) (*BookResult, error) {
	if coverImage == nil {
		if err := s.bookRepo.Update(book); err != nil {
			return nil, err
		}
		return &BookResult{Success: true}, nil
	}

	resized, err := ResizeImage(coverImage)
	if err != nil {
		return nil, err
	}

	safe, err := s.imageService.IsSafe(resized)
	if err != nil {
		return nil, err
	}
	if !safe {
		return &BookResult{Success: false, ErrorMessage: "The image failed the safety check. Please try another image."}, nil
	}

	imageURL, err := s.fileService.Save(resized, fileName)
	if err != nil {
		return nil, err
	}

	if book.CoverImageURL != "" {
		_ = s.fileService.Delete(book.CoverImageURL)
	}
	book.CoverImageURL = imageURL

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return &BookResult{Success: true}, nil
}
