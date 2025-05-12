# Bunny.net Go API Client

> **⚠️ UNDER CONSTRUCTION ⚠️**  
> This client library is currently under active development and is not yet considered stable. The API may change without notice. Use at your own risk.

A Go client library for the [Bunny.net API](https://docs.bunny.net/reference/bunnynet-api-overview).

## Installation

```bash
go get github.com/venom90/bunnynet-go
```

## Usage

### Creating a Client

```go
import (
    "github.com/venom90/bunnynet-go"
    "time"
)

// Create a client with default options
client := bunnynet.NewClient("your-api-key")

// Create a client with custom options
client := bunnynet.NewClient(
    "your-api-key",
    bunnynet.WithBaseURL("https://api.bunny.net"),
    bunnynet.WithTimeout(30 * time.Second),
    bunnynet.WithUserAgent("my-app/1.0.0"),
)
```

### Using the Country API

```go
import (
    "context"
    "fmt"
    "github.com/venom90/bunnynet-go"
    "github.com/venom90/bunnynet-go/common"
)

func main() {
    client := bunnynet.NewClient("your-api-key")

    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // List all countries
    countries, err := client.Country.List(ctx)
    if err != nil {
        panic(err)
    }

    for _, country := range countries {
        fmt.Printf("Country: %s (%s), EU: %t, Tax Rate: %.2f%%\n",
            country.Name, country.IsoCode, country.IsEU, country.TaxRate)
    }

    // Get a specific country by ISO code
    us, err := client.Country.Get(ctx, "US")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Country: %s (%s), EU: %t, Tax Rate: %.2f%%\n",
        us.Name, us.IsoCode, us.IsEU, us.TaxRate)
}
```

### Using the DNS Zone Resource

```go
import (
    "context"
    "fmt"
    "github.com/venom90/bunnynet-go"
    "github.com/venom90/bunnynet-go/common"
    "github.com/venom90/bunnynet-go/resources"
)

func main() {
    client := bunnynet.NewClient("your-api-key")
    ctx := context.Background()

    // List DNS zones with pagination
    pagination := common.NewPagination().WithPerPage(10)
    response, err := client.DNSZone.List(ctx, pagination, "")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d DNS zones\n", response.TotalItems)
    for _, zone := range response.Items {
        fmt.Printf("DNS Zone: ID=%d, Domain=%s, Records=%d\n",
            zone.Id, zone.Domain, len(zone.Records))
    }

    // Create a new DNS zone
    newZone, err := client.DNSZone.Add(ctx, resources.AddDNSZoneOptions{
        Domain: "example.com",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created new DNS zone: %s (ID: %d)\n", newZone.Domain, newZone.Id)

    // Add an A record
    record, err := client.DNSZone.AddRecord(ctx, newZone.Id, resources.AddDNSRecordOptions{
        Type:  resources.DNSRecordTypeA,
        Name:  "@",
        Value: "192.0.2.1",
        Ttl:   3600,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Added A record: ID=%d, Name=%s, Value=%s\n",
        record.Id, record.Name, record.Value)

    // Enable DNSSEC
    dnsSecInfo, err := client.DNSZone.EnableDNSSec(ctx, newZone.Id)
    if err != nil {
        panic(err)
    }

    fmt.Printf("DNSSEC enabled: %t, DS Record: %s\n",
        dnsSecInfo.Enabled, dnsSecInfo.DsRecord)

    // Export zone file
    zoneData, err := client.DNSZone.Export(ctx, newZone.Id)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Exported zone file (%d bytes)\n", len(zoneData))

    // Delete the zone
    err = client.DNSZone.Delete(ctx, newZone.Id)
    if err != nil {
        panic(err)
    }

    fmt.Println("DNS zone deleted successfully")
}
```

### Using the API Key Resource

```go
import (
    "context"
    "fmt"
    "github.com/venom90/bunnynet-go"
    "github.com/venom90/bunnynet-go/common"
)

func main() {
    client := bunnynet.NewClient("your-api-key")
    ctx := context.Background()

    // List API keys with pagination
    pagination := common.NewPagination().WithPerPage(10)
    response, err := client.APIKey.List(ctx, pagination)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d API keys\n", response.TotalItems)
    for _, apiKey := range response.Items {
        fmt.Printf("API Key: ID=%d, Roles=%v\n", apiKey.Id, apiKey.Roles)
    }

    // Create a new API key
    roles := []string{"PullZone.Read", "Statistics.Read"}
    newKey, err := client.APIKey.Create(ctx, roles)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created new API key: %s\n", newKey.Key)

    // Get an API key by ID
    apiKey, err := client.APIKey.Get(ctx, newKey.Id)
    if err != nil {
        panic(err)
    }

    fmt.Printf("API Key: ID=%d, Roles=%v\n", apiKey.Id, apiKey.Roles)

    // Delete an API key
    err = client.APIKey.Delete(ctx, newKey.Id)
    if err != nil {
        panic(err)
    }

    fmt.Println("API key deleted successfully")
}
```

## Using the Pull Zone Resource

```go
import (
    "context"
    "fmt"
    "github.com/venom90/bunnynet-go"
    "github.com/venom90/bunnynet-go/common"
    "github.com/venom90/bunnynet-go/resources"
    "time"
)

func main() {
    client := bunnynet.NewClient("your-api-key")
    ctx := context.Background()

    // List Pull Zones with pagination
    pagination := common.NewPagination().WithPerPage(10)
    response, err := client.PullZone.List(ctx, pagination, "", false)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d Pull Zones\n", response.TotalItems)
    for _, zone := range response.Items {
        fmt.Printf("Pull Zone: ID=%d, Name=%s, Origin=%s\n",
            zone.Id, zone.Name, zone.OriginUrl)
    }

    // Create a new Pull Zone
    newZone, err := client.PullZone.Add(ctx, resources.AddPullZoneOptions{
        Name:             "example-zone",
        OriginUrl:        "https://example.com",
        Type:             0, // Premium
        EnableGeoZoneUS:  true,
        EnableGeoZoneEU:  true,
        EnableGeoZoneASIA: true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created new Pull Zone: %s (ID: %d)\n", newZone.Name, newZone.Id)

    // Add a hostname
    err = client.PullZone.AddHostname(ctx, newZone.Id, resources.AddHostnameOptions{
        Hostname: "cdn.example.com",
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("Added hostname: cdn.example.com")

    // Load a free SSL certificate
    err = client.PullZone.LoadFreeCertificate(ctx, resources.LoadFreeCertificateOptions{
        Hostname: "cdn.example.com",
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("Loaded free SSL certificate for cdn.example.com")

    // Add an edge rule for forcing SSL
    err = client.PullZone.AddOrUpdateEdgeRule(ctx, newZone.Id, resources.AddOrUpdateEdgeRuleOptions{
        ActionType: 0, // ForceSSL
        Triggers: []resources.EdgeRuleTrigger{
            {
                Type:               0, // URL
                PatternMatches:     []string{"/*"},
                PatternMatchingType: 0, // MatchAny
                TriggerMatchingType: 0, // MatchAny
            },
        },
        Description: "Force SSL for all URLs",
        Enabled:     true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("Added edge rule to force SSL")

    // Purge cache
    err = client.PullZone.PurgeCache(ctx, newZone.Id, nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("Purged cache for Pull Zone")

    // Get the origin shield queue statistics
    now := time.Now()
    yesterday := now.AddDate(0, 0, -1)
    stats, err := client.PullZone.GetOriginShieldQueueStatistics(ctx, newZone.Id, &resources.StatisticsOptions{
        DateFrom: &yesterday,
        DateTo:   &now,
        Hourly:   true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("Retrieved origin shield queue statistics")

    // Update the Pull Zone configuration
    updatedZone, err := client.PullZone.Update(ctx, newZone.Id, &resources.PullZone{
        EnableWebPVary:     true,
        EnableAvifVary:     true,
        CacheControlMaxAgeOverride: 86400, // 1 day
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Updated Pull Zone: %s\n", updatedZone.Name)

    // Delete the Pull Zone
    err = client.PullZone.Delete(ctx, newZone.Id)
    if err != nil {
        panic(err)
    }

    fmt.Println("Pull Zone deleted successfully")
}
```

## Using the Purge Service

The Purge service allows you to purge a specific URL from the Bunny.net CDN cache to ensure that fresh content is delivered to your users.

```go
import (
    "context"
    "fmt"
    "github.com/venom90/bunnynet-go-client"
    "github.com/venom90/bunnynet-go-client/resources"
)

func main() {
    client := bunnynet.NewClient("your-api-key")
    ctx := context.Background()

    // Simple purge (synchronous)
    // This will wait for the purge to complete before returning
    err := client.Purge.Purge(ctx, "https://example.com/file.jpg", false)
    if err != nil {
        panic(err)
    }
    fmt.Println("Successfully purged URL")

    // Asynchronous purge
    // This will initiate the purge and return immediately
    err = client.Purge.Purge(ctx, "https://example.com/folder/", true)
    if err != nil {
        panic(err)
    }
    fmt.Println("Purge request initiated")

    // Using PurgeOptions for more control
    options := resources.PurgeOptions{
        URL:   "https://example.com/api/data.json",
        Async: true,
    }
    err = client.Purge.PurgeURL(ctx, options)
    if err != nil {
        panic(err)
    }
    fmt.Println("Purge request initiated with options")
}
```

The Purge service supports two modes:

1. **Synchronous purging:** The API call will wait for the purge operation to complete before returning.
2. **Asynchronous purging:** The API call will initiate the purge operation and return immediately, without waiting for completion.

### Purging can help when:

- You've updated content and want to ensure the latest version is being served
- You need to remove outdated or incorrect content from the cache
- You're troubleshooting caching issues

## Pagination

The client supports three approaches to pagination:

#### 1. Manual Pagination

```go
// Create pagination options
pagination := common.NewPagination().WithPage(1).WithPerPage(10)

// Get the first page
response, err := client.Country.ListPaginated(ctx, pagination)
if err != nil {
    panic(err)
}

// Process items in the first page
for _, item := range response.Items {
    // Process each item
}

// Check if there are more pages
if response.HasMoreItems {
    // Update pagination for the next page
    pagination.WithPage(response.CurrentPage + 1)

    // Get the next page
    nextPage, err := client.Country.ListPaginated(ctx, pagination)
    // ...
}
```

#### 2. Using the Page Iterator

```go
// Create a page iterator
iterator := common.NewPageIterator(
    func(page, perPage int) (*common.PaginatedResponse[resources.Country], error) {
        pagination := common.NewPagination().WithPage(page).WithPerPage(perPage)
        return client.Country.ListPaginated(ctx, pagination)
    },
    1, // starting page
    10, // items per page
)

// Iterate through all pages
for iterator.Next() {
    // Get the current page info
    pageInfo := iterator.PageInfo()
    fmt.Printf("Page %d of %d\n", pageInfo.CurrentPage, pageInfo.TotalItems)

    // Process items in the current page
    for _, item := range iterator.Items() {
        // Process each item
    }
}

// Check for errors after iteration
if err := iterator.Error(); err != nil {
    panic(err)
}
```

#### 3. Getting All Items at Once

```go
// Get all items across all pages
allItems, err := client.Country.ListAll(ctx, 25) // 25 items per page
if err != nil {
    panic(err)
}

// Process all items
for _, item := range allItems {
    // Process each item
}
```

## Features

- Idiomatic Go API client
- Full support for Bunny.net API resources
- Advanced pagination support with multiple approaches:
  - Manual pagination with page and perPage parameters
  - Page iterator for convenient page-by-page processing
  - Utility methods to fetch all items at once
- Comprehensive error handling with detailed error information
- Context support for cancellation and timeouts
- Customizable client options through functional options pattern
- Comprehensive test coverage with mocks for reliable testing

## Available Resources

- Country: List and retrieve country information
- API Key: List, create, retrieve, and delete API keys
- DNS Zone: Manage DNS zones and records
- More resources coming soon...

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This library is distributed under the MIT license.
