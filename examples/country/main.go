package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/venom90/bunnynet-go-client"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("BUNNYNET_API_KEY")
	if apiKey == "" {
		log.Fatal("BUNNYNET_API_KEY environment variable is not set")
	}

	// Create a new client with custom timeout
	client := bunnynet.NewClient(
		apiKey,
		bunnynet.WithTimeout(15*time.Second),
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List all countries
	countries, err := client.Country.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list countries: %v", err)
	}

	fmt.Printf("Found %d countries\n", len(countries))

	// Print the first 5 countries
	for i, country := range countries {
		if i >= 5 {
			break
		}
		fmt.Printf("Country: %s (%s), EU: %t, Tax Rate: %.2f%%\n",
			country.Name, country.IsoCode, country.IsEU, country.TaxRate)
	}

	// Get a specific country by ISO code
	us, err := client.Country.Get(ctx, "US")
	if err != nil {
		log.Fatalf("Failed to get country by ISO code: %v", err)
	}

	fmt.Printf("\nDetails for %s:\n", us.Name)
	fmt.Printf("  ISO Code: %s\n", us.IsoCode)
	fmt.Printf("  EU Member: %t\n", us.IsEU)
	fmt.Printf("  Tax Rate: %.2f%%\n", us.TaxRate)
	fmt.Printf("  Tax Prefix: %s\n", us.TaxPrefix)
	fmt.Printf("  POPs: %v\n", us.PopList)
}
