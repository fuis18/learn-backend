// Package logger configures the structured logger used by the whole backend.
//
// Two environment variables control it:
//
//	LOG_FORMAT  text (default) | json
//	LOG_LEVEL   debug | info (default) | warn | error
//
// Setup must be called at the very beginning of main so that every later
// log call (middlewares, controllers, db) goes through the same handler.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// loggerKey is an unexported context key: no other package can collide with it.
type loggerKey struct{}

// Setup builds the application logger from LOG_FORMAT and LOG_LEVEL,
// installs it as the slog default and returns it.
func Setup() *slog.Logger {
	// LOG_LEVEL: debug | info | warn | error. Empty means info.
	rawLevel := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if rawLevel == "" {
		rawLevel = "info"
	}
	level, levelErr := parseLevel(rawLevel)

	// LOG_FORMAT: text | json. Empty (or unknown) falls back to text.
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))

	opts := &slog.HandlerOptions{
		Level: level,
		// AddSource annotates every record with the file:line that emitted
		// it. Only enabled at debug, where it is worth the extra work.
		AddSource: level <= slog.LevelDebug,
	}

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l := slog.New(handler)
	slog.SetDefault(l)

	// Configuration feedback: warnings for values we had to ignore, and a
	// debug line that doubles as proof that LOG_LEVEL=debug works.
	if levelErr != nil {
		l.Warn("invalid LOG_LEVEL, using info", slog.String("value", rawLevel), slog.Any("error", levelErr))
	}
	if format != "text" && format != "json" {
		l.Warn("invalid LOG_FORMAT, using text", slog.String("value", format))
	}
	l.Debug("logger configured",
		slog.String("format", format),
		slog.String("level", level.String()),
	)

	return l
}

// parseLevel converts a LOG_LEVEL value into a slog.Level.
func parseLevel(v string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(strings.ToLower(v)))
	return level, err
}

// Into stores l in ctx so the rest of the request can recover it with From.
// The RequestID middleware uses it to attach the request-scoped logger.
func Into(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// From returns the logger stored in ctx by Into, or the default logger
// when there is none (startup code, tests, misrouted calls).
func From(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
