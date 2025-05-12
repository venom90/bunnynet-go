// Package resources provides API resource implementations for the Bunny.net API client
package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/venom90/bunnynet-go-client/common"
	"github.com/venom90/bunnynet-go-client/internal"
)

// DNSRecordType represents the type of a DNS record
type DNSRecordType int

const (
	// DNSRecordTypeA represents an A record
	DNSRecordTypeA DNSRecordType = 0
	// DNSRecordTypeAAAA represents an AAAA record
	DNSRecordTypeAAAA DNSRecordType = 1
	// DNSRecordTypeCNAME represents a CNAME record
	DNSRecordTypeCNAME DNSRecordType = 2
	// DNSRecordTypeTXT represents a TXT record
	DNSRecordTypeTXT DNSRecordType = 3
	// DNSRecordTypeMX represents an MX record
	DNSRecordTypeMX DNSRecordType = 4
	// DNSRecordTypeRedirect represents a Redirect record
	DNSRecordTypeRedirect DNSRecordType = 5
	// DNSRecordTypeFlatten represents a Flatten record
	DNSRecordTypeFlatten DNSRecordType = 6
	// DNSRecordTypePullZone represents a PullZone record
	DNSRecordTypePullZone DNSRecordType = 7
	// DNSRecordTypeSRV represents an SRV record
	DNSRecordTypeSRV DNSRecordType = 8
	// DNSRecordTypeCAA represents a CAA record
	DNSRecordTypeCAA DNSRecordType = 9
	// DNSRecordTypePTR represents a PTR record
	DNSRecordTypePTR DNSRecordType = 10
	// DNSRecordTypeScript represents a Script record
	DNSRecordTypeScript DNSRecordType = 11
	// DNSRecordTypeNS represents a NS record
	DNSRecordTypeNS DNSRecordType = 12
)

// MonitorStatus represents the status of a DNS record monitor
type MonitorStatus int

const (
	// MonitorStatusUnknown represents an unknown monitor status
	MonitorStatusUnknown MonitorStatus = 0
	// MonitorStatusOnline represents an online monitor status
	MonitorStatusOnline MonitorStatus = 1
	// MonitorStatusOffline represents an offline monitor status
	MonitorStatusOffline MonitorStatus = 2
)

// MonitorType represents the type of a DNS record monitor
type MonitorType int

const (
	// MonitorTypeNone represents no monitoring
	MonitorTypeNone MonitorType = 0
	// MonitorTypePing represents ping monitoring
	MonitorTypePing MonitorType = 1
	// MonitorTypeHTTP represents HTTP monitoring
	MonitorTypeHTTP MonitorType = 2
	// MonitorTypeMonitor represents other monitoring
	MonitorTypeMonitor MonitorType = 3
)

// SmartRoutingType represents the type of smart routing
type SmartRoutingType int

const (
	// SmartRoutingTypeNone represents no smart routing
	SmartRoutingTypeNone SmartRoutingType = 0
	// SmartRoutingTypeLatency represents latency-based smart routing
	SmartRoutingTypeLatency SmartRoutingType = 1
	// SmartRoutingTypeGeolocation represents geolocation-based smart routing
	SmartRoutingTypeGeolocation SmartRoutingType = 2
)

// LogAnonymizationType represents the type of log anonymization
type LogAnonymizationType int

const (
	// LogAnonymizationTypeOneDigit represents one-digit anonymization
	LogAnonymizationTypeOneDigit LogAnonymizationType = 0
	// LogAnonymizationTypeDrop represents drop anonymization
	LogAnonymizationTypeDrop LogAnonymizationType = 1
)

// GeoLocationInfo represents geolocation information for a DNS record
type GeoLocationInfo struct {
	// Country is the name of the country
	Country string `json:"Country"`

	// City is the name of the city
	City string `json:"City"`

	// Latitude is the latitude of the location
	Latitude float64 `json:"Latitude"`

	// Longitude is the longitude of the location
	Longitude float64 `json:"Longitude"`
}

