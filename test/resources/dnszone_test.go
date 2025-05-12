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

func TestDNSZoneService_List_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Items": [
			{
				"Id": 123,
				"Domain": "example.com",
				"Records": [
					{
						"Id": 456,
						"Type": 0,
						"Ttl": 3600,
						"Value": "192.0.2.1",
						"Name": "@",
						"Weight": 0,
						"Priority": 0,
						"Port": 0,
						"Flags": 0,
						"Accelerated": false,
						"AcceleratedPullZoneId": 0,
						"MonitorStatus": 0,
						"MonitorType": 0,
						"GeolocationLatitude": 0,
						"GeolocationLongitude": 0,
						"SmartRoutingType": 0,
						"Disabled": false
					}
				],
				"DateModified": "2023-01-01T00:00:00Z",
				"DateCreated": "2023-01-01T00:00:00Z",
				"NameserversDetected": true,
				"CustomNameserversEnabled": false,
				"NameserversNextCheck": "2023-01-02T00:00:00Z",
				"DnsSecEnabled": false,
				"LoggingEnabled": false,
				"LoggingIPAnonymizationEnabled": false,
				"LogAnonymizationType": 0
			},
			{
				"Id": 789,
				"Domain": "example.org",
				"Records": [],
				"DateModified": "2023-01-01T00:00:00Z",
				"DateCreated": "2023-01-01T00:00:00Z",
				"NameserversDetected": true,
				"CustomNameserversEnabled": false,
				"NameserversNextCheck": "2023-01-02T00:00:00Z",
				"DnsSecEnabled": false,
				"LoggingEnabled": false,
				"LoggingIPAnonymizationEnabled": false,
				"LogAnonymizationType": 0
			}
		],
		"CurrentPage": 1,
		"TotalItems": 2,
		"HasMoreItems": false
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/dnszone")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")

		// Verify pagination parameters
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("perPage"))

		// Verify search parameter
		assert.Equal(t, "example", r.URL.Query().Get("search"))
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Create pagination parameters
	pagination := common.NewPagination().WithPage(2).WithPerPage(10)

	// Call the List method with search
	response, err := client.DNSZone.List(context.Background(), pagination, "example")
	assert.NoError(t, err, "List should not return an error")
	assert.NotNil(t, response, "Response should not be nil")
	assert.Len(t, response.Items, 2, "Should return 2 DNS zones")

	// Check first DNS zone
	assert.Equal(t, int64(123), response.Items[0].Id)
	assert.Equal(t, "example.com", response.Items[0].Domain)
	assert.Len(t, response.Items[0].Records, 1)

	// Check DNS record in first zone
	record := response.Items[0].Records[0]
	assert.Equal(t, int64(456), record.Id)
	assert.Equal(t, resources.DNSRecordTypeA, record.Type)
	assert.Equal(t, int32(3600), record.Ttl)
	assert.Equal(t, "192.0.2.1", record.Value)
	assert.Equal(t, "@", record.Name)

	// Check second DNS zone
	assert.Equal(t, int64(789), response.Items[1].Id)
	assert.Equal(t, "example.org", response.Items[1].Domain)
	assert.Empty(t, response.Items[1].Records)

	// Check pagination info
	assert.Equal(t, 1, response.CurrentPage)
	assert.Equal(t, 2, response.TotalItems)
	assert.False(t, response.HasMoreItems)
}

