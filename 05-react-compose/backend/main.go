package main

import (
	"context"
	"log"
	"net/http"
	"react-compose/controllers"
	"react-compose/db"
	"react-compose/middleware"
	"time"
)

func main() {
	// Database: connect and migrate before serving
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.Connect(ctx); err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	defer db.Close()
	log.Printf("✅ Database connected and up to date")

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

	// Middlewares: Recover inside so Logger records the 500
	handler := middleware.Logger(middleware.Recover(mux))

	// Server configuration
	puerto := ":8080"
	srv := &http.Server{
		Addr:              puerto,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("🚀 Server listening at http://localhost%s\n", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("❌ Critical error starting the server: %v", err)
	}
}
