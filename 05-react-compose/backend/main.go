package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"react-compose/controllers"
	"react-compose/db"
	"react-compose/logger"
	"react-compose/middleware"
)

func main() {
	// Logger: structured logging with log/slog, configured by LOG_FORMAT and
	// LOG_LEVEL. Must run before anything else logs.
	l := logger.Setup()

	// Database: connect and migrate before serving
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.Connect(ctx); err != nil {
		l.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	if err := db.Migrate(ctx); err != nil {
		l.Error("migration failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()
	l.Info("database connected and up to date")

	// Routing
	mux := http.NewServeMux()

	// Handlers
	mux.HandleFunc("GET /health", controllers.Health)

	// API: products CRUD (JSON)
	mux.HandleFunc("GET /api/products", controllers.ListProducts)
	mux.HandleFunc("POST /api/products", controllers.CreateProduct)
	mux.HandleFunc("GET /api/products/{id}", controllers.GetProduct)
	mux.HandleFunc("PATCH /api/products/{id}", controllers.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", controllers.DeleteProduct)

	// Middlewares: RequestID outermost (every log line of a request carries
	// its request_id), then Logger so it records the 500 produced by Recover.
	handler := middleware.RequestID(middleware.Logger(middleware.Recover(mux)))

	// Server configuration
	puerto := ":8080"
	srv := &http.Server{
		Addr:              puerto,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	l.Info("server listening", slog.String("addr", "http://localhost"+srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		l.Error("critical error starting the server", slog.Any("error", err))
		os.Exit(1)
	}
}