func TestDNSZoneService_Get_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Id": 123,
		"Domain": "example.com",
		"Records": [
			{
				"Id": 456,
				"Type": 0,
				"Ttl": 3600,
				"Value": "192.0.2.1",
				"Name": "@",
				"Weight": 0,
				"Priority": 0,
				"Port": 0,
				"Flags": 0,
				"Accelerated": false,
				"AcceleratedPullZoneId": 0,
				"MonitorStatus": 0,
				"MonitorType": 0,
				"GeolocationLatitude": 0,
				"GeolocationLongitude": 0,
				"SmartRoutingType": 0,
				"Disabled": false
			}
		],
		"DateModified": "2023-01-01T00:00:00Z",
		"DateCreated": "2023-01-01T00:00:00Z",
		"NameserversDetected": true,
		"CustomNameserversEnabled": false,
		"NameserversNextCheck": "2023-01-02T00:00:00Z",
		"DnsSecEnabled": false,
		"LoggingEnabled": false,
		"LoggingIPAnonymizationEnabled": false,
		"LogAnonymizationType": 0
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/dnszone/123")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method
	dnsZone, err := client.DNSZone.Get(context.Background(), 123)
	assert.NoError(t, err, "Get should not return an error")
	assert.NotNil(t, dnsZone, "DNS zone should not be nil")

	// Check DNS zone details
	assert.Equal(t, int64(123), dnsZone.Id)
	assert.Equal(t, "example.com", dnsZone.Domain)
	assert.Len(t, dnsZone.Records, 1)
	assert.True(t, dnsZone.NameserversDetected)
	assert.False(t, dnsZone.CustomNameserversEnabled)
	assert.False(t, dnsZone.DnsSecEnabled)

	// Check record details
	record := dnsZone.Records[0]
	assert.Equal(t, int64(456), record.Id)
	assert.Equal(t, resources.DNSRecordTypeA, record.Type)
	assert.Equal(t, int32(3600), record.Ttl)
	assert.Equal(t, "192.0.2.1", record.Value)
	assert.Equal(t, "@", record.Name)
}

