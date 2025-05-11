package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents an error response from the Bunny.net API
type ErrorResponse struct {
	// ErrorKey is a machine-readable error code
	ErrorKey string `json:"ErrorKey"`

	// Field indicates which field caused the error
	Field string `json:"Field"`

	// Message is a human-readable error message
	Message string `json:"Message"`

	// StatusCode is the HTTP status code of the response
	StatusCode int `json:"-"`
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("[%d] %s: %s (%s)", e.StatusCode, e.ErrorKey, e.Message, e.Field)
}

// ParseErrorResponse attempts to parse an error response from the Bunny.net API
func ParseErrorResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	errorResponse := &ErrorResponse{
		StatusCode: resp.StatusCode,
	}

	// Try to decode the error response
	err := json.NewDecoder(resp.Body).Decode(errorResponse)
	if err != nil {
		// If we can't decode the error response, return a generic error
		return fmt.Errorf("bunnynet API error: status code %d", resp.StatusCode)
	}

	return errorResponse
}

// ClientError represents an error that occurred in the client
type ClientError struct {
	Message string
	Err     error
}

// Error implements the error interface
func (e *ClientError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ClientError) Unwrap() error {
	return e.Err
}

// NewClientError creates a new ClientError
func NewClientError(message string, err error) *ClientError {
	return &ClientError{
		Message: message,
		Err:     err,
	}
}
