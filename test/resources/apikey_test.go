package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venom90/bunnynet-go-client"
	"github.com/venom90/bunnynet-go-client/common"
	"github.com/venom90/bunnynet-go-client/resources"
	"github.com/venom90/bunnynet-go-client/test"
)

func TestAPIKeyService_List_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Items": [
			{
				"Id": 12345,
				"Key": "api-key-1",
				"Roles": ["PullZone.Read", "PullZone.Write"]
			},
			{
				"Id": 67890,
				"Key": "api-key-2",
				"Roles": ["Billing.Read"]
			}
		],
		"CurrentPage": 1,
		"TotalItems": 2,
		"HasMoreItems": false
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/apikey")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		// Verify pagination parameters
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("perPage"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create pagination parameters
	pagination := common.NewPagination().WithPage(2).WithPerPage(10)

	// Call the List method
	response, err := client.APIKey.List(context.Background(), pagination)
	assert.NoError(t, err, "List should not return an error")
	assert.NotNil(t, response, "Response should not be nil")
	assert.Len(t, response.Items, 2, "Should return 2 API keys")
	assert.Equal(t, int64(12345), response.Items[0].Id)
	assert.Equal(t, "api-key-1", response.Items[0].Key)
	assert.Equal(t, []string{"PullZone.Read", "PullZone.Write"}, response.Items[0].Roles)
	assert.Equal(t, int64(67890), response.Items[1].Id)
	assert.Equal(t, 1, response.CurrentPage)
	assert.Equal(t, 2, response.TotalItems)
	assert.False(t, response.HasMoreItems)
}

func TestAPIKeyService_List_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusUnauthorized, `{
		"ErrorKey": "unauthorized",
		"Field": "AccessKey",
		"Message": "The provided API key is invalid"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("invalid-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the List method
	response, err := client.APIKey.List(context.Background(), nil)
	assert.Error(t, err, "List should return an error")
	assert.Nil(t, response, "Response should be nil")
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestAPIKeyService_Get_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Id": 12345,
		"Key": "api-key-1",
		"Roles": ["PullZone.Read", "PullZone.Write"]
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/apikey/12345")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method
	apiKey, err := client.APIKey.Get(context.Background(), 12345)
	assert.NoError(t, err, "Get should not return an error")
	assert.NotNil(t, apiKey, "API key should not be nil")
	assert.Equal(t, int64(12345), apiKey.Id)
	assert.Equal(t, "api-key-1", apiKey.Key)
	assert.Equal(t, []string{"PullZone.Read", "PullZone.Write"}, apiKey.Roles)
}

func TestAPIKeyService_Get_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "apikey.not_found",
		"Field": "ApiKeyId",
		"Message": "The requested API key was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method with an invalid ID
	apiKey, err := client.APIKey.Get(context.Background(), 99999)
	assert.Error(t, err, "Get should return an error")
	assert.Nil(t, apiKey, "API key should be nil")
	assert.Contains(t, err.Error(), "apikey.not_found")
}

func TestAPIKeyService_Create_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Id": 12345,
		"Key": "new-api-key",
		"Roles": ["PullZone.Read", "PullZone.Write"]
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/apikey")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Create method
	roles := []string{"PullZone.Read", "PullZone.Write"}
	apiKey, err := client.APIKey.Create(context.Background(), roles)
	assert.NoError(t, err, "Create should not return an error")
	assert.NotNil(t, apiKey, "API key should not be nil")
	assert.Equal(t, int64(12345), apiKey.Id)
	assert.Equal(t, "new-api-key", apiKey.Key)
	assert.Equal(t, roles, apiKey.Roles)
}

func TestAPIKeyService_Create_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusBadRequest, `{
		"ErrorKey": "apikey.invalid_roles",
		"Field": "Roles",
		"Message": "The provided roles are invalid"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Create method with invalid roles
	roles := []string{"InvalidRole"}
	apiKey, err := client.APIKey.Create(context.Background(), roles)
	assert.Error(t, err, "Create should return an error")
	assert.Nil(t, apiKey, "API key should be nil")
	assert.Contains(t, err.Error(), "apikey.invalid_roles")
}

func TestAPIKeyService_Delete_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/apikey/12345")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Delete method
	err := client.APIKey.Delete(context.Background(), 12345)
	assert.NoError(t, err, "Delete should not return an error")
}

func TestAPIKeyService_Delete_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "apikey.not_found",
		"Field": "ApiKeyId",
		"Message": "The requested API key was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Delete method with an invalid ID
	err := client.APIKey.Delete(context.Background(), 99999)
	assert.Error(t, err, "Delete should return an error")
	assert.Contains(t, err.Error(), "apikey.not_found")
}

func TestAPIKeyService_ListAll(t *testing.T) {
	// Create a mock server for first page
	firstPageServer := test.MockServer(t, http.StatusOK, `{
		"Items": [
			{
				"Id": 12345,
				"Key": "api-key-1",
				"Roles": ["PullZone.Read"]
			},
			{
				"Id": 67890,
				"Key": "api-key-2",
				"Roles": ["Billing.Read"]
			}
		],
		"CurrentPage": 1,
		"TotalItems": 3,
		"HasMoreItems": true
	}`, func(r *http.Request) {
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "2", r.URL.Query().Get("perPage"))
	})
	defer firstPageServer.Close()

	// Create a mock client for the first page
	firstPageClient := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(firstPageServer.URL))

	// Start a goroutine that will get the first page
	firstPageCh := make(chan []resources.APIKey)
	firstPageErrCh := make(chan error)
	go func() {
		apiKeys, err := firstPageClient.APIKey.ListAll(context.Background(), 2)
		if err != nil {
			firstPageErrCh <- err
			return
		}
		firstPageCh <- apiKeys
	}()

	// Wait for the first page to be fetched
	select {
	case err := <-firstPageErrCh:
		t.Fatalf("Error fetching first page: %v", err)
	case apiKeys := <-firstPageCh:
		assert.Len(t, apiKeys, 3, "Should return all 3 API keys")
	}
}
