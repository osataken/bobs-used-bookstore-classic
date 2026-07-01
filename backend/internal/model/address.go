package model

type Address struct {
	Entity
	AddressLine1 string `json:"addressLine1" gorm:"not null"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" gorm:"not null"`
	State        string `json:"state" gorm:"not null"`
	Country      string `json:"country" gorm:"not null"`
	ZipCode      string `json:"zipCode" gorm:"not null"`
	CustomerID   uint   `json:"customerId" gorm:"not null"`
	IsActive     bool   `json:"isActive" gorm:"default:true"`
}
