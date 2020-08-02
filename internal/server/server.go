// Package server is used to create a new, initialized instance of
// a bfr server ready to run.
package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents bfr's http server, and holds all dependencies needed
// for a server to run.
type Server struct {
	config *Config     // Server's configuration options.
	router *mux.Router // Server's router.

	searchIn *SearchIn // Datastore to search in.
}

// SearchIn Contains the database and table's name to search in. It is a
// copy of package books's SearchIn struct. It is created to importing
// books package into main.
type SearchIn struct {
	Datastore *sql.DB // Datastore to search in.
	BookTable string  // Table to search in.
}

// Config holds server's configuration options.
type Config struct {
	Host string // Host to run the server on.
	Port string // Port to run the server on.
}

// New creates and returns a new, initialized instance of bfr's server
// to run.
func New(config *Config, searchIn *SearchIn) *Server {
	s := &Server{
		config:   config,
		searchIn: searchIn,
		router:   mux.NewRouter(),
	}

	s.router.HandleFunc("/books/{name}", s.searchByTitle).Methods("GET")
	s.router.HandleFunc("/books", s.search).Methods("GET")

	return s
}

// Run runs a server instance. s listens on the port specified by
// s.Config.Port.
func (s *Server) Run() {
	http.ListenAndServe(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port), s.router)
}
