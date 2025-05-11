// Package resources provides API resource implementations for the Bunny.net API client
package resources

import (
	"context"
	"net/http"

	"github.com/venom90/bunnynet-go-client/common"
	"github.com/venom90/bunnynet-go-client/internal"
)

// APIKey represents an API key in the Bunny.net API
type APIKey struct {
	// Id is the unique identifier of the API key
	Id int64 `json:"Id"`

	// Key is the API key value
	Key string `json:"Key"`

	// Roles is the list of roles assigned to the API key
	Roles []string `json:"Roles"`
}

// APIKeyService handles operations on API keys
type APIKeyService struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	userAgent string
}

// NewAPIKeyService creates a new APIKeyService
func NewAPIKeyService(client *http.Client, baseURL, apiKey, userAgent string) *APIKeyService {
	return &APIKeyService{
		client:    client,
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: userAgent,
	}
}

// SetAPIKey updates the API key used for authentication
func (s *APIKeyService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// List returns a paginated list of API keys
func (s *APIKeyService) List(ctx context.Context, pagination *common.Pagination) (*common.PaginatedResponse[APIKey], error) {
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, "/apikey", nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add pagination parameters
	if err := internal.AddQueryParams(req, pagination); err != nil {
		return nil, err
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var paginatedResponse common.PaginatedResponse[APIKey]
	if err := internal.ParsePaginatedResponse(resp, &paginatedResponse); err != nil {
		return nil, err
	}

	return &paginatedResponse, nil
}

// ListAll returns all API keys across all pages
func (s *APIKeyService) ListAll(ctx context.Context, perPage int) ([]APIKey, error) {
	if perPage <= 0 {
		perPage = common.DefaultPerPage
	}

	iterator := common.NewPageIterator(
		func(page, itemsPerPage int) (*common.PaginatedResponse[APIKey], error) {
			pagination := common.NewPagination().WithPage(page).WithPerPage(itemsPerPage)
			return s.List(ctx, pagination)
		},
		common.DefaultPage,
		perPage,
	)

	return iterator.AllItems()
}

// Get returns an API key by ID
func (s *APIKeyService) Get(ctx context.Context, id int64) (*APIKey, error) {
	path := "/apikey/" + internal.FormatInt64(id)
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var apiKey APIKey
	if err := internal.ParseResponse(resp, &apiKey); err != nil {
		return nil, err
	}

	return &apiKey, nil
}

// Create creates a new API key
func (s *APIKeyService) Create(ctx context.Context, roles []string) (*APIKey, error) {
	body := map[string]interface{}{
		"Roles": roles,
	}

	req, err := internal.NewRequest(http.MethodPost, s.baseURL, "/apikey", body, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var apiKey APIKey
	if err := internal.ParseResponse(resp, &apiKey); err != nil {
		return nil, err
	}

	return &apiKey, nil
}

// Delete deletes an API key
func (s *APIKeyService) Delete(ctx context.Context, id int64) error {
	path := "/apikey/" + internal.FormatInt64(id)
	req, err := internal.NewRequest(http.MethodDelete, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
