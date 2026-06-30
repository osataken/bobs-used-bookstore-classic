package model

// Offer represents a customer's resale offer
type Offer struct {
	Entity
	Author      string      `gorm:"type:varchar(255)" json:"author"`
	ISBN        string      `gorm:"type:varchar(50)" json:"isbn"`
	BookName    string      `gorm:"type:varchar(255);not null" json:"bookName"`
	FrontURL    string      `gorm:"type:varchar(500)" json:"frontUrl"`
	GenreID     int         `gorm:"not null" json:"genreId"`
	ConditionID int         `gorm:"not null" json:"conditionId"`
	PublisherID int         `gorm:"not null" json:"publisherId"`
	BookTypeID  int         `gorm:"not null" json:"bookTypeId"`
	Summary     string      `gorm:"type:text" json:"summary"`
	OfferStatus OfferStatus `gorm:"default:0" json:"offerStatus"`
	Comment     string      `gorm:"type:text" json:"comment"`
	CustomerID  int         `gorm:"not null" json:"customerId"`
	BookPrice   float64     `gorm:"type:decimal(18,2)" json:"bookPrice"`

	// Navigation properties
	Customer  Customer          `gorm:"foreignKey:CustomerID;constraint:OnDelete:NO ACTION" json:"customer,omitempty"`
	Genre     ReferenceDataItem `gorm:"foreignKey:GenreID;constraint:OnDelete:NO ACTION" json:"genre,omitempty"`
	Condition ReferenceDataItem `gorm:"foreignKey:ConditionID;constraint:OnDelete:NO ACTION" json:"condition,omitempty"`
	Publisher ReferenceDataItem `gorm:"foreignKey:PublisherID;constraint:OnDelete:NO ACTION" json:"publisher,omitempty"`
	BookType  ReferenceDataItem `gorm:"foreignKey:BookTypeID;constraint:OnDelete:NO ACTION" json:"bookType,omitempty"`
}

func (Offer) TableName() string {
	return "Offer"
}
