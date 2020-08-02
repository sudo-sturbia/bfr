package main

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

// newConfig returns a new configuration object with initialized fields,
// and sets global logger's configuration options.
func newConfig() *Config {
	configLogger()
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

// configLogger sets global logger's configuration options.
func configLogger() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(os.Stderr)
}
