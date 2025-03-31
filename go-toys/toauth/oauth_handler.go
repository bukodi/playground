package toauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UserInfo represents basic user information from OAuth providers
type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

func oauth2Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")

	// Get the provider configuration from the registry
	oauthCfg, err := DefaultRegistry.GetProvider(state)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid state: %s", state), http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := oauthCfg.Exchange(r.Context(), code /*, oauth2.SetAuthURLParam("code_verifier", "s256example")*/)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info from the provider
	userInfo, err := getUserInfo(r.Context(), state, token.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	// Get or create the user in the repository
	user, err := DefaultUserRepository.GetOrCreateUserByEmail(userInfo.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get or create user: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the token in the repository
	if err := DefaultUserRepository.StoreToken(user.ID, state, token); err != nil {
		http.Error(w, fmt.Sprintf("Failed to store token: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the user ID in a cookie for session management
	userIDCookie := &http.Cookie{
		Name:     "user_id",
		Value:    user.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   86400 * 30, // 30 days
	}
	http.SetCookie(w, userIDCookie)

	// Store the token in a cookie for backward compatibility
	tokenCookie := &http.Cookie{
		Name:     "oauth_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: false, // Set to false so JavaScript can access it
		Secure:   r.TLS != nil,
		MaxAge:   int(token.Expiry.Sub(time.Now()).Seconds()), // Set expiry to match token expiry
	}
	http.SetCookie(w, tokenCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// getUserInfo gets user information from the OAuth provider
func getUserInfo(ctx context.Context, provider, accessToken string) (*UserInfo, error) {
	var userInfoURL string

	// Determine the user info URL based on the provider
	switch provider {
	case "google":
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Create a request to the user info endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add the access token to the request
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", body)
	}

	// Parse the response
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}
