package bunnynet

// Package bunnynet provides a client for interacting with the Bunny.net API.
import (
	"net/http"
	"time"

	"github.com/venom90/bunnynet-go/resources"
)

const (
	// DefaultBaseURL is the default base URL for the Bunny.net API
	DefaultBaseURL = "https://api.bunny.net"
	// DefaultTimeout is the default timeout for API requests
	DefaultTimeout = 30 * time.Second
	// DefaultUserAgent is the default User-Agent header value
	DefaultUserAgent = "bunnynet-go/1.0.0"
)

// Client represents a Bunny.net API client
type Client struct {
	// HTTP client used to communicate with the API
	httpClient *http.Client

	// Base URL for API requests
	BaseURL string

	// API key for authenticating requests
	apiKey string

	// User agent used when communicating with the Bunny.net API
	UserAgent string

	// Resources
	Country  *resources.CountryService
	APIKey   *resources.APIKeyService
	DNSZone  *resources.DNSZoneService
	PullZone *resources.PullZoneService
	Purge    *resources.PurgeService
}

// NewClient returns a new Bunny.net API client
func NewClient(apiKey string, options ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		BaseURL:   DefaultBaseURL,
		apiKey:    apiKey,
		UserAgent: DefaultUserAgent,
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	// Initialize services
	client.Country = resources.NewCountryService(client.httpClient, client.BaseURL, client.apiKey, client.UserAgent)
	client.APIKey = resources.NewAPIKeyService(client.httpClient, client.BaseURL, client.apiKey, client.UserAgent)
	client.DNSZone = resources.NewDNSZoneService(client.httpClient, client.BaseURL, client.apiKey, client.UserAgent)
	client.PullZone = resources.NewPullZoneService(client.httpClient, client.BaseURL, client.apiKey, client.UserAgent)
	client.Purge = resources.NewPurgeService(client.httpClient, client.BaseURL, client.apiKey, client.UserAgent)

	return client
}

// SetAPIKey updates the API key used for authentication
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey

	// Update API key for all services
	c.Country.SetAPIKey(apiKey)
	c.APIKey.SetAPIKey(apiKey)
	c.DNSZone.SetAPIKey(apiKey)
	c.PullZone.SetAPIKey(apiKey)
	c.Purge.SetAPIKey(apiKey)
}
