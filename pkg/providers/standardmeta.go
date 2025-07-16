package providers

import (
	"strings"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

// StandardMetaProvider extracts standard meta tag metadata
type StandardMetaProvider struct {
	BaseProvider
}

// NewStandardMetaProvider creates a new standard meta provider
func NewStandardMetaProvider() *StandardMetaProvider {
	return &StandardMetaProvider{}
}

// Name returns the provider name
func (p *StandardMetaProvider) Name() string {
	return "meta"
}

// Priority returns the provider priority (third priority)
func (p *StandardMetaProvider) Priority() int {
	return 3
}

// CanHandle determines if this provider can handle the given element
func (p *StandardMetaProvider) CanHandle(node *html.Node) bool {
	if node.Type != html.ElementNode || node.Data != "meta" {
		return false
	}

	name := p.getAttribute(node, "name")
	property := p.getAttribute(node, "property")

	// Handle standard meta tags that don't have og: or twitter: prefixes
	return (name != "" || property != "") &&
		!strings.HasPrefix(name, OGPrefix) &&
		!strings.HasPrefix(name, TwitterPrefix) &&
		!strings.HasPrefix(property, OGPrefix) &&
		!strings.HasPrefix(property, TwitterPrefix)
}

// Scrape extracts standard meta data from the element
func (p *StandardMetaProvider) Scrape(node *html.Node) *metadata.ScrapedData {
	if !p.CanHandle(node) {
		return nil
	}

	return p.scrapeMetaTag(node, "")
}
