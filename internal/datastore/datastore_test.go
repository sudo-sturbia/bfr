package datastore

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// Test creation of a datastore.
func TestNew(t *testing.T) {
	config := &Config{
		Driver:    "sqlite3",
		Datastore: "testDatastore.db",
		BookTable: "books",
	}

	err := New("../../test-data/datastoreTest.csv", config, true)
	if err != nil {
		t.Errorf("Loading failed: %s.", err.Error())
		return
	}

	datastore, err := sql.Open(config.Driver, fmt.Sprintf("file:%s", config.Datastore))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer datastore.Close()
	defer os.Remove(config.Datastore)

	// Verify books.
	books := map[string]bool{
		"1,Harry Potter and the Half-Blood Prince (Harry Potter  #6),J.K. Rowling-Mary GrandPré,4.56,439785960,9780439785969,eng,652,1944099,26249":    true,
		"2,Harry Potter and the Order of the Phoenix (Harry Potter  #5),J.K. Rowling-Mary GrandPré,4.49,439358078,9780439358071,eng,870,1996446,27613": true,
		"3,Harry Potter and the Sorcerer's Stone (Harry Potter  #1),J.K. Rowling-Mary GrandPré,4.47,439554934,9780439554930,eng,320,5629932,70390":     true,
		"4,Harry Potter and the Chamber of Secrets (Harry Potter  #2),J.K. Rowling,4.41,439554896,9780439554893,eng,352,6267,272":                      true,
		"5,Harry Potter and the Prisoner of Azkaban (Harry Potter  #3),J.K. Rowling-Mary GrandPré,4.55,043965548X,9780439655484,eng,435,2149872,33964": true,
	}

	rows, err := datastore.Query("select * from books;")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer rows.Close()

	line := 0
	for rows.Next() {
		var (
			id            int
			title         string
			authors       string
			averageRating float32
			isbn          string
			isbn13        string
			languageCode  string
			pages         int
			ratingsCount  int
			reviewsCount  int
		)

		rows.Scan(
			&id,
			&title,
			&authors,
			&averageRating,
			&isbn,
			&isbn13,
			&languageCode,
			&pages,
			&ratingsCount,
			&reviewsCount,
		)

		row := fmt.Sprintf(
			"%d,%s,%s,%.2f,%s,%s,%s,%d,%d,%d",
			id,
			title,
			authors,
			averageRating,
			isbn,
			isbn13,
			languageCode,
			pages,
			ratingsCount,
			reviewsCount,
		)

		if !books[row] {
			t.Errorf("Incorrect row: %s.", row)
		}

		line++
	}

	if line != len(books) {
		t.Errorf("Expected %d rows, found %d.", len(books), line)
	}
}
