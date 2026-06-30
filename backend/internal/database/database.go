package database

import (
	"log/slog"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialize sets up the database connection and runs migrations
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.IsAWS("database") {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	slog.Info("Running database migrations")
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

	// Seed data if database is empty
	if err := Seed(db); err != nil {
		slog.Error("Failed to seed database", "error", err)
		return nil, err
	}

	return db, nil
}
