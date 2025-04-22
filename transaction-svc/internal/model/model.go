package model

type WebResponse[T any] struct {
	Success      bool          `json:"success"`
	Data         T             `json:"data,omitempty"`
	Message      string        `json:"message,omitempty"`
	PageMetadata *PageMetadata `json:"pagination,omitempty"`
}

type PageMetadata struct {
	Page            int    `json:"page"`
	Size            int    `json:"size"`
	Offset          int    `json:"offset"`
	TotalItem       int64  `json:"total_item"`
	TotalPage       int64  `json:"total_page"`
	HasNext         bool   `json:"has_next"`
	HasPrevious     bool   `json:"has_previous"`
	NextPageURL     string `json:"next_page_url"`
	PreviousPageURL string `json:"previous_page_url"`
}

type BodyParseErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ValidationErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
