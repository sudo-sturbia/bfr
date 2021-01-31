package frontend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/v2/pkg/books"
)

// searchForm serves the search form.
func (s *Server) searchForm(w http.ResponseWriter, r *http.Request) {
	s.tmpls[searchTmpl].Execute(w, nil)
}

// searchResults serves the search results acquired from search form.
func (s *Server) searchResults(w http.ResponseWriter, r *http.Request) {
	books, err := results(s.apiURL, r.URL.RawQuery)
	if err != nil {
		s.serveError(w, r, err)
	} else {
		s.tmpls[resultsTmpl].Execute(w, books)
	}
}

// serveBook serves a book based on an id.
func (s *Server) serveBook(w http.ResponseWriter, r *http.Request) {
	book, err := book(s.apiURL, mux.Vars(r)["id"])
	if err != nil {
		s.serveError(w, r, err)
	} else {
		s.tmpls[bookTmpl].Execute(w, book)
	}
}

// serveError serves a static error page.
func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithFields(
		log.Fields{
			"Address": r.RemoteAddr,
			"Method":  r.Method,
			"URL":     r.URL.String(),
		},
	).Info(err.Error())
	s.tmpls[errorTmpl].Execute(w, nil)
}

// results makes a request to the given api url and returns the response
// as a book slice, and an error.
func results(apiURL, query string) ([]*books.Book, error) {
	resp, err := http.Get(fmt.Sprintf("%s/books?%s", apiURL, query))
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %s", err.Error())
	}

	var books []*books.Book
	if err = json.Unmarshal(body, &books); err != nil {
		return nil, fmt.Errorf("invalid API reponse: %s", err.Error())
	}
	return books, nil
}

// book makes a request to the given api url and returns the response
// as a book, and an error.
func book(apiURL, id string) (*books.Book, error) {
	resp, err := http.Get(fmt.Sprintf("%s/book/%s", apiURL, id))
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %s", err.Error())
	}

	var book *books.Book
	if err = json.Unmarshal(body, &book); err != nil {
		return nil, fmt.Errorf("invalid API reponse: %s", err.Error())
	}
	return book, nil
}
