package dto

// Base represents the base response object.
// swagger:model
type Base struct {
	// Indicates whether the operation was successful.
	Success bool `json:"success"`
	// Error details if the operation failed.
	Error *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse represents an error response.
// swagger:model
type ErrorResponse struct {
	// The error code.
	Code string `json:"code"`
	// The error message.
	Message string `json:"message"`
}

// GetErrorResponse generates an error response.
// swagger:response ErrorResponse
func GetErrorResponse(code string, message string) *Base {
	return &Base{
		Success: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
		},
	}
}
