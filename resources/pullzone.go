// Package resources provides API resource implementations for the Bunny.net API client
package resources

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/venom90/bunnynet-go-client/common"
	"github.com/venom90/bunnynet-go-client/internal"
)

// PullZone represents a Pull Zone in the Bunny.net API
type PullZone struct {
	// Id is the unique identifier of the pull zone
	Id int64 `json:"Id"`

	// Name is the name of the pull zone
	Name string `json:"Name"`

	// OriginUrl is the origin URL of the pull zone where the files are fetched from
	OriginUrl string `json:"OriginUrl"`

	// Enabled determines if the Pull Zone is currently enabled, active and running
	Enabled bool `json:"Enabled"`

	// Hostnames is the list of hostnames linked to this Pull Zone
	Hostnames []Hostname `json:"Hostnames"`

	// StorageZoneId is the ID of the storage zone that the pull zone is linked to
	StorageZoneId int64 `json:"StorageZoneId"`

	// EdgeScriptId is the ID of the edge script that the pull zone is linked to
	EdgeScriptId int64 `json:"EdgeScriptId"`

	// AllowedReferrers is the list of referrer hostnames that are allowed to access the pull zone
	AllowedReferrers []string `json:"AllowedReferrers"`

	// BlockedReferrers is the list of referrer hostnames that are blocked from accessing the pull zone
	BlockedReferrers []string `json:"BlockedReferrers"`

	// BlockedIps is the list of IPs that are blocked from accessing the pull zone
	BlockedIps []string `json:"BlockedIps"`

	// EnableGeoZoneUS determines if the delivery from the North American region is enabled for this pull zone
	EnableGeoZoneUS bool `json:"EnableGeoZoneUS"`

	// EnableGeoZoneEU determines if the delivery from the European region is enabled for this pull zone
	EnableGeoZoneEU bool `json:"EnableGeoZoneEU"`

	// EnableGeoZoneASIA determines if the delivery from the Asian / Oceanian region is enabled for this pull zone
	EnableGeoZoneASIA bool `json:"EnableGeoZoneASIA"`

	// EnableGeoZoneSA determines if the delivery from the South American region is enabled for this pull zone
	EnableGeoZoneSA bool `json:"EnableGeoZoneSA"`

	// EnableGeoZoneAF determines if the delivery from the Africa region is enabled for this pull zone
	EnableGeoZoneAF bool `json:"EnableGeoZoneAF"`

	// ZoneSecurityEnabled is true if the URL secure token authentication security is enabled
	ZoneSecurityEnabled bool `json:"ZoneSecurityEnabled"`

	// ZoneSecurityKey is the security key used for secure URL token authentication
	ZoneSecurityKey string `json:"ZoneSecurityKey"`

	// ZoneSecurityIncludeHashRemoteIP is true if the zone security hash should include the remote IP
	ZoneSecurityIncludeHashRemoteIP bool `json:"ZoneSecurityIncludeHashRemoteIP"`

	// IgnoreQueryStrings is true if the Pull Zone is ignoring query strings when serving cached objects
	IgnoreQueryStrings bool `json:"IgnoreQueryStrings"`

	// MonthlyBandwidthLimit is the monthly limit of bandwidth in bytes that the pullzone is allowed to use
	MonthlyBandwidthLimit int64 `json:"MonthlyBandwidthLimit"`

	// MonthlyBandwidthUsed is the amount of bandwidth in bytes that the pull zone used this month
	MonthlyBandwidthUsed int64 `json:"MonthlyBandwidthUsed"`

	// MonthlyCharges is the total monthly charges for this so zone so far
	MonthlyCharges float64 `json:"MonthlyCharges"`

	// AddHostHeader determines if the Pull Zone should forward the current hostname to the origin
	AddHostHeader bool `json:"AddHostHeader"`

	// OriginHostHeader determines the host header that will be sent to the origin
	OriginHostHeader string `json:"OriginHostHeader"`

	// Type is the type of pull zone (0 = Premium, 1 = Volume)
	Type int `json:"Type"`

	// AccessControlOriginHeaderExtensions is the list of extensions that will return the CORS headers
	AccessControlOriginHeaderExtensions []string `json:"AccessControlOriginHeaderExtensions"`

	// EnableAccessControlOriginHeader determines if the CORS headers should be enabled
	EnableAccessControlOriginHeader bool `json:"EnableAccessControlOriginHeader"`

	// DisableCookies determines if the cookies are disabled for the pull zone
	DisableCookies bool `json:"DisableCookies"`

	// BudgetRedirectedCountries is the list of budget redirected countries with the two-letter Alpha2 ISO codes
	BudgetRedirectedCountries []string `json:"BudgetRedirectedCountries"`

	// BlockedCountries is the list of blocked countries with the two-letter Alpha2 ISO codes
	BlockedCountries []string `json:"BlockedCountries"`

	// EnableOriginShield if true the server will use the origin shield feature
	EnableOriginShield bool `json:"EnableOriginShield"`

	// CacheControlMaxAgeOverride is the override cache time for the pull zone
	CacheControlMaxAgeOverride int64 `json:"CacheControlMaxAgeOverride"`

	// CacheControlPublicMaxAgeOverride is the override cache time for the pull zone for the end client
	CacheControlPublicMaxAgeOverride int64 `json:"CacheControlPublicMaxAgeOverride"`

	// BurstSize - excessive requests are delayed until their number exceeds the maximum burst size
	BurstSize int32 `json:"BurstSize"`

	// RequestLimit - max number of requests per IP per second
	RequestLimit int32 `json:"RequestLimit"`

	// BlockRootPathAccess if true, access to root path will return a 403 error
	BlockRootPathAccess bool `json:"BlockRootPathAccess"`

	// BlockPostRequests if true, POST requests to the zone will be blocked
	BlockPostRequests bool `json:"BlockPostRequests"`

	// LimitRatePerSecond is the maximum rate at which the zone will transfer data in kb/s. 0 for unlimited
	LimitRatePerSecond float64 `json:"LimitRatePerSecond"`

	// LimitRateAfter is the amount of data after the rate limit will be activated
	LimitRateAfter float64 `json:"LimitRateAfter"`

	// ConnectionLimitPerIPCount is the number of connections limited per IP for this zone
	ConnectionLimitPerIPCount int32 `json:"ConnectionLimitPerIPCount"`

	// AddCanonicalHeader determines if the Add Canonical Header is enabled for this Pull Zone
	AddCanonicalHeader bool `json:"AddCanonicalHeader"`

	// EnableLogging determines if the logging is enabled for this Pull Zone
	EnableLogging bool `json:"EnableLogging"`

	// EnableCacheSlice determines if the cache slice (Optimize for video) feature is enabled for the Pull Zone
	EnableCacheSlice bool `json:"EnableCacheSlice"`

	// EnableSmartCache determines if smart caching is enabled for this zone
	EnableSmartCache bool `json:"EnableSmartCache"`

	// EdgeRules is the list of edge rules on this Pull Zone
	EdgeRules []EdgeRule `json:"EdgeRules"`

	// EnableWebPVary determines if the WebP Vary feature is enabled
	EnableWebPVary bool `json:"EnableWebPVary"`

	// EnableAvifVary determines if the AVIF Vary feature is enabled
	EnableAvifVary bool `json:"EnableAvifVary"`

	// EnableCountryCodeVary determines if the Country Code Vary feature is enabled
	EnableCountryCodeVary bool `json:"EnableCountryCodeVary"`

	// EnableMobileVary determines if the Mobile Vary feature is enabled
	EnableMobileVary bool `json:"EnableMobileVary"`

	// EnableCookieVary determines if the Cookie Vary feature is enabled
	EnableCookieVary bool `json:"EnableCookieVary"`

	// CookieVaryParameters contains the list of vary parameters that will be used for vary cache by cookie string
	CookieVaryParameters []string `json:"CookieVaryParameters"`

	// EnableHostnameVary determines if the Hostname Vary feature is enabled
	EnableHostnameVary bool `json:"EnableHostnameVary"`

	// CnameDomain is the CNAME domain of the pull zone for setting up custom hostnames
	CnameDomain string `json:"CnameDomain"`

	// AWSSigningEnabled determines if the AWS Signing is enabled
	AWSSigningEnabled bool `json:"AWSSigningEnabled"`

	// AWSSigningKey is the AWS Signing region key
	AWSSigningKey string `json:"AWSSigningKey"`

	// AWSSigningSecret is the AWS Signing region secret
	AWSSigningSecret string `json:"AWSSigningSecret"`

	// AWSSigningRegionName is the AWS Signing region name
	AWSSigningRegionName string `json:"AWSSigningRegionName"`

	// LoggingIPAnonymizationEnabled determines if IP anonymization is enabled for logs
	LoggingIPAnonymizationEnabled bool `json:"LoggingIPAnonymizationEnabled"`

	// EnableTLS1 determines if the TLS 1 is enabled on the Pull Zone
	EnableTLS1 bool `json:"EnableTLS1"`

	// EnableTLS1_1 determines if the TLS 1.1 is enabled on the Pull Zone
	EnableTLS1_1 bool `json:"EnableTLS1_1"`

	// VerifyOriginSSL determines if the Pull Zone should verify the origin SSL certificate
	VerifyOriginSSL bool `json:"VerifyOriginSSL"`

	// LogForwardingEnabled determines if the log forwarding is enabled
	LogForwardingEnabled bool `json:"LogForwardingEnabled"`

	// LogForwardingHostname is the log forwarding hostname
	LogForwardingHostname string `json:"LogForwardingHostname"`

	// LogForwardingPort is the log forwarding port
	LogForwardingPort int32 `json:"LogForwardingPort"`

	// LogForwardingToken is the log forwarding token value
	LogForwardingToken string `json:"LogForwardingToken"`

	// LogForwardingProtocol is the protocol used for log forwarding (0 = UDP, 1 = TCP, 2 = TCPEncrypted, 3 = DataDog)
	LogForwardingProtocol int `json:"LogForwardingProtocol"`

	// LoggingSaveToStorage determines if the permanent logging feature is enabled
	LoggingSaveToStorage bool `json:"LoggingSaveToStorage"`

	// LoggingStorageZoneId is the ID of the logging storage zone that is configured for this Pull Zone
	LoggingStorageZoneId int64 `json:"LoggingStorageZoneId"`

	// FollowRedirects determines if the zone will follow origin redirects
	FollowRedirects bool `json:"FollowRedirects"`

	// OriginRetries is the number of retries to the origin server
	OriginRetries int32 `json:"OriginRetries"`

	// OriginConnectTimeout is the amount of seconds to wait when connecting to the origin
	OriginConnectTimeout int32 `json:"OriginConnectTimeout"`

	// OriginResponseTimeout is the amount of seconds to wait when waiting for the origin reply
	OriginResponseTimeout int32 `json:"OriginResponseTimeout"`

	// UseStaleWhileUpdating determines if we should use stale cache while cache is updating
	UseStaleWhileUpdating bool `json:"UseStaleWhileUpdating"`

	// UseStaleWhileOffline determines if we should use stale cache while the origin is offline
	UseStaleWhileOffline bool `json:"UseStaleWhileOffline"`

	// OriginRetry5XXResponses determines if we should retry the request in case of a 5XX response
	OriginRetry5XXResponses bool `json:"OriginRetry5XXResponses"`

	// OriginRetryConnectionTimeout determines if we should retry the request in case of a connection timeout
	OriginRetryConnectionTimeout bool `json:"OriginRetryConnectionTimeout"`

	// OriginRetryResponseTimeout determines if we should retry the request in case of a response timeout
	OriginRetryResponseTimeout bool `json:"OriginRetryResponseTimeout"`

	// OriginRetryDelay determines the amount of time that the CDN should wait before retrying an origin request
	OriginRetryDelay int32 `json:"OriginRetryDelay"`

	// QueryStringVaryParameters contains the list of vary parameters for vary cache by query string
	QueryStringVaryParameters []string `json:"QueryStringVaryParameters"`

	// OriginShieldEnableConcurrencyLimit determines if the origin shield concurrency limit is enabled
	OriginShieldEnableConcurrencyLimit bool `json:"OriginShieldEnableConcurrencyLimit"`

	// OriginShieldMaxConcurrentRequests determines the number of maximum concurrent requests allowed to the origin
	OriginShieldMaxConcurrentRequests int32 `json:"OriginShieldMaxConcurrentRequests"`

	// EnableSafeHop enables the SafeHop feature
	EnableSafeHop bool `json:"EnableSafeHop"`

	// CacheErrorResponses determines if bunny.net should be caching error responses
	CacheErrorResponses bool `json:"CacheErrorResponses"`

	// OriginShieldQueueMaxWaitTime determines the max queue wait time
	OriginShieldQueueMaxWaitTime int32 `json:"OriginShieldQueueMaxWaitTime"`

	// OriginShieldMaxQueuedRequests determines the max number of origin requests that will remain in the queue
	OriginShieldMaxQueuedRequests int32 `json:"OriginShieldMaxQueuedRequests"`

	// UseBackgroundUpdate determines if cache update is performed in the background
	UseBackgroundUpdate bool `json:"UseBackgroundUpdate"`

	// EnableAutoSSL if set to true, any hostnames added to this Pull Zone will automatically enable SSL
	EnableAutoSSL bool `json:"EnableAutoSSL"`

	// EnableQueryStringOrdering if set to true the query string ordering property is enabled
	EnableQueryStringOrdering bool `json:"EnableQueryStringOrdering"`

	// LogAnonymizationType sets the type of log anonymization (0 = OneDigit, 1 = Drop)
	LogAnonymizationType int `json:"LogAnonymizationType"`

	// LogFormat sets the log format (0 = Plain, 1 = JSON)
	LogFormat int `json:"LogFormat"`

	// LogForwardingFormat sets the log forwarding format (0 = Plain, 1 = JSON)
	LogForwardingFormat int `json:"LogForwardingFormat"`

	// OriginType sets the origin type (0 = OriginUrl, 1 = DnsAccelerate, etc)
	OriginType int `json:"OriginType"`

	// EnableRequestCoalescing determines if request coalescing is currently enabled
	EnableRequestCoalescing bool `json:"EnableRequestCoalescing"`

	// RequestCoalescingTimeout determines the lock time for coalesced requests
	RequestCoalescingTimeout int32 `json:"RequestCoalescingTimeout"`

	// DisableLetsEncrypt if true, the built-in let's encrypt is disabled and requests are passed to the origin
	DisableLetsEncrypt bool `json:"DisableLetsEncrypt"`

	// PreloadingScreenEnabled determines if the preloading screen is currently enabled
	PreloadingScreenEnabled bool `json:"PreloadingScreenEnabled"`

	// PreloadingScreenLogoUrl is the preloading screen logo URL
	PreloadingScreenLogoUrl string `json:"PreloadingScreenLogoUrl"`

	// Additional fields can be added as needed
}

