package books

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sudo-sturbia/bfr/internal/datastore"
)

// Test searching for books using IDs.
func TestSearchByID(t *testing.T) {
	searchIn, deferFn := testingSearchIn(t)
	defer deferFn()

	for id, book := range map[int]*Book{
		4: &Book{
			ID:            4,
			Title:         "Harry Potter and the Chamber of Secrets (Harry Potter  #2)",
			Authors:       "J.K. Rowling",
			AverageRating: 4.41,
			ISBN:          "439554896",
			ISBN13:        "9780439554893",
			LanguageCode:  "eng",
			Pages:         352,
			RatingsCount:  6267,
			ReviewsCount:  272,
		},
		14: &Book{
			ID:            14,
			Title:         "The Hitchhiker's Guide to the Galaxy (Hitchhiker's Guide to the Galaxy  #1)",
			Authors:       "Douglas Adams",
			AverageRating: 4.22,
			ISBN:          "1400052920",
			ISBN13:        "9781400052929",
			LanguageCode:  "eng",
			Pages:         215,
			RatingsCount:  4416,
			ReviewsCount:  408,
		},
		30: nil,
	} {
		result, err := SearchByID(searchIn, id)
		if err != nil && result != nil {
			t.Errorf("search failed: %s", err.Error())
		} else if err == nil && result == nil {
			t.Errorf("expected error, got: %v", err)
		} else {
			if result != book && *result != *book {
				t.Errorf("expected: %v, got: %v", book, result)
			}
		}
	}
}

// Test searching for books by title.
func TestSearchByTitle(t *testing.T) {
	searchIn, deferFn := testingSearchIn(t)
	defer deferFn()

	for name, book := range map[string]*Book{
		"Harry Potter and the Chamber of Secrets (Harry Potter  #2)": &Book{
			ID:            4,
			Title:         "Harry Potter and the Chamber of Secrets (Harry Potter  #2)",
			Authors:       "J.K. Rowling",
			AverageRating: 4.41,
			ISBN:          "439554896",
			ISBN13:        "9780439554893",
			LanguageCode:  "eng",
			Pages:         352,
			RatingsCount:  6267,
			ReviewsCount:  272,
		},

		"A Short History of Nearly Everything": &Book{
			ID:            21,
			Title:         "A Short History of Nearly Everything",
			Authors:       "Bill Bryson-William Roberts",
			AverageRating: 4.2,
			ISBN:          "076790818X",
			ISBN13:        "9780767908184",
			LanguageCode:  "eng",
			Pages:         544,
			RatingsCount:  228522,
			ReviewsCount:  8840,
		},
	} {
		result, err := SearchByTitle(searchIn, name)
		if err != nil {
			t.Errorf("search failed: %s", err.Error())
		} else {
			if len(result) != 1 {
				t.Errorf("expected: %d search result, got: %d.", 1, len(result))
			} else if *result[0] != *book {
				t.Errorf("incorrect search result for %s.", name)
			}
		}
	}
}

