package main

import (
	"fmt"
	"os"

	"github.com/sudo-sturbia/bfr/internal/datastore"
	"github.com/sudo-sturbia/bfr/internal/server"
)

// Config holds all configuration needed to run the server.
type Config struct {
	Server    *server.Config    // Server's configuration options.
	Datastore *datastore.Config // Datastore's configuration options.
}

// New returns a new configuration object with initialized fields.
func New() *Config {
	return &Config{
		Server: &server.Config{
			Port: "6060",
		},

		Datastore: &datastore.Config{
			Driver:    "sqlite3",
			Dir:       fmt.Sprintf("%s/.config/bfr/", os.Getenv("HOME")),
			Datastore: "bfr.db",
			BookTable: "books",
		},
	}
}
