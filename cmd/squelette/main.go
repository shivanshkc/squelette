package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/shivanshkc/squelette/internal/config"
	"github.com/shivanshkc/squelette/internal/http"
	"github.com/shivanshkc/squelette/internal/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize basic dependencies.
	conf := config.Load()
	logger.Init(os.Stdout, conf.Logger.Level, conf.Logger.Pretty)

	// Initialize the HTTP server.
	server := &http.Server{}

	// Start the http server. The server will shut down when the context expires.
	if err := server.Start(ctx, conf.HTTPServer.Addr); err != nil {
		panic("error in server.Start call: " + err.Error())
	}
}
