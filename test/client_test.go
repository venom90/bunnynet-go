package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venom90/bunnynet-go"
)

func TestNewClient(t *testing.T) {
	// Test default client creation
	client := bunnynet.NewClient("test-api-key")
	assert.NotNil(t, client, "Client should not be nil")
	assert.Equal(t, bunnynet.DefaultBaseURL, client.BaseURL, "BaseURL should be the default")
	assert.Equal(t, bunnynet.DefaultUserAgent, client.UserAgent, "UserAgent should be the default")

	// Test client with options
	customClient := bunnynet.NewClient(
		"test-api-key",
		bunnynet.WithBaseURL("https://custom-api.bunny.net"),
		bunnynet.WithUserAgent("custom-agent/1.0"),
		bunnynet.WithTimeout(60*time.Second),
	)
	assert.NotNil(t, customClient, "Client should not be nil")
	assert.Equal(t, "https://custom-api.bunny.net", customClient.BaseURL, "BaseURL should be customized")
	assert.Equal(t, "custom-agent/1.0", customClient.UserAgent, "UserAgent should be customized")
}

func TestSetAPIKey(t *testing.T) {
	client := bunnynet.NewClient("initial-api-key")
	client.SetAPIKey("new-api-key")

	// Create a mock server to verify the API key is used
	server := MockServer(t, http.StatusOK, `[]`, func(r *http.Request) {
		AssertRequestHasHeader(t, r, "AccessKey", "new-api-key")
	})
	defer server.Close()

	// Override the base URL to use the mock server
	client.BaseURL = server.URL

	// Make a request to verify the API key is used
	_, err := client.Country.List(nil)
	assert.NoError(t, err, "Request should succeed")
}
