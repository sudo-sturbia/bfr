package frontend

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	searchTmpl  = "search"
	resultsTmpl = "results"
	bookTmpl    = "book"
)

type Server struct {
	cfg    *Config
	router *mux.Router
	apiURL string                        // URL of API to make requests to.
	tmpls  map[string]*template.Template // Map of HTML templates with names.
}

// Config holds server's configuration options.
type Config struct {
	Host   string // Host to run the server on.
	Port   string // Port to run the server on.
	Static string // Path to static files.
}

// New returns a new, initialized frontend server.
func New(cfg *Config, api string) (_ *Server, err error) {
	s := &Server{
		cfg:    cfg,
		router: mux.NewRouter(),
		apiURL: api,
	}

	s.tmpls, err = newTemplates(cfg)
	if err != nil {
		return nil, fmt.Errorf("New: %s", err.Error())
	}

	s.router.HandleFunc("/", s.searchForm).Methods("GET")
	s.router.HandleFunc("/search", s.searchResults).Methods("GET")
	s.router.HandleFunc("/book/{id}", s.serveBook).Methods("GET")
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.cfg.Static))))
	return s, nil
}

// Run runs a server instance on the host and port specified in its config.
func (s *Server) Run() {
	http.ListenAndServe(fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port), s.router)
}

// newTemplates parses templates and returns a new template map.
func newTemplates(cfg *Config) (_ map[string]*template.Template, err error) {
	base := fmt.Sprintf("%s/base.html", cfg.Static)

	tmpls := make(map[string]*template.Template)
	for _, name := range []string{searchTmpl, resultsTmpl, bookTmpl} {
		tmpls[name], err = template.ParseFiles(base, fmt.Sprintf("%s/%s.html", cfg.Static, name))
		if err != nil {
			return nil, err
		}
	}
	return tmpls, nil
}
