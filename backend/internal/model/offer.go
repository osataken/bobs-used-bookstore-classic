package model

// Offer entity
type Offer struct {
	Entity
	Author      string            `json:"author"`
	ISBN        string            `json:"isbn"`
	BookName    string            `gorm:"not null" json:"bookName"`
	FrontUrl    string            `json:"frontUrl"`
	GenreID     int               `gorm:"not null" json:"genreId"`
	Genre       ReferenceDataItem `gorm:"foreignKey:GenreID;constraint:OnDelete:RESTRICT" json:"genre,omitempty"`
	ConditionID int               `gorm:"not null" json:"conditionId"`
	Condition   ReferenceDataItem `gorm:"foreignKey:ConditionID;constraint:OnDelete:RESTRICT" json:"condition,omitempty"`
	PublisherID int               `gorm:"not null" json:"publisherId"`
	Publisher   ReferenceDataItem `gorm:"foreignKey:PublisherID;constraint:OnDelete:RESTRICT" json:"publisher,omitempty"`
	BookTypeID  int               `gorm:"not null" json:"bookTypeId"`
	BookType    ReferenceDataItem `gorm:"foreignKey:BookTypeID;constraint:OnDelete:RESTRICT" json:"bookType,omitempty"`
	Summary     string            `json:"summary"`
	OfferStatus OfferStatus       `gorm:"default:0" json:"offerStatus"`
	Comment     string            `json:"comment"`
	CustomerID  int               `gorm:"not null" json:"customerId"`
	Customer    Customer          `gorm:"foreignKey:CustomerID;constraint:OnDelete:RESTRICT" json:"customer,omitempty"`
	BookPrice   float64           `gorm:"type:decimal(18,2);not null" json:"bookPrice"`
}

func (Offer) TableName() string {
	return "Offer"
}
