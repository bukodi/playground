package toauth

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
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
	rootHandler.HandleFunc("/api/user", handleGetCurrentUser)
	rootHandler.HandleFunc("/api/user/token/{provider}", handleGetUserToken)

	srv := httptest.NewUnstartedServer(rootHandler)
	var err error
	srv.Listener, err = net.Listen("tcp", httpAddress)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	srv.Start()
	defer srv.Close()
	t.Logf("Server started on:\n%s\n", strings.Replace(srv.URL, "127.0.0.1", "localhost", -1))
	t.Logf("User repository file: %s\n", DefaultUserRepository.filePath)
	time.Sleep(600 * time.Second)
}
