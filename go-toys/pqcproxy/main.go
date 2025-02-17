package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

// See:

func main() {
	// Replace 'target' with the URL of the server you want to proxy to
	target, err := url.Parse("https://isitquantumsafe.info/")
	if err != nil {
		panic(err)
	}

	// Create a new ReverseProxy instance
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS13,
			GetConfigForClient: nil,
			CurvePreferences:   []tls.CurveID{tls.X25519},
			InsecureSkipVerify: true,
		},
	}
	// Configure the reverse proxy to use HTTPS

	serverCert, err := loadTLSKeyAndCer("test_tls_server")
	if err != nil {
		panic(err)
	}

	s := httptest.NewUnstartedServer(http.HandlerFunc(proxy.ServeHTTP))

	// Configure the server to present the certficate we created
	s.TLS = &tls.Config{
		Certificates:       []tls.Certificate{*serverCert},
		MinVersion:         tls.VersionTLS13,
		GetConfigForClient: nil,
		CurvePreferences:   []tls.CurveID{tls.X25519MLKEM768, tls.X25519},
	}

	tls.NewListener(s.Listener, s.TLS)

	// make a HTTPS request to the server
	s.StartTLS()
	defer s.Close()
	fmt.Printf("Server starting on:\n%s\n", s.URL)
	select {}
}

//go:embed testdata
var testdata embed.FS

func loadTLSKeyAndCer(name string) (*tls.Certificate, error) {
	cerBytes, err := testdata.ReadFile("testdata/" + name + ".cer")
	if err != nil {
		return nil, err
	}
	keyBytes, err := testdata.ReadFile("testdata/" + name + ".pkcs8")
	if err != nil {
		return nil, err
	}
	tlsCert, err := tls.X509KeyPair(cerBytes, keyBytes)
	return &tlsCert, err
}