// Hostname represents a hostname linked to a Pull Zone
type Hostname struct {
	// Id is the unique ID of the hostname
	Id int64 `json:"Id"`

	// Value is the hostname value for the domain name
	Value string `json:"Value"`

	// ForceSSL determines if the Force SSL feature is enabled
	ForceSSL bool `json:"ForceSSL"`

	// IsSystemHostname determines if this is a system hostname controlled by bunny.net
	IsSystemHostname bool `json:"IsSystemHostname"`

	// HasCertificate determines if the hostname has an SSL certificate configured
	HasCertificate bool `json:"HasCertificate"`

	// Certificate contains the Base64 encoded certificate for the hostname
	Certificate string `json:"Certificate,omitempty"`

	// CertificateKey contains the Base64 encoded certificate key for the hostname
	CertificateKey string `json:"CertificateKey,omitempty"`
}

// EdgeRule represents an edge rule on a Pull Zone
type EdgeRule struct {
	// Guid is the unique GUID of the edge rule
	Guid string `json:"Guid"`

	// ActionType is the type of action that the edge rule performs
	ActionType int `json:"ActionType"`

	// ActionParameter1 is the action parameter 1
	ActionParameter1 string `json:"ActionParameter1"`

	// ActionParameter2 is the action parameter 2
	ActionParameter2 string `json:"ActionParameter2"`

	// Triggers is the list of triggers for this edge rule
	Triggers []EdgeRuleTrigger `json:"Triggers"`

	// Description is the description of the edge rule
	Description string `json:"Description"`

	// Enabled determines if the edge rule is currently enabled or not
	Enabled bool `json:"Enabled"`
}

