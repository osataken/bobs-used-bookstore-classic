package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type BookService struct {
	bookRepo     *repository.BookRepository
	fileService  FileService
	imageService ImageValidationService
	uploadDir    string
}

type FileService interface {
	Save(reader io.Reader, filename string) (string, error)
	Delete(path string) error
}

type ImageValidationService interface {
	IsSafe(reader io.Reader) (bool, error)
}

func NewBookService(bookRepo *repository.BookRepository, fileService FileService, imageService ImageValidationService, uploadDir string) *BookService {
	return &BookService{
		bookRepo:     bookRepo,
		fileService:  fileService,
		imageService: imageService,
		uploadDir:    uploadDir,
	}
}

func (s *BookService) GetBook(id int) (*model.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookService) GetBooks(filters repository.BookFilters, pageIndex, pageSize int) (*repository.PaginatedList, error) {
	return s.bookRepo.List(filters, pageIndex, pageSize)
}

func (s *BookService) SearchBooks(searchString, sortBy string, pageIndex, pageSize int) (*repository.PaginatedList, error) {
	return s.bookRepo.Search(searchString, sortBy, pageIndex, pageSize)
}

func (s *BookService) ListBestSelling(count int) ([]model.Book, error) {
	return s.bookRepo.ListBestSelling(count)
}

func (s *BookService) GetStatistics() (*repository.BookStatistics, error) {
	return s.bookRepo.GetStatistics()
}

func (s *BookService) Add(dto CreateBookDTO) (*model.Book, error) {
	book := &model.Book{
		Name:        dto.Name,
		Author:      dto.Author,
		ISBN:        dto.ISBN,
		PublisherID: dto.PublisherID,
		BookTypeID:  dto.BookTypeID,
		GenreID:     dto.GenreID,
		ConditionID: dto.ConditionID,
		Price:       dto.Price,
		Quantity:    dto.Quantity,
		Year:        dto.Year,
		Summary:     dto.Summary,
	}

	if dto.CoverImage != nil {
		imageURL, err := s.processImage(dto.CoverImage, dto.CoverImageFileName)
		if err != nil {
			return nil, err
		}
		book.CoverImageURL = imageURL
	}

	err := s.bookRepo.Add(book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) Update(dto UpdateBookDTO) (*model.Book, error) {
	book, err := s.bookRepo.GetByID(dto.BookID)
	if err != nil {
		return nil, err
	}

	book.Name = dto.Name
	book.Author = dto.Author
	book.ISBN = dto.ISBN
	book.PublisherID = dto.PublisherID
	book.BookTypeID = dto.BookTypeID
	book.GenreID = dto.GenreID
	book.ConditionID = dto.ConditionID
	book.Price = dto.Price
	book.Quantity = dto.Quantity
	book.Year = dto.Year
	book.Summary = dto.Summary
	book.UpdatedOn = time.Now().UTC()

	if dto.CoverImage != nil {
		imageURL, err := s.processImage(dto.CoverImage, dto.CoverImageFileName)
		if err != nil {
			return nil, err
		}
		if book.CoverImageURL != "" {
			_ = s.fileService.Delete(book.CoverImageURL)
		}
		book.CoverImageURL = imageURL
	}

	err = s.bookRepo.Update(book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) processImage(image io.Reader, filename string) (string, error) {
	safe, err := s.imageService.IsSafe(image)
	if err != nil {
		return "", err
	}
	if !safe {
		return "", fmt.Errorf("the image failed the safety check. Please try another image")
	}

	imageURL, err := s.fileService.Save(image, filename)
	if err != nil {
		return "", err
	}
	return imageURL, nil
}

// LocalFileService implements FileService for local filesystem
type LocalFileService struct {
	uploadDir string
}

func NewLocalFileService(uploadDir string) *LocalFileService {
	os.MkdirAll(uploadDir, 0755)
	return &LocalFileService{uploadDir: uploadDir}
}

func (s *LocalFileService) Save(reader io.Reader, filename string) (string, error) {
	if reader == nil {
		return "", nil
	}

	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadDir, newFilename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}

	return "/uploads/" + newFilename, nil
}

func (s *LocalFileService) Delete(path string) error {
	if path == "" {
		return nil
	}
	fullPath := filepath.Join(s.uploadDir, filepath.Base(path))
	return os.Remove(fullPath)
}

// LocalImageValidationService always returns safe
type LocalImageValidationService struct{}

func NewLocalImageValidationService() *LocalImageValidationService {
	return &LocalImageValidationService{}
}

func (s *LocalImageValidationService) IsSafe(reader io.Reader) (bool, error) {
	return true, nil
}

type CreateBookDTO struct {
	Name               string
	Author             string
	ISBN               string
	PublisherID        int
	BookTypeID         int
	GenreID            int
	ConditionID        int
	Price              float64
	Quantity           int
	Year               *int
	Summary            string
	CoverImage         io.Reader
	CoverImageFileName string
}

type UpdateBookDTO struct {
	BookID             int
	Name               string
	Author             string
	ISBN               string
	PublisherID        int
	BookTypeID         int
	GenreID            int
	ConditionID        int
	Price              float64
	Quantity           int
	Year               *int
	Summary            string
	CoverImage         io.Reader
	CoverImageFileName string
}