// IPGeoLocationInfo represents IP geolocation information for a DNS record
type IPGeoLocationInfo struct {
	// CountryCode is the ISO country code of the location
	CountryCode string `json:"CountryCode"`

	// Country is the name of the country of the location
	Country string `json:"Country"`

	// ASN is the ASN of the IP organization
	ASN int64 `json:"ASN"`

	// OrganizationName is the name of the organization that owns the IP
	OrganizationName string `json:"OrganizationName"`

	// City is the name of the city of the location
	City string `json:"City"`
}

// EnvironmentalVariable represents an environmental variable for a DNS record
type EnvironmentalVariable struct {
	// Name is the name of the environmental variable
	Name string `json:"Name"`

	// Value is the value of the environmental variable
	Value string `json:"Value"`
}

// DNSRecord represents a DNS record in the Bunny.net API
type DNSRecord struct {
	// Id is the unique identifier of the DNS record
	Id int64 `json:"Id"`

	// Type is the type of the DNS record
	Type DNSRecordType `json:"Type"`

	// Ttl is the time to live of the DNS record
	Ttl int32 `json:"Ttl"`

	// Value is the value of the DNS record
	Value string `json:"Value"`

	// Name is the name of the DNS record
	Name string `json:"Name"`

	// Weight is the weight of the DNS record
	Weight int32 `json:"Weight"`

	// Priority is the priority of the DNS record
	Priority int32 `json:"Priority"`

	// Port is the port of the DNS record
	Port int32 `json:"Port"`

	// Flags is the flags of the DNS record
	Flags int `json:"Flags"`

	// Tag is the tag of the DNS record
	Tag string `json:"Tag"`

	// Accelerated indicates whether the DNS record is accelerated
	Accelerated bool `json:"Accelerated"`

	// AcceleratedPullZoneId is the ID of the accelerated pull zone
	AcceleratedPullZoneId int64 `json:"AcceleratedPullZoneId"`

	// LinkName is the link name of the DNS record
	LinkName string `json:"LinkName"`

	// IPGeoLocationInfo is the IP geolocation information of the DNS record
	IPGeoLocationInfo IPGeoLocationInfo `json:"IPGeoLocationInfo"`

	// GeolocationInfo is the geolocation information of the DNS record
	GeolocationInfo GeoLocationInfo `json:"GeolocationInfo"`

	// MonitorStatus is the monitor status of the DNS record
	MonitorStatus MonitorStatus `json:"MonitorStatus"`

	// MonitorType is the monitor type of the DNS record
	MonitorType MonitorType `json:"MonitorType"`

	// GeolocationLatitude is the geolocation latitude of the DNS record
	GeolocationLatitude float64 `json:"GeolocationLatitude"`

	// GeolocationLongitude is the geolocation longitude of the DNS record
	GeolocationLongitude float64 `json:"GeolocationLongitude"`

	// EnvironmentalVariables is the list of environmental variables of the DNS record
	EnvironmentalVariables []EnvironmentalVariable `json:"EnviromentalVariables"`

	// LatencyZone is the latency zone of the DNS record
	LatencyZone string `json:"LatencyZone"`

	// SmartRoutingType is the smart routing type of the DNS record
	SmartRoutingType SmartRoutingType `json:"SmartRoutingType"`

	// Disabled indicates whether the DNS record is disabled
	Disabled bool `json:"Disabled"`

	// Comment is the comment of the DNS record
	Comment string `json:"Comment"`
}

