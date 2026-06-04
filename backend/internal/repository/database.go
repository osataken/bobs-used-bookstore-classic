package repository

import (
	"bobs-used-bookstore-api/internal/model"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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
		&model.ShoppingCart{},
		&model.ShoppingCartItem{},
		&model.Offer{},
	)
	if err != nil {
		return nil, err
	}

	seedDatabase(db)

	return db, nil
}

func seedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&model.ReferenceDataItem{}).Count(&count)
	if count > 0 {
		return
	}

	slog.Info("Seeding database...")

	referenceData := []model.ReferenceDataItem{
		{Entity: model.Entity{ID: 1}, DataType: model.ReferenceDataTypeBookType, Text: "Hardcover"},
		{Entity: model.Entity{ID: 2}, DataType: model.ReferenceDataTypeBookType, Text: "Trade Paperback"},
		{Entity: model.Entity{ID: 3}, DataType: model.ReferenceDataTypeBookType, Text: "Mass Market Paperback"},
		{Entity: model.Entity{ID: 4}, DataType: model.ReferenceDataTypeCondition, Text: "New"},
		{Entity: model.Entity{ID: 5}, DataType: model.ReferenceDataTypeCondition, Text: "Like New"},
		{Entity: model.Entity{ID: 6}, DataType: model.ReferenceDataTypeCondition, Text: "Good"},
		{Entity: model.Entity{ID: 7}, DataType: model.ReferenceDataTypeCondition, Text: "Acceptable"},
		{Entity: model.Entity{ID: 8}, DataType: model.ReferenceDataTypeGenre, Text: "Biographies"},
		{Entity: model.Entity{ID: 9}, DataType: model.ReferenceDataTypeGenre, Text: "Children's Books"},
		{Entity: model.Entity{ID: 10}, DataType: model.ReferenceDataTypeGenre, Text: "History"},
		{Entity: model.Entity{ID: 11}, DataType: model.ReferenceDataTypeGenre, Text: "Literature & Fiction"},
		{Entity: model.Entity{ID: 12}, DataType: model.ReferenceDataTypeGenre, Text: "Mystery, Thriller & Suspense"},
		{Entity: model.Entity{ID: 13}, DataType: model.ReferenceDataTypeGenre, Text: "Science Fiction & Fantasy"},
		{Entity: model.Entity{ID: 14}, DataType: model.ReferenceDataTypeGenre, Text: "Travel"},
		{Entity: model.Entity{ID: 15}, DataType: model.ReferenceDataTypePublisher, Text: "Arcadia Books"},
		{Entity: model.Entity{ID: 16}, DataType: model.ReferenceDataTypePublisher, Text: "Astral Publishing"},
		{Entity: model.Entity{ID: 17}, DataType: model.ReferenceDataTypePublisher, Text: "Moonlight Publishing"},
		{Entity: model.Entity{ID: 18}, DataType: model.ReferenceDataTypePublisher, Text: "Dreamscape Press"},
		{Entity: model.Entity{ID: 19}, DataType: model.ReferenceDataTypePublisher, Text: "Enchanted Library"},
		{Entity: model.Entity{ID: 20}, DataType: model.ReferenceDataTypePublisher, Text: "Fantasia House"},
		{Entity: model.Entity{ID: 21}, DataType: model.ReferenceDataTypePublisher, Text: "Horizon Books"},
		{Entity: model.Entity{ID: 22}, DataType: model.ReferenceDataTypePublisher, Text: "Infinity Press"},
		{Entity: model.Entity{ID: 23}, DataType: model.ReferenceDataTypePublisher, Text: "Paradigm Publishing"},
		{Entity: model.Entity{ID: 24}, DataType: model.ReferenceDataTypePublisher, Text: "Aurora Publishing"},
	}
	db.Create(&referenceData)

	books := []model.Book{
		{Entity: model.Entity{ID: 1}, Name: "2020: The Apocalypse", Author: "Li Juan", ISBN: "6556784356", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 10.95, Quantity: 25, CoverImageURL: "/images/coverimages/apocalypse.png"},
		{Entity: model.Entity{ID: 2}, Name: "Children Of Iron", Author: "Nikki Wolf", ISBN: "7665438976", PublisherID: 16, BookTypeID: 1, GenreID: 11, ConditionID: 6, Price: 13.95, Quantity: 3, CoverImageURL: "/images/coverimages/childrenofiron.png"},
		{Entity: model.Entity{ID: 3}, Name: "Gold In The Dark", Author: "Richard Roe", ISBN: "5442280765", PublisherID: 17, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 6.50, Quantity: 10, CoverImageURL: "/images/coverimages/goldinthedark.png"},
		{Entity: model.Entity{ID: 4}, Name: "Leagues Of Smoke", Author: "Pat Candella", ISBN: "4556789542", PublisherID: 18, BookTypeID: 2, GenreID: 11, ConditionID: 7, Price: 3.00, Quantity: 1, CoverImageURL: "/images/coverimages/leaguesofsmoke.png"},
		{Entity: model.Entity{ID: 5}, Name: "Alone With The Stars", Author: "Carlos Salazar", ISBN: "4563358087", PublisherID: 19, BookTypeID: 2, GenreID: 12, ConditionID: 5, Price: 15.95, Quantity: 5, CoverImageURL: "/images/coverimages/alonewiththestars.png"},
		{Entity: model.Entity{ID: 6}, Name: "The Girl In The Polaroid", Author: "Terri Whitlock", ISBN: "2354435678", PublisherID: 20, BookTypeID: 1, GenreID: 12, ConditionID: 6, Price: 8.25, Quantity: 2, CoverImageURL: "/images/coverimages/girlinthepolaroid.png"},
		{Entity: model.Entity{ID: 7}, Name: "1001 Jokes", Author: "Mary Major", ISBN: "6554789632", PublisherID: 21, BookTypeID: 2, GenreID: 11, ConditionID: 5, Price: 13.95, Quantity: 7, CoverImageURL: "/images/coverimages/1001jokes.png"},
		{Entity: model.Entity{ID: 8}, Name: "My Search For Meaning", Author: "Mateo Jackson", ISBN: "4558786554", PublisherID: 22, BookTypeID: 3, GenreID: 8, ConditionID: 7, Price: 5.00, Quantity: 15, CoverImageURL: "/images/coverimages/mysearchformeaning.png"},
	}
	db.Create(&books)

	slog.Info("Database seeded successfully")
}
