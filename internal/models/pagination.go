package models

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"per_page"`
	Offset     int `json:"offset"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}