// AddDNSRecordOptions represents the options for adding a DNS record
type AddDNSRecordOptions struct {
	// Type is the type of the DNS record
	Type DNSRecordType `json:"Type"`

	// Ttl is the time to live of the DNS record
	Ttl int32 `json:"Ttl,omitempty"`

	// Value is the value of the DNS record
	Value string `json:"Value,omitempty"`

	// Name is the name of the DNS record
	Name string `json:"Name,omitempty"`

	// Weight is the weight of the DNS record
	Weight int32 `json:"Weight,omitempty"`

	// Priority is the priority of the DNS record
	Priority int32 `json:"Priority,omitempty"`

	// Flags is the flags of the DNS record
	Flags int `json:"Flags,omitempty"`

	// Tag is the tag of the DNS record
	Tag string `json:"Tag,omitempty"`

	// Port is the port of the DNS record
	Port int32 `json:"Port,omitempty"`

	// PullZoneId is the ID of the pull zone
	PullZoneId int64 `json:"PullZoneId,omitempty"`

	// ScriptId is the ID of the script
	ScriptId int64 `json:"ScriptId,omitempty"`

	// Accelerated indicates whether the DNS record should be accelerated
	Accelerated bool `json:"Accelerated,omitempty"`

	// MonitorType is the monitor type of the DNS record
	MonitorType MonitorType `json:"MonitorType,omitempty"`

	// GeolocationLatitude is the geolocation latitude of the DNS record
	GeolocationLatitude float64 `json:"GeolocationLatitude,omitempty"`

	// GeolocationLongitude is the geolocation longitude of the DNS record
	GeolocationLongitude float64 `json:"GeolocationLongitude,omitempty"`

	// LatencyZone is the latency zone of the DNS record
	LatencyZone string `json:"LatencyZone,omitempty"`

	// SmartRoutingType is the smart routing type of the DNS record
	SmartRoutingType SmartRoutingType `json:"SmartRoutingType,omitempty"`

	// Disabled indicates whether the DNS record should be disabled
	Disabled bool `json:"Disabled,omitempty"`

	// EnvironmentalVariables is the list of environmental variables of the DNS record
	EnvironmentalVariables []EnvironmentalVariable `json:"EnviromentalVariables,omitempty"`

	// Comment is the comment of the DNS record
	Comment string `json:"Comment,omitempty"`
}

// UpdateDNSRecordOptions represents the options for updating a DNS record
type UpdateDNSRecordOptions struct {
	// Id is the ID of the DNS record
	Id int64 `json:"Id"`

	// Type is the type of the DNS record
	Type DNSRecordType `json:"Type"`

	// Ttl is the time to live of the DNS record
	Ttl int32 `json:"Ttl,omitempty"`

	// Value is the value of the DNS record
	Value string `json:"Value,omitempty"`

	// Name is the name of the DNS record
	Name string `json:"Name,omitempty"`

	// Weight is the weight of the DNS record
	Weight int32 `json:"Weight,omitempty"`

	// Priority is the priority of the DNS record
	Priority int32 `json:"Priority,omitempty"`

	// Flags is the flags of the DNS record
	Flags int `json:"Flags,omitempty"`

	// Tag is the tag of the DNS record
	Tag string `json:"Tag,omitempty"`

	// Port is the port of the DNS record
	Port int32 `json:"Port,omitempty"`

	// PullZoneId is the ID of the pull zone
	PullZoneId int64 `json:"PullZoneId,omitempty"`

	// ScriptId is the ID of the script
	ScriptId int64 `json:"ScriptId,omitempty"`

	// Accelerated indicates whether the DNS record should be accelerated
	Accelerated bool `json:"Accelerated,omitempty"`

	// MonitorType is the monitor type of the DNS record
	MonitorType MonitorType `json:"MonitorType,omitempty"`

	// GeolocationLatitude is the geolocation latitude of the DNS record
	GeolocationLatitude float64 `json:"GeolocationLatitude,omitempty"`

	// GeolocationLongitude is the geolocation longitude of the DNS record
	GeolocationLongitude float64 `json:"GeolocationLongitude,omitempty"`

	// LatencyZone is the latency zone of the DNS record
	LatencyZone string `json:"LatencyZone,omitempty"`

	// SmartRoutingType is the smart routing type of the DNS record
	SmartRoutingType SmartRoutingType `json:"SmartRoutingType,omitempty"`

	// Disabled indicates whether the DNS record should be disabled
	Disabled bool `json:"Disabled,omitempty"`

	// EnvironmentalVariables is the list of environmental variables of the DNS record
	EnvironmentalVariables []EnvironmentalVariable `json:"EnviromentalVariables,omitempty"`

	// Comment is the comment of the DNS record
	Comment string `json:"Comment,omitempty"`
}

