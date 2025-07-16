package providers

import (
	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

// OtherElementsProvider extracts metadata from other HTML elements
type OtherElementsProvider struct {
	BaseProvider
}

// NewOtherElementsProvider creates a new other elements provider
func NewOtherElementsProvider() *OtherElementsProvider {
	return &OtherElementsProvider{}
}

// Name returns the provider name
func (p *OtherElementsProvider) Name() string {
	return "other"
}

// Priority returns the provider priority (lowest priority)
func (p *OtherElementsProvider) Priority() int {
	return 4
}

// CanHandle determines if this provider can handle the given element
func (p *OtherElementsProvider) CanHandle(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}

	switch node.Data {
	case "title", "h1":
		return true
	case "link":
		rel := p.getAttribute(node, "rel")
		return rel == "icon" || rel == "shortcut icon" || rel == "canonical"
	default:
		return false
	}
}

// Scrape extracts data from other HTML elements
func (p *OtherElementsProvider) Scrape(node *html.Node) *metadata.ScrapedData {
	if !p.CanHandle(node) {
		return nil
	}

	switch node.Data {
	case "title":
		content := p.getTextContent(node)
		if content != "" {
			return &metadata.ScrapedData{
				Key:   "title",
				Value: content,
			}
		}
	case "h1":
		content := p.getTextContent(node)
		if content != "" {
			return &metadata.ScrapedData{
				Key:   "firstHeading",
				Value: content,
			}
		}
	case "link":
		rel := p.getAttribute(node, "rel")
		href := p.getAttribute(node, "href")
		if rel != "" && href != "" {
			switch rel {
			case "icon", "shortcut icon":
				return &metadata.ScrapedData{
					Key:   rel,
					Value: href,
				}
			case "canonical":
				return &metadata.ScrapedData{
					Key:   "url",
					Value: href,
				}
			}
		}
	}

	return nil
}
