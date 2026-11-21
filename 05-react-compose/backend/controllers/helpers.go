package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"

	"react-compose/logger"
)

// writeJSON writes a Go value as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// readJSON decodes the request body into dst.
// It rejects unknown fields and extra data after the object.
func readJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	if dec.More() {
		return errors.New("body must contain a single JSON object")
	}
	return nil
}

// parseID extracts the {id} path variable and converts it to an int.
func parseID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

// internalError logs the real cause with the request-scoped logger (so the
// line carries the request_id) and answers with a generic 500,
// so no database details leak to the client.
func internalError(w http.ResponseWriter, r *http.Request, op string, err error) {
	logger.From(r.Context()).Error(op+" failed", slog.Any("error", err))
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

// isUniqueViolation reports whether err is a PostgreSQL unique_violation (23505),
// used to answer 409 when a product code is already taken.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
