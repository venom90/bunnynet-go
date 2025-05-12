package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venom90/bunnynet-go"
	"github.com/venom90/bunnynet-go/resources"
	"github.com/venom90/bunnynet-go/test"
)

func TestPurgeService_PurgeURL_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/purge")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		assert.Equal(t, "https://example.com/file.jpg", r.URL.Query().Get("url"))
		assert.Equal(t, "true", r.URL.Query().Get("async"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the PurgeURL method with the PurgeOptions
	options := resources.PurgeOptions{
		URL:   "https://example.com/file.jpg",
		Async: true,
	}
	err := client.Purge.PurgeURL(context.Background(), options)
	assert.NoError(t, err, "PurgeURL should not return an error")
}

func TestPurgeService_PurgeURL_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusUnauthorized, `{
		"ErrorKey": "unauthorized",
		"Field": "AccessKey",
		"Message": "The provided API key is invalid"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("invalid-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the PurgeURL method
	options := resources.PurgeOptions{
		URL: "https://example.com/file.jpg",
	}
	err := client.Purge.PurgeURL(context.Background(), options)
	assert.Error(t, err, "PurgeURL should return an error")
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestPurgeService_Purge_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/purge")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		assert.Equal(t, "https://example.com/file.jpg", r.URL.Query().Get("url"))
		assert.Equal(t, "", r.URL.Query().Get("async")) // Should be false/not present
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the simplified Purge method
	err := client.Purge.Purge(context.Background(), "https://example.com/file.jpg", false)
	assert.NoError(t, err, "Purge should not return an error")
}
