package internal

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/venom90/bunnynet-go-client/common"
)

// ParseResponse parses the response body into the given value
func ParseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	// For responses with no content
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// Check if we need to parse the response body
	if v == nil {
		return nil
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		return common.ParseErrorResponse(resp)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return common.NewClientError("failed to read response body", err)
	}

	// Parse the response body
	err = json.Unmarshal(body, v)
	if err != nil {
		return common.NewClientError("failed to parse response body", err)
	}

	return nil
}

// ParsePaginatedResponse parses a paginated response into the given value
func ParsePaginatedResponse[T any](resp *http.Response, v *common.PaginatedResponse[T]) error {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return common.NewClientError("failed to read response body", err)
	}

	// Parse the response body
	err = json.Unmarshal(body, v)
	if err != nil {
		return common.NewClientError("failed to parse paginated response", err)
	}

	return nil
}
