package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Define the path to the HTML directory
	htmlDir := "../http/frontend/src"

	// Serve HTML files from the "frontend/src" directory
	fs := http.FileServer(http.Dir(htmlDir))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	// Start the HTTP server on port 8080
	http.ListenAndServe("localhost:8080", r)
}