// EdgeRuleTrigger represents a trigger for an edge rule
type EdgeRuleTrigger struct {
	// Type is the type of trigger
	Type int `json:"Type"`

	// PatternMatches is the list of pattern matches that will trigger the edge rule
	PatternMatches []string `json:"PatternMatches"`

	// PatternMatchingType defines how patterns should be matched
	PatternMatchingType int `json:"PatternMatchingType"`

	// Parameter1 is the trigger parameter 1
	Parameter1 string `json:"Parameter1"`

	// TriggerMatchingType defines how triggers should be matched
	TriggerMatchingType int `json:"TriggerMatchingType"`
}

// Add PullZone request parameters
type AddPullZoneOptions struct {
	// Name is the name of the pull zone
	Name string `json:"Name"`

	// OriginUrl is the origin URL of the Pull Zone
	OriginUrl string `json:"OriginUrl"`

	// Type is the type of pull zone (0 = Premium, 1 = Volume)
	Type int `json:"Type,omitempty"`

	// Additional configuration parameters can be added here
	// The following are just some examples
	AllowedReferrers  []string `json:"AllowedReferrers,omitempty"`
	BlockedReferrers  []string `json:"BlockedReferrers,omitempty"`
	BlockedIps        []string `json:"BlockedIps,omitempty"`
	EnableGeoZoneUS   bool     `json:"EnableGeoZoneUS,omitempty"`
	EnableGeoZoneEU   bool     `json:"EnableGeoZoneEU,omitempty"`
	EnableGeoZoneASIA bool     `json:"EnableGeoZoneASIA,omitempty"`
	EnableGeoZoneSA   bool     `json:"EnableGeoZoneSA,omitempty"`
	EnableGeoZoneAF   bool     `json:"EnableGeoZoneAF,omitempty"`
	// Other options - can be expanded as needed
}

