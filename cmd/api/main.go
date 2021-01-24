// Package main runs a backend server instance.
package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/internal/api"
	"github.com/sudo-sturbia/bfr/internal/config"
	"github.com/sudo-sturbia/bfr/internal/datastore"
)

var (
	port    = flag.String("port", "", "Specify a port to run the server on.")
	dataset = flag.String("dataset", "", "Load a new csv dataset from specified path.")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	cfg := config.New()

	if *port != "" {
		cfg.Port = *port
	}
	if *dataset != "" {
		err := datastore.New(*dataset, cfg.Datastore, true)
		if err != nil {
			log.Fatalf("Failed to create a datastore: %s.", err.Error())
		}
	}

	datastore, err := datastore.Open(cfg.Datastore)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := api.New(
		&api.Config{
			Host: cfg.Host,
			Port: cfg.Port,
		},
		&api.SearchIn{
			Datastore: datastore,
			BookTable: cfg.Datastore.BookTable,
		},
	)

	server.Run()
}

// usage prints a help message.
func usage() {
	fmt.Println(
		"A REST API that enables searching for books using a set of parameters.\n",
		"Usage:\n",
		"    go run ./cmd/api                 Run a backend server at localhost:6060.\n",
		"    go run ./cmd/api -dataset path   Load a new csv dataset to use as a datastore, then run the server.\n",
		"    go run ./cmd/api -port <number>  Use the specified port to run the server.\n",
		"    go run ./cmd/api -h              Print a help message.\n",
		"\n",
		"See github.com/sudo-sturbia/bfr.",
	)
}
