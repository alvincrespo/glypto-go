package providers

import (
	"strings"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

const OGPrefix = "og:"

// OpenGraphProvider extracts OpenGraph metadata
type OpenGraphProvider struct {
	BaseProvider
}

// NewOpenGraphProvider creates a new OpenGraph provider
func NewOpenGraphProvider() *OpenGraphProvider {
	return &OpenGraphProvider{}
}

// Name returns the provider name
func (p *OpenGraphProvider) Name() string {
	return "openGraph"
}

// Priority returns the provider priority (highest priority)
func (p *OpenGraphProvider) Priority() int {
	return 1
}

// CanHandle determines if this provider can handle the given element
func (p *OpenGraphProvider) CanHandle(node *html.Node) bool {
	if node.Type != html.ElementNode || node.Data != "meta" {
		return false
	}

	property := p.getAttribute(node, "property")
	name := p.getAttribute(node, "name")

	return strings.HasPrefix(property, OGPrefix) || strings.HasPrefix(name, OGPrefix)
}

// Scrape extracts OpenGraph data from the element
func (p *OpenGraphProvider) Scrape(node *html.Node) *metadata.ScrapedData {
	if !p.CanHandle(node) {
		return nil
	}

	return p.scrapeMetaTag(node, OGPrefix)
}
