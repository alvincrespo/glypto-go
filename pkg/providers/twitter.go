package providers

import (
	"strings"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

const TwitterPrefix = "twitter:"

// TwitterProvider extracts Twitter Card metadata
type TwitterProvider struct {
	BaseProvider
}

// NewTwitterProvider creates a new Twitter provider
func NewTwitterProvider() *TwitterProvider {
	return &TwitterProvider{}
}

// Name returns the provider name
func (p *TwitterProvider) Name() string {
	return "twitter"
}

// Priority returns the provider priority (second highest priority)
func (p *TwitterProvider) Priority() int {
	return 2
}

// CanHandle determines if this provider can handle the given element
func (p *TwitterProvider) CanHandle(node *html.Node) bool {
	if node.Type != html.ElementNode || node.Data != "meta" {
		return false
	}

	property := p.getAttribute(node, "property")
	name := p.getAttribute(node, "name")

	return strings.HasPrefix(property, TwitterPrefix) || strings.HasPrefix(name, TwitterPrefix)
}

// Scrape extracts Twitter Card data from the element
func (p *TwitterProvider) Scrape(node *html.Node) *metadata.ScrapedData {
	if !p.CanHandle(node) {
		return nil
	}

	return p.scrapeMetaTag(node, TwitterPrefix)
}
