// Package models Product for the CRUD
package models

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type ProductPatch struct {
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
	Stock *int     `json:"stock,omitempty"`
}
