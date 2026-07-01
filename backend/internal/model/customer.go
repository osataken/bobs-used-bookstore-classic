package model

import "time"

type Customer struct {
	Entity
	Sub         string     `json:"sub" gorm:"uniqueIndex;size:450;not null"`
	Username    string     `json:"username"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Email       string     `json:"email"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	Phone       string     `json:"phone"`
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}
