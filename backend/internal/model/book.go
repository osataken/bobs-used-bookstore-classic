package model

const LowBookThreshold = 5

// Book entity
type Book struct {
	Entity
	Name          string            `gorm:"not null" json:"name"`
	Author        string            `gorm:"not null" json:"author"`
	Year          *int              `json:"year"`
	ISBN          string            `gorm:"not null" json:"isbn"`
	PublisherID   int               `gorm:"not null" json:"publisherId"`
	Publisher     ReferenceDataItem `gorm:"foreignKey:PublisherID;constraint:OnDelete:RESTRICT" json:"publisher,omitempty"`
	BookTypeID    int               `gorm:"not null" json:"bookTypeId"`
	BookType      ReferenceDataItem `gorm:"foreignKey:BookTypeID;constraint:OnDelete:RESTRICT" json:"bookType,omitempty"`
	GenreID       int               `gorm:"not null" json:"genreId"`
	Genre         ReferenceDataItem `gorm:"foreignKey:GenreID;constraint:OnDelete:RESTRICT" json:"genre,omitempty"`
	ConditionID   int               `gorm:"not null" json:"conditionId"`
	Condition     ReferenceDataItem `gorm:"foreignKey:ConditionID;constraint:OnDelete:RESTRICT" json:"condition,omitempty"`
	CoverImageUrl string            `json:"coverImageUrl"`
	Summary       string            `json:"summary"`
	Price         float64           `gorm:"type:decimal(18,2);not null" json:"price"`
	Quantity      int               `gorm:"not null;default:0" json:"quantity"`
}

func (Book) TableName() string {
	return "Book"
}

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
