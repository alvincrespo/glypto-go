package providers

import (
	"testing"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

// MockProvider for testing
type MockProvider struct {
	name     string
	priority int
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Priority() int {
	return m.priority
}

func (m *MockProvider) CanHandle(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "meta"
}

func (m *MockProvider) Scrape(node *html.Node) *metadata.ScrapedData {
	return &metadata.ScrapedData{
		Key:   "test",
		Value: "value",
	}
}

func (m *MockProvider) GetValue(key string, data map[string][]string) *string {
	if values, exists := data[key]; exists && len(values) > 0 {
		return &values[0]
	}
	return nil
}

func TestNewRegistry(t *testing.T) {
	provider1 := &MockProvider{name: "provider1", priority: 3}
	provider2 := &MockProvider{name: "provider2", priority: 1}
	provider3 := &MockProvider{name: "provider3", priority: 2}

	providers := []metadata.MetadataProvider{provider1, provider2, provider3}
	registry := NewRegistry(providers)

	if len(registry.providers) != 3 {
		t.Errorf("Expected 3 providers, got %d", len(registry.providers))
	}

	// Check that providers are sorted by priority
	priorities := []int{1, 2, 3}
	for i, provider := range registry.providers {
		if provider.Priority() != priorities[i] {
			t.Errorf("Expected priority %d at index %d, got %d", priorities[i], i, provider.Priority())
		}
	}
}

func TestProviderRegistry_GetProviders(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	providers := registry.GetProviders()
	if len(providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(providers))
	}

	if providers[0].Name() != "test" {
		t.Errorf("Expected provider name 'test', got '%s'", providers[0].Name())
	}
}

func TestProviderRegistry_ScrapeFromElement(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	// Create a test HTML node
	node := &html.Node{
		Type: html.ElementNode,
		Data: "meta",
	}

	result := registry.ScrapeFromElement(node)
	if result == nil {
		t.Error("Expected scraping result, got nil")
		return
	}

	if result.Data.Key != "test" {
		t.Errorf("Expected key 'test', got '%s'", result.Data.Key)
	}

	if result.Data.Value != "value" {
		t.Errorf("Expected value 'value', got '%s'", result.Data.Value)
	}
}

func TestProviderRegistry_ScrapeFromElement_NoHandler(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	// Create a node that can't be handled
	node := &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}

	result := registry.ScrapeFromElement(node)
	if result != nil {
		t.Error("Expected nil result for unhandled element")
	}
}

func TestProviderRegistry_ResolveValue(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	providerData := metadata.ProviderData{
		"test": map[string][]string{
			"title": {"Test Title"},
		},
	}

	result := registry.ResolveValue("title", providerData)
	if result == nil {
		t.Error("Expected value, got nil")
		return
	}

	if *result != "Test Title" {
		t.Errorf("Expected 'Test Title', got '%s'", *result)
	}
}

func TestProviderRegistry_ResolveValue_NotFound(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	providerData := metadata.ProviderData{}

	result := registry.ResolveValue("title", providerData)
	if result != nil {
		t.Error("Expected nil for non-existent value")
	}
}

func TestProviderRegistry_AddProvider(t *testing.T) {
	provider1 := &MockProvider{name: "provider1", priority: 2}
	registry := NewRegistry([]metadata.MetadataProvider{provider1})

	provider2 := &MockProvider{name: "provider2", priority: 1}
	registry.AddProvider(provider2)

	if len(registry.providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(registry.providers))
	}

	// Check that providers are re-sorted by priority
	if registry.providers[0].Priority() != 1 {
		t.Errorf("Expected first provider to have priority 1, got %d", registry.providers[0].Priority())
	}
}

func TestProviderRegistry_RemoveProvider(t *testing.T) {
	provider1 := &MockProvider{name: "provider1", priority: 1}
	provider2 := &MockProvider{name: "provider2", priority: 2}
	registry := NewRegistry([]metadata.MetadataProvider{provider1, provider2})

	registry.RemoveProvider("provider1")

	if len(registry.providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(registry.providers))
	}

	if registry.providers[0].Name() != "provider2" {
		t.Errorf("Expected remaining provider to be 'provider2', got '%s'", registry.providers[0].Name())
	}
}

func TestProviderRegistry_RemoveProvider_NotFound(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	registry.RemoveProvider("nonexistent")

	if len(registry.providers) != 1 {
		t.Error("Provider should not be removed when name doesn't match")
	}
}

func TestProviderRegistry_GetProvider(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	result := registry.GetProvider("test")
	if result == nil {
		t.Error("Expected provider, got nil")
	}

	if result.Name() != "test" {
		t.Errorf("Expected provider name 'test', got '%s'", result.Name())
	}
}

func TestProviderRegistry_GetProvider_NotFound(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1}
	registry := NewRegistry([]metadata.MetadataProvider{provider})

	result := registry.GetProvider("nonexistent")
	if result != nil {
		t.Error("Expected nil for non-existent provider")
	}
}
