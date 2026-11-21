package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"react-compose/db"
	"react-compose/logger"
	"react-compose/models"
)

// productColumns is reused by every SELECT/RETURNING to keep scans in sync.
const productColumns = "id, code, name, price"

// digitsOnly matches a barcode value (numeric content, leading zeros allowed).
var digitsOnly = regexp.MustCompile(`^[0-9]+$`)

// keepCodeSequenceAhead moves the fallback sequence past a manually entered
// code (best effort) so the next generated code cannot collide with it.
func keepCodeSequenceAhead(ctx context.Context, code string) {
	n, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return // too big for int64: the sequence can never reach it
	}
	_, _ = db.Pool.Exec(ctx,
		`SELECT setval('products_code_seq', GREATEST((SELECT last_value FROM products_code_seq), $1))`,
		n)
}

// validateProduct returns an error message when the payload is not valid.
// code is optional: when empty the database generates it (sequence fallback).
func validateProduct(p models.Product) string {
	if p.Code != "" && !digitsOnly.MatchString(p.Code) {
		return "code must contain only digits"
	}
	if strings.TrimSpace(p.Name) == "" {
		return "name is required"
	}
	if p.Price < 0 {
		return "price must be >= 0"
	}
	return ""
}

// ListProducts returns every product ordered by id.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(r.Context(),
		`SELECT `+productColumns+` FROM products ORDER BY id`)
	if err != nil {
		internalError(w, r, "list products", err)
		return
	}
	defer rows.Close()

	products := make([]models.Product, 0) // [] instead of null in JSON
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Price); err != nil {
			internalError(w, r, "scan product", err)
			return
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		internalError(w, r, "iterate products", err)
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
		Scan(&p.ID, &p.Code, &p.Name, &p.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		internalError(w, r, "get product", err)
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

	// code is optional: when empty let the sequence fallback generate it.
	explicitCode := p.Code != ""
	insert := `INSERT INTO products (name, price) VALUES ($1, $2)
		 RETURNING ` + productColumns
	args := []any{p.Name, p.Price}
	if explicitCode {
		insert = `INSERT INTO products (code, name, price) VALUES ($1, $2, $3)
		 RETURNING ` + productColumns
		args = []any{p.Code, p.Name, p.Price}
	}

	err := db.Pool.QueryRow(r.Context(), insert, args...).
		Scan(&p.ID, &p.Code, &p.Name, &p.Price)
	if isUniqueViolation(err) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "code already exists"})
		return
	}
	if err != nil {
		internalError(w, r, "create product", err)
		return
	}
	if explicitCode {
		keepCodeSequenceAhead(r.Context(), p.Code)
	}

	logger.From(r.Context()).Info("product created",
		slog.Int("id", p.ID), slog.String("code", p.Code))

	w.Header().Set("Location", "/api/products/"+strconv.Itoa(p.ID))
	writeJSON(w, http.StatusCreated, p)
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
	if v := patch.Code; v != nil {
		if !digitsOnly.MatchString(*v) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "code must contain only digits"})
			return
		}
		args = append(args, *v)
		set = append(set, fmt.Sprintf("code = $%d", len(args)))
	}
	if len(set) == 0 {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "no fields to update (send code, name and/or price)"})
		return
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE products SET %s WHERE id = $%d RETURNING %s`,
		strings.Join(set, ", "), len(args), productColumns)

	var p models.Product
	err = db.Pool.QueryRow(r.Context(), query, args...).
		Scan(&p.ID, &p.Code, &p.Name, &p.Price)
	if isUniqueViolation(err) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "code already exists"})
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		internalError(w, r, "update product", err)
		return
	}
	if patch.Code != nil {
		keepCodeSequenceAhead(r.Context(), *patch.Code)
	}

	logger.From(r.Context()).Info("product updated",
		slog.Int("id", p.ID), slog.String("code", p.Code))

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
		internalError(w, r, "delete product", err)
		return
	}
	if tag.RowsAffected() == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}

	logger.From(r.Context()).Info("product deleted", slog.Int("id", id))

	w.WriteHeader(http.StatusNoContent)
}
