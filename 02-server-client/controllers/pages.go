// Package controllers renders the HTML pages of the site
package controllers

import (
	"html/template"
	"log"
	"net/http"

	"server-client/models"
)

var templates = template.Must(template.ParseFiles("views/index.html"))

// pageData is what the template receives on each request
type pageData struct {
	Products []models.Product
	Error    string
	Success  string
}

// Home renders the main page with the product list.
func Home(w http.ResponseWriter, r *http.Request) {
	// Data
	query := r.URL.Query()
	data := pageData{
		Products: store.List(),
		Error:    query.Get("error"),
		Success:  query.Get("ok"),
	}

	// Headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Status Code
	w.WriteHeader(http.StatusOK)

	// Render the view
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("❌ Error rendering index.html: %v", err)
	}
}
