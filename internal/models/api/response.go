package model

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ListResponse[T any] struct {
	Skip       int `json:"skip"`
	Take       int `json:"take"`
	TotalCount int `json:"totalCount"`
	Data       T   `json:"data"`
}
