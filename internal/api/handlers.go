package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/pkg/books"
)

// A decoder to use for query parameters.
var (
	decoder = schema.NewDecoder()
)

// errorResponse is an error message wrapper meant to be used as a JSON response
// to requests in case of errors.
type errorResponse struct {
	Message string
}

// searchByID is a handler for /book/{id} endpoint.
func (s *Server) searchByID(w http.ResponseWriter, r *http.Request) {
	response, status := searchByIDResponse(s.searchIn, mux.Vars(r)["id"])
	write(w, r, response, status)
}

// searchByTitle is a handler for /books/{title} endpoint.
func (s *Server) searchByTitle(w http.ResponseWriter, r *http.Request) {
	response, status := searchByTitleResponse(s.searchIn, mux.Vars(r)["title"])
	write(w, r, response, status)
}

// search is a handler for /books endpoint.
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	response, status := searchResponse(
		r.URL.Query(),
		s.searchIn,
		&books.SearchBy{
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
		},
	)
	write(w, r, response, status)
}

// searchByIDResponse searchs the database for books based on given parameters and
// returns a response and a status code. It should be used by Server.searchByID.
func searchByIDResponse(searchIn *SearchIn, idString string) (interface{}, int) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return &errorResponse{fmt.Sprintf("Invalid id \"%s\".", idString)}, http.StatusBadRequest
	}

	book, err := books.SearchByID(
		&books.SearchIn{
			Datastore: searchIn.Datastore,
			BookTable: searchIn.BookTable,
		},
		id,
	)
	if err != nil {
		return &errorResponse{"Search failed."}, http.StatusBadRequest
	}
	return book, http.StatusOK
}

// searchByTitleResponse searchs the database for books based on given parameters and
// returns a response and a status code. It should be used by Server.searchByTitle.
func searchByTitleResponse(searchIn *SearchIn, title string) (interface{}, int) {
	books, err := books.SearchByTitle(
		&books.SearchIn{
			Datastore: searchIn.Datastore,
			BookTable: searchIn.BookTable,
		},
		title,
	)

	if err != nil {
		return &errorResponse{"Search failed."}, http.StatusBadRequest
	}
	return books, http.StatusOK
}

// searchResponse searchs the database for books based on given parameters and returns
// a response and a status code. It should be used by Server.search.
func searchResponse(query url.Values, searchIn *SearchIn, searchBy *books.SearchBy) (interface{}, int) {
	titlesOnly, parameters := isTitlesOnly(query)
	err := decoder.Decode(searchBy, parameters)
	if err != nil {
		return &errorResponse{err.Error()}, http.StatusBadRequest
	}

	in := &books.SearchIn{
		Datastore: searchIn.Datastore,
		BookTable: searchIn.BookTable,
	}

	if titlesOnly {
		titles, err := books.SearchForTitles(in, searchBy)
		if err != nil {
			return &errorResponse{"Search failed."}, http.StatusBadRequest
		}
		return titles, http.StatusOK
	}

	books, err := books.Search(in, searchBy)
	if err != nil {
		return &errorResponse{"Search failed."}, http.StatusBadRequest
	}
	return books, http.StatusOK
}

// write writes a JSON response to a request.
func write(w http.ResponseWriter, r *http.Request, response interface{}, status int) {
	if status >= 400 && status < 600 { // Error status.
		message, ok := response.(errorResponse)
		if ok {
			errorLog(r, message)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonRes, err := json.Marshal(response)
	if err == nil { // Just for safety as there are no cyclic structures.
		fmt.Fprint(w, string(jsonRes))
	}
}

// errorLog logs a request error.
func errorLog(r *http.Request, message errorResponse) {
	log.WithFields(
		log.Fields{
			"Address": r.RemoteAddr,
			"Method":  r.Method,
			"URL":     r.URL.String(),
		},
	).Info(message)
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
