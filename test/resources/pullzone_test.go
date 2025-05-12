package resources

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venom90/bunnynet-go-client"
	"github.com/venom90/bunnynet-go-client/resources"
	"github.com/venom90/bunnynet-go-client/test"
)

func TestPullZoneService_List_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Items": [
			{
				"Id": 12345,
				"Name": "test-zone-1",
				"OriginUrl": "https://example.com",
				"Enabled": true,
				"Hostnames": [
					{
						"Id": 111,
						"Value": "cdn.example.com",
						"ForceSSL": true,
						"IsSystemHostname": false,
						"HasCertificate": true
					}
				],
				"AllowedReferrers": ["example.com"],
				"BlockedReferrers": ["badsite.com"],
				"EnableGeoZoneUS": true,
				"EnableGeoZoneEU": true,
				"Type": 0
			},
			{
				"Id": 67890,
				"Name": "test-zone-2",
				"OriginUrl": "https://example2.com",
				"Enabled": true,
				"Hostnames": [],
				"EnableGeoZoneUS": true,
				"EnableGeoZoneEU": false,
				"Type": 1
			}
		],
		"CurrentPage": 1,
		"TotalItems": 2,
		"HasMoreItems": false
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/pullzone")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		// Verify pagination and search parameters
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("perPage"))
		assert.Equal(t, "test", r.URL.Query().Get("search"))
		assert.Equal(t, "true", r.URL.Query().Get("includeCertificate"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create the pull zone options with invalid data
	options := resources.AddPullZoneOptions{
		Name:      "test-zone-1",
		OriginUrl: "invalid-url",
	}

	// Call the Add method
	pullZone, err := client.PullZone.Add(context.Background(), options)
	assert.Error(t, err, "Add should return an error")
	assert.Nil(t, pullZone, "Pull zone should be nil")
	assert.Contains(t, err.Error(), "pullzone.invalid_origin")
}

func TestPullZoneService_Update_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "pullzone.not_found",
		"Field": "PullZoneId",
		"Message": "The requested Pull Zone was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create the pull zone options
	updateOptions := &resources.PullZone{
		OriginUrl: "https://updated-example.com",
	}

	// Call the Update method with an invalid ID
	pullZone, err := client.PullZone.Update(context.Background(), 99999, updateOptions)
	assert.Error(t, err, "Update should return an error")
	assert.Nil(t, pullZone, "Pull zone should be nil")
	assert.Contains(t, err.Error(), "pullzone.not_found")
}

func TestPullZoneService_Update_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Id": 12345,
		"Name": "test-zone-1",
		"OriginUrl": "https://updated-example.com",
		"Enabled": true,
		"Hostnames": [],
		"AllowedReferrers": ["example.com"],
		"BlockedReferrers": ["badsite.com"],
		"EnableGeoZoneUS": true,
		"EnableGeoZoneEU": false,
		"Type": 0
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone/12345")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create the pull zone options
	updateOptions := &resources.PullZone{
		OriginUrl:       "https://updated-example.com",
		EnableGeoZoneEU: false,
	}

	// Call the Update method
	pullZone, err := client.PullZone.Update(context.Background(), 12345, updateOptions)
	assert.NoError(t, err, "Update should not return an error")
	assert.NotNil(t, pullZone, "Pull zone should not be nil")

	// Verify pull zone
	assert.Equal(t, int64(12345), pullZone.Id)
	assert.Equal(t, "test-zone-1", pullZone.Name)
	assert.Equal(t, "https://updated-example.com", pullZone.OriginUrl)
	assert.True(t, pullZone.Enabled)
	assert.Empty(t, pullZone.Hostnames)
	assert.Equal(t, []string{"example.com"}, pullZone.AllowedReferrers)
	assert.Equal(t, []string{"badsite.com"}, pullZone.BlockedReferrers)
	assert.True(t, pullZone.EnableGeoZoneUS)
	assert.False(t, pullZone.EnableGeoZoneEU)
	assert.Equal(t, 0, pullZone.Type)
}

func TestPullZoneService_Delete_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "pullzone.not_found",
		"Field": "PullZoneId",
		"Message": "The requested Pull Zone was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Delete method with an invalid ID
	err := client.PullZone.Delete(context.Background(), 99999)
	assert.Error(t, err, "Delete should return an error")
	assert.Contains(t, err.Error(), "pullzone.not_found")
}

func TestPullZoneService_Delete_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/pullzone/12345")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Delete method
	err := client.PullZone.Delete(context.Background(), 12345)
	assert.NoError(t, err, "Delete should not return an error")
}

func TestPullZoneService_PurgeCache_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "pullzone.not_found",
		"Field": "PullZoneId",
		"Message": "The requested Pull Zone was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the PurgeCache method with an invalid ID
	err := client.PullZone.PurgeCache(context.Background(), 99999, nil)
	assert.Error(t, err, "PurgeCache should return an error")
	assert.Contains(t, err.Error(), "pullzone.not_found")
}

