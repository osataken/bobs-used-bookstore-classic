package dto

// AuthUserResponse represents the current authenticated user
type AuthUserResponse struct {
	Sub       string `json:"sub"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsAdmin   bool   `json:"isAdmin"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Items       interface{} `json:"items"`
	PageIndex   int         `json:"pageIndex"`
	TotalPages  int         `json:"totalPages"`
	TotalCount  int         `json:"totalCount"`
	HasNext     bool        `json:"hasNext"`
	HasPrevious bool        `json:"hasPrevious"`
}
