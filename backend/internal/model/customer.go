package model

// Customer represents a bookstore customer
type Customer struct {
	Entity
	Sub       string `gorm:"type:varchar(450);uniqueIndex;not null" json:"sub"`
	Username  string `gorm:"type:varchar(255)" json:"username"`
	FirstName string `gorm:"type:varchar(255)" json:"firstName"`
	LastName  string `gorm:"type:varchar(255)" json:"lastName"`
	Email     string `gorm:"type:varchar(255)" json:"email"`
	Phone     string `gorm:"type:varchar(50)" json:"phone"`
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}

func (Customer) TableName() string {
	return "Customer"
}
