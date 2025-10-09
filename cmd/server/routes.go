package main

import (
	"net/http"

	code_submissions "github.com/Bu1raj/byte-forge-backend/internal/api/code_submissions"
)

// RegisterRoutes registers the HTTP routes for the server.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/submit", code_submissions.SubmitHandler)
	mux.HandleFunc("/result/", code_submissions.ResultHandler)
}
