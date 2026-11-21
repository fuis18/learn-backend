package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"podman-postgres/db"
	"podman-postgres/models"
)

// productColumns is reused by every SELECT/RETURNING to keep scans in sync.
const productColumns = "id, name, price, stock"

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

// internalError logs the real cause and answers with a generic 500,
// so no database details leak to the client.
func internalError(w http.ResponseWriter, op string, err error) {
	log.Printf("❌ %s: %v", op, err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

// validateProduct returns an error message when the payload is not valid.
func validateProduct(p models.Product) string {
	if strings.TrimSpace(p.Name) == "" {
		return "name is required"
	}
	if p.Price < 0 {
		return "price must be >= 0"
	}
	if p.Stock < 0 {
		return "stock must be >= 0"
	}
	return ""
}

//------------
// API Request
//------------

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
		log.Printf("⚠️ health check failed (postgres): %v", err)
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

// ListProducts returns every product ordered by id.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(r.Context(),
		`SELECT `+productColumns+` FROM products ORDER BY id`)
	if err != nil {
		internalError(w, "list products", err)
		return
	}
	defer rows.Close()

	products := make([]models.Product, 0) // [] instead of null in JSON
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			internalError(w, "scan product", err)
			return
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		internalError(w, "iterate products", err)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// GetProduct returns one product by id.
func GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	var p models.Product
	err = db.Pool.QueryRow(r.Context(),
		`SELECT `+productColumns+` FROM products WHERE id = $1`, id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		internalError(w, "get product", err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// CreateProduct inserts a product and answers 201 with the persisted row
// (numeric rounding included) and a Location header.
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := readJSON(r, &p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body: " + err.Error()})
		return
	}
	if msg := validateProduct(p); msg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
		return
	}

	err := db.Pool.QueryRow(r.Context(),
		`INSERT INTO products (name, price, stock) VALUES ($1, $2, $3)
		 RETURNING `+productColumns,
		p.Name, p.Price, p.Stock).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		internalError(w, "create product", err)
		return
	}

	w.Header().Set("Location", "/api/products/"+strconv.Itoa(p.ID))
	writeJSON(w, http.StatusCreated, p)
}

// ReplaceProduct replaces every field of a product (PUT).
func ReplaceProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	var p models.Product
	if err := readJSON(r, &p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body: " + err.Error()})
		return
	}
	if msg := validateProduct(p); msg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
		return
	}

	err = db.Pool.QueryRow(r.Context(),
		`UPDATE products SET name = $2, price = $3, stock = $4
		 WHERE id = $1 RETURNING `+productColumns,
		id, p.Name, p.Price, p.Stock).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		internalError(w, "replace product", err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// UpdateProduct patches only the fields present in the body (PATCH):
// it builds the SET clause dynamically with $1..$n placeholders.
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

	set := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if v := patch.Name; v != nil {
		if strings.TrimSpace(*v) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
			return
		}
		args = append(args, *v)
		set = append(set, fmt.Sprintf("name = $%d", len(args)))
	}
	if v := patch.Price; v != nil {
		if *v < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "price must be >= 0"})
			return
		}
		args = append(args, *v)
		set = append(set, fmt.Sprintf("price = $%d", len(args)))
	}
	if v := patch.Stock; v != nil {
		if *v < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "stock must be >= 0"})
			return
		}
		args = append(args, *v)
		set = append(set, fmt.Sprintf("stock = $%d", len(args)))
	}
	if len(set) == 0 {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "no fields to update (send name, price and/or stock)"})
		return
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE products SET %s WHERE id = $%d RETURNING %s`,
		strings.Join(set, ", "), len(args), productColumns)

	var p models.Product
	err = db.Pool.QueryRow(r.Context(), query, args...).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		internalError(w, "update product", err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// DeleteProduct removes a product and answers 204 on success.
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
		return
	}

	tag, err := db.Pool.Exec(r.Context(), `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		internalError(w, "delete product", err)
		return
	}
	if tag.RowsAffected() == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
