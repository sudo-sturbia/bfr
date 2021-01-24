// Package api is used to create, and run a new, initialized
// backend server instance.
package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents bfr's http server, and holds all dependencies needed
// for a server to run.
type Server struct {
	cfg    *Config     // Server's configuration options.
	router *mux.Router // Server's router.

	searchIn *SearchIn // Datastore to search in.
}

// Config holds server's configuration options.
type Config struct {
	Host string // Host to run the server on.
	Port string // Port to run the server on.
}

// SearchIn Contains a database and name of the table to search in. It is a
// copy of package books's SearchIn struct created to prevent importing of
// books package into main to limit dependency.
type SearchIn struct {
	Datastore *sql.DB // Datastore to search in.
	BookTable string  // Table to search in.
}

// New creates and returns a new, initialized server instance with handlers
// pointing to correct routes.
func New(cfg *Config, searchIn *SearchIn) *Server {
	s := &Server{
		cfg:      cfg,
		searchIn: searchIn,
		router:   mux.NewRouter(),
	}

	s.router.HandleFunc("/books/{title}", s.searchByTitle).Methods("GET")
	s.router.HandleFunc("/books", s.search).Methods("GET")

	return s
}

// Run runs a server instance on the host and port specified in its config.
func (s *Server) Run() {
	http.ListenAndServe(fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port), s.router)
}
