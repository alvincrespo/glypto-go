package metadata

// Metadata represents the scraped metadata from a webpage
type Metadata struct {
	providerData ProviderData
	registry     Registry
	Feeds        []*Feed
}

// NewMetadata creates a new Metadata instance
func NewMetadata(registry Registry) *Metadata {
	m := &Metadata{
		providerData: make(ProviderData),
		registry:     registry,
		Feeds:        make([]*Feed, 0),
	}

	// Initialize provider data maps
	for _, provider := range registry.GetProviders() {
		m.providerData[provider.Name()] = make(map[string][]string)
	}

	return m
}

// AddData adds scraped data to the metadata
func (m *Metadata) AddData(providerName, key, value string) {
	if m.providerData[providerName] == nil {
		m.providerData[providerName] = make(map[string][]string)
	}

	data := m.providerData[providerName]
	data[key] = append(data[key], value)
}

// resolveValue resolves a value using the provider registry
func (m *Metadata) resolveValue(key string) *string {
	if m.registry == nil {
		return nil
	}
	return m.registry.ResolveValue(key, m.providerData)
}

// Favicon returns the favicon URL with fallback
func (m *Metadata) Favicon() string {
	if icon := m.resolveValue("icon"); icon != nil {
		return *icon
	}
	if shortcutIcon := m.resolveValue("shortcut icon"); shortcutIcon != nil {
		return *shortcutIcon
	}
	return "/favicon.ico"
}

// Title returns the page title
func (m *Metadata) Title() *string {
	if title := m.resolveValue("title"); title != nil {
		return title
	}
	return m.resolveValue("firstHeading")
}

// Description returns the page description
func (m *Metadata) Description() *string {
	return m.resolveValue("description")
}

// Image returns the page image URL
func (m *Metadata) Image() *string {
	return m.resolveValue("image")
}

// URL returns the canonical URL
func (m *Metadata) URL() *string {
	return m.resolveValue("url")
}

// SiteName returns the site name
func (m *Metadata) SiteName() *string {
	if siteName := m.resolveValue("site_name"); siteName != nil {
		return siteName
	}
	// Twitter uses 'site' instead of 'site_name'
	return m.resolveValue("site")
}

// GetProviderData returns the raw provider data for a specific provider
func (m *Metadata) GetProviderData(providerName string) map[string][]string {
	if data, exists := m.providerData[providerName]; exists {
		return data
	}
	return make(map[string][]string)
}

// OpenGraph returns OpenGraph data for backward compatibility
func (m *Metadata) OpenGraph() map[string][]string {
	return m.GetProviderData("openGraph")
}

// TwitterCard returns Twitter Card data for backward compatibility
func (m *Metadata) TwitterCard() map[string][]string {
	return m.GetProviderData("twitter")
}

// Meta returns standard meta data for backward compatibility
func (m *Metadata) Meta() map[string][]string {
	return m.GetProviderData("meta")
}

// Other returns other elements data for backward compatibility
func (m *Metadata) Other() map[string][]string {
	return m.GetProviderData("other")
}
