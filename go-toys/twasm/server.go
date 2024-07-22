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
	fs := http.FileServer(http.FS(content))
	http.Handle("/", fs)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
