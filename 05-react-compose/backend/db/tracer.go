// queryTracer logs every SQL statement at DEBUG level.
//
// pgx calls the tracer around each statement executed through the pool, and
// passes the caller's context through: that is the request context, so the
// SQL lines are tagged with the same request_id as the rest of the request.
// The tracer is attached unconditionally and relies on the slog level to
// decide whether anything is actually printed (LOG_LEVEL=debug to enable).
package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"react-compose/logger"
)

// queryTracer implements pgx.QueryTracer.
type queryTracer struct{}

// queryTrace carries what Start collected so End can log it: pgx only
// passes SQL/Args on Start and the result on End.
type queryTrace struct {
	start time.Time
	sql   string
	args  []any
}

// traceKey stores the queryTrace value in the tracer context.
type traceKey struct{}

// TraceQueryStart runs before the statement is sent to PostgreSQL.
func (queryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, traceKey{}, queryTrace{
		start: time.Now(),
		sql:   data.SQL,
		args:  data.Args,
	})
}

// TraceQueryEnd runs once the statement finishes (success or error).
func (queryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	trace, _ := ctx.Value(traceKey{}).(queryTrace)

	attrs := []slog.Attr{
		slog.String("sql", trace.sql),
		slog.Any("args", trace.args),
		slog.Duration("duration", time.Since(trace.start)),
	}
	if data.Err != nil {
		// Logged at DEBUG on purpose: the handler already reports failures
		// at ERROR with the operation name; this line adds the SQL context.
		attrs = append(attrs, slog.Any("error", data.Err))
	} else {
		attrs = append(attrs, slog.Int64("rows", data.CommandTag.RowsAffected()))
	}

	logger.From(ctx).LogAttrs(ctx, slog.LevelDebug, "query", attrs...)
}
