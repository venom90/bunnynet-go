package common

import (
	"fmt"
	"strconv"
)

const (
	// DefaultPage is the default page number for paginated requests
	DefaultPage = 1

	// DefaultPerPage is the default number of items per page
	DefaultPerPage = 100

	// MaxPerPage is the maximum number of items per page
	MaxPerPage = 1000
)

// Pagination represents pagination parameters for API requests
type Pagination struct {
	// Page is the current page number (1-based)
	Page int `url:"page,omitempty"`

	// PerPage is the number of items per page
	PerPage int `url:"perPage,omitempty"`
}

// NewPagination creates a new Pagination with default values
func NewPagination() *Pagination {
	return &Pagination{
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
	}
}

// WithPage sets the page number and returns the Pagination for chaining
func (p *Pagination) WithPage(page int) *Pagination {
	if page < 1 {
		page = DefaultPage
	}
	p.Page = page
	return p
}

// WithPerPage sets the number of items per page and returns the Pagination for chaining
func (p *Pagination) WithPerPage(perPage int) *Pagination {
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}
	p.PerPage = perPage
	return p
}

// ToQueryParams converts the Pagination to query parameters
func (p *Pagination) ToQueryParams() map[string]string {
	if p == nil {
		return map[string]string{}
	}

	params := map[string]string{}

	if p.Page > 0 {
		params["page"] = strconv.Itoa(p.Page)
	}

	if p.PerPage > 0 {
		params["perPage"] = strconv.Itoa(p.PerPage)
	}

	return params
}

// String returns a string representation of the Pagination
func (p *Pagination) String() string {
	if p == nil {
		return "Pagination{}"
	}
	return fmt.Sprintf("Pagination{Page: %d, PerPage: %d}", p.Page, p.PerPage)
}

// PageInfo represents pagination information from a response
type PageInfo struct {
	// CurrentPage is the current page number
	CurrentPage int

	// TotalItems is the total number of items
	TotalItems int

	// HasMoreItems indicates whether there are more pages
	HasMoreItems bool
}

// TotalPages calculates the total number of pages based on total items and items per page
func (p *PageInfo) TotalPages(perPage int) int {
	if perPage <= 0 {
		return 0
	}

	totalPages := p.TotalItems / perPage
	if p.TotalItems%perPage > 0 {
		totalPages++
	}

	return totalPages
}

// NextPage returns the next page number, or 0 if there are no more pages
func (p *PageInfo) NextPage() int {
	if !p.HasMoreItems {
		return 0
	}

	return p.CurrentPage + 1
}

// PreviousPage returns the previous page number, or 0 if there is no previous page
func (p *PageInfo) PreviousPage() int {
	if p.CurrentPage <= 1 {
		return 0
	}

	return p.CurrentPage - 1
}

// IsFirstPage returns true if the current page is the first page
func (p *PageInfo) IsFirstPage() bool {
	return p.CurrentPage <= 1
}

// IsLastPage returns true if the current page is the last page
func (p *PageInfo) IsLastPage() bool {
	return !p.HasMoreItems
}

// PageInfoFromResponse extracts pagination information from a PaginatedResponse
func PageInfoFromResponse[T any](response *PaginatedResponse[T]) *PageInfo {
	if response == nil {
		return nil
	}

	return &PageInfo{
		CurrentPage:  response.CurrentPage,
		TotalItems:   response.TotalItems,
		HasMoreItems: response.HasMoreItems,
	}
}

// PageIterator is a utility for iterating through pages of results
type PageIterator[T any] struct {
	// client is the function that fetches a page of results
	client func(page, perPage int) (*PaginatedResponse[T], error)

	// pagination is the current pagination state
	pagination *Pagination

	// currentResponse is the current page of results
	currentResponse *PaginatedResponse[T]

	// err is the last error that occurred
	err error
}

// NewPageIterator creates a new PageIterator
func NewPageIterator[T any](
	client func(page, perPage int) (*PaginatedResponse[T], error),
	page, perPage int,
) *PageIterator[T] {
	return &PageIterator[T]{
		client:     client,
		pagination: NewPagination().WithPage(page).WithPerPage(perPage),
	}
}

// Next fetches the next page of results
// Returns true if there are more results, false otherwise
func (i *PageIterator[T]) Next() bool {
	// If we've already encountered an error, don't continue
	if i.err != nil {
		return false
	}

	// If we've already fetched a page and there are no more items, don't continue
	if i.currentResponse != nil && !i.currentResponse.HasMoreItems {
		return false
	}

	// Fetch the next page
	response, err := i.client(i.pagination.Page, i.pagination.PerPage)
	if err != nil {
		i.err = err
		return false
	}

	// Update the current response
	i.currentResponse = response

	// Increment the page for the next fetch
	i.pagination.Page++

	// Return true if we have items in this page
	return len(response.Items) > 0
}

// Items returns the items in the current page
func (i *PageIterator[T]) Items() []T {
	if i.currentResponse == nil {
		return nil
	}

	return i.currentResponse.Items
}

// Error returns the last error that occurred
func (i *PageIterator[T]) Error() error {
	return i.err
}

// PageInfo returns pagination information for the current page
func (i *PageIterator[T]) PageInfo() *PageInfo {
	return PageInfoFromResponse(i.currentResponse)
}

// Reset resets the iterator to the first page
func (i *PageIterator[T]) Reset() {
	i.pagination.Page = DefaultPage
	i.currentResponse = nil
	i.err = nil
}

// AllItems fetches all items across all pages
// Warning: This may result in a large number of API requests and items
func (i *PageIterator[T]) AllItems() ([]T, error) {
	// Reset the iterator to ensure we start from the first page
	i.Reset()

	var allItems []T

	// Fetch all pages
	for i.Next() {
		allItems = append(allItems, i.Items()...)
	}

	// Check if we encountered an error
	if i.Error() != nil {
		return nil, i.Error()
	}

	return allItems, nil
}
