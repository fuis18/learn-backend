package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"api-json/models"
)

// store holds the products in memory.
var store = models.NewStore()

// writeJSON writes a Go value as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// readJSON decodes the request body into dst.
// It rejects unknown fields and extra data after the object.
func readJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	if dec.More() {
		return errors.New("body must contain a single JSON object")
	}
	return nil
}

// parseID extracts the {id} path variable and converts it to an int.
func parseID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

//------------
// API Request
//------------

type healthResponse struct {
	Status string `json:"status"`
}

// Health responds OK to check that the server is alive.
func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := store.List()
	writeJSON(w, http.StatusOK, products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	product, ok := store.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := readJSON(r, &product); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body: " + err.Error()})
		return
	}

	created := store.Create(product)
	w.Header().Set("Location", "/api/products/"+strconv.Itoa(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

func ReplaceProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	var product models.Product
	if err := readJSON(r, &product); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body: " + err.Error()})
		return
	}

	replaced, ok := store.Replace(id, product)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}

	writeJSON(w, http.StatusOK, replaced)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	var patch models.ProductPatch
	if err := readJSON(r, &patch); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body: " + err.Error()})
		return
	}

	updated, ok := store.Update(id, patch)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	if !store.Delete(id) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
