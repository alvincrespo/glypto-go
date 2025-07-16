package metadata

import (
	"golang.org/x/net/html"
	"testing"
)

func TestMetadata_Favicon(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Metadata
		expected string
	}{
		{
			name: "returns default favicon when no icon found",
			setup: func() *Metadata {
				return &Metadata{
					providerData: make(ProviderData),
				}
			},
			expected: "/favicon.ico",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			result := m.Favicon()
			if result != tt.expected {
				t.Errorf("Favicon() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMetadata_AddData(t *testing.T) {
	m := &Metadata{
		providerData: make(ProviderData),
	}

	m.AddData("test", "key", "value")

	if m.providerData["test"] == nil {
		t.Error("Expected provider data to be initialized")
	}

	if len(m.providerData["test"]["key"]) != 1 {
		t.Error("Expected one value in provider data")
	}

	if m.providerData["test"]["key"][0] != "value" {
		t.Error("Expected value to be stored correctly")
	}
}

func TestMetadata_GetProviderData(t *testing.T) {
	m := &Metadata{
		providerData: ProviderData{
			"test": map[string][]string{
				"key": {"value"},
			},
		},
	}

	result := m.GetProviderData("test")
	if len(result) != 1 {
		t.Error("Expected one key in provider data")
	}

	if result["key"][0] != "value" {
		t.Error("Expected value to be returned correctly")
	}

	// Test non-existent provider
	empty := m.GetProviderData("nonexistent")
	if len(empty) != 0 {
		t.Error("Expected empty map for non-existent provider")
	}
}

// MockRegistry for testing
type MockRegistry struct {
	providers []MetadataProvider
}

func (m *MockRegistry) GetProviders() []MetadataProvider {
	return m.providers
}

func (m *MockRegistry) ScrapeFromElement(node *html.Node) *ScrapingResult {
	return nil
}

func (m *MockRegistry) ResolveValue(key string, providerData ProviderData) *string {
	for _, provider := range m.providers {
		if data, exists := providerData[provider.Name()]; exists {
			if value := provider.GetValue(key, data); value != nil {
				return value
			}
		}
	}
	return nil
}

func (m *MockRegistry) AddProvider(provider MetadataProvider) {
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

func (m *MockRegistry) GetProvider(name string) MetadataProvider {
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
	data     map[string][]string
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Priority() int {
	return m.priority
}

func (m *MockProvider) CanHandle(node *html.Node) bool {
	return true
}

func (m *MockProvider) Scrape(node *html.Node) *ScrapedData {
	return nil
}

func (m *MockProvider) GetValue(key string, data map[string][]string) *string {
	if m.data != nil {
		if values, exists := m.data[key]; exists && len(values) > 0 {
			return &values[0]
		}
	}
	if values, exists := data[key]; exists && len(values) > 0 {
		return &values[0]
	}
	return nil
}

func TestNewMetadata(t *testing.T) {
	mockProvider1 := &MockProvider{name: "test1", priority: 1}
	mockProvider2 := &MockProvider{name: "test2", priority: 2}
	registry := &MockRegistry{providers: []MetadataProvider{mockProvider1, mockProvider2}}

	m := NewMetadata(registry)

	if m == nil {
		t.Error("NewMetadata() returned nil")
	}

	if m.registry == nil {
		t.Error("NewMetadata() did not set registry correctly")
	}

	if len(m.providerData) != 2 {
		t.Errorf("Expected 2 provider data entries, got %d", len(m.providerData))
	}

	if m.Feeds == nil {
		t.Error("NewMetadata() did not initialize Feeds slice")
	}

	if len(m.Feeds) != 0 {
		t.Error("NewMetadata() should initialize empty Feeds slice")
	}
}

func TestMetadata_Title(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Metadata
		expected *string
	}{
		{
			name: "returns title from provider data",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"title": {"Test Title"}}}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				m := NewMetadata(registry)
				m.AddData("test", "title", "Test Title")
				return m
			},
			expected: stringPtr("Test Title"),
		},
		{
			name: "returns firstHeading when no title",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"firstHeading": {"First Heading"}}}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				m := NewMetadata(registry)
				m.AddData("test", "firstHeading", "First Heading")
				return m
			},
			expected: stringPtr("First Heading"),
		},
		{
			name: "returns nil when no title or heading",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				return NewMetadata(registry)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			result := m.Title()

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Title() = %v, want nil", *result)
				}
				return
			}

			if result == nil {
				t.Error("Title() = nil, want non-nil")
				return
			}

			if *result != *tt.expected {
				t.Errorf("Title() = %v, want %v", *result, *tt.expected)
			}
		})
	}
}

func TestMetadata_Description(t *testing.T) {
	mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"description": {"Test Description"}}}
	registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
	m := NewMetadata(registry)
	m.AddData("test", "description", "Test Description")

	result := m.Description()
	if result == nil {
		t.Error("Description() = nil, want non-nil")
		return
	}

	if *result != "Test Description" {
		t.Errorf("Description() = %v, want %v", *result, "Test Description")
	}
}

