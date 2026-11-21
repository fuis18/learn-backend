package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Routing
	mux := http.NewServeMux()

	// Handlers
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /health", handleHealth)

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

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Headers
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Status Code
	w.WriteHeader(http.StatusOK)

	_, _ = fmt.Fprintf(w, "Welcome to Home\n")
	_, _ = fmt.Fprintf(w, "Method: %s\n", r.Method)
	_, _ = fmt.Fprintf(w, "Query: %v\n", r.URL.Query())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Headers
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Status Code
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK\n"))
}
