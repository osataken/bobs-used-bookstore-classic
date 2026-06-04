package model

type Customer struct {
	Entity
	Sub       string `gorm:"column:sub;size:450;uniqueIndex" json:"sub"`
	Username  string `gorm:"column:username" json:"username"`
	FirstName string `gorm:"column:first_name" json:"firstName"`
	LastName  string `gorm:"column:last_name" json:"lastName"`
	Email     string `gorm:"column:email" json:"email"`
	Phone     string `gorm:"column:phone" json:"phone"`
}

func (Customer) TableName() string {
	return "customer"
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}
