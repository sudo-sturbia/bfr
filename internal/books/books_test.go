package books

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sudo-sturbia/bfr/internal/datastore"
)

// Test searching for books by title.
func TestSearchByTitle(t *testing.T) {
	config := &datastore.Config{
		Driver:    "sqlite3",
		Datastore: "testDatastore.db",
		BookTable: "books",
	}

	err := datastore.New("../../test-data/booksTest.csv", config, true)
	if err != nil {
		t.Errorf("Couldn't load datastore: %s.", err.Error())
		return
	}

	datastore, err := sql.Open(config.Driver, fmt.Sprintf("file:%s", config.Datastore))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer datastore.Close()
	defer os.Remove(config.Datastore)

	searchIn := &SearchIn{
		Datastore: datastore,
		BookTable: config.BookTable,
	}

	books := map[string]*Book{
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
	}

	for name, book := range books {
		result, err := SearchByTitle(searchIn, name)
		if err != nil {
			t.Errorf("Search failed.")
			continue
		}

		if len(result) != 1 {
			t.Errorf("Expected %d search result, Found %d.", 1, len(result))
			continue
		}

		if *result[0] != *book {
			t.Errorf("Incorrect search result for %s.", name)
		}
	}
}

// Test searching using a SearchBy.
func TestSearch(t *testing.T) {
	config := &datastore.Config{
		Driver:    "sqlite3",
		Datastore: "testDatastore.db",
		BookTable: "books",
	}

	err := datastore.New("../../test-data/booksTest.csv", config, true)
	if err != nil {
		t.Errorf("Couldn't load datastore: %s.", err.Error())
		return
	}

	datastore, err := sql.Open(config.Driver, fmt.Sprintf("file:%s", config.Datastore))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer datastore.Close()
	defer os.Remove(config.Datastore)

	searchIn := &SearchIn{
		Datastore: datastore,
		BookTable: config.BookTable,
	}

	books := map[*SearchBy][]*Book{
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
	}

	for searchBy, book := range books {
		result, err := Search(searchIn, searchBy)
		if err != nil {
			t.Errorf("Search failed.")
			continue
		}

		if len(result) != len(book) {
			t.Errorf("Expected %d search result, Found %d.", 1, len(result))
			continue
		}

		for i, res := range result {
			if *res != *book[i] {
				t.Errorf("Incorrect search result for %s.", book[i].Title)
			}
		}
	}
}

// Test query generation based on SearchBy.
func TestQuery(t *testing.T) {
	searchIn := &SearchIn{
		Datastore: nil,
		BookTable: "books",
	}

	searchQueries := map[*SearchBy]string{
		&SearchBy{
			TitleHas:          "aaa",
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
		}: "select * from books where title like ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           []string{"a", "b", "c"},
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
		}: "select * from books where (authors like ? or authors like ? or authors like ?);",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      []string{"a", "b", "c"},
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
		}: "select * from books where (languageCode like ? or languageCode like ? or languageCode like ?);",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "123456789",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where isbn = ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "123456789abc",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where isbn13 = ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        3,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where averageRating <= ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where averageRating > ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         100,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where pages <= ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        50,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where pages > ?;",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  100,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where ratingsCount <= ?;",

		&SearchBy{
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
			RatingsCountFloor: 50,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where ratingsCount > ?;",

		&SearchBy{
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
			ReviewsCountCeil:  100,
			ReviewsCountFloor: -1,
		}: "select * from books where reviewsCount <= ?;",

		&SearchBy{
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
			ReviewsCountFloor: 50,
		}: "select * from books where reviewsCount > ?;",

		&SearchBy{
			TitleHas:          "aaa",
			Authors:           nil,
			LanguageCode:      []string{"a", "b"},
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
		}: "select * from books where title like ? and (languageCode like ? or languageCode like ?);",

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "0123456789abc",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        150,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  200,
			ReviewsCountFloor: -1,
		}: "select * from books where isbn13 = ? and pages > ? and reviewsCount <= ?;",

		&SearchBy{
			TitleHas:          "aaa",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "0123456789",
			ISBN13:            "",
			RatingCeil:        4.8,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  500,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: "select * from books where title like ? and isbn = ? and averageRating <= ? and ratingsCount <= ?;",

		&SearchBy{
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
		}: "select * from books;",

		&SearchBy{
			TitleHas:          "aaa",
			Authors:           []string{"a", "b", "c"},
			LanguageCode:      []string{"a", "b"},
			ISBN:              "0123456789",
			ISBN13:            "0123456789abc",
			RatingCeil:        4.5,
			RatingFloor:       2,
			PagesCeil:         200,
			PagesFloor:        20,
			RatingsCountCeil:  1000,
			RatingsCountFloor: 500,
			ReviewsCountCeil:  1000,
			ReviewsCountFloor: 500,
		}: "select * from books where title like ? and (authors like ? or authors like ? or authors like ?) and (languageCode like ? or languageCode like ?) and " +
			"isbn = ? and isbn13 = ? and " +
			"averageRating <= ? and averageRating > ? and " +
			"pages <= ? and pages > ? and " +
			"ratingsCount <= ? and ratingsCount > ? and " +
			"reviewsCount <= ? and reviewsCount > ?;",
	}

	for s, sq := range searchQueries {
		if q, _ := query(searchIn, s); q != sq {
			t.Errorf("Expected \"%s\", Found \"%s\".", sq, q)
		}
	}
}

