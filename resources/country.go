package resources

import (
	"context"
	"net/http"

	"github.com/venom90/bunnynet-go/common"
	"github.com/venom90/bunnynet-go/internal"
)

// Country represents a country in the Bunny.net API
type Country struct {
	// Name is the name of the country
	Name string `json:"Name"`

	// IsoCode is the ISO 3166-1 alpha-2 code of the country
	IsoCode string `json:"IsoCode"`

	// IsEU indicates whether the country is in the European Union
	IsEU bool `json:"IsEU"`

	// TaxRate is the tax rate of the country
	TaxRate float64 `json:"TaxRate"`

	// TaxPrefix is the tax prefix of the country
	TaxPrefix string `json:"TaxPrefix"`

	// FlagUrl is the URL of the country's flag
	FlagUrl string `json:"FlagUrl"`

	// PopList is a list of POPs in the country
	PopList []string `json:"PopList"`
}

// CountryService handles operations on countries
type CountryService struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	userAgent string
}

// NewCountryService creates a new CountryService
func NewCountryService(client *http.Client, baseURL, apiKey, userAgent string) *CountryService {
	return &CountryService{
		client:    client,
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: userAgent,
	}
}

// SetAPIKey updates the API key used for authentication
func (s *CountryService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// List returns a list of all countries
func (s *CountryService) List(ctx context.Context) ([]Country, error) {
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, "/country", nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var countries []Country
	if err := internal.ParseResponse(resp, &countries); err != nil {
		return nil, err
	}

	return countries, nil
}

// ListPaginated returns a paginated list of countries
func (s *CountryService) ListPaginated(ctx context.Context, pagination *common.Pagination) (*common.PaginatedResponse[Country], error) {
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, "/country", nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add pagination parameters
	if err := internal.AddQueryParams(req, pagination); err != nil {
		return nil, err
	}

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}

	var paginatedResponse common.PaginatedResponse[Country]
	if err := internal.ParsePaginatedResponse(resp, &paginatedResponse); err != nil {
		return nil, err
	}

	return &paginatedResponse, nil
}

// ListAll returns all countries across all pages
func (s *CountryService) ListAll(ctx context.Context, perPage int) ([]Country, error) {
	if perPage <= 0 {
		perPage = common.DefaultPerPage
	}

	iterator := common.NewPageIterator(
		func(page, itemsPerPage int) (*common.PaginatedResponse[Country], error) {
			pagination := common.NewPagination().WithPage(page).WithPerPage(itemsPerPage)
			return s.ListPaginated(ctx, pagination)
		},
		common.DefaultPage,
		perPage,
	)

	return iterator.AllItems()
}

// Get returns a country by ISO code
func (s *CountryService) Get(ctx context.Context, isoCode string) (*Country, error) {
	path := "/country/" + isoCode
	req, err := internal.NewRequest(http.MethodGet, s.baseURL, path, nil, s.apiKey, s.userAgent)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := internal.DoRequest(s.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var country Country
	if err := internal.ParseResponse(resp, &country); err != nil {
		return nil, err
	}

	return &country, nil
}
