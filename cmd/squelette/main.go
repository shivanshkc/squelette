package main

import (
	"github.com/shivanshkc/squelette/internal/http"
	"github.com/shivanshkc/squelette/pkg/config"
	"github.com/shivanshkc/squelette/pkg/logger"
)

func main() {
	// Initialize basic dependencies.
	cfg := config.Load()
	log := logger.New(cfg)

	// Initialize the HTTP server.
	server := http.Server{
		Config:     cfg,
		Logger:     log,
		Middleware: &http.Middleware{Logger: log},
	}

	// This internally calls ListenAndServe.
	// This is a blocking call and will panic if the server is unable to start.
	server.Start()
}
