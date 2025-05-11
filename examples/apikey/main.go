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
	client := bunnynet.NewClient(
		apiKey,
		bunnynet.WithTimeout(15*time.Second),
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example 1: List API keys
	fmt.Println("Example 1: Listing API keys")
	listAPIKeys(ctx, client)

	// Example 2: Create a new API key
	fmt.Println("\nExample 2: Creating a new API key")
	newAPIKey := createAPIKey(ctx, client)
	if newAPIKey != nil {
		// Example 3: Get an API key by ID
		fmt.Println("\nExample 3: Getting an API key by ID")
		getAPIKey(ctx, client, newAPIKey.Id)

		// Example 4: Delete an API key
		fmt.Println("\nExample 4: Deleting an API key")
		deleteAPIKey(ctx, client, newAPIKey.Id)
	}
}

func listAPIKeys(ctx context.Context, client *bunnynet.Client) {
	// Create pagination options
	pagination := common.NewPagination().WithPerPage(10)

	// List all API keys
	response, err := client.APIKey.List(ctx, pagination)
	if err != nil {
		log.Fatalf("Failed to list API keys: %v", err)
	}

	// Display API keys
	fmt.Printf("Found %d API keys (Page %d, Has more: %t)\n",
		response.TotalItems, response.CurrentPage, response.HasMoreItems)

	for i, apiKey := range response.Items {
		fmt.Printf("%d. API Key: ID=%d, Key=%s, Roles=%v\n",
			i+1, apiKey.Id, maskAPIKey(apiKey.Key), apiKey.Roles)
	}
}

func createAPIKey(ctx context.Context, client *bunnynet.Client) *resources.APIKey {
	// Define roles for the new API key
	roles := []string{"PullZone.Read", "Statistics.Read"}
	fmt.Printf("Creating a new API key with roles: %v\n", roles)

	// Create a new API key
	apiKey, err := client.APIKey.Create(ctx, roles)
	if err != nil {
		log.Printf("Failed to create API key: %v", err)
		return nil
	}

	fmt.Printf("Created new API key:\n")
	fmt.Printf("  ID: %d\n", apiKey.Id)
	fmt.Printf("  Key: %s\n", apiKey.Key)
	fmt.Printf("  Roles: %v\n", apiKey.Roles)

	return apiKey
}

func getAPIKey(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Getting API key with ID: %d\n", id)

	// Get API key by ID
	apiKey, err := client.APIKey.Get(ctx, id)
	if err != nil {
		log.Printf("Failed to get API key: %v", err)
		return
	}

	fmt.Printf("API Key details:\n")
	fmt.Printf("  ID: %d\n", apiKey.Id)
	fmt.Printf("  Key: %s\n", maskAPIKey(apiKey.Key))
	fmt.Printf("  Roles: %v\n", apiKey.Roles)
}

func deleteAPIKey(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Deleting API key with ID: %d\n", id)

	// Delete API key
	err := client.APIKey.Delete(ctx, id)
	if err != nil {
		log.Printf("Failed to delete API key: %v", err)
		return
	}

	fmt.Printf("Successfully deleted API key with ID: %d\n", id)
}

// maskAPIKey masks an API key for display purposes
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return key
	}
	return key[:4] + "..." + key[len(key)-4:]
}
