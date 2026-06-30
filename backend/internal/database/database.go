package database

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/model"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.IsLocalDB() {
		db, err = gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.ReferenceDataItem{},
		&model.Book{},
		&model.Customer{},
		&model.Address{},
		&model.Order{},
		&model.OrderItem{},
		&model.Offer{},
		&model.ShoppingCart{},
		&model.ShoppingCartItem{},
	)
	if err != nil {
		return nil, err
	}

	if err := Seed(db); err != nil {
		slog.Error("failed to seed database", "error", err)
	}

	return db, nil
}
