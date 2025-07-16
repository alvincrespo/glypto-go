package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"golang.org/x/net/html"
)

func TestGetURLFromInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expected    string
		expectError bool
	}{
		{
			name:        "URL provided as argument",
			args:        []string{"https://example.com"},
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:        "empty args",
			args:        []string{},
			expected:    "",
			expectError: true, // Will fail because we can't simulate stdin in this test
		},
		{
			name:        "multiple args (takes first)",
			args:        []string{"https://example.com", "https://test.com"},
			expected:    "https://example.com",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args) == 0 {
				// Skip stdin test as it's hard to mock in unit tests
				t.Skip("Skipping stdin test")
				return
			}

			result, err := getURLFromInput(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFetchWebpage(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Test</title></head></html>"))
	}))
	defer server.Close()

	// Test successful fetch
	resp, err := fetchWebpage(server.URL)
	if err != nil {
		t.Errorf("fetchWebpage() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestFetchWebpage_HTTPError(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	resp, err := fetchWebpage(server.URL)

	if err == nil {
		if resp != nil {
			resp.Body.Close()
		}
		t.Error("Expected error for 404 response")
	}

	if !strings.Contains(err.Error(), "HTTP error! status: 404") {
		t.Errorf("Expected HTTP error message, got: %v", err)
	}
}

func TestFetchWebpage_InvalidURL(t *testing.T) {
	resp, err := fetchWebpage("invalid-url")

	if err == nil {
		if resp != nil {
			resp.Body.Close()
		}
		t.Error("Expected error for invalid URL")
	}
}

func TestParseHTML(t *testing.T) {
	htmlContent := "<html><head><title>Test</title></head><body><h1>Hello</h1></body></html>"

	// Create a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get test response: %v", err)
	}
	defer resp.Body.Close()

	doc, err := parseHTML(resp)
	if err != nil {
		t.Errorf("parseHTML() failed: %v", err)
	}

	if doc == nil {
		t.Error("parseHTML() returned nil document")
	}

	if doc.Type != html.DocumentNode {
		t.Error("parseHTML() did not return a document node")
	}
}

func TestParseHTML_InvalidHTML(t *testing.T) {
	// Even invalid HTML should parse successfully with html.Parse
	invalidHTML := "<html><head><title>Test</head><body><h1>Hello</body></html>"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(invalidHTML))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get test response: %v", err)
	}
	defer resp.Body.Close()

	doc, err := parseHTML(resp)
	if err != nil {
		t.Errorf("parseHTML() failed on invalid HTML: %v", err)
	}

	if doc == nil {
		t.Error("parseHTML() returned nil document for invalid HTML")
	}
}

func TestScrapeMetadata(t *testing.T) {
	// Create a simple HTML document
	doc := &html.Node{
		Type: html.DocumentNode,
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "html",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "head",
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: "title",
					FirstChild: &html.Node{
						Type: html.TextNode,
						Data: "Test Page",
					},
				},
			},
		},
	}

	result, err := scrapeMetadata(doc)
	if err != nil {
		t.Errorf("scrapeMetadata() failed: %v", err)
	}

	if result == nil {
		t.Error("scrapeMetadata() returned nil result")
	}
}

func TestDisplayResults(t *testing.T) {
	// Create test metadata
	testMetadata := &metadata.Metadata{}

	// Capture output by redirecting stdout (this is a basic test)
	// In a real scenario, you might want to use dependency injection
	// to make the display function more testable

	// This test mainly ensures the function doesn't panic
	displayResults(testMetadata)
}

func TestPrintField(t *testing.T) {
	// Capture output
	old := bytes.NewBuffer(nil)

	tests := []struct {
		name  string
		field string
		value *string
	}{
		{
			name:  "field with value",
			field: "Title",
			value: stringPtr("Test Title"),
		},
		{
			name:  "field with nil value",
			field: "Description",
			value: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test mainly ensures the function doesn't panic
			printField(tt.field, tt.value)

			// Reset buffer
			old.Reset()
		})
	}
}

func TestPrintProviderData(t *testing.T) {
	tests := []struct {
		name  string
		title string
		data  map[string][]string
	}{
		{
			name:  "data with values",
			title: "Test Data",
			data: map[string][]string{
				"title":       {"Test Title"},
				"description": {"Test Description", "Alternative Description"},
			},
		},
		{
			name:  "empty data",
			title: "Empty Data",
			data:  map[string][]string{},
		},
		{
			name:  "nil data",
			title: "Nil Data",
			data:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test mainly ensures the function doesn't panic
			printProviderData(tt.title, tt.data)
		})
	}
}

func TestScrapeCmd(t *testing.T) {
	if scrapeCmd.Use != "scrape [URL]" {
		t.Errorf("Expected Use to be 'scrape [URL]', got '%s'", scrapeCmd.Use)
	}

	if scrapeCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}

	if scrapeCmd.Long == "" {
		t.Error("Expected Long description to be set")
	}

	if scrapeCmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

// Helper function for tests
func stringPtr(s string) *string {
	return &s
}
