// Package main runs a new server instance.
package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sudo-sturbia/bfr/internal/config"
	"github.com/sudo-sturbia/bfr/internal/datastore"
	"github.com/sudo-sturbia/bfr/internal/server"
)

var (
	cfg = config.New()
)

func main() {
	parseFlags()

	datastore, err := datastore.Open(cfg.Datastore)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := server.New(cfg.Server, &server.SearchIn{
		Datastore: datastore,
		BookTable: cfg.Datastore.BookTable,
	})

	server.Run()
}

// parseFlags parses, and handles command line flags.
func parseFlags() {
	help := flag.Bool("h", false, "Print a help message.")
	port := flag.String("port", "", "Specify a port to run the server on.")
	dataset := flag.String("dataset", "", "Load a new csv dataset from specified path.")
	flag.Parse()

	if *help {
		description()
	}
	if *port != "" {
		specifyPort(*port)
	}
	if *dataset != "" {
		createDatastore(*dataset)
	}
}

// description prints a help message.
func description() {
	defer os.Exit(0)
	fmt.Printf(
		"%s\n%s\n%s\n%s\n%s\n%s\n\n%s\n",
		"bfr is a REST API to search for books using a set of parameters.",
		"Usage:",
		"    bfr                 Runs as a web server at :6060.",
		"    bfr -dataset path   Loads a new csv dataset to use as a datastore, then runs the server.",
		"    bfr -port <number>  Specifies a different port to run the server on.",
		"    bfr -h              Prints a help message.",
		"See github.com/sudo-sturbia/bfr.",
	)
}

// specifyPort updates the port that the server runs on.
func specifyPort(port string) {
	cfg.Server.Port = port
}

// createDatastore creates a new datastore using dataset at specified
// path. Datastore is created at position specified by Config.
func createDatastore(dataset string) {
	err := datastore.New(dataset, cfg.Datastore, true)
	if err != nil {
		log.Fatalf("Failed to create a datastore: %s.", err.Error())
	}
}
