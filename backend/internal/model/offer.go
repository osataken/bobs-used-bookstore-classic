package model

type Offer struct {
	Entity
	CustomerID  uint              `json:"customerId" gorm:"not null"`
	Customer    Customer          `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	BookName    string            `json:"bookName" gorm:"not null"`
	Author      string            `json:"author" gorm:"not null"`
	ISBN        string            `json:"isbn"`
	FrontURL    string            `json:"frontUrl"`
	Summary     string            `json:"summary"`
	Comment     string            `json:"comment"`
	GenreID     uint              `json:"genreId" gorm:"not null"`
	Genre       ReferenceDataItem `json:"genre" gorm:"foreignKey:GenreID;constraint:OnDelete:SET NULL"`
	ConditionID uint              `json:"conditionId" gorm:"not null"`
	Condition   ReferenceDataItem `json:"condition" gorm:"foreignKey:ConditionID;constraint:OnDelete:SET NULL"`
	PublisherID uint              `json:"publisherId" gorm:"not null"`
	Publisher   ReferenceDataItem `json:"publisher" gorm:"foreignKey:PublisherID;constraint:OnDelete:SET NULL"`
	BookTypeID  uint              `json:"bookTypeId" gorm:"not null"`
	BookType    ReferenceDataItem `json:"bookType" gorm:"foreignKey:BookTypeID;constraint:OnDelete:SET NULL"`
	OfferStatus OfferStatus       `json:"offerStatus" gorm:"default:0"`
	BookPrice   float64           `json:"bookPrice" gorm:"not null"`
}
