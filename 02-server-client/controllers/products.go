package controllers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"server-client/models"
)

// store holds the products in memory
var store = models.NewStore()

func redirectToHome(w http.ResponseWriter, r *http.Request, key, message string) {
	query := url.Values{}
	query.Set(key, message)
	http.Redirect(w, r, "/?"+query.Encode(), http.StatusSeeOther)
}

// CreateProduct reads the form and stores a new product.
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Form fields validation
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		redirectToHome(w, r, "error", "The name is required")
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil || price < 0 {
		redirectToHome(w, r, "error", "Invalid price")
		return
	}

	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil || stock < 0 {
		redirectToHome(w, r, "error", "Invalid stock")
		return
	}

	// Save
	store.Create(models.Product{Name: name, Price: price, Stock: stock})

	// PRG: redirect so reloading the page does not repeat the POST
	redirectToHome(w, r, "ok", "Product created")
}
