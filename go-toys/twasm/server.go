//go:build !js

package main

import (
	"embed"
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed index.html
//go:embed wasm_exec.js
//go:embed main.wasm
var content embed.FS

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(content)))

	fmt.Println("Starting server on :9090")
	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}

}
