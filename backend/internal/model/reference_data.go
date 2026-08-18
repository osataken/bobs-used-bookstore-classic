package model

// ReferenceDataItem entity
type ReferenceDataItem struct {
	Entity
	DataType ReferenceDataType `gorm:"not null" json:"dataType"`
	Text     string            `gorm:"not null" json:"text"`
}

func (ReferenceDataItem) TableName() string {
	return "ReferenceData"
}
