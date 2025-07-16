package scraper

import (
	"strings"
	"testing"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

// MockRegistry for testing
type MockRegistry struct {
	providers []metadata.MetadataProvider
}

func (m *MockRegistry) GetProviders() []metadata.MetadataProvider {
	return m.providers
}

func (m *MockRegistry) ScrapeFromElement(node *html.Node) *metadata.ScrapingResult {
	for _, provider := range m.providers {
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

func (m *MockRegistry) ResolveValue(key string, providerData metadata.ProviderData) *string {
	for _, provider := range m.providers {
		if data, exists := providerData[provider.Name()]; exists {
			if value := provider.GetValue(key, data); value != nil {
				return value
			}
		}
	}
	return nil
}

func (m *MockRegistry) AddProvider(provider metadata.MetadataProvider) {
	m.providers = append(m.providers, provider)
}

func (m *MockRegistry) RemoveProvider(name string) {
	for i, provider := range m.providers {
		if provider.Name() == name {
			m.providers = append(m.providers[:i], m.providers[i+1:]...)
			return
		}
	}
}

func (m *MockRegistry) GetProvider(name string) metadata.MetadataProvider {
	for _, provider := range m.providers {
		if provider.Name() == name {
			return provider
		}
	}
	return nil
}

// MockProvider for testing
type MockProvider struct {
	name     string
	priority int
	element  string
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Priority() int {
	return m.priority
}

func (m *MockProvider) CanHandle(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == m.element
}

func (m *MockProvider) Scrape(node *html.Node) *metadata.ScrapedData {
	if !m.CanHandle(node) {
		return nil
	}

	switch m.element {
	case "meta":
		content := ""
		for _, attr := range node.Attr {
			if attr.Key == "content" {
				content = attr.Val
				break
			}
		}
		if content != "" {
			return &metadata.ScrapedData{
				Key:   "test",
				Value: content,
			}
		}
	case "title":
		if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
			return &metadata.ScrapedData{
				Key:   "title",
				Value: node.FirstChild.Data,
			}
		}
	}
	return nil
}

func (m *MockProvider) GetValue(key string, data map[string][]string) *string {
	if values, exists := data[key]; exists && len(values) > 0 {
		return &values[0]
	}
	return nil
}

func TestNewScraper(t *testing.T) {
	registry := &MockRegistry{}
	scraper := NewScraper(registry)

	if scraper == nil {
		t.Error("NewScraper() returned nil")
	}

	if scraper.registry != registry {
		t.Error("NewScraper() did not set registry correctly")
	}
}

func TestScraper_Scrape_NilDocument(t *testing.T) {
	registry := &MockRegistry{}
	scraper := NewScraper(registry)

	result, err := scraper.Scrape(nil)

	if err == nil {
		t.Error("Expected error for nil document, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for nil document")
	}

	expectedError := "HTML document cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestScraper_Scrape_ValidDocument(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1, element: "meta"}
	registry := &MockRegistry{providers: []metadata.MetadataProvider{provider}}
	scraper := NewScraper(registry)

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
					{Key: "name", Val: "description"},
					{Key: "content", Val: "Test Description"},
				},
			},
		},
	}

	result, err := scraper.Scrape(doc)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestScraper_scrapeMetaTags(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1, element: "meta"}
	registry := &MockRegistry{providers: []metadata.MetadataProvider{provider}}
	scraper := NewScraper(registry)
	scraper.result = metadata.NewMetadata(registry)

	// Create HTML with meta tag
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
	scraper.doc = doc

	result := scraper.scrapeMetaTags()

	if result != scraper {
		t.Error("scrapeMetaTags() should return scraper for chaining")
	}
}

func TestScraper_scrapeTitleTag(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1, element: "title"}
	registry := &MockRegistry{providers: []metadata.MetadataProvider{provider}}
	scraper := NewScraper(registry)
	scraper.result = metadata.NewMetadata(registry)

	// Create HTML with title tag
	titleNode := &html.Node{
		Type: html.ElementNode,
		Data: "title",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: "Test Title",
		},
	}
	doc := &html.Node{
		Type:       html.ElementNode,
		Data:       "html",
		FirstChild: titleNode,
	}
	scraper.doc = doc

	result := scraper.scrapeTitleTag()

	if result != scraper {
		t.Error("scrapeTitleTag() should return scraper for chaining")
	}
}

func TestScraper_scrapeHeadingTags(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1, element: "h1"}
	registry := &MockRegistry{providers: []metadata.MetadataProvider{provider}}
	scraper := NewScraper(registry)
	scraper.result = metadata.NewMetadata(registry)

	// Create HTML with h1 tag
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "h1",
			FirstChild: &html.Node{
				Type: html.TextNode,
				Data: "Test Heading",
			},
		},
	}
	scraper.doc = doc

	result := scraper.scrapeHeadingTags()

	if result != scraper {
		t.Error("scrapeHeadingTags() should return scraper for chaining")
	}
}

