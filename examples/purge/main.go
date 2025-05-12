package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/venom90/bunnynet-go"
	"github.com/venom90/bunnynet-go/resources"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("BUNNYNET_API_KEY")
	if apiKey == "" {
		log.Fatal("BUNNYNET_API_KEY environment variable is not set")
	}

	// Create a new client
	client := bunnynet.NewClient(
		apiKey,
		bunnynet.WithTimeout(15*time.Second),
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example 1: Purge a URL using the simplified method
	fmt.Println("Example 1: Purging a URL (synchronously)")
	purgeURL(ctx, client, "https://example.com/file.jpg", false)

	// Example 2: Purge a URL asynchronously
	fmt.Println("\nExample 2: Purging a URL asynchronously")
	purgeURLAsync(ctx, client, "https://example.com/folder/", true)

	// Example 3: Purge a URL using PurgeOptions
	fmt.Println("\nExample 3: Purging a URL with PurgeOptions")
	purgeURLWithOptions(ctx, client, resources.PurgeOptions{
		URL:   "https://example.com/api/data.json",
		Async: true,
	})
}

func purgeURL(ctx context.Context, client *bunnynet.Client, url string, async bool) {
	fmt.Printf("Purging URL: %s (async: %t)\n", url, async)

	// Purge the URL
	err := client.Purge.Purge(ctx, url, async)
	if err != nil {
		log.Printf("Failed to purge URL: %v", err)
		return
	}

	fmt.Printf("Successfully purged URL: %s\n", url)
}

func purgeURLAsync(ctx context.Context, client *bunnynet.Client, url string, async bool) {
	fmt.Printf("Purging URL asynchronously: %s\n", url)

	// Purge the URL
	err := client.Purge.Purge(ctx, url, async)
	if err != nil {
		log.Printf("Failed to purge URL: %v", err)
		return
	}

	fmt.Printf("Successfully initiated asynchronous purge for URL: %s\n", url)
}

func purgeURLWithOptions(ctx context.Context, client *bunnynet.Client, options resources.PurgeOptions) {
	fmt.Printf("Purging URL with options: %s (async: %t)\n", options.URL, options.Async)

	// Purge the URL with options
	err := client.Purge.PurgeURL(ctx, options)
	if err != nil {
		log.Printf("Failed to purge URL: %v", err)
		return
	}

	fmt.Printf("Successfully purged URL: %s\n", options.URL)
}
