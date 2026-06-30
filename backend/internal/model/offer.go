package model

type Offer struct {
	Entity
	Author      string      `gorm:"type:text" json:"author"`
	ISBN        string      `gorm:"type:text" json:"isbn"`
	BookName    string      `gorm:"type:text" json:"bookName"`
	FrontUrl    string      `gorm:"type:text" json:"frontUrl"`
	GenreID     int         `gorm:"not null" json:"genreId"`
	ConditionID int         `gorm:"not null" json:"conditionId"`
	PublisherID int         `gorm:"not null" json:"publisherId"`
	BookTypeID  int         `gorm:"not null" json:"bookTypeId"`
	Summary     string      `gorm:"type:text" json:"summary"`
	OfferStatus OfferStatus `gorm:"not null;default:0" json:"offerStatus"`
	Comment     string      `gorm:"type:text" json:"comment"`
	CustomerID  int         `gorm:"not null" json:"customerId"`
	BookPrice   float64     `gorm:"type:decimal(18,2);not null" json:"bookPrice"`

	Customer  *Customer          `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	Genre     *ReferenceDataItem `gorm:"foreignKey:GenreID" json:"genre,omitempty"`
	Condition *ReferenceDataItem `gorm:"foreignKey:ConditionID" json:"condition,omitempty"`
	Publisher *ReferenceDataItem `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	BookType  *ReferenceDataItem `gorm:"foreignKey:BookTypeID" json:"bookType,omitempty"`
}
