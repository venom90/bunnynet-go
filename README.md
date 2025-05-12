# Bunny.net Go API Client

> **⚠️ UNDER CONSTRUCTION ⚠️**  
> This client library is currently under active development and is not yet considered stable. The API may change without notice. Use at your own risk.

A Go client library for the [Bunny.net API](https://docs.bunny.net/reference/bunnynet-api-overview).

## Installation

```bash
go get github.com/venom90/bunnynet-go-client
```

## Usage

### Creating a Client

```go
import (
    "github.com/venom90/bunnynet-go-client"
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
    "github.com/venom90/bunnynet-go-client"
    "github.com/venom90/bunnynet-go-client/common"
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
    "github.com/venom90/bunnynet-go-client"
    "github.com/venom90/bunnynet-go-client/common"
    "github.com/venom90/bunnynet-go-client/resources"
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
    "github.com/venom90/bunnynet-go-client"
    "github.com/venom90/bunnynet-go-client/common"
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

### Pagination

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
