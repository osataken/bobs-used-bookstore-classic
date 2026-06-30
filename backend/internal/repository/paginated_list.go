package repository

import "math"

type PaginatedList[T any] struct {
	Items      []T `json:"items"`
	PageIndex  int `json:"pageIndex"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

func NewPaginatedList[T any](items []T, totalCount, pageIndex, pageSize int) PaginatedList[T] {
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return PaginatedList[T]{
		Items:      items,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
