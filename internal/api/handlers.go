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
	"github.com/sudo-sturbia/bfr/v2/pkg/books"
)

// A decoder to use for query parameters.
var (
	decoder = schema.NewDecoder()
)

// searchByID is a handler for /book/{id} endpoint.
func (s *Server) searchByID(w http.ResponseWriter, r *http.Request) {
	response, status, ok := searchByIDResponse(s.searchIn, mux.Vars(r)["id"])
	if ok {
		write(w, r, response, status)
	} else {
		message, ok := response.(string)
		if ok {
			writeError(w, r, message, status)
		}
	}
}

// searchByTitle is a handler for /books/{title} endpoint.
func (s *Server) searchByTitle(w http.ResponseWriter, r *http.Request) {
	response, status, ok := searchByTitleResponse(s.searchIn, mux.Vars(r)["title"])
	if ok {
		write(w, r, response, status)
	} else {
		message, ok := response.(string)
		if ok {
			writeError(w, r, message, status)
		}
	}
}

// search is a handler for /books endpoint.
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	response, status, ok := searchResponse(
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
	if ok {
		write(w, r, response, status)
	} else {
		message, ok := response.(string)
		if ok {
			writeError(w, r, message, status)
		}
	}
}

// searchByIDResponse searchs the database for books based on given parameters and
// returns a response, a status code, and bool indicating if the operation was performed
// successfully. It should be used by Server.searchByID.
func searchByIDResponse(searchIn *SearchIn, idString string) (interface{}, int, bool) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return fmt.Sprintf("Invalid id \"%s\".", idString), http.StatusBadRequest, false
	}

	book, err := books.SearchByID(
		&books.SearchIn{
			Datastore: searchIn.Datastore,
			BookTable: searchIn.BookTable,
		},
		id,
	)
	if err != nil {
		return "Search failed.", http.StatusBadRequest, false
	}
	return book, http.StatusOK, true
}

// searchByTitleResponse searchs the database for books based on given parameters and
// returns a response, a status code, and bool indicating if the operation was performed
// successfully. It should be used by Server.searchByTitle.
func searchByTitleResponse(searchIn *SearchIn, title string) (interface{}, int, bool) {
	books, err := books.SearchByTitle(
		&books.SearchIn{
			Datastore: searchIn.Datastore,
			BookTable: searchIn.BookTable,
		},
		title,
	)

	if err != nil {
		return "Search failed.", http.StatusBadRequest, false
	}
	return books, http.StatusOK, true
}

// searchResponse searchs the database for books based on given parameters and
// returns a response, a status code, and bool indicating if the operation was performed
// successfully. It should be used by Server.search.
func searchResponse(query url.Values, searchIn *SearchIn, searchBy *books.SearchBy) (interface{}, int, bool) {
	titlesOnly, parameters := isTitlesOnly(query)
	err := decoder.Decode(searchBy, parameters)
	if err != nil {
		return "Unable to decode search query.", http.StatusBadRequest, false
	}

	in := &books.SearchIn{
		Datastore: searchIn.Datastore,
		BookTable: searchIn.BookTable,
	}

	if titlesOnly {
		titles, err := books.SearchForTitles(in, searchBy)
		if err != nil {
			return "Search failed.", http.StatusBadRequest, false
		}
		return titles, http.StatusOK, true
	}

	books, err := books.Search(in, searchBy)
	if err != nil {
		return "Search failed.", http.StatusBadRequest, false
	}
	return books, http.StatusOK, true
}

// write writes a JSON response to a request.
func write(w http.ResponseWriter, r *http.Request, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonRes, err := json.MarshalIndent(response, "", "\t")
	if err == nil { // Just for safety as there are no cyclic structures.
		fmt.Fprint(w, string(jsonRes))
	}
}

// writeError logs the error and writes it to request.
func writeError(w http.ResponseWriter, r *http.Request, message string, status int) {
	log.WithFields(
		log.Fields{
			"Address": r.RemoteAddr,
			"Method":  r.Method,
			"URL":     r.URL.String(),
		},
	).Info(message)
	http.Error(w, message, status)
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