// AddHostnameOptions represents the options for adding a hostname to a pull zone
type AddHostnameOptions struct {
	// Hostname is the hostname that will be added
	Hostname string `json:"Hostname"`
}

// RemoveHostnameOptions represents the options for removing a hostname from a pull zone
type RemoveHostnameOptions struct {
	// Hostname is the hostname that will be removed
	Hostname string `json:"Hostname"`
}

// AddCertificateOptions represents the options for adding a certificate to a hostname
type AddCertificateOptions struct {
	// Hostname is the hostname to which the certificate will be added
	Hostname string `json:"Hostname"`

	// Certificate is the Base64 encoded binary data of the certificate file
	Certificate string `json:"Certificate"`

	// CertificateKey is the Base64 encoded binary data of the certificate key file
	CertificateKey string `json:"CertificateKey"`
}

// RemoveCertificateOptions represents the options for removing a certificate from a hostname
type RemoveCertificateOptions struct {
	// Hostname is the hostname from which the certificate will be removed
	Hostname string `json:"Hostname"`
}

// SetForceSSLOptions represents the options for setting the Force SSL option on a hostname
type SetForceSSLOptions struct {
	// Hostname is the hostname that will be updated
	Hostname string `json:"Hostname"`

	// ForceSSL set to true to force SSL on the given pull zone hostname
	ForceSSL bool `json:"ForceSSL"`
}

