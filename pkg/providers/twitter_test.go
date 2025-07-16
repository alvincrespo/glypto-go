package providers

import (
	"testing"

	"golang.org/x/net/html"
)

func TestTwitterProvider_Name(t *testing.T) {
	provider := NewTwitterProvider()
	if provider.Name() != "twitter" {
		t.Errorf("Expected name 'twitter', got '%s'", provider.Name())
	}
}

func TestTwitterProvider_Priority(t *testing.T) {
	provider := NewTwitterProvider()
	if provider.Priority() != 2 {
		t.Errorf("Expected priority 2, got %d", provider.Priority())
	}
}

func TestTwitterProvider_CanHandle(t *testing.T) {
	provider := NewTwitterProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected bool
	}{
		{
			name: "meta tag with twitter: property",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "twitter:card"},
					{Key: "content", Val: "summary"},
				},
			},
			expected: true,
		},
		{
			name: "meta tag with twitter: name",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "twitter:title"},
					{Key: "content", Val: "Test Title"},
				},
			},
			expected: true,
		},
		{
			name: "meta tag without twitter: prefix",
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

func TestTwitterProvider_Scrape(t *testing.T) {
	provider := NewTwitterProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected *struct {
			key   string
			value string
		}
	}{
		{
			name: "valid twitter:card meta tag",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "twitter:card"},
					{Key: "content", Val: "summary"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "card", value: "summary"},
		},
		{
			name: "valid twitter:title with name attribute",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "name", Val: "twitter:title"},
					{Key: "content", Val: "Test Title"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "title", value: "Test Title"},
		},
		{
			name: "meta tag without content",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "meta",
				Attr: []html.Attribute{
					{Key: "property", Val: "twitter:card"},
				},
			},
			expected: nil,
		},
		{
			name: "non-twitter meta tag",
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

func TestTwitterProvider_GetValue(t *testing.T) {
	provider := NewTwitterProvider()

	tests := []struct {
		name     string
		key      string
		data     map[string][]string
		expected *string
	}{
		{
			name: "existing key with value",
			key:  "card",
			data: map[string][]string{
				"card": {"summary"},
			},
			expected: stringPtr("summary"),
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
			key:  "card",
			data: map[string][]string{
				"title": {"Test Title"},
			},
			expected: nil,
		},
		{
			name: "existing key with empty values",
			key:  "card",
			data: map[string][]string{
				"card": {},
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
