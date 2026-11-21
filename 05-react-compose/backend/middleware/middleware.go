// Package middleware provides HTTP middlewares: RequestID, Logger and Recover.
package middleware

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"react-compose/logger"
)

// requestIDHeader is the conventional header used to propagate the request
// identifier between the client, the proxies (nginx) and the backend.
const requestIDHeader = "X-Request-ID"

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

// RequestID gives every request an identifier: it reuses a valid incoming
// X-Request-ID (set by a proxy or the client) or generates a UUIDv4. The ID
// travels back in the response header and, more importantly, a request-scoped
// logger carrying it is stored in the context: every log line emitted from
// From(ctx) for this request is automatically tagged with request_id.
// It must be the outermost middleware so everything below gets the logger.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := newRequestID(r)

		l := slog.Default().With("request_id", id)
		w.Header().Set(requestIDHeader, id)

		next.ServeHTTP(w, r.WithContext(logger.Into(r.Context(), l)))
	})
}

// newRequestID reuses a valid incoming header value or generates a fresh UUIDv4.
func newRequestID(r *http.Request) string {
	if id := r.Header.Get(requestIDHeader); validRequestID(id) {
		return id
	}
	return newUUIDv4()
}

// validRequestID keeps only short printable values: the header is
// user-controlled, so it must be bounded before being echoed and logged.
func validRequestID(id string) bool {
	if id == "" || len(id) > 64 {
		return false
	}
	for _, c := range id {
		if c <= 0x20 || c >= 0x7f {
			return false
		}
	}
	return true
}

// newUUIDv4 returns a random RFC 4122 version-4 UUID string.
func newUUIDv4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])      // never fails on the supported platforms
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10xx
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// Logger emits one structured record per request with method, path, status,
// duration and remote address, tagged with the request_id set by RequestID.
// Failed statuses stand out: 5xx logs at ERROR, 4xx at WARN, the rest at INFO.
// It must sit outside Recover so it also records the 500 produced by it.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		level := slog.LevelInfo
		switch {
		case rw.status >= 500:
			level = slog.LevelError
		case rw.status >= 400:
			level = slog.LevelWarn
		}

		logger.From(r.Context()).LogAttrs(r.Context(), level, "request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.status),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}

// Recover turns a panic into a 500 instead of killing the server.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.From(r.Context()).Error("panic recovered",
					slog.Any("panic", err),
					slog.String("stack", string(debug.Stack())),
				)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
