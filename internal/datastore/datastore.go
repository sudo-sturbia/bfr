// Package datastore contains configuration options, and functions to
// create a books' database before running the server.
package datastore

// Config holds datastore's configuration options.
type Config struct {
	Driver string // DBMS's driver.

	Dir       string // Directory containing the datastore.
	Datastore string // Datastore's name.
	BookTable string // Table containing books.

	logTo string // Name of the file to write logs to.
}
