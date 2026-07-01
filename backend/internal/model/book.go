package model

const LowBookThreshold = 5

type Book struct {
	Entity
	Name          string            `json:"name" gorm:"not null"`
	Author        string            `json:"author" gorm:"not null"`
	Year          *int              `json:"year"`
	ISBN          string            `json:"isbn"`
	PublisherID   uint              `json:"publisherId" gorm:"not null"`
	Publisher     ReferenceDataItem `json:"publisher" gorm:"foreignKey:PublisherID;constraint:OnDelete:SET NULL"`
	BookTypeID    uint              `json:"bookTypeId" gorm:"not null"`
	BookType      ReferenceDataItem `json:"bookType" gorm:"foreignKey:BookTypeID;constraint:OnDelete:SET NULL"`
	GenreID       uint              `json:"genreId" gorm:"not null"`
	Genre         ReferenceDataItem `json:"genre" gorm:"foreignKey:GenreID;constraint:OnDelete:SET NULL"`
	ConditionID   uint              `json:"conditionId" gorm:"not null"`
	Condition     ReferenceDataItem `json:"condition" gorm:"foreignKey:ConditionID;constraint:OnDelete:SET NULL"`
	CoverImageURL string            `json:"coverImageUrl"`
	Summary       string            `json:"summary"`
	Price         float64           `json:"price" gorm:"not null"`
	Quantity      int               `json:"quantity" gorm:"not null;default:0"`
}

func (b *Book) IsInStock() bool {
	return b.Quantity > 0
}

func (b *Book) IsLowInStock() bool {
	return b.Quantity <= LowBookThreshold
}

func (b *Book) ReduceStockLevel(quantity int) {
	b.Quantity -= quantity
	if b.Quantity < 0 {
		b.Quantity = 0
	}
}
