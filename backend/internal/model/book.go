package model

type Book struct {
	Entity
	Name          string            `gorm:"column:name" json:"name"`
	Author        string            `gorm:"column:author" json:"author"`
	Year          *int              `gorm:"column:year" json:"year"`
	ISBN          string            `gorm:"column:isbn" json:"isbn"`
	PublisherID   int               `gorm:"column:publisher_id" json:"publisherId"`
	Publisher     ReferenceDataItem `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	BookTypeID    int               `gorm:"column:book_type_id" json:"bookTypeId"`
	BookType      ReferenceDataItem `gorm:"foreignKey:BookTypeID" json:"bookType,omitempty"`
	GenreID       int               `gorm:"column:genre_id" json:"genreId"`
	Genre         ReferenceDataItem `gorm:"foreignKey:GenreID" json:"genre,omitempty"`
	ConditionID   int               `gorm:"column:condition_id" json:"conditionId"`
	Condition     ReferenceDataItem `gorm:"foreignKey:ConditionID" json:"condition,omitempty"`
	CoverImageURL string            `gorm:"column:cover_image_url" json:"coverImageUrl"`
	Summary       string            `gorm:"column:summary" json:"summary"`
	Price         float64           `gorm:"column:price;type:decimal(18,2)" json:"price"`
	Quantity      int               `gorm:"column:quantity" json:"quantity"`
}

func (Book) TableName() string {
	return "book"
}

const LowBookThreshold = 5

func (b *Book) IsInStock() bool {
	return b.Quantity > 0
}

func (b *Book) IsLowInStock() bool {
	return b.Quantity > LowBookThreshold
}

func (b *Book) ReduceStockLevel(quantity int) {
	b.Quantity -= quantity
	if b.Quantity < 0 {
		b.Quantity = 0
	}
}
