package repository

import (
	"bobs-used-bookstore-api/internal/model"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&model.ReferenceDataItem{}).Count(&count)
	if count > 0 {
		return nil
	}

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

	if err := db.Create(&referenceData).Error; err != nil {
		return err
	}

	year2020 := 2020
	year2019 := 2019
	year2018 := 2018
	year2017 := 2017
	year2016 := 2016
	year2015 := 2015
	year2014 := 2014
	year2013 := 2013

	books := []model.Book{
		{Entity: model.Entity{ID: 1}, Name: "2020: The Apocalypse", Author: "Li Juan", Year: &year2020, ISBN: "978-0-13-468599-1", PublisherID: 15, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 10.95, Quantity: 25, Summary: "A thrilling science fiction novel about survival in a post-apocalyptic world."},
		{Entity: model.Entity{ID: 2}, Name: "Children Of Iron", Author: "Nikki Wolf", Year: &year2019, ISBN: "978-0-13-468599-2", PublisherID: 16, BookTypeID: 1, GenreID: 11, ConditionID: 6, Price: 13.95, Quantity: 3, Summary: "A captivating tale of resilience and determination."},
		{Entity: model.Entity{ID: 3}, Name: "Gold In The Dark", Author: "Richard Roe", Year: &year2018, ISBN: "978-0-13-468599-3", PublisherID: 17, BookTypeID: 1, GenreID: 13, ConditionID: 5, Price: 6.50, Quantity: 10, Summary: "An adventure through dark times seeking the light."},
		{Entity: model.Entity{ID: 4}, Name: "Leagues Of Smoke", Author: "Pat Candella", Year: &year2017, ISBN: "978-0-13-468599-4", PublisherID: 18, BookTypeID: 2, GenreID: 11, ConditionID: 7, Price: 3.00, Quantity: 1, Summary: "A mysterious journey through smoky leagues."},
		{Entity: model.Entity{ID: 5}, Name: "Alone With The Stars", Author: "Carlos Salazar", Year: &year2016, ISBN: "978-0-13-468599-5", PublisherID: 19, BookTypeID: 2, GenreID: 12, ConditionID: 5, Price: 15.95, Quantity: 5, Summary: "A suspenseful thriller set under the stars."},
		{Entity: model.Entity{ID: 6}, Name: "The Girl In The Polaroid", Author: "Terri Whitlock", Year: &year2015, ISBN: "978-0-13-468599-6", PublisherID: 20, BookTypeID: 1, GenreID: 12, ConditionID: 6, Price: 8.25, Quantity: 2, Summary: "A mystery unraveled through a single photograph."},
		{Entity: model.Entity{ID: 7}, Name: "1001 Jokes", Author: "Mary Major", Year: &year2014, ISBN: "978-0-13-468599-7", PublisherID: 21, BookTypeID: 2, GenreID: 11, ConditionID: 5, Price: 13.95, Quantity: 7, Summary: "A collection of jokes for all ages."},
		{Entity: model.Entity{ID: 8}, Name: "My Search For Meaning", Author: "Mateo Jackson", Year: &year2013, ISBN: "978-0-13-468599-8", PublisherID: 22, BookTypeID: 3, GenreID: 8, ConditionID: 7, Price: 5.00, Quantity: 15, Summary: "An autobiographical journey of self-discovery."},
	}

	return db.Create(&books).Error
}
