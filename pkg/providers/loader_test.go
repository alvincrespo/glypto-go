package providers

import (
	"testing"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()

	if loader == nil {
		t.Error("NewLoader() returned nil")
		return
	}

	if len(loader.defaultProviders) == 0 {
		t.Error("NewLoader() created loader with no default providers")
	}

	// Check that all expected default providers are present
	expectedProviders := []string{"openGraph", "twitter", "meta", "other"}
	if len(loader.defaultProviders) != len(expectedProviders) {
		t.Errorf("Expected %d default providers, got %d", len(expectedProviders), len(loader.defaultProviders))
	}
}

func TestLoader_LoadDefaults(t *testing.T) {
	loader := NewLoader()
	providers := loader.LoadDefaults()

	if len(providers) != 4 {
		t.Errorf("Expected 4 default providers, got %d", len(providers))
	}

	// Check provider names and priorities
	expectedData := []struct {
		name     string
		priority int
	}{
		{"openGraph", 1},
		{"twitter", 2},
		{"meta", 3},
		{"other", 4},
	}

	for i, provider := range providers {
		if provider.Name() != expectedData[i].name {
			t.Errorf("Expected provider %d to have name '%s', got '%s'", i, expectedData[i].name, provider.Name())
		}

		if provider.Priority() != expectedData[i].priority {
			t.Errorf("Expected provider %d to have priority %d, got %d", i, expectedData[i].priority, provider.Priority())
		}
	}
}

func TestLoader_LoadFromDirectory_EmptyDir(t *testing.T) {
	loader := NewLoader()
	providers, err := loader.LoadFromDirectory("")

	if err != nil {
		t.Errorf("LoadFromDirectory(\"\") returned error: %v", err)
	}

	if len(providers) != 4 {
		t.Errorf("Expected 4 default providers for empty directory, got %d", len(providers))
	}
}

func TestLoader_LoadFromDirectory_NonexistentDir(t *testing.T) {
	loader := NewLoader()
	providers, err := loader.LoadFromDirectory("/nonexistent/directory")

	// Should return an error but we expect it to fallback to defaults in the factory
	if err == nil {
		// If no error, should have returned defaults
		if len(providers) != 4 {
			t.Error("Expected default providers when directory doesn't exist")
		}
	}
}

func TestLoader_LoadFromList(t *testing.T) {
	loader := NewLoader()

	tests := []struct {
		name          string
		providerNames []string
		expectError   bool
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "all valid providers",
			providerNames: []string{"openGraph", "twitter", "meta", "other"},
			expectError:   false,
			expectedCount: 4,
			expectedNames: []string{"openGraph", "twitter", "meta", "other"},
		},
		{
			name:          "subset of providers",
			providerNames: []string{"openGraph", "twitter"},
			expectError:   false,
			expectedCount: 2,
			expectedNames: []string{"openGraph", "twitter"},
		},
		{
			name:          "single provider",
			providerNames: []string{"meta"},
			expectError:   false,
			expectedCount: 1,
			expectedNames: []string{"meta"},
		},
		{
			name:          "invalid provider",
			providerNames: []string{"invalid"},
			expectError:   true,
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name:          "mixed valid and invalid",
			providerNames: []string{"openGraph", "invalid", "twitter"},
			expectError:   true,
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name:          "empty list",
			providerNames: []string{},
			expectError:   false,
			expectedCount: 4, // Should return defaults
			expectedNames: []string{"openGraph", "twitter", "meta", "other"},
		},
		{
			name:          "duplicate providers",
			providerNames: []string{"openGraph", "openGraph"},
			expectError:   false,
			expectedCount: 2,
			expectedNames: []string{"openGraph", "openGraph"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			providers, err := loader.LoadFromList(tt.providerNames)

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

			if len(providers) != tt.expectedCount {
				t.Errorf("Expected %d providers, got %d", tt.expectedCount, len(providers))
				return
			}

			for i, provider := range providers {
				if i < len(tt.expectedNames) {
					if provider.Name() != tt.expectedNames[i] {
						t.Errorf("Expected provider %d to have name '%s', got '%s'", i, tt.expectedNames[i], provider.Name())
					}
				}
			}
		})
	}
}

func TestLoader_LoadFromList_UnknownProvider(t *testing.T) {
	loader := NewLoader()

	providers, err := loader.LoadFromList([]string{"unknown"})

	if err == nil {
		t.Error("Expected error for unknown provider")
	}

	if providers != nil {
		t.Error("Expected nil providers for unknown provider")
	}

	expectedError := "unknown provider: unknown"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoader_GetAvailableProviders(t *testing.T) {
	loader := NewLoader()
	available := loader.GetAvailableProviders()

	expected := []string{"openGraph", "twitter", "meta", "other"}

	if len(available) != len(expected) {
		t.Errorf("Expected %d available providers, got %d", len(expected), len(available))
	}

	for i, name := range expected {
		if available[i] != name {
			t.Errorf("Expected provider %d to be '%s', got '%s'", i, name, available[i])
		}
	}
}