func TestPullZoneService_PurgeCache_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone/12345/purgeCache")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the PurgeCache method
	options := &resources.PurgeCacheOptions{
		CacheTag: "tag1",
	}
	err := client.PullZone.PurgeCache(context.Background(), 12345, options)
	assert.NoError(t, err, "PurgeCache should not return an error")
}

func TestPullZoneService_AddHostname_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusBadRequest, `{
		"ErrorKey": "hostname.invalid",
		"Field": "Hostname",
		"Message": "The provided hostname is invalid"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddHostname method with an invalid hostname
	options := resources.AddHostnameOptions{
		Hostname: "invalid hostname with spaces",
	}
	err := client.PullZone.AddHostname(context.Background(), 12345, options)
	assert.Error(t, err, "AddHostname should return an error")
	assert.Contains(t, err.Error(), "hostname.invalid")
}

func TestPullZoneService_AddHostname_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone/12345/addHostname")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddHostname method
	options := resources.AddHostnameOptions{
		Hostname: "cdn.example.com",
	}
	err := client.PullZone.AddHostname(context.Background(), 12345, options)
	assert.NoError(t, err, "AddHostname should not return an error")
}

func TestPullZoneService_RemoveHostname_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/pullzone/12345/removeHostname")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the RemoveHostname method
	options := resources.RemoveHostnameOptions{
		Hostname: "cdn.example.com",
	}
	err := client.PullZone.RemoveHostname(context.Background(), 12345, options)
	assert.NoError(t, err, "RemoveHostname should not return an error")
}

func TestPullZoneService_AddCertificate_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusBadRequest, `{
		"ErrorKey": "certificate.invalid",
		"Field": "Certificate",
		"Message": "The provided certificate is invalid"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddCertificate method with invalid certificate data
	options := resources.AddCertificateOptions{
		Hostname:       "cdn.example.com",
		Certificate:    "INVALID_CERT_DATA",
		CertificateKey: "INVALID_KEY_DATA",
	}
	err := client.PullZone.AddCertificate(context.Background(), 12345, options)
	assert.Error(t, err, "AddCertificate should return an error")
	assert.Contains(t, err.Error(), "certificate.invalid")
}

func TestPullZoneService_AddCertificate_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone/12345/addCertificate")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddCertificate method
	options := resources.AddCertificateOptions{
		Hostname:       "cdn.example.com",
		Certificate:    "BASE64_CERT_DATA",
		CertificateKey: "BASE64_KEY_DATA",
	}
	err := client.PullZone.AddCertificate(context.Background(), 12345, options)
	assert.NoError(t, err, "AddCertificate should not return an error")
}

