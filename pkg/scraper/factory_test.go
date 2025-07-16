package scraper

import (
	"testing"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

func TestCreateScraper(t *testing.T) {
	scraper, err := CreateScraper()

	if err != nil {
		t.Errorf("CreateScraper() returned error: %v", err)
	}

	if scraper == nil {
		t.Error("CreateScraper() returned nil scraper")
		return
	}

	if scraper.registry == nil {
		t.Error("CreateScraper() created scraper without registry")
	}
}

func TestCreateScraperWithProviders(t *testing.T) {
	mockProvider := &MockProvider{name: "test", priority: 1, element: "meta"}
	providers := []metadata.MetadataProvider{mockProvider}

	scraper := CreateScraperWithProviders(providers)

	if scraper == nil {
		t.Error("CreateScraperWithProviders() returned nil scraper")
		return
	}

	if scraper.registry == nil {
		t.Error("CreateScraperWithProviders() created scraper without registry")
		return
	}

	registryProviders := scraper.registry.GetProviders()
	if len(registryProviders) != 1 {
		t.Errorf("Expected 1 provider in registry, got %d", len(registryProviders))
	}

	if registryProviders[0].Name() != "test" {
		t.Errorf("Expected provider name 'test', got '%s'", registryProviders[0].Name())
	}
}

func TestCreateScraperWithProviderNames(t *testing.T) {
	tests := []struct {
		name          string
		providerNames []string
		expectError   bool
		expectedCount int
	}{
		{
			name:          "valid provider names",
			providerNames: []string{"openGraph", "twitter"},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "single valid provider",
			providerNames: []string{"meta"},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:          "invalid provider name",
			providerNames: []string{"nonexistent"},
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "mixed valid and invalid",
			providerNames: []string{"openGraph", "invalid"},
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "empty list",
			providerNames: []string{},
			expectError:   false,
			expectedCount: 4, // Should return defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := CreateScraperWithProviderNames(tt.providerNames)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if scraper != nil {
					t.Error("Expected nil scraper when error occurs")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if scraper == nil {
				t.Error("Expected non-nil scraper")
				return
			}

			providers := scraper.registry.GetProviders()
			if len(providers) != tt.expectedCount {
				t.Errorf("Expected %d providers, got %d", tt.expectedCount, len(providers))
			}
		})
	}
}

func TestScrapeMetadata(t *testing.T) {
	// Create a simple HTML document
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "head",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "title",
				FirstChild: &html.Node{
					Type: html.TextNode,
					Data: "Test Page",
				},
			},
		},
	}

	result, err := ScrapeMetadata(doc)

	if err != nil {
		t.Errorf("ScrapeMetadata() returned error: %v", err)
	}

	if result == nil {
		t.Error("ScrapeMetadata() returned nil result")
	}
}

func TestScrapeMetadata_NilDocument(t *testing.T) {
	result, err := ScrapeMetadata(nil)

	if err == nil {
		t.Error("Expected error for nil document")
	}

	if result != nil {
		t.Error("Expected nil result for nil document")
	}
}

func TestScrapeMetadataWithProviders(t *testing.T) {
	mockProvider := &MockProvider{name: "test", priority: 1, element: "meta"}
	providers := []metadata.MetadataProvider{mockProvider}

	// Create HTML document with meta tag
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "meta",
			Attr: []html.Attribute{
				{Key: "content", Val: "Test Content"},
			},
		},
	}

	result, err := ScrapeMetadataWithProviders(doc, providers)

	if err != nil {
		t.Errorf("ScrapeMetadataWithProviders() returned error: %v", err)
	}

	if result == nil {
		t.Error("ScrapeMetadataWithProviders() returned nil result")
	}
}

func TestScrapeMetadataWithProviderNames(t *testing.T) {
	providerNames := []string{"openGraph", "twitter"}

	// Create a simple HTML document
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "head",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "og:title"},
					{Key: "content", Val: "Test Title"},
				},
			},
		},
	}

	result, err := ScrapeMetadataWithProviderNames(doc, providerNames)

	if err != nil {
		t.Errorf("ScrapeMetadataWithProviderNames() returned error: %v", err)
	}

	if result == nil {
		t.Error("ScrapeMetadataWithProviderNames() returned nil result")
	}
}

func TestScrapeMetadataWithProviderNames_InvalidProvider(t *testing.T) {
	providerNames := []string{"invalid"}

	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
	}

	result, err := ScrapeMetadataWithProviderNames(doc, providerNames)

	if err == nil {
		t.Error("Expected error for invalid provider name")
	}

	if result != nil {
		t.Error("Expected nil result for invalid provider name")
	}
}
