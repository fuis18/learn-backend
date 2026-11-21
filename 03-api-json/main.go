package main

import (
	"log"
	"net/http"
	"time"

	"api-json/controllers"
	"api-json/middleware"
)

func main() {
	// Routing
	mux := http.NewServeMux()

	// Static files (css, js)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handlers
	mux.HandleFunc("GET /", controllers.Home)
	mux.HandleFunc("GET /health", controllers.Health)

	// API: products CRUD (JSON)
	mux.HandleFunc("GET /api/products", controllers.ListProducts)
	mux.HandleFunc("POST /api/products", controllers.CreateProduct)
	mux.HandleFunc("GET /api/products/{id}", controllers.GetProduct)
	mux.HandleFunc("PUT /api/products/{id}", controllers.ReplaceProduct)
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
