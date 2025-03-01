package main

import (
	"os"

	"github.com/shivanshkc/squelette/internal/config"
	"github.com/shivanshkc/squelette/internal/handlers"
	"github.com/shivanshkc/squelette/internal/http"
	"github.com/shivanshkc/squelette/internal/logger"
	"github.com/shivanshkc/squelette/internal/middleware"
	"github.com/shivanshkc/squelette/pkg/signals"
)

func main() {
	// All signals.Xxx calls are for interruption (SIGINT, SIGTERM) handling.
	// Wait blocks until all actions (registered by signals.OnSignal) have executed.
	defer signals.Wait()
	// Manually trigger cleanup whenever main exits.
	// This MUST run before signals.Wait and so it is deferred after it.
	defer signals.Manual()

	// Initialize basic dependencies.
	conf := config.Load()
	logger.Init(os.Stdout, conf.Logger.Level, conf.Logger.Pretty)

	// Initialize the HTTP server.
	server := &http.Server{Config: conf, Middleware: middleware.Middleware{}, Handler: &handlers.Handler{}}
	// Shutdown server upon interruption or exit.
	signals.OnSignal(func(_ os.Signal) { server.Shutdown() })

	// This internally calls ListenAndServe.
	// This is a blocking call and will panic if the server is unable to start.
	if err := server.Start(); err != nil {
		panic("error in server.Start call: " + err.Error())
	}
}
