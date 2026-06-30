package model

// Book represents a book in the bookstore inventory
type Book struct {
	Entity
	Name          string  `gorm:"type:varchar(255);not null" json:"name"`
	Author        string  `gorm:"type:varchar(255);not null" json:"author"`
	Year          *int    `json:"year"`
	ISBN          string  `gorm:"type:varchar(50)" json:"isbn"`
	PublisherID   int     `gorm:"not null" json:"publisherId"`
	BookTypeID    int     `gorm:"not null" json:"bookTypeId"`
	GenreID       int     `gorm:"not null" json:"genreId"`
	ConditionID   int     `gorm:"not null" json:"conditionId"`
	CoverImageURL string  `gorm:"type:varchar(500)" json:"coverImageUrl"`
	Summary       string  `gorm:"type:text" json:"summary"`
	Price         float64 `gorm:"type:decimal(18,2);not null" json:"price"`
	Quantity      int     `gorm:"not null" json:"quantity"`

	// Navigation properties
	Publisher ReferenceDataItem `gorm:"foreignKey:PublisherID;constraint:OnDelete:NO ACTION" json:"publisher,omitempty"`
	BookType  ReferenceDataItem `gorm:"foreignKey:BookTypeID;constraint:OnDelete:NO ACTION" json:"bookType,omitempty"`
	Genre     ReferenceDataItem `gorm:"foreignKey:GenreID;constraint:OnDelete:NO ACTION" json:"genre,omitempty"`
	Condition ReferenceDataItem `gorm:"foreignKey:ConditionID;constraint:OnDelete:NO ACTION" json:"condition,omitempty"`
}

const LowBookThreshold = 5

func (b *Book) IsInStock() bool {
	return b.Quantity > 0
}

func (b *Book) IsLowInStock() bool {
	return b.Quantity > 0 && b.Quantity <= LowBookThreshold
}

func (b *Book) ReduceStockLevel(quantity int) {
	b.Quantity -= quantity
	if b.Quantity < 0 {
		b.Quantity = 0
	}
}

func (Book) TableName() string {
	return "Book"
}
