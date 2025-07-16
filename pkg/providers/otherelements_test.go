package providers

import (
	"testing"

	"golang.org/x/net/html"
)

func TestOtherElementsProvider_Name(t *testing.T) {
	provider := NewOtherElementsProvider()
	if provider.Name() != "other" {
		t.Errorf("Expected name 'other', got '%s'", provider.Name())
	}
}

func TestOtherElementsProvider_Priority(t *testing.T) {
	provider := NewOtherElementsProvider()
	if provider.Priority() != 4 {
		t.Errorf("Expected priority 4, got %d", provider.Priority())
	}
}

func TestOtherElementsProvider_CanHandle(t *testing.T) {
	provider := NewOtherElementsProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected bool
	}{
		{
			name: "title element",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "title",
			},
			expected: true,
		},
		{
			name: "h1 element",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "h1",
			},
			expected: true,
		},
		{
			name: "link element with icon rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "icon"},
					{Key: "href", Val: "/favicon.ico"},
				},
			},
			expected: true,
		},
		{
			name: "link element with shortcut icon rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "shortcut icon"},
					{Key: "href", Val: "/favicon.ico"},
				},
			},
			expected: true,
		},
		{
			name: "link element with canonical rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "canonical"},
					{Key: "href", Val: "https://example.com"},
				},
			},
			expected: true,
		},
		{
			name: "link element with stylesheet rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "stylesheet"},
					{Key: "href", Val: "/style.css"},
				},
			},
			expected: false,
		},
		{
			name: "div element",
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

func TestOtherElementsProvider_Scrape(t *testing.T) {
	provider := NewOtherElementsProvider()

	tests := []struct {
		name     string
		node     *html.Node
		expected *struct {
			key   string
			value string
		}
	}{
		{
			name: "title element with text content",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "title",
				FirstChild: &html.Node{
					Type: html.TextNode,
					Data: "Test Page Title",
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "title", value: "Test Page Title"},
		},
		{
			name: "h1 element with text content",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "h1",
				FirstChild: &html.Node{
					Type: html.TextNode,
					Data: "Main Heading",
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "firstHeading", value: "Main Heading"},
		},
		{
			name: "link element with icon rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "icon"},
					{Key: "href", Val: "/favicon.ico"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "icon", value: "/favicon.ico"},
		},
		{
			name: "link element with shortcut icon rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "shortcut icon"},
					{Key: "href", Val: "/favicon.png"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "shortcut icon", value: "/favicon.png"},
		},
		{
			name: "link element with canonical rel",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "canonical"},
					{Key: "href", Val: "https://example.com/page"},
				},
			},
			expected: &struct {
				key   string
				value string
			}{key: "url", value: "https://example.com/page"},
		},
		{
			name: "empty title element",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "title",
			},
			expected: nil,
		},
		{
			name: "link element without href",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "link",
				Attr: []html.Attribute{
					{Key: "rel", Val: "icon"},
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

func TestOtherElementsProvider_GetValue(t *testing.T) {
	provider := NewOtherElementsProvider()

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
			name: "non-existing key",
			key:  "title",
			data: map[string][]string{
				"description": {"Test Description"},
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

func TestOtherElementsProvider_getTextContent(t *testing.T) {
	provider := NewOtherElementsProvider()

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
			name: "element with single text child",
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
			name: "element with multiple text children",
			node: func() *html.Node {
				parent := &html.Node{
					Type: html.ElementNode,
					Data: "div",
				}
				child1 := &html.Node{
					Type: html.TextNode,
					Data: "First ",
				}
				child2 := &html.Node{
					Type: html.TextNode,
					Data: "Second",
				}
				parent.FirstChild = child1
				child1.NextSibling = child2
				child1.Parent = parent
				child2.Parent = parent
				return parent
			}(),
			expected: "First Second",
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
			result := provider.getTextContent(tt.node)
			if result != tt.expected {
				t.Errorf("getTextContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}