func TestPullZoneService_AddOrUpdateEdgeRule_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusBadRequest, `{
		"ErrorKey": "edgerule.invalid",
		"Field": "ActionType",
		"Message": "The provided action type is invalid"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddOrUpdateEdgeRule method with an invalid action type
	options := resources.AddOrUpdateEdgeRuleOptions{
		ActionType: 999, // Invalid action type
		Triggers: []resources.EdgeRuleTrigger{
			{
				Type:                0,
				PatternMatches:      []string{"example.com/*"},
				PatternMatchingType: 0,
				TriggerMatchingType: 0,
			},
		},
		Enabled: true,
	}
	err := client.PullZone.AddOrUpdateEdgeRule(context.Background(), 12345, options)
	assert.Error(t, err, "AddOrUpdateEdgeRule should return an error")
	assert.Contains(t, err.Error(), "edgerule.invalid")
}

func TestPullZoneService_AddOrUpdateEdgeRule_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusCreated, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone/12345/edgerules/addOrUpdate")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddOrUpdateEdgeRule method
	options := resources.AddOrUpdateEdgeRuleOptions{
		ActionType: 0, // ForceSSL
		Triggers: []resources.EdgeRuleTrigger{
			{
				Type:                0, // Url
				PatternMatches:      []string{"example.com/*"},
				PatternMatchingType: 0, // MatchAny
				TriggerMatchingType: 0, // MatchAny
			},
		},
		Description: "Force SSL for example.com",
		Enabled:     true,
	}
	err := client.PullZone.AddOrUpdateEdgeRule(context.Background(), 12345, options)
	assert.NoError(t, err, "AddOrUpdateEdgeRule should not return an error")
}

func TestPullZoneService_DeleteEdgeRule_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/pullzone/12345/edgerules/abcd1234")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the DeleteEdgeRule method
	err := client.PullZone.DeleteEdgeRule(context.Background(), 12345, "abcd1234")
	assert.NoError(t, err, "DeleteEdgeRule should not return an error")
}

func TestPullZoneService_GetOriginShieldQueueStatistics_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "pullzone.not_found",
		"Field": "PullZoneId",
		"Message": "The requested Pull Zone was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the GetOriginShieldQueueStatistics method with an invalid ID
	stats, err := client.PullZone.GetOriginShieldQueueStatistics(context.Background(), 99999, nil)
	assert.Error(t, err, "GetOriginShieldQueueStatistics should return an error")
	assert.Nil(t, stats, "Statistics should be nil")
	assert.Contains(t, err.Error(), "pullzone.not_found")
}

func TestPullZoneService_GetOriginShieldQueueStatistics_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"ConcurrentRequestsChart": {
			"2023-01-01": 10,
			"2023-01-02": 15
		},
		"QueuedRequestsChart": {
			"2023-01-01": 5,
			"2023-01-02": 8
		}
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/pullzone/12345/originshield/queuestatistics")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		assert.Equal(t, "true", r.URL.Query().Get("hourly"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create statistics options
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	options := &resources.StatisticsOptions{
		DateFrom: &yesterday,
		DateTo:   &now,
		Hourly:   true,
	}

	// Call the GetOriginShieldQueueStatistics method
	stats, err := client.PullZone.GetOriginShieldQueueStatistics(context.Background(), 12345, options)
	assert.NoError(t, err, "GetOriginShieldQueueStatistics should not return an error")
	assert.NotNil(t, stats, "Statistics should not be nil")

	// Verify statistics
	assert.NotNil(t, stats.ConcurrentRequestsChart)
	assert.NotNil(t, stats.QueuedRequestsChart)
	assert.Equal(t, float64(10), stats.ConcurrentRequestsChart["2023-01-01"])
	assert.Equal(t, float64(15), stats.ConcurrentRequestsChart["2023-01-02"])
	assert.Equal(t, float64(5), stats.QueuedRequestsChart["2023-01-01"])
	assert.Equal(t, float64(8), stats.QueuedRequestsChart["2023-01-02"])
}

func TestPullZoneService_CheckAvailability_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusBadRequest, `{
		"ErrorKey": "pullzone.name_required",
		"Field": "Name",
		"Message": "The pull zone name is required"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the CheckAvailability method with an empty name
	options := resources.CheckAvailabilityOptions{
		Name: "",
	}
	response, err := client.PullZone.CheckAvailability(context.Background(), options)
	assert.Error(t, err, "CheckAvailability should return an error")
	assert.Nil(t, response, "Response should be nil")
	assert.Contains(t, err.Error(), "pullzone.name_required")
}

func TestPullZoneService_CheckAvailability_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Available": true
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone/checkavailability")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the CheckAvailability method
	options := resources.CheckAvailabilityOptions{
		Name: "test-zone-1",
	}
	response, err := client.PullZone.CheckAvailability(context.Background(), options)
	assert.NoError(t, err, "CheckAvailability should not return an error")
	assert.NotNil(t, response, "Response should not be nil")
	assert.True(t, response.Available, "Name should be available")
}

