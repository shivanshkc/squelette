package main

import (
	"os"

	"github.com/shivanshkc/squelette/internal/http"
	"github.com/shivanshkc/squelette/pkg/config"
	"github.com/shivanshkc/squelette/pkg/logger"
)

func main() {
	// Initialize basic dependencies.
	conf := config.Load()
	logger.Init(os.Stdout, conf.Logger.Level, conf.Logger.Pretty)

	// Initialize the HTTP server.
	server := &http.Server{
		Config:     conf,
		Middleware: http.Middleware{},
	}

	// This internally calls ListenAndServe.
	// This is a blocking call and will panic if the server is unable to start.
	server.Start()
}
