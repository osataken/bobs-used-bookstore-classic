package model

import "testing"

func TestBookIsInStock(t *testing.T) {
	book := &Book{Quantity: 5}
	if !book.IsInStock() {
		t.Error("expected book to be in stock")
	}

	book.Quantity = 0
	if book.IsInStock() {
		t.Error("expected book to be out of stock")
	}
}

func TestBookIsLowInStock(t *testing.T) {
	book := &Book{Quantity: 5}
	if !book.IsLowInStock() {
		t.Error("expected book to be low in stock")
	}

	book.Quantity = 6
	if book.IsLowInStock() {
		t.Error("expected book to not be low in stock")
	}

	book.Quantity = 0
	if book.IsLowInStock() {
		t.Error("expected out of stock to not be low in stock")
	}
}

func TestBookReduceStockLevel(t *testing.T) {
	book := &Book{Quantity: 10}
	book.ReduceStockLevel(3)
	if book.Quantity != 7 {
		t.Errorf("expected quantity 7, got %d", book.Quantity)
	}

	book.ReduceStockLevel(10)
	if book.Quantity != 0 {
		t.Errorf("expected quantity 0, got %d", book.Quantity)
	}
}