// DNSZone represents a DNS zone in the Bunny.net API
type DNSZone struct {
	// Id is the unique identifier of the DNS zone
	Id int64 `json:"Id"`

	// Domain is the domain of the DNS zone
	Domain string `json:"Domain"`

	// Records is the list of DNS records in the zone
	Records []DNSRecord `json:"Records"`

	// DateModified is the date and time when the DNS zone was last modified
	DateModified time.Time `json:"DateModified"`

	// DateCreated is the date and time when the DNS zone was created
	DateCreated time.Time `json:"DateCreated"`

	// NameserversDetected indicates whether nameservers have been detected
	NameserversDetected bool `json:"NameserversDetected"`

	// CustomNameserversEnabled indicates whether custom nameservers are enabled
	CustomNameserversEnabled bool `json:"CustomNameserversEnabled"`

	// Nameserver1 is the first custom nameserver
	Nameserver1 string `json:"Nameserver1"`

	// Nameserver2 is the second custom nameserver
	Nameserver2 string `json:"Nameserver2"`

	// SoaEmail is the SOA email of the DNS zone
	SoaEmail string `json:"SoaEmail"`

	// NameserversNextCheck is the date and time of the next nameserver check
	NameserversNextCheck time.Time `json:"NameserversNextCheck"`

	// DnsSecEnabled indicates whether DNSSEC is enabled
	DnsSecEnabled bool `json:"DnsSecEnabled"`

	// LoggingEnabled indicates whether logging is enabled
	LoggingEnabled bool `json:"LoggingEnabled"`

	// LoggingIPAnonymizationEnabled indicates whether IP anonymization is enabled for logging
	LoggingIPAnonymizationEnabled bool `json:"LoggingIPAnonymizationEnabled"`

	// LogAnonymizationType is the type of log anonymization
	LogAnonymizationType LogAnonymizationType `json:"LogAnonymizationType"`
}

// AddDNSZoneOptions represents the options for adding a DNS zone
type AddDNSZoneOptions struct {
	// Domain is the domain of the DNS zone to add
	Domain string `json:"Domain"`
}

// UpdateDNSZoneOptions represents the options for updating a DNS zone
type UpdateDNSZoneOptions struct {
	// CustomNameserversEnabled indicates whether custom nameservers should be enabled
	CustomNameserversEnabled bool `json:"CustomNameserversEnabled,omitempty"`

	// Nameserver1 is the first custom nameserver
	Nameserver1 string `json:"Nameserver1,omitempty"`

	// Nameserver2 is the second custom nameserver
	Nameserver2 string `json:"Nameserver2,omitempty"`

	// SoaEmail is the SOA email of the DNS zone
	SoaEmail string `json:"SoaEmail,omitempty"`

	// LoggingEnabled indicates whether logging should be enabled
	LoggingEnabled bool `json:"LoggingEnabled,omitempty"`

	// LogAnonymizationType is the type of log anonymization
	LogAnonymizationType LogAnonymizationType `json:"LogAnonymizationType,omitempty"`

	// LoggingIPAnonymizationEnabled indicates whether IP anonymization should be enabled for logging
	LoggingIPAnonymizationEnabled bool `json:"LoggingIPAnonymizationEnabled,omitempty"`
}

// DNSSecInfo represents DNSSEC information for a DNS zone
type DNSSecInfo struct {
	// Enabled indicates whether DNSSEC is enabled
	Enabled bool `json:"Enabled"`

	// DsRecord is the DS record
	DsRecord string `json:"DsRecord"`

	// Digest is the digest
	Digest string `json:"Digest"`

	// DigestType is the digest type
	DigestType string `json:"DigestType"`

	// Algorithm is the algorithm
	Algorithm int32 `json:"Algorithm"`

	// PublicKey is the public key
	PublicKey string `json:"PublicKey"`

	// KeyTag is the key tag
	KeyTag int32 `json:"KeyTag"`

	// Flags are the flags
	Flags int32 `json:"Flags"`
}

