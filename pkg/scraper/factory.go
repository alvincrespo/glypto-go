package scraper

import (
	"golang.org/x/net/html"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"github.com/alvincrespo/glypto-go/pkg/providers"
)

// CreateScraper creates a scraper with auto-loaded providers
func CreateScraper() (*Scraper, error) {
	loader := providers.NewLoader()

	// Try to load from directory first, fallback to defaults
	providerList, err := loader.LoadFromDirectory("")
	if err != nil {
		// If loading from directory fails, use defaults
		providerList = loader.LoadDefaults()
	}

	registry := providers.NewRegistry(providerList)
	return NewScraper(registry), nil
}

// CreateScraperWithProviders creates a scraper with custom providers
func CreateScraperWithProviders(providerList []metadata.MetadataProvider) *Scraper {
	registry := providers.NewRegistry(providerList)
	return NewScraper(registry)
}

// CreateScraperWithProviderNames creates a scraper with specific provider names
func CreateScraperWithProviderNames(providerNames []string) (*Scraper, error) {
	loader := providers.NewLoader()

	providerList, err := loader.LoadFromList(providerNames)
	if err != nil {
		return nil, err
	}

	registry := providers.NewRegistry(providerList)
	return NewScraper(registry), nil
}

// ScrapeMetadata is a convenience function to scrape metadata from a document
func ScrapeMetadata(doc *html.Node) (*metadata.Metadata, error) {
	scraper, err := CreateScraper()
	if err != nil {
		return nil, err
	}

	return scraper.Scrape(doc)
}

// ScrapeMetadataWithProviders is a convenience function to scrape with custom providers
func ScrapeMetadataWithProviders(doc *html.Node, providerList []metadata.MetadataProvider) (*metadata.Metadata, error) {
	scraper := CreateScraperWithProviders(providerList)
	return scraper.Scrape(doc)
}

// ScrapeMetadataWithProviderNames is a convenience function to scrape with specific provider names
func ScrapeMetadataWithProviderNames(doc *html.Node, providerNames []string) (*metadata.Metadata, error) {
	scraper, err := CreateScraperWithProviderNames(providerNames)
	if err != nil {
		return nil, err
	}

	return scraper.Scrape(doc)
}
