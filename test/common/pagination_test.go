package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venom90/bunnynet-go/common"
)

func TestPagination(t *testing.T) {
	// Test default pagination
	p := common.NewPagination()
	assert.Equal(t, common.DefaultPage, p.Page)
	assert.Equal(t, common.DefaultPerPage, p.PerPage)

	// Test with custom values
	p.WithPage(5).WithPerPage(50)
	assert.Equal(t, 5, p.Page)
	assert.Equal(t, 50, p.PerPage)

	// Test with invalid values
	p.WithPage(-1).WithPerPage(-10)
	assert.Equal(t, common.DefaultPage, p.Page)
	assert.Equal(t, common.DefaultPerPage, p.PerPage)

	// Test with values exceeding limits
	p.WithPerPage(common.MaxPerPage + 100)
	assert.Equal(t, common.MaxPerPage, p.PerPage)

	// Test ToQueryParams
	params := p.ToQueryParams()
	assert.Equal(t, "1", params["page"])
	assert.Equal(t, "1000", params["perPage"])

	// Test String representation
	assert.Contains(t, p.String(), "Page: 1")
	assert.Contains(t, p.String(), "PerPage: 1000")
}

func TestPageInfo(t *testing.T) {
	// Create a page info
	pageInfo := &common.PageInfo{
		CurrentPage:  2,
		TotalItems:   25,
		HasMoreItems: true,
	}

	// Test TotalPages
	assert.Equal(t, 5, pageInfo.TotalPages(5))
	assert.Equal(t, 3, pageInfo.TotalPages(10))
	assert.Equal(t, 0, pageInfo.TotalPages(0))

	// Test NextPage
	assert.Equal(t, 3, pageInfo.NextPage())

	// Test with no more items
	pageInfo.HasMoreItems = false
	assert.Equal(t, 0, pageInfo.NextPage())

	// Test PreviousPage
	assert.Equal(t, 1, pageInfo.PreviousPage())

	// Test first page
	pageInfo.CurrentPage = 1
	assert.Equal(t, 0, pageInfo.PreviousPage())
	assert.True(t, pageInfo.IsFirstPage())

	// Test last page
	assert.True(t, pageInfo.IsLastPage())
}

func TestPageInfoFromResponse(t *testing.T) {
	// Create a paginated response
	response := &common.PaginatedResponse[string]{
		Items:        []string{"item1", "item2"},
		CurrentPage:  3,
		TotalItems:   50,
		HasMoreItems: true,
	}

	// Extract page info
	pageInfo := common.PageInfoFromResponse(response)
	assert.Equal(t, 3, pageInfo.CurrentPage)
	assert.Equal(t, 50, pageInfo.TotalItems)
	assert.True(t, pageInfo.HasMoreItems)

	// Test with nil response
	nilPageInfo := common.PageInfoFromResponse[string](nil)
	assert.Nil(t, nilPageInfo)
}

type mockItem struct {
	ID   int
	Name string
}

func TestPageIterator(t *testing.T) {
	// Mock data for testing
	mockPages := []common.PaginatedResponse[mockItem]{
		{
			Items: []mockItem{
				{ID: 1, Name: "Item 1"},
				{ID: 2, Name: "Item 2"},
			},
			CurrentPage:  1,
			TotalItems:   5,
			HasMoreItems: true,
		},
		{
			Items: []mockItem{
				{ID: 3, Name: "Item 3"},
				{ID: 4, Name: "Item 4"},
			},
			CurrentPage:  2,
			TotalItems:   5,
			HasMoreItems: true,
		},
		{
			Items: []mockItem{
				{ID: 5, Name: "Item 5"},
			},
			CurrentPage:  3,
			TotalItems:   5,
			HasMoreItems: false,
		},
	}

	// Mock client function
	clientFn := func(page, perPage int) (*common.PaginatedResponse[mockItem], error) {
		if page < 1 || page > len(mockPages) {
			return nil, errors.New("page out of range")
		}
		return &mockPages[page-1], nil
	}

	// Create iterator
	iterator := common.NewPageIterator(clientFn, 1, 2)

	// Test first page
	assert.True(t, iterator.Next())
	assert.Equal(t, mockPages[0].Items, iterator.Items())
	assert.Equal(t, 1, iterator.PageInfo().CurrentPage)
	assert.Equal(t, 5, iterator.PageInfo().TotalItems)
	assert.True(t, iterator.PageInfo().HasMoreItems)

	// Test second page
	assert.True(t, iterator.Next())
	assert.Equal(t, mockPages[1].Items, iterator.Items())
	assert.Equal(t, 2, iterator.PageInfo().CurrentPage)

	// Test third page
	assert.True(t, iterator.Next())
	assert.Equal(t, mockPages[2].Items, iterator.Items())
	assert.Equal(t, 3, iterator.PageInfo().CurrentPage)
	assert.False(t, iterator.PageInfo().HasMoreItems)

	// Test no more pages
	assert.False(t, iterator.Next())
	assert.Nil(t, iterator.Error())

	// Test error handling
	errorIterator := common.NewPageIterator(
		func(page, perPage int) (*common.PaginatedResponse[mockItem], error) {
			return nil, errors.New("test error")
		},
		1, 10,
	)
	assert.False(t, errorIterator.Next())
	assert.Error(t, errorIterator.Error())
	assert.Equal(t, "test error", errorIterator.Error().Error())

	// Test AllItems
	iterator.Reset()
	allItems, err := iterator.AllItems()
	assert.NoError(t, err)
	assert.Len(t, allItems, 5)
	assert.Equal(t, 1, allItems[0].ID)
	assert.Equal(t, 5, allItems[4].ID)

	// Test AllItems with error
	errorIterator.Reset()
	_, err = errorIterator.AllItems()
	assert.Error(t, err)
}
