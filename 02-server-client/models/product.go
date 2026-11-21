// Package models contains the entities and the in-memory store
package models

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}
