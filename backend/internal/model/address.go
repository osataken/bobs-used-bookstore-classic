package model

type Address struct {
	Entity
	AddressLine1 string `gorm:"type:text" json:"addressLine1"`
	AddressLine2 string `gorm:"type:text" json:"addressLine2"`
	City         string `gorm:"type:text" json:"city"`
	State        string `gorm:"type:text" json:"state"`
	Country      string `gorm:"type:text" json:"country"`
	ZipCode      string `gorm:"type:text" json:"zipCode"`
	CustomerID   int    `gorm:"not null" json:"customerId"`
	IsActive     bool   `gorm:"not null;default:true" json:"isActive"`

	Customer *Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
}
