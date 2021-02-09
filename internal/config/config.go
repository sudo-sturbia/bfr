// Package config handles server's configuration.
package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/v2/internal/datastore"
)

// Config holds all configuration needed to run the server.
type Config struct {
	Host      string            // Host to run the server on.
	Port      string            // Port to run the server on.
	Datastore *datastore.Config // Datastore's configuration options.
}

// New returns a new Config object with Port as 6060.
func New() *Config {
	return NewOnPort("6060")
}

// NewOnPort returns a new Config object with Port set to given.
func NewOnPort(port string) *Config {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(os.Stderr)

	return &Config{
		Host: "",
		Port: port,
		Datastore: &datastore.Config{
			Driver:    "sqlite3",
			Dir:       fmt.Sprintf("%s/.config/bfr/", os.Getenv("HOME")),
			Datastore: "bfr.db",
			BookTable: "books",
		},
	}
}
