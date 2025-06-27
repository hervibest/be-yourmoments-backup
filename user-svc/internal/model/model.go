package model

type WebResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
	// Token        *TokenResponse `json:"token,omitempty"`
	PageMetadata *PageMetadata `json:"pagination,omitempty"`
}

type PageMetadata struct {
	Page            int
	Size            int
	Offset          int
	TotalItem       int64
	TotalPage       int64
	HasNext         bool
	HasPrevious     bool
	NextPageURL     string
	PreviousPageURL string
}

type ValidationError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type BodyParseErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