// Test parameters of queries generated based on SearchBy.
func TestQueryParameters(t *testing.T) {
	searchIn := &SearchIn{
		Datastore: nil,
		BookTable: "books",
	}

	searchQueries := map[*SearchBy][]interface{}{
		&SearchBy{
			TitleHas:          "aaa",
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
		}: []interface{}{"%aaa%"},

		&SearchBy{
			TitleHas:          "",
			Authors:           []string{"a", "b", "c"},
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
		}: []interface{}{"%a%", "%b%", "%c%"},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      []string{"a", "b", "c"},
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
		}: []interface{}{"%a%", "%b%", "%c%"},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "123456789",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{"123456789"},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "123456789abc",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{"123456789abc"},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        3,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{float32(3)},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{float32(1)},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         100,
			PagesFloor:        -1,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{100},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        50,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{50},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  100,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{100},

		&SearchBy{
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
			RatingsCountFloor: 50,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{50},

		&SearchBy{
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
			ReviewsCountCeil:  100,
			ReviewsCountFloor: -1,
		}: []interface{}{100},

		&SearchBy{
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
			ReviewsCountFloor: 50,
		}: []interface{}{50},

		&SearchBy{
			TitleHas:          "aaa",
			Authors:           nil,
			LanguageCode:      []string{"a", "b"},
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
		}: []interface{}{"%aaa%", "%a%", "%b%"},

		&SearchBy{
			TitleHas:          "",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "",
			ISBN13:            "0123456789abc",
			RatingCeil:        -1,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        150,
			RatingsCountCeil:  -1,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  200,
			ReviewsCountFloor: -1,
		}: []interface{}{"0123456789abc", 150, 200},

		&SearchBy{
			TitleHas:          "aaa",
			Authors:           nil,
			LanguageCode:      nil,
			ISBN:              "0123456789",
			ISBN13:            "",
			RatingCeil:        4.8,
			RatingFloor:       -1,
			PagesCeil:         -1,
			PagesFloor:        -1,
			RatingsCountCeil:  500,
			RatingsCountFloor: -1,
			ReviewsCountCeil:  -1,
			ReviewsCountFloor: -1,
		}: []interface{}{"%aaa%", "0123456789", float32(4.8), 500},

		&SearchBy{
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
		}: []interface{}{},

		&SearchBy{
			TitleHas:          "aaa",
			Authors:           []string{"a", "b", "c"},
			LanguageCode:      []string{"a", "b"},
			ISBN:              "0123456789",
			ISBN13:            "0123456789abc",
			RatingCeil:        4.5,
			RatingFloor:       2,
			PagesCeil:         200,
			PagesFloor:        20,
			RatingsCountCeil:  1000,
			RatingsCountFloor: 500,
			ReviewsCountCeil:  1000,
			ReviewsCountFloor: 500,
		}: []interface{}{"%aaa%", "%a%", "%b%", "%c%", "%a%", "%b%", "0123456789", "0123456789abc", float32(4.5), float32(2), 200, 20, 1000, 500, 1000, 500},
	}

	for s, sp := range searchQueries {
		if _, p := query(searchIn, s); !compareSlices(p, sp) {
			t.Errorf("Expected \"%s\", Found \"%s\".", sp, p)
		}
	}
}

// compareSlices returns true if a and b are equal, false otherwise.
func compareSlices(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