// Test searching using a SearchBy.
func TestSearch(t *testing.T) {
	searchIn, deferFn := testingSearchIn(t)
	defer deferFn()

	for searchBy, book := range map[*SearchBy][]*Book{
		&SearchBy{
			TitleHas:          "Secrets",
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
		}: []*Book{
			&Book{
				ID:            4,
				Title:         "Harry Potter and the Chamber of Secrets (Harry Potter  #2)",
				Authors:       "J.K. Rowling",
				AverageRating: 4.41,
				ISBN:          "439554896",
				ISBN13:        "9780439554893",
				LanguageCode:  "eng",
				Pages:         352,
				RatingsCount:  6267,
				ReviewsCount:  272,
			},
		},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       4.7,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []*Book{
			&Book{
				ID:            8,
				Title:         "Harry Potter Boxed Set  Books 1-5 (Harry Potter  #1-5)",
				Authors:       "J.K. Rowling-Mary GrandPré",
				AverageRating: 4.78,
				ISBN:          "439682584",
				ISBN13:        "9780439682589",
				LanguageCode:  "eng",
				Pages:         2690,
				RatingsCount:  38872,
				ReviewsCount:  154,
			},

			&Book{
				ID:            10,
				Title:         "Harry Potter Collection (Harry Potter  #1-6)",
				Authors:       "J.K. Rowling",
				AverageRating: 4.73,
				ISBN:          "439827604",
				ISBN13:        "9780439827607",
				LanguageCode:  "eng",
				Pages:         3342,
				RatingsCount:  27410,
				ReviewsCount:  820,
			},
		},

		&SearchBy{
			TitleHas:          "",
			Authors:           []string{"Bill", "Adams"},
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         400,
			PagesFloor:        200,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: 2000,
		}: []*Book{
			&Book{
				ID:            24,
				Title:         "In a Sunburned Country",
				Authors:       "Bill Bryson",
				AverageRating: 4.07,
				ISBN:          "767903862",
				ISBN13:        "9780767903868",
				LanguageCode:  "eng",
				Pages:         335,
				RatingsCount:  68213,
				ReviewsCount:  4077,
			},

			&Book{
				ID:            25,
				Title:         "I'm a Stranger Here Myself: Notes on Returning to America After Twenty Years Away",
				Authors:       "Bill Bryson",
				AverageRating: 3.9,
				ISBN:          "076790382X",
				ISBN13:        "9780767903820",
				LanguageCode:  "eng",
				Pages:         304,
				RatingsCount:  47490,
				ReviewsCount:  2153,
			},
		},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      []string{"fre"},
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
		}: []*Book{},
	} {
		result, err := Search(searchIn, searchBy)
		if err != nil {
			t.Errorf("search failed: %s", err.Error())
		} else {
			if len(result) != len(book) {
				t.Errorf("expected: %d search result, got: %d", 1, len(result))
			} else {
				for i, res := range result {
					if *res != *book[i] {
						t.Errorf("incorrect search result for %s", book[i].Title)
					}
				}
			}
		}
	}
}

// Test searching for books' titles using SearchBy.
func TestSearchForTitles(t *testing.T) {
	searchIn, deferFn := testingSearchIn(t)
	defer deferFn()

	for searchBy, title := range map[*SearchBy][]string{
		&SearchBy{
			TitleHas:          "Secrets",
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
		}: []string{
			"Harry Potter and the Chamber of Secrets (Harry Potter  #2)",
		},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       4.7,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []string{
			"Harry Potter Boxed Set  Books 1-5 (Harry Potter  #1-5)",
			"Harry Potter Collection (Harry Potter  #1-6)",
		},

		&SearchBy{
			TitleHas:          "",
			Authors:           []string{"Bill", "Adams"},
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         400,
			PagesFloor:        200,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: 2000,
		}: []string{
			"In a Sunburned Country",
			"I'm a Stranger Here Myself: Notes on Returning to America After Twenty Years Away",
		},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      []string{"fre"},
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
		}: []string{},
	} {
		result, err := SearchForTitles(searchIn, searchBy)
		if err != nil {
			t.Errorf("search failed: %s", err.Error())
		} else {
			if len(result) != len(title) {
				t.Errorf("expected: %d search result, got: %d", 1, len(result))
			} else {
				for i, res := range result {
					if res != title[i] {
						t.Errorf("incorrect search result for %s", title[i])
					}
				}
			}
		}
	}
}

// testingSearchIn returns a SearchIn to use for tests, and a func
// that frees resources, and should be defered.
func testingSearchIn(t *testing.T) (*SearchIn, func()) {
	t.Helper()
	config := &datastore.Config{
		Driver:    "sqlite3",
		Datastore: "testDatastore.db",
		BookTable: "books",
	}

	err := datastore.New("../../test-data/booksTest.csv", config, true)
	if err != nil {
		t.Fatalf("couldn't load datastore: %s.", err.Error())
	}

	datastore, err := sql.Open(config.Driver, fmt.Sprintf("file:%s", config.Datastore))
	if err != nil {
		t.Fatalf("failed to open database: %s", err.Error())
	}

	deferFn := func() {
		datastore.Close()
		os.Remove(config.Datastore)
	}

	return &SearchIn{
		Datastore: datastore,
		BookTable: config.BookTable,
	}, deferFn
}
