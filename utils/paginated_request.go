package utils

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
