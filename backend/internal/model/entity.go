package model

import "time"

type Entity struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedBy string    `json:"createdBy" gorm:"default:System"`
	CreatedOn time.Time `json:"createdOn" gorm:"autoCreateTime"`
	UpdatedOn time.Time `json:"updatedOn" gorm:"autoUpdateTime"`
}
