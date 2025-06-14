package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   message,
	}
}

// JSONResponse sends a JSON response
func JSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// ValidationError represents a validation error response
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(errors []ValidationError) Response {
	return Response{
		Success: false,
		Message: "Validation failed",
		Error:   "Validation failed",
		Data:    errors,
	}
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(data interface{}, pagination PaginationResponse) Response {
	return Response{
		Success: true,
		Message: "Data retrieved successfully",
		Data: map[string]interface{}{
			"data":       data,
			"pagination": pagination,
		},
	}
}
