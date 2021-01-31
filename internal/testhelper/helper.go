package testhelper

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sudo-sturbia/bfr/v2/internal/datastore"
)

// SearchIn returns a datastore, and a table name to use for constructing a SearchIn to
// be used in testing. It also returns a function that frees resources, and should be
// defered immediately after return.
func SearchIn(t *testing.T) (*sql.DB, string, func()) {
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

	return datastore, config.BookTable, deferFn
}
