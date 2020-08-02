// Package datastore contains configuration options, and functions to
// create a books' database before running the server.
package datastore

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// Number of dataset's columns.
const Columns = 10

// Config holds datastore's configuration options.
type Config struct {
	Driver string // DBMS's driver.

	Dir       string // Directory containing the datastore.
	Datastore string // Datastore's name.
	BookTable string // Table containing books.
}

// Open opens a connection to a database specified by given configuration.
func Open(config *Config) (*sql.DB, error) {
	return sql.Open(config.Driver, fmt.Sprintf("file:%s%s", config.Dir, config.Datastore))
}

// New creates a new datastore to be used by the server. The datastore is created
// at the path specified by config (Dir + Name). overwriteIfExists specifies what
// to do if a datastore with the same path exists. The datastore is created using
// dataset at the specified path. Dataset should be a csv file with the following
// columns (id, title, authors, averageRating, isbn, isbn13, languageCode, pages,
// ratingsCount, textReviewsCount).
// See https://www.kaggle.com/jealousleopard/goodreadsbooks
func New(datasetPath string, config *Config, overwriteIfExists bool) error {
	dataset, err := os.Open(datasetPath)
	if err != nil {
		return err
	}
	defer dataset.Close()

	if _, err := os.Stat(config.Dir + config.Datastore); !os.IsNotExist(err) && !overwriteIfExists {
		return err
	}

	os.MkdirAll(config.Dir, 0777)
	datastore, err := sql.Open(config.Driver, fmt.Sprintf("file:%s%s", config.Dir, config.Datastore))
	if err != nil {
		return err
	}
	defer datastore.Close()

	create := fmt.Sprintf(
		"create table %s ("+
			"id integer not null primary key, "+
			"title text, "+
			"authors text, "+
			"averageRating float, "+
			"isbn string, "+
			"isbn13 string, "+
			"languageCode text, "+
			"pages integer, "+
			"ratingsCount integer, "+
			"reviewsCount integer);", config.BookTable)

	_, err = datastore.Exec(create)
	if err != nil {
		return err
	}

	tx, err := datastore.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = insertBooks(dataset, tx, config)
	if err != nil {
		return err
	}

	return nil
}

// insertBooks inserts books from a dataset (csv file) into a table using
// given transaction. Corrupt lines are logged and skipped.
func insertBooks(dataset *os.File, tx *sql.Tx, config *Config) error {
	// Use prepared statements to populate books' table in bfr's database.
	insert := fmt.Sprintf("insert into %s values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", config.BookTable)
	stmt, err := tx.Prepare(insert)
	if err != nil {
		return err
	}

	line := 0
	scanner := bufio.NewScanner(dataset)
	for scanner.Scan() {
		fields := strings.FieldsFunc(scanner.Text(), func(c rune) bool {
			return c == ','
		})

		if len(fields) != Columns {
			log.WithFields(
				log.Fields{
					"line":   line,
					"fields": len(fields),
				},
			).Error("Corrupted line.")
			continue
		}

		_, err := stmt.Exec(
			fields[0],
			fields[1],
			fields[2],
			fields[3],
			fields[4],
			fields[5],
			fields[6],
			fields[7],
			fields[8],
			fields[9],
		)
		if err != nil {
			log.WithFields(
				log.Fields{
					"line":   line,
					"fields": len(fields),
				},
			).Error(err)
		}

		line++
	}

	return nil
}
