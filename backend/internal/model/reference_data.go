package model

type ReferenceDataItem struct {
	Entity
	DataType ReferenceDataType `gorm:"not null" json:"dataType"`
	Text     string            `gorm:"type:text" json:"text"`
}

func (ReferenceDataItem) TableName() string {
	return "reference_data"
}