func TestScraper_scrapeLinkTags(t *testing.T) {
	provider := &MockProvider{name: "test", priority: 1, element: "link"}
	registry := &MockRegistry{providers: []metadata.MetadataProvider{provider}}
	scraper := NewScraper(registry)
	scraper.result = metadata.NewMetadata(registry)

	// Create HTML with link tag
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "link",
			Attr: []html.Attribute{
				{Key: "rel", Val: "canonical"},
				{Key: "href", Val: "https://example.com"},
			},
		},
	}
	scraper.doc = doc

	result := scraper.scrapeLinkTags()

	if result != scraper {
		t.Error("scrapeLinkTags() should return scraper for chaining")
	}
}

func TestScraper_scrapeFeedLinks(t *testing.T) {
	registry := &MockRegistry{}
	scraper := NewScraper(registry)
	scraper.result = metadata.NewMetadata(registry)

	// Create HTML with RSS feed link
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "link",
			Attr: []html.Attribute{
				{Key: "rel", Val: "alternate"},
				{Key: "type", Val: "application/rss+xml"},
				{Key: "title", Val: "RSS Feed"},
				{Key: "href", Val: "/feed.rss"},
			},
		},
	}
	scraper.doc = doc

	result := scraper.scrapeFeedLinks()

	if result != scraper {
		t.Error("scrapeFeedLinks() should return scraper for chaining")
	}

	// Check if feed was added
	if len(scraper.result.Feeds) != 1 {
		t.Errorf("Expected 1 feed, got %d", len(scraper.result.Feeds))
	}

	if scraper.result.Feeds[0].Type != "application/rss+xml" {
		t.Errorf("Expected feed type 'application/rss+xml', got '%s'", scraper.result.Feeds[0].Type)
	}

	if scraper.result.Feeds[0].Href != "/feed.rss" {
		t.Errorf("Expected feed href '/feed.rss', got '%s'", scraper.result.Feeds[0].Href)
	}

	if scraper.result.Feeds[0].Title == nil || *scraper.result.Feeds[0].Title != "RSS Feed" {
		t.Error("Expected feed title 'RSS Feed'")
	}
}

func TestScraper_scrapeFeedLinks_NoTitle(t *testing.T) {
	registry := &MockRegistry{}
	scraper := NewScraper(registry)
	scraper.result = metadata.NewMetadata(registry)

	// Create HTML with RSS feed link without title
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "link",
			Attr: []html.Attribute{
				{Key: "rel", Val: "alternate"},
				{Key: "type", Val: "application/atom+xml"},
				{Key: "href", Val: "/feed.atom"},
			},
		},
	}
	scraper.doc = doc

	scraper.scrapeFeedLinks()

	if len(scraper.result.Feeds) != 1 {
		t.Errorf("Expected 1 feed, got %d", len(scraper.result.Feeds))
	}

	if scraper.result.Feeds[0].Title != nil {
		t.Error("Expected feed title to be nil")
	}
}

func TestScraper_getAttribute(t *testing.T) {
	scraper := &Scraper{}

	node := &html.Node{
		Type: html.ElementNode,
		Data: "meta",
		Attr: []html.Attribute{
			{Key: "name", Val: "description"},
			{Key: "content", Val: "Test Content"},
		},
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"name", "description"},
		{"content", "Test Content"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		result := scraper.getAttribute(node, tt.key)
		if result != tt.expected {
			t.Errorf("getAttribute(%s) = %v, want %v", tt.key, result, tt.expected)
		}
	}
}

func TestScraper_hasAttribute(t *testing.T) {
	scraper := &Scraper{}

	node := &html.Node{
		Type: html.ElementNode,
		Data: "meta",
		Attr: []html.Attribute{
			{Key: "name", Val: "description"},
			{Key: "content", Val: "Test Content"},
		},
	}

	tests := []struct {
		key      string
		expected bool
	}{
		{"name", true},
		{"content", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		result := scraper.hasAttribute(node, tt.key)
		if result != tt.expected {
			t.Errorf("hasAttribute(%s) = %v, want %v", tt.key, result, tt.expected)
		}
	}
}

func TestScraper_getTextContent(t *testing.T) {
	scraper := &Scraper{}

	tests := []struct {
		name     string
		node     *html.Node
		expected string
	}{
		{
			name: "text node",
			node: &html.Node{
				Type: html.TextNode,
				Data: "Simple text",
			},
			expected: "Simple text",
		},
		{
			name: "element with text child",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "span",
				FirstChild: &html.Node{
					Type: html.TextNode,
					Data: "Child text",
				},
			},
			expected: "Child text",
		},
		{
			name: "element with multiple children",
			node: func() *html.Node {
				parent := &html.Node{
					Type: html.ElementNode,
					Data: "div",
				}
				child1 := &html.Node{
					Type: html.TextNode,
					Data: "  First  ",
				}
				child2 := &html.Node{
					Type: html.TextNode,
					Data: "  Second  ",
				}
				parent.FirstChild = child1
				child1.NextSibling = child2
				return parent
			}(),
			expected: "First    Second",
		},
		{
			name: "empty element",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "div",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scraper.getTextContent(tt.node)
			if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
				t.Errorf("getTextContent() = '%v', want '%v'", result, tt.expected)
			}
		})
	}
}
