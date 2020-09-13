package books

import (
	"fmt"
	"strings"
)

// query generates a SQL select query based on fields specified in SearchBy.
// Returns a prepared statement, and a list of parameters to use with the prepared
// statement when executing. If titles is true, then select statement only selects
// books' titles.
func query(searchIn *SearchIn, searchBy *SearchBy, titles bool) (string, []interface{}) {
	queryParts := make([]string, 0)
	fields := make([]interface{}, 0)

	if ok, q, f := titleHas(searchBy.TitleHas); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := authors(searchBy.Authors); ok {
		queryParts, fields = append(queryParts, q), append(fields, f...)
	}

	if ok, q, f := languageCode(searchBy.LanguageCode); ok {
		queryParts, fields = append(queryParts, q), append(fields, f...)
	}

	if ok, q, f := isbn(searchBy.ISBN); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := isbn13(searchBy.ISBN13); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := ratingCeil(searchBy.RatingCeil); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := ratingFloor(searchBy.RatingFloor); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := pagesCeil(searchBy.PagesCeil); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := pagesFloor(searchBy.PagesFloor); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := ratingsCountCeil(searchBy.RatingsCountCeil); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := ratingsCountFloor(searchBy.RatingsCountFloor); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := reviewsCountCeil(searchBy.ReviewsCountCeil); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	if ok, q, f := reviewsCountFloor(searchBy.ReviewsCountFloor); ok {
		queryParts, fields = append(queryParts, q), append(fields, f)
	}

	return buildQuery(queryParts, searchIn, titles), fields
}

// buildQuery builds a sql select query using given a list of search string.
// If titles is true, then buildQuery returns a sql statement that selects
// titles only, otherwise a statement that selects all columns is returned.
func buildQuery(queryParts []string, searchIn *SearchIn, titles bool) string {
	builder := new(strings.Builder)

	if titles {
		builder.WriteString(fmt.Sprintf("select title from %s", searchIn.BookTable))
	} else {
		builder.WriteString(fmt.Sprintf("select * from %s", searchIn.BookTable))
	}

	if len(queryParts) != 0 {
		builder.WriteString(" where ")
		for i, field := range queryParts {
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
