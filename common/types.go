// Package common provides common types and utilities for the Bunny.net API client
package common

// PaginationOptions contains options for paginated API requests
type PaginationOptions struct {
	// Page is the page number to retrieve (starting from 1)
	Page int `url:"page,omitempty"`

	// PerPage is the number of items per page
	PerPage int `url:"perPage,omitempty"`
}

// PaginatedResponse is a generic response type for paginated API responses
type PaginatedResponse[T any] struct {
	// Items is the list of items in the current page
	Items []T `json:"Items"`

	// CurrentPage is the current page number
	CurrentPage int `json:"CurrentPage"`

	// TotalItems is the total number of items across all pages
	TotalItems int `json:"TotalItems"`

	// HasMoreItems indicates whether there are more pages of items
	HasMoreItems bool `json:"HasMoreItems"`
}

// RequestParams interface represents a type that can be converted to URL query parameters
type RequestParams interface {
	ToQueryParams() map[string]string
}
