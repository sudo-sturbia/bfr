// Package datastore contains configuration options, and functions to
// create a books' database before running the server.
package datastore

// Config holds datastore's configuration options.
type Config struct {
	Driver string // Database's driver.
	Path   string // Datastore to use.
}
