package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venom90/bunnynet-go-client"
	"github.com/venom90/bunnynet-go-client/test"
)

func TestCountryList(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `[
		{
			"Name": "United States",
			"IsoCode": "US",
			"IsEU": false,
			"TaxRate": 0,
			"TaxPrefix": "",
			"FlagUrl": "https://bunnycdn.com/flags/us.png",
			"PopList": ["NY", "LA", "CHI"]
		},
		{
			"Name": "Germany",
			"IsoCode": "DE",
			"IsEU": true,
			"TaxRate": 19,
			"TaxPrefix": "DE123456789",
			"FlagUrl": "https://bunnycdn.com/flags/de.png",
			"PopList": ["FRA", "BER"]
		}
	]`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/country")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the List method
	countries, err := client.Country.List(context.Background())
	assert.NoError(t, err, "List should not return an error")
	assert.Len(t, countries, 2, "Should return 2 countries")

	// Check the first country
	assert.Equal(t, "United States", countries[0].Name)
	assert.Equal(t, "US", countries[0].IsoCode)
	assert.False(t, countries[0].IsEU)
	assert.Equal(t, 0.0, countries[0].TaxRate)
	assert.Empty(t, countries[0].TaxPrefix)
	assert.Equal(t, "https://bunnycdn.com/flags/us.png", countries[0].FlagUrl)
	assert.Equal(t, []string{"NY", "LA", "CHI"}, countries[0].PopList)

	// Check the second country
	assert.Equal(t, "Germany", countries[1].Name)
	assert.Equal(t, "DE", countries[1].IsoCode)
	assert.True(t, countries[1].IsEU)
	assert.Equal(t, 19.0, countries[1].TaxRate)
	assert.Equal(t, "DE123456789", countries[1].TaxPrefix)
	assert.Equal(t, "https://bunnycdn.com/flags/de.png", countries[1].FlagUrl)
	assert.Equal(t, []string{"FRA", "BER"}, countries[1].PopList)
}

func TestCountryGet(t *testing.T) {
	// Create a mock server
	server := test.MockServer(t, http.StatusOK, `{
		"Name": "United States",
		"IsoCode": "US",
		"IsEU": false,
		"TaxRate": 0,
		"TaxPrefix": "",
		"FlagUrl": "https://bunnycdn.com/flags/us.png",
		"PopList": ["NY", "LA", "CHI"]
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/country/US")
		test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method
	country, err := client.Country.Get(context.Background(), "US")
	assert.NoError(t, err, "Get should not return an error")
	assert.NotNil(t, country, "Country should not be nil")

	// Check the country
	assert.Equal(t, "United States", country.Name)
	assert.Equal(t, "US", country.IsoCode)
	assert.False(t, country.IsEU)
	assert.Equal(t, 0.0, country.TaxRate)
	assert.Empty(t, country.TaxPrefix)
	assert.Equal(t, "https://bunnycdn.com/flags/us.png", country.FlagUrl)
	assert.Equal(t, []string{"NY", "LA", "CHI"}, country.PopList)
}

func TestCountryError(t *testing.T) {
	// Create a mock server that returns an error
	server := test.MockServer(t, http.StatusNotFound, `{
		"ErrorKey": "country.not_found",
		"Field": "Country",
		"Message": "The requested country was not found"
	}`, func(r *http.Request) {
		test.AssertRequestMethod(t, r, http.MethodGet)
		test.AssertRequestPath(t, r, "/country/XX")
	})
	defer server.Close()

	// Create a client that uses the mock server
	client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

	// Call the Get method with an invalid ISO code
	country, err := client.Country.Get(context.Background(), "XX")
	assert.Error(t, err, "Get should return an error")
	assert.Nil(t, country, "Country should be nil")
	assert.Contains(t, err.Error(), "country.not_found")
	assert.Contains(t, err.Error(), "The requested country was not found")
}
