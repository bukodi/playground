package toauth

import (
	"fmt"
	"net/http"
)

func handleOauthLoginURL(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "google" {
		clone := *r.URL
		clone.Path = "/oauth2"
		cfg := googleOauthConfig(clone.String())
		url := cfg.AuthCodeURL("google")
		http.Redirect(w, r, url, 302)
	} else {
		http.Error(w, fmt.Sprintf("Unknown provider: %s", provider), 404)
	}
}
