package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/venom90/bunnynet-go"
	"github.com/venom90/bunnynet-go/common"
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

	// Example 1: List Pull Zones
	fmt.Println("Example 1: Listing Pull Zones")
	listPullZones(ctx, client)

	// Example 2: Create a new Pull Zone
	fmt.Println("\nExample 2: Creating a new Pull Zone")
	newPullZone := createPullZone(ctx, client)
	if newPullZone != nil {
		// Example 3: Get a Pull Zone by ID
		fmt.Println("\nExample 3: Getting a Pull Zone by ID")
		getPullZone(ctx, client, newPullZone.Id)

		// Example 4: Add a hostname to the Pull Zone
		fmt.Println("\nExample 4: Adding a hostname to the Pull Zone")
		addHostname(ctx, client, newPullZone.Id, "cdn.example.com")

		// Example 5: Add an Edge Rule to the Pull Zone
		fmt.Println("\nExample 5: Adding an Edge Rule to the Pull Zone")
		addEdgeRule(ctx, client, newPullZone.Id)

		// Example 6: Purge cache for the Pull Zone
		fmt.Println("\nExample 6: Purging cache for the Pull Zone")
		purgeCache(ctx, client, newPullZone.Id)

		// Example 7: Update the Pull Zone
		fmt.Println("\nExample 7: Updating the Pull Zone")
		updatePullZone(ctx, client, newPullZone.Id)

		// Example 8: Delete the Pull Zone (uncomment to actually delete)
		fmt.Println("\nExample 8: Deleting the Pull Zone")
		deletePullZone(ctx, client, newPullZone.Id)
	}
}

func listPullZones(ctx context.Context, client *bunnynet.Client) {
	// Create pagination options
	pagination := common.NewPagination().WithPerPage(10)

	// List all Pull Zones
	response, err := client.PullZone.List(ctx, pagination, "", false)
	if err != nil {
		log.Fatalf("Failed to list Pull Zones: %v", err)
	}

	// Display Pull Zones
	fmt.Printf("Found %d Pull Zones (Page %d, Has more: %t)\n",
		response.TotalItems, response.CurrentPage, response.HasMoreItems)

	for i, pullZone := range response.Items {
		fmt.Printf("%d. Pull Zone: ID=%d, Name=%s, Origin=%s, Enabled=%t, Hostnames=%d\n",
			i+1, pullZone.Id, pullZone.Name, pullZone.OriginUrl, pullZone.Enabled, len(pullZone.Hostnames))
	}
}

func createPullZone(ctx context.Context, client *bunnynet.Client) *resources.PullZone {
	// Define options for the new Pull Zone
	options := resources.AddPullZoneOptions{
		Name:              "example-pull-zone",
		OriginUrl:         "https://example.com",
		Type:              0, // Premium
		EnableGeoZoneUS:   true,
		EnableGeoZoneEU:   true,
		EnableGeoZoneASIA: true,
		EnableGeoZoneSA:   true,
		EnableGeoZoneAF:   true,
	}

	fmt.Printf("Creating a new Pull Zone with name: %s, origin: %s\n", options.Name, options.OriginUrl)

	// Create a new Pull Zone
	pullZone, err := client.PullZone.Add(ctx, options)
	if err != nil {
		log.Printf("Failed to create Pull Zone: %v", err)
		return nil
	}

	fmt.Printf("Created new Pull Zone:\n")
	fmt.Printf("  ID: %d\n", pullZone.Id)
	fmt.Printf("  Name: %s\n", pullZone.Name)
	fmt.Printf("  Origin URL: %s\n", pullZone.OriginUrl)
	fmt.Printf("  Enabled: %t\n", pullZone.Enabled)
	fmt.Printf("  Default Hostname: %s\n", pullZone.Hostnames[0].Value)

	return pullZone
}

