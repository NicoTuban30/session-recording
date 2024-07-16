package main

import (
	"log"
	"net/http"

	"cassette/config"
	"cassette/pkg/server"
)

func main() {
	cfg, err := config.FromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Create the server with the config, repository, and storage
	s := server.New(cfg, cfg.Repository, cfg.Storage)

	// Start the HTTP server
	if err := http.ListenAndServe(":3000", s); err != nil {
		log.Fatal(err)
	}
}
