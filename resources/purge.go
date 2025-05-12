// Package resources provides API resource implementations for the Bunny.net API client
package resources

import (
	"context"
	"net/http"

	"github.com/venom90/bunnynet-go/internal"
)

// PurgeService handles URL purging operations
type PurgeService struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	userAgent string
}

// NewPurgeService creates a new PurgeService
func NewPurgeService(client *http.Client, baseURL, apiKey, userAgent string) *PurgeService {
	return &PurgeService{
		client:    client,
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: userAgent,
	}
}

// SetAPIKey updates the API key used for authentication
func (s *PurgeService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// PurgeOptions represents the options for purging a URL
type PurgeOptions struct {
	// URL is the URL that will be purged from cache
	URL string `url:"url" json:"url"`

	// Async determines if the call should wait for the purge logic to complete
	Async bool `url:"async,omitempty" json:"async,omitempty"`
}

// ToQueryParams converts the PurgeOptions to query parameters
func (o *PurgeOptions) ToQueryParams() map[string]string {
	params := make(map[string]string)

	params["url"] = o.URL

	if o.Async {
		params["async"] = "true"
	}

	return params
}

// PurgeURL purges a URL from the CDN cache
func (s *PurgeService) PurgeURL(ctx context.Context, options PurgeOptions) error {
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, "/purge", nil, s.apiKey, s.userAgent)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	// Add query parameters
	if err := internal.AddQueryParams(req, options); err != nil {
		return err
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Purge purges the specified URL from the CDN cache
// Shorthand for PurgeURL with simpler parameters
func (s *PurgeService) Purge(ctx context.Context, url string, async bool) error {
	return s.PurgeURL(ctx, PurgeOptions{
		URL:   url,
		Async: async,
	})
}
