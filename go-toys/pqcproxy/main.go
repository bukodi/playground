package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// See:

func main1() {
	// Replace 'target' with the URL of the server you want to proxy to
	target, err := url.Parse("https://noreg.hu")
	if err != nil {
		panic(err)
	}

	// Create a new ReverseProxy instance
	proxy := httputil.NewSingleHostReverseProxy(target)
	// Configure the reverse proxy to use HTTPS
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Create a handler function that logs the URL and forwards the request to the proxy
	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		r.Host = target.Host

		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}
		conn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = conn

		//w.Header().Set("X-Ben", "radi")
		proxy.ServeHTTP(w, r)
	}

	// Register the handler function with the HTTP server
	http.HandleFunc("/", handler)

	// Start the HTTP server
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}

}
