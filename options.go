package bunnynet

import (
	"net/http"
	"time"
)

// Option is a function that configures a Client
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for API requests
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the base URL for API requests
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithUserAgent sets the User-Agent header for API requests
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// WithTimeout sets the timeout for API requests
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}
