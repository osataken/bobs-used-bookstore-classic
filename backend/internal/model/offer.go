package model

type OfferStatus int

const (
	OfferStatusPendingApproval OfferStatus = 0
	OfferStatusApproved        OfferStatus = 1
	OfferStatusReceived        OfferStatus = 2
	OfferStatusPaid            OfferStatus = 3
	OfferStatusRejected        OfferStatus = 4
)

type Offer struct {
	Entity
	BookName    string            `gorm:"column:book_name" json:"bookName"`
	Author      string            `gorm:"column:author" json:"author"`
	ISBN        string            `gorm:"column:isbn" json:"isbn"`
	FrontURL    string            `gorm:"column:front_url" json:"frontUrl"`
	GenreID     int               `gorm:"column:genre_id" json:"genreId"`
	Genre       ReferenceDataItem `gorm:"foreignKey:GenreID" json:"genre,omitempty"`
	ConditionID int               `gorm:"column:condition_id" json:"conditionId"`
	Condition   ReferenceDataItem `gorm:"foreignKey:ConditionID" json:"condition,omitempty"`
	PublisherID int               `gorm:"column:publisher_id" json:"publisherId"`
	Publisher   ReferenceDataItem `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	BookTypeID  int               `gorm:"column:book_type_id" json:"bookTypeId"`
	BookType    ReferenceDataItem `gorm:"foreignKey:BookTypeID" json:"bookType,omitempty"`
	Summary     string            `gorm:"column:summary" json:"summary"`
	OfferStatus OfferStatus       `gorm:"column:offer_status;default:0" json:"offerStatus"`
	Comment     string            `gorm:"column:comment" json:"comment"`
	CustomerID  int               `gorm:"column:customer_id" json:"customerId"`
	Customer    Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	BookPrice   float64           `gorm:"column:book_price;type:decimal(18,2)" json:"bookPrice"`
}

func (Offer) TableName() string {
	return "offer"
}
