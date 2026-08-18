package model

// Address entity
type Address struct {
	Entity
	AddressLine1 string   `gorm:"not null" json:"addressLine1"`
	AddressLine2 string   `json:"addressLine2"`
	City         string   `gorm:"not null" json:"city"`
	State        string   `gorm:"not null" json:"state"`
	Country      string   `gorm:"not null" json:"country"`
	ZipCode      string   `gorm:"not null" json:"zipCode"`
	CustomerID   int      `gorm:"not null" json:"customerId"`
	Customer     Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:RESTRICT" json:"-"`
	IsActive     bool     `gorm:"default:true;not null" json:"isActive"`
}

func (Address) TableName() string {
	return "Address"
}
