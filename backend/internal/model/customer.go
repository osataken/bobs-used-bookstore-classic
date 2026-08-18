package model

import "time"

// Customer entity
type Customer struct {
	Entity
	Sub         string     `gorm:"size:450;uniqueIndex;not null" json:"sub"`
	Username    string     `json:"username"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Email       string     `json:"email"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	Phone       string     `json:"phone"`
}

func (Customer) TableName() string {
	return "Customer"
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}
