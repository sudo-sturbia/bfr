// Package books contains book-related logic, and provides functions to search
// for books in a datastore.
//
// books also serves to expose an API to be used by servers.
package books

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Book represents a searchable book object.
type Book struct {
	Id            int     // A different number for each book.
	Title         string  // Book's title.
	Authors       string  // Book's authors.
	AverageRating float32 // Average rating (out of 5.)
	ISBN          string  // 10 digit ISBN.
	ISBN13        string  // 13 digit ISBN.
	LanguageCode  string  // 3-character language code.
	Pages         int     // Number of book's pages.
	RatingsCount  int     // Number of ratings (out of 5.)
	ReviewsCount  int     // Number of text reviews.
}

// SearchIn contains the database and table's name to search in.
type SearchIn struct {
	Datastore *sql.DB // Datastore to search in.
	BookTable string  // Table to search in.
}

// SearchBy is a set of parameters to use when searching the datastore
// for books.
// Not all fields have to be specifed, a search can be performed using
// only a sub-set of the fields. To ignore a string field when searching,
// leave it empty, to ignore a number set it to < 0.
// For floor/ceil values floor is exclusive, ceil is inclusive.
type SearchBy struct {
	TitleHas string // A sub-string that must exist in the title.

	Authors      []string // Book's authors. Ignored if nil or empty.
	LanguageCode []string // Must be one of these languages. Ignored if nil or empty.

	ISBN   string // 10 digit ISBN.
	ISBN13 string // 13 digit ISBN.

	RatingCeil        float32 // Rating must be less than or equal.
	RatingFloor       float32 // Rating must be higher than.
	PagesCeil         int     // Number of pages must be less than or equal.
	PagesFloor        int     // Number of pages must be higher than.
	RatingsCountCeil  int     // Number of ratings must be less than or equal.
	RatingsCountFloor int     // Number of ratings must be higher than.
	ReviewsCountCeil  int     // Number of reviews must be less than or equal.
	ReviewsCountFloor int     // Number of reviews must be higher than.
}

// SearchByTitle searchs for books with a specific given title in the
// given datastore and table. Returns an empty list if non is found.
func SearchByTitle(searchIn *SearchIn, title string) ([]*Book, error) {
	books := make([]*Book, 0)

	search := fmt.Sprintf("select * from %s where title = ?;", searchIn.BookTable)
	rows, err := searchIn.Datastore.Query(search, title)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		book := new(Book)
		rows.Scan(
			&book.Id,
			&book.Title,
			&book.Authors,
			&book.AverageRating,
			&book.ISBN,
			&book.ISBN13,
			&book.LanguageCode,
			&book.Pages,
			&book.RatingsCount,
			&book.ReviewsCount,
		)

		books = append(books, book)
	}

	return books, nil
}

// Search searchs the specified table in the given database using SearchBy's
// fields, and returns a list of books that match the given parameters.
func Search(searchIn *SearchIn, searchBy *SearchBy) ([]*Book, error) {
	return nil, nil
}
