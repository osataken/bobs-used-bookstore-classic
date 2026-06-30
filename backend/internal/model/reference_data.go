package model

// ReferenceDataItem represents a reference data entry (publisher, condition, book type, genre)
type ReferenceDataItem struct {
	Entity
	DataType ReferenceDataType `gorm:"not null" json:"dataType"`
	Text     string            `gorm:"type:varchar(255);not null" json:"text"`
}

func (ReferenceDataItem) TableName() string {
	return "ReferenceData"
}
