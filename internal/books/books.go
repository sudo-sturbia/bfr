// Package books contains book-related logic, and provides functions to search
// for books in a datastore.
// books also serves to expose an API to be used by servers.
package books

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Used with sql package.
)

// Book represents a searchable book object.
type Book struct {
	ID            int     // A different number for each book.
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

// SearchBy is a set of parameters to use when searching for books in
// a datastore.
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
	search := fmt.Sprintf("select * from %s where title = ?;", searchIn.BookTable)
	rows, err := searchIn.Datastore.Query(search, title)
	if err != nil {
		return nil, err
	}

	books := make([]*Book, 0)
	for rows.Next() {
		book := new(Book)
		rows.Scan(
			&book.ID,
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
	query, parameters := query(searchIn, searchBy)
	rows, err := searchIn.Datastore.Query(query, parameters...)
	if err != nil {
		return nil, err
	}

	books := make([]*Book, 0)
	for rows.Next() {
		book := new(Book)
		rows.Scan(
			&book.ID,
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

// query generates a SQL search query based on fields specified in SearchBy.
// Returns a sql query, and a list of parameters to use with the prepared
// statement when executing.
func query(searchIn *SearchIn, searchBy *SearchBy) (string, []interface{}) {
	queryFields := make([]string, 0)
	fields := make([]interface{}, 0)

	if ok, q, f := titleHas(searchBy.TitleHas); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := authors(searchBy.Authors); ok {
		queryFields, fields = append(queryFields, q), append(fields, f...)
	}

	if ok, q, f := languageCode(searchBy.LanguageCode); ok {
		queryFields, fields = append(queryFields, q), append(fields, f...)
	}

	if ok, q, f := isbn(searchBy.ISBN); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := isbn13(searchBy.ISBN13); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := ratingCeil(searchBy.RatingCeil); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := ratingFloor(searchBy.RatingFloor); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := pagesCeil(searchBy.PagesCeil); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := pagesFloor(searchBy.PagesFloor); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := ratingsCountCeil(searchBy.RatingsCountCeil); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := ratingsCountFloor(searchBy.RatingsCountFloor); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := reviewsCountCeil(searchBy.ReviewsCountCeil); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	if ok, q, f := reviewsCountFloor(searchBy.ReviewsCountFloor); ok {
		queryFields, fields = append(queryFields, q), append(fields, f)
	}

	return buildQuery(queryFields, searchIn), fields
}

// buildQuery builds a sql select query using given a list of search string.
func buildQuery(queryFields []string, searchIn *SearchIn) string {
	builder := new(strings.Builder)
	builder.WriteString(fmt.Sprintf("select * from %s", searchIn.BookTable))

	if len(queryFields) != 0 {
		builder.WriteString(" where ")
		for i, field := range queryFields {
			if i != 0 {
				builder.WriteString(" and ")
			}
			builder.WriteString(field)
		}
	}

	builder.WriteByte(';')

	return builder.String()
}

// titleHas checks if SearchIn's TitleHas field was specified, and if so
// returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func titleHas(titleHas string) (bool, string, string) {
	if titleHas != "" {
		return true, "title like ?", fmt.Sprintf("%%%s%%", titleHas)
	}

	return false, "", ""
}

// authors checks if SearchIn's Authors field was specified, and if so
// returns a search string to add to the sql query (prepared statement),
// and a slice of parameters needed for the prepared statement to execute.
func authors(authors []string) (bool, string, []interface{}) {
	if len(authors) != 0 {
		parameters := make([]interface{}, len(authors))
		builder := new(strings.Builder)
		builder.WriteByte('(')

		for i, author := range authors {
			parameters[i] = fmt.Sprintf("%%%s%%", author)
			if i != 0 {
				builder.WriteString(" or ")
			}
			builder.WriteString("authors like ?")
		}
		builder.WriteByte(')')

		return true, builder.String(), parameters
	}

	return false, "", nil
}

// languageCode checks if SearchIn's LanguageCode field was specified, and
// if so returns a search string to add to the sql query (prepared statement),
// and a slice of parameters needed for the prepared statement to execute.
func languageCode(codes []string) (bool, string, []interface{}) {
	if len(codes) != 0 {
		parameters := make([]interface{}, len(codes))
		builder := new(strings.Builder)
		builder.WriteByte('(')

		for i, code := range codes {
			parameters[i] = fmt.Sprintf("%%%s%%", code)

			if i != 0 {
				builder.WriteString(" or ")
			}
			builder.WriteString("languageCode like ?")
		}
		builder.WriteByte(')')

		return true, builder.String(), parameters
	}

	return false, "", nil
}

// isbn checks if SearchIn's ISBN field was specified, and if so returns
// a search string to add to the sql query (prepared statement), and a
// parameter needed for the prepared statement to execute.
func isbn(isbn string) (bool, string, string) {
	if isbn != "" {
		return true, "isbn = ?", isbn
	}

	return false, "", ""
}

// isbn13 checks if SearchIn's ISBN13 field was specified, and if so
// returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func isbn13(isbn13 string) (bool, string, string) {
	if isbn13 != "" {
		return true, "isbn13 = ?", isbn13
	}

	return false, "", ""
}

// ratingCeil checks if SearchIn's RatingCeil field was specified, and if
// so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func ratingCeil(ceil float32) (bool, string, float32) {
	if ceil >= 0 {
		return true, "averageRating <= ?", ceil
	}

	return false, "", 0
}

// ratingFloor checks if SearchIn's RatingFloor field was specified, and if
// so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func ratingFloor(floor float32) (bool, string, float32) {
	if floor >= 0 {
		return true, "averageRating > ?", floor
	}

	return false, "", 0
}

// pagesCeil checks if SearchIn's PagesCeil field was specified, and if
// so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func pagesCeil(ceil int) (bool, string, int) {
	if ceil >= 0 {
		return true, "pages <= ?", ceil
	}

	return false, "", 0
}

// pagesFloor checks if SearchIn's PagesFloor field was specified, and if
// so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func pagesFloor(floor int) (bool, string, int) {
	if floor >= 0 {
		return true, "pages > ?", floor
	}

	return false, "", 0
}

// ratingsCountCeil checks if SearchIn's RatingsCountCeil field was specified,
// and if so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func ratingsCountCeil(ceil int) (bool, string, int) {
	if ceil >= 0 {
		return true, "ratingsCount <= ?", ceil
	}

	return false, "", 0
}

// ratingsCountFloor checks if SearchIn's RatingsCountFloor field was specified,
// and if so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func ratingsCountFloor(floor int) (bool, string, int) {
	if floor >= 0 {
		return true, "ratingsCount > ?", floor
	}

	return false, "", 0
}

// reviewsCountCeil checks if SearchIn's ReviewsCountCeil field was specified,
// and if so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func reviewsCountCeil(ceil int) (bool, string, int) {
	if ceil >= 0 {
		return true, "reviewsCount <= ?", ceil
	}

	return false, "", 0
}

// reviewsCountFloor checks if SearchIn's ReviewsCountFloor field was specified,
// and if so returns a search string to add to the sql query (prepared statement),
// and a parameter needed for the prepared statement to execute.
func reviewsCountFloor(floor int) (bool, string, int) {
	if floor >= 0 {
		return true, "reviewsCount > ?", floor
	}

	return false, "", 0
}
