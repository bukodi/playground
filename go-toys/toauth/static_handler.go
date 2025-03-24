package toauth

import (
	"bytes"
	_ "embed"
	"net/http"
	"strings"
	"time"
)

//go:embed index.html
var indexHtml []byte

//go:embed auth.html
var authHtml []byte

//go:embed login.html
var loginHtml []byte

var staticModTime = time.Now()

func staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" || strings.HasSuffix(r.RequestURI, "index.html") {
		http.ServeContent(w, r, "index.html", staticModTime, bytes.NewReader(indexHtml))
	} else if strings.HasSuffix(r.RequestURI, "auth.html") {
		http.ServeContent(w, r, "auth.html", staticModTime, bytes.NewReader(authHtml))
	} else if strings.HasSuffix(r.RequestURI, "login.html") {
		http.ServeContent(w, r, "login.html", staticModTime, bytes.NewReader(loginHtml))
	} else {
		http.Error(w, "Not found", 404)
	}
	return
}
