package model

import "time"

// Entity is the base struct for all domain entities
type Entity struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy string    `gorm:"type:varchar(255);default:System" json:"createdBy"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"createdOn"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"updatedOn"`
}