// ImportResult represents the result of importing DNS records
type ImportResult struct {
	// RecordsSuccessful is the number of successfully imported records
	RecordsSuccessful int32 `json:"RecordsSuccessful"`

	// RecordsFailed is the number of failed imported records
	RecordsFailed int32 `json:"RecordsFailed"`

	// RecordsSkipped is the number of skipped imported records
	RecordsSkipped int32 `json:"RecordsSkipped"`
}

// CheckZoneAvailabilityOptions represents the options for checking zone availability
type CheckZoneAvailabilityOptions struct {
	// Name is the name of the zone to check
	Name string `json:"Name"`
}

// ZoneAvailabilityResult represents the result of a zone availability check
type ZoneAvailabilityResult struct {
	// Available indicates whether the zone is available
	Available bool `json:"Available"`

	// Message is a message about the zone availability
	Message string `json:"Message"`
}

// DNSZoneService handles operations on DNS zones
type DNSZoneService struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	userAgent string
}

// NewDNSZoneService creates a new DNSZoneService
func NewDNSZoneService(client *http.Client, baseURL, apiKey, userAgent string) *DNSZoneService {
	return &DNSZoneService{
		client:    client,
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: userAgent,
	}
}

// SetAPIKey updates the API key used for authentication
func (s *DNSZoneService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// List returns a paginated list of DNS zones
func (s *DNSZoneService) List(ctx context.Context, pagination *common.Pagination, search string) (*common.PaginatedResponse[DNSZone], error) {
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, "/dnszone", nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add pagination parameters
	if err := internal.AddQueryParams(req, pagination); err != nil {
		return nil, err
	}

	// Add search parameter if provided
	if search != "" {
		q := req.URL.Query()
		q.Add("search", search)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var paginatedResponse common.PaginatedResponse[DNSZone]
	if err := internal.ParsePaginatedResponse(resp, &paginatedResponse); err != nil {
		return nil, err
	}

	return &paginatedResponse, nil
}

// ListAll returns all DNS zones across all pages
func (s *DNSZoneService) ListAll(ctx context.Context, perPage int, search string) ([]DNSZone, error) {
	if perPage <= 0 {
		perPage = common.DefaultPerPage
	}

	iterator := common.NewPageIterator(
		func(page, itemsPerPage int) (*common.PaginatedResponse[DNSZone], error) {
			pagination := common.NewPagination().WithPage(page).WithPerPage(itemsPerPage)
			return s.List(ctx, pagination, search)
		},
		common.DefaultPage,
		perPage,
	)

	return iterator.AllItems()
}

// Get returns a DNS zone by ID
func (s *DNSZoneService) Get(ctx context.Context, id int64) (*DNSZone, error) {
	path := "/dnszone/" + internal.FormatInt64(id)
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var dnsZone DNSZone
	if err := internal.ParseResponse(resp, &dnsZone); err != nil {
		return nil, err
	}

	return &dnsZone, nil
}

// Add creates a new DNS zone
func (s *DNSZoneService) Add(ctx context.Context, options AddDNSZoneOptions) (*DNSZone, error) {
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, "/dnszone", options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var dnsZone DNSZone
	if err := internal.ParseResponse(resp, &dnsZone); err != nil {
		return nil, err
	}

	return &dnsZone, nil
}

// Update updates a DNS zone
func (s *DNSZoneService) Update(ctx context.Context, id int64, options UpdateDNSZoneOptions) (*DNSZone, error) {
	path := "/dnszone/" + internal.FormatInt64(id)
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, path, options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var dnsZone DNSZone
	if err := internal.ParseResponse(resp, &dnsZone); err != nil {
		return nil, err
	}

	return &dnsZone, nil
}

// Delete deletes a DNS zone
func (s *DNSZoneService) Delete(ctx context.Context, id int64) error {
	path := "/dnszone/" + internal.FormatInt64(id)
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

// EnableDNSSec enables DNSSEC for a DNS zone
func (s *DNSZoneService) EnableDNSSec(ctx context.Context, id int64) (*DNSSecInfo, error) {
	path := "/dnszone/" + internal.FormatInt64(id) + "/dnssec"
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var dnsSecInfo DNSSecInfo
	if err := internal.ParseResponse(resp, &dnsSecInfo); err != nil {
		return nil, err
	}

	return &dnsSecInfo, nil
}

// DisableDNSSec disables DNSSEC for a DNS zone
func (s *DNSZoneService) DisableDNSSec(ctx context.Context, id int64) (*DNSSecInfo, error) {
	path := "/dnszone/" + internal.FormatInt64(id) + "/dnssec"
	req, err := internal.NewRequest(http.MethodDelete, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var dnsSecInfo DNSSecInfo
	if err := internal.ParseResponse(resp, &dnsSecInfo); err != nil {
		return nil, err
	}

	return &dnsSecInfo, nil
}

// Export exports a DNS zone
func (s *DNSZoneService) Export(ctx context.Context, id int64) ([]byte, error) {
	path := "/dnszone/" + internal.FormatInt64(id) + "/export"
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// CheckAvailability checks if a DNS zone is available
func (s *DNSZoneService) CheckAvailability(ctx context.Context, options CheckZoneAvailabilityOptions) (*ZoneAvailabilityResult, error) {
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, "/dnszone/checkavailability", options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var availabilityResult ZoneAvailabilityResult
	if err := internal.ParseResponse(resp, &availabilityResult); err != nil {
		return nil, err
	}

	return &availabilityResult, nil
}

// AddRecord adds a DNS record to a DNS zone
func (s *DNSZoneService) AddRecord(ctx context.Context, zoneId int64, options AddDNSRecordOptions) (*DNSRecord, error) {
	path := fmt.Sprintf("/dnszone/%d/records", zoneId)
	req, err := internal.NewRequest(http.MethodPut, s.baseURL, path, options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var dnsRecord DNSRecord
	if err := internal.ParseResponse(resp, &dnsRecord); err != nil {
		return nil, err
	}

	return &dnsRecord, nil
}

// UpdateRecord updates a DNS record in a DNS zone
func (s *DNSZoneService) UpdateRecord(ctx context.Context, zoneId, recordId int64, options UpdateDNSRecordOptions) error {
	path := fmt.Sprintf("/dnszone/%d/records/%d", zoneId, recordId)
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, path, options, s.apiKey, s.userAgent)
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

// DeleteRecord deletes a DNS record from a DNS zone
func (s *DNSZoneService) DeleteRecord(ctx context.Context, zoneId, recordId int64) error {
	path := fmt.Sprintf("/dnszone/%d/records/%d", zoneId, recordId)
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

// ImportRecords imports DNS records to a DNS zone
func (s *DNSZoneService) ImportRecords(ctx context.Context, zoneId int64, data []byte) (*ImportResult, error) {
	path := fmt.Sprintf("/dnszone/%d/import", zoneId)

	// Create a multipart form request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file part
	part, err := writer.CreateFormFile("file", "import.txt")
	if err != nil {
		return nil, common.NewClientError("failed to create form file", err)
	}

	if _, err := part.Write(data); err != nil {
		return nil, common.NewClientError("failed to write data to form file", err)
	}

	if err := writer.Close(); err != nil {
		return nil, common.NewClientError("failed to close multipart writer", err)
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPost, s.baseURL+path, body)
	if err != nil {
		return nil, common.NewClientError("failed to create request", err)
	}

	req = req.WithContext(ctx)

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("AccessKey", s.apiKey)
	req.Header.Set("User-Agent", s.userAgent)

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, common.NewClientError("failed to send request", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		err := common.ParseErrorResponse(resp)
		return nil, err
	}

	// Parse the response
	var importResult ImportResult
	if err := json.NewDecoder(resp.Body).Decode(&importResult); err != nil {
		return nil, common.NewClientError("failed to parse response", err)
	}

	return &importResult, nil
}
