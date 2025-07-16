package providers

import (
	"testing"

	"golang.org/x/net/html"
)

func TestStandardMetaProvider_Name(t *testing.T) {
	provider := NewStandardMetaProvider()
	if provider.Name() != "meta" {
		t.Errorf("Expected name 'meta', got '%s'", provider.Name())
	}
}

func TestStandardMetaProvider_Priority(t *testing.T) {
	provider := NewStandardMetaProvider()
	if provider.Priority() != 3 {
		t.Errorf("Expected priority 3, got %d", provider.Priority())
	}
}

func TestStandardMetaProvider_CanHandle(t *testing.T) {
	provider := NewStandardMetaProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected bool
	}{
		{
			name: "standard meta tag with name",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "description"},
					{Key: "content", Val: "Test Description"},
				},
			},
			expected: true,
		},
		{
			name: "standard meta tag with property",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "author"},
					{Key: "content", Val: "John Doe"},
				},
			},
			expected: true,
		},
		{
			name: "og: meta tag",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "og:title"},
					{Key: "content", Val: "Test Title"},
				},
			},
			expected: false,
		},
		{
			name: "twitter: meta tag",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "twitter:card"},
					{Key: "content", Val: "summary"},
				},
			},
			expected: false,
		},
		{
			name: "meta tag without name or property",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "content", Val: "Test Content"},
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

func TestStandardMetaProvider_Scrape(t *testing.T) {
	provider := NewStandardMetaProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected *struct {
			key   string
			value string
		}
	}{
		{
			name: "valid meta tag with name",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "description"},
					{Key: "content", Val: "Test Description"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "description", value: "Test Description"},
		},
		{
			name: "valid meta tag with property",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "author"},
					{Key: "content", Val: "John Doe"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "author", value: "John Doe"},
		},
		{
			name: "meta tag without content",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "description"},
				},
			},
			expected: nil,
		},
		{
			name: "og: meta tag (should not handle)",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "og:title"},
					{Key: "content", Val: "Test Title"},
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

func TestStandardMetaProvider_GetValue(t *testing.T) {
	provider := NewStandardMetaProvider()

	tests := []struct {
		name     string
		key      string
		data     map[string][]string
		expected *string
	}{
		{
			name: "existing key with value",
			key:  "description",
			data: map[string][]string{
				"description": {"Test Description"},
			},
			expected: stringPtr("Test Description"),
		},
		{
			name: "non-existing key",
			key:  "description",
			data: map[string][]string{
				"title": {"Test Title"},
			},
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
