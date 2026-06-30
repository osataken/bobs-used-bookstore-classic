package dto

// AddressResponse represents an address in API responses
type AddressResponse struct {
	ID           int    `json:"id"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	ZipCode      string `json:"zipCode"`
}

// AddressCreateRequest for creating an address
type AddressCreateRequest struct {
	AddressLine1 string `json:"addressLine1" binding:"required"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	Country      string `json:"country" binding:"required"`
	ZipCode      string `json:"zipCode" binding:"required"`
}

// AddressUpdateRequest for updating an address
type AddressUpdateRequest struct {
	AddressLine1 string `json:"addressLine1" binding:"required"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	Country      string `json:"country" binding:"required"`
	ZipCode      string `json:"zipCode" binding:"required"`
}
