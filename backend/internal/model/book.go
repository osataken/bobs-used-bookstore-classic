package model

type Book struct {
	Entity
	Name          string  `gorm:"type:text" json:"name"`
	Author        string  `gorm:"type:text" json:"author"`
	Year          *int    `json:"year"`
	ISBN          string  `gorm:"type:text" json:"isbn"`
	PublisherID   int     `gorm:"not null" json:"publisherId"`
	BookTypeID    int     `gorm:"not null" json:"bookTypeId"`
	GenreID       int     `gorm:"not null" json:"genreId"`
	ConditionID   int     `gorm:"not null" json:"conditionId"`
	CoverImageUrl string  `gorm:"type:text" json:"coverImageUrl"`
	Summary       string  `gorm:"type:text" json:"summary"`
	Price         float64 `gorm:"type:decimal(18,2);not null" json:"price"`
	Quantity      int     `gorm:"not null" json:"quantity"`

	Publisher *ReferenceDataItem `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	BookType  *ReferenceDataItem `gorm:"foreignKey:BookTypeID" json:"bookType,omitempty"`
	Genre     *ReferenceDataItem `gorm:"foreignKey:GenreID" json:"genre,omitempty"`
	Condition *ReferenceDataItem `gorm:"foreignKey:ConditionID" json:"condition,omitempty"`
}

func (b *Book) IsInStock() bool {
	return b.Quantity > 0
}

func (b *Book) IsLowInStock() bool {
	return b.Quantity <= 5
}

func (b *Book) ReduceStockLevel(quantity int) {
	b.Quantity -= quantity
	if b.Quantity < 0 {
		b.Quantity = 0
	}
}

const LowBookThreshold = 5
