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
		bunnynet.WithTimeout(30*time.Second),
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Example 1: List DNS Zones
	fmt.Println("Example 1: Listing DNS Zones")
	listDNSZones(ctx, client)

	// Example 2: Add a new DNS Zone
	fmt.Println("\nExample 2: Add a new DNS Zone")
	dnsZone := addDNSZone(ctx, client, "example-"+fmt.Sprintf("%d", time.Now().Unix())+".com")
	if dnsZone != nil {
		// Example 3: Add a DNS Record
		fmt.Println("\nExample 3: Add DNS Records")
		addDNSRecords(ctx, client, dnsZone.Id)

		// Example 4: Update DNS Zone
		fmt.Println("\nExample 4: Update DNS Zone")
		updateDNSZone(ctx, client, dnsZone.Id)

		// Example 5: Check and Enable DNSSEC
		fmt.Println("\nExample 5: Enable DNSSEC")
		enableDNSSEC(ctx, client, dnsZone.Id)

		// Example 6: Export DNS Zone
		fmt.Println("\nExample 6: Export DNS Zone")
		exportDNSZone(ctx, client, dnsZone.Id)

		// Wait a moment before cleanup
		time.Sleep(2 * time.Second)

		// Example 7: Delete DNS Zone
		fmt.Println("\nExample 7: Delete DNS Zone")
		deleteDNSZone(ctx, client, dnsZone.Id)
	}

	// Example 8: Check Zone Availability
	fmt.Println("\nExample 8: Check Zone Availability")
	checkZoneAvailability(ctx, client, "example.com")
}

func listDNSZones(ctx context.Context, client *bunnynet.Client) {
	// Create pagination options
	pagination := common.NewPagination().WithPerPage(5)

	// List DNS zones with pagination
	response, err := client.DNSZone.List(ctx, pagination, "")
	if err != nil {
		log.Printf("Failed to list DNS zones: %v", err)
		return
	}

	// Display DNS zones
	fmt.Printf("Found %d DNS zones (Page %d, Has more: %t)\n",
		response.TotalItems, response.CurrentPage, response.HasMoreItems)

	for i, zone := range response.Items {
		fmt.Printf("%d. DNS Zone: ID=%d, Domain=%s, Records=%d\n",
			i+1, zone.Id, zone.Domain, len(zone.Records))
	}

	// If there are more zones than the first page shows, demonstrate getting all zones
	if response.TotalItems > len(response.Items) {
		fmt.Println("\nFetching all zones (this may take a moment)...")
		allZones, err := client.DNSZone.ListAll(ctx, 100, "")
		if err != nil {
			log.Printf("Failed to list all DNS zones: %v", err)
			return
		}
		fmt.Printf("Successfully retrieved all %d zones\n", len(allZones))
	}
}

func addDNSZone(ctx context.Context, client *bunnynet.Client, domain string) *resources.DNSZone {
	// Check if the domain is available first
	availabilityResult, err := client.DNSZone.CheckAvailability(ctx, resources.CheckZoneAvailabilityOptions{
		Name: domain,
	})
	if err != nil {
		log.Printf("Failed to check domain availability: %v", err)
		return nil
	}

	if !availabilityResult.Available {
		log.Printf("Domain %s is not available: %s", domain, availabilityResult.Message)
		return nil
	}

	fmt.Printf("Domain %s is available. Adding...\n", domain)

	// Add the DNS zone
	options := resources.AddDNSZoneOptions{
		Domain: domain,
	}

	dnsZone, err := client.DNSZone.Add(ctx, options)
	if err != nil {
		log.Printf("Failed to add DNS zone: %v", err)
		return nil
	}

	fmt.Printf("Added DNS Zone: ID=%d, Domain=%s\n", dnsZone.Id, dnsZone.Domain)
	return dnsZone
}

