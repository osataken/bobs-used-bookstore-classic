package service

import (
	"testing"

	"bobs-used-bookstore-api/internal/database"
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/repository"
)

func setupTestDB(t *testing.T) *repository.Repositories {
	t.Helper()
	cfg := &config.Config{
		Database:    "local",
		DatabaseDSN: ":memory:",
	}
	db, err := database.Initialize(cfg)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	return repository.NewRepositories(db)
}

func TestBookServiceGetBook(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewBookService(repos.Book, repos.Order, NewImageResizeService(), NewLocalImageValidationService(), NewLocalFileService("/tmp/test"))

	book, err := svc.GetBook(1)
	if err != nil {
		t.Fatalf("failed to get book: %v", err)
	}
	if book.Name != "2020: The Apocalypse" {
		t.Errorf("expected '2020: The Apocalypse', got '%s'", book.Name)
	}
	if book.Price != 10.95 {
		t.Errorf("expected price 10.95, got %.2f", book.Price)
	}
}

func TestBookServiceSearchBooks(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewBookService(repos.Book, repos.Order, NewImageResizeService(), NewLocalImageValidationService(), NewLocalFileService("/tmp/test"))

	result, err := svc.SearchBooks("", "Name", 1, 10)
	if err != nil {
		t.Fatalf("failed to search books: %v", err)
	}
	if len(result.Items) != 8 {
		t.Errorf("expected 8 books, got %d", len(result.Items))
	}
}

func TestBookServiceSearchBooksWithFilter(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewBookService(repos.Book, repos.Order, NewImageResizeService(), NewLocalImageValidationService(), NewLocalFileService("/tmp/test"))

	result, err := svc.SearchBooks("Apocalypse", "Name", 1, 10)
	if err != nil {
		t.Fatalf("failed to search books: %v", err)
	}
	if len(result.Items) != 1 {
		t.Errorf("expected 1 book, got %d", len(result.Items))
	}
}

func TestBookServiceListBestSelling(t *testing.T) {
	repos := setupTestDB(t)
	svc := NewBookService(repos.Book, repos.Order, NewImageResizeService(), NewLocalImageValidationService(), NewLocalFileService("/tmp/test"))

	// No orders yet, so should fall back to first N books
	books, err := svc.ListBestSellingBooks(4)
	if err != nil {
		t.Fatalf("failed to list best selling books: %v", err)
	}
	if len(books) != 4 {
		t.Errorf("expected 4 books, got %d", len(books))
	}
}
