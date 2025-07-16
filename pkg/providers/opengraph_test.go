package providers

import (
	"testing"

	"golang.org/x/net/html"
)

func TestOpenGraphProvider_Name(t *testing.T) {
	provider := NewOpenGraphProvider()
	if provider.Name() != "openGraph" {
		t.Errorf("Expected name 'openGraph', got '%s'", provider.Name())
	}
}

func TestOpenGraphProvider_Priority(t *testing.T) {
	provider := NewOpenGraphProvider()
	if provider.Priority() != 1 {
		t.Errorf("Expected priority 1, got %d", provider.Priority())
	}
}

func TestOpenGraphProvider_CanHandle(t *testing.T) {
	provider := NewOpenGraphProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected bool
	}{
		{
			name: "meta tag with og: property",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "og:title"},
					{Key: "content", Val: "Test Title"},
				},
			},
			expected: true,
		},
		{
			name: "meta tag with og: name",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "og:description"},
					{Key: "content", Val: "Test Description"},
				},
			},
			expected: true,
		},
		{
			name: "meta tag without og: prefix",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "description"},
					{Key: "content", Val: "Test Description"},
				},
			},
			expected: false,
		},
		{
			name: "non-meta element",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "div",
			},
			expected: false,
		},
		{
			name: "text node",
			node: &html.Node{
				Type: html.TextNode,
				Data: "text content",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.CanHandle(tt.node)
			if result != tt.expected {
				t.Errorf("CanHandle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOpenGraphProvider_Scrape(t *testing.T) {
	provider := NewOpenGraphProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected *struct {
			key   string
			value string
		}
	}{
		{
			name: "valid og:title meta tag",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "og:title"},
					{Key: "content", Val: "Test Title"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "title", value: "Test Title"},
		},
		{
			name: "valid og:description with name attribute",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "og:description"},
					{Key: "content", Val: "Test Description"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "description", value: "Test Description"},
		},
		{
			name: "meta tag without content",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "og:title"},
				},
			},
			expected: nil,
		},
		{
			name: "meta tag without property/name",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "content", Val: "Test Content"},
				},
			},
			expected: nil,
		},
		{
			name: "non-og meta tag",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "description"},
					{Key: "content", Val: "Test Description"},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.Scrape(tt.node)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Scrape() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Error("Scrape() = nil, want non-nil result")
				return
			}

			if result.Key != tt.expected.key {
				t.Errorf("Scrape().Key = %v, want %v", result.Key, tt.expected.key)
			}

			if result.Value != tt.expected.value {
				t.Errorf("Scrape().Value = %v, want %v", result.Value, tt.expected.value)
			}
		})
	}
}

func TestOpenGraphProvider_GetValue(t *testing.T) {
	provider := NewOpenGraphProvider()

	tests := []struct {
		name     string
		key      string
		data     map[string][]string
		expected *string
	}{
		{
			name: "existing key with value",
			key:  "title",
			data: map[string][]string{
				"title": {"Test Title"},
			},
			expected: stringPtr("Test Title"),
		},
		{
			name: "existing key with multiple values",
			key:  "title",
			data: map[string][]string{
				"title": {"First Title", "Second Title"},
			},
			expected: stringPtr("First Title"),
		},
		{
			name: "non-existing key",
			key:  "title",
			data: map[string][]string{
				"description": {"Test Description"},
			},
			expected: nil,
		},
		{
			name: "existing key with empty values",
			key:  "title",
			data: map[string][]string{
				"title": {},
			},
			expected: nil,
		},
		{
			name:     "empty data",
			key:      "title",
			data:     map[string][]string{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.GetValue(tt.key, tt.data)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("GetValue() = %v, want nil", *result)
				}
				return
			}

			if result == nil {
				t.Error("GetValue() = nil, want non-nil result")
				return
			}

			if *result != *tt.expected {
				t.Errorf("GetValue() = %v, want %v", *result, *tt.expected)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
