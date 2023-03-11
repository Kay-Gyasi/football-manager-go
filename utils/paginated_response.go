package utils

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	TotalItems int         `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
}
