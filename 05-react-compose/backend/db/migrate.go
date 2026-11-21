// Package db provides a schema migration runner using embedded SQL files.
package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// migrationFS embeds the migration SQL files directly into the binary,
//
//go:embed migrations/*.sql
var migrationFS embed.FS

// Migrate applies pending *.sql migration files in alphabetical order and records
// them in the schema_migrations table.
func Migrate(ctx context.Context) error {
	// Each migration runs in its own transaction along with its version insert.
	if Pool == nil {
		return fmt.Errorf("db: Migrate called before Connect")
	}

	_, err := Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	if err != nil {
		return fmt.Errorf("db: creating schema_migrations: %w", err)
	}

	applied, err := appliedVersions(ctx)
	if err != nil {
		return err
	}

	// sorted alphabetically by name
	files, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("db: reading embedded migrations: %w", err)
	}

	// Calling Migrate multiple times is safe and acts as a no-op if everything is up to date.
	pending := 0
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || applied[name] {
			continue
		}
		if err := applyOne(ctx, name); err != nil {
			return err
		}
		slog.Info("migration applied", slog.String("version", name))
		pending++
	}

	if pending == 0 {
		slog.Info("database up to date", slog.Int("applied", len(applied)))
	}
	return nil
}

// appliedVersions queries the database and returns a set of already recorded migration versions.
func appliedVersions(ctx context.Context) (map[string]bool, error) {
	rows, err := Pool.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("db: reading schema_migrations: %w", err)
	}
	defer rows.Close()

	versions := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("db: scanning version: %w", err)
		}
		versions[version] = true
	}
	return versions, rows.Err()
}

// applyOne reads a single migration file and executes its content alongside
// the schema_migrations record insertion within an atomic transaction.
func applyOne(ctx context.Context, name string) error {
	content, err := migrationFS.ReadFile("migrations/" + name)
	if err != nil {
		return fmt.Errorf("db: reading %s: %w", name, err)
	}

	tx, err := Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db: beginning transaction for %s: %w", name, err)
	}

	// Safely handle rollback, ignoring ErrTxClosed if the transaction was already committed.
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Warn("migration rollback failed", slog.String("version", name), slog.Any("error", rbErr))
		}
	}()

	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("db: applying SQL for %s: %w", name, err)
	}

	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", name); err != nil {
		return fmt.Errorf("db: recording version %s: %w", name, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("db: committing transaction for %s: %w", name, err)
	}

	return nil
}