func TestDNSZoneService_Add_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusCreated, `{
		"Id": 123,
		"Domain": "example.com",
		"Records": [],
		"DateModified": "2023-01-01T00:00:00Z",
		"DateCreated": "2023-01-01T00:00:00Z",
		"NameserversDetected": false,
		"CustomNameserversEnabled": false,
		"NameserversNextCheck": "2023-01-02T00:00:00Z",
		"DnsSecEnabled": false,
		"LoggingEnabled": false,
		"LoggingIPAnonymizationEnabled": false,
		"LogAnonymizationType": 0
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/dnszone")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Add method
	options := resources.AddDNSZoneOptions{
		Domain: "example.com",
	}
	dnsZone, err := client.DNSZone.Add(context.Background(), options)
	assert.NoError(t, err, "Add should not return an error")
	assert.NotNil(t, dnsZone, "DNS zone should not be nil")

	// Check DNS zone details
	assert.Equal(t, int64(123), dnsZone.Id)
	assert.Equal(t, "example.com", dnsZone.Domain)
	assert.Empty(t, dnsZone.Records)
	assert.False(t, dnsZone.NameserversDetected)
	assert.False(t, dnsZone.CustomNameserversEnabled)
	assert.False(t, dnsZone.DnsSecEnabled)
}

func TestDNSZoneService_Update_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Id": 123,
		"Domain": "example.com",
		"Records": [],
		"DateModified": "2023-01-01T00:00:00Z",
		"DateCreated": "2023-01-01T00:00:00Z",
		"NameserversDetected": true,
		"CustomNameserversEnabled": true,
		"Nameserver1": "ns1.example.com",
		"Nameserver2": "ns2.example.com",
		"SoaEmail": "admin@example.com",
		"NameserversNextCheck": "2023-01-02T00:00:00Z",
		"DnsSecEnabled": false,
		"LoggingEnabled": true,
		"LoggingIPAnonymizationEnabled": true,
		"LogAnonymizationType": 1
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/dnszone/123")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Update method
	options := resources.UpdateDNSZoneOptions{
		CustomNameserversEnabled:      true,
		Nameserver1:                   "ns1.example.com",
		Nameserver2:                   "ns2.example.com",
		SoaEmail:                      "admin@example.com",
		LoggingEnabled:                true,
		LogAnonymizationType:          resources.LogAnonymizationTypeDrop,
		LoggingIPAnonymizationEnabled: true,
	}
	dnsZone, err := client.DNSZone.Update(context.Background(), 123, options)
	assert.NoError(t, err, "Update should not return an error")
	assert.NotNil(t, dnsZone, "DNS zone should not be nil")

	// Check DNS zone details
	assert.Equal(t, int64(123), dnsZone.Id)
	assert.Equal(t, "example.com", dnsZone.Domain)
	assert.True(t, dnsZone.CustomNameserversEnabled)
	assert.Equal(t, "ns1.example.com", dnsZone.Nameserver1)
	assert.Equal(t, "ns2.example.com", dnsZone.Nameserver2)
	assert.Equal(t, "admin@example.com", dnsZone.SoaEmail)
	assert.True(t, dnsZone.LoggingEnabled)
	assert.True(t, dnsZone.LoggingIPAnonymizationEnabled)
	assert.Equal(t, resources.LogAnonymizationTypeDrop, dnsZone.LogAnonymizationType)
}

func TestDNSZoneService_Delete_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/dnszone/123")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Delete method
	err := client.DNSZone.Delete(context.Background(), 123)
	assert.NoError(t, err, "Delete should not return an error")
}

func TestDNSZoneService_EnableDNSSec_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Enabled": true,
		"DsRecord": "example.com. 3600 IN DS 12345 8 2 ABCDEF123456789ABCDEF123456789ABCDEF123456789ABCDEF123456789",
		"Digest": "ABCDEF123456789ABCDEF123456789ABCDEF123456789ABCDEF123456789",
		"DigestType": "SHA-256",
		"Algorithm": 8,
		"PublicKey": "ABCDEF123456789ABCDEF123456789ABCDEF123456789ABCDEF123456789",
		"KeyTag": 12345,
		"Flags": 256
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/dnszone/123/dnssec")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the EnableDNSSec method
	dnsSecInfo, err := client.DNSZone.EnableDNSSec(context.Background(), 123)
	assert.NoError(t, err, "EnableDNSSec should not return an error")
	assert.NotNil(t, dnsSecInfo, "DNSSEC info should not be nil")

	// Check DNSSEC info
	assert.True(t, dnsSecInfo.Enabled)
	assert.Equal(t, "example.com. 3600 IN DS 12345 8 2 ABCDEF123456789ABCDEF123456789ABCDEF123456789ABCDEF123456789", dnsSecInfo.DsRecord)
	assert.Equal(t, "ABCDEF123456789ABCDEF123456789ABCDEF123456789ABCDEF123456789", dnsSecInfo.Digest)
	assert.Equal(t, "SHA-256", dnsSecInfo.DigestType)
	assert.Equal(t, int32(8), dnsSecInfo.Algorithm)
	assert.Equal(t, "ABCDEF123456789ABCDEF123456789ABCDEF123456789ABCDEF123456789", dnsSecInfo.PublicKey)
	assert.Equal(t, int32(12345), dnsSecInfo.KeyTag)
	assert.Equal(t, int32(256), dnsSecInfo.Flags)
}

func TestDNSZoneService_AddRecord_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusCreated, `{
		"Id": 456,
		"Type": 0,
		"Ttl": 3600,
		"Value": "192.0.2.1",
		"Name": "www",
		"Weight": 0,
		"Priority": 0,
		"Port": 0,
		"Flags": 0,
		"Accelerated": false,
		"AcceleratedPullZoneId": 0,
		"MonitorStatus": 0,
		"MonitorType": 0,
		"GeolocationLatitude": 0,
		"GeolocationLongitude": 0,
		"SmartRoutingType": 0,
		"Disabled": false
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPut)
		test.AssertRequestPath(t, r, "/dnszone/123/records")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the AddRecord method
	options := resources.AddDNSRecordOptions{
		Type:  resources.DNSRecordTypeA,
		Ttl:   3600,
		Value: "192.0.2.1",
		Name:  "www",
	}
	record, err := client.DNSZone.AddRecord(context.Background(), 123, options)
	assert.NoError(t, err, "AddRecord should not return an error")
	assert.NotNil(t, record, "DNS record should not be nil")

	// Check DNS record details
	assert.Equal(t, int64(456), record.Id)
	assert.Equal(t, resources.DNSRecordTypeA, record.Type)
	assert.Equal(t, int32(3600), record.Ttl)
	assert.Equal(t, "192.0.2.1", record.Value)
	assert.Equal(t, "www", record.Name)
}

func TestDNSZoneService_UpdateRecord_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/dnszone/123/records/456")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the UpdateRecord method
	options := resources.UpdateDNSRecordOptions{
		Id:    456,
		Type:  resources.DNSRecordTypeA,
		Ttl:   7200,
		Value: "192.0.2.2",
		Name:  "www",
	}
	err := client.DNSZone.UpdateRecord(context.Background(), 123, 456, options)
	assert.NoError(t, err, "UpdateRecord should not return an error")
}

