package toauth

import (
	"fmt"
	"golang.org/x/oauth2"
	"sync"
)

// OAuthProviderRegistry is a registry for OAuth2 providers
type OAuthProviderRegistry struct {
	providers map[string]*oauth2.Config
	mu        sync.RWMutex
}

// NewOAuthProviderRegistry creates a new registry for OAuth2 providers
func NewOAuthProviderRegistry() *OAuthProviderRegistry {
	return &OAuthProviderRegistry{
		providers: make(map[string]*oauth2.Config),
	}
}

// RegisterProvider registers an OAuth2 provider with the registry
func (r *OAuthProviderRegistry) RegisterProvider(name string, config *oauth2.Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = config
}

// GetProvider returns the OAuth2 configuration for the specified provider
func (r *OAuthProviderRegistry) GetProvider(name string) (*oauth2.Config, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	config, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return config, nil
}

// HasProvider checks if a provider is registered
func (r *OAuthProviderRegistry) HasProvider(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.providers[name]
	return ok
}

// ListProviders returns a list of registered provider names
func (r *OAuthProviderRegistry) ListProviders() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var providers []string
	for name := range r.providers {
		providers = append(providers, name)
	}
	return providers
}

// DefaultRegistry is the default OAuth2 provider registry
var DefaultRegistry = NewOAuthProviderRegistry()

// Initialize the default registry with the Google provider
func init() {
	// Register Google provider
	DefaultRegistry.RegisterProvider("google", &oauth2.Config{
		ClientID:     DEMO_APP_GOOGLE_CLIENT_ID,
		ClientSecret: DEMO_APP_GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	})
}
