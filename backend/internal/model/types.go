package model

// PaginatedList represents a paginated result set
type PaginatedList struct {
	Items      interface{} `json:"items"`
	PageIndex  int         `json:"pageIndex"`
	PageSize   int         `json:"pageSize"`
	TotalCount int64       `json:"totalCount"`
	TotalPages int         `json:"totalPages"`
}

// BookStatistics for admin dashboard
type BookStatistics struct {
	LowStock   int64 `json:"lowStock"`
	OutOfStock int64 `json:"outOfStock"`
	StockTotal int64 `json:"stockTotal"`
}

// OrderStatistics for admin dashboard
type OrderStatistics struct {
	Pending   int64 `json:"pending"`
	Ordered   int64 `json:"ordered"`
	Shipped   int64 `json:"shipped"`
	Delivered int64 `json:"delivered"`
	Cancelled int64 `json:"cancelled"`
}

// OfferStatistics for admin dashboard
type OfferStatistics struct {
	PendingApproval int64 `json:"pendingApproval"`
	Approved        int64 `json:"approved"`
	Received        int64 `json:"received"`
	Paid            int64 `json:"paid"`
	Rejected        int64 `json:"rejected"`
}
