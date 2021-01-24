// package config handles server's configuration.
package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/internal/datastore"
)

// Config holds all configuration needed to run the server.
type Config struct {
	Host      string            // Host to run the server on.
	Port      string            // Port to run the server on.
	Datastore *datastore.Config // Datastore's configuration options.
}

// New returns a new Config object with initialized fields, and sets
// logger's options.
func New() *Config {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(os.Stderr)

	return &Config{
		Host: "",
		Port: "6060",
		Datastore: &datastore.Config{
			Driver:    "sqlite3",
			Dir:       fmt.Sprintf("%s/.config/bfr/", os.Getenv("HOME")),
			Datastore: "bfr.db",
			BookTable: "books",
		},
	}
}
