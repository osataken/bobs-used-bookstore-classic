package model

// Address represents a customer shipping address
type Address struct {
	Entity
	AddressLine1 string `gorm:"type:varchar(255);not null" json:"addressLine1"`
	AddressLine2 string `gorm:"type:varchar(255)" json:"addressLine2"`
	City         string `gorm:"type:varchar(100);not null" json:"city"`
	State        string `gorm:"type:varchar(100);not null" json:"state"`
	Country      string `gorm:"type:varchar(100);not null" json:"country"`
	ZipCode      string `gorm:"type:varchar(20);not null" json:"zipCode"`
	CustomerID   int    `gorm:"not null" json:"customerId"`
	IsActive     bool   `gorm:"default:true" json:"isActive"`

	// Navigation properties
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Address) TableName() string {
	return "Address"
}
