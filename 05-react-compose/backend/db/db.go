// Package db manages the PostgreSQL connection using pgx (connection pool).
package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Pool is the global connection pool shared across the application.
// It must be initialized by calling Connect before executing queries.
var Pool *pgxpool.Pool

// Connect initializes the global database connection pool using DB_* environment variables.
func Connect(ctx context.Context) error {
	// Connect is idempotent: subsequent calls are no-ops if the pool is already established.
	if Pool != nil {
		return nil
	}

	loadDevEnv()

	target, dsn, err := buildDSN()
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("db: invalid config for %s: %w", target, err)
	}

	// It verifies connectivity with a 5-second timeout ping.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return fmt.Errorf("db: no answer from %s: %w", target, err)
	}

	Pool = pool
	return nil
}

// Close gracefully closes the active connection pool and resets Pool to nil.
func Close() {
	if Pool != nil {
		Pool.Close()
		Pool = nil
	}
}

// buildDSN constructs a secure PostgreSQL connection string (DSN) using net/url
// to properly escape special characters in credentials from the DB_* environment variables.
func buildDSN() (target string, dsn string, err error) {
	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	user := env("DB_USER", "app")
	pass := env("DB_PASSWORD", "app123")
	name := env("DB_NAME", "app")
	sslmode := env("DB_SSLMODE", "disable")

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + name,
	}
	q := url.Values{}
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	// It returns a sanitized, human-readable target string (omitting the password) for safe error logging.
	target = fmt.Sprintf("%s/%s (sslmode=%s)", u.Host, name, sslmode)
	return target, u.String(), nil
}

// loadDevEnv loads the .env file at the repo root so `go run .` works on the
// host without exporting variables by hand. godotenv never overrides variables
// already present, so in containers the compose-provided environment wins.
// A missing file is expected (FROM scratch images) and is not an error.
func loadDevEnv() {
	for _, file := range []string{"../.env", ".env"} {
		if err := godotenv.Load(file); err == nil {
			return
		}
	}
}

// env returns the value of the environment variable for the given key,
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
