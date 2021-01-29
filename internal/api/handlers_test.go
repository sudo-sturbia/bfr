package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sudo-sturbia/bfr/internal/testhelper"
)

// TestSearchByID tests searching for a book using id.
func TestSearchByID(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	server := New(nil, &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	})

	for _, test := range []struct {
		id       int
		response string
		status   int
	}{
		{
			id: -1,
			response: fmt.Sprint(
				"{\n",
				"\t\"Message\": \"Search failed.\"\n",
				"}",
			),
			status: 400,
		},
		{
			id: 40,
			response: fmt.Sprint(
				"{\n",
				"\t\"Message\": \"Search failed.\"\n",
				"}",
			),
			status: 400,
		},
		{
			id: 1,
			response: fmt.Sprint(
				"{\n",
				"\t\"ID\": 1,\n",
				"\t\"Title\": \"Harry Potter and the Half-Blood Prince (Harry Potter  #6)\",\n",
				"\t\"Authors\": \"J.K. Rowling-Mary GrandPré\",\n",
				"\t\"AverageRating\": 4.56,\n",
				"\t\"ISBN\": \"439785960\",\n",
				"\t\"ISBN13\": \"9780439785969\",\n",
				"\t\"LanguageCode\": \"eng\",\n",
				"\t\"Pages\": 652,\n",
				"\t\"RatingsCount\": 1944099,\n",
				"\t\"ReviewsCount\": 26249\n",
				"}",
			),
			status: 200,
		},
	} {
		t.Run(
			fmt.Sprintf("id:%d", test.id),
			func(*testing.T) {
				recorder := recordResponse(t, fmt.Sprintf("/book/%d", test.id), "/book/{id}", server.searchByID)
				if recorder.Code != test.status {
					t.Errorf("incorrect status, want: %d, got: %d", test.status, recorder.Code)
				}
				if recorder.Body.String() != test.response {
					t.Errorf("incorrect response, want: %s, got: %s", test.response, recorder.Body.String())
				}
			},
		)
	}
}

// TestSearchByTitle tests searching for a book using title.
func TestSearchByTitle(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	server := New(nil, &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	})

	for _, test := range []struct {
		title    string
		response string
		status   int
	}{
		{
			title:    "NoTitle",
			response: "[]",
			status:   200,
		},
		{
			title:    "1",
			response: "[]",
			status:   200,
		},
		{
			title: "The Adventures of Sherlock Holmes",
			response: fmt.Sprint(
				"[\n",
				"\t{\n",
				"\t\t\"ID\": 3588,\n",
				"\t\t\"Title\": \"The Adventures of Sherlock Holmes\",\n",
				"\t\t\"Authors\": \"Arthur Conan Doyle-Eoin Colfer\",\n",
				"\t\t\"AverageRating\": 4.31,\n",
				"\t\t\"ISBN\": \"439574285\",\n",
				"\t\t\"ISBN13\": \"9780439574280\",\n",
				"\t\t\"LanguageCode\": \"eng\",\n",
				"\t\t\"Pages\": 336,\n",
				"\t\t\"RatingsCount\": 811,\n",
				"\t\t\"ReviewsCount\": 86\n",
				"\t}\n",
				"]",
			),
			status: 200,
		},
	} {
		t.Run(
			fmt.Sprintf("title:%s", test.title),
			func(*testing.T) {
				recorder := recordResponse(t, fmt.Sprintf("/books/%s", test.title), "/books/{title}", server.searchByTitle)
				if recorder.Code != test.status {
					t.Errorf("incorrect status, want: %d, got: %d", test.status, recorder.Code)
				}
				if recorder.Body.String() != test.response {
					t.Errorf("incorrect response, want: %s, got: %s", test.response, recorder.Body.String())
				}
			},
		)
	}
}

// TestSearch tests searching for a book using a set of parameters.
func TestSearch(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	server := New(nil, &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	})

	for _, test := range []struct {
		queryParams string
		response    string
		status      int
	}{
		{
			queryParams: "Wrong=10",
			response: fmt.Sprint(
				"{\n",
				"\t\"Message\": \"schema: invalid path \\\"Wrong\\\"\"\n",
				"}",
			),
			status: 400,
		},
		{
			queryParams: "Authors=Arthur&RatingFloor=4.3",
			response: fmt.Sprint(
				"[\n",
				"\t{\n",
				"\t\t\"ID\": 3588,\n",
				"\t\t\"Title\": \"The Adventures of Sherlock Holmes\",\n",
				"\t\t\"Authors\": \"Arthur Conan Doyle-Eoin Colfer\",\n",
				"\t\t\"AverageRating\": 4.31,\n",
				"\t\t\"ISBN\": \"439574285\",\n",
				"\t\t\"ISBN13\": \"9780439574280\",\n",
				"\t\t\"LanguageCode\": \"eng\",\n",
				"\t\t\"Pages\": 336,\n",
				"\t\t\"RatingsCount\": 811,\n",
				"\t\t\"ReviewsCount\": 86\n",
				"\t}\n",
				"]",
			),
			status: 200,
		},
	} {
		t.Run(
			fmt.Sprintf("title:%s", test.queryParams),
			func(*testing.T) {
				recorder := recordResponse(t, fmt.Sprintf("/books?%s", test.queryParams), "/books", server.search)
				if recorder.Code != test.status {
					t.Errorf("incorrect status, want: %d, got: %d", test.status, recorder.Code)
				}
				if recorder.Body.String() != test.response {
					t.Errorf("incorrect response, want: %s, got: %s", test.response, recorder.Body.String())
				}
			},
		)
	}
}

// recordResponse performs a test request with a ResponseRecoder and returns the
// recoder.
func recordResponse(t *testing.T, url, muxURL string, helper http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("failed to create request: %s", err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc(muxURL, helper).Methods("GET")
	router.ServeHTTP(recorder, request)
	return recorder
}
