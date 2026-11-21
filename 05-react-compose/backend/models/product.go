// Package models Product for the CRUD
package models

type Product struct {
	ID    int     `json:"id"`
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductPatch struct {
	Code  *string  `json:"code,omitempty"`
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
}
