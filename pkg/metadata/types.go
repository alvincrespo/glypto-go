package metadata

import "golang.org/x/net/html"

// MetadataProvider defines the interface for metadata extraction providers
type MetadataProvider interface {
	// Name returns the provider's unique name
	Name() string

	// Priority returns the provider's priority (lower numbers = higher priority)
	Priority() int

	// CanHandle determines if this provider can extract data from the given element
	CanHandle(node *html.Node) bool

	// Scrape extracts metadata from the given element
	Scrape(node *html.Node) *ScrapedData

	// GetValue resolves a value for a given key from the provider's data
	GetValue(key string, data map[string][]string) *string
}

// ScrapedData represents extracted metadata from a provider
type ScrapedData struct {
	Key   string
	Value string
}

// ProviderData aggregates data from all providers
type ProviderData map[string]map[string][]string

// Feed represents an RSS/Atom feed link
type Feed struct {
	Title *string `json:"title,omitempty"`
	Type  string  `json:"type"`
	Href  string  `json:"href"`
}

// ScrapingResult represents the result of a scraping operation
type ScrapingResult struct {
	Provider *MetadataProvider
	Data     *ScrapedData
}

// Registry interface for provider management
type Registry interface {
	GetProviders() []MetadataProvider
	ScrapeFromElement(node *html.Node) *ScrapingResult
	ResolveValue(key string, providerData ProviderData) *string
	AddProvider(provider MetadataProvider)
	RemoveProvider(name string)
	GetProvider(name string) MetadataProvider
}
