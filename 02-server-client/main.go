package main

import (
	"log"
	"net/http"
	"time"

	"server-client/controllers"
)

func main() {
	// Routing
	mux := http.NewServeMux()

	// Static files (css)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handlers
	mux.HandleFunc("GET /", controllers.Home)
	mux.HandleFunc("POST /products", controllers.CreateProduct)
	mux.HandleFunc("GET /health", Health)

	// Server configuration
	puerto := ":8080"
	srv := &http.Server{
		Addr:              puerto,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("🚀 Server listening at http://localhost%s\n", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("❌ Critical error starting the server: %v", err)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK\n"))
}
