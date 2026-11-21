// Package middleware provides HTTP middlewares: Logger and Recover
package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps the original http.ResponseWriter to intercept
// and record the HTTP status code, since the standard interface does not allow
// it to be read once it has been written.
// Flow: The wrapper travels down toward the handlers; the status travels back up.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader intercepts the assignment of the status code to store it locally
// within the structure before delegating the call to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Middlewares

// Logger prints method, path, status and duration of every request.
// It must be the outer middleware so it also logs the 500 from Recover.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf("%s %s → %d (%s)", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

// Recover turns a panic into a 500 instead of killing the server.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ panic: %v", err)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
