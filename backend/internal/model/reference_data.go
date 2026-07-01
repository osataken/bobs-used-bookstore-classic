package model

type ReferenceDataItem struct {
	Entity
	DataType ReferenceDataType `json:"dataType" gorm:"not null"`
	Text     string            `json:"text" gorm:"not null"`
}