func TestDNSZoneService_DeleteRecord_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusNoContent, ``, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/dnszone/123/records/456")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the DeleteRecord method
	err := client.DNSZone.DeleteRecord(context.Background(), 123, 456)
	assert.NoError(t, err, "DeleteRecord should not return an error")
}

func TestDNSZoneService_CheckAvailability_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Available": true,
		"Message": "Domain is available"
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/dnszone/checkavailability")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		test.AssertRequestHasHeader(t, r, "Content-Type", "application/json")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the CheckAvailability method
	options := resources.CheckZoneAvailabilityOptions{
		Name: "example.com",
	}
	result, err := client.DNSZone.CheckAvailability(context.Background(), options)
	assert.NoError(t, err, "CheckAvailability should not return an error")
	assert.NotNil(t, result, "Availability result should not be nil")
	assert.True(t, result.Available)
	assert.Equal(t, "Domain is available", result.Message)
}

func TestDNSZoneService_Export_Success(t *testing.T) {
	// Sample zone file content
	zoneFileContent := `$ORIGIN example.com.
$TTL 3600
@       IN      SOA     ns1.example.com. admin.example.com. (
                        2023010101      ; serial
                        3600            ; refresh
                        1800            ; retry
                        604800          ; expire
                        86400 )         ; minimum
@       IN      NS      ns1.example.com.
@       IN      NS      ns2.example.com.
@       IN      A       192.0.2.1
www     IN      A       192.0.2.1
`

	// Create a mock server
	server := test.MockServer(t, http.StatusOK, zoneFileContent, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/dnszone/123/export")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Export method
	data, err := client.DNSZone.Export(context.Background(), 123)
	assert.NoError(t, err, "Export should not return an error")
	assert.NotNil(t, data, "Exported data should not be nil")
	assert.Equal(t, zoneFileContent, string(data))
}

func TestDNSZoneService_Error_Handling(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "dnszone.not_found",
		"Field": "ZoneId",
		"Message": "The requested DNS zone was not found"
	}`, nil)
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Test Get method error handling
	dnsZone, err := client.DNSZone.Get(context.Background(), 999)
	assert.Error(t, err, "Get should return an error for non-existent zone")
	assert.Nil(t, dnsZone, "DNS zone should be nil")
	assert.Contains(t, err.Error(), "dnszone.not_found")
	assert.Contains(t, err.Error(), "The requested DNS zone was not found")
}

func TestDNSZoneService_ListAll(t *testing.T) {
	// Create a mock server for first page
	firstPageServer := test.MockServer(t, http.StatusOK, `{
		"Items": [
			{
				"Id": 123,
				"Domain": "example1.com",
				"Records": [],
				"DateModified": "2023-01-01T00:00:00Z",
				"DateCreated": "2023-01-01T00:00:00Z",
				"NameserversDetected": true,
				"CustomNameserversEnabled": false,
				"NameserversNextCheck": "2023-01-02T00:00:00Z",
				"DnsSecEnabled": false,
				"LoggingEnabled": false,
				"LoggingIPAnonymizationEnabled": false,
				"LogAnonymizationType": 0
			},
			{
				"Id": 456,
				"Domain": "example2.com",
				"Records": [],
				"DateModified": "2023-01-01T00:00:00Z",
				"DateCreated": "2023-01-01T00:00:00Z",
				"NameserversDetected": true,
				"CustomNameserversEnabled": false,
				"NameserversNextCheck": "2023-01-02T00:00:00Z",
				"DnsSecEnabled": false,
				"LoggingEnabled": false,
				"LoggingIPAnonymizationEnabled": false,
				"LogAnonymizationType": 0
			}
		],
		"CurrentPage": 1,
		"TotalItems": 3,
		"HasMoreItems": true
	}`, func(r *http.Request) {
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "2", r.URL.Query().Get("perPage"))
		assert.Equal(t, "test", r.URL.Query().Get("search"))
	})
	defer firstPageServer.Close()

	// Create a second server for the second page
	secondPageServer := test.MockServer(t, http.StatusOK, `{
		"Items": [
			{
				"Id": 789,
				"Domain": "example3.com",
				"Records": [],
				"DateModified": "2023-01-01T00:00:00Z",
				"DateCreated": "2023-01-01T00:00:00Z",
				"NameserversDetected": true,
				"CustomNameserversEnabled": false,
				"NameserversNextCheck": "2023-01-02T00:00:00Z",
				"DnsSecEnabled": false,
				"LoggingEnabled": false,
				"LoggingIPAnonymizationEnabled": false,
				"LogAnonymizationType": 0
			}
		],
		"CurrentPage": 2,
		"TotalItems": 3,
		"HasMoreItems": false
	}`, func(r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "2", r.URL.Query().Get("perPage"))
		assert.Equal(t, "test", r.URL.Query().Get("search"))
	})
	defer secondPageServer.Close()

	// Create a client for the first page
	firstPageClient := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(firstPageServer.URL))

	// Test ListAll with the first page
	zones, err := firstPageClient.DNSZone.ListAll(context.Background(), 2, "test")
	assert.NoError(t, err, "ListAll should not return an error")
	assert.NotNil(t, zones, "Zones should not be nil")

	// We can't easily test pagination across multiple mocked servers in this test framework,
	// so we'll just verify the first page results
	assert.Len(t, zones, 2, "Should return 2 zones from first page")
	assert.Equal(t, "example1.com", zones[0].Domain)
	assert.Equal(t, "example2.com", zones[1].Domain)

	// Create a client for the second page
	secondPageClient := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(secondPageServer.URL))

	// Test the second page
	secondPageResponse, err := secondPageClient.DNSZone.List(context.Background(), common.NewPagination().WithPage(2).WithPerPage(2), "test")
	assert.NoError(t, err, "List should not return an error")
	assert.Len(t, secondPageResponse.Items, 1, "Should return 1 zone from second page")
	assert.Equal(t, "example3.com", secondPageResponse.Items[0].Domain)
}

func TestDNSZoneService_DisableDNSSec_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Enabled": false,
		"DsRecord": null,
		"Digest": null,
		"DigestType": null,
		"Algorithm": 0,
		"PublicKey": null,
		"KeyTag": 0,
		"Flags": 0
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodDelete)
		test.AssertRequestPath(t, r, "/dnszone/123/dnssec")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the DisableDNSSec method
	dnsSecInfo, err := client.DNSZone.DisableDNSSec(context.Background(), 123)
	assert.NoError(t, err, "DisableDNSSec should not return an error")
	assert.NotNil(t, dnsSecInfo, "DNSSEC info should not be nil")

	// Check DNSSEC info
	assert.False(t, dnsSecInfo.Enabled)
	assert.Empty(t, dnsSecInfo.DsRecord)
	assert.Empty(t, dnsSecInfo.Digest)
	assert.Empty(t, dnsSecInfo.DigestType)
	assert.Zero(t, dnsSecInfo.Algorithm)
	assert.Empty(t, dnsSecInfo.PublicKey)
	assert.Zero(t, dnsSecInfo.KeyTag)
	assert.Zero(t, dnsSecInfo.Flags)
}

func TestDNSZoneService_ImportRecords_Success(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"RecordsSuccessful": 5,
		"RecordsFailed": 1,
		"RecordsSkipped": 2
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodPost)
		test.AssertRequestPath(t, r, "/dnszone/123/import")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Sample zone file data to import
	zoneData := []byte(`$ORIGIN example.com.
$TTL 3600
@       IN      A       192.0.2.1
www     IN      A       192.0.2.1
mail    IN      A       192.0.2.2
`)

	// Call the ImportRecords method
	result, err := client.DNSZone.ImportRecords(context.Background(), 123, zoneData)
	assert.NoError(t, err, "ImportRecords should not return an error")
	assert.NotNil(t, result, "Import result should not be nil")

	// Check import result
	assert.Equal(t, int32(5), result.RecordsSuccessful)
	assert.Equal(t, int32(1), result.RecordsFailed)
	assert.Equal(t, int32(2), result.RecordsSkipped)
}