func addDNSRecords(ctx context.Context, client *bunnynet.Client, zoneId int64) {
	// Add A record
	aRecordOptions := resources.AddDNSRecordOptions{
		Type:  resources.DNSRecordTypeA,
		Name:  "@",
		Value: "192.0.2.1",
		Ttl:   3600,
	}

	aRecord, err := client.DNSZone.AddRecord(ctx, zoneId, aRecordOptions)
	if err != nil {
		log.Printf("Failed to add A record: %v", err)
		return
	}
	fmt.Printf("Added A record: ID=%d, Name=%s, Value=%s\n",
		aRecord.Id, aRecord.Name, aRecord.Value)

	// Add CNAME record
	cnameRecordOptions := resources.AddDNSRecordOptions{
		Type:  resources.DNSRecordTypeCNAME,
		Name:  "www",
		Value: "@",
		Ttl:   3600,
	}

	cnameRecord, err := client.DNSZone.AddRecord(ctx, zoneId, cnameRecordOptions)
	if err != nil {
		log.Printf("Failed to add CNAME record: %v", err)
		return
	}
	fmt.Printf("Added CNAME record: ID=%d, Name=%s, Value=%s\n",
		cnameRecord.Id, cnameRecord.Name, cnameRecord.Value)

	// Add MX record
	mxRecordOptions := resources.AddDNSRecordOptions{
		Type:     resources.DNSRecordTypeMX,
		Name:     "@",
		Value:    "mail.example.com",
		Ttl:      3600,
		Priority: 10,
	}

	mxRecord, err := client.DNSZone.AddRecord(ctx, zoneId, mxRecordOptions)
	if err != nil {
		log.Printf("Failed to add MX record: %v", err)
		return
	}
	fmt.Printf("Added MX record: ID=%d, Name=%s, Value=%s, Priority=%d\n",
		mxRecord.Id, mxRecord.Name, mxRecord.Value, mxRecord.Priority)

	// Update the CNAME record
	updateRecordOptions := resources.UpdateDNSRecordOptions{
		Id:    cnameRecord.Id,
		Type:  resources.DNSRecordTypeCNAME,
		Name:  "blog",
		Value: "@",
		Ttl:   7200,
	}

	err = client.DNSZone.UpdateRecord(ctx, zoneId, cnameRecord.Id, updateRecordOptions)
	if err != nil {
		log.Printf("Failed to update CNAME record: %v", err)
		return
	}
	fmt.Printf("Updated CNAME record (ID=%d) to Name=%s, TTL=%d\n",
		cnameRecord.Id, updateRecordOptions.Name, updateRecordOptions.Ttl)

	// Delete the MX record
	err = client.DNSZone.DeleteRecord(ctx, zoneId, mxRecord.Id)
	if err != nil {
		log.Printf("Failed to delete MX record: %v", err)
		return
	}
	fmt.Printf("Deleted MX record (ID=%d)\n", mxRecord.Id)
}

func updateDNSZone(ctx context.Context, client *bunnynet.Client, zoneId int64) {
	// Update the DNS zone settings
	options := resources.UpdateDNSZoneOptions{
		CustomNameserversEnabled:      true,
		Nameserver1:                   "ns1.example.com",
		Nameserver2:                   "ns2.example.com",
		SoaEmail:                      "admin@example.com",
		LoggingEnabled:                true,
		LogAnonymizationType:          resources.LogAnonymizationTypeOneDigit,
		LoggingIPAnonymizationEnabled: true,
	}

	updatedZone, err := client.DNSZone.Update(ctx, zoneId, options)
	if err != nil {
		log.Printf("Failed to update DNS zone: %v", err)
		return
	}

	fmt.Printf("Updated DNS Zone: ID=%d, Domain=%s\n", updatedZone.Id, updatedZone.Domain)
	fmt.Printf("  Custom Nameservers: %t\n", updatedZone.CustomNameserversEnabled)
	fmt.Printf("  Nameserver1: %s\n", updatedZone.Nameserver1)
	fmt.Printf("  Nameserver2: %s\n", updatedZone.Nameserver2)
	fmt.Printf("  SOA Email: %s\n", updatedZone.SoaEmail)
	fmt.Printf("  Logging Enabled: %t\n", updatedZone.LoggingEnabled)
	fmt.Printf("  IP Anonymization: %t\n", updatedZone.LoggingIPAnonymizationEnabled)
}

