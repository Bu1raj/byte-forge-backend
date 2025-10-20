package api

import (
	"net/http"
)

// RegisterRoutes registers the HTTP routes for the server.
func (srv *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/submit", srv.SubmitHandler)
	mux.HandleFunc("/result/", srv.ResultHandler)
}
