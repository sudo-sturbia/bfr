// Package main runs a server instance.
package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/internal/config"
	"github.com/sudo-sturbia/bfr/internal/datastore"
	"github.com/sudo-sturbia/bfr/internal/server"
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
		cfg.Server.Port = *port
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

	server := server.New(
		cfg.Server, &server.SearchIn{
			Datastore: datastore,
			BookTable: cfg.Datastore.BookTable,
		},
	)

	server.Run()
}

// usage prints a help message.
func usage() {
	fmt.Printf(
		"%s\n%s\n%s\n%s\n%s\n%s\n\n%s\n",
		"bfr is a web server that enables searching for books using a set of parameters.",
		"Usage:",
		"    bfr                 Run as a web server at :6060.",
		"    bfr -dataset path   Load a new csv dataset to use as a datastore, then runs the server.",
		"    bfr -port <number>  Use the specified port to run the server.",
		"    bfr -h              Print a help message.",
		"See github.com/sudo-sturbia/bfr.",
	)
}
