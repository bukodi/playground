package toauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"time"
)

func handleOauthLoginURL(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	cfg, err := DefaultRegistry.GetProvider(provider)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unknown provider: %s", provider), 404)
		return
	}

	// Clone the URL and set the path to /oauth2
	rurl := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   "/oauth2",
	}

	// Use the provider's configuration to generate the auth URL
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("redirect_uri", rurl.String()),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("prompt", "select_account"),
		//	oauth2.SetAuthURLParam("login_hint", "bukodi@gmail.com"),
	}

	loginUrl := cfg.AuthCodeURL(provider, opts...)

	http.Redirect(w, r, loginUrl, 302)
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

// handleGetCurrentUser returns information about the current user
func handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	// Get the user from the repository
	user, err := DefaultUserRepository.GetUserByID(cookie.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		return
	}

	// Create a response object with user information
	type UserResponse struct {
		ID     string   `json:"id"`
		Email  string   `json:"email"`
		Logins []string `json:"logins"`
	}

	// Extract provider names from tokens
	var logins []string
	for _, token := range user.Tokens {
		logins = append(logins, token.ProviderName)
	}

	response := UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Logins: logins,
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the user as JSON and write it to the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding user: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleGetUserToken returns a token for a specific provider
func handleGetUserToken(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	// Get the provider from the URL
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "Provider not specified", http.StatusBadRequest)
		return
	}

	// Get the token from the repository
	token, err := DefaultUserRepository.GetToken(cookie.Value, provider)
	if err != nil {
		http.Error(w, fmt.Sprintf("Token not found: %v", err), http.StatusNotFound)
		return
	}

	// Create a response object with token information
	type TokenResponse struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token,omitempty"`
		Expiry       time.Time `json:"expiry"`
	}

	response := TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the token as JSON and write it to the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding token: %v", err), http.StatusInternalServerError)
		return
	}
}
