// Package controllers renders the HTML pages of the site
package controllers

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("views/index.html"))

// Home renders the main page.
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("❌ Error rendering index.html: %v", err)
	}
}
