package service

import (
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetByID(id int) (*model.Book, error) {
	return s.repo.GetByID(id)
}

func (s *BookService) Search(pageIndex, pageSize int, searchString, sortBy string) (repository.PaginatedList[model.Book], error) {
	return s.repo.GetAll(pageIndex, pageSize, searchString, sortBy)
}

func (s *BookService) GetFiltered(name, author string, publisherID, genreID, bookTypeID, conditionID int, lowStock bool, pageIndex, pageSize int) (repository.PaginatedList[model.Book], error) {
	return s.repo.GetFiltered(name, author, publisherID, genreID, bookTypeID, conditionID, lowStock, pageIndex, pageSize)
}

func (s *BookService) Create(book *model.Book) error {
	return s.repo.Create(book)
}

func (s *BookService) Update(book *model.Book) error {
	return s.repo.Update(book)
}

func (s *BookService) GetBestSellers(limit int) ([]model.Book, error) {
	return s.repo.GetBestSellers(limit)
}

func (s *BookService) GetTotalCount() (int64, error) {
	return s.repo.GetTotalCount()
}

func (s *BookService) GetLowStockCount() (int64, error) {
	return s.repo.GetLowStockCount()
}
