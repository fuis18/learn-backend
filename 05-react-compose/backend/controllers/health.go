package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"react-compose/db"
	"react-compose/logger"
)

type checkResult struct {
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
}

type healthResponse struct {
	Status string                 `json:"status"`
	Checks map[string]checkResult `json:"checks,omitempty"`
}

// Health verifies that the server and PostgreSQL are reachable.
func Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := db.Pool.Ping(ctx)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		logger.From(r.Context()).Warn("health check failed",
			slog.String("check", "postgres"),
			slog.Any("error", err),
			slog.Int64("latency_ms", latency),
		)
		writeJSON(w, http.StatusServiceUnavailable, healthResponse{
			Status: "down",
			Checks: map[string]checkResult{
				"postgres": {Status: "down", LatencyMs: latency},
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
		Checks: map[string]checkResult{
			"postgres": {Status: "up", LatencyMs: latency},
		},
	})
}
