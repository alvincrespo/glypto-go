# Glypto Go

[![CI](https://github.com/alvincrespo/glypto-go/workflows/CI/badge.svg)](https://github.com/alvincrespo/glypto-go/actions)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go implementation of Glypto - a CLI tool for scraping metadata from websites using a provider-based architecture. Extract Open Graph tags, Twitter Cards, standard meta tags, and RSS/Atom feeds from web pages.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [CLI Usage](#cli-usage)
  - [Programmatic Usage](#programmatic-usage)
- [Architecture](#architecture)
- [Built-in Providers](#built-in-providers)
- [Development](#development)
- [Testing](#testing)
- [License](#license)

## Overview

Glypto Go extracts comprehensive metadata from websites including:

- **Page metadata**: Titles, descriptions, images, favicons
- **Social media**: Open Graph tags, Twitter Cards
- **Feed discovery**: RSS/Atom feeds with automatic detection
- **Site information**: Site names, canonical URLs

The tool features a modular provider system with priority-based resolution, making it easy to extend and customize metadata extraction.

## Features

- 🔍 **Comprehensive Metadata Scraping**: Open Graph, Twitter Cards, standard meta tags, and RSS/Atom feeds
- 🧩 **Extensible Provider System**: Plug-and-play architecture for adding new metadata sources
- 🚀 **Priority-Based Resolution**: Intelligent fallback system for metadata values (OpenGraph → Twitter → Standard → Other)
- ⚡ **Fast HTML Parsing**: Built on `golang.org/x/net/html` for efficient parsing
- 📦 **Multiple Usage Patterns**: CLI tool and programmatic Go API
- 🎯 **Type-Safe**: Full Go type safety with interfaces and structured data
- 🎨 **Colorized Output**: Beautiful CLI output with color-coded results
- 📝 **Feed Discovery**: Automatic detection and parsing of RSS/Atom feeds

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/alvincrespo/glypto-go.git
cd glypto-go

# Build the project
go build -o bin/glypto ./cmd/glypto

# Run the CLI
./bin/glypto --help
```

### Using Go Install

```bash
go install github.com/alvincrespo/glypto-go/cmd/glypto@latest
```

## Usage

### CLI Usage

```bash
# Scrape metadata from a URL
./bin/glypto scrape https://example.com

# Interactive mode (will prompt for URL)
./bin/glypto scrape

# Get help
./bin/glypto --help
./bin/glypto scrape --help
```

#### Example Output

```bash
$ ./bin/glypto scrape https://github.com

✓ Metadata scraped successfully:
Title: GitHub · Build and ship software on a single, collaborative platform
Description: Join the world's most widely adopted, AI-powered developer platform...
Image: https://github.githubassets.com/assets/home24-5939032587c9.jpg
URL: https://github.com/
Site Name: GitHub
Favicon: https://github.githubassets.com/favicons/favicon.svg

Feeds:
  1. Untitled () - https://github.com/?locale=ja
  2. Untitled () - https://github.com/?locale=ko

Open Graph Tags:
  site_name: GitHub
  type: object
  title: GitHub · Build and ship software on a single, collaborative platform
  url: https://github.com/
  image: https://github.githubassets.com/assets/home24-5939032587c9.jpg

Twitter Card Tags:
  card: summary_large_image
  site: @github
  title: GitHub · Build and ship software on a single, collaborative platform
```

### Programmatic Usage

#### Simple Usage with Factory

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/alvincrespo/glypto-go/pkg/scraper"
    "golang.org/x/net/html"
)

func main() {
    // Fetch webpage
    resp, err := http.Get("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // Parse HTML
    doc, err := html.Parse(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Scrape metadata
    metadata, err := scraper.ScrapeMetadata(doc)
    if err != nil {
        log.Fatal(err)
    }

    if title := metadata.Title(); title != nil {
        fmt.Printf("Title: %s\n", *title)
    }
    if description := metadata.Description(); description != nil {
        fmt.Printf("Description: %s\n", *description)
    }
    if image := metadata.Image(); image != nil {
        fmt.Printf("Image: %s\n", *image)
    }

    // Access provider-specific data
    ogData := metadata.OpenGraph()
    twitterData := metadata.TwitterCard()
    fmt.Printf("Found %d Open Graph tags\n", len(ogData))
    fmt.Printf("Found %d Twitter Card tags\n", len(twitterData))
}
```

#### Custom Providers

```go
package main

import (
    "log"
    "net/http"

    "github.com/alvincrespo/glypto-go/pkg/metadata"
    "github.com/alvincrespo/glypto-go/pkg/providers"
    "github.com/alvincrespo/glypto-go/pkg/scraper"
    "golang.org/x/net/html"
)

func main() {
    // Create custom provider list (only OpenGraph and Twitter)
    providerList := []metadata.MetadataProvider{
        providers.NewOpenGraphProvider(),
        providers.NewTwitterProvider(),
    }

    // Create scraper with custom providers
    scraperInstance := scraper.CreateScraperWithProviders(providerList)

    // Or use provider names for convenience
    scraperByNames, err := scraper.CreateScraperWithProviderNames([]string{
        "opengraph", "twitter", "standardmeta",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Fetch and parse HTML...
    resp, _ := http.Get("https://example.com")
    defer resp.Body.Close()
    doc, _ := html.Parse(resp.Body)

    // Scrape with custom configuration
    metadata, err := scraperInstance.Scrape(doc)
    if err != nil {
        log.Fatal(err)
    }

    // Process results...
}
```

## Architecture

Glypto Go uses a modular provider architecture with clear separation of concerns:

### Core Components

- **`Scraper`**: Main scraping engine with fluent method chaining
- **`ProviderRegistry`**: Manages and prioritizes metadata providers
- **`Metadata`**: Result object with intelligent value resolution
- **`MetadataProvider`**: Interface for implementing custom providers

### Project Structure

```
glypto-go/
├── .github/             # GitHub Actions workflows and configuration
│   ├── workflows/       # CI/CD pipelines
│   ├── dependabot.yml   # Dependency management
│   └── labeler.yml      # PR auto-labeling
├── cmd/glypto/          # CLI entry point
│   └── main.go          # Application main function
├── pkg/
│   ├── cli/             # Cobra CLI commands and logic
│   ├── metadata/        # Core metadata types and interfaces
│   ├── providers/       # Provider implementations and registry
│   └── scraper/         # Scraping engine and factory functions
├── bin/                 # Compiled binaries (created on build)
├── CLAUDE.md           # AI coding assistant instructions
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Built-in Providers

The following providers are included by default, listed by priority:

1. **OpenGraph Provider** (Priority 1): Extracts `og:*` properties
2. **Twitter Provider** (Priority 2): Extracts `twitter:*` properties
3. **Standard Meta Provider** (Priority 3): Extracts standard meta tags
4. **Other Elements Provider** (Priority 4): Extracts from `<title>`, `<h1>`, `<link>` tags

## Development

### Prerequisites

- Go 1.24 or higher
- Git (for cloning the repository)

### Building

```bash
# Build the CLI
go build -o bin/glypto ./cmd/glypto

# Build and run
go run ./cmd/glypto scrape https://example.com

# Install dependencies
go mod tidy
```

### Project Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

## Testing

The project includes comprehensive tests using Go's built-in testing framework:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/metadata -v

# Run tests with race detection
go test -race ./...
```

### Test Structure

The project includes comprehensive test coverage with:

- **Unit tests** for all packages (`*_test.go` files)
- **Table-driven tests** for comprehensive coverage
- **Interface-based testing** for provider system
- **Integration tests** for CLI commands
- **Mock implementations** for testing provider behavior

**Test Coverage by Package:**
- `pkg/cli/` - CLI command functionality and HTTP handling
- `pkg/metadata/` - Metadata structure and value resolution
- `pkg/providers/` - All provider implementations and registry
- `pkg/scraper/` - Scraping engine and factory functions

## CI/CD

The project includes GitHub Actions workflows for:

- **Continuous Integration**: Automated testing, linting, and building on every push/PR
- **Security Scanning**: Vulnerability checking with `govulncheck`
- **Code Quality**: `golangci-lint` for comprehensive code analysis
- **Dependency Management**: Dependabot for automated dependency updates
- **Releases**: Automated multi-platform binary builds on version tags
- **Auto-labeling**: Automatic PR labeling based on changed files

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## Acknowledgments

This project is a Go translation of the original [Glypto](https://github.com/alvincrespo/glypto) TypeScript project.

## License

MIT License - see [LICENSE](LICENSE) file for details.
