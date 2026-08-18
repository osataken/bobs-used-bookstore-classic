package model

import "time"

// Entity is the base model for all domain entities
type Entity struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy string    `gorm:"default:System" json:"createdBy"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"createdOn"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"updatedOn"`
	Version   int       `gorm:"default:1" json:"-"`
}
