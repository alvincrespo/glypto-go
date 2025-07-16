package providers

import (
	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
	"strings"
)

// BaseProvider provides common functionality for all metadata providers
type BaseProvider struct{}

// getAttribute gets an attribute value from a node
func (b *BaseProvider) getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getTextContent extracts text content from a node
func (b *BaseProvider) getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result.WriteString(b.getTextContent(c))
	}
	return strings.TrimSpace(result.String())
}

// GetValue resolves a value for a given key
func (b *BaseProvider) GetValue(key string, data map[string][]string) *string {
	if values, exists := data[key]; exists && len(values) > 0 {
		return &values[0]
	}
	return nil
}

// scrapeMetaTag provides common meta tag scraping logic for providers
func (b *BaseProvider) scrapeMetaTag(node *html.Node, prefixToRemove string) *metadata.ScrapedData {
	property := b.getAttribute(node, "property")
	if property == "" {
		property = b.getAttribute(node, "name")
	}

	content := b.getAttribute(node, "content")

	if property == "" || content == "" {
		return nil
	}

	key := strings.TrimPrefix(property, prefixToRemove)

	return &metadata.ScrapedData{
		Key:   key,
		Value: content,
	}
}
