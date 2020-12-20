// package config handles server's configuration.
package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/internal/datastore"
	"github.com/sudo-sturbia/bfr/internal/server"
)

// Config holds all configuration needed to run the server.
type Config struct {
	Server    *server.Config    // Server's configuration options.
	Datastore *datastore.Config // Datastore's configuration options.
}

// New returns a new Config object with initialized fields, and sets
// logger's options.
func New() *Config {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(os.Stderr)

	return &Config{
		Server: &server.Config{
			Host: "",
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
