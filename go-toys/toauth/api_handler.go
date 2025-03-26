package toauth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleOauthLoginURL(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	cfg, err := DefaultRegistry.GetProvider(provider)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unknown provider: %s", provider), 404)
		return
	}

	// Clone the URL and set the path to /oauth2
	clone := *r.URL
	clone.Path = "/oauth2"

	// Use the provider's configuration to generate the auth URL
	url := cfg.AuthCodeURL(provider)
	http.Redirect(w, r, url, 302)
}

// handleListOauthProviders returns a list of all registered OAuth providers
func handleListOauthProviders(w http.ResponseWriter, r *http.Request) {
	// Get the list of providers from the registry
	providers := DefaultRegistry.ListProviders()

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the list of providers as JSON and write it to the response
	if err := json.NewEncoder(w).Encode(providers); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding providers: %v", err), http.StatusInternalServerError)
		return
	}
}
