package model

type ReferenceDataType int

const (
	ReferenceDataTypePublisher ReferenceDataType = 0
	ReferenceDataTypeCondition ReferenceDataType = 1
	ReferenceDataTypeBookType  ReferenceDataType = 2
	ReferenceDataTypeGenre     ReferenceDataType = 3
)

type ReferenceDataItem struct {
	Entity
	DataType ReferenceDataType `gorm:"column:data_type" json:"dataType"`
	Text     string            `gorm:"column:text" json:"text"`
}

func (ReferenceDataItem) TableName() string {
	return "reference_data"
}
