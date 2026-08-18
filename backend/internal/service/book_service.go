package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type BookService struct {
	repo                   *repository.BookRepository
	fileService            FileServiceInterface
	imageResizeService     ImageResizeServiceInterface
	imageValidationService ImageValidationServiceInterface
}

func NewBookService(repo *repository.BookRepository, fs FileServiceInterface, irs ImageResizeServiceInterface, ivs ImageValidationServiceInterface) *BookService {
	return &BookService{
		repo:                   repo,
		fileService:            fs,
		imageResizeService:     irs,
		imageValidationService: ivs,
	}
}

func (s *BookService) GetBook(id int) (*model.Book, error) {
	return s.repo.GetByID(id)
}

func (s *BookService) SearchBooks(searchString, sortBy string, pageIndex, pageSize int) (*model.PaginatedList, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.List(searchString, sortBy, pageIndex, pageSize)
}

func (s *BookService) ListBooks(filters map[string]interface{}, pageIndex, pageSize int) (*model.PaginatedList, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListWithFilters(filters, pageIndex, pageSize)
}

func (s *BookService) ListBestSellingBooks(limit int) ([]model.Book, error) {
	return s.repo.ListBestSelling(limit)
}

func (s *BookService) GetStatistics() (*model.BookStatistics, error) {
	return s.repo.GetStatistics()
}

func (s *BookService) CreateBook(book *model.Book, imageData []byte, imageFilename string) error {
	if len(imageData) > 0 {
		// Resize image
		resized, err := s.imageResizeService.ResizeImage(imageData, 400, 600)
		if err != nil {
			return err
		}

		// Validate image
		safe, err := s.imageValidationService.IsSafe(resized)
		if err != nil {
			return err
		}
		if !safe {
			return ErrImageUnsafe
		}

		// Save image
		url, err := s.fileService.Save(resized, imageFilename)
		if err != nil {
			return err
		}
		book.CoverImageUrl = url
	}

	return s.repo.Create(book)
}

func (s *BookService) UpdateBook(book *model.Book, imageData []byte, imageFilename string) error {
	if len(imageData) > 0 {
		// Resize image
		resized, err := s.imageResizeService.ResizeImage(imageData, 400, 600)
		if err != nil {
			return err
		}

		// Validate image
		safe, err := s.imageValidationService.IsSafe(resized)
		if err != nil {
			return err
		}
		if !safe {
			return ErrImageUnsafe
		}

		// Delete old image if exists
		if book.CoverImageUrl != "" {
			_ = s.fileService.Delete(book.CoverImageUrl)
		}

		// Save new image
		url, err := s.fileService.Save(resized, imageFilename)
		if err != nil {
			return err
		}
		book.CoverImageUrl = url
	}

	return s.repo.Update(book)
}
