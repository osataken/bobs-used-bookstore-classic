package model

import "time"

type Entity struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy string    `gorm:"default:System" json:"createdBy"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"createdOn"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"updatedOn"`
}