func TestMetadata_Image(t *testing.T) {
	mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"image": {"https://example.com/image.jpg"}}}
	registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
	m := NewMetadata(registry)
	m.AddData("test", "image", "https://example.com/image.jpg")

	result := m.Image()
	if result == nil {
		t.Error("Image() = nil, want non-nil")
		return
	}

	if *result != "https://example.com/image.jpg" {
		t.Errorf("Image() = %v, want %v", *result, "https://example.com/image.jpg")
	}
}

func TestMetadata_URL(t *testing.T) {
	mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"url": {"https://example.com"}}}
	registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
	m := NewMetadata(registry)
	m.AddData("test", "url", "https://example.com")

	result := m.URL()
	if result == nil {
		t.Error("URL() = nil, want non-nil")
		return
	}

	if *result != "https://example.com" {
		t.Errorf("URL() = %v, want %v", *result, "https://example.com")
	}
}

func TestMetadata_SiteName(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Metadata
		expected *string
	}{
		{
			name: "returns site_name when available",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"site_name": {"Example Site"}}}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				m := NewMetadata(registry)
				m.AddData("test", "site_name", "Example Site")
				return m
			},
			expected: stringPtr("Example Site"),
		},
		{
			name: "returns site when site_name not available",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"site": {"@example"}}}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				m := NewMetadata(registry)
				m.AddData("test", "site", "@example")
				return m
			},
			expected: stringPtr("@example"),
		},
		{
			name: "returns nil when neither available",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				return NewMetadata(registry)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			result := m.SiteName()

			if tt.expected == nil {
				if result != nil {
					t.Errorf("SiteName() = %v, want nil", *result)
				}
				return
			}

			if result == nil {
				t.Error("SiteName() = nil, want non-nil")
				return
			}

			if *result != *tt.expected {
				t.Errorf("SiteName() = %v, want %v", *result, *tt.expected)
			}
		})
	}
}

func TestMetadata_Favicon_WithValues(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Metadata
		expected string
	}{
		{
			name: "returns icon when available",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"icon": {"/custom-icon.png"}}}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				m := NewMetadata(registry)
				m.AddData("test", "icon", "/custom-icon.png")
				return m
			},
			expected: "/custom-icon.png",
		},
		{
			name: "returns shortcut icon when icon not available",
			setup: func() *Metadata {
				mockProvider := &MockProvider{name: "test", priority: 1, data: map[string][]string{"shortcut icon": {"/shortcut.ico"}}}
				registry := &MockRegistry{providers: []MetadataProvider{mockProvider}}
				m := NewMetadata(registry)
				m.AddData("test", "shortcut icon", "/shortcut.ico")
				return m
			},
			expected: "/shortcut.ico",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			result := m.Favicon()

			if result != tt.expected {
				t.Errorf("Favicon() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMetadata_OpenGraph(t *testing.T) {
	m := &Metadata{
		providerData: ProviderData{
			"openGraph": map[string][]string{
				"title": {"OG Title"},
			},
		},
	}

	result := m.OpenGraph()
	if len(result) != 1 {
		t.Errorf("Expected 1 key in OpenGraph data, got %d", len(result))
	}

	if result["title"][0] != "OG Title" {
		t.Errorf("Expected 'OG Title', got '%s'", result["title"][0])
	}
}

func TestMetadata_TwitterCard(t *testing.T) {
	m := &Metadata{
		providerData: ProviderData{
			"twitter": map[string][]string{
				"card": {"summary"},
			},
		},
	}

	result := m.TwitterCard()
	if len(result) != 1 {
		t.Errorf("Expected 1 key in TwitterCard data, got %d", len(result))
	}

	if result["card"][0] != "summary" {
		t.Errorf("Expected 'summary', got '%s'", result["card"][0])
	}
}

func TestMetadata_Meta(t *testing.T) {
	m := &Metadata{
		providerData: ProviderData{
			"meta": map[string][]string{
				"description": {"Meta Description"},
			},
		},
	}

	result := m.Meta()
	if len(result) != 1 {
		t.Errorf("Expected 1 key in Meta data, got %d", len(result))
	}

	if result["description"][0] != "Meta Description" {
		t.Errorf("Expected 'Meta Description', got '%s'", result["description"][0])
	}
}

func TestMetadata_Other(t *testing.T) {
	m := &Metadata{
		providerData: ProviderData{
			"other": map[string][]string{
				"title": {"Other Title"},
			},
		},
	}

	result := m.Other()
	if len(result) != 1 {
		t.Errorf("Expected 1 key in Other data, got %d", len(result))
	}

	if result["title"][0] != "Other Title" {
		t.Errorf("Expected 'Other Title', got '%s'", result["title"][0])
	}
}

func TestMetadata_resolveValue_NilRegistry(t *testing.T) {
	m := &Metadata{
		registry:     nil,
		providerData: make(ProviderData),
	}

	result := m.resolveValue("title")
	if result != nil {
		t.Error("Expected nil for resolveValue with nil registry")
	}
}

func stringPtr(s string) *string {
	return &s
}
