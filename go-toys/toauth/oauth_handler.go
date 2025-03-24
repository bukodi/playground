package toauth

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

func oauth2Handler(w http.ResponseWriter, r *http.Request) {

	var oauthCfg *oauth2.Config

	r.ParseForm()
	state := r.Form.Get("state")
	if state == "google" {
		oauthCfg = googleOauthConfig(r.URL.String())
	} else {
		http.Error(w, "State invalid", http.StatusBadRequest)
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
	fmt.Sprintf("Token: %v", token)
}
