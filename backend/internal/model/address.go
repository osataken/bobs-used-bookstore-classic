package model

type Address struct {
	Entity
	AddressLine1 string   `gorm:"column:address_line1" json:"addressLine1"`
	AddressLine2 string   `gorm:"column:address_line2" json:"addressLine2"`
	City         string   `gorm:"column:city" json:"city"`
	State        string   `gorm:"column:state" json:"state"`
	Country      string   `gorm:"column:country" json:"country"`
	ZipCode      string   `gorm:"column:zip_code" json:"zipCode"`
	CustomerID   int      `gorm:"column:customer_id" json:"customerId"`
	Customer     Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	IsActive     bool     `gorm:"column:is_active;default:true" json:"isActive"`
}

func (Address) TableName() string {
	return "address"
}
