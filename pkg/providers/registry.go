package providers

import (
	"sort"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

// ProviderRegistry manages metadata providers with priority-based resolution
type ProviderRegistry struct {
	providers []metadata.MetadataProvider
}

// NewRegistry creates a new provider registry
func NewRegistry(providers []metadata.MetadataProvider) *ProviderRegistry {
	// Sort providers by priority (lower numbers = higher priority)
	sortedProviders := make([]metadata.MetadataProvider, len(providers))
	copy(sortedProviders, providers)

	sort.Slice(sortedProviders, func(i, j int) bool {
		return sortedProviders[i].Priority() < sortedProviders[j].Priority()
	})

	return &ProviderRegistry{
		providers: sortedProviders,
	}
}

// GetProviders returns all registered providers
func (r *ProviderRegistry) GetProviders() []metadata.MetadataProvider {
	return r.providers
}

// ScrapeFromElement attempts to scrape metadata from an element using all providers
func (r *ProviderRegistry) ScrapeFromElement(node *html.Node) *metadata.ScrapingResult {
	for _, provider := range r.providers {
		if provider.CanHandle(node) {
			if data := provider.Scrape(node); data != nil {
				return &metadata.ScrapingResult{
					Provider: &provider,
					Data:     data,
				}
			}
		}
	}
	return nil
}

// ResolveValue resolves a value using provider priority
func (r *ProviderRegistry) ResolveValue(key string, providerData metadata.ProviderData) *string {
	for _, provider := range r.providers {
		if data, exists := providerData[provider.Name()]; exists {
			if value := provider.GetValue(key, data); value != nil {
				return value
			}
		}
	}
	return nil
}

// AddProvider adds a new provider to the registry
func (r *ProviderRegistry) AddProvider(provider metadata.MetadataProvider) {
	r.providers = append(r.providers, provider)

	// Re-sort providers by priority
	sort.Slice(r.providers, func(i, j int) bool {
		return r.providers[i].Priority() < r.providers[j].Priority()
	})
}

// RemoveProvider removes a provider from the registry by name
func (r *ProviderRegistry) RemoveProvider(name string) {
	for i, provider := range r.providers {
		if provider.Name() == name {
			r.providers = append(r.providers[:i], r.providers[i+1:]...)
			return
		}
	}
}

// GetProvider returns a provider by name
func (r *ProviderRegistry) GetProvider(name string) metadata.MetadataProvider {
	for _, provider := range r.providers {
		if provider.Name() == name {
			return provider
		}
	}
	return nil
}
