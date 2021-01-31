package books

import (
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sudo-sturbia/bfr/v2/internal/testhelper"
)

// Test searching for books using IDs.
func TestSearchByID(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	searchIn := &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	}

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
		t.Run(
			fmt.Sprintf("id: %d", id),
			func(*testing.T) {
				result, err := SearchByID(searchIn, id)
				if err != nil && result != nil {
					t.Fatalf("search failed: %s", err.Error())
				}

				if err == nil && result == nil {
					t.Fatalf("expected error, got: %v", err)
				}

				if result != book && *result != *book {
					t.Fatalf("expected: %v, got: %v", book, result)
				}
			},
		)
	}
}

// Test searching for books by title.
func TestSearchByTitle(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	searchIn := &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	}

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
		t.Run(
			fmt.Sprintf("name: %s", name),
			func(*testing.T) {
				result, err := SearchByTitle(searchIn, name)
				if err != nil {
					t.Fatalf("search failed: %s", err.Error())
				}

				if len(result) != 1 {
					t.Fatalf("expected: %d search result, got: %d.", 1, len(result))
				}

				if *result[0] != *book {
					t.Fatalf("incorrect search result for %s.", name)
				}
			},
		)
	}
}

// Test searching using a SearchBy.
func TestSearch(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	searchIn := &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	}

	i := 0
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
		t.Run(
			fmt.Sprintf("test: %d", i),
			func(*testing.T) {
				result, err := Search(searchIn, searchBy)
				if err != nil {
					t.Fatalf("search failed: %s", err.Error())
				}

				if len(result) != len(book) {
					t.Fatalf("expected: %d search result, got: %d", 1, len(result))
				}

				for i, res := range result {
					if *res != *book[i] {
						t.Fatalf("incorrect search result for %s", book[i].Title)
					}
				}
			},
		)
		i++
	}
}

// Test searching for books' titles using SearchBy.
func TestSearchForTitles(t *testing.T) {
	datastore, bookTable, deferFn := testhelper.SearchIn(t)
	defer deferFn()

	searchIn := &SearchIn{
		Datastore: datastore,
		BookTable: bookTable,
	}

	i := 0
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
		t.Run(
			fmt.Sprintf("test: %d", i),
			func(*testing.T) {
				result, err := SearchForTitles(searchIn, searchBy)
				if err != nil {
					t.Fatalf("search failed: %s", err.Error())
				}

				if len(result) != len(title) {
					t.Fatalf("expected: %d search result, got: %d", 1, len(result))
				}

				for i, res := range result {
					if res != title[i] {
						t.Fatalf("incorrect search result for %s", title[i])
					}
				}
			},
		)
		i++
	}
}
