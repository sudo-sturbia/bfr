package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/pkg/books"
)

// A decoder to use for query parameters.
var (
	decoder = schema.NewDecoder()
)

// messageResponse is a string response meant to be used as an error
// message wrapper to be sent as a JSON object if a request fails to
// execute.
type messageResponse struct {
	Message string
}

// searchByTitle is a handler for /books/:title endpoint.
func (s *Server) searchByTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	books, err := books.SearchByTitle(
		&books.SearchIn{
			Datastore: s.searchIn.Datastore,
			BookTable: s.searchIn.BookTable,
		},
		vars["title"],
	)

	if err != nil {
		writeError("Search failed.", w, r)
	} else {
		writeResponse(books, w, r)
	}
}

// search is a handler for /books endpoint.
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	searchBy := &books.SearchBy{
		TitleHas:          "",
		Authors:           nil,
		LanguageCode:      nil,
		ISBN:              "",
		ISBN13:            "",
		RatingCeil:        -1,
		RatingFloor:       -1,
		PagesCeil:         -1,
		PagesFloor:        -1,
		RatingsCountCeil:  -1,
		RatingsCountFloor: -1,
		ReviewsCountCeil:  -1,
		ReviewsCountFloor: -1,
	}

	titlesOnly, parameters := isTitlesOnly(r.URL.Query())

	err := decoder.Decode(searchBy, parameters)
	if err != nil {
		writeError(err.Error(), w, r)
	} else {
		searchIn := &books.SearchIn{
			Datastore: s.searchIn.Datastore,
			BookTable: s.searchIn.BookTable,
		}

		if titlesOnly {
			titles, err := books.SearchForTitles(searchIn, searchBy)
			if err != nil {
				writeError("Search failed.", w, r)
			} else {
				writeResponse(titles, w, r)
			}
		} else {
			books, err := books.Search(searchIn, searchBy)
			if err != nil {
				writeError("Search failed.", w, r)
			} else {
				writeResponse(books, w, r)
			}
		}
	}
}

// writeResponse writes a JSON response to a request.
func writeResponse(toWrite interface{}, w http.ResponseWriter, r *http.Request) {
	response, err := json.MarshalIndent(toWrite, "", "\t")
	if err != nil {
		writeError("Search failed.", w, r)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	}
}

// writeError responds to a request with an error message using
// the given writer, and logs the error.
func writeError(message string, w http.ResponseWriter, r *http.Request) {
	log.WithFields(
		log.Fields{
			"Address": r.RemoteAddr,
			"Method":  r.Method,
			"URL":     r.URL.String(),
		},
	).Info(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response, err := json.MarshalIndent(&messageResponse{message}, "", "\t")
	if err == nil { // Just for safety as there are no cyclic structures.
		fmt.Fprint(w, string(response))
	}
}

// isTitlesOnly checks query parameters to see if TitlesOnly was specified. If
// so returns true and the query parameters with "TitlesOnly" key removed, else
// returns false and parameters without change.
func isTitlesOnly(query url.Values) (bool, url.Values) {
	titlesOnly := query.Get("TitlesOnly") == "true" || query.Get("TitlesOnly") == "True"
	if titlesOnly {
		query.Del("TitlesOnly")
	}

	return titlesOnly, query
}
