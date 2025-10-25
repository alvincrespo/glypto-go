# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Glypto Go is a CLI tool for scraping metadata from websites using a provider-based architecture. It extracts Open Graph tags, Twitter Cards, standard meta tags, and RSS/Atom feeds from web pages. The project is a Go translation of the original TypeScript Glypto project.

## Development Commands

### Building and Running
- `go build -o bin/glypto ./cmd/glypto` - Build the CLI binary
- `go run ./cmd/glypto scrape https://example.com` - Build and run directly
- `./bin/glypto scrape https://example.com` - Run built binary with URL
- `./bin/glypto scrape` - Run in interactive mode (prompts for URL)

### Testing
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output
- `go test -cover ./...` - Run tests with coverage report
- `go test ./pkg/metadata -v` - Run specific package tests
- `go test -race ./...` - Run tests with race detection
- `go test ./pkg/scraper -v` - Test scraper engine
- `go test ./pkg/providers -v` - Test provider implementations

### Code Quality
- `go fmt ./...` - Format all code
- `go mod tidy` - Clean up dependencies
- `golangci-lint run` - Run linter (if installed)

## Prerequisites

- **Go 1.24+** (checked in `go.mod`)
- **Dependencies**: `golang.org/x/net/html` (parsing), `spf13/cobra` (CLI), `fatih/color` (output formatting)

## Architecture

The project uses a **provider-based architecture** with the following key design patterns:

### Core Flow
1. **CLI Entry** (`cmd/glypto/main.go` + `pkg/cli/`): Cobra-based CLI with interactive URL prompts
2. **HTTP Fetching** (`pkg/cli/scrape.go`): Fetches web pages using standard `net/http`
3. **HTML Parsing**: Uses `golang.org/x/net/html` to parse response body into DOM tree
4. **Provider System**: Each provider implements `MetadataProvider` interface to extract specific metadata types
5. **Priority Resolution**: Providers have priority numbers (lower = higher priority) for intelligent fallback
6. **Method Chaining**: Scraper uses fluent interface for sequential HTML element type processing
7. **Colored Output** (`pkg/cli/scrape.go`): Uses `fatih/color` for CLI formatting

### Key Components

**MetadataProvider Interface** (`pkg/metadata/types.go`):
- `Name()` - Unique provider identifier
- `Priority()` - Resolution priority (1=highest, 4=lowest)
- `CanHandle(node *html.Node)` - Determines if provider can process element
- `Scrape(node *html.Node)` - Extracts key-value data from element
- `GetValue(key, data)` - Resolves final value from provider's data

**Provider Registry** (`pkg/providers/registry.go`):
- `ProviderRegistry` struct manages all providers
- Automatically sorts providers by priority
- `ScrapeFromElement()` tries providers in priority order until one succeeds
- `ResolveValue()` uses provider priority to resolve metadata values

**Scraper Engine** (`pkg/scraper/scraper.go`):
- Uses method chaining: `scrapeMetaTags().scrapeTitleTag().scrapeHeadingTags().scrapeLinkTags().scrapeFeedLinks()`
- Each method walks HTML DOM tree targeting specific element types (`<meta>`, `<title>`, `<h1>`, `<link>`)
- Delegates extraction to provider registry for priority-based provider resolution
- Builds final `Metadata` result object with aggregated provider data

**CLI Package** (`pkg/cli/`):
- `root.go`: Main command setup with Cobra
- `scrape.go`: HTTP fetching, CLI output formatting, interactive URL prompting
- Uses `fatih/color` for colored console output

**Built-in Providers** (priority order):
1. **OpenGraph** (priority 1): Extracts `og:*` properties
2. **Twitter** (priority 2): Extracts `twitter:*` properties  
3. **StandardMeta** (priority 3): Extracts standard meta tags
4. **OtherElements** (priority 4): Extracts `<title>`, `<h1>`, `<link>` elements

### Package Structure
- `pkg/metadata/` - Core types, interfaces, and metadata result object
- `pkg/providers/` - Provider implementations and registry
- `pkg/scraper/` - Main scraping engine and factory functions
- `pkg/cli/` - Cobra CLI commands
- `cmd/glypto/` - CLI entry point

### Import Cycle Prevention
The architecture carefully avoids import cycles:
- `metadata` package defines core interfaces (`MetadataProvider`, `Registry`)
- `providers` package implements providers and registry
- `scraper` package depends on both metadata and providers
- `cli` package only depends on scraper

### Adding New Providers
1. Implement `MetadataProvider` interface from `pkg/metadata/types.go`
2. Define unique `Name()` (string identifier) and `Priority()` (int: 1=highest, 4+=lowest)
3. Implement `CanHandle(node *html.Node)` to identify relevant elements
4. Extract key-value data in `Scrape(node *html.Node)` method
5. Implement `GetValue(key string, data map[string][]string)` for value resolution
6. Add to provider loader in `pkg/providers/loader.go`

**Reference Implementations**: See `pkg/providers/opengraph.go`, `twitter.go`, or `standardmeta.go` for examples of different extraction patterns. Also see `pkg/providers/base.go` for shared provider utilities.

### Factory Pattern
`pkg/scraper/factory.go` provides convenience functions:
- `CreateScraper()` - Auto-loads all default providers (OpenGraph, Twitter, StandardMeta, OtherElements)
- `CreateScraperWithProviders(providerList)` - Create scraper with custom `[]MetadataProvider` instances
- `CreateScraperWithProviderNames(names)` - Create scraper by provider name strings (e.g., `[]string{"opengraph", "twitter"}`)
- `ScrapeMetadata(doc)` - One-shot scraping using default providers; returns `*Metadata`

Use `CreateScraperWithProviderNames()` for CLI scenarios where users specify providers by name. Use `CreateScraperWithProviders()` for programmatic APIs with custom provider instances.

## Testing Notes

Tests use table-driven patterns typical for Go. When adding tests:
- Use `pkg/metadata/metadata_test.go` as a reference
- Test both success and error cases
- Use `*html.Node` for DOM testing (from `golang.org/x/net/html`)
- Mock providers by implementing `MetadataProvider` interface
- Run specific test file: `go test ./pkg/metadata -v -run TestMetadata`

## Common Development Workflows

### Adding a New Metadata Field to Results
1. Add getter method to `Metadata` struct in `pkg/metadata/metadata.go`
2. Implement extraction logic in relevant providers (e.g., add to `opengraph.go` for new `og:*` property)
3. Add test cases to provider test files
4. Update CLI output formatting in `pkg/cli/scrape.go` if needed

### Creating a Custom Provider
1. Create new file `pkg/providers/customprovider.go`
2. Implement `MetadataProvider` interface with unique name and priority
3. Write tests in `pkg/providers/customprovider_test.go`
4. Register in `pkg/providers/loader.go` if it should be auto-loaded

### Debugging Scraper Output
- Run with verbose test output: `go test -v ./pkg/scraper`
- Check provider priority resolution in `pkg/providers/registry.go`
- Verify `CanHandle()` logic in individual provider implementations
- Use `go run ./cmd/glypto scrape URL` to test live with real websites

### Modifying CLI Output
- Format code is in `pkg/cli/scrape.go`
- Uses `github.com/fatih/color` for styling
- Check existing color scheme before adding new outputs