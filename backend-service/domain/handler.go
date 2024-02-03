package domain

// SuccessResponse represents a standard success response
// swagger:model
type SuccessResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents a standard error response
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}
