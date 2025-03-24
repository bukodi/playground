package toauth

import (
	"golang.org/x/oauth2"
	"os"
)

var (
	//  Open the https://console.cloud.google.com/auth/clients page select theproject and the demo-app client
	//  and set these enviroment variables copy the client id and secret
	DEMO_APP_GOOGLE_CLIENT_ID     = os.Getenv("DEMO_APP_GOOGLE_CLIENT_ID")
	DEMO_APP_GOOGLE_CLIENT_SECRET = os.Getenv("DEMO_APP_GOOGLE_CLIENT_SECRET")
)

func googleOauthConfig(redirectURL string) *oauth2.Config {
	config := oauth2.Config{
		ClientID:     DEMO_APP_GOOGLE_CLIENT_ID,
		ClientSecret: DEMO_APP_GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	return &config
}
