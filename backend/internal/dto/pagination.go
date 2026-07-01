package dto

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
	PageIndex  int         `json:"pageIndex"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}
