// Package main runs a frontend server instance.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sudo-sturbia/bfr/v2/internal/config"
	"github.com/sudo-sturbia/bfr/v2/internal/frontend"
)

var (
	api    = flag.String("api", "http://localhost:6060", "URL to use for API calls.")
	port   = flag.String("port", "5050", "Port number to run the server on.")
	static = flag.String("static", "content/static", "Path to static files.")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	cfg := config.NewOnPort(*port)
	server, err := frontend.New(
		&frontend.Config{
			Host:   cfg.Host,
			Port:   cfg.Port,
			Static: *static,
		},
		*api,
	)
	if err != nil {
		log.Fatalf("failed to create server: %s", err.Error())
	}

	server.Run()
}

// usage prints a help message.
func usage() {
	fmt.Println(
		"A frontend web server that utilizes bfr API to search for and find books.\n",
		"Usage:\n",
		"    go run ./cmd/frontend                 Run a frontend server at localhost:5050.\n",
		"    go run ./cmd/frontend -api <url>      Use given URL for API calls, default is localhost:6060.\n",
		"    go run ./cmd/frontend -port <number>  Run server at specified port.\n",
		"    go run ./cmd/frontend -static <path>  Use given path for static files.\n",
		"    go run ./cmd/frontend -h              Print this help message.\n",
		"\n",
		"See github.com/sudo-sturbia/bfr.",
	)
}