func getPullZone(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Getting Pull Zone with ID: %d\n", id)

	// Get Pull Zone by ID
	pullZone, err := client.PullZone.Get(ctx, id, false)
	if err != nil {
		log.Printf("Failed to get Pull Zone: %v", err)
		return
	}

	fmt.Printf("Pull Zone details:\n")
	fmt.Printf("  ID: %d\n", pullZone.Id)
	fmt.Printf("  Name: %s\n", pullZone.Name)
	fmt.Printf("  Origin URL: %s\n", pullZone.OriginUrl)
	fmt.Printf("  Enabled: %t\n", pullZone.Enabled)
	fmt.Printf("  Monthly Bandwidth Used: %d bytes\n", pullZone.MonthlyBandwidthUsed)
	fmt.Printf("  Monthly Bandwidth Limit: %d bytes\n", pullZone.MonthlyBandwidthLimit)
	fmt.Printf("  Hostnames: %d\n", len(pullZone.Hostnames))

	for i, hostname := range pullZone.Hostnames {
		fmt.Printf("    %d. %s (Force SSL: %t, System: %t)\n",
			i+1, hostname.Value, hostname.ForceSSL, hostname.IsSystemHostname)
	}

	fmt.Printf("  Edge Rules: %d\n", len(pullZone.EdgeRules))
	for i, rule := range pullZone.EdgeRules {
		fmt.Printf("    %d. %s (Type: %d, Enabled: %t)\n",
			i+1, rule.Description, rule.ActionType, rule.Enabled)
	}
}

func addHostname(ctx context.Context, client *bunnynet.Client, id int64, hostname string) {
	fmt.Printf("Adding hostname %s to Pull Zone ID: %d\n", hostname, id)

	// Add hostname
	err := client.PullZone.AddHostname(ctx, id, resources.AddHostnameOptions{
		Hostname: hostname,
	})
	if err != nil {
		log.Printf("Failed to add hostname: %v", err)
		return
	}

	fmt.Printf("Successfully added hostname %s to Pull Zone ID: %d\n", hostname, id)
}

func addEdgeRule(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Adding Edge Rule to Pull Zone ID: %d\n", id)

	// Create Edge Rule options
	options := resources.AddOrUpdateEdgeRuleOptions{
		ActionType: 0, // ForceSSL
		Triggers: []resources.EdgeRuleTrigger{
			{
				Type:                0, // URL
				PatternMatches:      []string{"/*"},
				PatternMatchingType: 0, // MatchAny
				TriggerMatchingType: 0, // MatchAny
			},
		},
		Description: "Force SSL for all URLs",
		Enabled:     true,
	}

	// Add Edge Rule
	err := client.PullZone.AddOrUpdateEdgeRule(ctx, id, options)
	if err != nil {
		log.Printf("Failed to add Edge Rule: %v", err)
		return
	}

	fmt.Printf("Successfully added Edge Rule to Pull Zone ID: %d\n", id)
}

func purgeCache(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Purging cache for Pull Zone ID: %d\n", id)

	// Purge cache
	err := client.PullZone.PurgeCache(ctx, id, nil)
	if err != nil {
		log.Printf("Failed to purge cache: %v", err)
		return
	}

	fmt.Printf("Successfully purged cache for Pull Zone ID: %d\n", id)
}

func updatePullZone(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Updating Pull Zone ID: %d\n", id)

	// Create update options
	options := &resources.PullZone{
		EnableHostnameVary:         true,
		EnableQueryStringOrdering:  true,
		CacheControlMaxAgeOverride: 86400, // 1 day
	}

	// Update Pull Zone
	pullZone, err := client.PullZone.Update(ctx, id, options)
	if err != nil {
		log.Printf("Failed to update Pull Zone: %v", err)
		return
	}

	fmt.Printf("Successfully updated Pull Zone ID: %d\n", id)
	fmt.Printf("  EnableHostnameVary: %t\n", pullZone.EnableHostnameVary)
	fmt.Printf("  EnableQueryStringOrdering: %t\n", pullZone.EnableQueryStringOrdering)
	fmt.Printf("  CacheControlMaxAgeOverride: %d seconds\n", pullZone.CacheControlMaxAgeOverride)
}

func deletePullZone(ctx context.Context, client *bunnynet.Client, id int64) {
	fmt.Printf("Deleting Pull Zone ID: %d\n", id)

	// Delete Pull Zone
	err := client.PullZone.Delete(ctx, id)
	if err != nil {
		log.Printf("Failed to delete Pull Zone: %v", err)
		return
	}

	fmt.Printf("Successfully deleted Pull Zone ID: %d\n", id)
}
