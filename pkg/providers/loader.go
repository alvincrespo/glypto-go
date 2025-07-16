package providers

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"plugin"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
)

// Loader manages dynamic loading of metadata providers
type Loader struct {
	defaultProviders []metadata.MetadataProvider
}

// NewLoader creates a new provider loader
func NewLoader() *Loader {
	return &Loader{
		defaultProviders: []metadata.MetadataProvider{
			NewOpenGraphProvider(),
			NewTwitterProvider(),
			NewStandardMetaProvider(),
			NewOtherElementsProvider(),
		},
	}
}

// LoadFromDirectory loads providers from a directory (plugin-based)
func (l *Loader) LoadFromDirectory(dir string) ([]metadata.MetadataProvider, error) {
	var providers []metadata.MetadataProvider

	if dir == "" {
		return l.defaultProviders, nil
	}

	// Walk through directory looking for .so files (plugins)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".so" {
			return nil
		}

		// Load the plugin
		p, err := plugin.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open plugin %s: %w", path, err)
		}

		// Look for the NewProvider function
		sym, err := p.Lookup("NewProvider")
		if err != nil {
			return fmt.Errorf("plugin %s does not export NewProvider function: %w", path, err)
		}

		// Assert that it's a function that returns MetadataProvider
		newProvider, ok := sym.(func() metadata.MetadataProvider)
		if !ok {
			return fmt.Errorf("plugin %s NewProvider function has wrong signature", path)
		}

		// Create the provider instance
		provider := newProvider()
		providers = append(providers, provider)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load providers from directory %s: %w", dir, err)
	}

	// If no providers were loaded from directory, return defaults
	if len(providers) == 0 {
		return l.defaultProviders, nil
	}

	return providers, nil
}

// LoadDefaults returns the default built-in providers
func (l *Loader) LoadDefaults() []metadata.MetadataProvider {
	return l.defaultProviders
}

// LoadFromList loads providers from a provided list
func (l *Loader) LoadFromList(providerNames []string) ([]metadata.MetadataProvider, error) {
	var providers []metadata.MetadataProvider

	providerMap := map[string]metadata.MetadataProvider{
		"openGraph": NewOpenGraphProvider(),
		"twitter":   NewTwitterProvider(),
		"meta":      NewStandardMetaProvider(),
		"other":     NewOtherElementsProvider(),
	}

	for _, name := range providerNames {
		if provider, exists := providerMap[name]; exists {
			providers = append(providers, provider)
		} else {
			return nil, fmt.Errorf("unknown provider: %s", name)
		}
	}

	if len(providers) == 0 {
		return l.defaultProviders, nil
	}

	return providers, nil
}

// GetAvailableProviders returns a list of available built-in provider names
func (l *Loader) GetAvailableProviders() []string {
	return []string{"openGraph", "twitter", "meta", "other"}
}
