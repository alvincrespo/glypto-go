# Glypto Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go implementation of Glypto - a CLI tool for scraping metadata from websites using a provider-based architecture.

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

Glypto Go extracts metadata from websites including titles, descriptions, images, Open Graph data, Twitter Cards, and RSS/Atom feeds. It features a modular provider system that makes it easy to add support for new metadata formats.

## Features

- 🔍 **Comprehensive Metadata Scraping**: Open Graph, Twitter Cards, standard meta tags, and more
- 🧩 **Extensible Provider System**: Plug-and-play architecture for adding new metadata sources
- 🚀 **Priority-Based Resolution**: Intelligent fallback system for metadata values
- ⚡ **Fast HTML Parsing**: Built on golang.org/x/net/html for efficient parsing
- 📦 **Multiple Usage Patterns**: CLI tool and programmatic API
- 🎯 **Type-Safe**: Full Go type safety with interfaces

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
    
    fmt.Printf("Title: %s\n", *metadata.Title())
    fmt.Printf("Description: %s\n", *metadata.Description())
    fmt.Printf("Image: %s\n", *metadata.Image())
}
```

#### Custom Providers

```go
package main

import (
    "github.com/alvincrespo/glypto-go/pkg/metadata"
    "github.com/alvincrespo/glypto-go/pkg/providers"
    "github.com/alvincrespo/glypto-go/pkg/scraper"
)

func main() {
    // Create custom provider list
    providerList := []metadata.MetadataProvider{
        providers.NewOpenGraphProvider(),
        providers.NewTwitterProvider(),
    }
    
    // Create scraper with custom providers
    scraperInstance := scraper.CreateScraperWithProviders(providerList)
    
    // Use scraper...
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
├── cmd/glypto/          # CLI entry point
├── pkg/
│   ├── metadata/        # Core metadata types and logic
│   ├── providers/       # Provider implementations and registry
│   ├── scraper/         # Scraping engine and factory functions
│   └── cli/             # CLI commands
├── bin/                 # Compiled binaries
└── go.mod              # Go module definition
```

## Built-in Providers

The following providers are included by default, listed by priority:

1. **OpenGraph Provider** (Priority 1): Extracts `og:*` properties
2. **Twitter Provider** (Priority 2): Extracts `twitter:*` properties  
3. **Standard Meta Provider** (Priority 3): Extracts standard meta tags
4. **Other Elements Provider** (Priority 4): Extracts from `<title>`, `<h1>`, `<link>` tags

## Development

### Prerequisites

- Go 1.20 or higher
- Make (optional, for convenience commands)

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

- `pkg/metadata/metadata_test.go` - Tests for metadata functionality
- Table-driven tests for comprehensive coverage
- Interface-based testing for provider system

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## Acknowledgments

This project is a Go translation of the original [Glypto](https://github.com/alvincrespo/glypto) TypeScript project.