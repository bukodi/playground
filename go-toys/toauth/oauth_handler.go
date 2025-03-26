package toauth

import (
	"fmt"
	"net/http"
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
	fmt.Fprintf(w, "Token: %v", token)
}
