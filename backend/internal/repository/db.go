package repository

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.DatabaseMode == "aws" {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.ReferenceDataItem{},
		&model.Customer{},
		&model.Address{},
		&model.Book{},
		&model.ShoppingCart{},
		&model.ShoppingCartItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.Offer{},
	)
	if err != nil {
		return nil, err
	}

	if err := Seed(db); err != nil {
		return nil, err
	}

	return db, nil
}

type Repositories struct {
	Book          *BookRepository
	Order         *OrderRepository
	Customer      *CustomerRepository
	Address       *AddressRepository
	Offer         *OfferRepository
	Cart          *CartRepository
	ReferenceData *ReferenceDataRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Book:          NewBookRepository(db),
		Order:         NewOrderRepository(db),
		Customer:      NewCustomerRepository(db),
		Address:       NewAddressRepository(db),
		Offer:         NewOfferRepository(db),
		Cart:          NewCartRepository(db),
		ReferenceData: NewReferenceDataRepository(db),
	}
}
