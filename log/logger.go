package log

import (
	"context"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/exp/slog"

	"github.com/belak/toolkit/internal"
)

const loggerContextKey internal.ContextKey = "Logger"

type loggerConfig struct {
	Level slog.Level `envconfig:"level"`
}

// NewLogger returns an slog.Logger configured for the current environment. The
// log level can be overridden with the LOG_LEVEL environment variable.
func NewLogger() (*slog.Logger, error) {
	var config loggerConfig

	err := envconfig.Process("LOG", &config)

	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(config.Level),
	}

	// Yes, there are instances where you could have a terminal and be in prod
	// (or not have a terminal and be in dev mode), but because of how this is
	// set up, that shouldn't be an issue.
	if internal.IsATTY() {
		return slog.New(opts.NewTextHandler(os.Stdout)), err
	} else {
		return slog.New(opts.NewJSONHandler(os.Stdout)), err
	}
}

// LoggerMiddleware returns a middleware which inserts the *slog.Logger into the
// context.
func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return internal.ContextValueMiddleware(loggerContextKey, logger)
}

// ExtractLogger extracts the logger from the given context.Context or panics if
// none exists.
func ExtractLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*slog.Logger); ok {
		return logger
	}

	panic("no logger in context")
}
