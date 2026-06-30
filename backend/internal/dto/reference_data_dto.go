package dto

// ReferenceDataResponse represents a reference data item
type ReferenceDataResponse struct {
	ID       int    `json:"id"`
	DataType int    `json:"dataType"`
	TypeText string `json:"typeText"`
	Text     string `json:"text"`
}

// ReferenceDataCreateRequest for creating reference data
type ReferenceDataCreateRequest struct {
	DataType int    `json:"dataType" binding:"required"`
	Text     string `json:"text" binding:"required"`
}

// ReferenceDataUpdateRequest for updating reference data
type ReferenceDataUpdateRequest struct {
	DataType int    `json:"dataType" binding:"required"`
	Text     string `json:"text" binding:"required"`
}

// ReferenceDataFiltersRequest for filtering reference data
type ReferenceDataFiltersRequest struct {
	DataType *int `form:"dataType"`
}
