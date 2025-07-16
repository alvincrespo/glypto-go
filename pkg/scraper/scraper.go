package scraper

import (
	"fmt"
	"strings"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

// Scraper provides metadata extraction functionality
type Scraper struct {
	registry metadata.Registry
	doc      *html.Node
	result   *metadata.Metadata
}

// NewScraper creates a new scraper instance
func NewScraper(registry metadata.Registry) *Scraper {
	return &Scraper{
		registry: registry,
	}
}

// Scrape extracts metadata from an HTML document
func (s *Scraper) Scrape(doc *html.Node) (*metadata.Metadata, error) {
	if doc == nil {
		return nil, fmt.Errorf("HTML document cannot be nil")
	}

	s.doc = doc
	s.result = metadata.NewMetadata(s.registry)

	return s.scrapeMetaTags().
		scrapeTitleTag().
		scrapeHeadingTags().
		scrapeLinkTags().
		scrapeFeedLinks().
		getResult(), nil
}

// scrapeMetaTags extracts metadata from <meta> tags
func (s *Scraper) scrapeMetaTags() *Scraper {
	s.walkNodes(s.doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "meta" {
			s.scrapeFromElement(n)
		}
		return true
	})
	return s
}

// scrapeTitleTag extracts data from <title> tag
func (s *Scraper) scrapeTitleTag() *Scraper {
	s.walkNodes(s.doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "title" {
			s.scrapeFromElement(n)
		}
		return true
	})
	return s
}

// scrapeHeadingTags extracts data from <h1> tags
func (s *Scraper) scrapeHeadingTags() *Scraper {
	s.walkNodes(s.doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "h1" {
			s.scrapeFromElement(n)
		}
		return true
	})
	return s
}

// scrapeLinkTags extracts data from <link> tags with rel attribute
func (s *Scraper) scrapeLinkTags() *Scraper {
	s.walkNodes(s.doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "link" && s.hasAttribute(n, "rel") {
			s.scrapeFromElement(n)
		}
		return true
	})
	return s
}

// scrapeFeedLinks extracts RSS/Atom feed links
func (s *Scraper) scrapeFeedLinks() *Scraper {
	s.walkNodes(s.doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "link" {
			rel := s.getAttribute(n, "rel")
			if rel == "alternate" {
				title := s.getAttribute(n, "title")
				feedType := s.getAttribute(n, "type")
				href := s.getAttribute(n, "href")

				if href != "" {
					feed := &metadata.Feed{
						Type: feedType,
						Href: href,
					}
					if title != "" {
						feed.Title = &title
					}
					s.result.Feeds = append(s.result.Feeds, feed)
				}
			}
		}
		return true
	})
	return s
}

// scrapeFromElement attempts to scrape metadata from an element
func (s *Scraper) scrapeFromElement(node *html.Node) {
	if extraction := s.registry.ScrapeFromElement(node); extraction != nil {
		s.result.AddData(
			(*extraction.Provider).Name(),
			extraction.Data.Key,
			extraction.Data.Value,
		)
	}
}

// walkNodes recursively walks through HTML nodes
func (s *Scraper) walkNodes(n *html.Node, fn func(*html.Node) bool) {
	if !fn(n) {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.walkNodes(c, fn)
	}
}

// getAttribute gets an attribute value from a node
func (s *Scraper) getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// hasAttribute checks if a node has an attribute
func (s *Scraper) hasAttribute(n *html.Node, key string) bool {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return true
		}
	}
	return false
}

// getTextContent extracts text content from a node
func (s *Scraper) getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result.WriteString(s.getTextContent(c))
	}
	return strings.TrimSpace(result.String())
}

// getResult returns the scraping result
func (s *Scraper) getResult() *metadata.Metadata {
	return s.result
}
