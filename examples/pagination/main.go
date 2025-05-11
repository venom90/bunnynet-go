package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/venom90/bunnynet-go-client"
	"github.com/venom90/bunnynet-go-client/common"
	"github.com/venom90/bunnynet-go-client/resources"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("BUNNYNET_API_KEY")
	if apiKey == "" {
		log.Fatal("BUNNYNET_API_KEY environment variable is not set")
	}

	// Create a new client
	client := bunnynet.NewClient(apiKey)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example 1: Manual pagination
	fmt.Println("Example 1: Manual pagination")
	manualPagination(ctx, client)

	// Example 2: Using the iterator
	fmt.Println("\nExample 2: Using the iterator")
	iteratorPagination(ctx, client)

	// Example 3: Get all items at once
	fmt.Println("\nExample 3: Get all items at once")
	getAllItems(ctx, client)
}

func manualPagination(ctx context.Context, client *bunnynet.Client) {
	// Create pagination options for 5 items per page
	pagination := common.NewPagination().WithPerPage(5)

	// Get the first page
	response, err := client.Country.ListPaginated(ctx, pagination)
	if err != nil {
		log.Fatalf("Failed to list countries: %v", err)
	}

	// Process the first page
	fmt.Printf("Page %d of countries (Total: %d, Has more: %t)\n",
		response.CurrentPage, response.TotalItems, response.HasMoreItems)

	for i, country := range response.Items {
		fmt.Printf("  %d. %s (%s)\n", i+1, country.Name, country.IsoCode)
	}

	// If there are more pages, get the next page
	if response.HasMoreItems {
		// Update the pagination for the next page
		pagination.WithPage(response.CurrentPage + 1)

		// Get the next page
		nextResponse, err := client.Country.ListPaginated(ctx, pagination)
		if err != nil {
			log.Fatalf("Failed to list countries: %v", err)
		}

		// Process the next page
		fmt.Printf("Page %d of countries (Total: %d, Has more: %t)\n",
			nextResponse.CurrentPage, nextResponse.TotalItems, nextResponse.HasMoreItems)

		for i, country := range nextResponse.Items {
			fmt.Printf("  %d. %s (%s)\n", i+1, country.Name, country.IsoCode)
		}
	}
}

func iteratorPagination(ctx context.Context, client *bunnynet.Client) {
	// Create an iterator for 5 items per page
	iterator := common.NewPageIterator(
		func(page, perPage int) (*common.PaginatedResponse[resources.Country], error) {
			pagination := common.NewPagination().WithPage(page).WithPerPage(perPage)
			return client.Country.ListPaginated(ctx, pagination)
		},
		common.DefaultPage,
		5,
	)

	// Iterate through all pages
	pageCount := 0
	for iterator.Next() {
		pageCount++
		pageInfo := iterator.PageInfo()

		fmt.Printf("Page %d of countries (Total: %d, Has more: %t)\n",
			pageInfo.CurrentPage, pageInfo.TotalItems, pageInfo.HasMoreItems)

		for i, country := range iterator.Items() {
			fmt.Printf("  %d. %s (%s)\n", i+1, country.Name, country.IsoCode)
		}

		// For this example, only process the first 2 pages
		if pageCount >= 2 {
			fmt.Println("Limiting to 2 pages for this example...")
			break
		}
	}

	// Check for errors
	if err := iterator.Error(); err != nil {
		log.Fatalf("Failed to iterate countries: %v", err)
	}
}

func getAllItems(ctx context.Context, client *bunnynet.Client) {
	// Get all countries at once
	countries, err := client.Country.ListAll(ctx, 10)
	if err != nil {
		log.Fatalf("Failed to list all countries: %v", err)
	}

	fmt.Printf("Retrieved all %d countries\n", len(countries))

	// Print the first 5 countries
	for i, country := range countries {
		if i >= 5 {
			fmt.Printf("... and %d more\n", len(countries)-5)
			break
		}
		fmt.Printf("  %d. %s (%s)\n", i+1, country.Name, country.IsoCode)
	}
}