// HostnameOptions represents the options for operations on a hostname
type HostnameOptions struct {
	// Hostname is the hostname to operate on
	Hostname string `json:"Hostname"`
}

// BlockedIPOptions represents the options for operations on blocked IPs
type BlockedIPOptions struct {
	// BlockedIp is the IP address to block or unblock
	BlockedIp string `json:"BlockedIp"`
}

// AddOrUpdateEdgeRuleOptions represents the options for adding or updating an edge rule
type AddOrUpdateEdgeRuleOptions struct {
	// Guid is the unique GUID of the edge rule
	Guid string `json:"Guid,omitempty"`

	// ActionType is the type of action that the edge rule performs
	ActionType int `json:"ActionType"`

	// ActionParameter1 is the action parameter 1
	ActionParameter1 string `json:"ActionParameter1,omitempty"`

	// ActionParameter2 is the action parameter 2
	ActionParameter2 string `json:"ActionParameter2,omitempty"`

	// Triggers is the list of triggers for this edge rule
	Triggers []EdgeRuleTrigger `json:"Triggers"`

	// Description is the description of the edge rule
	Description string `json:"Description,omitempty"`

	// Enabled determines if the edge rule is currently enabled or not
	Enabled bool `json:"Enabled"`
}

// SetEdgeRuleEnabledOptions represents the options for enabling or disabling an edge rule
type SetEdgeRuleEnabledOptions struct {
	// Id is the ID of the object to update
	Id int64 `json:"Id"`

	// Value is the boolean value to set
	Value bool `json:"Value"`
}

// PurgeCacheOptions represents the options for purging a cache
type PurgeCacheOptions struct {
	// CacheTag is an optional tag that can be used to target specific cached objects
	CacheTag string `json:"CacheTag,omitempty"`
}

