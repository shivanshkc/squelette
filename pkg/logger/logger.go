package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/shivanshkc/squelette/pkg/config"
	"github.com/shivanshkc/squelette/pkg/utils/ctxutils"
)

// Logger is a wrapper around zerolog.Logger to provide custom methods on it.
type Logger struct {
	*zerolog.Logger

	Config *config.Config
}

// WithContext creates a new logger with this logger as the base.
// The new logger by default logs the request metadata present in the given context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Get the loggable data out of the context.
	ctxInfo := ctxutils.GetRequestCtxInfo(ctx)
	if ctxInfo == nil {
		return l
	}

	// Add the required fields to the subLogger.
	subLogger := l.With().
		Str("trace_id", ctxInfo.TraceID).
		Str("request_id", ctxInfo.RequestID).
		// More fields can be added here.
		Logger()

	return &Logger{Logger: &subLogger}
}

// New creates a new Logger instance.
func New(conf *config.Config) *Logger {
	// Decide the formatting based on the config.
	var logOutput io.Writer
	if conf.Logger.Pretty {
		logOutput = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		logOutput = os.Stdout
	}

	// Determine the log level.
	level, err := zerolog.ParseLevel(conf.Logger.Level)
	if err != nil {
		panic("unknown log level provided: " + conf.Logger.Level)
	}

	// Instantiate the logger.
	zLogger := zerolog.New(logOutput).
		Level(level).With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{Logger: &zLogger, Config: conf}
}
