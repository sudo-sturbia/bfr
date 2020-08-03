package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/internal/books"
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

	books, err := books.SearchByTitle(&books.SearchIn{
		Datastore: s.searchIn.Datastore,
		BookTable: s.searchIn.BookTable,
	}, vars["title"])

	if err != nil {
		errorResponse("Search failed.", w, r)
	} else {
		response, err := json.MarshalIndent(books, "", "\t")
		if err != nil {
			errorResponse("Search failed.", w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, string(response))
		}
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

	err := decoder.Decode(searchBy, r.URL.Query())
	if err != nil {
		errorResponse("Invalid request url.", w, r)
	} else {
		books, err := books.Search(&books.SearchIn{
			Datastore: s.searchIn.Datastore,
			BookTable: s.searchIn.BookTable,
		}, searchBy)

		if err != nil {
			errorResponse("Search failed.", w, r)
		} else {
			response, err := json.MarshalIndent(books, "", "\t")
			if err != nil {
				errorResponse("Search failed.", w, r)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, string(response))
			}
		}
	}
}

// errorResponse responds to a request with an error message using
// the given writer, and logs the error.
func errorResponse(message string, w http.ResponseWriter, r *http.Request) {
	log.WithFields(
		log.Fields{
			"IP":     r.RemoteAddr,
			"Method": r.Method,
			"URL":    r.URL.String(),
		},
	).Info(message)

	response, _ := json.MarshalIndent(&messageResponse{message}, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	fmt.Fprint(w, string(response))
}