// CheckAvailabilityOptions represents the options for checking pull zone name availability
type CheckAvailabilityOptions struct {
	// Name is the name of the zone to check
	Name string `json:"Name"`
}

// CheckAvailabilityResponse represents the response from a check availability request
type CheckAvailabilityResponse struct {
	// Available indicates if the name is available
	Available bool `json:"Available"`
}

// OriginShieldQueueStatistics represents the statistics for the origin shield queue
type OriginShieldQueueStatistics struct {
	// ConcurrentRequestsChart is the constructed chart of origin shield concurrent requests
	ConcurrentRequestsChart map[string]interface{} `json:"ConcurrentRequestsChart"`

	// QueuedRequestsChart is the constructed chart of origin shield requests chart
	QueuedRequestsChart map[string]interface{} `json:"QueuedRequestsChart"`
}

// OptimizerStatistics represents the statistics for the optimizer
type OptimizerStatistics struct {
	// RequestsOptimizedChart is the constructed chart of optimized requests
	RequestsOptimizedChart map[string]interface{} `json:"RequestsOptimizedChart"`

	// AverageCompressionChart is the average compression chart of the responses
	AverageCompressionChart map[string]interface{} `json:"AverageCompressionChart"`

	// TrafficSavedChart is the constructed chart of saved traffic
	TrafficSavedChart map[string]interface{} `json:"TrafficSavedChart"`

	// AverageProcessingTimeChart is the constructed chart of processing time
	AverageProcessingTimeChart map[string]interface{} `json:"AverageProcessingTimeChart"`

	// TotalRequestsOptimized is the total number of optimized requests
	TotalRequestsOptimized float64 `json:"TotalRequestsOptimized"`

	// TotalTrafficSaved is the total requests saved
	TotalTrafficSaved float64 `json:"TotalTrafficSaved"`

	// AverageProcessingTime is the average processing time of each request
	AverageProcessingTime float64 `json:"AverageProcessingTime"`

	// AverageCompressionRatio is the average compression ratio of CDN responses
	AverageCompressionRatio float64 `json:"AverageCompressionRatio"`
}

// StatisticsOptions represents the options for requesting statistics
type StatisticsOptions struct {
	// DateFrom is the start date of the statistics
	DateFrom *time.Time `url:"dateFrom,omitempty" json:"dateFrom,omitempty"`

	// DateTo is the end date of the statistics
	DateTo *time.Time `url:"dateTo,omitempty" json:"dateTo,omitempty"`

	// Hourly if true, the statistics data will be returned in hourly grouping
	Hourly bool `url:"hourly,omitempty" json:"hourly,omitempty"`
}

// LoadFreeCertificateOptions represents the options for loading a free certificate
type LoadFreeCertificateOptions struct {
	// Hostname is the hostname that the certificate will be loaded for
	Hostname string `url:"hostname" json:"hostname"`
}

// ToQueryParams converts the StatisticsOptions to query parameters
func (o *StatisticsOptions) ToQueryParams() map[string]string {
	params := make(map[string]string)

	if o.DateFrom != nil {
		params["dateFrom"] = o.DateFrom.UTC().Format(time.RFC3339)
	}

	if o.DateTo != nil {
		params["dateTo"] = o.DateTo.UTC().Format(time.RFC3339)
	}

	if o.Hourly {
		params["hourly"] = "true"
	}

	return params
}

// PullZoneService handles operations on pull zones
type PullZoneService struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	userAgent string
}

// NewPullZoneService creates a new PullZoneService
func NewPullZoneService(client *http.Client, baseURL, apiKey, userAgent string) *PullZoneService {
	return &PullZoneService{
		client:    client,
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: userAgent,
	}
}

