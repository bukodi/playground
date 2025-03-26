package toauth

import (
	"golang.org/x/oauth2"
	"testing"
)

func TestOAuthProviderRegistry(t *testing.T) {
	// Create a new registry
	registry := NewOAuthProviderRegistry()

	// Test registering a provider
	registry.RegisterProvider("test", &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Scopes:       []string{"test-scope"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://test.com/auth",
			TokenURL:  "https://test.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	})

	// Test HasProvider
	if !registry.HasProvider("test") {
		t.Errorf("Expected registry to have provider 'test'")
	}
	if registry.HasProvider("nonexistent") {
		t.Errorf("Expected registry to not have provider 'nonexistent'")
	}

	// Test GetProvider
	provider, err := registry.GetProvider("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider.ClientID != "test-client-id" {
		t.Errorf("Expected ClientID to be 'test-client-id', got %s", provider.ClientID)
	}

	// Test GetProvider with nonexistent provider
	_, err = registry.GetProvider("nonexistent")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test ListProviders
	providers := registry.ListProviders()
	if len(providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(providers))
	}
	if providers[0] != "test" {
		t.Errorf("Expected provider name to be 'test', got %s", providers[0])
	}

	// Test DefaultRegistry
	if !DefaultRegistry.HasProvider("google") {
		t.Errorf("Expected DefaultRegistry to have provider 'google'")
	}
}
