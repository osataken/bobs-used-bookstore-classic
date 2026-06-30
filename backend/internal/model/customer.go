package model

import "time"

type Customer struct {
	Entity
	Sub         string     `gorm:"type:varchar(450);uniqueIndex" json:"sub"`
	Username    string     `gorm:"type:text" json:"username"`
	FirstName   string     `gorm:"type:text" json:"firstName"`
	LastName    string     `gorm:"type:text" json:"lastName"`
	Email       string     `gorm:"type:text" json:"email"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	Phone       string     `gorm:"type:text" json:"phone"`
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}