func enableDNSSEC(ctx context.Context, client *bunnynet.Client, zoneId int64) {
	// Get the current DNS zone to check DNSSEC status
	zone, err := client.DNSZone.Get(ctx, zoneId)
	if err != nil {
		log.Printf("Failed to get DNS zone: %v", err)
		return
	}

	// Enable DNSSEC if not already enabled
	if !zone.DnsSecEnabled {
		fmt.Println("Enabling DNSSEC...")
		dnsSecInfo, err := client.DNSZone.EnableDNSSec(ctx, zoneId)
		if err != nil {
			log.Printf("Failed to enable DNSSEC: %v", err)
			return
		}

		fmt.Printf("DNSSEC Enabled: %t\n", dnsSecInfo.Enabled)
		fmt.Printf("  DS Record: %s\n", dnsSecInfo.DsRecord)
		fmt.Printf("  Algorithm: %d\n", dnsSecInfo.Algorithm)
		fmt.Printf("  Key Tag: %d\n", dnsSecInfo.KeyTag)
	} else {
		fmt.Println("DNSSEC is already enabled. Disabling...")

		// Disable DNSSEC
		dnsSecInfo, err := client.DNSZone.DisableDNSSec(ctx, zoneId)
		if err != nil {
			log.Printf("Failed to disable DNSSEC: %v", err)
			return
		}

		fmt.Printf("DNSSEC Disabled: %t\n", !dnsSecInfo.Enabled)
	}
}

func exportDNSZone(ctx context.Context, client *bunnynet.Client, zoneId int64) {
	// Export the DNS zone
	zoneData, err := client.DNSZone.Export(ctx, zoneId)
	if err != nil {
		log.Printf("Failed to export DNS zone: %v", err)
		return
	}

	// Display the zone file (first 500 characters)
	zoneString := string(zoneData)
	if len(zoneString) > 500 {
		zoneString = zoneString[:500] + "...[truncated]"
	}

	fmt.Printf("Exported Zone File (%d bytes):\n%s\n", len(zoneData), zoneString)

	// Example of importing the zone file back
	fmt.Println("Importing zone file back...")
	importResult, err := client.DNSZone.ImportRecords(ctx, zoneId, zoneData)
	if err != nil {
		log.Printf("Failed to import DNS records: %v", err)
		return
	}

	fmt.Printf("Import Result:\n")
	fmt.Printf("  Records Successful: %d\n", importResult.RecordsSuccessful)
	fmt.Printf("  Records Failed: %d\n", importResult.RecordsFailed)
	fmt.Printf("  Records Skipped: %d\n", importResult.RecordsSkipped)
}

func deleteDNSZone(ctx context.Context, client *bunnynet.Client, zoneId int64) {
	// Delete the DNS zone
	err := client.DNSZone.Delete(ctx, zoneId)
	if err != nil {
		log.Printf("Failed to delete DNS zone: %v", err)
		return
	}

	fmt.Printf("Deleted DNS Zone ID=%d\n", zoneId)
}

func checkZoneAvailability(ctx context.Context, client *bunnynet.Client, domain string) {
	// Check if a domain is available for adding as a DNS zone
	result, err := client.DNSZone.CheckAvailability(ctx, resources.CheckZoneAvailabilityOptions{
		Name: domain,
	})
	if err != nil {
		log.Printf("Failed to check zone availability: %v", err)
		return
	}

	fmt.Printf("Domain %s availability: %t\n", domain, result.Available)
	if !result.Available && result.Message != "" {
		fmt.Printf("  Message: %s\n", result.Message)
	}
}
