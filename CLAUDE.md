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

### Code Quality
- `go fmt ./...` - Format all code
- `go mod tidy` - Clean up dependencies
- `golangci-lint run` - Run linter (if installed)

## Architecture

The project uses a **provider-based architecture** with the following key design patterns:

### Core Flow
1. **HTML Parsing**: Uses `golang.org/x/net/html` to parse web pages
2. **Provider System**: Each provider implements `MetadataProvider` interface to extract specific types of metadata
3. **Priority Resolution**: Providers have priority numbers (lower = higher priority) for value resolution
4. **Method Chaining**: Scraper uses fluent interface for sequential HTML element processing

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
- Uses method chaining: `scrapeMetaTags().scrapeTitleTag().scrapeHeadingTags()`
- Walks HTML DOM tree for each element type
- Delegates extraction to provider registry
- Builds final `Metadata` result object

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
1. Implement `MetadataProvider` interface
2. Define unique `Name()` and `Priority()`
3. Implement `CanHandle()` logic for HTML elements
4. Extract data in `Scrape()` method
5. Add to provider loader in `pkg/providers/loader.go`

### Factory Pattern
`pkg/scraper/factory.go` provides convenience functions:
- `CreateScraper()` - Auto-loads default providers
- `CreateScraperWithProviders()` - Use custom provider list
- `ScrapeMetadata()` - One-shot scraping function

## Testing Notes

Tests use table-driven patterns typical for Go. When adding tests:
- Use `pkg/metadata/metadata_test.go` as a reference
- Test both success and error cases
- Use `*html.Node` for DOM testing (from `golang.org/x/net/html`)
- Mock providers by implementing `MetadataProvider` interface