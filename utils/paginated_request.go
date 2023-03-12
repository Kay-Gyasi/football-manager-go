package utils

type PaginationRequest struct {
	PageNumber int `json:"page" validate:"required"`
	PageSize   int `json:"pageSize" validate:"required"`
}
