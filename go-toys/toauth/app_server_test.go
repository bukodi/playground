package toauth

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const httpAddress = "localhost:9094"

func TestOauth2(t *testing.T) {
	rootHandler := http.NewServeMux()
	rootHandler.HandleFunc("/oauth2", oauth2Handler)
	rootHandler.HandleFunc("/", staticHandler)
	rootHandler.HandleFunc("/api/oauthLoginURL/{provider}", handleOauthLoginURL)
	rootHandler.HandleFunc("/api/oauthProviders", handleListOauthProviders)

	srv := httptest.NewUnstartedServer(rootHandler)
	var err error
	srv.Listener, err = net.Listen("tcp", httpAddress)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	srv.Start()
	defer srv.Close()
	t.Logf("Server started on:\n%s\n", srv.URL)
	time.Sleep(60 * time.Second)

}