func TestPullZoneService_LoadFreeCertificate_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusBadRequest, `{
		"ErrorKey": "certificate.hostname_not_found",
		"Field": "Hostname",
		"Message": "The provided hostname was not found in your account"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the LoadFreeCertificate method with a hostname that doesn't exist
	options := resources.LoadFreeCertificateOptions{
		Hostname: "nonexistent.example.com",
	}
	err := client.PullZone.LoadFreeCertificate(context.Background(), options)
	assert.Error(t, err, "LoadFreeCertificate should return an error")
	assert.Contains(t, err.Error(), "certificate.hostname_not_found")
}

func TestPullZoneService_LoadFreeCertificate_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/pullzone/loadFreeCertificate")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		assert.Equal(t, "cdn.example.com", r.URL.Query().Get("hostname"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the LoadFreeCertificate method
	options := resources.LoadFreeCertificateOptions{
		Hostname: "cdn.example.com",
	}
	err := client.PullZone.LoadFreeCertificate(context.Background(), options)
	assert.NoError(t, err, "LoadFreeCertificate should not return an error")
}
func TestPullZoneService_List_Error(t *testing.T) {
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
	response, err := client.PullZone.List(context.Background(), nil, "", false)
	assert.Error(t, err, "List should return an error")
	assert.Nil(t, response, "Response should be nil")
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestPullZoneService_Get_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Id": 12345,
		"Name": "test-zone-1",
		"OriginUrl": "https://example.com",
		"Enabled": true,
		"Hostnames": [
			{
				"Id": 111,
				"Value": "cdn.example.com",
				"ForceSSL": true,
				"IsSystemHostname": false,
				"HasCertificate": true
			}
		],
		"AllowedReferrers": ["example.com"],
		"BlockedReferrers": ["badsite.com"],
		"EnableGeoZoneUS": true,
		"EnableGeoZoneEU": true,
		"Type": 0,
		"EdgeRules": [
			{
				"Guid": "abcd1234",
				"ActionType": 0,
				"ActionParameter1": "",
				"ActionParameter2": "",
				"Triggers": [
					{
						"Type": 0,
						"PatternMatches": ["example.com/*"],
						"PatternMatchingType": 0,
						"Parameter1": "",
						"TriggerMatchingType": 0
					}
				],
				"Description": "Force SSL for example.com",
				"Enabled": true
			}
		]
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/pullzone/12345")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		assert.Equal(t, "true", r.URL.Query().Get("includeCertificate"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method
	pullZone, err := client.PullZone.Get(context.Background(), 12345, true)
	assert.NoError(t, err, "Get should not return an error")
	assert.NotNil(t, pullZone, "Pull zone should not be nil")

	// Verify pull zone
	assert.Equal(t, int64(12345), pullZone.Id)
	assert.Equal(t, "test-zone-1", pullZone.Name)
	assert.Equal(t, "https://example.com", pullZone.OriginUrl)
	assert.True(t, pullZone.Enabled)
	assert.Len(t, pullZone.Hostnames, 1)
	assert.Equal(t, int64(111), pullZone.Hostnames[0].Id)
	assert.Equal(t, "cdn.example.com", pullZone.Hostnames[0].Value)
	assert.True(t, pullZone.Hostnames[0].ForceSSL)
	assert.True(t, pullZone.Hostnames[0].HasCertificate)
	assert.Equal(t, []string{"example.com"}, pullZone.AllowedReferrers)
	assert.Equal(t, []string{"badsite.com"}, pullZone.BlockedReferrers)
	assert.True(t, pullZone.EnableGeoZoneUS)
	assert.True(t, pullZone.EnableGeoZoneEU)
	assert.Equal(t, 0, pullZone.Type)

	// Verify edge rules
	assert.Len(t, pullZone.EdgeRules, 1)
	assert.Equal(t, "abcd1234", pullZone.EdgeRules[0].Guid)
	assert.Equal(t, 0, pullZone.EdgeRules[0].ActionType)
	assert.Equal(t, "Force SSL for example.com", pullZone.EdgeRules[0].Description)
	assert.True(t, pullZone.EdgeRules[0].Enabled)
	assert.Len(t, pullZone.EdgeRules[0].Triggers, 1)
	assert.Equal(t, 0, pullZone.EdgeRules[0].Triggers[0].Type)
	assert.Equal(t, []string{"example.com/*"}, pullZone.EdgeRules[0].Triggers[0].PatternMatches)
}

func TestPullZoneService_Get_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "pullzone.not_found",
		"Field": "PullZoneId",
		"Message": "The requested Pull Zone was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method with an invalid ID
	pullZone, err := client.PullZone.Get(context.Background(), 99999, false)
	assert.Error(t, err, "Get should return an error")
	assert.Nil(t, pullZone, "Pull zone should be nil")
	assert.Contains(t, err.Error(), "pullzone.not_found")
}

func TestPullZoneService_Add_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusCreated, `{
		"Id": 12345,
		"Name": "test-zone-1",
		"OriginUrl": "https://example.com",
		"Enabled": true,
		"Hostnames": [],
		"AllowedReferrers": ["example.com"],
		"BlockedReferrers": ["badsite.com"],
		"EnableGeoZoneUS": true,
		"EnableGeoZoneEU": true,
		"Type": 0
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/pullzone")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create the pull zone options
	options := resources.AddPullZoneOptions{
		Name:             "test-zone-1",
		OriginUrl:        "https://example.com",
		Type:             0,
		AllowedReferrers: []string{"example.com"},
		BlockedReferrers: []string{"badsite.com"},
		EnableGeoZoneUS:  true,
		EnableGeoZoneEU:  true,
	}

	// Call the Add method
	pullZone, err := client.PullZone.Add(context.Background(), options)
	assert.NoError(t, err, "Add should not return an error")
	assert.NotNil(t, pullZone, "Pull zone should not be nil")

	// Verify pull zone
	assert.Equal(t, int64(12345), pullZone.Id)
	assert.Equal(t, "test-zone-1", pullZone.Name)
	assert.Equal(t, "https://example.com", pullZone.OriginUrl)
	assert.True(t, pullZone.Enabled)
	assert.Empty(t, pullZone.Hostnames)
	assert.Equal(t, []string{"example.com"}, pullZone.AllowedReferrers)
	assert.Equal(t, []string{"badsite.com"}, pullZone.BlockedReferrers)
	assert.True(t, pullZone.EnableGeoZoneUS)
	assert.True(t, pullZone.EnableGeoZoneEU)
	assert.Equal(t, 0, pullZone.Type)
}
