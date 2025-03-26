package toauth

import (
	"fmt"
	"net/http"
	"time"
)

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

	// Store the token in a cookie
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