// SetAPIKey updates the API key used for authentication
func (s *PullZoneService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// List returns a paginated list of pull zones
func (s *PullZoneService) List(ctx context.Context, pagination *common.Pagination, search string, includeCertificate bool) (*common.PaginatedResponse[PullZone], error) {
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, "/pullzone", nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add pagination parameters
	if err := internal.AddQueryParams(req, pagination); err != nil {
		return nil, err
	}

	// Add additional query parameters
	q := req.URL.Query()
	if search != "" {
		q.Add("search", search)
	}
	if includeCertificate {
		q.Add("includeCertificate", "true")
	}
	req.URL.RawQuery = q.Encode()

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var paginatedResponse common.PaginatedResponse[PullZone]
	if err := internal.ParsePaginatedResponse(resp, &paginatedResponse); err != nil {
		return nil, err
	}

	return &paginatedResponse, nil
}

// ListAll returns all pull zones across all pages
func (s *PullZoneService) ListAll(ctx context.Context, perPage int, search string, includeCertificate bool) ([]PullZone, error) {
	if perPage <= 0 {
		perPage = common.DefaultPerPage
	}

	iterator := common.NewPageIterator(
		func(page, itemsPerPage int) (*common.PaginatedResponse[PullZone], error) {
			pagination := common.NewPagination().WithPage(page).WithPerPage(itemsPerPage)
			return s.List(ctx, pagination, search, includeCertificate)
		},
		common.DefaultPage,
		perPage,
	)

	return iterator.AllItems()
}

// Get returns a pull zone by ID
func (s *PullZoneService) Get(ctx context.Context, id int64, includeCertificate bool) (*PullZone, error) {
	path := fmt.Sprintf("/pullzone/%d", id)
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add query parameters
	if includeCertificate {
		q := req.URL.Query()
		q.Add("includeCertificate", "true")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var pullZone PullZone
	if err := internal.ParseResponse(resp, &pullZone); err != nil {
		return nil, err
	}

	return &pullZone, nil
}

// Add creates a new pull zone
func (s *PullZoneService) Add(ctx context.Context, options AddPullZoneOptions) (*PullZone, error) {
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, "/pullzone", options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var pullZone PullZone
	if err := internal.ParseResponse(resp, &pullZone); err != nil {
		return nil, err
	}

	return &pullZone, nil
}

// Update updates an existing pull zone
func (s *PullZoneService) Update(ctx context.Context, id int64, options *PullZone) (*PullZone, error) {
	path := fmt.Sprintf("/pullzone/%d", id)
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, path, options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var updatedPullZone PullZone
	if err := internal.ParseResponse(resp, &updatedPullZone); err != nil {
		return nil, err
	}

	return &updatedPullZone, nil
}

// Delete deletes a pull zone
func (s *PullZoneService) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/pullzone/%d", id)
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

// PurgeCache purges the cache for a pull zone
func (s *PullZoneService) PurgeCache(ctx context.Context, id int64, options *PurgeCacheOptions) error {
	path := fmt.Sprintf("/pullzone/%d/purgeCache", id)
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

// AddHostname adds a hostname to a pull zone
func (s *PullZoneService) AddHostname(ctx context.Context, id int64, options AddHostnameOptions) error {
	path := fmt.Sprintf("/pullzone/%d/addHostname", id)
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

// RemoveHostname removes a hostname from a pull zone
func (s *PullZoneService) RemoveHostname(ctx context.Context, id int64, options RemoveHostnameOptions) error {
	path := fmt.Sprintf("/pullzone/%d/removeHostname", id)
	req, err := internal.NewRequest(http.MethodDelete, s.baseURL, path, options, s.apiKey, s.userAgent)
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

// AddCertificate adds a custom certificate to a hostname
func (s *PullZoneService) AddCertificate(ctx context.Context, id int64, options AddCertificateOptions) error {
	path := fmt.Sprintf("/pullzone/%d/addCertificate", id)
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

// RemoveCertificate removes a certificate from a hostname
func (s *PullZoneService) RemoveCertificate(ctx context.Context, id int64, options RemoveCertificateOptions) error {
	path := fmt.Sprintf("/pullzone/%d/removeCertificate", id)
	req, err := internal.NewRequest(http.MethodDelete, s.baseURL, path, options, s.apiKey, s.userAgent)
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

// SetForceSSL sets the Force SSL option on a hostname
func (s *PullZoneService) SetForceSSL(ctx context.Context, id int64, options SetForceSSLOptions) error {
	path := fmt.Sprintf("/pullzone/%d/setForceSSL", id)
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

// ResetSecurityKey resets the token key for a pull zone
func (s *PullZoneService) ResetSecurityKey(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/pullzone/%d/resetSecurityKey", id)
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, path, nil, s.apiKey, s.userAgent)
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

// AddAllowedReferrer adds an allowed referrer to a pull zone
func (s *PullZoneService) AddAllowedReferrer(ctx context.Context, id int64, options HostnameOptions) error {
	path := fmt.Sprintf("/pullzone/%d/addAllowedReferrer", id)
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

// RemoveAllowedReferrer removes an allowed referrer from a pull zone
func (s *PullZoneService) RemoveAllowedReferrer(ctx context.Context, id int64, options HostnameOptions) error {
	path := fmt.Sprintf("/pullzone/%d/removeAllowedReferrer", id)
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

// AddBlockedReferrer adds a blocked referrer to a pull zone
func (s *PullZoneService) AddBlockedReferrer(ctx context.Context, id int64, options HostnameOptions) error {
	path := fmt.Sprintf("/pullzone/%d/addBlockedReferrer", id)
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

// RemoveBlockedReferrer removes a blocked referrer from a pull zone
func (s *PullZoneService) RemoveBlockedReferrer(ctx context.Context, id int64, options HostnameOptions) error {
	path := fmt.Sprintf("/pullzone/%d/removeBlockedReferrer", id)
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

// AddBlockedIP adds a blocked IP to a pull zone
func (s *PullZoneService) AddBlockedIP(ctx context.Context, id int64, options BlockedIPOptions) error {
	path := fmt.Sprintf("/pullzone/%d/addBlockedIp", id)
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

// RemoveBlockedIP removes a blocked IP from a pull zone
func (s *PullZoneService) RemoveBlockedIP(ctx context.Context, id int64, options BlockedIPOptions) error {
	path := fmt.Sprintf("/pullzone/%d/removeBlockedIp", id)
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

// AddOrUpdateEdgeRule adds or updates an edge rule on a pull zone
func (s *PullZoneService) AddOrUpdateEdgeRule(ctx context.Context, pullZoneId int64, options AddOrUpdateEdgeRuleOptions) error {
	path := fmt.Sprintf("/pullzone/%d/edgerules/addOrUpdate", pullZoneId)
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

// DeleteEdgeRule deletes an edge rule from a pull zone
func (s *PullZoneService) DeleteEdgeRule(ctx context.Context, pullZoneId int64, edgeRuleId string) error {
	path := fmt.Sprintf("/pullzone/%d/edgerules/%s", pullZoneId, edgeRuleId)
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

// SetEdgeRuleEnabled enables or disables an edge rule
func (s *PullZoneService) SetEdgeRuleEnabled(ctx context.Context, pullZoneId int64, edgeRuleId string, options SetEdgeRuleEnabledOptions) error {
	path := fmt.Sprintf("/pullzone/%d/edgerules/%s/setEdgeRuleEnabled", pullZoneId, edgeRuleId)
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

// GetOriginShieldQueueStatistics retrieves the origin shield queue statistics for a pull zone
func (s *PullZoneService) GetOriginShieldQueueStatistics(ctx context.Context, pullZoneId int64, options *StatisticsOptions) (*OriginShieldQueueStatistics, error) {
	path := fmt.Sprintf("/pullzone/%d/originshield/queuestatistics", pullZoneId)
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add query parameters
	if options != nil {
		if err := internal.AddQueryParams(req, options); err != nil {
			return nil, err
		}
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var stats OriginShieldQueueStatistics
	if err := internal.ParseResponse(resp, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetOptimizerStatistics retrieves the optimizer statistics for a pull zone
func (s *PullZoneService) GetOptimizerStatistics(ctx context.Context, pullZoneId int64, options *StatisticsOptions) (*OptimizerStatistics, error) {
	path := fmt.Sprintf("/pullzone/%d/optimizer/statistics", pullZoneId)
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add query parameters
	if options != nil {
		if err := internal.AddQueryParams(req, options); err != nil {
			return nil, err
		}
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var stats OptimizerStatistics
	if err := internal.ParseResponse(resp, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// LoadFreeCertificate loads a free SSL certificate for a hostname
func (s *PullZoneService) LoadFreeCertificate(ctx context.Context, options LoadFreeCertificateOptions) error {
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, "/pullzone/loadFreeCertificate", nil, s.apiKey, s.userAgent)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	// Add query parameters
	q := req.URL.Query()
	q.Add("hostname", options.Hostname)
	req.URL.RawQuery = q.Encode()

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// CheckAvailability checks if a pull zone name is available
func (s *PullZoneService) CheckAvailability(ctx context.Context, options CheckAvailabilityOptions) (*CheckAvailabilityResponse, error) {
	req, err := internal.NewRequest(http.MethodPost, s.baseURL, "/pullzone/checkavailability", options, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var response CheckAvailabilityResponse
	if err := internal.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
